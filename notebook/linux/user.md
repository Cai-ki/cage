# 新建用户
```
# 基本用法（以 root 身份执行）
useradd -m username

# 为新用户设置密码
passwd username
```

# root 用户修改自己的密码
```
passwd
```

# root 修改其他用户的密码
```
passwd username
```

# 将用户加入 sudo
```
usermod -aG sudo username

su - username
sudo whoami
```