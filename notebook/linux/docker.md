# 安装 docker
```
docker --version

sudo apt update
sudo apt install -y docker.io
sudo usermod -aG docker caiki  # 把 caiki 加入 docker 组，避免每次用 sudo

sudo systemctl enable docker # 开机自启
```