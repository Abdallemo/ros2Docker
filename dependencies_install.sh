#!/usr/bin/env bash

source /etc/os-release

dependencies=("curl" "docker")
curl_download_dir="/usr/bin"
docker_download_dir=/var/lib


# 📂 Check for dependencies
echo "📂 Checking for dependencies..."

if [[ "$ID_LIKE" =~ (ubuntu|debian) || "$ID" =~ (ubuntu|debian) ]]; then
    os="ubuntu"
elif [[ "$ID_LIKE" =~ (arch) || "$ID" =~ (arch) ]]; then
    os="arch"
    sh_pkg="sudo pacman -S --noconfirm"
    d_pkg="yay -S --noconfirm"
else
    os="unknown"
fi

for file in "${dependencies[@]}"; do
  
  if [[ -f "/usr/bin/$file" ]]; then
    echo "✅ Found: $file"
  else
    echo "❌ Missing: $file"

    if [ $file == 'curl' ];then
      echo "⬇ Downloading $file to $curl_download_dir"

        if [ $os == "arch" ]; then
          sudo pacman -Syu
          sudo pacman -S --noconfirm curl
        elif [ $os == 'ubuntu' ]; then 
          sudo apt-get update
          sudo apt-get install ca-certificates curl
        fi
    else

      echo "⬇ Downloading $file to $docker_download_dir"

      if [ $os = 'ubuntu' ]; then

        for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done

        sudo apt-get install ca-certificates curl
        sudo install -m 0755 -d /etc/apt/keyrings
        sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
        sudo chmod a+r /etc/apt/keyrings/docker.asc

        # Add the repository to Apt sources:
        echo \
          "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
          $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
          sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
          sudo apt-get update
          sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
      
      elif [ $os == 'arch' ]; then
        yay -Syu
        sudo pacman -S docker --noconfirm
        sudo systemctl enable docker.service
        sudo systemctl start docker.service

        yay -S docker-compose --noconfirm
      fi
    fi
  fi
  done


