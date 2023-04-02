---
title: "使用Windows自带工具进行端口转发"
date: 2023-04-02T23:24:56+08:00
draft: false
tags: 
  - 教程
  - Windows
summary: netsh真香
---

由于Jetson TX2的内核版本低于5.15，不能直接使用wireguard进行组网，为了能够在外网的情况下通过wireguard连接到我的开发板，需要对开发板的ssh端口进行转发

有人会问，为什么不直接连接到一台Linux电脑然后直接进行ssh呢，那是因为我现在只有一台性能勉强够用的Windows小主机，总不可能为了这个功能而装一个虚拟机吧😂

谷歌了一下，发现Windows其实从XP版本开始就自带了一个CLI的端口转发工具`netsh`，使用方式非常简单，直接上CRUD

- `netsh interface portproxy show all` 查看所有转发规则
- `netsh interface portproxy dump` 查看转发设置
- `netsh interface portproxy delete v4tov4 listenport=$LOCALPORT listenaddress=$LOCALADDRESS`删除一个v4转发
- `netsh interface portproxy add v4tov4 listenaddress=0.0.0.0 listenport=8080 connectaddress=$REMOTE_ADDRESS connectport=$REMOTE_PORT`添加一个v4转发

非常方便😍