# XDP SynProxy Vulnerability Testing - ISO-IEC TS 17961-2013

This directory contains vulnerability patches that implement various ISO-IEC TS 17961-2013 secure coding rule violations within the XDP SynProxy implementation.

## Vulnerability Patches

All vulnerability patches target `xdp_synproxy_kern.c`, an XDP-based SYN proxy implementation taken from the Linux kernel selftests. The `patches/` directory contains vulnerability patches for each applicable rule from ISO-IEC TS 17961-2013.

For a complete overview of all patches including compilation results, verifier status, and exploitability analysis, see **[patches.csv](patches.csv)**.

Each vulnerability rule is implemented as a Git commit patch that modifies the base `xdp_synproxy_kern.c` file. These patches can be:
- **Applied manually**: Use `git apply` or `git am` to apply individual patches for manual testing
- **Applied automatically**: Use the XVTLAS tool which handles patch application, compilation, verification, and output export

The patches are designed to demonstrate specific ISO-IEC TS 17961-2013 rule violations while maintaining the core SYN proxy functionality.

#### Legend

- **No** : Blocked by verifier
- **Limited** : Passed by verifier but limited explotability, not really an exploit
- **Yes** : Passed by verifer and can lead to further exploitability

## Rules Not Applicable to XDP/eBPF

The following ISO-IEC TS 17961-2013 rules are **not applicable** to XDP/eBPF environments due to fundamental limitations and architectural differences:

| Rule | Title | Category | Reason | Author |
|------|-------|----------|---------|--------|
| **5.2** | Accessing freed memory | Memory Management | No `free()` function or dynamic memory allocation | @all |
| **5.3** | Accessing shared objects in signal handlers | Signal Handling | BPF helper `bpf_send_signal()` is present but cannot be used to implement this vulnerability as it sends the signal to the user space and woruld not create the race condition related to  the `err_msg` pointer as described in the PDF. | @all |
| **5.5** | Calling functions from signal handlers except abort, _Exit, signal | Signal Handling | `bpf_send_signal()` cannot take custom handler as an argument. | @all |
| **5.7** | Calling signal from interruptible signal handlers | Signal Handling | No custom signal handler possible | @all |
| **5.8** | Calling system | Library Functions | No `system()` function available | @all |
| **5.12** | Copying a FILE object | File Operations | No file structures or FILE type available | @all |
| **5.18** | Failing to close files or free memory | Memory Management | No `malloc()` or file close operations available | @all |
| **5.19** | Failing to detect and handle stdlib errors | Library Functions | Limited standard library support | @all |
| **5.21** | Allocating insufficient memory | Memory Management | No dynamic memory allocation (`malloc`) in eBPF | @all |
| **5.23** | Freeing memory multiple times | Memory Management | No `free()` function available | @all |
| **5.25** | Incorrect use of errno | Discarded | Not particularly relevant to eBPF testing | @all |
| **5.27** | Interleaving stream I/O without flush or positioning | File Operations | No buffered stdio operations | @all |
| **5.29** | Modifying getenv/localeconv/etc. return values | Library Functions | No `getenv()` or `setlocale()` functions | @all |
| **5.32** | Invalid chars to character-handling functions | Library Functions | Can't import `ctype.h` library (error: failed to load: -13 ), No character-handling functions available in eBPF environment | @all |
| **5.34** | Re/freeing non-dynamically allocated memory | Memory Management | No dynamic memory allocation or `free()` operations | @all |
| **5.37** | Tainted strings are passed to a string copying function | Format String | No `strcpy()` | @all |
| **5.38** | Taking size of pointer to get pointed-to size | Discarded | Not useful for eBPF vulnerability testing scenarios | @all |
| **5.41** | Invalid value for fsetpos | File Operations | No file operations available | @all |
| **5.42** | Using object overwritten by getenv/localeconv/etc. | Library Functions | No libc environment functions | @all |
| **5.43** | Char values indistinguishable from EOF | File Operations | No file operations or EOF handling | @all |
| **5.44** | Using reserved identifiers | Discarded | Not relevant for security testing focus | @all |

#### Warning
Rules 5.25, 5.38, and 5.44 are technically applicable to XDP/eBPF but were intentionally discarded as not useful for practical vulnerability testing.

**Rule 5.25**
Rule 5.25 addresses the correct use of errno when interacting with Standard C Library functions. In the XDP/eBPF context, this rule is not relevant for security testing because eBPF programs do not link against the C standard library and do not rely on errno for error signaling. Instead, BPF helpers communicate errors via explicit return values. Misuse of errno cannot occur in this environment and therefore does not introduce security risks. For the security of eBPF/XDP code, this rule can be considered not applicable.

(Signed-by: Giorgio Fardo, Francesco Rollo)

**Rule 5.38**
While 5.38 rule is useful for preventing functional bugs in standard C code, it is not relevant for eBPF/XDP vulnerability testing. eBPF programs array bounds must be explicitly tracked, and pointer size misuse does not introduce exploitable security vulnerabilities. At worst, such code leads to logical errors rather than security issues. Therefore, this rule can be considered not useful in a security review.

(Signed-by: Giorgio Fardo, Francesco Rollo)

**Rule 5.44**
Rule 5.44 is about avoiding the use of reserved identifiers (like errno, or identifiers starting with _ followed by uppercase) to ensure portability and prevent undefined behavior in standard C environments.
In the context of XDP/eBPF, this rule has little relevance for security testing because eBPF programs operate in a restricted environment without the C standard library, and reserved identifier clashes cannot lead to security vulnerabilities. At worst, such violations cause compilation issues, not exploitable conditions.

(Signed-by: Giorgio Fardo, Francesco Rollo)

These exclusions reflect the constrained execution environment of eBPF programs, which operate in kernel space with:
- No dynamic memory allocation
- No signal handling
- Limited standard library access
- No file system operations
- Restricted system call access

## Patches Report

## Base warnings

Warnings present in the base file compilation attempt, so present in all the compilation attempts of the tests.

```
xdp_synproxy_kern.c:208:16: warning: comparison of distinct pointer types ('__u8 *' (aka 'unsigned char *') and 'void *') [-Wcompare-distinct-pointer-types]
  208 |         if (data + sz >= ctx->data_end)
      |             ~~~~~~~~~ ^  ~~~~~~~~~~~~~
xdp_synproxy_kern.c:377:19: warning: comparison of distinct pointer types ('struct ethhdr *' and 'void *') [-Wcompare-distinct-pointer-types]
  377 |         if (hdr->eth + 1 > data_end)
      |             ~~~~~~~~~~~~ ^ ~~~~~~~~
xdp_synproxy_kern.c:385:21: warning: comparison of distinct pointer types ('struct iphdr *' and 'void *') [-Wcompare-distinct-pointer-types]
  385 |                 if (hdr->ipv4 + 1 > data_end)
      |                     ~~~~~~~~~~~~~ ^ ~~~~~~~~
xdp_synproxy_kern.c:401:21: warning: comparison of distinct pointer types ('struct ipv6hdr *' and 'void *') [-Wcompare-distinct-pointer-types]
  401 |                 if (hdr->ipv6 + 1 > data_end)
      |                     ~~~~~~~~~~~~~ ^ ~~~~~~~~
xdp_synproxy_kern.c:419:19: warning: comparison of distinct pointer types ('struct tcphdr *' and 'void *') [-Wcompare-distinct-pointer-types]
  419 |         if (hdr->tcp + 1 > data_end)
      |             ~~~~~~~~~~~~ ^ ~~~~~~~~
5 warnings generated.

```

## Rules
### [5.1 ptrcomp]: Accessing an object through a pointer to an incompatible type

Accessing memory through a pointer to an incompatible type (other than `unsigned char`) is undefined behavior in C. This breaks the **strict aliasing rules** defined in ISO/IEC 9899:2011, 6.5§7, which restricts how objects can be accessed through lvalues of different types. Violating these rules can lead to unpredictable behavior, since the compiler may assume incompatible types do not alias and apply optimizations accordingly.

In eBPF programs, such violations may not always be caught by the verifier, since the verifier mainly ensures memory safety (bounds checking) rather than C language aliasing rules. As such, UB injections of this kind often **compile and pass verifier checks**, but they still represent undefined behavior from the C standard’s perspective.

---

#### [5_01a_ptrcomp]

**Implementation Details:**
- The Ethernet header (`hdr->eth`) is **reinterpreted as an IPv4 header**:
  ```c
  struct iphdr *fake_ipv4 = (struct iphdr *)hdr->eth;
  __u8 ttl = fake_ipv4->ttl;
  ```

  This breaks strict aliasing since an `ethhdr` is not compatible with an `iphdr`.

  The `ttl` value is then compared against `DEFAULT_TTL` to decide whether to drop the packet.

This is undefined behavior because the effective type of the object (`struct ethhdr`) is accessed through an incompatible type (`struct iphdr`).

- **Verifier**: Passed.
- **Extra warnings**: None.
- **Exploitable**: Not possible in eBPF due to verifier-enforced memory bounds (cannot go beyondhdr->eth bound), but semantically invalid under C aliasing rules.

---

#### [5_01b_ptrcomp]

**Implementation Details:**
- The Ethernet header is **reinterpreted as a TCP header**:
  ```c
  struct tcphdr *fake_tcp = (struct tcphdr *)hdr->eth;
  __u16 fake_src_port = fake_tcp->source;
  ```

  The code then checks if the fake source port is `0` to potentially drop the packet.

  Since the memory layout of `struct ethhdr` and `struct tcphdr` are incompatible, accessing the Ethernet data as a TCP header violates the aliasing rule.

- **Verifier**: Passed.
- **Extra warnings**: None.
- **Exploitable**: Not possible in eBPF due to verifier-enforced memory bounds , but logically incorrect under the C standard.

---

#### [5_01c_ptrcomp]

