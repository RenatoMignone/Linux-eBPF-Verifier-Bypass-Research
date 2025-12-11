# eBPF Vulnerability Testing Suite - ISO-IEC TS 17961-2013

This project provides a comprehensive testing environment for eBPF/XDP programs with **intentional vulnerability patches** based on the **ISO-IEC TS 17961-2013** standard. The suite includes automated VM provisioning, vulnerability injection, and verifier analysis tools.

## Overview

The project demonstrates security vulnerabilities in eBPF/XDP kernel code by implementing various rules from ISO-IEC TS 17961-2013. Each vulnerability is carefully crafted to show:
- How coding standard violations manifest in eBPF programs
- eBPF verifier behavior and limitations
- Compiler diagnostic capabilities
- Real-world security implications

## Project Structure

```
ebpf-tests-3-1/
├── virt/                    # VM management scripts
│   └── vmctl.sh            # VM creation, destruction, and connection
├── XDPs/
│   └── xdp_synproxy/       # XDP SYN proxy implementation (from Linux selftests)
│       ├── apply_rules     # Network rules configuration script
│       ├── start_session.sh # Tmux session setup for testing
│       └── patches/        # ISO-IEC TS 17961-2013 vulnerability patches
├── xvtlas/                 # XDP Verifier Launch Automation Suite
├── pretty-verifier/        # Python verifier output formatter
└── docs/                   # Documentation and references
```

## Quick Start

### 1. VM Environment Setup

The project uses a standardized Ubuntu VM environment for consistent testing:

```bash
cd virt/
./vmctl.sh create ~/.ssh/id_rsa.pub    # Create VM with your SSH key
./vmctl.sh connect                      # Connect to the VM
```

**VM Management Commands:**
- `./vmctl.sh create <ssh_pubkey_file>` - Create and configure new VM
- `./vmctl.sh destroy` - Destroy the VM
- `./vmctl.sh connect` - Connect/reconnect to existing VM

**VM Connection Details:**
- After VM creation, you may need to press **`ENTER`** a few times to reach the login prompt
- Default login credentials (as specified in `virt/user-data.yaml`):
  - Username: `user`
  - Password: `` (empty password)
**VM Connection Details:**
- After VM creation, you can connect either via:
  - `./vmctl.sh connect` - Direct console connection
  - `ssh user@<vm-ip>` - SSH connection using the pubkey configured during creation

The script uses cloud-config with a modified `user-data.yaml` file to provision an Ubuntu environment with all necessary development tools pre-installed. The SSH public key is automatically copied into the VM during the creation process.

### 2. Repository Setup (Inside VM)

```bash
# Clone the repository inside the VM (with submodules for pretty-verifier)
git clone --recurse-submodules <repository-url> ebpf-tests-3-1
cd ebpf-tests-3-1
```

### 3. XDP SynProxy Configuration

Configure the necessary network rules for XDP SynProxy operation:

```bash
cd XDPs/xdp_synproxy/
./apply_rules
```

This script configures:
- Network interface settings
- iptables rules for SYN proxy operation
- Kernel parameters for eBPF program loading

### 4. Testing Environment Setup

The project includes a testing script to start a tmux session with the necessary components:

```bash
cd XDPs/xdp_synproxy/
./start_session.sh
```

To close the session and clean the env :

```bash
//From inside the tmux session :
./kill_session.sh
```

**Network Topology:**
The current testing setup uses a simplified topology:
- **XDP SynProxy**: Running inside the VM on the network interface
- **Netcat Server**: Running inside the VM to test connections
- **Netcat Client**: Running on the host machine (outside VM) connecting to the server

**Design Decision:**
Initially, we followed the Linux kernel selftests approach using a 3-veth (virtual Ethernet) interface topology. However, virtual interfaces conflicted with the SYN cookie functionality, due to checksum problems, as noted in the Linux kernel source code comments. This simplified approach provides a realistic testing environment without these conflicts.

### 5. Vulnerability Patches

All vulnerability patches and detailed documentation can be found in the [XDP SynProxy README](./XDPs/xdp_synproxy/README.md).

This section contains:
- Complete patch directory listing with vulnerability types
- Detailed implementation explanations for each rule
- Verifier behavior analysis
- Progress tracking and examples

The patches target `xdp_synproxy_kern.c` and demonstrate ISO-IEC TS 17961-2013 rule violations while maintaining core SYN proxy functionality.

