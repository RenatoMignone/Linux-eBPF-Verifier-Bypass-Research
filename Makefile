# Root Makefile for VM lifecycle management

.PHONY: create destroy

# Path to the VM control script
VMCTL=ebpf-exploitability-test/codebase/virt/vmctl.sh

# Default SSH key (edit as needed)
SSH_KEY?=$(HOME)/.ssh/SVT_Project/svt_vm.pub

create:
	$(VMCTL) create $(SSH_KEY)

destroy:
	$(VMCTL) destroy
