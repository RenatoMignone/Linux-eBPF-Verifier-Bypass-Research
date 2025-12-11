# XVTlas - XDP Verifier Launch Automation Suite

**XVTlas** automates compilation, patching, loading, and verification of XDP eBPF programs. It leverages `bpftool`, `make`, and a Python "pretty verifier" script to provide structured loading and output for BPF/XDP developers.

---

## 🧰 Requirements

- Go (1.18+ recommended)
- `bpftool` (from `iproute2` package)
- `clang`, `llvm`, and `make` for compiling eBPF programs
- Python 3 with your `pretty_verifier.py` script
- PrettyVerifier
- `tmux` (for interactive session mode)

---

## ⚙️ Build

```bash
go build -o xvtlas .
```

This will produce an `xvtlas` binary in the current directory.

---

## 🚀 Usage

### ▶️ Standard Directory Mode

Run on a directory of eBPF programs:

```bash
./xvtlas/xvtlas \
  --export "./output/" \
  --kernel "6.8.58" \
  --path "./template_folders/" \
  --pretty "./pretty-verifier/pretty_verifier.py" \
  --verbose \
  --interactive
```

This compiles all programs found under `--path`, loads and verifies them using `bpftool`, and saves logs to `--export`.

---

### 🩹 Patch-Based Mode

```bash
./xvtlas/xvtlas \
  --export "./output/" \
  --kernel "6.8.58" \
  --patch-path "./XDPs/xdp_synproxy/patches/" \
  --base-file "./linux-kernel/tools/testing/selftests/bpf/xdp_synproxy_kern.c" \
  --pretty "./pretty-verifier/pretty_verifier.py" \
  --verbose \
  --interactive
```

This will:
- Walk the patch folders under `--patch-path`
- Apply each folder's `patch.diff` to the `--base-file`
- Compile using the provided Makefile
- Load, verify, and unpin each resulting program
- Clean up the patched file (unless `--keep-patched` is set)

---

### 🧪 Single Patch Interactive Mode

```bash
 ./xvtlas/xvtlas \
  --run-single "./XDPs/xdp_synproxy/patches/path/to/patch/*.patch" \
  --base-file "./XDPs/xdp_synproxy/xdp_synproxy_kern.c"
```

This:
- Applies the patch to the base file
- Saves the state to `/tmp/xvtlas.swp`
- Compiles the kernel code
- Launches `start_session.sh` inside a `tmux` session for manual interaction

You can detach from the tmux session (`Ctrl + b`, then `d`) to return to the CLI and continue cleanup.

**❗ Only `--run-single` and `--base-file` are allowed together. All other flags are ignored.**

---

### Multi patch report creation 

```bash
./xvtlas/xvtlas \
  --export "./output/" \
  --kernel "v6.8" \
  --patch-path ./XDPs/xdp_synproxy/patches/ \
  --base-file ./XDPs/xdp_synproxy/xdp_synproxy_kern.c \
  --pretty ./pretty-verifier/pretty_verifier.py \
  --save-logs \
  --interactive \
  --verbose
```

This will:
- Walk the patch folders under `--patch-path`
- Apply each folder's `<patch-name>.patch` to the `--base-file`
- Compile using the provided Makefile
- Load, verify, and unpin each resulting program
- Clean up the patched file (unless `--keep-patched` is set)


---

### 💣 Restore / Cleanup Previous Session

To restore the git state and clean up a previous `--run-single` session:

```bash
 ./xvtlas/xvtlas --destroy
```

This will:
- Read `/tmp/xvtlas.swp`
- Run `git reset --hard` to the saved commit
- Run `make clean` in the base directory
- Run `./destroy_session.sh` in the same directory (to close tmux)
- Delete the `.swp` file

**❗ This flag must be used alone (no other flags allowed)**

---

## 🧾 Command Line Options

| Flag               | Description                                                                 |
|--------------------|-----------------------------------------------------------------------------|
| `--path`           | Root path for standard eBPF program directories                             |
| `--patch-path`     | Folder containing subfolders with `patch.diff`, `config.yaml`, Makefile     |
| `--base-file`      | Path to the base file (used with `--patch-path` or `--run-single`)          |
| `--run-single`     | Apply one patch interactively (must be used only with `--base-file`)        |
| `--destroy`        | Reset state and clean up previous interactive run                           |
| `--export`         | Output directory for logs, reports, and CSVs                                |
| `--pretty`         | Path to the Python pretty verifier script                                   |
| `--kernel`         | Target kernel version                                                       |
| `--verbose`        | Show detailed logs from build/patching/verifier                             |
| `--interactive`    | Prompt on failure before continuing                                         |
| `--keep-patched`   | Do not revert patched files after run                                       |

---

## 📁 Folder Structure Examples

### Structure for `--path` (Standard Mode)

```
template_folders/
├── example1/
│   ├── main.c
│   ├── config.yaml
│   └── Makefile
├── example2/
│   ├── ...
```

### Structure for `--patch-path` (Patch Mode)

```
rules/
├── 001/
│   ├── patch.diff
│   ├── config.yaml
│   └── Makefile
├── 002/
│   ├── ...
```

Each folder contains a patch to be applied to the same `--base-file`.

---

## Report format and columns

The report is a CSV file that is both showed on screen after the full run and saved in the `output` folder.

In order we have : 
- id : progressive number of the test
- filename : name of the `.patch` file applied
- load parameters : `unused`
- vuln_number : `unused`
- compiled : if it has compiled without errors
- verified : if it has passed the verifier
- loaded : if it correcly loaded and pinned the maps
- load_output : `unused`

---

## 📌 Notes

- Programs are pinned under `/sys/fs/bpf/{program_name}` and unpinned automatically
- Logs are saved to the export folder with `.log` and `.csv` extensions
- The tool restores original state after running patches (unless `--keep-patched` is used)
- `--run-single` is interactive and launches `tmux`
- `--destroy` safely resets repo and closes `tmux`

---

## 🆘 Help

```bash
./xvtlas --help
```

If no arguments are provided, help is printed and the tool exits.

---

