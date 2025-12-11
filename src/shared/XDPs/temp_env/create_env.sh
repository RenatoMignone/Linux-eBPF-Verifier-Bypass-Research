#!/bin/bash
sudo modprobe nf_conntrack

# Clean existing setup
sudo ip netns del ns-client 2>/dev/null || true
sudo ip netns del ns-router 2>/dev/null || true
sudo ip netns del ns-server 2>/dev/null || true
sudo ip link del veth-client 2>/dev/null || true
sudo ip link del veth-router2 2>/dev/null || true

# Create namespaces
sudo ip netns add ns-client
sudo ip netns add ns-router
sudo ip netns add ns-server

# Create veth pairs
sudo ip link add veth-client type veth peer name veth-router1
sudo ip link add veth-router2 type veth peer name veth-server

# Assign to namespaces
sudo ip link set veth-client netns ns-client
sudo ip link set veth-router1 netns ns-router
sudo ip link set veth-router2 netns ns-router
sudo ip link set veth-server netns ns-server

# Setup client namespace
sudo ip netns exec ns-client ip addr add 10.3.3.9/24 dev veth-client
sudo ip netns exec ns-client ip link set veth-client up
sudo ip netns exec ns-client ip link set lo up
sudo ip netns exec ns-client ip route add default via 10.3.3.8
sudo ip netns exec ns-client ip route add 10.6.6.0/24 via 10.3.3.8

# Setup server namespace
sudo ip netns exec ns-server ip addr add 10.6.6.6/24 dev veth-server
sudo ip netns exec ns-server ip link set veth-server up
sudo ip netns exec ns-server ip link set lo up
sudo ip netns exec ns-server ip route add default via 10.6.6.8
sudo ip netns exec ns-server ip route add 10.3.3.0/24 via 10.6.6.8

# Setup router namespace
sudo ip netns exec ns-router ip addr add 10.3.3.8/24 dev veth-router1
sudo ip netns exec ns-router ip addr add 10.6.6.8/24 dev veth-router2
sudo ip netns exec ns-router ip link set veth-router1 up
sudo ip netns exec ns-router ip link set veth-router2 up
sudo ip netns exec ns-router ip link set lo up

# Enable forwarding and SYNPROXY sysctls
sudo ip netns exec ns-router sysctl -w net.ipv4.ip_forward=1
sudo ip netns exec ns-router sysctl -w net.ipv4.tcp_syncookies=2
sudo ip netns exec ns-router sysctl -w net.ipv4.tcp_timestamps=1
sudo ip netns exec ns-router sysctl -w net.netfilter.nf_conntrack_tcp_loose=0

# Setup iptables rules for SYNPROXY
sudo ip netns exec ns-router iptables -t raw -I PREROUTING -i veth-router1 -p tcp --syn --dport 80 -j CT --notrack
sudo ip netns exec ns-router iptables -t filter -A FORWARD -i veth-router1 -p tcp --dport 80 -m state --state INVALID,UNTRACKED -j SYNPROXY --sack-perm --timestamp --wscale 7 --mss 1460
sudo ip netns exec ns-router iptables -t filter -A FORWARD -i veth-router1 -m state --state INVALID -j DROP

# Print the network schema
cat << "EOF"

Routing table entries:
----------------------
ns-client:  ip route add 10.6.6.0/24 via 10.3.3.8
ns-server:  ip route add 10.3.3.0/24 via 10.6.6.8

Topology:
---------

+----------------+      +-----------------------------+      +----------------+
|                |      |                             |      |                |
|   ns-client    |      |         ns-router           |      |   ns-server    |
|   10.3.3.9     |<---> | veth-router1   veth-router2 |<---> |   10.6.6.6     |
|                |      | 10.3.3.8        10.6.6.8    |      |                |
+----------------+      +-----------------------------+      +----------------+

EOF

