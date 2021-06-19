---
title: "利用fusuma为Linux添加触控板手势"
date: 2021-06-19T21:07:52+08:00
draft: false
toc: false
images:
tags: 
  - Linux
---

之前一直在新笔记本上用虚拟机开Linux编写毕设，现在毕业答辩已经完成，所以就直接将Win10更换到Linux系统。但是众所周知，Linux并没有自带触控板手势驱动，但我早已习惯于win上的一些手势，所以我开始寻找一个合适的触控板手势驱动。

在最初，我发现了一个名为`libinput-gesture`的库，其Github上有3k赞，于是我立马就尝试了这个库。安装过程一切畅通，但最后运行的阶段却出了问题。在启动该库时提示我启动失败，但并没有给其他错误信息。我检查了`xinput`和`libinput`对触控板的驱动，并没有任何问题。谷歌了一下发现这玩意对arch用户非常良好，但对我这使用的ubuntu来说会出现这种情况。于是我把目光转向另一个库，`fusuma`。

`fusuma`这个包非常完美的解决了Ubuntu触控板手势的问题，并且安装过程极其简单，配置也十分简洁明了。

## 安装

#### 配置用户权限
首先需要配置一下用户状态组，将用户添加到`input`组中

    sudo gpasswd -a $USER input
    newgrp input // 立刻刷新

#### 对于Debian系用户(Ubuntu,Debian,Mint,Pop!OS)
1. 安装libinput
```sh
sudo apt install libinput-tools
```
2. 安装Ruby

由于fusuma需要在ruby上运行，因此需要安装ruby。
```sh
sudo apt install ruby
```
3. 安装Fusuma
```sh
sudo gem install fusuma
```
4. 安装xdotool(可选，但是对后续非常有用建议安装)
```sh
sudo apt install xdotool
```

#### 对于Arch用户

    sudo pacman -S libinput ruby xdotool
    sudo gem install fusuma

#### 启动

    fusuma

## 配置
Fusuma的配置非常简单，首先需要在用户目录下创建配置文件`~/.config/fusuma/config.yml`

    mkdir -p $HOME/.config/fusuma
    touch $HOME/.config/fusuma/config.yml

这里简单贴出我使用的配置，可以在配置中发现，其实手势是通过`xdotool`工具间接调用系统的快捷键来实现手势控制的

```yml
swipe:
  3:
    left:
      command: "xdotool key alt+Right" # History forward
    right:
      command: "xdotool key alt+Left" # History back
      #up:
      #command: "xdotool key super" # Activity
      #down:
      #command: "xdotool key super" # Activity
  4:
    left:
      command: "xdotool key alt+u" # Switch to next workspace
    right:
      command: "xdotool key alt+d" # Switch to previous workspace
    up:
      command: "xdotool key alt+u" # Switch to next workspace
    down:
      command: "xdotool key alt+d" # Switch to previous workspace
pinch:
  in:
    command: "xdotool keydown ctrl click 4 keyup ctrl" # Zoom in
  out:
    command: "xdotool keydown ctrl click 5 keyup ctrl" # Zoom out
```

在添加完配置文件后，重新启动fusuma即可使用手势动作了

> 配置文件中的alt+u和alt+d为本人自定义快捷键，请更换为自己系统快捷键

## 开机自启动
Fusuma的开机自启动也比较简单，有两种方式

#### 通过桌面管理器自启
大部分用户使用的桌面都带有Autostart的功能，如官方提供的教程所说(Gnome用户)：
1. 使用`which fusuma`查看fusuma位置
2. 打开`gnome-session-properties`
3. 添加Fusuma相关信息
4. 在命令末尾添加`-d`选项

#### 通过systemd开机启动
这种方式适用于所有用户，但需要注意的是，由于触控板驱动是用户层使用的，在service文件中如果采取系统级别的启动需要额外添加环境配置信息，如果是以用户服务启动则无需添加

1. 系统服务
```sh
# /etc/systemd/system/fusuma.service
[Unit]
Description=Fusuma touchpad gestures

[Service]
EnvironmentFile="DISPLAY=:0"
ExecStart=/usr/local/bin/fusuma
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
Restart=on-failure
RestartSec=42s
User=youruser
Group=yourusergroup

[Install]
WantedBy=multi-user.target
```
> 参考[Github Fusuma Issue](https://github.com/iberianpig/fusuma/issues/52)

然后启动服务，并设置自启

    sudo systemctl enable fusuma
    sudo systemctl start fusuma
    # 查看启动状态
    sudo systemctl status fusuma

2. 用户服务
```sh
# $HOME/.config/systemd/user/fusuma.service
[Unit]
Description=Fusuma touchpad gestures

[Service]
ExecStart=/usr/local/bin/fusuma -d
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
然后启动服务，并设置自启

    systemctl --user enable fusuma
    systemctl --user start fusuma
    # 查看启动状态
    systemctl --user status fusuma

最后就可以通过手势来操控桌面了~