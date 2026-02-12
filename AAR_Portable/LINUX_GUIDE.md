# Linux Deployment Guide for AAR

## 1. Copy Files
Copy this entire folder (`AAR_Portable`) to your Linux machine (e.g., to `~/AAR_Portable`).

## 2. Prepare the Environment
Open a terminal on your Linux machine and navigate to the folder:
```bash
cd ~/AAR_Portable
```

Make the router executable:
```bash
chmod +x aar_router_linux
```

## 3. Dependencies
The router needs these standard Linux tools to control the network:
- `ip` (iproute2)
- `iw` (wireless tools)
- `wg-quick` (WireGuard)

You can install them (on Debian/Ubuntu/Kali) with:
```bash
sudo apt update
sudo apt install iproute2 iw wireguard-tools
```

## 4. Run the Router
The router needs `sudo` privileges to change network settings and switch Wi-Fi channels.

```bash
sudo ./aar_router_linux
```

## 5. (Optional) Set Up VPN Privacy
For the **Gaming VPN** and **Privacy Shield** to work, your Linux laptop needs a WireGuard configuration.

1.  **Get Configuration**: Sign up for a VPN (like Mullvad, ProtonVPN, or your own VPS) and download a `wg0.conf` file.
2.  **Use Our Template**: If you have the keys but no file, specific details are in `wg0.conf.example` in this folder.
3.  **Install Config**: Copy your config file to the system folder:
    ```bash
    sudo cp wg0.conf.example /etc/wireguard/wg0.conf
    # (Edit the file first to add your real keys!)
    ```
4.  **Test It**:
    Click "Enable VPN" on the AAR Dashboard. The router will use this file to encrypt your traffic.

## 6. Enable Wi-Fi for Phones/Friends
Your Linux laptop will act as the "Tower". You need to turn on the Hotspot signal.

**Option A: Automated Script (Best)**
I have included a script to do this for you.
```bash
chmod +x start_hotspot.sh
./start_hotspot.sh
```
This will create a network named **"AAR-Link"** (Password: `secureconnect`).

**Option B: Manual Setup**
If the script fails, just use your standard Linux settings:
1. Open **Settings** > **Wi-Fi**.
2. Click the menu/dots and select **"Turn On Wi-Fi Hotspot"**.
3. Use the network name/password shown on your screen.

## 6. Access the Dashboard
Once running, open your web browser on the Linux machine (or any connected phone!) and go to:
- **Local:** http://localhost:8080
- **Phone:** http://<YOUR_LINUX_IP>:8080

## Troubleshooting
- **Interface Name**: The router assumes your Wi-Fi card is `wlan0`. If `start_hotspot.sh` fails, run `ip link` to find your real name (e.g., `wlp2s0`) and edit the script (`nano start_hotspot.sh`).
- **Permission denied**: Make sure you used `chmod +x` and `sudo`.
