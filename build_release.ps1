Write-Host "🚧 Building AAR for Production (Linux/ARM64 - Raspberry Pi) 🚧"

# 1. Build Backend
Write-Host "Building Go Backend..."
$env:GOOS = "linux"
$env:GOARCH = "arm64"
cd backend
go build -o aar-daemon-linux-arm64 main.go
cd ..

Write-Host "✅ Backend Build Complete: backend/aar-daemon-linux-arm64"
Write-Host "To deploy via Docker, transfer files to your Pi and run:"
Write-Host "  docker-compose -f docker-compose.prod.yml up --build -d"