**Implementation Details:**
- The Ethernet header is accessed as a raw `__u32` pointer:
  ```c
  __u32 *fake_eth = (__u32 *)hdr->eth;
  __u32 eth_value = *fake_eth;
  ```

  This violates strict aliasing rules because the memory originally declared as a `struct ethhdr` is accessed as a plain integer pointer.

  The retrieved value is printed with `bpf_printk`.

- **Verifier**: Passed.
- **Extra warnings**: None.
- **Exploitable**: Not possible in eBPF due to verifier-enforced memory bounds, though accessing structured header fields as raw integers is UB in standard C.

---

### Summary

All three UB injections compile and pass the eBPF verifier since they remain within memory bounds and do not trigger invalid pointer dereferencing. However, they are undefined behavior under ISO C, as they access objects through incompatible pointer types. The verifier does not diagnose this class of UB.

*Signed-by*: Gianfranco Trad

---

### [5.4 boolasgn]: No assignment in conditional expressions

Frequent mistake in C/C++ is typing `if (x = y)` instead of `if (x == y)`. The assignment expression `(x = y)` evaluates to the value assigned to `x`. If `y` is non-zero, the condition is always true, regardless of `x`'s initial value. This can lead to bugs where code branches are taken unexpectedly or loops become infinite.

#### [5_4a_boolasgn]

**Implementation Details:**
- The patch introduces a `while` loop with a direct assignment in its conditional: `while (processed_len = current_tcp_len) { ... }`. The purpose is to emulate a potential mistake in the while condition for processing a TCP header.
- `current_tcp_len` is derived from `hdr->tcp_len`, which for a valid TCP header, will always be a non-zero value (minimum 20 bytes).
- Inside the loop, `processed_len` is incremented. However, in the next iteration, `processed_len` is **re-assigned** the non-zero `current_tcp_len`, effectively resetting its value for the loop condition. Due to `current_tcp_len` always being a non-zero value, the condition `(processed_len = current_tcp_len)` will **always evaluate to true**.

The eBPF verifier, through its static analysis of register states and control flow, will correctly identify this `while` loop as an **infinite loop**. As a result, the eBPF verifier will reject the program load.

**Verifier:** Not passed (`infinite loop detected`).

**Extra warnings**:

```
xdp_synproxy_kern.c:631:27: warning: using the result of an assignment as a condition without parentheses [-Wparentheses]
  631 |         while (processed_len = current_tcp_len) {
      |                ~~~~~~~~~~~~~~^~~~~~~~~~~~~~~~~
xdp_synproxy_kern.c:631:27: note: place parentheses around the assignment to silence this warning
  631 |         while (processed_len = current_tcp_len) {
      |                              ^
      |                (                              )
xdp_synproxy_kern.c:631:27: note: use '==' to turn this assignment into an equality comparison
  631 |         while (processed_len = current_tcp_len) {
      |                              ^
      |                              ==
```

**Exploitable:** Not possible.

#### [5_4b_boolasgn]

**Implementation Details:**
- This version introduces a `do ... while` loop with a direct assignment in its conditional:
  `do { ... } while (processed_len = current_tcp_len);`.
- Unlike the `while` version, the `do ... while` construct guarantees that the loop body will execute **at least once**, even if `hdr->tcp_len` (`current_tcp_len`) is zero.
- Since `hdr->tcp_len` represents the TCP header length, it is normally non-zero (minimum 20 bytes). This means the condition `(processed_len = current_tcp_len)` always evaluates to true.

The eBPF verifier, through its static analysis of register states and control flow, will correctly identify this `do..while` loop as an **infinite loop**. As a result, the eBPF verifier will reject the program load.

**Verifier:** Not passed (`infinite loop detected`).

**Extra warnings**:

```
xdp_synproxy_kern.c:626:25: warning: using the result of an assignment as a condition without parentheses [-Wparentheses]
  626 |         } while (processed_len = current_tcp_len);
      |                  ~~~~~~~~~~~~~~^~~~~~~~~~~~~~~~~
xdp_synproxy_kern.c:626:25: note: place parentheses around the assignment to silence this warning
  626 |         } while (processed_len = current_tcp_len);
      |                                ^
      |                  (                              )
xdp_synproxy_kern.c:626:25: note: use '==' to turn this assignment into an equality comparison
  626 |         } while (processed_len = current_tcp_len);
      |                                ^
      |                                ==
```

**Exploitable:** Not possible.

*Signed-by: Francesco Rollo*

---

### [5.6 argcomp]: Calling a function with the wrong number or type of arguments

Calling functions with incorrect arguments, incompatible types, or mismatched prototypes leads to undefined behavior. This commonly occurs in multi-file projects where function declarations and definitions don't match.

#### [5_06a_argcomp]: Function pointer type incompatibility

**Implementation Details:**
- A function pointer `fp_wrong_type` is declared with signature `__u16 ()()` (no arguments, returns `__u16`).
- The pointer is assigned to `csum_fold`, which actually has signature `__u16 (__u32)` (takes one `__u32` argument).
- The function is called through the incompatible pointer with two arguments, creating type incompatibility.
- Demonstrates UB 26: "A pointer is used to call a function whose type is not compatible with the pointed-to type."
- In eBPF context, the BPF calling convention limits arguments to 5, but the type mismatch still violates C standards.

**Verifier:** Passed (compiles with warnings, but type incompatibility remains).

**Extra warnings:**
```
xdp_synproxy_kern.c:449:35: warning: passing arguments to a function without a prototype is deprecated in all versions of C and is not supported in C23 [-Wdeprecated-non-prototype]
```

**Exploitable:** **Yes**, potentially dangerous - Function pointer type incompatibility disrupts the calling convention, potentially corrupting registers or stack data when arguments are passed incorrectly. This can lead to unpredictable program behavior or memory corruption.

#### [5_06b_argcomp]: Wrong number of arguments

**Implementation Details:**
- Forward declaration `network_copy_helper()` is made without parameters to simulate separate compilation units.
- The function is called with 3 arguments (`src`, `dst`, `extra_buf`) but is actually defined to take only 2.
- Simulates the "separate source file" scenario where the caller doesn't have the correct prototype.
- Demonstrates UB 38: "For a call to a function without a function prototype in scope, the number of arguments does not equal the number of parameters."
- The extra argument is passed but ignored by the actual function, potentially causing stack corruption in other contexts.

**Verifier:** Passed (function call succeeds but with undefined argument handling).

**Extra warnings:**
```
xdp_synproxy_kern.c:444:21: warning: passing arguments to 'network_copy_helper' without a prototype is deprecated in all versions of C and is not supported in C23 [-Wdeprecated-non-prototype]
xdp_synproxy_kern.c:374:6: warning: a function declaration without a prototype is deprecated in all versions of C and is treated as a zero-parameter prototype in C23, conflicting with a subsequent definition [-Wdeprecated-non-prototype]
```

**Exploitable:** **Yes**, potentially dangerous - Passing extra arguments can overwrite adjacent stack memory since the function only expects two parameters. The third argument gets pushed onto the stack but has no designated storage, potentially corrupting nearby data structures.

#### [5_06c_argcomp]: Variadic function without prototype

**Implementation Details:**
- Forward declaration `debug_print_helper()` is made without parameters, hiding its variadic nature.
- The function is called with format string and arguments like printf (`"Port value: %u at line %d", port_val, __LINE__`).
- The actual function definition DOES support variable arguments, but the caller lacks the proper prototype.
- Demonstrates UB 39: "For a call to a function without a function prototype in scope where the function is defined with a function prototype, either the prototype ends with an ellipsis or the types of the arguments after promotion are not compatible."
- In eBPF, variadic functions are limited, but the principle of prototype mismatch applies.

**Verifier:** Not passed (compilation fails due to conflicting function prototypes).

**Extra warnings:**
```
xdp_synproxy_kern.c:445:20: warning: passing arguments to 'debug_print_helper' without a prototype is deprecated in all versions of C and is not supported in C23 [-Wdeprecated-non-prototype]
xdp_synproxy_kern.c:374:6: warning: a function declaration without a prototype is deprecated in all versions of C and is treated as a zero-parameter prototype in C23, conflicting with a subsequent definition [-Wdeprecated-non-prototype]
xdp_synproxy_kern.c:453:6: error: conflicting types for 'debug_print_helper'
```

**Exploitable:** **No** - Compilation fails due to conflicting function declarations, preventing the code from running. No runtime security risk exists since the program cannot be built.

#### [5_06d_argcomp]: Wrong argument types

**Implementation Details:**
- Forward declaration `helper_function()` is made without parameters to simulate cross-file compilation.
- The function is called with an `int` argument (`int_value`) but is actually defined to expect a `long`.
- Simulates cross-file scenarios where caller uses wrong type assumptions.
- Demonstrates UB 41: "A function is defined with a type that is not compatible with the type pointed to by the expression that denotes the called function."
- The type mismatch between `int` and `long` can cause incorrect value interpretation, especially on systems where they differ in size.

**Verifier:** Passed (function call succeeds but with potential data truncation/extension issues).

**Extra warnings:**
```
xdp_synproxy_kern.c:445:31: warning: passing arguments to 'helper_function' without a prototype is deprecated in all versions of C and is not supported in C23 [-Wdeprecated-non-prototype]
xdp_synproxy_kern.c:374:6: warning: a function declaration without a prototype is deprecated in all versions of C and is treated as a zero-parameter prototype in C23, conflicting with a subsequent definition [-Wdeprecated-non-prototype]
```

**Exploitable:** **Limited** - The int/long type mismatch causes value truncation or sign extension issues, but the impact remains localized to the affected variable without broader memory safety implications.

#### [5_06e_argcomp]: BPF helper function with incompatible argument types

**Implementation Details:**
- `bpf_map_lookup_elem()` expects `(void *map, const void *key)` with a valid map pointer.
- We pass an integer `0xDEADBEEF` casted to `void*` as the map argument instead of a real map pointer.
- This violates the BPF helper function contract and can crash the kernel if not caught by the verifier.
- Demonstrates UB 41: "A function is defined with a type that is not compatible with the type pointed to by the expression that denotes the called function."
- More dangerous than regular function calls because BPF helpers operate in kernel space where invalid pointers can cause immediate kernel panic.

