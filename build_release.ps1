Write-Host "🚧 Building AAR for Production (Linux/ARM64 - Raspberry Pi) 🚧"

# Cleanup
if (Test-Path "release") { Remove-Item "release" -Recurse -Force }
New-Item -ItemType Directory -Path "release" | Out-Null
New-Item -ItemType Directory -Path "release/backend" | Out-Null

# 1. Build Backend (Cross-Compile)
Write-Host "go: Building Backend..."
$env:GOOS = "linux"

# Build for Raspberry Pi (ARM64)
Write-Host "  -> Compiling for Raspberry Pi (ARM64)..."
$env:GOARCH = "arm64"
cd backend
go build -o ../release/backend/aar-daemon-pi main.go

# Build for Standard Laptop/PC (AMD64)
Write-Host "  -> Compiling for Standard PC/Laptop (AMD64)..."
$env:GOARCH = "amd64"
go build -o ../release/backend/aar-daemon-pc main.go
cd ..

# 2. Build Frontend
Write-Host "npm: Building Frontend..."
cd frontend
npm run build
cd ..

# 3. Package
Write-Host "pkg: Assembling Release..."
Copy-Item "frontend/dist" -Destination "release/backend/dist" -Recurse
Copy-Item "backend/deploy/*" -Destination "release/" -Recurse

Write-Host "✅ BUILD COMPLETE!"
Write-Host "Files are in the 'release/' folder."
Write-Host "To deploy:"
Write-Host "1. Copy 'release' folder to your Raspberry Pi."
Write-Host "2. Run 'sudo bash setup_router.sh'"
Write-Host "3. Run './backend/aar-daemon'"
