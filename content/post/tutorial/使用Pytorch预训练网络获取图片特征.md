---
title: 使用Pytorch预训练网络获取图片特征
date: 2021-12-17T17:52:42+08:00
draft: true
toc: true
images:
tags: 
  - pytorch
---

由于Carla训练过程中需要使用到图像数据，而通常深度强化学习(DRL)算法中，输入的数据都为一个
向量(Vector)，如果需要使用图像这种二维数据，通常的做法是使用一个CNN网络，提取图像的特征，
将二维数据转化为一维向量数据，这样就能够直接使用现有的DRL算法

例如在[stable_baseline3](https://github.com/DLR-RM/stable-baselines3)中，支持MLP、CNN和
多种不同类型观测对象混合输入。对于图像数据，都使用了一个基于CNN深度卷积网络来提取特征。

在pytorch中已经有了很多优秀的CNN特征提取网络，在这里我将会使用到`MobileNetV3`这个网络用
于提取Carla的摄像头图片特征。

# 读取模型

首先需要安装`torchvision`这个包，里面包含了现有非常流行的图像处理网络

```python
pip install torchvision
```

然后打开一个终端，导入`torchvision.models`模型包

```python
from torchvision import models
```

# 导入MobileNetV3并初始化

Pytorch为每个网络都提供了一个预训练好的结构，可以自动下载导入

```
mobilenet = models.mobilenet.mobilenet_v3_small(pretrained=True)
mobilenet.eval()
```
下面就是MobileNetV3的网络结构，从最终输出层可以得知，MobileNetV3会提取1000个特征值
```shell
MobileNetV3(
  (features): Sequential(
    (0): ConvNormActivation(
      (0): Conv2d(3, 16, kernel_size=(3, 3), stride=(2, 2), padding=(1, 1), bias=False)
      (1): BatchNorm2d(16, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
      (2): Hardswish()
    )
    (1): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(16, 16, kernel_size=(3, 3), stride=(2, 2), padding=(1, 1), groups=16, bias=False)
          (1): BatchNorm2d(16, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): ReLU(inplace=True)
        )
        (1): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(16, 8, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(8, 16, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (2): ConvNormActivation(
          (0): Conv2d(16, 16, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(16, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (2): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(16, 72, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(72, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): ReLU(inplace=True)
        )
        (1): ConvNormActivation(
          (0): Conv2d(72, 72, kernel_size=(3, 3), stride=(2, 2), padding=(1, 1), groups=72, bias=False)
          (1): BatchNorm2d(72, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): ReLU(inplace=True)
        )
        (2): ConvNormActivation(
          (0): Conv2d(72, 24, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(24, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (3): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(24, 88, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(88, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): ReLU(inplace=True)
        )
        (1): ConvNormActivation(
          (0): Conv2d(88, 88, kernel_size=(3, 3), stride=(1, 1), padding=(1, 1), groups=88, bias=False)
          (1): BatchNorm2d(88, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): ReLU(inplace=True)
        )
        (2): ConvNormActivation(
          (0): Conv2d(88, 24, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(24, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (4): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(24, 96, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(96, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (1): ConvNormActivation(
          (0): Conv2d(96, 96, kernel_size=(5, 5), stride=(2, 2), padding=(2, 2), groups=96, bias=False)
          (1): BatchNorm2d(96, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (2): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(96, 24, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(24, 96, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (3): ConvNormActivation(
          (0): Conv2d(96, 40, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(40, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (5): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(40, 240, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(240, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (1): ConvNormActivation(
          (0): Conv2d(240, 240, kernel_size=(5, 5), stride=(1, 1), padding=(2, 2), groups=240, bias=False)
          (1): BatchNorm2d(240, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (2): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(240, 64, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(64, 240, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (3): ConvNormActivation(
          (0): Conv2d(240, 40, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(40, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (6): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(40, 240, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(240, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (1): ConvNormActivation(
          (0): Conv2d(240, 240, kernel_size=(5, 5), stride=(1, 1), padding=(2, 2), groups=240, bias=False)
          (1): BatchNorm2d(240, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (2): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(240, 64, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(64, 240, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (3): ConvNormActivation(
          (0): Conv2d(240, 40, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(40, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (7): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(40, 120, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(120, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (1): ConvNormActivation(
          (0): Conv2d(120, 120, kernel_size=(5, 5), stride=(1, 1), padding=(2, 2), groups=120, bias=False)
          (1): BatchNorm2d(120, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (2): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(120, 32, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(32, 120, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (3): ConvNormActivation(
          (0): Conv2d(120, 48, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(48, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (8): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(48, 144, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(144, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (1): ConvNormActivation(
          (0): Conv2d(144, 144, kernel_size=(5, 5), stride=(1, 1), padding=(2, 2), groups=144, bias=False)
          (1): BatchNorm2d(144, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (2): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(144, 40, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(40, 144, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (3): ConvNormActivation(
          (0): Conv2d(144, 48, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(48, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (9): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(48, 288, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(288, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (1): ConvNormActivation(
          (0): Conv2d(288, 288, kernel_size=(5, 5), stride=(2, 2), padding=(2, 2), groups=288, bias=False)
          (1): BatchNorm2d(288, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (2): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(288, 72, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(72, 288, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (3): ConvNormActivation(
          (0): Conv2d(288, 96, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(96, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (10): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(96, 576, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(576, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (1): ConvNormActivation(
          (0): Conv2d(576, 576, kernel_size=(5, 5), stride=(1, 1), padding=(2, 2), groups=576, bias=False)
          (1): BatchNorm2d(576, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (2): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(576, 144, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(144, 576, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (3): ConvNormActivation(
          (0): Conv2d(576, 96, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(96, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (11): InvertedResidual(
      (block): Sequential(
        (0): ConvNormActivation(
          (0): Conv2d(96, 576, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(576, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (1): ConvNormActivation(
          (0): Conv2d(576, 576, kernel_size=(5, 5), stride=(1, 1), padding=(2, 2), groups=576, bias=False)
          (1): BatchNorm2d(576, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
          (2): Hardswish()
        )
        (2): SqueezeExcitation(
          (avgpool): AdaptiveAvgPool2d(output_size=1)
          (fc1): Conv2d(576, 144, kernel_size=(1, 1), stride=(1, 1))
          (fc2): Conv2d(144, 576, kernel_size=(1, 1), stride=(1, 1))
          (activation): ReLU()
          (scale_activation): Hardsigmoid()
        )
        (3): ConvNormActivation(
          (0): Conv2d(576, 96, kernel_size=(1, 1), stride=(1, 1), bias=False)
          (1): BatchNorm2d(96, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
        )
      )
    )
    (12): ConvNormActivation(
      (0): Conv2d(96, 576, kernel_size=(1, 1), stride=(1, 1), bias=False)
      (1): BatchNorm2d(576, eps=0.001, momentum=0.01, affine=True, track_running_stats=True)
      (2): Hardswish()
    )
  )
  (avgpool): AdaptiveAvgPool2d(output_size=1)
  (classifier): Sequential(
    (0): Linear(in_features=576, out_features=1024, bias=True)
    (1): Hardswish()
    (2): Dropout(p=0.2, inplace=True)
    (3): Linear(in_features=1024, out_features=1000, bias=True)
  )
)
```

# 测试模型
pytorch官方训练的模型对输入有严格限制，必须是Tensor类型的数组，包含多幅图像。为了测试一张图像，我们可以使用PIL读取图片之后
使用官方的transform工具转为可输入的Tensor类型数据后，插入到空数组中得到特征值

```python
from PIL import Image

# 读取图像
img = Image.open('image/path').convert('RGB')

# 定义transform
import torch
from torchvision import transforms

transform = transforms.Compose([
  transforms.Resize(128),
  transforms.ToTensor(),
  transforms.Normalize(
    mean=[0.485, 0.456, 0.405],
    std=[0.229, 0.224, 0.225]
  )
])

# 获取Tensor
img_tensor = transform(img)

# 转为batch
batch_t = torch.unsqueeze(img_tensor, 0)

# 测试
features = mobilenet(batch_t)
print(features[0].shape) # torch.Size([1000])
```

这样，就有了大小为1000的特征向量，可以作为环境的状态进行DRL训练了
