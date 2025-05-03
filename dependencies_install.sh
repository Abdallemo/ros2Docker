#!/usr/bin/env bash

source /etc/os-release

dependencies=("docker" "curl")

# 📂 Check for dependencies
echo "📂 Checking for dependencies..."



for file in "${dependencies[@]}"; do
  if [[ -f "/usr/bin/$file" ]]; then
    echo "✅ Found: $file"
  else
    echo "❌ Missing: $file"
    exit 1
  fi
done