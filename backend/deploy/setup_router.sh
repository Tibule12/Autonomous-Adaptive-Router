#!/bin/bash
# AAR Router Initialization Script (Phase 1)
# RUN AS ROOT (sudo)

set -e

echo "[AAR] 🚀 Initializing Base Router Configuration..."

# 1. Install Dependencies
echo "[AAR] 📦 Installing hostapd, dnsmasq, wireguard..."
apt-get update
apt-get install -y hostapd dnsmasq wireguard nftables iptables-persistent

# 2. Stop conflicting services
echo "[AAR] 🛑 Stopping interfering services..."
systemctl stop dnsmasq || true
systemctl stop hostapd || true
systemctl stop wpa_supplicant || true

# 3. Configure Wi-Fi Interface (Static IP)
echo "[AAR] 📡 Configuring wlan0 static IP (192.168.50.1)..."
ip link set wlan0 down
ip addr flush dev wlan0
ip addr add 192.168.50.1/24 dev wlan0
ip link set wlan0 up

# 4. Enable IP Forwarding (The "Router" part)
echo "[AAR] 🛣️ Enabling IP Forwarding..."
sysctl -w net.ipv4.ip_forward=1
sed -i 's/#net.ipv4.ip_forward=1/net.ipv4.ip_forward=1/g' /etc/sysctl.conf

# 5. Configure NAT (Share Internet from eth0 to wlan0)
echo "[AAR] 🛡️ Configuring Firewall (NAT)..."
# Clear old rules
iptables -F
iptables -t nat -F

# Masquerade (NAT) traffic leaving via eth0 (WAN)
iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE

# Forward traffic between interfaces
iptables -A FORWARD -i eth0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -A FORWARD -i wlan0 -o eth0 -j ACCEPT

# Save rules
netfilter-persistent save

# 6. Start Services
echo "[AAR] 🏁 Starting Network Services..."
dnsmasq -C ./dnsmasq.conf
hostapd -B ./hostapd.conf

echo "[AAR] ✅ Router Active!"
echo "[AAR] 📶 SSID: AAR_Secure_Network"
echo "[AAR] 🔑 Pass: aar_password_123"