**Verifier:** Failed (should reject program with invalid map pointer).

**Extra warnings:** None (only base warnings present - compiles successfully)

**Exploitable:** **No** - The eBPF verifier detects the invalid map pointer and rejects the program entirely. While this would be extremely dangerous in kernel space, the verification prevents execution.

*Signed-by: Giovanni Nicosia*

---

### [5.9 padcomp]: Comparison of padding data

C compilers often insert padding bytes into structures to ensure proper alignment of fields, especially when mixing data types of different sizes. The values of these padding bytes may contain arbitrary data from previous stack usage, or garbage.

If two instances of an identical-looking struct have their data fields set to the same values, their padding bytes might still differ. A `memcmp` between these two structs might result in a mismatch introducing non-deterministic behavior into the program.

**Implementation Details:**
- A new struct `padded_config_data` is defined with `__attribute__((aligned(8)))`. This attribute, combined with chosen different field types, forces the compiler to insert padding bytes for alignment.
- Two instances of this struct, `config1` and `config2`, are created within the `syncookie_xdp` function:
    - `config1` is **fully aggregate-initialized**. This often causes compilers to zero-initialize any padding bytes within `config1`.
    - `config2` is **declared without aggregate initialization**, and its member fields are then individually assigned values. This leaves its padding bytes *explicitly uninitialized*.
- A conditional `__builtin_memcmp(&config1, &config2, sizeof(struct padded_config_data))` checks if the structs (including padding) are identical. If they are an `XDP_DROP` is returned. This demonstrates the potential runtime impact: if the uninitialized padding bytes happen to align, an otherwise valid packet could be unexpectedly dropped, demonstrating the non-deterministic and dangerous consequences of violating `[padcomp]`.

**Verifier:** Passed.

**Extra warnings**: None.

**Exploitable**: The issue only affects stack padding bytes within the local `padded_config_data` struct. Even though memcmp may behave non-deterministically due to uninitialized padding, this does not expose arbitrary memory outside the eBPF program’s stack. The “exploitation” is limited to logic non-determinism inside the program (e.g., occasional XDP_DROP), not a security vulnerability.

**Note:** *Additional tests for this rule would produce the same result as long as the struct layout remains controlled and the padding bytes are left uninitialized. Writing more tests for this case would be redundant, as the behavior is deterministic with respect to how padding is treated by the compiler and initialization method.*

Signed-by: Francesco Rollo*

---

### [5.10 intptrconv]: Converting a pointer type to an integer type or integer type to a pointer type

Converting pointers to integers and back can lead to undefined behavior if the resulting pointer is incorrectly aligned, doesn't point to an entity of the referenced type, or creates invalid memory references. This is particularly dangerous in eBPF where pointer arithmetic is strictly controlled by the verifier.

#### [5_10a_intptrconv]: Naive pointer truncation approach (blocked by verifier)

**Implementation Details:**
- Demonstrates the straightforward but naive approach to pointer-to-integer conversion and arithmetic.
- Attempts direct conversion: `calculated_ptr = (void *)(unsigned long)(data_base_truncated + payload_offset)`.
- The verifier blocks this because: (1) it tracks `data_base_truncated` as "derived from pointer", (2) arithmetic is still considered "pointer arithmetic", (3) conversion requires bitwise operations prohibited on pointers.
- Shows why the naive approach fails and demonstrates the verifier's protective mechanisms.

**Verifier:** Failed (correctly blocks unsafe pointer arithmetic).

**Extra warnings:** None (only base warnings present - compiles successfully)

**Exploitable:** **No** - The verifier recognizes and blocks the pointer arithmetic pattern before program execution. This demonstrates that the security mechanisms effectively prevent this class of attack.

#### [5_10a_exploit_intptrconv]: Information disclosure via advanced pointer truncation bypass

**Implementation Details:**
- Comprehensive exploit research demonstrating real information disclosure through pointer truncation.
- Includes detailed analysis of verifier bypass mechanisms and memory layout confusion.
- Shows how reconstructed pointers can access unauthorized kernel memory or other packet data.
- Features extensive logging and diagnostic output to demonstrate the attack process step-by-step.
- Proves that bounds check bypass can lead to actual information leakage from unintended memory regions.

**Verifier:** Passed (successful information disclosure exploit through verifier bypass).

**Extra warnings:** None (only base warnings present)

**Exploitable:** **Yes** - Critical information disclosure vulnerability that bypasses verifier protections and can leak kernel memory contents.

#### [5_10b_intptrconv]: Hardcoded integer to pointer conversion

**Implementation Details:**
- Direct conversion of hardcoded integer `0xDEADBEEF` to pointer (`magic_ptr = (void *)MAGIC_ADDR`).
- Implements EXAMPLE 2 from ISO-IEC TS 17961-2013 standard.
- Demonstrates arbitrary address access where hardcoded addresses may point to sensitive memory regions.
- Attempts bounds-checked memory access with the invalid pointer.

**Verifier:** Failed (rejects program due to invalid pointer usage).

**Extra warnings:** None (only base warnings present - compiles successfully)

**Exploitable:** **No** - Converting hardcoded integers to pointers would allow arbitrary memory access in normal programs, but the eBPF verifier catches this pattern and blocks execution entirely.

*Signed-by: Giovanni Nicosia*

---

### [5.11 alignconv]: Converting pointer values to more strictly aligned pointer types

Converting a pointer value to a type that requires stricter alignment than the object actually provides is **undefined behavior** in C. According to MISRA C:2012 Rule 5.11 and ISO/IEC 9899:2011 §6.3.2.3, such conversions are invalid because the destination type may impose stricter alignment requirements than the source. If the underlying memory is not suitably aligned, dereferencing the pointer triggers undefined behavior.

In **eBPF**, the verifier enforces bounds and provenance checks but does not track alignment requirements at the C standard level. This means that UB injections of this form often **compile and pass verifier checks**. However, from the perspective of ISO C semantics, they remain undefined.

---

#### [5_11a_alignconv]

**Implementation Details:**
- A local buffer is intentionally misaligned:
  ```c
  char unaligned_buf[sizeof(struct iphdr) + 1];
  char *unaligned_ptr = &unaligned_buf[1]; // intentionally unaligned
  struct iphdr *misaligned_iph = (struct iphdr *)unaligned_ptr;
  __u8 dummy_version = misaligned_iph->version;
  ```

  The `unaligned_ptr` is offset by one byte, ensuring it is not aligned to `struct iphdr`’s natural boundary.

  Accessing `misaligned_iph->version` is UB due to stricter alignment requirements.

- **Verifier**: Passed.
- **Extra warnings**: None.
- **Exploitable**: Not memory exploitable in eBPF. The verifier still enforces bounds safety; the misalignment only causes logical misbehavior, not arbitrary memory access. But this could incorrectly satisfy or bypass control-flow conditions.

---

#### [5_11b_alignconv]

**Implementation Details:**
- The field `tcp_len` (declared as `__u16`) is accessed through a `__u8` pointer:
  ```c
  __u8 *tcp_len_as_u8 = (__u8 *)&hdr->tcp_len;
  __u8 fake_tcp_len = *tcp_len_as_u8;
  if (fake_tcp_len == 0)
      return XDP_DROP;
  ```

  This violates alignment and effective type rules by reinterpreting a `__u16` as a stricter `__u8`.

  Although logically meaningless, the code compiles and passes verifier checks.

- **Verifier**: Passed.
- **Extra warnings**: None.
- **Exploitable**: Not memory exploitable. However, this type of UB can be logically exploitable:
  - If program logic checks the full `__u16` but the UB cast inspects only the low byte, attackers can craft inputs where the LSBs match expected values (e.g., `0x0100 → 0x00` when truncated to `__u8`).
  - This could incorrectly satisfy or bypass control-flow conditions.

---

#### [5_11c_alignconv]

**Implementation Details:**
- The Ethernet header pointer (`hdr->eth`, aligned for `struct ethhdr`) is converted to a `__u64 *`:
  ```c
  __u64 *strictly_aligned_ptr = (__u64 *)hdr->eth;
  __u64 fake_eth_value = *strictly_aligned_ptr;
  bpf_printk("[alignconv]: Accessed fake_eth_value: %llu", fake_eth_value);
  ```

  Since `__u64` may impose stricter alignment requirements than `struct ethhdr`, this conversion is undefined behavior if the pointer is not suitably aligned.

- **Verifier**: Passed.
- **Extra warnings**: None.
- **Exploitable**: Not memory exploitable. Same as above, potential for logical misinterpretation if the misaligned cast influences how header values are checked.

---

### Summary

All three UB injections pass the verifier since the eBPF memory model only checks bounds and provenance, not alignment constraints. From the ISO C perspective, these are undefined behaviors because they convert pointers to more strictly aligned types than allowed.

**Exploitability considerations**:
- Not security-exploitable for memory corruption.
- Potentially logically exploitable:
  - Misaligned casts or narrowing reinterpretations (e.g., `__u16 → __u8`) can alter validation logic.
  - Attackers may craft packets where only the low-order bytes match expectations, bypassing checks that should fail.

*Signed-by*: Gianfranco Trad

---

### [5.13 objdec]: Declaring the same function or object in incompatible ways

Declaring the same function or object multiple times with **incompatible types** is **undefined behavior** in C. According to ISO/IEC 9899:2011 §6.2.7, two or more incompatible declarations of the same object or function in the same program must be diagnosed.

Undefined behavior arises when:
- An object is accessed through an lvalue of an incompatible type.
- A function is called through a pointer whose type is incompatible with the function’s actual definition.

In **eBPF**, this UB manifests at **compile time**, because the compiler enforces type consistency for objects and functions. The eBPF verifier never sees the code, as it cannot be loaded if compilation fails.

---

#### [5_13a_objdec]

