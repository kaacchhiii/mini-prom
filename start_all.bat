@echo off
start "Mini Prometheus Target" cmd /k "go run target/main.go"
timeout /t 2 >nul
start "Mini Prometheus Server" cmd /k "go run main.go"
echo Services started!
echo Target running on http://localhost:8081/metrics
echo Server running on http://localhost:9090
