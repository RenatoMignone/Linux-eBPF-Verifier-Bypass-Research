# Practical Exploitation Work

> **Context**: This directory contains the actual exploitation work — PoC development, testing, and analysis.  
> **Prerequisite**: Read the [main README](../README.md) first to understand eBPF, the verifier, and the inherited vulnerability assessment.

---

## What is This Practical Work About?

This section documents the full, step-by-step process for moving from a clean VM environment to a completed exploitation workflow for selected eBPF verifier vulnerabilities. The goal is to demonstrate, with working Proofs of Concept (PoCs), how specific vulnerabilities can be exploited in practice on Linux LTS 6.8 using the inherited infrastructure.

---

## End-to-End Exploitation Workflow

### 1. Set Up and Start the Test VM

- Go to the VM management scripts:
  ```bash
  cd ../ebpf-exploitability-test/codebase/virt/
  ./vmctl.sh create ~/.ssh/id_rsa.pub
  ./vmctl.sh connect
  ```
- This will create and provision a VM with all required dependencies (clang, bpftool, etc.).

### 2. Enter the XDP SynProxy Directory in the VM

- Once inside the VM, navigate to the shared folder:
  ```bash
  cd /mnt/shared/XDPs/xdp_synproxy/
  ```

### 3. Apply a Vulnerability Patch (if testing a specific bug)

- To test a vulnerability, apply the relevant patch:
  ```bash
  git apply patches/<patch-folder>/<patch-file>.patch
  # Example:
  git apply patches/5_06a_argcomp/0001-feat-5_06a_argcomp.patch
  ```

### 4. Compile the Code

- Build the kernel and user-space programs:
  ```bash
  make clean && make
  ```
- The Makefile only compiles the code. It does **not** load the eBPF program into the kernel.

### 5. Start the Test Session and Load the eBPF Program

- Use the provided script to load the eBPF program and set up the test environment:
  ```bash
  ./start_session.sh
  ```
- This script:
  - Extracts `vmlinux.h` if needed
  - Loads the eBPF program into the kernel and attaches it to the network interface
  - Sets up a tmux session with:
    - Pane 1: XDP SynProxy loader output
    - Pane 2: Netcat server listening on port 80
    - Pane 3: Kernel trace output

### 6. Trigger the Vulnerability (from Host)

- On your host machine, connect to the VM using netcat to send crafted packets:
  ```bash
  nc <VM_IP> 80 -v
  ```
- This simulates an attack and allows you to observe the eBPF program's behavior and test for exploitation.

### 7. Analyze and Document Results

- Observe the tmux panes for output, kernel messages, and exploitation evidence.
- Use tools like `bpftool`, `llvm-objdump`, and logs to analyze verifier output and bytecode.
- Document your findings, PoC code, and exploitation steps in the appropriate subdirectory for each vulnerability.

---

## Applying and Reverting Vulnerability Patches

To test each vulnerability, apply the corresponding patch using the following commands (run these in `/mnt/shared/XDPs/xdp_synproxy/` inside the VM):

```bash
# 5.6 argcomp - Function pointer type incompatibility
git apply patches/5_06a_argcomp/0001-feat-5_06a.patch

# 5.6 argcomp - Wrong number of arguments
git apply patches/5_06b_argcomp/0001-feat-5_06b.patch

# 5.6 argcomp - Wrong argument types
git apply patches/5_06d_argcomp/0001-feat-5_06d.patch

# 5.39 taintnoproto - Using tainted values as function pointers without prototypes
git apply patches/5_39_taintnoproto/0001-feat-5.39-taintnoproto-Using-a-tainted-value-as-poin.patch

# 5.46b taintsink - Memory copy with tainted length
git apply patches/5_46b_taintsink/0001-feat-5.46-taintsink_2-Tainted-potentially-mutilated-.patch
```

### Rolling Back (Reverting) a Patch

To revert a patch and return to the clean version of the code, use:

```bash
# Example: Revert the 5.6a argcomp patch
git apply -R patches/5_06a_argcomp/0001-feat-5_06a_argcomp.patch
```

You can use the same `-R` option for any patch to roll back to the previous (unpatched) state.

---

## Assigned Vulnerabilities (Focus)

You are focusing on the following vulnerabilities for practical exploitation:

- **5.6 argcomp** — Function pointer type incompatibility
- **5.6 argcomp** — Wrong number of arguments
- **5.6 argcomp** — Wrong argument types
- **5.39 taintnoproto** — Using tainted values as function pointers without prototypes
- **5.46b taintsink** — Memory copy with tainted length

---

## Directory Structure and Workflow

- Each vulnerability should have its own folder (e.g., `5_06a_argcomp/`) containing:
  - Analysis notes and logs
  - PoC code and Makefile
  - Evidence (screenshots, logs)
  - README documenting the exploitation process

- Use the provided templates and tools for consistent documentation and analysis.

---

## Summary: From Start to Finish

1. **Provision and connect to the VM**
2. **Navigate to the XDP SynProxy directory**
3. **Apply a patch for the vulnerability you are testing**
4. **Compile the code with `make`**
5. **Run `./start_session.sh` to load and activate the eBPF program**
6. **Trigger the vulnerability from the host using netcat or crafted packets**
7. **Analyze, document, and report your findings**

This workflow ensures a reproducible, step-by-step process for practical eBPF exploitation research.