**Implementation Details:**
- An object `fake_var` is declared with conflicting types:
  ```c
  extern int fake_var;
  short fake_var = 42;  // Conflicting type (int vs short)
  ```

  The compiler produces a type mismatch error (`-Wint-conversion` or similar).

- **Compilation**: Fails.
- **Extra warnings**: None.
- **Verifier**: Not reached.
- **Exploitable**: Not exploitable; code does not compile, so no runtime behavior occurs.

---

#### [5_13b_objdec]

**Implementation Details:**
- A function is declared with incompatible types:
  ```c
  extern int h(int a);
  long h(long a) { return a * 2; }
  ```

  The compiler rejects this due to incompatible function types.

- **Compilation**: Fails.
- **Extra warnings**: None.
- **Verifier**: Not reached.
- **Exploitable**: Not exploitable; compilation prevents runtime execution.

---

### Summary

UB from incompatible object or function declarations is caught at compile time, preventing the program from being loaded or executed.

*Signed-by*: Gianfranco Trad

---

### [5.14 nullref]: Dereferencing a possibly null or invalid pointer

Dereferencing pointers derived from potentially tainted input (e.g., packet headers) without validating them can result in undefined behavior, including invalid memory accesses or crashes.


#### [5_14_nullref]
**Implementation Details:**
- The helper function `null_copy_address()` is introduced to simulate a dereference on an input pointer.
- A pointer `copy_from` is initialized to `&hdr->ipv4->saddr`, which assumes that `hdr->ipv4` is valid.
- A buffer `input_string` of matching size (`sizeof(__be32)`) is allocated on the stack.
- The function attempts to `__builtin_memcpy` from `copy_from` into `input_string`.
- If `hdr->ipv4` is `NULL`, this dereference results in an invalid memory access, highlighting the risk of dereferencing unvalidated or tainted pointers in packet parsing logic.

**Extra warnings**: None.

**Verifier:** Passed (invalid dereference not detected by static analysis).

**Exploitable:** If an attacker can craft a packet that omits or corrupts the IPv4 header, `hdr->ipv4` may resolve to `NULL` or an invalid pointer. In such a case, the dereference of `hdr->ipv4->saddr` would trigger an invalid memory access, which in eBPF could lead to verifier bypasses being overlooked in other contexts, or in non-BPF C code could cause kernel crashes or privilege escalation through controlled faulting behavior.


#### [5_14a_nullref]

**Implementation Details:**
- The helper function `null_copy_address_v6()` is introduced to simulate dereferencing of an IPv6 field without checking whether `hdr->ipv6` is valid.
- A pointer `copy_from` is initialized to `&hdr->ipv6->daddr`, assuming `hdr->ipv6` is non-`NULL`.
- A stack buffer `sink` of size `sizeof(struct in6_addr)` is allocated.
- The function attempts to `__builtin_memcpy` from `copy_from` into `sink`.
- If the packet being processed is IPv4, then `hdr->ipv6` is `NULL`, and the dereference of `hdr->ipv6->daddr` results in an invalid access.

**Extra warnings**: None.

**Verifier:** Passed (no rejection by the BPF verifier). The verifier tracks packet parsing paths but does not detect that `hdr->ipv6` may be `NULL` here.

**Exploitable:** If an attacker can send IPv4 traffic, `hdr->ipv6` will be `NULL`, and the dereference inside `null_copy_address_v6()` leads to an invalid pointer access. No real escape BPF attack path.


#### [5_14b_nullref]

**Implementation Details:**
- The helper function `unsafe_values_peek()` is introduced to simulate dereferencing the result of a BPF map lookup without checking for `NULL`.
- A key (`__u32 key = 1234`) is chosen outside of the defined range of the map on purpose, so the lookup will usually fail.
- A pointer `value` is initialized to the return value of `bpf_map_lookup_elem(&values, &key)`.
- A stack buffer `buf` of matching size (`sizeof(__u64)`) is allocated.
- The function attempts to `__builtin_memcpy` from `value` into `buf` without verifying whether `value` is valid.
- If `bpf_map_lookup_elem()` returns `NULL` (no element present), this results in an invalid dereference.


**Verifier:** Not passed :
```
; __builtin_memcpy(buf, value, sizeof(__u64));
170: (71) r3 = *(u8 *)(r7 +0)
R7 invalid mem access 'map_value_or_null'
processed 122 insns (limit 1000000) max_states_per_insn 0 total_states 9 peak_states 9 mark_read 5
-- END PROG LOAD LOG --
```
**Observed Behavior:** The verifier rejects this code due to potential null pointer dereference from map lookup.

**Extra warnings**: None.

**Exploitable:** Not exploitable. The eBPF verifier correctly detects and prevents this null pointer dereference at load time.

*Signed-by: Giorgio Fardo*

---

### [5.15 addrescape]: Escaping the address of an automatic object

Automatic (stack-allocated) variables exist only for the lifetime of the function in which they are defined. Returning or storing their address beyond that lifetime results in undefined behavior, as the memory may be overwritten or invalidated.

#### [5_15_addrescape]

**Implementation Details:**
- The function `set_pointer()` defines a local string `char str[] = "TEst1"`.
- It assigns the address of this local string to a pointer argument `*ptr_param`.
- After `set_pointer()` returns, the pointer `ptr` in the caller still references `str`, which is now out of scope.
- A call to `bpf_printk("Res: %s", ptr)` prints garbage or nothing, as the pointer references invalid memory.
- Demonstrates how escaping stack addresses can result in use-after-scope bugs and potential memory corruption.

**Extra warnings**: None.

**Verifier:** Passed (stack lifetime violations are not detected).

**Exploitable:** Not really — while this results in a dangling pointer, in eBPF the stack frame is strictly managed and reallocated per packet. The pointer cannot outlive the helper call context, so an attacker cannot reliably control or reuse the memory for malicious purposes beyond producing garbage logs.

#### [5_15a_addrescape]

**Implementation Details:**
- The function `init_array()` defines a local array `int array[5] = { 1, 2, 3, 4, 5 }`.
- It incorrectly returns the address of this local array.
- In `syncookie_part1()`, the pointer `arr_ptr` receives this address and is later dereferenced in a `bpf_printk` call.
- After `init_array()` returns, `array` is out of scope, so `arr_ptr` references invalid memory.
- The program still logs `"Leaked array[0]: 1"`, but this is undefined behavior and may produce garbage in other contexts.
- Demonstrates how escaping stack addresses can result in dangling pointers and potential use-after-scope bugs.

**Verifier:** Passed (stack lifetime violations are not detected).

**Extra Warnings:**
```
xdp_synproxy_kern.c:761:9: warning: address of stack memory associated with local variable 'array' returned [-Wreturn-stack-address]
761 | return array; // diagnostic required
| ^~~~~
6 warnings generated.
```

**Exploitable:** Not really — while this results in a dangling pointer, in eBPF the stack frame is strictly managed and reallocated per packet. The pointer cannot outlive the helper call context, so an attacker cannot reliably control or reuse the memory for malicious purposes beyond producing garbage logs.


#### [5_15b_addrescape]

**Implementation Details:**
- The helper function `squirrel_away()` defines a local string `char fmt[] = "Error: %s\n"`.
- It stores the address of this local array into a pointer argument `*ptr_param`.
- In `syncookie_part1()`, the caller receives this escaped pointer into `fmt_ptr`.
- After `squirrel_away()` returns, `fmt` is out of scope, so `fmt_ptr` references invalid memory.
- A call to `bpf_printk("Escaped fmt string: %s", fmt_ptr)` may appear to work, but the pointer is dangling and the behavior is undefined.
- This demonstrates how stack addresses can escape through function parameters, leading to use-after-scope bugs.

**Verifier:** Passed (stack lifetime violations are not detected).

**Extra warnings**: None.

**Exploitable:** Not really — as with the first case, although this creates a dangling pointer, in eBPF the stack frame is strictly managed and reset per packet. The pointer cannot persist across contexts, so attackers cannot exploit it beyond producing garbage or misleading logs.

*Signed-off-by: Giorgio Fardo*

---

### [5.16 signconv]: Converting a tainted value of type char or signed char to a larger integer type without first casting to unsigned char

When `char` is signed (implementation-defined), converting directly to `int` without first casting to `unsigned char` can cause 0xFF bytes to be sign-extended to -1 (EOF), leading to false positives in EOF checks.

#### [5_16a_signconv]: Raw TCP payload access with signed conversion (unsafe memory access)

**Implementation Details:**
- Directly accesses TCP payload data (`char *tcp_payload = (char *)hdr->tcp + (hdr->tcp->doff * 4)`) without proper bounds checking.
- Performs unsafe signed char to int conversion (`int c = raw_char`) where 0xFF becomes -1 instead of 255.
- The problematic EOF comparison (`if (c == EOF)`) triggers false positives on legitimate 0xFF bytes.
- This version is expected to fail verifier due to unbounded memory access of tainted network data.

**Verifier:** Not passed (unsafe memory access).

**Extra warnings:** None (only base warnings present - compiles successfully)

**Exploitable:** **No** - While the signed conversion vulnerability is conceptually valid, the verifier prevents unsafe access to TCP payload data, blocking the attack vector before it can cause memory corruption.

#### [5_16b_signconv]: Controlled demonstration of signed char conversion vulnerability

**Implementation Details:**
- Uses controlled test data (`char test_data[4] = {0x41, 0x42, 0xFF, 0x44}`) to demonstrate the same vulnerability while passing verifier checks.
- Shows how 0xFF bytes in legitimate data get confused with EOF due to sign extension.
- Demonstrates the core vulnerability in a verifier-compatible way while maintaining the security implications.

**Verifier:** Passed (controlled demonstration).

**Extra warnings:** None (only base warnings present)

**Exploitable:** **Limited** - This controlled demonstration shows how signed conversion can cause logic errors (false EOF detection), but the impact is confined to program logic rather than memory safety violations.

*Signed-by: Giovanni Nicosia*

---

