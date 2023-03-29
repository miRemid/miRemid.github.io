---
title: "Jetson TX2 NX 配置笔记"
date: 2023-03-05T22:01:11+08:00
draft: false
toc: false
images:
  - https://images.unsplash.com/photo-1653038282408-13b605af0ef7?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=774&q=80
tags: 
  - Jetson
  - Linux
---
> Jetson真的好贵啊> <

{{< figure src="/images/jetsontx2nx.jpg" alt="image" caption="tx2nx开发板" >}}

<!-- 原由 -->
由于科研的需要，需要在边缘设备上对深度学习模型进行量化测试。实验室内刚好有闲置的Nvidia Jetson TX2开发板，所以就直接拿来测试了

原本以为在开发板上的配置会比较简单，无非就是从x86平台转换到arm上面，但实际上坑非常多例如前面提到的arm架构，其实在Jetson上面压根不是传统的arm架构，而是从armv8独立分支出来的aarch架构

虽然坑非常多，但总体上配置流程也非常简洁，基本配置过程也和x86平台上面大致相同，只是在安装的时候需要注意。本文大纲如下
- 烧录安装系统
- 安装conda虚拟环境
- 安装Pytorch和torchvision
- 安装torch2trt

<!-- 内容 -->
## 安装系统

首先需要烧录镜像到开发板中，由于我这块板子在购买的时候代理商就已经烧录好了系统到固件中，因此直接开机即可

如果是裸板安装，则可以参考知乎的教程进行烧录安装

连接电源启动开发板，在安装界面中可以选择启用CPU核心的选项，默认是4核，对应的SWAP大小为2GB，可以按照自己的需求进行更改。当然，如果已经安装完了系统也可以在系统中进行修改。

tx2自带的存储大小非常小只有16GB，按需求进行扩容，我这里是外接了一块128G的SSD硬盘

## 安装conda虚拟环境

按照本人的习惯，需要一个虚拟环境用于区分不同项目的环境，在tx2也不例外。这里就出现了第一个坑，那就是官方的conda并不能在tx2上完美运行！
就算你下载的是aarch架构的脚本，你在安装过程中也可能报错，又或者即使安装成功你在创建虚拟环境选择低版本的Python时也会出现错误

