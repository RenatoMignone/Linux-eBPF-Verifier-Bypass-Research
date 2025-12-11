# Modifications to Inherited Codebase

> This document tracks all modifications made to the previous team's work in `ebpf-exploitability-test/`.

---

## Summary of Changes

| Date | File | Modification | Reason |
|------|------|--------------|--------|
| 2025-12-09 | `virt/user-data.yaml` | Created file | Missing from repository |
| 2025-12-09 | `virt/vmctl.sh` | Replaced sed with awk | SSH key special characters breaking substitution |
| 2025-12-09 | Host system | Installed virtualization packages | `virt-install` command not found |
| 2025-12-09 | `virt/vmctl.sh` | Added virtiofs shared folder support | Enable host↔VM file sharing |
| 2025-12-09 | `virt/user-data.yaml` | Added shared folder auto-mount | Auto-mount `/mnt/shared` at boot |
| 2025-12-09 | `virt/user-data.yaml` | Added password authentication | SSH key format issues (PEM vs OpenSSH) |
| 2025-12-09 | `virt/user-data.yaml` | Set hostname to `ubuntu`, password to `u` | Simplify VM login credentials |
| 2025-12-11 | `virt/vmctl.sh` | Added codebase auto-copy to shared folder | Automate file transfer on VM creation |
| 2025-12-11 | `virt/user-data.yaml` | Added full eBPF/XDP auto-setup | Automate module loading, vmlinux.h, compilation |
| 2025-12-11 | Host system | Installed `virtiofsd` | Required for virtiofs shared folder |

---

## Detailed Modifications

### 1. Created `virt/user-data.yaml`

**Date**: 2025-12-09  
**File**: `ebpf-exploitability-test/codebase/virt/user-data.yaml`  
**Status**: NEW FILE

**Problem**: The `vmctl.sh` script referenced `user-data.yaml` but the file was missing from the repository. Running `./vmctl.sh create` failed because the cloud-init configuration didn't exist.

**Solution**: Created the cloud-init user-data file with:

```yaml
#cloud-config

# User configuration
users:
  - name: ubuntu
    sudo: ALL=(ALL) NOPASSWD:ALL
    shell: /bin/bash
    ssh_authorized_keys:
      - __SSH_PUBLIC_KEY__

# Package updates and installation
package_update: true
package_upgrade: true

packages:
  # Build essentials
  - build-essential
  - git
  - make
  - pkg-config
  
  # LLVM/Clang for eBPF compilation
  - clang
  - llvm
  - libclang-dev
  
  # BPF tools and libraries
  - libbpf-dev
  - linux-tools-common
  - linux-tools-generic
  - bpftool
  
  # Kernel headers
  - linux-headers-generic
  
  # Network tools for testing
  - tcpdump
  - wireshark-common
  - tshark
  - netcat-openbsd
  - hping3
  - iperf3
  - net-tools
  
  # Development utilities
  - tmux
  - vim
  - curl
  - wget
  - jq

# Run commands after boot
runcmd:
  - ln -sf /usr/lib/linux-tools/*/bpftool /usr/local/bin/bpftool || true
  - mkdir -p /home/ubuntu/ebpf-tests
  - chown -R ubuntu:ubuntu /home/ubuntu/ebpf-tests
  - echo 1 > /proc/sys/net/core/bpf_jit_enable
  - echo "net.core.bpf_jit_enable = 1" >> /etc/sysctl.conf

final_message: "eBPF development VM ready! Boot took $UPTIME seconds."
```

---

### 2. Fixed `virt/vmctl.sh` — SSH Key Substitution

**Date**: 2025-12-09  
**File**: `ebpf-exploitability-test/codebase/virt/vmctl.sh`  
**Status**: MODIFIED (line ~51)