### [5.17 swtchdflt]: Switch statement missing default case or incomplete enumeration coverage

A switch statement with an enumerated controlling expression that lacks a default case and doesn't handle all enumeration constants can lead to undefined behavior when unhandled values are encountered.

#### [5_17_swtchdflt]: Switch statement missing default case for firewall actions

**Implementation Details:**
- Defines a `firewall_action` enum with four values: `FIREWALL_ALLOW`, `FIREWALL_BLOCK`, `FIREWALL_REDIRECT`, and `FIREWALL_LOG`.
- Integrates firewall classification into the `tcp_dissect` function based on destination port ranges, making the vulnerability realistic in network processing context.
- The switch statement handles only 3 out of 4 enum values (missing `FIREWALL_REDIRECT` case) and lacks a default case.
- When packets are classified with `FIREWALL_REDIRECT` action (ports 8000-8999), execution falls through with undefined behavior, potentially returning garbage values that could be interpreted as `XDP_ABORTED`.
- Demonstrates how incomplete switch coverage can lead to security policy violations in real network filtering scenarios.

**Verifier:** Passed (but causes undefined behavior on missing cases).

**Extra warnings:**
```
xdp_synproxy_kern.c:509:10: warning: enumeration value 'FIREWALL_REDIRECT' not handled in switch [-Wswitch]
```

**Exploitable:** **Yes**, potentially dangerous - The missing FIREWALL_REDIRECT case creates undefined behavior when packets trigger that action. This can result in arbitrary return values, potentially causing legitimate traffic to be dropped or malicious traffic to pass through.

*Signed-by: Giovanni Nicosia*

---

### [5.20 libptr]: Forming invalid pointers by library function

Invoking a function with arguments that cause it to form pointers that do not point into or just past the end of an object violates Rule 5.20. While eBPF doesn't have access to standard C library functions, it still uses memory manipulation functions like `__builtin_memcpy` and BPF helpers that can form invalid pointers through incorrect size calculations.

#### [5_20a_libptr]: Buffer overflow with __builtin_memcpy oversized copy

**Implementation Details:**
- Allocates a 16-byte buffer but attempts to copy 24 bytes using `__builtin_memcpy`
- The library function forms pointers beyond the allocated object bounds when performing the oversized copy
- Includes TCP options extraction scenario where options can be up to 40 bytes but buffer is only 16 bytes
- Shows how size calculation errors lead to buffer overruns that violate pointer validity rules
- Demonstrates the security implications when fixed-size copy operations exceed buffer boundaries

**Verifier:** **Passed** - The eBPF verifier does not detect this buffer overflow, allowing the violation to execute

**Extra warnings:**
```
xdp_synproxy_kern.c:469:25: warning: comparison of distinct pointer types ('char *' and 'void *') [-Wcompare-distinct-pointer-types]
xdp_synproxy_kern.c:476:4: warning: 'memcpy' will always overflow; destination buffer has size 16, but size argument is 24 [-Wfortify-source]
```

**Exploitable:** **Yes**, potentially dangerous - Buffer overflow of 8 bytes can corrupt stack variables adjacent to the buffer, potentially causing program crashes or memory corruption

#### [5_20b_libptr]: BPF helper with invalid size parameters (verifier rejected)

**Implementation Details:**
- Uses `bpf_probe_read_kernel` with size parameter (32 bytes) larger than destination buffer (12 bytes)
- The BPF helper attempts to form invalid pointers when accessing beyond the actual buffer boundaries
- Demonstrates how helper function parameter mismatches can lead to out-of-bounds memory access
- Shows realistic scenarios where helper size parameters don't match buffer expectations
- Based on Rule 5.20 subpoint 5.20.1: functions taking (pointer, size_bytes) parameters

**Verifier:** **Rejected** - The eBPF verifier blocks this pattern, detecting the invalid buffer size parameter

**Extra warnings:** None (only base warnings present - compiles successfully)

**Exploitable:** **No** - The verifier prevents this violation from executing, demonstrating better protection for BPF helpers compared to `__builtin_memcpy`

#### [5_20c_libptr]: Type confusion in size calculations causing buffer overflow

**Implementation Details:**
- Performs size calculations using wrong data types (e.g., `sizeof(int) * 8 = 32` instead of `sizeof(char) * 8 = 8`)
- Uses `__builtin_memcpy` with miscalculated sizes that exceed buffer boundaries (32 bytes into 20-byte buffer)
- Shows how type assumption errors lead to incorrect size calculations in packet processing
- Demonstrates pointer formation violations where type confusion causes out-of-bounds access
- Based on Rule 5.20 Example 2: sizeof(int) vs sizeof(float) type confusion

**Verifier:** **Passed** - The verifier does not detect type confusion in size calculations, allowing the buffer overflow

**Extra warnings:**
```
xdp_synproxy_kern.c:476:33: warning: comparison of distinct pointer types ('char *' and 'void *') [-Wcompare-distinct-pointer-types]
```

**Exploitable:** **Yes**, potentially dangerous - Type confusion causing 12-byte buffer overflow may corrupt adjacent stack memory, leading to program instability or information leakage

*Signed-by: Giovanni Nicosia*

---

### [5.22 invptr]: Using out-of-bounds pointers or array subscripts

Pointer arithmetic or array indexing that goes beyond the bounds of an object is **undefined behavior** in C (ISO/IEC 9899:2011 §6.5.6). This includes:

- Addition or subtraction of pointers that result in addresses outside the same object or just past its end.
- Dereferencing pointers outside the valid object bounds.
- Array subscripts that access elements outside the declared array size.
- Accessing flexible array members when no elements exist.

In eBPF, the verifier enforces **memory safety** (bounds checking) for packet and map memory, but does not prevent all forms of logical OOB pointer manipulations if they remain within verifier-allowed memory regions. UB injections of this type may compile and pass the verifier, but still represent undefined behavior in C.

---

#### [5_22a_invptr]

**Implementation Details:**
- Using a negative offset to access a map element:
  ```c
  int ub_offset = -MAX_ALLOWED_PORTS;
  __u16 *value = bpf_map_lookup_elem(&allowed_ports, &ub_offset);  // UB
  __u16 ub_trigger = *value;  // optional dereference
  ```

  This forms an out-of-bounds pointer relative to the map key space.

- **Compilation**: Passed.
- **Extra warnings**: None.
- **Verifier**: Passed.
- **Exploitable**: Not exploitable in eBPF due to map key checking; UB is logical.

---

#### [5_22b_invptr]

**Implementation Details:**
- Pointer arithmetic beyond a TCP header field:
  ```c
  __u16 *tcp_ports = (__u16 *)&hdr->tcp->source;
  __u16 *invalid_ptr = tcp_ports + 10;
  bpf_printk("[invptr-1]: Invalid pointer value: %p", invalid_ptr);
  ```

  This produces a pointer outside the allocated TCP structure.

- **Compilation**: Passed.
- **Extra warnings**: None.
- **Verifier**: Passed.
- **Exploitable**: Not exploitable for memory corruption; may produce misleading logical values if used in calculations.

---

#### [5_22c_invptr]

**Implementation Details:**
- Dereferencing an out-of-bounds pointer:
  ```c
  __u16 *tcp_ports = (__u16 *)&hdr->tcp->source;
  __u16 invalid_value = *(tcp_ports + 10);
  bpf_printk("[invptr-2]: Invalid value: %u", invalid_value);
  ```

  Access is UB because it points past the structure.

- **Compilation**: Passed.
- **Extra warnings**: None.
- **Verifier**: Not passed.
```
Invalid access to context parameter
```

- **Exploitable**: Only logically exploitable; no memory corruption possible in eBPF due to verifier bounds.

---

#### [5_22d_invptr]

**Implementation Details:**
- Out-of-bounds array indexing:
  ```c
  __u16 tcp_array_info[3] = {hdr->tcp->source, hdr->tcp->dest, hdr->tcp_len};
  __u16 out_of_bounds_value = tcp_array_info[5];  // UB
  bpf_printk("[invptr-3]: Out-of-bounds array value: %u", out_of_bounds_value);
  ```

  This patch does not compile, as the compiler detects the OOB array access.

- **Compilation**: Fails.
- **Extra warnings**: None.
- **Verifier**: Not reached.
- **Exploitable**: Not applicable.

---

#### [5_22e_invptr]

**Implementation Details:**
- Pointer just past the end of TCP header:
  ```c
  __u8 *tcp_header_end = (__u8 *)hdr->tcp + hdr->tcp_len;
  __u8 *invalid_access = tcp_header_end + 1;
  bpf_printk("[invptr-4]: Invalid access pointer: %p", invalid_access);
  ```

  UB occurs when using `invalid_access`, pointing outside the object.

- **Compilation**: Passed.
- **Extra warnings**: None.
- **Verifier**: Passed.
- **Exploitable**: Not exploitable for memory corruption; could mislead logic if low-byte checks are used.

---

#### [5_22f_invptr]

**Implementation Details:**
- Flexible array member access with no elements:
  ```c
  struct {
      __u16 len;
      __u8 data[];
  } *flexible_struct = (__u8 *)hdr->tcp;
  __u8 invalid_flex_access = flexible_struct->data[0]; // UB
  bpf_printk("[invptr-5]: Invalid flexible array access: %u", invalid_flex_access);
  ```

  Accessing a non-existent element is undefined.

- **Compilation**: Passed.
- **Extra warnings**: None.
- **Verifier**: Passed.
- **Exploitable**: Only logically exploitable; memory corruption prevented by verifier bounds.

---

### Summary

- **Compiling and verifier behavior**: All patches except 5.22d compile and pass the eBPF verifier.
- **Memory safety**: eBPF verifier prevents out-of-bounds memory access, so these UB injections cannot corrupt memory.
- **Logical exploitability**: If the code performs partial checks (e.g., checking only LSBs of truncated values), OOB accesses could lead to unexpected logical results, potentially bypassing validation.

*Signed-by*: Gianfranco Trad

---

