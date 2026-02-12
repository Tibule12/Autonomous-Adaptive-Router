AAR Portable Deployment
========================

This folder contains everything needed to run the Autonomous Adaptive Router (AAR) on Windows or Linux.

Windows Usage:
--------------
1. Double-click run.bat to start the router.
2. The browser will open automatically to the dashboard.

Linux Usage:
------------
1. Copy this folder to your Linux machine.
2. Follow the detailed instructions in LINUX_GUIDE.md.
   - Run 'chmod +x aar_router_linux'
   - Broadcast Signal: './start_hotspot.sh'
   - Start Router: 'sudo ./aar_router_linux'
   - (Optional) Configure VPN: Edit 'wg0.conf.example' with your provider's keys.

Network Requirements:
---------------------
- Windows: No special requirements (simulation mode included).
- Linux: Needs 'sudo' permissions and 'wireguard-tools' installed for real operations.

