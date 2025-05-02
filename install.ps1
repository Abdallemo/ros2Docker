$ErrorActionPreference = "Stop"

Write-Host "📥 Downloading ros2docker script..."

$installDir = "$env:USERPROFILE\.config\.ros2docker"
$confDir = "$env:USERPROFILE\.config"
$executableName = "ros2dockr.ps1"
$targetPath = Join-Path $installDir $executableName
$shortcutPath = Join-Path $confDir "ros2docker.cmd"

New-Item -ItemType Directory -Force -Path $installDir | Out-Null
New-Item -ItemType Directory -Force -Path $confDir | Out-Null

Invoke-WebRequest -Uri "https://raw.githubusercontent.com/xaatim/ROS2-Docker-Launcher/refs/heads/main/src/ros2dockr.ps1" -OutFile $targetPath


Set-Content -Path $shortcutPath -Value "@echo off`nPowerShell -ExecutionPolicy Bypass -File `"$targetPath`" %*"

if (-not ($env:PATH -split ";" | Where-Object { $_ -eq $confDir })) {
    Write-Host "🧩 Adding $confDir to user PATH"
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$confDir", "User")
}

[Environment]::SetEnvironmentVariable("DISPLAY", "host.docker.internal:0.0", "User")
[Environment]::SetEnvironmentVariable("XAUTHORITY", "/tmp/.docker.xauth", "User")

Write-Host "✅ ros2docker installed successfully!"
Write-Host "🔄 Restart your shell or log out/in to use 'ros2docker'"