### [5.24 usrfmt]: Including tainted or out-of-domain input in a format string

Using **tainted or unvalidated input** in a format string for formatted I/O functions (e.g., `printf`, `vfprintf`, `bpf_printk`) is **undefined behavior** in C (ISO/IEC 9899:2011 §7.21.6). This can lead to:

- Crashes or segmentation faults.
- Reading unintended memory (stack or heap).
- Writing to arbitrary memory locations (e.g., via `%n` specifier).
- Potential arbitrary code execution if an attacker controls part of the format string.

**Notes:**

- An empty string is not considered tainted.
- Any comparison of a character to a value other than null may sanitize the string, but full control over the format string remains dangerous.

In eBPF, the verifier ensures memory safety but does **not validate format string contents**. UB injections with tainted format strings may compile and pass the verifier, but still represent undefined behavior in C.

---

#### [5_24a_usrfmt]

**Implementation Details:**
- A format string is retrieved from a BPF map:
  ```c
  __u32 key = hdr->ipv4->daddr;
  char *tainted_fmt = bpf_map_lookup_elem(&values, &key); // Potentially user-controlled
  if (tainted_fmt) {
      bpf_printk(tainted_fmt); // UB: tainted format string
  }
  ```

  Using a user-controlled string directly as a format argument is UB.

- **Compilation**: Fails.
  - **Reason**: `bpf_map_lookup_elem` returns `void *` in eBPF C, which cannot be implicitly converted to `char *` without a cast. Strict type rules in kernel BPF programs prevent direct compilation.
- **Extra warnings**: None.
- **Verifier**: Not reached due to compilation failure.
- **Exploitable**: If compiled with proper casting, this could be logically exploitable, e.g., a `%n` specifier could allow writing to arbitrary memory locations in standard C. In eBPF, memory safety prevents actual memory corruption, but logic or information leakage could occur.

---

#### [5_24b_usrfmt]

**Implementation Details:**
- Tainted input derived from the TCP header is inserted into a format string:
  ```c
  char tainted_input[16];
  __builtin_memcpy(tainted_input, &hdr->tcp->source, sizeof(hdr->tcp->source));
  tainted_input[sizeof(hdr->tcp->source)] = '\0';
  bpf_printk("Tainted input: %s", tainted_input); // UB
  ```

  The input could be partially controlled by an attacker (e.g., through network packet data). Using it in a formatted string is UB.

- **Compilation**: Passed.
- **Extra warnings**: None.
- **Verifier**: Not passed.
```
Invalid access to context parameter
```

- **Exploitable**: Logic-level exploit possible.
  - **Example**: If subsequent code parses the string or assumes format compliance, malformed input could bypass checks or corrupt logical processing.
  - Memory corruption is not possible due to eBPF verifier bounds checking.

---

### Summary

- **Patch 5.24a**: Does not compile due to strict pointer type mismatch from `bpf_map_lookup_elem`.
- **Patch 5.24b**: Compiles and passes verifier; demonstrates UB via tainted input in format string.
- **Exploitable scenarios**: Logical or information leakage; memory corruption prevented by eBPF verifier.

*Signed-by*: Gianfranco Trad

---

### [5.26 diverr]: Integer division errors

In standard C, division by zero and modulo by zero result in **undefined behavior**. This means the compiler is not required to handle such cases predictably.

The eBPF verifier is extremely **strict** about preventing undefined behavior and ensuring program safety. It performs static analysis to determine the possible range of values for any register that might be used as a divisor.

#### [5_26a_diverr] [5_26b_diverr] [5_26c_diverr]
**Implementation Details:**
- The value of `hdr->tcp->ack_seq` is extracted from the incoming TCP header and assigned to `tainted_divisor_val`. `ack_seq` is chosen because, in a TCP header, it can legitimately carry a value of zero, making it a suitable "tainted" variable that could cause a division-by-zero.
- A constant `numerator_val` is defined as `100`.
- Two operations are performed: `numerator_val / tainted_divisor_val` and `numerator_val % tainted_divisor_val;`. Since `tainted_divisor_val` can be zero, these operations violate the `[diverr]` rule.

Why this approach? If the verifier statically determines that a divisor *could* evaluate to zero during program execution, it will prevent the program from loading. Using a a value not known at compile time, might stop the Verifier from preventing the load.

***Unfortunately the eBPF runtime environment defines specific behavior for division/modulo by zero by setting the destination register to zero.***

**NB**: The same example is split across two different test files (example 2 and 3) reproducing the same operations in two different files for completeness.

**Verifier:** Passed (but not an issue at runtime).

**Extra warnings**: None.

**Explotable:** Not possible.

#### [5_26d_diverr]

**Implementation Details:**

- UB Case: Division result not representable in two’s complement arithmetic.
  Specifically: `INT_MIN / -1 (0x80000000 / -1)` overflows because `-INT_MIN`
  cannot be represented in a 32-bit signed integer.

- Implementation:
  - `tainted_dividend` is set arbitrarily to **INT_MIN** (0x80000000).
  - `divisor_neg_one` is set to -1.
  - Division is performed: `result_overflow = tainted_dividend / divisor_neg_one`.
  - `bpf_printk` prints the dividend, divisor, and result to force evaluation.

Standard C does not define the behavior for `INT_MIN / -1`. The result is theoretically non-representable, but memory safety is not violated.

**Verifier:** Passed.

**Extra warnings**: None.

**Exploitable**: Not possible. The eBPF runtime returns 0 in practice, and the operation cannot be leveraged for arbitrary memory access.

#### [5_26e_diverr]

**Implementation Details:**

- UB Case: Modulo result not representable in two’s complement arithmetic.
  Specifically: `INT_MIN % -1 (0x80000000 % -1)` is undefined in standard C.

- Implementation:
  - `tainted_dividend` is set arbitrarily to **INT_MIN** (0x80000000).
  - `divisor_neg_one` is set to -1.
  - Modulo operation is performed: `result_overflow = tainted_dividend % divisor_neg_one`.
  - `bpf_printk` prints the dividend, divisor, and result to force evaluation.

In C, the modulo requires the result to fit within the signed integer range. For `INT_MIN % -1`, the standard does not define a value.

**Verifier:** Passed

**Extra warnings**: None.

**Exploitabile**: Not possible. The eBPF runtime produces 0 as the result, so this cannot be used to read or write memory.

*Signed-by: Francesco Rollo*

---

### [5.28 strmod]: Modifying string literals

String literals in C are stored in read-only memory. Attempting to modify them results in undefined behavior, typically leading to a segmentation fault or silent failure.

**Implementation Details:**
- The helper function `setStringIndex()` assigns a string literal to a `char *` pointer: `char *str_literal = "This is a string literal";`.
- It then attempts to modify the literal with `str_literal[loc / 100000000] = 'A';`, using a computed offset.
- This operation is undefined behavior: string literals must not be written to.
- Despite the violation, the verifier allows the code to pass since it does not track mutability of string literal memory.

**Extra warnings**: None.

**Verifier:** Passed (modification of string literals not checked).

**Exploitable:** Not really — attempts to modify read-only memory holding literals will do nothing in this case. In eBPF, this results in program terminatio rather than memory corruption, so it cannot be weaponized by an attacker.

No extra tests as are pointless.

*Signed-by: Giorgio Fardo*

---

### [5.30 intoflow]: Overflowing signed integers

Integer overflow of signed types is undefined behavior in C. While unsigned integer overflow is well-defined, signed overflow can result in unpredictable behavior, especially if optimized away or miscompiled.

**Implementation Details:**
- The helper function `checkOverflow()` takes an `int` value and adds a large constant: `int result = value + 2147483647;`.
- When `value` is positive, the addition causes a signed integer overflow.
- The tainted value passed in is `bpf_htons(hdr->tcp->seq)`, which typically holds large values.
- Overflows and underflows are managed with wrap so they are ignored by the verifier

**Extra warnings**: None.

**Verifier:** Passed, wrap used.

**Exploitable:** Signed integer overflow in eBPF is not exploitable in practice, since the verifier tracks scalar ranges and arithmetic is defined modulo two’s complement in the JITed code path. At most, it causes incorrect logic branches (e.g., treating a valid sequence number as negative), but does not yield memory safety violations.

No extra tests as **wrap** is always  there no out of bound memory access are possible using overload as overflow is not possible.

*Signed-by: Giorgio Fardo*

---

### [5.31 nonnullcs]: Passing a non-null-terminated character sequence to a library function that expects a string

A C string is a sequence of characters terminated by a null character (`\0`). For example, when `bpf_printk` is given a `%s` format specifier, it reads bytes sequentially from the provided pointer until it encounters a null character.

If a character array passed to such a function is *not* null-terminated within its allocated bounds, the function will continue reading past the end of the intended buffer.This constitutes an **out-of-bounds read**, leading to UB and potentially sensitive information leakage.

#### [5_31a_nonnullcs]

**Implementation Details:**
- A struct `test_memory_layout` is declared on the stack within `syncookie_handle_syn`. This struct is specifically designed to control the memory layout, ensuring the data is stored contiguously. This is important because the Verifier may place guards between individual stack variables.

- The struct contains:
    -  `char string_buffer[8]`: An 8-byte character array intended to act as a non-null-terminated string.
    - `__u64 filler_data_1`, `__u32 filler_data_2`, `char padding_byte`: These fields are placed immediately after `string_buffer` and are explicitly initialized with known non-zero valus. This allows to clearly observe what bytes are read if the string is not null-terminated.
- A `memset` is used to initially zero out the entire struct.
- A `#pragma unroll` loop fills `string_buffer` with 'A' through 'H', **deliberately omitting the null terminator**.
- The core violation is demonstrated when the non-null-terminated `string_buffer` is passed to `bpf_printk` with the `%s` format specifier. `bpf_printk` will attempt to read bytes from `test_memory_layout.string_buffer` until it encounters a null byte. In this controlled scenario, it would eventually read into the `filler_data` and potentially beyond until a zero byte is found.

