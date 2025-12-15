# Root Makefile for VM lifecycle management
# Wraps the underlying vmctl.sh script for easy access

.PHONY: help create connect destroy

# Path to the VM control script
VMCTL := ./ebpf-exploitability-test/codebase/virt/vmctl.sh

# Default SSH key (can be overridden via command line: make create SSH_KEY=~/.ssh/other.pub)
SSH_KEY ?= $(HOME)/.ssh/SVT_Project/svt_vm.pub

# Default target: Print usage help
help:
	@echo "VM Management Commands:"
	@echo "  make create   - Create and provision the VM (uses SSH_KEY)"
	@echo "  make connect  - Connect to the VM console via SSH"
	@echo "  make destroy  - Power off and delete the VM"
	@echo ""
	@echo "Current Configuration:"
	@echo "  SSH_KEY       = $(SSH_KEY)"
	@echo "  Script Path   = $(VMCTL)"

# Create the VM
create:
	@echo "[*] Creating VM with key: $(SSH_KEY)..."
	@$(VMCTL) create $(SSH_KEY)

# Connect to the VM
connect:
	@echo "[*] Connecting to VM..."
	@$(VMCTL) connect

# Destroy the VM
destroy:
	@echo "[!] Destroying VM..."
	@$(VMCTL) destroy