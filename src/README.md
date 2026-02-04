# Practical Exploitation Work

> **Prerequisite**: Read the [main README](../README.md) first to understand eBPF fundamentals, the verifier, and the inherited vulnerability assessment from the previous team.
>
> **Context**: This directory contains the actual exploitation work — PoC development, testing, and analysis.
>
> **Documentation Hierarchy**: 
> - **Parent**: [Root README](../README.md) — eBPF background and project overview
> - **Current**: This file — Exploitation workflow and environment setup
> - **Next**: [src/exploits/README.md](exploits/README.md) — Multi-stage analysis methodology and vulnerability details

> **Note:** See [src/MODs.md](MODs.md) for a complete list of modifications made relative to the previous team's repository to ensure a consistent and reproducible setup.

---

## 1. Scope & Objective

This section documents the transition from theoretical vulnerability assessment to **practical exploitation**. While the previous team identified potential issues based on ISO-IEC TS 17961-2013, this work engineers functional **Proofs of Concept (PoCs)** to demonstrate real kernel impact on Linux LTS 6.8.

The goal is to answer: *Does this verifier bypass actually allow memory corruption or privilege escalation in a real runtime environment?*

---

## 2. Infrastructure Setup

A `Makefile` in the repository root simplifies the lifecycle management of the test VM.

### VM Management

- **Create & Provision:**
    
    ```bash
    make create
    
    ```
    
    *Creates a KVM/QEMU virtual machine with all required eBPF dependencies (clang, bpftool, etc.) pre-installed.*
    
- **Connect to the VM:**
    
    ```bash
    make connect
    
    ```
    
- **Destroy/Reset:**
    
    ```bash
    make destroy
    
    ```
    

---

## 3. End-to-End Exploitation Workflow

Follow this standard procedure for testing any of the assigned vulnerabilities.

### Step 1: Inject the Vulnerability (VM)

Inside the VM, navigate to the target directory and apply the specific patch you want to test.

**Important:** Run `git apply` from the root of the other team’s codebase (e.g., `/mnt/shared/`). The patch path must be **relative to that root**; otherwise `git apply` will fail to locate files.

```bash
# Example: Apply the 5.46b vulnerability patch
git apply XDPs/xdp_synproxy/patches/5_46b_taintsink/0001-feat-5.46-taintsink_2-Tainted-potentially-mutilated-.patch
```

### Step 2: Compile the Code (VM)

Compile the modified XDP program. This step converts the C code into eBPF bytecode.

```bash
cd /mnt/shared/XDPs/xdp_synproxy/

make clean && make
```

> Note: This only compiles the code. It does NOT load it into the kernel yet.
> 

### Step 3: Start the Test Session (VM)

Run the session script to load the eBPF program and attach it to the network interface.

```bash
./start_session.sh
```

This launches a **tmux session** with three panes:

1. **Loader Output:** Shows the status of the userspace XDP loader.
2. **Netcat Server:** Listens on Port 80 to receive traffic.
3. **Kernel Trace Pipe:** Displays real-time `bpf_printk` debug logs from the kernel.