**Verifier**: Passed (Under controlled memory layout).

**Extra warnings**: None.

**Exploitable**: Not really in practice. It depends on what type of information is disclosed in the controlled memory layout.

#### [5_31b_nonnullcs]

**Implementation Details**
- A small character array (e.g., char bad_str[3] = {'a', 'b', 'c'}) is allocated on the stack inside `syncookie_handle_syn`.
- The array is deliberately not null-terminated.
- When `bpf_printk` of the array is called, the `%s` specifier causes `bpf_printk` to read memory sequentially until it encounters a `\0` byte.
- Since no terminator exists in the array, bpf_printk should keep reading into adjacent stack memory until a zero byte happens to be found.

**Verifier**: Passed.

**Extra warnings**: None.

**Exploitable**: Not really in practice. At worst, you might accidentally log adjacent stack contents, which is a form of information disclosure but is limited to what the eBPF program itself already has access to.
The Verifier prevents the program to read adjacent memory content as it always zero initialize each stack frame.

*Signed-by: Francesco Rollo*

---

### [5.33 restrict]: Passing pointers into the same object as arguments to different restrict-qualified parameters

The `restrict` keyword is a promise to the compiler that a pointer is the sole means of accessing a particular memory region for the duration of its scope. If two `restrict`-qualified pointers are used to access **overlapping** memory, or if a `restrict` pointer **aliases** with another pointer that modifies the same memory, this promise is then violated.

When the `restrict` rule is broken, the compiler is free to perform aggressive optimizations based on the false assumption of non-aliasing. This can lead to **undefined and dangerous runtime behavior**, such as **superseded data reads**. The compiler might cache a value from memory and then reuse that cached value even after the memory has been modified by an aliasing `restrict` pointer, leading to incorrect program logic.

#### [5_33a_restrict]

**Implementation Details:**
- A new helper function, `simulate_restrict_ub` takes two `char *restrict` pointers (`read_ptr_base`, `write_ptr_base`) and a new value.

- Inside `syncookie_handle_syn`, `simulate_restrict_ub` is called with `(char *)&hdr->tcp->seq` passed as *both* `read_ptr_base` and `write_ptr_base`. This is the direct violation of the `restrict` contract, as two `restrict` pointers are made to alias (point to the same memory location).

- The **UB** is then triggered:
    -  `simulate_restrict_ub` first reads the original value of `hdr->tcp->seq` (`initial_read_host_val`).
    -  It then explicitly writes a new value (the `cookie`) to `hdr->tcp->seq` via `write_ptr_base`.
    -  It then attempts to read the value from `hdr->tcp->seq` *again* (`final_read_host_val`) using `read_ptr_base`.
    -  If the compiler, due to `restrict` optimization, reuses the cached `initial_read_host_val` instead of performing a fresh memory read after the write, then `initial_read_host_val` will equal `final_read_host_val` even though the memory was just modified. This equality is the symptom of the UB.

The patch links this UB symptom to a direct change in the program's flow. Even though the tcp->seq of the generated SYN-ACK will be the same, the now changed value will
affect the SYN-ACK's ack_seq value, forcing the client to retry the connection.
By targeting `hdr->tcp->seq`, a legitimate SYN packet, which should have led to a SYN-ACK, is instead dropped due to an unpredictable internal state caused by the `restrict` violation.

**Verifier**: Passed.

**Extra warnings**: None.

**Exploitable**: Not in a security sense. Only causes **logic/data corruption** in local eBPF stack memory, since the memory is fully controlled by the program.

#### [5_33b_restrict]

**Implementation Details**
- A local stack buffer `tcp_options[16]` is declared inside `syncookie_handle_syn`.
- Two pointers are defined:
  - `ptr1` points to the start of the buffer.
  - `ptr2` points 4 bytes into the buffer, overlapping with `ptr1`.
- The noncompliant operation uses `__builtin_memcpy(ptr2, ptr1, 8)`:
  - This copies 8 bytes from `ptr1` to `ptr2`, causing the source and destination
    regions to overlap.
  - Because `memcpy` semantics assume `restrict` pointers do not alias, this violates
    the restrict contract.
  - The compiler may optimize under the assumption of non-overlapping pointers,
    leading to undefined behavior and possible corruption of the copied data.

In the eBPF context only local stack memory is affected, so no arbitrary memory read/write occurs. The corruption is limited to the `tcp_options` buffer used for demonstration.

**Verifier:** Passed.

**Extra warnings**: None.

**Exploitable:** Not in a security sense. It only causes local stack data corruption.

*Signed-by: Francesco Rollo*

---

### [5.35 unint_mem]: Referencing uninitialized memory

Using uninitialized memory results in undefined behavior. It can expose garbage values, leak data, or corrupt program logic depending on the compiler and runtime context.

#### [5_35_unint_mem]

**Implementation Details:**
- The function `uninitializedRead()` declares a pointer `char *uninit_ptr` without initializing it.
- It then uses `__builtin_memcpy(buf, uninit_ptr, 64);`, reading from an uninitialized memory location.
- The copied content (`buf`) is printed byte-by-byte with `bpf_printk()`, demonstrating random or garbage values.
- This shows how lack of initialization can lead to unpredictable outcomes and violate memory safety.

**Extra warnings**: None.

**Verifier:** Passed.

**Exploitable:** If the stack slot is not zeroed, uninitialized reads may leak kernel stack data to user space via `bpf_printk`, providing attackers with information disclosure. If zero-initialization happens at runtime, it reduces to benign behavior, but where disclosure occurs, it could aid in bypassing ASLR or building further attacks.

#### [5_35a_unint_mem]

**Implementation Details:**
- The function `uninitialized_auto_var_read()` declares an integer `int uninit_int;` without initializing it.
- It immediately checks `if (uninit_int == 0)` and logs via `bpf_printk()` whether the "uninitialized value" appeared as zero or not.
- Because `uninit_int` is not given a defined value, its contents come directly from the kernel stack, making the comparison unpredictable.
- The call to `uninitialized_auto_var_read()` was inserted in `syncookie_part1()`, ensuring that the function is exercised whenever this path is triggered.

**Warnings:**
```
xdp_synproxy_kern.c:762:9: warning: variable 'uninit_int' is uninitialized when used here [-Wuninitialized]
  762 |     if (uninit_int == 0) { // Reading indeterminate value
      |         ^~~~~~~~~~
xdp_synproxy_kern.c:760:19: note: initialize the variable 'uninit_int' to silence this warning
  760 |     int uninit_int; // uninitialized automatic variable
      |                   ^
      |                    = 0

```

**Verifier:** Passed.

**Exploitable:** This pattern risks leaking kernel stack data to user space via `bpf_printk()`, depending on how the verifier and runtime handle uninitialized stack slots. In contexts where stack slots are not cleared, this can expose sensitive information, potentially aiding exploitation strategies such as ASLR bypass or kernel memory disclosure. If the compiler or runtime zero-initializes the stack, the behavior reduces to a benign but misleading test case.

#### [5_35b_unint_mem]

**Implementation Details:**
- The helper `uninitialized_packet_read()` attempts to grow the packet buffer using `bpf_xdp_adjust_tail(ctx, add_len)`.
- If successful, a new tail region is exposed but not initialized by the kernel.
- The function then iterates over this extended area, reading each byte and logging its content with `bpf_printk()`.
- This simulates a scenario where uninitialized packet data could be observed, potentially leaking sensitive information or introducing nondeterministic behavior.

**Extra warnings**: None.

**Verifier:** Failed with `"R1 invalid mem access 'scalar'"`.

**Exploitable:** Not exploitable in this state.

*Signed-off: Giorgio Fardo*

---

### [5.36 ptrobj]: Subtracting or comparing pointers from different array objects

Subtracting or relationally comparing pointers that don't refer to the same array object results in undefined behavior. This commonly occurs when accidentally mixing pointers from different memory regions in packet processing scenarios.

#### [5_36a_ptrobj]: Local buffer vs packet data pointer comparison

**Implementation Details:**
- Creates two distinct objects: packet data from the network (Object 1) and a local stack buffer (Object 2).
- Demonstrates comparison violations (`if ((char *)hdr->eth != local_buffer)`) between completely separate memory regions.
- Performs pointer subtraction between different objects (`ptrdiff_t wrong_distance = (char *)hdr->tcp - local_buffer`), producing meaningless results.
- Shows relational comparisons (`>`, `<`) that have no defined meaning since pointers reference unrelated memory regions.
- Uses undefined results in program logic, demonstrating how violations lead to unpredictable behavior.

**Verifier:** Passed (but produces undefined results that may leak memory layout information).

**Extra warnings:** None (only base warnings present)

**Exploitable:** **Limited** - Comparing pointers from different objects produces undefined results that may leak memory layout information through comparison outcomes, but cannot directly access unauthorized memory regions.

#### [5_36b_ptrobj]: Context pointers vs packet data comparison

**Implementation Details:**
- Focuses on violations between XDP context structure pointers and packet data pointers.
- Tests real XDP context pointer comparisons (`struct xdp_md *xdp_ctx`) versus packet buffer data.
- Demonstrates context field violations (data_meta vs data) from different memory regions.
- Shows how eBPF verifier handles context-specific pointer arithmetic violations.
- Modified function signature to include context parameter for realistic testing.

**Verifier:** Passed (but produces undefined results that may leak memory layout information).

**Extra warnings:** None (only base warnings present)

**Exploitable:** **Limited** - Comparing XDP context pointers with packet data may reveal kernel memory layout relationships, but eBPF's memory protection prevents this from escalating to unauthorized access.

#### [5_36c_ptrobj]: Map pointers vs packet data comparison

**Implementation Details:**
- Targets violations between eBPF map value pointers (heap objects) and packet data.
- Tests comparisons between different map types (hash vs array maps) representing different heap objects.
- Demonstrates map value vs packet boundary violations across distinct memory regions.
- Shows realistic packet processing scenarios where map lookup results are incorrectly compared with network data.
- Includes additional hash map for testing cross-map pointer violations.

