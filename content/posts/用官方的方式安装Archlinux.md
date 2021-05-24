---
title: "用官方的方式安装Archlinux"
date: 2021-05-23T22:47:09+08:00
draft: false
toc: true
images:
tags: 
  - Archlinux
---

Archlinux是一个非常干净的Linux发行版，干净到以致于图形化界面都没有，本文将按照[官方wiki](https://wiki.archlinux.org/title/Installation_guide)，逐步安装Archlinux

由于国内网络问题，可以到一些镜像网站中来下载iso镜像，这里选择[清华镜像](https://mirrors.tuna.tsinghua.edu.cn/archlinux/iso/2021.04.01/)，下载最新的5.15版本

> 关于启动盘，可以使用[USBWriter](https://sourceforge.net/projects/usbwriter/)或者[Rufus](https://rufus.ie)来制作，这里就不多加赘述

## 准备过程

### 检查引导模式

要安装Arch，首先要确认自己电脑的引导模式，可以通过以下指令来检查

    ls /sys/firmware/efi/efivars

如果该指令没有出错并列出了文件列表，则你的电脑是使用UEFI模式引导的。如果文件不存在，则系统将会使用BIOS(或CSM)模式来引导启动。如果电脑并不是以你想要的方式启动，则需要到主板BIOS中手动设置引导模式

### 连接网络

Arch的安装一定需要联网，因此在确认引导模式之后，必须要检查网络连接状态。你可以通过以下步骤来检查自己的网络状态:
- 确保网卡正确插入，例如通过`ip link`来检查
- 对于无线用户，确保无线网卡工作正常，可以通过`rfkill list`来检查
- 连接到互联网
    - 有线用户：插入网线
    - 无线用户：通过`iwctl`接入WIFI
    - 移动热点：通过`mmcli`连接热点
- 配置网络连接
    - DHCP：动态IP、DNS服务，适用于有线网络、无线网络……
    - 静态IP地址
- 确认网络连接`ping www.baidu.com`

### 更改时间配置
首先通过`date -R`来检查时间是否正确，如果不正确则可以通过`timedatectl`来更改时区

- 首先通过`timedatectl list-timezones`来获取时区列表
- 选择自己所在时区，以国内为例，通过`timedatectl set-timezone Asia/Shanghai`来更改时区
- 同步时区`timedatectl set-ntp true`

### 分区
电脑中的硬盘将会以设备块的形式出现在系统之中，例如`/dev/sda`,`/dev/nvmeOn1`，可以通过`fdisk`来查看

    fdisk -l

以下列表的分区是必要的
- 用于存放根目录的`/`分区
- 对`UEFI`模式的用户来说还需要一个`EFI system partition`

对于`BIOS`用户，需要将`root`分区挂载到`/mnt`中，而`UEFI`还需要将`EFI system partition`挂载到`/mnt/boot`或`/mnt/efi`分区之中

假设当前为`BIOS`模式，且硬盘为`/dev/sda`大小为15G，我们可以用`fdisk`来对硬盘进行分区

    fdisk /dev/sda # 进入分区模式
    g # 创建分区表
    n # 创建swap分区
    # 选择任意一个编号（通常默认）
    # 选择起始区域（通常默认）
    +2G # 
    n # 创建根目录分区
    # 选择任意一个编号（通常默认）
    # 选择起始区域（通常默认）
    +13G
    t # 更改分区类型
    1 # 选择swap分区
    L # 查看所有类型
    19 # 选择Linux swap

    # 如果是efi用户，还需要更改efi分区类型

### 格式化
分区完后，需要将所有分区进行格式化。通常来说，分区时efi分区为第一个区，其次是swap区然后是根目录区，按照这个顺序以此格式化。在这里，只有两个分区分别为
- /dev/sda1: swap
- /dev/sda2: 根目录

因此格式化操作如下

    mkswap /dev/sda1 
    swapon /dev/sda1
    mkfs.ext4 /dev/sda2

### 挂载

分区完之后，将硬盘挂载到`/mnt`目录中，efi用户还需要将引导分区挂载到`/mnt/boot`中

    mount /dev/sda2 /mnt

## 正式安装

### 更换镜像源

在挂载完文件系统后，就是Archlinux的正式安装了。不过由于网络问题，国内推荐更换镜像源来加快安装速度

使用`vim`或`nano`打开`/etc/pacman.d/mirrorlist`，将国内镜像地址填写在文件的最开头处，以清华镜像为例

    Server = https://mirrors.tuna.tsinghua.edu.cn/archlinux/$repo/os/$arch

### 基本安装
使用`pacstrap`将基础包安装进文件系统之中，官方要求安装以下三个

    pacstrap /mnt base linux linux-firmware

但是，为了后续系统顺利使用，还需要安装文本处理工具和网络连接工具和dhcpcd，如果需要编译还可以安装base-devel，这些都可以按需进行安装

    pacstrap /mnt vim dhcpcd base-devel sudo

## 配置系统

### Fstab
生成自动挂载分区文件，可以使用`fstab`进行

    genfstab -U /mnt >> /mnt/etc/fstab

接下来就需要在新系统内部进行操作了

    arch-chroot /mnt

### 配置时间
新系统需要定义好时间信息，在这选择`Asia/Shanghai`

    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
    hwclock --systohc # 生成/etc/adjtime

### 位置信息
我们需要对Arch系统的地理位置信息进行更改，通常使用`en_US.UTF-8`。修改`/etc/locale.gen`将含有`zh_CN.UTF-8 UTF-8`、`zh_HK.UTF-8 UTF-8`、`zh_TW.UTF-8 UTF-8`、`en_US.UTF-8 UTF-8`的字段去除注释后保存

运行下面的指令配置地理位置文件

    locale-gen

随后创建`locale.conf`文件，并配置`LANG`变量

    # /etc/locale.conf
    LANG=en_US.UTF-8

### 网络配置
创建`hostname`文件，填入你想要的主机名称

    # /etc/hostname
    archlinux

然后打开`hosts`文件，填入以下信息

    # /etc/hosts
    127.0.0.1   localhost
    ::1         localhost
    127.0.1.1   archlinux.localdomain archlinux

如果主机拥有永久性IP地址，可以替换`127.0.1.1`

### 配置Root密码

    passwd

### 添加用户
一般来说不会使用Root用户来运行系统，因为权限太高了，因此需要创建一个普通用户来使用Archlinux

可以通过`useradd`来创建一个用户，并通过`passwd`来更改密码

    useradd -m coder
    passwd coder

在创建完用户之后，还需要将用户添加到`wheel`用户组中，这样后续就能够使用`sudo`来运行较高权限的指令了

    usermod -aG wheel,audio,video,optical,storage coder

随后，通过`visudo`修改文件中的用户组部分，找到`wheel`那一行并删除注释，保存退出

    %wheel ALL=(ALL) ALL

### 安装Intel-ucode(非Intel CPU跳过)
这个包是用来优化intel的，其他cpu可以查看wiki

    pacman -S intel-ucode

### 安装Bootloader
这里选择安装目前最为流行的`grub2`，在安装之前可以安装另外两个包来自动设置启动选项

    pacman -S os-prober ntfs-3g

如果你使用的是BIOS引导

- 安装grub

        pacman -S grub

- 部署grub

        grub-install --target=i386-pc /dev/sda

- 生成配置文件

        grub-mkconfig -o /boot/grub/grub.cfg

如果你使用UEFI引导

- 安装grub和efibootmgr等

        pacman -S grub efibootmgr dosfstools mtools
    
- 挂载efi分区

        mkdir /boot/efi
        mount /dev/sda1(你的efi分区) /boot/efi

- 部署grub

        grub-install --target=x86_64-efi --bootloader-id=grub_uefi --recheck

- 生成配置文件

        grub-mkconfg -o /boot/grub/grub.cfg

到此为止，Arch的安装基本完成，`exit`退出后，关机，拔出引导U盘后重启即可进入系统

### AUR
除了pacman之外，还有第三方软件包`aur`可以供arch用户使用，这里使用[yay](https://github.com/Jguer/yay)工具来使用`aur`

首先下载`yay`的`PKGBUILD`文件

    sudo pacman -S git
    git clone https://aur.archlinux.org/yay.git
    cd yay

由于网络问题，在`clone`和编译安装`yay`时会非常缓慢，这里可以修改`PKGBUILD`文件和设置环境变量来解决

修改`PKGBUILD`文件，将其中的`https://github.com/Jguer/yay`地址替换为`https://hub.fastgit.org/Jguer/yay`

然后执行安装指令

    GO111MODULE=on GOPROXY=https://goproxy.cn makepkg -si
    # 如果不想保留golang则
    # GO111MODULE=on GOPROXY=https://goproxy.cn makepkg -sir

安装完`yay`后，添加国内aur`archlinuxcn`。修改`/etc/pacman.conf`文件，在最后添加

    # /etc/pacman.conf
    [archlinuxcn]
    Server = https://mirrors.tuna.tsinghua.edu.cn/archlinuxcn/$arch

然后更新软件源并安装GPG key

    sudo pacman -Syy && sudo pacman -S archlinuxcn-keyring

如果遇到下述错误

    ERROR: 5984EA8F3C could not be locally signed

可以通过清除缓存来解决

    sudo rm -rf /etc/pacman.d/gnupg
    sudo pacman-key --init
    sudo pacman-key --populate archlinux

## 桌面环境

### 图形驱动
按照[官方Wiki](https://wiki.archlinux.org/index.php/Xorg#Driver_installation)的表格自行选择驱动安装，例如本人是Intel集显，则执行

    sudo pacman -S xf86-xvideo-intel

对于Nividia独显驱动，由于某原因并不推荐安装，如果要安装则在上述链接中自行查询

### X11
安装完图形驱动后，就可以正式安装桌面环境了，这里选择`X11`，直接在终端执行

    sudo pacman -S xorg xorg-init

安装完之后，选择一个你喜欢的桌面例如`Xfce`、`KDE`、`Gnome`等等，这里我选择的是桌面管理器`sway`

    sudo pacman -S sway alacritty adobe-source-code-pro-fonts dmenu

安装完之后，需要配置一下`X11`的启动项，将初始配置文件复制到用户目录中，然后在文件末尾添加`exec sway`的选项

    cp /etc/X11/xinit/xinitrc $HOME/.xinitrc

最后在终端中输入`startx`即可进入`sway`桌面了，输入`win+shift+enter`即可打开终端

具体`sway`的配置，可以从`etc/sway/config`中查看，将该文件复制到`$HOME/.config/sway/config`就可以自行修改配置`sway`了，例如添加

    font pongo:SourceCodePro Medium 12

即可修改`sway`桌面的字体和大小了

如果想在登录后自动启动桌面，则可以在终端配置文件中如bash终端的`.bash_profile`或zsh的`.zshrc`中，添加下述内容即可

    [[ $(fgconsole 2>/dev/null)==1 ]] && exec startx --vt1