为了解决这个问题，我们需要使用一个修改版的conda，那就是[Mambaforge](https://github.com/conda-forge/miniforge)。从官方的[Release](https://github.com/conda-forge/miniforge/releases)中下载所需要的脚本文件(**注意得是aarch架构的文件！**)，安装方式则和普通版本的conda如出一辙，需要注意的是安装的位置需要选择一个容量更大的硬盘上面

```shell
wget https://github.com/conda-forge/miniforge/releases/download/22.11.1-4/Mambaforge-22.11.1-4-Linux-aarch64.sh
sh Mambaforge-22.11.1-4-Linux-aarch64.sh
```

随后创建你的第一个环境吧

```shell
conda create -n hello python=3.6
conda activate hello
```

## 安装Pytorch和torchvision

在安装Pytorch前，你需要安装Jetson的CUDA环境，默认情况下是没有的，当然在Nvidia自家产品上安装CUDA环境非常简单，不像Linux那样惹人厌。安装只需要一条命令即可
```shell
sudo apt install nvidia-jetpack
```
这时，查看`/usr/local`中应该就会出现对应的CUDA环境了，Jetpack会帮助你安装CUDA工具包、TensorRT和其他工具，你可以通过`jetson_release`来查看当前安装的CUDA版本
```shell
(base) tx2-1@tx2-1-desktop:/usr/local$ jetson_release
Software part of jetson-stats 4.1.5 - (c) 2023, Raffaello Bonghi
Model: lanai-3636 - Jetpack 4.6.1 [L4T 32.7.1]
NV Power Mode: MAXP_CORE_ARM - Type: 3
jtop:
 - Version: 4.1.5
 - Service: Active
Libraries:
 - CUDA: 10.2.300
 - cuDNN: 8.2.1.32
 - TensorRT: 8.2
 - VPI: 1.2.3
 - OpenCV: 4.1.1 - with CUDA: NO
```
此时，你也可以将CUDA的环境添加到终端中，为后续编译做准备
```shell
# $HOME/.bashrc
export CUDA_HOME=/usr/local/cuda
export PATH=$PATH:$CUDA_HOME/bin
```

接下来就要开始安装Pytorch了，注意的是在Jetson设备上并不能像传统conda环境那样直接通过conda或pip安装，需要到Nvidia官方下载对应的安装包。而torchvision也需要自己手动编译安装或下载第三方编译好的包进行安装。以下流程均来自[官方教程](https://forums.developer.nvidia.com/t/pytorch-for-jetson/72048)

1. 安装Pytorch
在这里，由于虚拟环境选择的Python版本为3.6，而最高支持3.6版本的Pytorch版本为1.8.0，所以这里安装1.8.0版本
```shell
wget https://nvidia.box.com/shared/static/p57jwntv436lfrd78inwl7iml6p13fzh.whl -O torch-1.8.0-cp36-cp36m-linux_aarch64.whl
sudo apt-get install python3-pip libopenblas-base libopenmpi-dev libomp-dev
conda activate hello
pip install cython numpy==1.19.4 torch-1.8.0-cp36-cp36m-linux_aarch64.whl
```
注意，这里`numpy`选择安装`1.19.4`版本，默认的`1.19.5`版本会出现莫名其妙的bug

2. 安装torchvision  
安装torchvison需要注意所安装的torch版本，对应列表如下
```
PyTorch v1.0 - torchvision v0.2.2
PyTorch v1.1 - torchvision v0.3.0
PyTorch v1.2 - torchvision v0.4.0
PyTorch v1.3 - torchvision v0.4.2
PyTorch v1.4 - torchvision v0.5.0
PyTorch v1.5 - torchvision v0.6.0
PyTorch v1.6 - torchvision v0.7.0
PyTorch v1.7 - torchvision v0.8.1
PyTorch v1.8 - torchvision v0.9.0
PyTorch v1.9 - torchvision v0.10.0
PyTorch v1.10 - torchvision v0.11.1
PyTorch v1.11 - torchvision v0.12.0
PyTorch v1.12 - torchvision v0.13.0
```
我们先看看自行编译安装方式，首先确认所需要的版本，这里为`v0.9.0`，然后按照下述过程进行安装
```shell
sudo apt-get install libjpeg-dev zlib1g-dev libpython3-dev libavcodec-dev libavformat-dev libswscale-dev
git clone --branch 0.9.0 https://github.com/pytorch/vision torchvision
cd torchvision
export BUILD_VERSION=0.9.0
python3 setup.py install --user
```
如果你一切顺利的话，通过`pip list | grep torch`则可以看到对应的安装文件了，但是本人运气不咋行，在编译中遇到`nvcc`的错误，什么找不到目标文件之类的，但是也没有其他错误信息，让我非常头疼

幸运的是，在谷歌上搜到别人为nano编译好了的torchvision。虽然是在nano编译的，但是都是aarch架构应该可以用吧，实际上也确实可以用，[参考链接](https://qengineering.eu/install-pytorch-on-jetson-nano.html)

```shell
sudo apt-get install libjpeg-dev zlib1g-dev libpython3-dev
sudo apt-get install libavcodec-dev libavformat-dev libswscale-dev
pip install pillow gdown
gdown https://drive.google.com/uc?id=1BdvXkwUGGTTamM17Io4kkjIT6zgvf4BJ
pip install torchvision-0.9.0a0+01dfa8e-cp36-cp36m-linux_aarch64.whl
pip list | grep torch
torch               1.8.0
torch2trt           0.4.0 # 这个后面会进行安装
torchvision         0.9.0a0+01dfa8e
```
ok，这就安装完事了，非常简单

最后打开终端测试是否安装成功
```python
import torch
import torchvision
torch.cuda.is_avaliable()
>>> True
```

## 安装torch2trt

安装torch2trt的过程非常简单，但也有两个点需要注意

首先，我们可以发现在上面安装`jetpack`的时候可以发现已经安装好了`TensorRT`（甚至`docker`），但是安装包并不存在于虚拟环境中，而我们并不想要在重新安装一遍这玩意，因此需要对其进行复用，不然后续编译`torch2trt`时会直接报错（

```shell
export PYTHONPATH=$PYTHONPATH:/usr/lib/python3.6/dist-packages # 目录按照实际情况更改
```
这时你再到虚拟环境的终端中测试python是否能正常导入tensorrt
```python
import tensorrt
```
如果一切正常就可以进行下一步的安装了，首先需要clone所需要的仓库

```shell
git clone https://github.com/NVIDIA-AI-IOT/torch2trt
cd torch2trt
```

在这里我们直接选择最新版本，避免和tensorrt新api冲突

```shell
git checkout v0.4.0
```

随后就可以直接进行安装了，如果安装过程中碰到缺少包的情况直接pip安装即可
```shell
python setup.py install
pip list | grep torch
torch               1.8.0
torch2trt           0.4.0 # 这个后面会进行安装
torchvision         0.9.0a0+01dfa8e
```

<!-- 结尾 -->
## 结语

总的来说Tx2的配置没有我想象中的简单但也没有那么复杂，Nvidia社区中有大量的解决方案（虽然nano居多）仔细谷歌也都可以寻找到答案。我在测试了几个开源的项目后发现TX2的性能也还不错，不过CPU非常垃圾，后续也会写写关于开源项目如何在TX2上测试的例子，更好的理解和熟悉这块开发板




















