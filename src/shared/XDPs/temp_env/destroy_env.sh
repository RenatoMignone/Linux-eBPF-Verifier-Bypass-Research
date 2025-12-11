sudo ip netns del ns-client
sudo ip netns del ns-router
sudo ip netns del ns-server
sudo ip link del veth-client || true
sudo ip link del veth-router2 || true
