# Setup USB Portable Version
$ErrorActionPreference = "Stop"

Write-Host "Building Backend..."
Push-Location backend
go build -o server.exe main.go
if ($LASTEXITCODE -ne 0) { Write-Error "Backend build failed"; exit 1 }
Pop-Location

Write-Host "Building Frontend..."
Push-Location frontend
npm run build
if ($LASTEXITCODE -ne 0) { Write-Error "Frontend build failed"; exit 1 }
Pop-Location

Write-Host "Creating Portable Folder..."
$usbDir = "AAR_Portable"
if (Test-Path $usbDir) { Remove-Item $usbDir -Recurse -Force }
New-Item -ItemType Directory -Path $usbDir | Out-Null
New-Item -ItemType Directory -Path "$usbDir/frontend_dist" | Out-Null

Copy-Item "backend/server.exe" "$usbDir/server.exe"
Copy-Item "frontend/dist/*" "$usbDir/frontend_dist" -Recurse
Copy-Item "stress_test.html" "$usbDir/stress_test.html"
Copy-Item "PORTABLE_README.txt" "$usbDir/README.txt"

# Create a runner script
$batContent = @"
@echo off
echo Starting Autonomous Adaptive Router (Portable)...
echo Opening Dashboard...
start http://localhost:8080
echo Starting Backend...
server.exe
pause
"@
Set-Content "$usbDir/run.bat" $batContent

Write-Host "Success! Copy the '$usbDir' folder to your USB drive."
