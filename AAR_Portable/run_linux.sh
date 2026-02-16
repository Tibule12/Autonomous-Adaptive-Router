#!/bin/bash

# AAR "Double-Click" Launcher for Linux
# This script tries to open a terminal window to run the router.

cd "$(dirname "$0")"

echo "=============================================="
echo "   Autonomous Adaptive Router (Linux Mode)"
echo "=============================================="
echo ""

# 1. Start Hotspot (Optional, based on user choice)
echo ">> Step 1: Wi-Fi Signal"
echo "   Do you want to start the Hotspot? (y/n)"
read -r -t 10 response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])+$ ]] || [[ -z "$response" ]]; then
    ./start_hotspot.sh
else
    echo "   Skipping Hotspot..."
fi

echo ""
echo ">> Step 2: Router Core"
echo "   We need Admin (Sudo) permissions to control the network."
echo "   Please enter your password if asked."
echo ""

# 2. Start Router Engine
sudo ./aar_router_linux

echo ""
echo "Router stopped. Press Enter to exit."
read
