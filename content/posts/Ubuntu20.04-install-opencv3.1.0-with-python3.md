---
title: "Ubuntu20.04 Install Opencv3.1.0 With Python3"
date: 2022-03-28T16:40:37+08:00
draft: false
toc: false
images:
tags: 
  - linux
---

# 安装依赖
Ubuntu20.04需要安装以下依赖

```bash
sudo apt install -y --no-install-recommends \ 
    build-essential \
    git \
    wget \
    unzip \
    yasm \
    pkg-config \
    libswscale-dev \
    libtbb2 \
    libtbb-dev \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    libopenjp2-7-dev \
    libavformat-dev \
    libpq-dev \
    libgtk2.0-dev libgtk-3-dev \
    libgphoto2-dev \
    libtiff5-dev libjpeg8-dev libpng-dev cmake make \
    libavformat-dev libavcodec-dev libswscale-dev libdc1394-22-dev libavresample-dev \
    libxine2-dev libv4l-dev \
    libatlas-base-dev \
    libfaac-dev libmp3lame-dev libtheora-dev \
    libvorbis-dev libxvidcore-dev \
    libeigen3-dev \
    libgstreamer1.0-0 gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \ 
    gstreamer1.0-plugins-ugly \ 
    gstreamer1.0-libav \
    gstreamer1.0-tools \
    python3.8-dev python3-pip \

pip3 install numpy
```

# 下载opencv3.1.0和contrib
假设下载目录为`OP_DIR=$HOME/opencv_project`
```bash
cd $OP_DIR
git clone https://github.com/opencv/opencv.git
cd opencv && git checkout 3.1.0 && cd ..
git clone https://github.com/opencv/opencv_contrib.git
cd opencv_contrib && git checkout 3.1.0 && cd ..
```

# 修复链接
如果直接cmake检查，在Ubuntu20.04会出现依赖错误
## 1. linux/videodev.h not found
```bash
sudo ln -sf /usr/include/libv4l1-videodev.h /usr/include/linux/videodev.h
```
## 2. ffmpeg/avformat.h not found
```bash
sudo mkdir -p /usr/include/ffmpeg 
sudo ln -sf /usr/include/x86_64-linux-gnu/libavcodec/*.h /usr/include/ffmpeg 
sudo ln -sf /usr/include/x86_64-linux-gnu/libavformat/*.h /usr/include/ffmpeg 
sudo ln -sf /usr/include/x86_64-linux-gnu/libswscale/*.h /usr/include/ffmpeg 
```
## 3. sys/videoio.h not found
```bash
sudo mkdir -p /usr/include/sys && sudo touch /usr/include/sys/videoio.h
```

# Cmake检查并编译安装
```
cd $OP_DIR
mkdir build && cd build
cmake \
-D BUILD_TIFF=ON \
-D BUILD_opencv_java=OFF \
-D WITH_CUDA=OFF \
-D WITH_OPENGL=ON \
-D WITH_OPENCL=ON \
-D WITH_IPP=ON \
-D WITH_TBB=ON \
-D WITH_EIGEN=ON \
-D WITH_V4L=ON \
-D WITH_LIBV4L=OFF \
-D BUILD_TESTS=OFF \
-D BUILD_PERF_TESTS=OFF \
-D CMAKE_BUILD_TYPE=RELEASE \
-D OPENCV_EXTRA_MODULES_PATH=$OP_DIR/opencv_contrib/modules \
-D ENABLE_PRECOMPILED_HEADERS=OFF \ # stdlib.io not found
-D PYTHON_EXECUTABLE=$(which python3) \
-D CMAKE_INSTALL_PREFIX=$(python3 -c "import sys; print(sys.prefix)") \
-D PYTHON_INCLUDE_DIR=$(python3 -c "from distutils.sysconfig import get_python_inc; print(get_python_inc())") \
-D PYTHON_PACKAGES_PATH=$(python3 -c "from distutils.sysconfig import get_python_lib; print(get_python_lib())") \
$OP_DIR/opencv \
make -j${nproc}
sudo make install
```