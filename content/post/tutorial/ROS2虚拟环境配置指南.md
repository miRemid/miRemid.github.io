---
title: "ROS2虚拟环境配置指南"
date: 2022-02-22T14:35:57+08:00
draft: false
toc: false
images:
tags: 
  - 教程
---

在使用CARLA的ROS2包的过程中，需要导入torch对图像进行处理，按照正常情况下，激活conda环境之后应该就能够直接使用虚拟环境中安装的第三方包，但是在ROS2启动的过程中却提示`torch not found`的错误

开始怀疑是conda环境导致的，于是更换为官方的venv，但是在安装完依赖后编译启动显示还是同样的错误。

在Github搜索后发现，原来通过colcon编译的ROS2节点在启动的时候默认使用系统环境的Python环境运行，并不会加载编译时使用的虚拟环境，因此在虚拟环境中安装的第三方依赖也不会正常使用了。根据[rotu](https://github.com/colcon/colcon-core/pull/183)的说法，由于ROS2节点的运行文件是由colcon编译生成，而在colcon编译过程中解释器的选择是写死为系统环境的Python的，那是不是ROS2就用不了虚拟环境运行呢？

当然不是，根据上个链接的说法，可以通过修改ROS2的配置文件来间接的让系统使用虚拟环境中的Python解释器，按照[theunkn0wn1](https://github.com/ros2/ros2/issues/1094#issuecomment-927179578)的描述，总共分为4步

1. 修改setup.cfg文件，添加如下配置
```
# src/{node}/setup.cfg
[build_scripts]
executable = /usr/bin/env python3
```
2. 创建虚拟环境
在创建虚拟环境时，需要将系统环境中的ROS2包链接到虚拟环境中
```
python3 -m venv venv --system-site-packages --symlinks
```
这时在目录中会生成venv文件

3. 激活ros2环境配置和虚拟环境配置
```
source /opt/ros/foxy/setup.zsh
source ./venv/bin/activate
```

4. colcon build并运行
```
colcon build
ros2 run package_name node_name
```

这时就可以使用到虚拟环境中的第三方包了
