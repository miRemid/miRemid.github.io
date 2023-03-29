---
title: "WSL2更改桥接网络"
date: 2022-05-05T13:34:24+08:00
draft: false
toc: false
images:
tags: 
  - 教程
  - Linux
---

WSL2默认情况下网络是处于NAT模式之下，在正常情况下联网是没有问题的，但是遇到一些特殊情况下时NAT网络就非常的鸡肋。
例如，NAT的WSL2并不能直接使用宿主机的代理服务器，需要手动获取本机IP和宿主机在WSL2中的IP地址才能进行局域网连接。
在例如本人需要在WSL2中远程连接实验室车辆上的ROS2节点，由于NAT的存在导致小车上ROS2节点发布的Topic并不能在WSL2中接收到，因此我需要将WSL2默认的NAT网络
更改为桥接网络，这样在不影响宿主机上网的同时也能分配给WSL2一个局域网地址

> 本方法需要通过Hyper-V虚拟网卡进行，因此首先需要安装并开启Windows的Hyper-V功能

修改桥接网络步骤如下(以下步骤均在管理员模式的Powershell中进行)：
1. 从Powershell中开启wsl，生成网卡
```powershell
wsl.exe
```
如果只是开机自启的wsl会出现找不到WSL网卡的错误
2. 获取网卡信息
```powershell
Get-NetAdapter
```
运行之后可以查看当前Windows中的网卡设备，例如WLAN
3. 桥接网卡
选择WSL2需要桥接的网卡设备，以WLAN为例
```powershell
Set-VMSwitch WSL -NetAdapterName WLAN
```
4. 配置WSL2地址
这时，已经将WSL2的网卡桥接到物理网卡之上，需要手动配置WSL2的静态地址，假设路由器的网关为`192.168.8.1`，需要配置静态地址为`192.168.8.123`，则可以按照以下步骤进行配置(在WSL2中进行)
```shell
sudo ip addr flush dev eth0
sudo ip addr add 192.168.8.123/24 dev eth0
sudo ip route add 0.0.0.0/0 via 192.168.8.1 dev eth0
```
更改过后还需手动更改DNS地址
```shell
# /etc/resolv.conf
nameserver 192.168.8.1
```
保存并退出

以上就完成了WSL2的桥接网络配置，需要注意的是，每次都需要手动更新DNS文件，如果连接到别的路由器之下还需要手动更新静态地址和网关地址