# Compiler and flags
CLANG := clang
LLC := llc
ARCH ?= x86

CFLAGS := -O2 -g -Wall -target bpf \
          -I. \
          -I../tools \
          -I.. \
          -I/usr/include/bpf \
          -I/usr/include/x86_64-linux-gnu/linux \
          -D__TARGET_ARCH_$(ARCH)

OBJ := $(SRC:.c=.o)

all: $(OBJ)

$(OBJ): $(SRC)
	$(CLANG) $(CFLAGS) -c $< -o $@

clean:
	rm -f $(OBJ)

.PHONY: all clean

