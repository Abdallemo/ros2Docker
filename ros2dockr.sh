#!/usr/bin/env bash

set -e
export PWD="$(pwd)"
source /etc/os-release

docDir="$HOME/.ros2docker"
mkdir -p "$docDir"

required_files=("compose.yaml" "Dockerfile")

echo "📂 Checking for required files in $docDir..."
for file in "${required_files[@]}"; do
  if [[ -f "$docDir/$file" ]]; then
    echo "✅ Found: $file"
  else
    echo "❌ Missing: $file"

    echo "⬇ Downloading $file to $docDir"
    curl -L -o "$docDir/$file" "https://raw.githubusercontent.com/xaatim/ROS2-Docker-Launcher/refs/heads/main/ros2dockr/$file"

  fi
done

if ! command -v xhost &>/dev/null; then
  echo "❌ 'xhost' command not found. Please install it."

  if [[ "$ID_LIKE" =~ (ubuntu|debian) || "$ID" =~ (ubuntu|debian) ]]; then
    echo "🧩 Installing on Debian-based system..."
    sudo apt install -y x11-xserver-utils
  elif [[ "$ID_LIKE" =~ (arch) || "$ID" =~ (arch) ]]; then
    echo "🧩 Installing on Arch-based system..."
    sudo pacman -S --noconfirm xorg-xhost
  else
    echo "⚠ Unsupported distro. Manual installation required."
  fi

  exit 1
fi

if [[ -z "$DISPLAY" || ! -S /tmp/.X11-unix/X0 ]]; then
  echo "❌ X11 display is not running. Please start your graphical session."
  exit 1
fi

case "$1" in
-clean)
  echo "♻ Rebuilding and cleaning volumes"
  if docker ps -a --format '{{.Names}}' | grep -q '^ros2$'; then
    echo "🧹 Removing old container 'ros2'"
    docker rm -f ros2
  fi
  xhost +local:root
  docker-compose -f "$docDir/compose.yaml" down
  docker-compose -f "$docDir/compose.yaml" up --build -d
  ;;
-start)
  echo "🚀 Starting (without rebuild)"
  xhost +local:root
  docker-compose -f "$docDir/compose.yaml" up -d
  ;;
-stop)
  echo "🛑 Stopping services"
  docker-compose -f "$docDir/compose.yaml" down
  ;;
-shell)
  echo "🔧 Entering container shell"
  docker exec -it ros2 bash
  ;;
-logs)
  echo "📜 Streaming logs"
  docker-compose -f "$docDir/compose.yaml" logs -f
  ;;
-h | --help)
  echo "🛠 ros2docker usage:"
  echo "  -clean     🔄 Rebuild image and restart container"
  echo "  -start     🚀 Start container (without rebuild)"
  echo "  -stop      🛑 Stop and remove container"
  echo "  -shell     🔧 Enter interactive container shell"
  echo "  -logs      📜 Follow container logs"
  echo "  -h|--help  🆘 Show this help message"
  ;;
*)
  echo "ℹ Usage: $0 {-clean|-start|-stop|-shell|-logs}"
  ;;
esac
