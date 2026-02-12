#!/bin/bash
# AAR Hotspot Launcher
# This script uses NetworkManager to broadcast a Wi-Fi signal.

# CONFIGURATION
INTERFACE="wlan0"  # Change this if your wifi card is named differently (e.g., wlp2s0)
SSID="AAR-Link"
PASSWORD="secureconnect"

echo "----------------------------------------"
echo "  AAR Autonomous Hotspot Creator"
echo "----------------------------------------"

# Check for NetworkManager
if ! command -v nmcli &> /dev/null
then
    echo "❌ Error: 'nmcli' tool not found."
    echo "Please open your System Settings > Wi-Fi and create a Hotspot manually."
    exit 1
fi

echo ">> Config: Interface=$INTERFACE | SSID=$SSID"

# Create/Enable the hotspot
# Note: 'con delete' removes old profiles to ensure clean start
echo ">> Cleaning old connections..."
sudo nmcli con delete "$SSID" > /dev/null 2>&1

echo ">> Starting Hotspot..."
if sudo nmcli dev wifi hotspot ifname "$INTERFACE" ssid "$SSID" password "$PASSWORD"; then
    echo ""
    echo "✅ SUCCESS! Network is Active."
    echo "----------------------------------------"
    echo "Connect your friends/family to:"
    echo "📡 Name: $SSID"
    echo "🔑 Pass: $PASSWORD"
    echo "----------------------------------------"
else
    echo ""
    echo "❌ Failed to start hotspot."
    echo "Try checking your interface name with: ip link show"
    echo "Then edit this script: nano start_hotspot.sh"
    echo "Or use your Linux Settings menu to create a Hotspot."
fi