**Verifier:** Passed (but produces undefined results that may leak memory layout information).

**Extra warnings:** None (only base warnings present)

**Exploitable:** **Limited** - Cross-object pointer comparisons may leak information about kernel heap organization, but eBPF's isolation mechanisms prevent exploitation beyond information disclosure.

*Signed-by: Giovanni Nicosia*

---

### [5.39 taintnoproto]: Calling a function through a pointer without a prototype using tainted input

Calling a function through a pointer without a proper prototype leads to undefined behavior. If the function expects arguments and the caller provides incompatible or tainted input, the results are unpredictable.

**Implementation Details:**
- A function `restricted_sink(int i)` writes into an array using a tainted index: `s.array1[i] = 42;`.
- A function pointer `pf` is assigned the address of `restricted_sink` but declared with no prototype (`void (*pf)()`).
- The function is called with `(*pf)(tainted_val)`, where `tainted_val` is derived from `hdr->tcp->seq`, a value from the packet.
- This violates UB 39 and UB 41 by combining a tainted input with a call to a function without a prototype.

**Extra warnings:** Deprecated passing argument to function without prototype :
```
xdp_synproxy_kern.c:788:7: warning: passing arguments to a function without a prototype is deprecated in all versions of C and is not supported in C23 [-Wdeprecated-non-prototype]
  788 |         (*pf)(bpf_htons(hdr->tcp->seq)/1000); //This is the tainted input into unproto function call
      |              ^
6 warnings generated.
```

**Verifier:** Passed (compiler allows call, type mismatch undetected).

**Exploitable:** Limited — although the call is undefined, in practice the compiler will generate a call instruction with a fixed calling convention. The tainted value may corrupt stack arguments or registers, but within eBPF’s restricted environment the damage is confined and cannot be steered toward arbitrary memory writes. It primarily results in unpredictable logic, not exploitable memory corruption.

No extra tests due to unknown other possible callings.

*Signed-by: Giorgio Fardo*

---

### [5.40 taintformatio]: Tainted format string usage in eBPF helpers

Calls to the `sprintf` function that can result in writes outside the bounds of the destinzation array shall be diagnosed when any of its variadic arguments are tainted.

In this case we use the helper `bpf_snprintf`, on the other hand `bpf_trace_printk` cannot be exploited as the arguments are restricted. So no extra tests needed.

**Implementation Details:**
- A helper function `taintedBufPrint()` reads attacker-controlled values directly from the TCP header:
  - `tainted_value = hdr->tcp->seq;`
  - `tainted_len = hdr->tcp->window;`
- These values are passed to `bpf_snprintf`:
  ```c
  format = bpf_snprintf(buf, sizeof(buf), "%d", &args[0], tainted_len);
  ```
- Here, `tainted_len` is used as the *string length parameter*, directly influencing how much data `bpf_snprintf` attempts to write.
- The target buffer `buf[4]` is intentionally undersized, making the call unsafe if the verifier allowed it.
- The verifier detects this as **invalid indirect access to stack** and rejects the program (`call bpf_snprintf#165 invalid indirect access to stack`).

**Extra warnings:** None at compile time.

**Verifier:** **Rejected.** The verifier identifies an unbounded tainted access (`invalid indirect access to stack`).

**Exploitable:** Not exploitable. The verifier fully rejects the program before it can be JITed or run. Unlike user-space format string bugs, this does not result in buffer overflows or memory corruption in kernel/eBPF contexts, but simply prevents the program from loading

*Signed-by: Giorgio Fardo*


---

### [5.45 invfmtstr]: Invalid format strings in formatted I/O functions

Using format strings with conversion specifiers that don't match the provided arguments, invalid flag combinations, or incorrect argument counts leads to undefined behavior and potentially exploitable vulnerabilities.
#### [5_45_invfmtstr]: Invalid format strings with mismatched arguments

**Implementation Details:**
- **Type Mismatch (UB 160):** `bpf_printk("Parsing packet at offset %s\n", (long)data)` - %s expects string but receives integer, resulting in empty/garbage output.
- **Invalid Precision (UB 155):** `bpf_printk("IPv4 bounds check failed: %.5x\n", ...)` - precision with %x conversion may produce undefined formatting.
- **Argument Count Mismatch (UB 156):** `bpf_printk("IPv4 TCP header at %p, IHL=%d, proto=%d\n", hdr->tcp, hdr->ipv4->ihl)` - format has 3 specifiers but only 2 arguments provided.
- **Invalid Flag Combination (UB 157):** `bpf_printk("IPv6 nexthdr: %#d\n", hdr->ipv6->nexthdr)` - # flag not valid with %d conversion specifier.
- These violations don't cause verifier rejection but result in corrupted logging output, potentially hiding security events or leaking memory addresses through malformed prints.

**Verifier:** Passed (format string errors not detected by verifier, manifest at runtime).

**Extra warnings:** None (only base warnings present)

**Exploitable:** **Limited** - Format string mismatches corrupt logging output and may leak memory addresses through malformed prints, but eBPF's restricted environment prevents escalation to arbitrary memory access.

*Signed-by: Giovanni Nicosia*

---

### [5.46 taintsink]: Tainted, potentially mutilated, or out-of-domain integer values are used in a restricted sink

Using tainted, potentially mutilated, or out-of-domain integers in an integer restricted sink can result in accessing memory that is outside the bounds of existing objects.

In the context of **xdp_synproxy** using data received from external source, can be considered "tainted" because an attacker could craft it to contain arbitrary or malicious values. Certain operations are "restricted sinks" because they rely on the input value being within a specific, safe range (e.g., array indices, memory allocation sizes, loop counters, pointer arithmetic offsets).

This scenario is illustrated by two examples, each demonstrating different behaviors of the verifier. In the first case, the verifier **correctly rejects** the program. In the second case, however, it allows the program to pass and **permits an out-of-bounds write** under certain conditions that go undetected.

#### [5_46a_taintsink]

**Implementation Details:**
- A small `char` array `policy_flags[32]` is declared on the stack. Its size is intentionally limited to `32` bytes to make it highly susceptible to out-of-bounds access by typical network values.
- The destination port (`hdr->tcp->dest`) from the incoming TCP packet is extracted and stored in `tainted_dest_port`. For the purpose of our test we can consider this value "tainted" as it can range from `0` to `65535`, and highly probable to trigger an out-of-bounds read if used as array index to access our small buffer.
- The core violation is demonstrated by the line: `char accessed_flag_value = policy_flags[tainted_dest_port]`. Here, the `tainted_dest_port` is used directly as an array index into `policy_flags` **without any bounds checking**.

In this particular case the eBPF verifier performs a correct memory validation. Since `tainted_dest_port` can clearly exceed the array's bounds, and there is no bounds checking, the verifier will detect a potential out-of-bounds memory read. As a result, the verifier will reject the eBPF program load.

**Verifier:** Not passed:
```
math between fp pointer and register with unbounded min value is not allowed
```

**Extra warnings:** None.

**Exploitable:** Not possible.

#### [5_46b_taintsink]

**Implementation Details:**
- The function `print_field` is added as a helper to visually dump the contents of byte arrays for debugging the memory state.
- While Variable Length Arrays (VLAs) are not supported in eBPF, this example mimics a common VLA-related vulnerability by, instead, using a fixed-size target buffer (`hdr->eth->h_dest`, which is a 6-byte array within the `ethhdr` structure) and a tainted length for a loop-based write operation.
- The `tainted_length_from_packet` is derived from `hdr->tcp->doff * 4`. The `doff` field (TCP data offset) has a valid range of 5 to 15 (representing 20 to 60 bytes for the TCP header length). This means `tainted_length_from_packet` can be up to 60 bytes, which is significantly larger than the 6-byte `h_dest` buffer.
- The code explicitly makes a copy of `hdr->eth->h_source` (which immediately follows `h_dest` in `struct ethhdr`) before the loop.
- A `#pragma unroll` loop attempts to write `i` into `hdr->eth->h_dest[i]` for `i` up to `tainted_length_from_packet`. This will attempt to write beyond the 6-byte boundary of `hdr->eth->h_dest`, overflowing into `hdr->eth->h_source` and potentially further.
- A `if (i > (h_source_size + hdr->eth->h_proto)) { break; }` condition is an attempt to introduce a verifier-friendly exit path to ensure loop termination, but its threshold is intentionally bogus (well beyond where the actual overflow occurs) that it **does not prevent the bounds violation itself**.
- After the loop, `__builtin_memcmp` is used to compare the current `hdr->eth->h_source` with its original copy. If the `memcmp` reveals a mismatch, it confirms that the out-of-bounds write to `h_dest` successfully corrupted `h_source`, directly demonstrating the overflow.

Under this specific condition, the verifier allows the eBPF program to load, even though it permits an **out-of-bounds write** within a controlled memory layout. This vulnerability is confirmed at runtime.

**Verifier:** Passed.

**Extra warnings:** None.

**Exploitable:**
- Memory safety exploitation (kernel R/W): No, not possible.
- Logic exploitation (attacker-controlled packet alteration): Yes, possible.

#### [5_46c_taintsink]

**Implementation Details**

- Inside `syncookie_handle_syn`, a `__u32 tainted_vla_size` variable is declared and initialized with a value derived from `hdr->tcp->doff * 4`.
- The line `char vla_buffer[tainted_vla_size]` attempts to declare a Variable Length Array `vla_buffer` using this runtime-determined size.
- This program **will not pass verification**, regardless of the taintedness of `tainted_vla_size`. The eBPF verifier, as mentioned in the previous example, prohibits VLAs. The compilation succeeds, but the attempt to load such a BPF program into the kernel will result in a clear rejection message from the verifier.

**Verifier:** Not passed:

```
Address R11 is invalid (result of VLA attempt)
```

**Extra warnings:** None.

**Exploitable:** Not possible.

*Signed-by: Francesco Rollo*