For tmux navigation and session controls, see the **Tips & Troubleshooting** section in [Part 5](#5-tips--troubleshooting).

### Step 4: Trigger the Exploit (Host)

Switch to your **Host machine**. You can trigger the vulnerability using standard tools (like `netcat`) or specific exploit scripts.

- **Basic Trigger (Netcat):**

    `nc <VM_IP> 80 -v`

- **Advanced Exploitation (Python/Scapy or Raw Sockets):**

    Specific exploit scripts are located in the src/exploits/ directory. These scripts craft malicious packets (e.g., with manipulated checksums or headers) to bypass initial checks and trigger specific memory corruption paths.

    `# Example
    cd src/exploits/5_46b_taintsink/
    sudo python3 poc_546b_final.py`
    

### Step 5: Analyze Results

Observe the **Kernel Trace Pipe (Pane 3)** in the VM.

- **Success:** Look for logs indicating the success of the exploit (e.g., `index: 60` or `Content: 42`).
- **Failure:** If no logs appear, check if the packet was dropped by the XDP parser (checksum errors) or if the VM interface is down.

**Additional analysis steps (as used in the reports):**

1. **Bytecode dump (`llvm-objdump`)** — Confirms the vulnerable path exists in the compiled `.o` and wasn’t optimized away. This validates the issue at the compiler output stage.
2. **Verifier log (`bpftool prog load -d`)** — Captures loader/verifier output to confirm the program is accepted and to gather reasoning/range data when available.
3. **Translated bytecode (`bpftool prog dump xlated`)** — Shows the verifier‑processed program, proving whether the unsafe path survived verification.
4. **JIT dump (`bpftool prog dump jited`)** — If available, confirms the **native** instructions still implement the unsafe path after JIT translation.

These steps establish a full evidence chain: **C source → .o bytecode → verifier view → xlated → JIT**, making it clear where the issue is preserved or mitigated.

---

## 4. Assigned Vulnerabilities & Analysis

From the original dataset of 60+ patches, the following 5 vulnerabilities were selected for in-depth practical exploitation. These were chosen because they represent critical failures in the eBPF verifier's ability to track types, bounds, or memory safety contexts.

**Detailed exploitation reports, Proof-of-Concept (PoC) scripts, and evidence logs are located in the specific subdirectory for each vulnerability.** Each vulnerability folder includes its own comprehensive README to explain the issue, evidence, and reproduction steps in depth (see the per‑vulnerability folders under src/exploits/).

| **ID** | **Vulnerability Name** | **ISO Rule** | **Description of Flaw** | **Analysis Location** |
| --- | --- | --- | --- | --- |
| **1** | **5.46b taintsink** | **5.46** | **Tainted Length Memory Copy** Trusts a user-controlled TCP header length to perform a `memcpy` into a fixed-size stack buffer. Causes a stack buffer overflow. | src/exploits/5_46b_taintsink/ |
| **2** | **5.39 taintnoproto** | **5.39** | **Unprototyped Function Pointer** Passes a tainted value (TCP Seq Num) to a function pointer declared without arguments. Bypasses type checking to perform an OOB write. | src/exploits/5_39_taintnoproto/ |
| **3** | **5.06a argcomp** | **5.06** | **Function Type Mismatch** Calls a function via a pointer with an incompatible signature (expects different types). Causes register/ABI confusion. | src/exploits/5_06a_argcomp/
| **4** | **5.06b argcomp** | **5.06** | **Wrong Argument Count** Calls a function with more arguments than it expects. Exploits how the BPF JIT handles extra registers. | src/exploits/5_06b_argcomp/ |
| **5** | **5.06d argcomp** | **5.06** | **Wrong Argument Types** Passes an `int` to a function expecting a `long`. Investigates potential truncation or sign-extension vulnerabilities. | src/exploits/5_06d_argcomp/ |

---

## 5. Tips & Troubleshooting

- **Session Management:** If logs stop appearing or the behavior seems inconsistent, the XDP program might be "stuck" from a previous run.
    - **Fix:** Kill the tmux session and force detach the program:Bash
        
        `sudo ip link set dev enp2s0 xdp off`
        
- **Tmux Navigation:**
    - `Ctrl+b` then `Arrow Keys` to switch panes.
    - `Ctrl+b` then `c` to create a new tmux window.
    - `Ctrl+b` then `n`/`p` to move to next/previous window, or `Ctrl+b` then `w` to choose a window.
    - `Ctrl+b` then `z` to zoom into the trace pipe for better visibility.
- **Reverting Changes:** Always revert a patch before applying a new one to ensure a clean testing state.Bash
    
    `git apply -R patches/<previous_patch>.patch`
    
---

## Summary

1. **Provision** the VM: `make create`. This opens the VM terminal.
2. **Patch** the target code in `/mnt/shared/XDPs/xdp_synproxy/`.
3. **Compile** with `make`.
4. **Load** with `./start_session.sh`.
5. **Exploit** from the host using scripts in `src/exploits/`.
6. **Verify** using the kernel trace output.
7. **Analyze** with bytecode, verifier, xlated, and (if supported) JIT dumps to complete the evidence chain.

This workflow ensures a reproducible, step-by-step process for validating eBPF verifier failures.

---

## 📖 Next Steps

Now that you understand the environment setup and workflow, proceed to the detailed analysis:

**→ [Continue to src/exploits/README.md](exploits/README.md)** to learn about the multi-stage analysis methodology and detailed descriptions of each of the 5 vulnerabilities under investigation.