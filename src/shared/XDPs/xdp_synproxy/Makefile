# Define source files for BPF kernel program and user-space application
BPF_SRC := xdp_synproxy_kern.c
USER_SRC := xdp_synproxy.c

# Define the name of the generated BPF skeleton header file
# This file will be included by the user-space program to interact with the BPF program.
SKEL_H := $(BPF_SRC:.c=.skel.h) # e.g., xdp_synproxy_kern.skel.h

# Define the output BPF object file and user-space executable
BPF_OBJ := $(BPF_SRC:.c=.bpf.o)
USER_BIN := $(USER_SRC:.c=) # Removes the .c extension to get the executable name (e.g., xdp_synproxy)

# Compiler and tools
CLANG := clang
LLC := llc
BPFTOOL := bpftool

# Architecture for BPF compilation (default to x86)
ARCH ?= x86

# CFLAGS for BPF programs (kernel-side)
# These flags target the BPF bytecode and include necessary kernel headers.
BPF_CFLAGS := -O2 -g -Wall -target bpf \
              -I. \
              -I../tools \
              -I.. \
              -I/usr/include/bpf \
              -I/usr/include/x86_64-linux-gnu/linux \
              -D__TARGET_ARCH_$(ARCH)

# CFLAGS for user-space programs
# These flags are for standard C compilation and link against libbpf and other common libraries.
# Added -lbfd to link the libbpf library.
USER_CFLAGS := -O2 -g -Wall \
               -I. \
               -I../tools \
               -I.. \
               -I/usr/include/bpf \
	       -I/usr/include \
               -I/usr/include/x86_64-linux-gnu/linux \
               -lbpf -lbfd -lelf -lz -lcap -lmnl # Added -lbfd here

# Default target: build everything
all: $(BPF_OBJ) $(USER_BIN)

# Rule to compile BPF .c files into .o (BPF object files)
# This rule applies to any .c file that needs to be compiled for the BPF target.
%.bpf.o: %.c
	$(CLANG) $(BPF_CFLAGS) -c $< -o $@

# Rule to compile user-space .c files into executables
# This rule applies to any .c file that needs to be compiled into a user-space binary.
# This target depends on the user-space source file.
$(USER_BIN): $(USER_SRC)
	@echo "Compiling user-space program: $@"
	$(CLANG) $(USER_CFLAGS) $< -o $@

# Clean target: remove all generated files
clean:
	@echo "Cleaning compiled files..."
	rm -f $(BPF_OBJ) $(USER_BIN)
	@echo "Clean complete."

# Phony targets to ensure Make doesn't confuse them with files of the same name
.PHONY: all clean
