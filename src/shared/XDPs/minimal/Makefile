# Minimal eBPF Exploit PoC Makefile

CLANG = clang
LLC = llc
ARCH = $(shell uname -m | sed 's/x86_64/x86/')

# eBPF program
KERN_SOURCES = 5_10a_exploit.c
KERN_OBJECTS = $(KERN_SOURCES:.c=.o)

# Compiler flags for eBPF
CFLAGS = -O2 -g -Wall -Werror \
	-D__KERNEL__ -D__BPF_TRACING__ \
	-Wno-unused-value -Wno-pointer-sign \
	-Wno-compare-distinct-pointer-types \
	-Wno-gnu-variable-sized-type-not-at-end \
	-Wno-address-of-packed-member -Wno-tautological-compare \
	-Wno-unknown-warning-option

# No optimization flags for preserving exploit code
CFLAGS_NO_OPT = -O0 -g -Wall \
	-D__KERNEL__ -D__BPF_TRACING__ \
	-Wno-unused-value -Wno-pointer-sign \
	-Wno-compare-distinct-pointer-types \
	-Wno-gnu-variable-sized-type-not-at-end \
	-Wno-address-of-packed-member -Wno-tautological-compare \
	-Wno-unknown-warning-option

# Include paths
INCLUDES = -I/usr/include/$(shell uname -m)-linux-gnu

all: $(KERN_OBJECTS)

# Build without optimizations to preserve exploit code
no-opt: 5_10a_exploit_no_opt.o

%.o: %.c
	$(CLANG) $(CFLAGS) $(INCLUDES) -target bpf -c $< -o $@

%_no_opt.o: %.c
	$(CLANG) $(CFLAGS_NO_OPT) $(INCLUDES) -target bpf -c $< -o $@

clean:
	rm -f *.o

load: 5_10a_exploit.o
	sudo ip link set dev lo xdp obj 5_10a_exploit.o sec xdp verbose

unload:
	sudo ip link set dev lo xdp off

# Analyze bytecode
disasm: 5_10a_exploit.o
	llvm-objdump -S 5_10a_exploit.o

# Analyze bytecode without optimizations
disasm-no-opt: 5_10a_exploit_no_opt.o
	llvm-objdump -S 5_10a_exploit_no_opt.o

# Show verifier log
verify: 5_10a_exploit.o
	bpftool prog load 5_10a_exploit.o /sys/fs/bpf/5_10a_exploit

.PHONY: all clean load unload disasm verify no-opt disasm-no-opt
