# AAR Portable - Setup Guide

This folder contains the self-contained Autonomous Adaptive Router simulation.

## Interpretation
You can run this on any Windows laptop without installing anything (no Go, no Node.js required).

## Instructions
1. Copy this entire `AAR_Portable` folder to your laptop (Desktop or Documents).
2. Double-click `run.bat`.
3. A command window will open (starting the AI Brain).
4. A browser window should automatically open to `http://localhost:8080`.

## ⚡ Multi-Laptop "Stress Test" Mode
Turn your second laptop into a traffic generator!
1. Run the Router on Laptop A as usual.
2. Find Laptop A's IP address (Open CMD, type `ipconfig`).
3. On Laptop B, open the file `stress_test.html` (inside this folder) in Chrome/Edge.
4. Enter Laptop A's IP address.
5. Click **START ATTACK**.
6. Watch the Dashboard on Laptop A - you will see latency spike and packet loss increase as the AI tries to fight it!

## Features
- **Simulation**: Full logic is active.
- **Persistence**: Metrics are saved to `metrics_history.json` automatically.
- **Dashboard**: Full visualization capability.

## Networking Note
Currently, this runs in "Standalone Mode". If you want this laptop to talk to your other laptop (Mesh Network), ensure they are on the same Wi-Fi and update the IP address in the configuration.
