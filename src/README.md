# Practical Exploitation Work

> **Context**: This directory contains the actual exploitation work — PoC development, testing, and analysis.  
> **Prerequisite**: Read the [main README](../README.md) first to understand eBPF, the verifier, and the inherited vulnerability assessment.

---

## Step-by-Step Progress Log

This section documents the actual steps taken to set up the environment and begin exploitation work.

> **Note**: Modifications made to the inherited codebase are documented in [MODs.md](./MODs.md).

---

### Step 1: Setting Up the Test VM

The previous team provided a VM management script (`vmctl.sh`) to create an isolated testing environment. However, several issues needed to be resolved before it worked.

#### 1.1 Missing Files and Script Fixes

The inherited `vmctl.sh` script had issues:
1. **Missing `user-data.yaml`**: The cloud-init configuration file was not in the repository
2. **SSH key substitution bug**: The `sed` command failed with special characters in SSH keys

These were fixed — see [MODs.md](./MODs.md) for details.

#### 1.2 Installing Host Dependencies

Running the script initially failed due to missing virtualization tools:

```bash
./vmctl.sh create ~/.ssh/SVT_Project/public.pub
# Error: sudo: virt-install: command not found
```

**Solution**: Install the required virtualization packages on the host system:

```bash
# Update package list
sudo apt update

# Install virtualization tools
# Note: On Ubuntu 24.04, the package is 'virtinst' not 'virt-install'
sudo apt install -y virtinst libvirt-daemon-system libvirt-clients qemu-system qemu-utils
```

**Packages installed:**
| Package | Purpose |
|---------|---------|
| `virtinst` | Provides `virt-install` command for VM creation |
| `libvirt-daemon-system` | Libvirt daemon for managing VMs |
| `libvirt-clients` | CLI tools (`virsh`, etc.) |
| `qemu-system` | QEMU emulator |
| `qemu-utils` | QEMU disk image utilities |

#### 1.3 Creating the VM

After installing dependencies and applying fixes:

```bash
cd ebpf-exploitability-test/codebase/virt/
./vmctl.sh create ~/.ssh/SVT_Project/public.pub
```

The script will:
1. Download Ubuntu 24.04 cloud image (~597 MB)
2. Substitute your SSH public key into `user-data.yaml`
3. Create a VM with 4GB RAM, 2 vCPUs, 20GB disk
4. Boot and provision with cloud-init (installs clang, bpftool, etc.)

#### 1.4 Connecting to the VM

```bash
./vmctl.sh connect
```

