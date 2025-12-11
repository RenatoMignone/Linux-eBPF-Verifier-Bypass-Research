#!/bin/bash

set -euo pipefail

SESSION="synproxy_setup"
IFACE="enp2s0"

# Get IPv4 address
IP_ADDR=$(ip -4 addr show "$IFACE" | grep -oP '(?<=inet\s)\d+(\.\d+){3}' | head -1)

if [ -z "$IP_ADDR" ]; then
    echo "Error: Could not determine IP address for interface $IFACE"
    exit 1
fi

# Kill any existing session
tmux has-session -t "$SESSION" 2>/dev/null && tmux kill-session -t "$SESSION"

# Create session and first window
tmux new-session -d -s "$SESSION" -n main
tmux set-option -t "$SESSION" mouse on

# Step 1: Run xdp_synproxy in the first pane
tmux send-keys -t "$SESSION":0.0 "sudo ./xdp_synproxy --iface $IFACE --mss4 1460 --mss6 1440 --wscale 7 --ttl 64 --ports 80" C-m

# Step 2: Split horizontally to create right pane (future nc listener)
tmux split-window -h -t "$SESSION":0.0

# Step 3: Select left pane again and split vertically (trace_pipe below)
tmux select-pane -L  # move to left pane
tmux split-window -v

# Now assign pane variables explicitly to avoid confusion
# Layout:
# ┌────────────┬────────────┐
# │    PANE0   │   PANE2    │
# ├────────────┘            │
# │    PANE1                │
# └─────────────────────────┘

# Get pane IDs
PANE0=$(tmux list-panes -t "$SESSION" -F "#{pane_index}" | sed -n '1p')  # top-left
PANE1=$(tmux list-panes -t "$SESSION" -F "#{pane_index}" | sed -n '2p')  # bottom-left
PANE2=$(tmux list-panes -t "$SESSION" -F "#{pane_index}" | sed -n '3p')  # right

# Send trace_pipe to bottom-left (PANE1)
tmux send-keys -t "$SESSION":0.$PANE1 "sudo cat /sys/kernel/debug/tracing/trace_pipe" C-m

# Send nc commands to right pane (PANE2)
tmux send-keys -t "$SESSION":0.$PANE2 "echo \"Run this command on the host: nc $IP_ADDR 80 -v\"" C-m
tmux send-keys -t "$SESSION":0.$PANE2 "echo \"Started server listening on port 80\"" C-m
tmux send-keys -t "$SESSION":0.$PANE2 "sudo nc -lvnp 80" C-m

# Focus top-left pane
tmux select-pane -t "$SESSION":0.$PANE0
tmux attach-session -t "$SESSION"