For a complete list of rules not applicable to XDP/eBPF environments, see the [XDP SynProxy README](./XDPs/xdp_synproxy/README.md#rules-not-applicable-to-xdpebpf).


### 6. XVTLAS - XDP Verifier Launch Automation Suite

**XVTLAS** automates the entire process of:
- Compiling eBPF programs
- Applying vulnerability patches
- Loading programs with bpftool
- Analyzing verifier output
- Generating structured reports

#### Installation

```bash
cd xvtlas/
# See xvtlas/README.md for detailed compilation instructions
# The tool is written in Go and requires compilation
go build -o xvtlas .
```

#### Usage

```bash
# Apply and test a single vulnerability patch interactively
./xvtlas --run-single "./XDPs/xdp_synproxy/patches/5_45_invfmtstr/*.patch" \
         --base-file "./XDPs/xdp_synproxy/xdp_synproxy_kern.c"

# Run comprehensive testing on all patches
./xvtlas --export "./output/" \
         --kernel "6.8.58" \
         --patch-path "./XDPs/xdp_synproxy/patches/" \
         --base-file "./XDPs/xdp_synproxy/xdp_synproxy_kern.c" \
         --pretty "./pretty-verifier/pretty_verifier.py" \
         --save-logs \
         --verbose \
         --interactive

# Clean up after interactive session
./xvtlas --destroy

# Manual patch application (alternative approach)
cd XDPs/xdp_synproxy/
git apply patches/5_45_invfmtstr/0001-feat-5.45-invfmtstr-*.patch
make  # Compile manually
```

For detailed usage instructions, refer to `xvtlas/README.md`.


## Educational Purpose

This suite is designed for:
- **Security Research**: Understanding eBPF/XDP vulnerability patterns
- **Secure Coding Training**: Learning to avoid common pitfalls
- **Verifier Analysis**: Understanding eBPF verifier capabilities and limitations
- **Compiler Diagnostics**: Testing static analysis tool effectiveness

## Architecture

- **VM Environment**: Standardized Ubuntu cloud-config setup
- **XDP Target**: Real-world SYN proxy implementation
- **Patch System**: Git-based vulnerability injection
- **Automation**: Go-based testing orchestration
- **Analysis**: Python-based verifier output formatting

## Requirements

- Host system with KVM/QEMU support
- SSH key pair for VM access
- Go compiler (for XVTLAS)
- Python 3.x (for pretty-verifier)


## Developer Flow

For developers who want to contribute new vulnerability patches:

### Local Development Workflow

1. **Set up your preferred development environment** on the host machine
2. **Identify the target vulnerability** from ISO-IEC TS 17961-2013 standard
3. **Modify the source code** locally:
   ```bash
   # Edit the file in your preferred editor/IDE
   vim XDPs/xdp_synproxy/xdp_synproxy_kern.c

   # Focus on functions like tcp_dissect() for realistic vulnerability injection
   # Add clear comments explaining the vulnerability and expected behavior
   ```

4. **Transfer modified file to VM** using one of these methods:
   ```bash
   # Option 1: SCP copy
   scp XDPs/xdp_synproxy/xdp_synproxy_kern.c user@<vm-ip>:~/ebpf-tests-3-1/XDPs/xdp_synproxy/

   # Option 2: Mount via SSHFS (recommended for iterative development)
   mkdir ~/vm-mount
   sshfs user@<vm-ip>:~/ebpf-tests-3-1 ~/vm-mount
   # Now you can edit files directly in ~/vm-mount/
   ```

5. **Test inside the VM**:
   ```bash
   # Connect to VM
   ./virt/vmctl.sh connect

   # Compile with appropriate warning flags
   cd ~/ebpf-tests-3-1/XDPs/xdp_synproxy/
   make

   # Test the vulnerability behavior
   ./start_session.sh  # Start tmux testing environment
   ```

6. **Create patch** when satisfied with the vulnerability implementation:
   ```bash
   # Commit your changes
   git add xdp_synproxy_kern.c
   git commit -m "feat: 5.XX-rulename - Description of vulnerability"

   # Generate patch file
   git format-patch HEAD~1 -o patches/5_XX_rulename/
   ```

### Testing Your Vulnerability

- **Compiler Diagnostics**: Test with various warning flags (`-Wall`, `-Wformat`, `-Wpointer-arith`)
- **eBPF Verifier**: Load the program and observe verifier behavior
- **Runtime Behavior**: Use the tmux session to test actual packet processing
- **Documentation**: Update comments to explain expected vs actual behavior

## Contributing

When adding new vulnerability patches:
1. Follow the ISO-IEC TS 17961-2013 standard classification
2. Include detailed comments explaining the vulnerability
3. Document expected compiler/verifier behavior
4. Test in isolated VM environment
5. Update patch directory structure accordingly
6. Use the Developer Flow above for consistent development process

## References

- [ISO-IEC TS 17961-2013](https://www.iso.org/standard/61134.html) - C Static Analysis Standard
- [eBPF Documentation](https://ebpf.io/) - Extended Berkeley Packet Filter
- [XDP Documentation](https://www.kernel.org/doc/html/latest/networking/filter.html) - eXpress Data Path
- [BPF Verifier](https://www.kernel.org/doc/html/latest/bpf/verifier.html) - eBPF Program Verification
- [XDP SynProxy Examples](https://github.com/xdp-project/bpf-examples/tree/main/xdp-synproxy) - XDP-project BPF examples
- [Linux Kernel SynProxy Selftests](https://github.com/torvalds/linux/blob/v6.8/tools/testing/selftests/bpf/prog_tests/xdp_synproxy.c) - Linux kernel selftests

## License

TODO add license

## Authors

- Francesco Rollo
- Gianfranco Trad
- Giorgio Fardo
- Giovanni Nicosia

Developed for security research and educational purposes in the context of eBPF/XDP vulnerability analysis.

---
