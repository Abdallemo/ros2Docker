#!/bin/bash

set -e
echo "📥 Downloading ros2docker script..."

INSTALL_DIR="$HOME/.ros2docker"
BIN_DIR="$HOME/.local/bin"
EXECUTABLE_NAME="ros2dockr"
baseUrl="https://raw.githubusercontent.com/xaatim/ROS2-Docker-Launcher/main/bash--deprecated/src"
mkdir -p "$INSTALL_DIR"
curl -L -o "$INSTALL_DIR/$EXECUTABLE_NAME" "$baseUrl/$EXECUTABLE_NAME"
chmod +x "$INSTALL_DIR/$EXECUTABLE_NAME"

mkdir -p "$BIN_DIR"
ln -sf "$INSTALL_DIR/$EXECUTABLE_NAME" "$BIN_DIR/$EXECUTABLE_NAME"

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    echo "🧩 Adding $BIN_DIR to PATH"

    case "$SHELL" in
    */zsh)
        echo "⚙️  Detected Zsh"
        echo "export PATH=\"\$PATH:$BIN_DIR\"" >>"$HOME/.zshrc"
        echo "export DISPLAY=:0" >>"$HOME/.zshrc"
        echo "export WAYLAND_DISPLAY=wayland-1" >>"$HOME/.zshrc"
        echo "export XAUTHORITY=/tmp/.docker.xauth" >>"$HOME/.zshrc"
        ;;
    */bash)
        echo "⚙️  Detected Bash"
        echo "export PATH=\"\$PATH:$BIN_DIR\"" >>"$HOME/.bashrc"
        echo "export DISPLAY=:0" >>"$HOME/.bashrc"
        echo "export WAYLAND_DISPLAY=wayland-1" >>"$HOME/.bashrc"
        echo "export XAUTHORITY=/tmp/.docker.xauth" >>"$HOME/.bashrc"
        ;;
    *)
        echo "⚠️  Unknown shell. Please add $BIN_DIR to your PATH manually."
        echo "export DISPLAY=:0" >>"$HOME/.profile"
        echo "export WAYLAND_DISPLAY=wayland-1" >>"$HOME/.profile"
        echo "export XAUTHORITY=/tmp/.docker.xauth" >>"$HOME/.profile"
        ;;
    esac

    export PATH="$PATH:$BIN_DIR"
fi

echo "✅ ros2docker installed successfully!"
echo "🔄 Please restart your shell or run 'source ~/.bashrc' to update your PATH."
echo "💡 Try running: ros2dockr -clean"
