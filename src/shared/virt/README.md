# vmctl
This bash script provides a simple CLI to **create**, **destroy**, or **connect** to the virtual machine needed for this test environment.

## 🔧 Requirements

- `virt-install` and `libvirt` (`sudo apt install virt-manager libvirt-daemon-system`)
- `wget`
- A valid SSH **public key file**

## 📦 Usage

```bash
chmod +x vm-tool.sh
./vmctl.sh <command> [args]
```

### Commands

#### Create a VM

```bash
./vm-tool.sh create ~/.ssh/id_rsa.pub
```

- Prompts for confirmation before creating.
- Copies the public SSH key into `user-data.yaml`.
- Downloads Ubuntu 24.04 cloud image if not already present.
- Automatically provisions the VM using `virt-install`.

#### Destroy a VM

```bash
./vm-tool.sh destroy
```

- Prompts for confirmation before destroying.
- Destroys and undefines the VM.
- Optionally deletes the downloaded disk image.

#### Connect to VM Console

```bash
./vm-tool.sh connect
```

- Starts the VM if it's not running.
- Opens a direct console via `virsh`.

### ⚠️ Notes

- The script assumes the VM name is `ubuntu-xdp-24.04`.
- Disk image is stored at: `/var/lib/libvirt/images/noble-server-cloudimg-amd64.img`

### 🧪 Example

```bash
./vm-tool.sh create ~/.ssh/id_ed25519.pub
./vm-tool.sh connect
./vm-tool.sh destroy
```
