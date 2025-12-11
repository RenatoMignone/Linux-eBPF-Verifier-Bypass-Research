#!/usr/bin/env bash

set -e

# Constants
IMAGE_URL="https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img"
IMAGE_NAME="noble-server-cloudimg-amd64.img"
IMAGE_PATH="/var/lib/libvirt/images/$IMAGE_NAME"
VM_NAME="ubuntu-xdp-24.04"
MEMORY="4096"
VCPUS="2"
DISK_SIZE="20"
OS_VARIANT="ubuntu24.04"
NETWORK="default"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
USER_DATA="$SCRIPT_DIR/user-data.yaml"
META_DATA="$SCRIPT_DIR/meta-data.yaml"

# Shared folder configuration
# This folder on the host will be accessible inside the VM at /mnt/shared
SHARED_FOLDER="/home/ren/PERSONAL_DIRECTORY/CyberSecurity/Second_Year/First_Semester/SVT/Project_EBPF/src/shared"
SHARE_NAME="host_shared"
CODEBASE_DIR="$SCRIPT_DIR/.."

# Functions
function usage {
  echo "Usage: $0 [create <ssh_pubkey_file> | destroy | connect]"
  echo ""
  echo "Shared folder: $SHARED_FOLDER"
  echo "  - Place files here on the host"
  echo "  - Access them in the VM at: /mnt/shared"
  exit 1
}

function confirm {
  read -p "Are you sure you want to $1 the VM '$VM_NAME'? [y/N]: " CONFIRM
  CONFIRM=${CONFIRM:-n}
  [[ "$CONFIRM" =~ ^[Yy]$ ]] || { echo "Aborted."; exit 0; }
}

function create_vm {
  SSH_KEY_FILE="$1"
  if [ -z "$SSH_KEY_FILE" ]; then
    echo "Error: SSH pubkey file required for VM creation."
    usage
  fi
  if [ ! -f "$SSH_KEY_FILE" ]; then
    echo "Error: File not found: $SSH_KEY_FILE"
    exit 2
  fi

  # Check if it's a .pub file
  if [[ "$SSH_KEY_FILE" != *.pub ]]; then
    echo "Error: The SSH key file must be a public key."
    exit 3
  fi

  confirm "create"

  SSH_KEY=$(<"$SSH_KEY_FILE")
  # Use awk instead of sed to avoid escaping issues with SSH keys
  awk -v key="$SSH_KEY" '{gsub(/__SSH_PUBLIC_KEY__/, key); print}' "$USER_DATA" > "${USER_DATA}.tmp" && mv "${USER_DATA}.tmp" "$USER_DATA"

  if virsh --connect qemu:///system list --all | grep -qw "$VM_NAME"; then
    echo "Error: A VM with the name '$VM_NAME' already exists."
    exit 1
  fi

  # Create shared folder if it doesn't exist
  if [ ! -d "$SHARED_FOLDER" ]; then
    echo "Creating shared folder: $SHARED_FOLDER"
    mkdir -p "$SHARED_FOLDER"
  fi

  # Copy codebase to shared folder
  echo "Copying codebase to shared folder..."
    cp -r "$CODEBASE_DIR"/* "$SHARED_FOLDER/" 2>/dev/null || true
  # Remove the nested shared folder to avoid recursion
  rm -rf "$SHARED_FOLDER/shared" 2>/dev/null || true
  echo "Codebase copied to: $SHARED_FOLDER"

  if [ -f "$IMAGE_PATH" ]; then
    read -p "Image '$IMAGE_PATH' already exists. Re-download it? [y/N]: " REDOWNLOAD
    REDOWNLOAD=${REDOWNLOAD:-n}
    if [[ "$REDOWNLOAD" =~ ^[Yy]$ ]]; then
      wget -O "$IMAGE_NAME" "$IMAGE_URL"
      sudo mv "$IMAGE_NAME" "$IMAGE_PATH"
    else
      echo "Using existing image."
    fi
  else
    wget -O "$IMAGE_NAME" "$IMAGE_URL"
    sudo mv "$IMAGE_NAME" "$IMAGE_PATH"
  fi

  echo "Creating VM with shared folder..."
  sudo virt-install \
    --name "$VM_NAME" \
    --memory "$MEMORY" \
    --vcpus "$VCPUS" \
    --disk size="$DISK_SIZE",backing_store="$IMAGE_PATH" \
    --os-variant "$OS_VARIANT" \
    --network network="$NETWORK",model=virtio \
    --filesystem source="$SHARED_FOLDER",target="$SHARE_NAME",driver.type=virtiofs \
    --memorybacking source.type=memfd,access.mode=shared \
    --cloud-init user-data="$USER_DATA",meta-data="$META_DATA"

  echo ""
  echo "VM creation complete."
  echo ""
  echo "=========================================="
  echo "SHARED FOLDER SETUP"
  echo "=========================================="
  echo "Host path:   $SHARED_FOLDER"
  echo "VM mount:    /mnt/shared"
  echo ""
  echo "To mount inside VM, run:"
  echo "  sudo mkdir -p /mnt/shared"
  echo "  sudo mount -t virtiofs $SHARE_NAME /mnt/shared"
  echo ""
  echo "Or add to /etc/fstab for auto-mount:"
  echo "  $SHARE_NAME /mnt/shared virtiofs defaults 0 0"
  echo "=========================================="

  IP=$(virsh --connect qemu:///system domifaddr "$VM_NAME" | grep ipv4 | awk '{print $4}' | cut -d'/' -f1)
  if [ -n "$IP" ]; then
    echo "VM IP: $IP"
  else
    echo "Could not retrieve IP."
  fi
}

function destroy_vm {
  if ! virsh --connect qemu:///system dominfo "$VM_NAME" &>/dev/null; then
    echo "Error: VM '$VM_NAME' does not exist."
    exit 1
  fi

  confirm "destroy"

  echo "Destroying VM..."
  virsh --connect qemu:///system destroy "$VM_NAME" || echo "VM not running."
  virsh --connect qemu:///system undefine "$VM_NAME"

  if [ -f "$IMAGE_NAME" ]; then
    read -p "Remove disk image '$IMAGE_NAME'? [y/N]: " CONFIRM_DELETE
    CONFIRM_DELETE=${CONFIRM_DELETE:-n}
    if [[ "$CONFIRM_DELETE" =~ ^[Yy]$ ]]; then
      rm -f "$IMAGE_NAME"
      echo "Disk image removed."
    else
      echo "Disk image not removed."
    fi
  fi

  echo "VM destroyed and undefined."
}

function connect_vm {
  if ! virsh --connect qemu:///system dominfo "$VM_NAME" &>/dev/null; then
    echo "Error: VM '$VM_NAME' does not exist."
    exit 1
  fi

  STATE=$(virsh --connect qemu:///system domstate "$VM_NAME")
  if [[ "$STATE" != "running" ]]; then
    echo "Starting VM..."
    virsh --connect qemu:///system start "$VM_NAME"
  fi

  echo "Connecting to VM console..."
  virsh --connect qemu:///system console "$VM_NAME"
}

case "$1" in
  create)
    create_vm "$2"
    ;;
  destroy)
    destroy_vm
    ;;
  connect)
    connect_vm
    ;;
  *)
    usage
    ;;
esac
