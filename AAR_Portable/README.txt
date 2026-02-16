AAR Portable Deployment
========================

This folder contains everything needed to run the Autonomous Adaptive Router (AAR).

Windows Usage:
--------------
Double-click 'run.bat'.

Linux Usage:
------------
Linux prevents double-clicking programs by default for security.
You usually need to utilize the Terminal.

Option 1: The Terminal (Recommended)
1. Right-click inside this folder > 'Open in Terminal'
2. Type: ./run_linux.sh

Option 2: Try Double-Clicking
1. Right-click 'run_linux.sh' -> Properties -> Permissions -> Check 'Allow executing file as program'.
2. Now you can try double-clicking it (select 'Run in Terminal' if asked).

VPN Configuration:
------------------
See 'wg0.conf.example' to set up privacy features.

DOWNLOAD LINKS (Before you travel):
----------------------------------
If your target laptop doesn't have Linux yet:
- Ubuntu Desktop: https://ubuntu.com/download/desktop

Most tools (commands like 'ip', 'iw') are built-in, but 'WireGuard' might be missing.
If that laptop has NO internet, download this package now and put it on the USB:
- WireGuard Tools: https://packages.ubuntu.com/jammy/wireguard-tools