Or via SSH (after getting the VM's IP):
```bash
# Get VM IP
virsh --connect qemu:///system domifaddr ubuntu-xdp-24.04

# SSH in
ssh ubuntu@<VM_IP>
```

---

### Step 2: Verifying the Environment (TODO)

Once inside the VM, verify the development tools:

```bash
# Check kernel version (should be 6.8.x for our target)
uname -r

# Check clang/LLVM
clang --version
llvm-objdump --version

# Check bpftool
bpftool version

# Check libbpf
pkg-config --modversion libbpf
```

---

### Step 3: Testing Baseline Compilation (TODO)

```bash
# Clone or copy the test files to the VM
cd ~/ebpf-tests/XDPs/xdp_synproxy/

# Compile clean version (no patches)
make clean && make

# Test loading
sudo bpftool prog load xdp_synproxy_kern.o /sys/fs/bpf/test
sudo bpftool prog list
sudo rm /sys/fs/bpf/test
```

---

### Step 4: Testing Patch Workflow (TODO)

*To be documented as we progress...*

---

## Where We Are Now

The previous team completed the **theoretical assessment phase**:

The previous team completed the **theoretical assessment phase**:

```
✅ Identified 60+ potential vulnerabilities (ISO-IEC TS 17961-2013)
✅ Created patches that inject each vulnerability into xdp_synproxy_kern.c
✅ Tested each patch against the eBPF verifier
✅ Classified results: YES (9) / LIMITED (12) / NO (39)
✅ Built automation tooling (xvtlas)
✅ Documented everything
```

**Our task**: Transform theoretical findings into **working exploits**.

---

## Vulnerability Targets

### Primary Targets (9 "Exploitable")

These vulnerabilities **passed the verifier** and have clear exploitation potential:

| # | Patch | ISO Rule | Vulnerability Type | Exploitation Primitive |
|---|-------|----------|-------------------|------------------------|
| 1 | `5_06a_argcomp` | 5.6 | Function pointer mismatch | Register/stack corruption |
| 2 | `5_06b_argcomp` | 5.6 | Wrong argument count | Stack memory overwrite |
| 3 | `5_10a_exploit_intptrconv` | 5.10 | Pointer truncation bypass | **OOB-R (Info disclosure)** |
| 4 | `5_14_nullref` | 5.14 | Null pointer dereference | Invalid memory access |
| 5 | `5_17_swtchdflt` | 5.17 | Missing switch default | Undefined control flow |
| 6 | `5_20a_libptr` | 5.20 | Buffer overflow (8 bytes) | Adjacent stack corruption |
| 7 | `5_20c_libptr` | 5.20 | Type confusion overflow | 12-byte buffer overflow |
| 8 | `5_35_uninit_mem` | 5.35 | Uninitialized memory read | **Kernel stack data leak** |
| 9 | `5_35a_unint_mem` | 5.35 | Uninitialized memory read | **Kernel stack data leak** |

### Secondary Targets (12 "Limited")

These passed the verifier but have **constrained** exploitation potential:

| # | Patch | ISO Rule | Vulnerability Type | Notes |
|---|-------|----------|-------------------|-------|
| 10 | `5_06d_argcomp` | 5.6 | Wrong argument types | Value truncation |
| 11 | `5_14a_nullref` | 5.14 | Null pointer dereference | No escape path |
| 12 | `5_16b_signconv` | 5.16 | Signed conversion | Logic errors only |
| 13 | `5_33a_restrict` | 5.33 | Restrict pointer violation | Stack corruption |
| 14 | `5_33b_restrict` | 5.33 | Restrict pointer violation | Local stack only |
| 15 | `5_36a_ptrobj` | 5.36 | Pointer comparison UB | Memory layout leak |
| 16 | `5_36b_ptrobj` | 5.36 | Context pointer comparison | Kernel layout leak |
| 17 | `5_36c_ptrobj` | 5.36 | Map pointer comparison | Heap organization leak |
| 18 | `5_39_taintnoproto` | 5.39 | Tainted function pointer | Unpredictable behavior |
| 19 | `5_45_invfmtstr` | 5.45 | Invalid format strings | Address leak via logging |
| 20 | `5_46b_taintsink` | 5.46 | Tainted memory copy | Packet alteration |

**Total: 21 vulnerabilities to investigate**

---

## Team Distribution

With **4 team members** and **21 vulnerabilities**:

| Member | Primary (YES) | Secondary (LIMITED) | Total |
|--------|---------------|---------------------|-------|
| Person 1 | 2-3 | 3 | ~5-6 |
| Person 2 | 2-3 | 3 | ~5-6 |
| Person 3 | 2-3 | 3 | ~5-6 |
| Person 4 | 2-3 | 3 | ~5-6 |

### Suggested Distribution by Category

```
Option A: By Vulnerability Type
┌─────────────────────────────────────────────────────────────────┐
│ Person 1: Memory Bugs                                           │
│   → 5_20a_libptr, 5_20c_libptr, 5_35_uninit_mem, 5_35a_unint_mem│
│   → Focus: Buffer overflows, uninitialized memory reads         │
│                                                                 │
│ Person 2: Type Confusion & Conversions                          │
│   → 5_06a_argcomp, 5_06b_argcomp, 5_06d_argcomp                 │
│   → 5_10a_exploit_intptrconv, 5_16b_signconv                    │
│   → Focus: Argument mismatches, pointer/integer conversions     │
│                                                                 │
│ Person 3: Pointer Operations                                    │
│   → 5_14_nullref, 5_14a_nullref                                 │
│   → 5_36a_ptrobj, 5_36b_ptrobj, 5_36c_ptrobj                    │
│   → Focus: Null derefs, pointer comparisons, layout leaks       │
│                                                                 │
│ Person 4: Control Flow & Data Tainting                          │
│   → 5_17_swtchdflt, 5_33a_restrict, 5_33b_restrict              │
│   → 5_39_taintnoproto, 5_45_invfmtstr, 5_46b_taintsink          │
│   → Focus: Missing defaults, tainted data propagation           │
└─────────────────────────────────────────────────────────────────┘
```

---

## Exploitation Workflow

### Phase 1: Environment Setup (Do First)

Before analyzing any vulnerability, ensure your environment works:

```bash
# 1. Create and start the test VM
cd ../ebpf-exploitability-test/codebase/virt/
./vmctl.sh create ~/.ssh/id_rsa.pub
./vmctl.sh connect

# 2. Inside VM: Verify tools
clang --version          # Should have clang 14+
bpftool version          # Should work
uname -r                 # Should be 6.8.x

# 3. Test baseline compilation
cd ~/ebpf-tests/XDPs/xdp_synproxy/
make clean && make       # Should compile with only warnings

# 4. Test baseline loading
sudo bpftool prog load xdp_synproxy_kern.o /sys/fs/bpf/test
sudo bpftool prog list   # Should show loaded program
sudo rm /sys/fs/bpf/test # Cleanup
```

### Phase 2: Test Patch Workflow

Verify you can apply, test, and revert patches:

```bash
# Apply a patch
git apply patches/5_10a_exploit_intptrconv/0001-feat-5_10a_exploit.patch

# Compile
make clean && make

# Check verifier (verbose)
sudo bpftool prog load xdp_synproxy_kern.o /sys/fs/bpf/test verbose 2>&1 | tee verifier.log

# Analyze bytecode
llvm-objdump -d xdp_synproxy_kern.o > bytecode.txt

# Revert patch
git checkout -- xdp_synproxy_kern.c
# OR
git apply -R patches/5_10a_exploit_intptrconv/0001-feat-5_10a_exploit.patch
```

---

## Per-Vulnerability Analysis Template

For **each assigned vulnerability**, follow this workflow:

### Step 1: Understand the Bug

```markdown
## Vulnerability: [PATCH_NAME]

### ISO Rule
- **Rule Number**: 5.XX
- **Rule Title**: [From ISO-IEC TS 17961-2013]
- **Violation**: [What the code does wrong]

### Source Code Analysis
- **File Modified**: xdp_synproxy_kern.c
- **Function Affected**: [function name]
- **Lines Changed**: [line numbers]

### What the Patch Does
[Explain in your own words what vulnerable code is injected]
```

### Step 2: Verify Compilation & Loading

```markdown
### Compilation
- **Result**: ✅ Success / ❌ Failed
- **Warnings**: [list any new warnings]

### Verifier
- **Result**: ✅ Passed / ❌ Rejected
- **Verifier Output**: [key lines from verbose output]
```

### Step 3: Bytecode Analysis (The Translation Gap)

```markdown
### Bytecode Analysis
- **Does vulnerability survive compilation?**: YES / NO
- **Relevant bytecode instructions**:
  ```
  [paste relevant llvm-objdump output]
  ```
- **Verifier's register state belief**:
  ```
  [paste verifier's understanding of registers]
  ```
- **Actual runtime state**:
  [explain what actually happens]
```

### Step 4: Exploitation Strategy

```markdown
### Exploitation Primitive
- **Type**: OOB-R / OOB-W / Type Confusion / Info Leak / Other
- **What can be read/written**: [specific memory regions]
- **Trigger mechanism**: [how to trigger via network packet]

### Attack Plan
1. [Step 1]
2. [Step 2]
3. [Step 3]
...
```

### Step 5: PoC Development

```markdown
### Proof of Concept

#### Trigger Code
```c
// Code that triggers the vulnerability
```

#### Expected vs Actual Behavior
- **Expected (safe)**: [what should happen]
- **Actual (vulnerable)**: [what actually happens]

#### Evidence
- **Screenshot/Log**: [proof of exploitation]
- **Impact demonstrated**: [what we achieved]
```

### Step 6: Impact Assessment

```markdown
### Security Impact

| Aspect | Assessment |
|--------|------------|
| **Confidentiality** | HIGH / MEDIUM / LOW / NONE |
| **Integrity** | HIGH / MEDIUM / LOW / NONE |
| **Availability** | HIGH / MEDIUM / LOW / NONE |

### Exploitability
- **Attack Vector**: Local (requires CAP_BPF)
- **Complexity**: HIGH / MEDIUM / LOW
- **Privileges Required**: CAP_BPF / Root

### Real-World Impact
[Explain what an attacker could achieve with this vulnerability]
```

---

## Directory Structure

Organize your work as follows:

```
src/
├── README.md                    # This file
├── common/                      # Shared exploitation utilities
│   ├── exploit_helpers.h        # Common macros, structures
│   └── trigger_utils.c          # Packet crafting utilities
│
├── 5_06a_argcomp/              # One folder per vulnerability
│   ├── README.md               # Analysis following template above
│   ├── analysis/
│   │   ├── bytecode.txt        # llvm-objdump output
│   │   └── verifier.log        # Verifier verbose output
│   ├── poc/
│   │   ├── exploit.c           # PoC code
│   │   └── Makefile
│   └── evidence/
│       └── screenshots/        # Proof of exploitation
│
├── 5_10a_exploit_intptrconv/
│   └── ...
│
└── [other vulnerabilities]/
```

---

## Tools Reference

### Inherited Tools

| Tool | Location | Purpose |
|------|----------|---------|
| `xvtlas` | `../ebpf-exploitability-test/codebase/xvtlas/` | Automated patch testing |
| `vmctl.sh` | `../ebpf-exploitability-test/codebase/virt/` | VM lifecycle management |
| `start_session.sh` | `../ebpf-exploitability-test/codebase/XDPs/xdp_synproxy/` | tmux test session |

### Analysis Tools

```bash
# Bytecode disassembly
llvm-objdump -d xdp_synproxy_kern.o

# Loaded program dump
sudo bpftool prog dump xlated id <ID>

# Verbose verifier output
sudo bpftool prog load prog.o /sys/fs/bpf/test verbose

# JIT disassembly (if JIT enabled)
sudo bpftool prog dump jited id <ID>

# Map inspection
sudo bpftool map list
sudo bpftool map dump id <ID>
```

### Packet Crafting (for triggering vulnerabilities)

```bash
# Using scapy (Python)
sudo scapy

# Using netcat
nc -v <target> <port>

# Using hping3 (SYN floods)
sudo hping3 -S -p 80 --flood <target>
```

---

## Progress Tracking

### My Assigned Vulnerabilities

| # | Vulnerability | Status | Notes |
|---|---------------|--------|-------|
| 1 | `____________` | ⬜ Not Started | |
| 2 | `____________` | ⬜ Not Started | |
| 3 | `____________` | ⬜ Not Started | |
| 4 | `____________` | ⬜ Not Started | |
| 5 | `____________` | ⬜ Not Started | |

**Status Legend:**
- ⬜ Not Started
- 🔄 In Progress — Environment Setup
- 🔍 In Progress — Analysis
- 🛠️ In Progress — PoC Development
- ✅ Complete
- ❌ Not Exploitable (with justification)

---

## Quick Reference: Patch Locations

All patches are in: `../ebpf-exploitability-test/codebase/XDPs/xdp_synproxy/patches/`

```bash
# List all patches
ls ../ebpf-exploitability-test/codebase/XDPs/xdp_synproxy/patches/

# View patch content
cat ../ebpf-exploitability-test/codebase/XDPs/xdp_synproxy/patches/5_10a_exploit_intptrconv/*.patch

# Previous team's detailed analysis
cat ../ebpf-exploitability-test/codebase/XDPs/xdp_synproxy/README.md
```

---

## Next Steps

1. ✅ Read main README (understand eBPF, verifier, threat model)
2. ⬜ Set up VM environment (Phase 1 above)
3. ⬜ Test patch workflow (Phase 2 above)
4. ⬜ Coordinate with team on vulnerability distribution
5. ⬜ Begin analysis of assigned vulnerabilities
6. ⬜ Develop PoCs
7. ⬜ Document findings
8. ⬜ Write individual report
