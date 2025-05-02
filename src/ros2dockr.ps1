$requiredFiles = @("compose.yaml", "Dockerfile")
$baseUrl = "https://raw.githubusercontent.com/xaatim/ROS2-Docker-Launcher/refs/heads/main"
$docDir = "$env:USERPROFILE\.ros2docker"
New-Item -ItemType Directory -Force -Path $docDir | Out-Null

Write-Host "📂 Checking for required files..."
foreach ($file in $requiredFiles) {
    if (Test-Path "$docDir\$file") {
        Write-Host "✅ Found: $file"
    }
    else {
        Write-Host "❌ Missing: $file"
        Write-Host "⬇️  Downloading $file..."
        try {
            Invoke-WebRequest -Uri "$baseUrl/$file" -OutFile "$docDir\$file" -UseBasicParsing
            Write-Host "✅ Downloaded: $file"
        }
        catch {
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
        docker compose -f "$docDir\compose.yaml" down
        docker compose -f "$docDir\compose.yaml" up --build -d
    }
    "-start" {
        Write-Host "🚀 Starting (without rebuild)"
        docker compose -f "$docDir\compose.yaml" up -d
    }
    "-stop" {
        Write-Host "🛑 Stopping services"
        docker compose -f "$docDir\compose.yaml" down
    }
    "-shell" {
        Write-Host "🔧 Entering container shell"
        docker exec -it ros2 bash
    }
    "-logs" {
        Write-Host "📜 Streaming logs"
        docker compose -f "$docDir\compose.yaml" logs -f
    }
    { $_ -eq "-h" -or $_ -eq "--help" } {
        Write-Host "🛠 ros2dock.ps1 usage:"
        Write-Host "  -clean     🔄 Rebuild image and restart container"
        Write-Host "  -start     🚀 Start container (without rebuild)"
        Write-Host "  -stop      🛑 Stop and remove container"
        Write-Host "  -shell     🔧 Enter interactive container shell"
        Write-Host "  -logs      📜 Follow container logs"
        Write-Host "  -h|--help  🆘 Show this help message"
    }
    default {
        Write-Host "ℹ️ Usage: ./ros2dock.ps1 [-help | -h]"
    }
}
