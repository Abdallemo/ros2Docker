# 🐳 ROS2 Docker Launcher

A simple cross-platform Docker Compose tool to streamline building and running ROS2 environments with graphical support on **Linux** and **Windows**.

# Installation
## for Linux/Mac run this command
``` bash
curl -sSL https://raw.githubusercontent.com/xaatim/ROS2-Docker-Launcher/refs/heads/main/install.sh | bash
```
## for Windows run this command
```bash
iwr -useb  https://raw.githubusercontent.com/xaatim/ROS2-Docker-Launcher/refs/heads/main/install.ps1 | iex

```

---

## 📁 Repository Structure

| File              | Purpose                                      |
|-------------------|----------------------------------------------|
| `compose.yaml`    | Docker Compose configuration                 |
| `Dockerfile`      | Container setup with ROS2                    |
| `.env`            | Environment variables for host display       |
| `ros2dock.sh`     | Linux launcher script                        |
| `ros2dock.ps1`    | Windows PowerShell launcher script           |

---

## ⚙️ Prerequisites

Make sure you have:

- [Docker](https://www.docker.com/)
- X11 server (e.g., [VcXsrv](https://sourceforge.net/projects/vcxsrv/) or [X410](https://x410.dev/)) for Windows

---

~~## 🖥️ Host Display Setup via `.env`~~

~~The `.env` file determines how graphical output (GUI apps) is forwarded from the container to the host.~~



## 🚀 Usage

From your terminal:

### 🐧 Linux/macOS (Bash)
```bash
chmod +x ros2dock.sh
ros2dockr -clean      # Clean, rebuild and restart
ros2dockr -start      # Start the container
ros2dockr -stop       # Stop and remove containers
ros2dockr -shell      # Open container shell
ros2dockr -logs       # View logs
ros2dockr -h          # for manual
```
### 🪟 Windows PowerShell
```
.\ros2dockr.ps1 -start
.\ros2dockr.ps1 -clean
.\ros2dockr.ps1 -stop
.\ros2dockr.ps1 -shell
.\ros2dockr.ps1 -logs

```

## ❓ What Each Command Does
| Command     | Description                              |
|-------------|------------------------------------------|
| `-start`    | Launch containers without rebuilding      |
| `-clean`    | Rebuilds the container and restarts       |
| `-stop`     | Stops and removes the containers          |
| `-shell`    | Enters the running container shell        |
| `-logs`     | Shows the live container logs             |

## 💡 Tips

- For **Windows users**, start your X11 server **before** launching the container.
- If using **Wayland**, ensure your session supports `xhost` or use tools like [XWayland](https://wiki.archlinux.org/title/XWayland).
- This setup is ideal for both GUI-based and CLI-based  ROS2 tools  (like `rqt`, `gazebo`, etc.).

---

## 📎 Related Files

- [`compose.yaml`](https://github.com/xaatim/ROS2-Docker-Launcher/compose.yaml)
- [`Dockerfile`](https://github.com/xaatim/ROS2-Docker-Launcher/Dockerfile)
- [`ros2dock.sh`](https://github.com/xaatim/ROS2-Docker-Launcher/ros2dock.sh)
- [`ros2dock.ps1`](https://github.com/xaatim/ROS2-Docker-Launcher/ros2dock.ps1)
~~[`.env`](https://github.com/Abdallemo/ros2Docker/blob/main/.env)~~

---

