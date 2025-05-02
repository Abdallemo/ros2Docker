#!/usr/bin/env bash

source /etc/os-release

required_files=("compose.yaml" "Dockerfile" ".env")

# 📂 Check for required files
echo "📂 Checking for required files..."
for file in "${required_files[@]}"; do
  if [[ -f "$file" ]]; then
    echo "✅ Found: $file"
  else
    echo "❌ Missing: $file"
    echo "⬇️ Downloading $file"
    # as you know i want to download it in the userProject
    curl -L -o $file https://raw.githubusercontent.com/xaatim/ROS2-Docker-Launcher/refs/heads/main/ros2dockr$file
    
    
  fi
done


if ! command -v xhost &>/dev/null; then
  echo "❌ 'xhost' command not found. Please install it."

  if [[ "$ID_LIKE" =~ (ubuntu|debian) || "$ID" =~ (ubuntu|debian) ]]; then
    echo "🧩 Installing on Debian-based system..."
    sudo apt install x11-xserver-utils
  elif [[ "$ID_LIKE" =~ (arch) || "$ID" =~ (arch) ]]; then
    echo "🧩 Installing on Arch-based system..."
    sudo pacman -S xorg-xhost
  else
    echo "⚠️ Unsupported distro. Manual installation required."
  fi

  exit 1
fi

if [[ -z "$DISPLAY" || ! -S /tmp/.X11-unix/X0 ]]; then
  echo "❌ X11 display is not running. Please start your graphical session."
  exit 1
fi

# 🚦 Command line options
case "$1" in
  -clean)
    echo "♻️  Rebuilding and cleaning volumes"
    if docker ps -a --format '{{.Names}}' | grep -q '^ros2$'; then
      echo "🧹 Removing old container 'ros2'"
      docker rm -f ros2
    fi
    xhost +local:root
    docker-compose down 
    docker-compose up --build -d
    ;;
  -start)
    echo "🚀 Starting (without rebuild)"
    xhost +local:root
    docker-compose up -d
    ;;
  -stop)
    echo "🛑 Stopping services"
    docker-compose down
    ;;
  -shell)
    echo "🔧 Entering container shell"
    docker exec -it ros2 bash
    ;;
  -logs)
    echo "📜 Streaming logs"
    docker-compose logs -f
    ;;
  *)
    echo "ℹ️ Usage: $0 {-clean|-start|-stop|-shell|-logs}"
    ;;
esac

