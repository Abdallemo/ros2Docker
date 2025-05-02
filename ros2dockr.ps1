
$requiredFiles = @("compose.yaml", "Dockerfile", ".env")
$baseUrl = "https://raw.githubusercontent.com/xaatim/ROS2-Docker-Launcher/refs/heads/main"
Write-Host "📂 Checking for required files..."
foreach ($file in $requiredFiles) {
    if (Test-Path $file) {
        Write-Host "✅ Found: $file"
    } else {
        Write-Host "❌ Missing: $file"
        Write-Host "⬇️  Downloading $file..."
        try {
            Invoke-WebRequest -Uri "$baseUrl/$file" -OutFile $file -UseBasicParsing
            Write-Host "✅ Downloaded: $file"
        } catch {
            Write-Error "❌ Failed to download $file"
            exit 1
        }
    }
}
$Command = $args[0]

switch ($Command) {
    "-clean" {
        Write-Host "♻️  Rebuilding and cleaning volumes"
        if (docker ps -a --format '{{.Names}}' | Select-String -Pattern "^ros2$") {
            Write-Host "🧹 Removing old container 'ros2'"
            docker rm -f ros2
        }
        docker compose down
        docker compose up --build -d
    }
    "-start" {
        Write-Host "🚀 Starting (without rebuild)"
        docker compose up -d
    }
    "-stop" {
        Write-Host "🛑 Stopping services"
        docker compose down
    }
    "-shell" {
        Write-Host "🔧 Entering container shell"
        docker exec -it ros2 bash
    }
    "-logs" {
        Write-Host "📜 Streaming logs"
        docker compose logs -f
    }
    default {
        Write-Host "ℹ️ Usage: ./ros2dock.ps1 [-clean|-start|-stop|-shell|-logs]"
    }
}
