#!/bin/bash

set -e

sudo modprobe nf_conntrack
sudo sysctl -w net.ipv4.tcp_syncookies=2
sudo sysctl -w net.ipv4.tcp_timestamps=1
sudo sysctl -w net.netfilter.nf_conntrack_tcp_loose=0
sudo iptables -t raw -I PREROUTING  -i enp1s0 -p tcp -m tcp --syn --dport 80 -j CT --notrack
sudo iptables -t filter -A INPUT -i enp1s0 -p tcp -m tcp --dport 80 -m state --state INVALID,UNTRACKED -j SYNPROXY --sack-perm --timestamp --wscale 7 --mss 1460
sudo iptables -t filter -A INPUT -i enp1s0 -m state --state INVALID -j DROP