**Problem**: The `sed` command failed with error:
```
sed: -e expression #1, char 47: unterminated `s' command
```

This occurred because SSH public keys contain special characters (`/`, `+`, `=`, `&`) that break the `sed` substitution pattern, even when using `|` as a delimiter or attempting to escape characters.

**Original Code**:
```bash
SSH_KEY=$(<"$SSH_KEY_FILE")
sed -i "s|__SSH_PUBLIC_KEY__|$SSH_KEY|" "$USER_DATA"
```

**First Attempt (Still Failed)**:
```bash
SSH_KEY=$(<"$SSH_KEY_FILE")
# Escape special characters in SSH key for sed replacement
SSH_KEY_ESCAPED=$(printf '%s\n' "$SSH_KEY" | sed -e 's/[&/\]/\\&/g')
sed -i "s|__SSH_PUBLIC_KEY__|$SSH_KEY_ESCAPED|" "$USER_DATA"
```

**Final Fix (Working)**:
```bash
SSH_KEY=$(<"$SSH_KEY_FILE")
# Use awk instead of sed to avoid escaping issues with SSH keys
awk -v key="$SSH_KEY" '{gsub(/__SSH_PUBLIC_KEY__/, key); print}' "$USER_DATA" > "${USER_DATA}.tmp" && mv "${USER_DATA}.tmp" "$USER_DATA"
```

**Explanation**: `awk` handles variable substitution natively without requiring special character escaping. The `-v key="$SSH_KEY"` passes the key as an awk variable, and `gsub()` performs the replacement safely regardless of the key's content.

---

### 3. Installed Host Virtualization Dependencies

**Date**: 2025-12-09  
**Location**: Host system (not inherited codebase)  
**Status**: PREREQUISITE

**Problem**: Running `vmctl.sh create` failed with:
```
sudo: virt-install: command not found
```

The previous team's README mentioned the dependencies but didn't include complete installation instructions.

**Solution**: Install virtualization packages on Ubuntu 24.04:

```bash
sudo apt update
sudo apt install -y virtinst libvirt-daemon-system libvirt-clients qemu-system qemu-utils
```

**Note**: On Ubuntu 24.04, the package is `virtinst` (not `virt-install`). The README's suggested command `apt install virt-install` would fail with:
```
E: Package 'virt-install' has no installation candidate
```

**Packages and their purposes**:
| Package | Purpose |
|---------|---------|
| `virtinst` | Provides `virt-install` command |
| `libvirt-daemon-system` | Libvirt virtualization daemon |
| `libvirt-clients` | CLI tools (virsh, etc.) |
| `qemu-system` | QEMU system emulators |
| `qemu-utils` | Disk image utilities (qemu-img) |

---

### 4. Added Virtiofs Shared Folder Support to `vmctl.sh`

**Date**: 2025-12-09  
**File**: `ebpf-exploitability-test/codebase/virt/vmctl.sh`  
**Status**: MODIFIED

**Problem**: No way to share files between host and VM. Docker is unsuitable for eBPF development because it shares the host kernel, making it unsafe for kernel exploitation testing. A VM with a shared folder is needed for:
- Editing code on host with familiar tools
- Running/testing in isolated VM kernel
- Easy transfer of compiled eBPF programs

**Solution**: Added virtiofs filesystem support to the VM creation command:

```bash
# Added variable
SHARED_FOLDER="${CODEBASE_DIR}/shared"

# Added to virt-install command
--filesystem source="${SHARED_FOLDER}",target=shared,driver.type=virtiofs \
--memorybacking source.type=memfd,access.mode=shared \
```

**Explanation**:
- `--filesystem`: Creates a virtiofs mount point from host `codebase/shared/` to VM tag `shared`
- `--memorybacking`: Required for virtiofs to work (shared memory access)
- The VM can then mount this at `/mnt/shared`

---

### 5. Added Shared Folder Auto-Mount to `user-data.yaml`

**Date**: 2025-12-09  
**File**: `ebpf-exploitability-test/codebase/virt/user-data.yaml`  
**Status**: MODIFIED

**Problem**: After adding virtiofs support to vmctl.sh, the VM still needed to mount the shared folder manually each boot.

**Solution**: Added automatic mount configuration via cloud-init:

```yaml
# Auto-mount shared folder from host
mounts:
  - [ shared, /mnt/shared, virtiofs, "rw,relatime", "0", "0" ]

# Ensure mount point exists
runcmd:
  - mkdir -p /mnt/shared
```

**Usage**: 
- Host: Place files in `ebpf-exploitability-test/codebase/shared/`
- VM: Access files at `/mnt/shared/`

---

### 6. Added Password Authentication to `user-data.yaml`

**Date**: 2025-12-09  
**File**: `ebpf-exploitability-test/codebase/virt/user-data.yaml`  
**Status**: MODIFIED

**Problem**: SSH key authentication failed due to key format mismatch. The generated SSH key was in PEM format:
```
-----BEGIN RSA PRIVATE KEY-----
```

But cloud-init/OpenSSH on Ubuntu 24.04 expects OpenSSH format:
```
-----BEGIN OPENSSH PRIVATE KEY-----
```

**Solution**: Added password authentication as a fallback:

```yaml
# SSH settings - allow password authentication
ssh_pwauth: true

# User configuration with password
users:
  - name: ubuntu
    lock_passwd: false

chpasswd:
  expire: false
  users:
    - name: ubuntu
      type: text
      password: u
```

**Note**: SSH key authentication still configured but password provides reliable console/SSH access regardless of key format.

---

### 7. Simplified VM Credentials

**Date**: 2025-12-09  
**File**: `ebpf-exploitability-test/codebase/virt/user-data.yaml`  
**Status**: MODIFIED

**Problem**: User requested simpler login credentials for easier VM access during development.

**Solution**: 
- **Hostname**: Changed to `ubuntu`
- **Username**: `ubuntu` (unchanged)
- **Password**: `u` (single character)

```yaml
hostname: ubuntu

chpasswd:
  expire: false
  users:
    - name: ubuntu
      type: text
      password: u
```

**Note**: Cloud-init does not support truly empty passwords for security reasons. Single character `u` is the simplest valid password that works reliably.

**Final VM Access**:
```bash
# Console login
Username: ubuntu
Password: u

# SSH (if key works)
ssh ubuntu@<vm-ip>

# Or with password
ssh ubuntu@<vm-ip>  # then enter 'u'
```

---

## Future Modifications

*Document any additional changes here as the project progresses.*
