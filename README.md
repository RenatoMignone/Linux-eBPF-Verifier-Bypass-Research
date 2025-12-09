# Practical Exploitation of Linux eBPF Verifier Vulnerabilities (LTS 6.8)

![Polito Logo](resources/images/logo_polito.jpg)

> **Course**: Security Verification and Testing - Second Year, First Semester  
> **Author**: Renato Mignone  
> **Date**: December 2025

---

## Abstract

This project explores the security boundaries of **eBPF (extended Berkeley Packet Filter)**, a revolutionary technology that allows user-defined programs to run safely inside the Linux kernel. While eBPF's verifier is designed to guarantee memory safety and prevent malicious code execution, certain logic bugs can potentially bypass these protections.

The work is structured in three phases:
1. **Understanding eBPF** — A comprehensive introduction to eBPF architecture, its verifier, and why it matters for kernel security.
2. **Prior Theoretical Assessment** — Analysis of inherited work from a previous team who identified 60+ potential vulnerabilities based on ISO-IEC TS 17961-2013, classifying ~10 as theoretically exploitable.
3. **Practical Exploitation** — Engineering functional Proofs of Concept (PoCs) that demonstrate real kernel impact on LTS 6.8.

---

## Table of Contents

1. [What is eBPF?](#what-is-ebpf)
   - [Origins and Evolution](#origins-and-evolution)
   - [Architecture Overview](#architecture-overview)
   - [The eBPF Verifier](#the-ebpf-verifier)
   - [Program Types and Use Cases](#program-types-and-use-cases)
   - [Security Implications](#security-implications)
2. [Inherited Work: What the Previous Team Did](#inherited-work-what-the-previous-team-did)
3. [My Objective: Practical Exploitation](#my-objective)
4. [Repository Structure](#repository-structure)
5. [Quick Start](#quick-start-using-inherited-infrastructure)
6. [Progress Tracking](#progress-tracking)
7. [References](#references)

---

## What is eBPF?

### Origins and Evolution

**eBPF (extended Berkeley Packet Filter)** is a revolutionary technology that originated from the classic BPF (Berkeley Packet Filter) created in 1992 for network packet filtering. While classic BPF was limited to simple packet inspection (used by tools like `tcpdump`), eBPF has evolved into a general-purpose **in-kernel virtual machine** capable of running sandboxed programs in the Linux kernel.

<p align="center">
  <img src="resources/images/ebpf_timeline.jpg" width="700">
</p>

**Timeline:**
* **1992**: BPF created (packet filtering only).
* **2014**: eBPF introduced in Linux 3.18 (Extended instruction set, Maps, Helpers).
* **2016**: XDP (eXpress Data Path) added for programmable network processing at the driver level.
* **2020+**: eBPF becomes ubiquitous for Observability, Security, and Networking.

### Architecture Overview

eBPF programs are written in a restricted C dialect, compiled to eBPF bytecode, verified by the kernel, and JIT-compiled to native machine code for execution.

<p align="center">
  <img src="resources/images/ebpf_architecture.jpg" width="600">
</p>

**Key Components:**

| Component | Description |
|-----------|-------------|
| **eBPF Bytecode** | Platform-independent instruction set (11 registers, 64-bit) |
| **Verifier** | Static analyzer ensuring program safety before execution |
| **JIT Compiler** | Translates bytecode to native CPU instructions |
| **Maps** | Key-value stores for sharing data between eBPF programs and user space |
| **Helpers** | Kernel functions callable from eBPF (e.g., `bpf_map_lookup_elem`) |
| **Hooks** | Attachment points in kernel (kprobes, tracepoints, XDP, TC, LSM) |

### The eBPF Verifier

The verifier is the **security cornerstone** of eBPF. It performs static analysis on every program before allowing it to run in the kernel. Its goal is to guarantee:

<p align="center">
  <img src="resources/images/verifier_logic.jpg" width="300">
</p>

1. **Memory Safety** — No out-of-bounds reads/writes
2. **Termination** — Program will always complete (no infinite loops)
3. **Valid Operations** — Only allowed helpers and operations
4. **Type Safety** — Correct use of pointers and data types

**Verifier Analysis Stages:**
1.  **CFG Construction**: Build Control Flow Graph and identify execution paths.
2.  **Instruction Walk**: Simulate each instruction and track register states.
3.  **Bounds Checking**: Verify memory accesses and track pointer arithmetic.
4.  **Loop Detection**: Ensure bounded iteration.
5.  **Helper Validation**: Check argument types.

**Verifier Limitations (The Attack Surface):**

The verifier is complex (~20,000 lines of code in `kernel/bpf/verifier.c`) and must balance security, usability, and performance. This complexity creates potential for **logic bugs** where:
- The verifier's model of program behavior differs from actual execution.
- Edge cases in type tracking allow invalid operations.
- Pointer arithmetic confuses the verifier's bounds tracking.

**These verifier bugs are the focus of this project.**

### How the Verifier Tracks State (The Logic)

To ensure safety without executing the code, the verifier simulates every instruction and maintains a **state** for each of the 11 registers (`R0`-`R10`). A vulnerability occurs when this internal state gets out of sync with the actual runtime machine state.

**What the Verifier tracks per register:**
1.  **Type**: Is it a **SCALAR** (number) or a **POINTER** (to memory/map)?
    * *Exploit Path:* If we can trick the verifier into thinking a **SCALAR** is a **POINTER**, we can read arbitrary kernel memory.
2.  **Range (`tnum`)**: What is the min/max value this register holds?
    * *Exploit Path:* If the verifier thinks a value is `0-100` but at runtime it is `1000`, bounds checks are bypassed (OOB R/W).
3.  **Offset**: Where does this pointer point relative to the base object?

> **The "Pruning" Problem:** To save time, if the verifier sees a code path that looks similar to one it already checked, it skips it ("pruning"). Many logic bugs hide here—where the verifier prunes a path it *thought* was safe, but actually wasn't.

### Program Types and Use Cases

eBPF supports multiple program types, each attaching to different kernel hooks:

| Program Type | Hook Point | Use Case |
|--------------|------------|----------|
| `BPF_PROG_TYPE_XDP` | Network driver (pre-stack) | DDoS mitigation, load balancing |
| `BPF_PROG_TYPE_SCHED_CLS` | Traffic Control (TC) | Packet filtering, NAT |
| `BPF_PROG_TYPE_KPROBE` | Kernel functions | Tracing, debugging |
| `BPF_PROG_TYPE_TRACEPOINT` | Static kernel events | Performance monitoring |
| `BPF_PROG_TYPE_LSM` | Linux Security Modules | Security policies |
| `BPF_PROG_TYPE_SOCKET_FILTER` | Socket layer | Packet inspection |

**XDP (eXpress Data Path)** is particularly interesting for security research because it runs at the **earliest point** in the network stack, has **direct memory access**, and is **performance-critical**.

### Security Implications

eBPF's ability to run custom code in kernel space makes it a **high-value target** for attackers.

<p align="center">
  <img src="resources/images/security_sandbox.jpg" width="600">
</p>

* **Intended Boundaries**: Cannot access arbitrary memory, cannot call arbitrary functions.
* **If Verifier Bypassed**: Arbitrary kernel read (Info disclosure), Arbitrary kernel write (Privilege escalation).

**Historical eBPF Vulnerabilities:**

| CVE | Year | Description |
|-----|------|-------------|
| CVE-2020-8835 | 2020 | Verifier bounds tracking error, arbitrary R/W |
| CVE-2021-3490 | 2021 | ALU32 bounds tracking, privilege escalation |
| CVE-2021-31440 | 2021 | Bounds propagation bug |
| CVE-2022-23222 | 2022 | Pointer arithmetic confusion |
| CVE-2023-2163 | 2023 | Verifier log buffer overflow |

This project focuses on identifying and exploiting similar **verifier logic bugs** in LTS 6.8.

---

## Inherited Work: What the Previous Team Did

The repository `ebpf-exploitability-test/` contains work from a previous team (Francesco Rollo, Gianfranco Trad, Giorgio Fardo, Giovanni Nicosia) who performed a **theoretical vulnerability assessment** based on **ISO-IEC TS 17961-2013** (C Secure Coding Standard).

### What is XDP SynProxy?

The previous team focused **exclusively** on testing vulnerabilities within the XDP SynProxy implementation. To understand why this is a significant target, we need to understand the problem it solves.

#### The SYN Flood Attack

<p align="center">
  <img src="resources/images/syn_flood_diagram.jpg" width="600">
</p>

A **SYN flood** is a type of Denial-of-Service (DoS) attack that exploits the TCP three-way handshake:

**The problem**: Each SYN packet causes the server to allocate a **Transmission Control Block (TCB)** and wait for the ACK. Attackers send thousands of SYN packets with spoofed source IPs, exhausting the server's connection table while never completing the handshake.

#### The SYN Proxy Solution

<p align="center">
  <img src="resources/images/syn_proxy_flow.jpg" width="600">
</p>

A **SYN Proxy** acts as an intermediary that validates clients before forwarding connections to the real server:

**SYN Cookies**: Instead of storing state, the proxy encodes connection information (source IP, port, timestamp) into the **sequence number** of the SYN-ACK. When the client responds with ACK, the proxy can validate the cookie mathematically without having stored any state. Fake clients with spoofed IPs never receive the SYN-ACK, so they can't complete the handshake.

#### Why XDP for SYN Proxy?

<p align="center">
  <img src="resources/images/xdp_stack_comparison.jpg" width="600">
</p>

Traditional SYN proxies operate in user space or higher in the kernel stack, which introduces latency. **XDP (eXpress Data Path)** solves this by running the SYN proxy logic at the **earliest possible point** — directly in the network driver:

**XDP SynProxy advantages**:
- **Line-rate processing**: Can handle millions of packets per second
- **Zero-copy**: Operates directly on packet memory in the driver
- **Early drop**: Malicious packets never enter the kernel stack
- **CPU efficient**: No context switches, minimal overhead

#### The Target: `xdp_synproxy_kern.c`

The file `xdp_synproxy_kern.c` is an eBPF/XDP implementation of a SYN proxy, taken from the **Linux kernel selftests**. It performs:
1. **Packet parsing**: Ethernet → IP → TCP header dissection
2. **SYN detection**: Identifies incoming TCP SYN packets
3. **Cookie generation**: Creates cryptographic SYN cookies
4. **SYN-ACK crafting**: Builds response packets with embedded cookies
5. **ACK validation**: Verifies returning ACKs contain valid cookies
6. **Connection forwarding**: Passes validated connections to the kernel

The previous team chose this target because:
- It's **real production code** from Linux kernel selftests
- It involves **complex pointer arithmetic** for packet parsing
- It uses **multiple eBPF features** (maps, helpers, tail calls)
- It's **security-critical** — bugs here could disable DDoS protection or leak information
- It's **performance-critical** — must process packets at wire speed

### Their Testing Environment

The team created an isolated VM-based testing environment using KVM/QEMU and Ubuntu 24.04.

<p align="center">
  <img src="resources/images/VM_Structure.jpg" width="700">
</p>

**Key design decisions:**
- **Single VM topology**: Simplified from a 3-veth setup to avoid checksum issues.
- **Cloud-init provisioning**: The VM is auto-configured via `user-data.yaml`.
- **tmux testing sessions**: `start_session.sh` creates a split terminal with the XDP program, netcat server, and packet monitoring.

### Their Methodology

<p align="center">
  <img src="resources/images/old_work_structure.jpg" width="700">
</p>

They applied 46 C vulnerability rules to the target code to test the verifier's response.

### What They Produced

| Artifact | Description |
|----------|-------------|
| **60+ vulnerability patches** | Each implements a specific ISO-IEC TS 17961-2013 rule violation |
| **Base target: `xdp_synproxy_kern.c`** | Real XDP SYN proxy from Linux kernel selftests |
| **xvtlas tool** | Go automation suite for patch/compile/verify workflow |
| **patches.csv** | Summary table with exploitability classification |
| **Detailed documentation** | 1500+ lines explaining each vulnerability |

### Their Results Summary

* **60 patches total**
    * ~10 compilation failures (rejected before verifier).
    * ~15 blocked by eBPF verifier ✓ (security works).
    * ~25 passed verifier but "not exploitable" (memory bounds still enforced).
    * **~10 passed verifier AND marked "exploitable"** (These are my targets).

### Vulnerabilities Marked as "Exploitable"

| Patch | ISO Rule | Vulnerability Type | Their Assessment |
|-------|----------|-------------------|------------------|
| `5_06a_argcomp` | 5.6 | Function pointer mismatch | Register/stack corruption |
| `5_06b_argcomp` | 5.6 | Wrong argument count | Stack memory overwrite |
| `5_10a_exploit_intptrconv` | 5.10 | Pointer truncation bypass | **Info disclosure** |
| `5_14_nullref` | 5.14 | Null pointer dereference | Invalid memory access |
| `5_17_swtchdflt` | 5.17 | Missing switch default | Undefined firewall behavior |
| `5_20a_libptr` | 5.20 | Buffer overflow (8 bytes) | Adjacent stack corruption |
| `5_20c_libptr` | 5.20 | Type confusion overflow | 12-byte buffer overflow |
| `5_35_uninit_mem` | 5.35 | Uninitialized memory read | Kernel stack data leak |
| `5_35a_unint_mem` | 5.35 | Uninitialized memory read | Kernel stack data leak |

### Vulnerabilities Marked as "Limited" Exploitability

These vulnerabilities passed the verifier but have uncertain or constrained exploitation potential:

| Patch | ISO Rule | Vulnerability Type | Their Assessment |
|-------|----------|-------------------|------------------|
| `5_06d_argcomp` | 5.6 | Wrong argument types | Value truncation (localized impact) |
| `5_14a_nullref` | 5.14 | Null pointer dereference | Invalid access, no escape path |
| `5_16b_signconv` | 5.16 | Signed conversion | Logic errors only |
| `5_33a_restrict` | 5.33 | Restrict pointer violation | Logic/data corruption in stack |
| `5_33b_restrict` | 5.33 | Restrict pointer violation | Local stack data corruption |
| `5_36a_ptrobj` | 5.36 | Pointer comparison UB | Memory layout information leak |
| `5_36b_ptrobj` | 5.36 | Context pointer comparison | Kernel layout info disclosure |
| `5_36c_ptrobj` | 5.36 | Map pointer comparison | Heap organization leak |
| `5_39_taintnoproto` | 5.39 | Tainted function pointer | Unpredictable logic behavior |
| `5_45_invfmtstr` | 5.45 | Invalid format strings | Address leak via logging |
| `5_46b_taintsink` | 5.46 | Tainted memory copy | Attacker-controlled packet alteration |

### Vulnerabilities Marked as "Not Exploitable"

These vulnerabilities were either blocked by the compiler, blocked by the verifier, or passed but had no security impact:

**Blocked by Compilation Errors (8 patches):**
| Patch | Reason |
|-------|--------|
| `5_06c_argcomp` | Conflicting types for function |
| `5_13_objdec` | Conflicting types for variable |
| `5_13b_objdec` | Conflicting types for variable |
| `5_22d_invptr` | Array subscript out of bounds |
| `5_24a_usrfmt` | Incompatible pointer types |

**Blocked by eBPF Verifier (12 patches):**
| Patch | Verifier Error |
|-------|----------------|
| `5_4a_boolasgn` | Infinite loop detected |
| `5_4b_boolasgn` | Infinite loop detected |
| `5_06e_argcomp` | R1 type=scalar expected=map_ptr |
| `5_10a_intptrconv` | Pointer arithmetic with <<= prohibited |
| `5_10b_intptrconv` | R1 invalid mem access 'scalar' |
| `5_14b_nullref` | R7 invalid mem access 'map_value_or_null' |
| `5_16a_signconv` | R8 offset outside packet |
| `5_20b_libptr` | Invalid indirect access to stack |
| `5_22c_invptr` | Invalid access to context parameter |
| `5_24b_usrfmt` | Invalid access to context parameter |
| `5_35b_unint_mem` | R1 invalid mem access 'scalar' |
| `5_40_taintformatio` | Invalid indirect access to stack |
| `5_46a_taintsink` | Unbounded min value not allowed |
| `5_46c_taintsink` | Address R11 invalid (VLA attempt) |

**Passed Verifier but Not Exploitable (19 patches):**
| Patch | Reason |
|-------|--------|
| `5_01a/b/c_ptrcomp` | Memory bounds still enforced by verifier |
| `5_9_padcomp` | Logic non-determinism only, no memory exposure |
| `5_11_alignconv` | Logical misinterpretation, not memory exploitable |
| `5_11a/c_alignconv` | Logical misinterpretation, not memory exploitable |
| `5_15_addrescape` | Dangling pointer but stack managed per-packet |
| `5_15a/b_addrescape` | Dangling pointer but stack managed per-packet |
| `5_22_invptr` | Verifier-enforced memory bounds |
| `5_22b/e/f_invptr` | Logical misinterpretation only |
| `5_26a-e_diverr` | Division by zero results in 0, no crash |
| `5_28_strmod` | Read-only memory, program terminates |
| `5_30_intoflow` | Two's complement arithmetic, no memory violation |
| `5_31a/b_nonnullcs` | Controlled memory layout required |

### Complete Vulnerability Summary

```
Total: 60 patches tested
         │
         ├─── YES (Exploitable): 9 patches ────────────► PRIMARY TARGETS
         │    └── Verifier bypassed, real security impact possible
         │
         ├─── LIMITED (Uncertain): 12 patches ─────────► SECONDARY TARGETS  
         │    └── Passed verifier, constrained exploitation potential
         │
         └─── NO (Not Exploitable): 39 patches
              ├── 8 blocked by compiler
              ├── 14 blocked by verifier ✓ (security works)
              └── 17 passed but no security impact
```

**Potential investigation scope: 21 patches (9 YES + 12 LIMITED)**

### What They Did NOT Do

- ❌ Create actual working exploits
- ❌ Achieve arbitrary kernel read/write
- ❌ Demonstrate privilege escalation
- ❌ Validate impact on real kernel versions

They proved vulnerabilities **can bypass the verifier**, but stopped at theoretical classification.

---

## Repository Structure

```text
Project_EBPF/
├── README.md                          # This file
├── report/                            # My exploitation reports (TODO)
├── src/                               # My PoC code (TODO)
│
└── ebpf-exploitability-test/          # INHERITED FROM PREVIOUS TEAM
    └── codebase/
        ├── README.md                  # Their documentation
        ├── docs/
        │   └── ISO-IEC-TS-17961-2013.pdf
        │
        ├── virt/                      # VM management
        │   ├── vmctl.sh               # Create/destroy/connect to test VM
        │   └── meta-data.yaml
        │
        ├── XDPs/
        │   ├── xdp_synproxy/          # Main target
        │   │   ├── xdp_synproxy_kern.c    # Base vulnerable program
        │   │   ├── Makefile
        │   │   ├── patches/           # 60+ vulnerability patches
        │   │   │   ├── 5_10a_exploit_intptrconv/
        │   │   │   ├── 5_14_nullref/
        │   │   │   ├── 5_20a_libptr/
        │   │   │   └── ...
        │   │   ├── patches.csv        # Summary of all patches
        │   │   ├── apply_rules.sh
        │   │   ├── start_session.sh
        │   │   └── kill_session.sh
        │   │
        │   ├── minimal/               # Minimal exploit examples
        │   │   └── 5_10a_exploit.c
        │   │
        │   └── tools/                 # BPF headers
        │
        ├── xvtlas/                    # Automation tool (Go)
        │   ├── main.go
        │   ├── xvtlas                 # Pre-compiled binary
        │   └── README.md
        │
        └── xvtlas_output/             # Previous test results
            ├── report.csv
            └── [per-patch logs]/
```

---

## My Objective

Transform the theoretical exploitability assessments into **functional Proof of Concept exploits** targeting the **LTS 6.8 kernel**.

### Goals

1. **Analyze** the ~10 "exploitable" vulnerabilities in depth
2. **Develop** working PoCs that demonstrate real kernel impact
3. **Document** exploitation techniques and verifier bypass methods
4. **Validate** severity on LTS 6.8 kernel

---

## Quick Start (Using Inherited Infrastructure)

### 1. Set Up Test VM

```bash
cd ebpf-exploitability-test/codebase/virt/
./vmctl.sh create ~/.ssh/id_rsa.pub
./vmctl.sh connect
```

### 2. Test a Vulnerability Patch

```bash
# Inside VM
cd ~/ebpf-tests/XDPs/xdp_synproxy/

# Apply a patch
git apply patches/5_10a_exploit_intptrconv/0001-feat-5_10a_exploit.patch

# Compile
make

# Test
./start_session.sh
```

### 3. Use xvtlas Automation

```bash
cd xvtlas/
./xvtlas --run-single "./patches/5_10a_exploit_intptrconv/*.patch" \
         --base-file "./xdp_synproxy_kern.c"
```

---

## Progress Tracking

- [ ] Environment setup (VM, kernel 6.8)
- [ ] Study inherited codebase
- [ ] Analyze priority vulnerabilities
- [ ] Develop PoC #1: (TBD)
- [ ] Develop PoC #2: (TBD)
- [ ] Write exploitation report
- [ ] Final documentation

---

## References

- [ISO-IEC TS 17961-2013](https://www.iso.org/standard/61134.html) - C Secure Coding Standard
- [eBPF Documentation](https://ebpf.io/)
- [Linux Kernel BPF Verifier](https://www.kernel.org/doc/html/latest/bpf/verifier.html)
- Previous team's detailed documentation: `ebpf-exploitability-test/codebase/XDPs/xdp_synproxy/README.md`