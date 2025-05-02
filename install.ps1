$ErrorActionPreference = "Stop"

Write-Host "📥 Downloading ros2docker script..."

$installDir = "$env:USERPROFILE\.config\.ros2docker"
$confDir = "$env:USERPROFILE\.config"
$executableName = "ros2dockr.ps1"
$targetPath = Join-Path $installDir $executableName
$shortcutPath = Join-Path $confDir "ros2docker.cmd"

New-Item -ItemType Directory -Force -Path $installDir | Out-Null
New-Item -ItemType Directory -Force -Path $confDir | Out-Null

Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Abdallemo/ros2Docker/main/ros2dockr" -OutFile $targetPath


Set-Content -Path $shortcutPath -Value "@echo off`nPowerShell -ExecutionPolicy Bypass -File `"$targetPath`" %*"

if (-not ($env:PATH -split ";" | Where-Object { $_ -eq $confDir })) {
    Write-Host "🧩 Adding $confDir to user PATH"
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$confDir", "User")
}

Write-Host "✅ ros2docker installed successfully!"
Write-Host "🔄 Restart your shell or log out/in to use 'ros2docker'"
