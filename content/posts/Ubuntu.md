---
title: "Ubuntu server 20.04自定义IP地址及其网关"
date: 2021-08-19T14:52:19+08:00
draft: true
toc: false
images:
tags: 
  - Linux
---

Ubuntu从20.04版本开始使用`netplan`对网络进行管理，百度上使用的`/etc/network/interfaces`文件已不再适用，需要修改`netplan`配置文件进行更新

## 修改配置文件
`netplan`的配置文件存放在`/etc/netplan`中，其中有诸如`00-installer-config.yaml`的配置文件，打开可以发现

```yaml
# This is the network config written by 'subiquity'
network:
  ethernets:
    # 网卡名称，我这是虚拟机ens33
    ens33:
      # 是否开启dhcp服务，设置为true后后续的ip和网关可以不用设置
      dhcp4: no
      # IP地址，一般设置一个就行
      # 后面的24是子网掩码数目，代表255.255.255.0(11111111.11111111.11111111.00000000)
      addresses: [192.168.1.106/24]
      optional: true
      # 网关地址
      gateway4: 192.168.1.253
      # dns地址
      nameservers:
        addresses: [233.5.5.5, 192.168.1.253]
  version: 2
```

当然，这只是最为基础的网络配置，对于多张网卡还可以进行其他的配置，详情可以参考[Netplan官方文档]("https://netplan.io/examples/")

在修改完文件之后，终端输入以下命令应用配置，Ping外网测试连通性

    sudo netplay apply
    ping 8.8.8.8