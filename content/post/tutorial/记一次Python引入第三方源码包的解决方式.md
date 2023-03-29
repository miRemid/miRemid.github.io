---
title: "记一次Python引入第三方源码包的解决方式"
date: 2022-08-16T17:13:03+08:00
draft: false
toc: false
images:
tags: 
  - 教程
  - Python
---

> Python虽然挺好用的，但也挺难用的

众所周知，Python由于其“先进”的包管理功能，让你几乎可以通过一条命令安装所有打包好的第三方包。但是当我们想从第三方的项目里面引入少许函数时，这种方式就行不得通了.

例如我想从一个Pytorch项目中将模型定义的函数引入（说的就是你Yolo）到我自己的函数中，项目目录结构如下：
```shell
.
├── main.py
├── src
│   └── module
│       ├── a.py
│       ├── __init__.py
│       └── __pycache__
│           ├── a.cpython-38.pyc
│           └── __init__.cpython-38.pyc
└── thirdparty
    └── codes
        ├── modules
        │   ├── b.py
        │   ├── __init__.py
        │   ├── __pycache__
        │   │   ├── b.cpython-38.pyc
        │   │   ├── c.cpython-38.pyc
        │   │   └── __init__.cpython-38.pyc
        │   └── submodule
        │       ├── c.py
        │       └── __pycache__
        │           └── c.cpython-38.pyc
        ├── __pycache__
        │   └── b.cpython-38.pyc
        └── run
            └── run.py
```
其中`a.py`和`b.py`的代码如下
```python
# src/module/a.py
from thirdparty.codes.modules import b

def print_a():
    print("this is a")
    b.print_b()

if __name__ == '__main__':
    print_a()

# thirdparty/codes/modules/b.py
from codes.modules.submodule import c

def print_b():
    print("this is b")
    c.print_c()

# thirdparty/codes/modules/submodule/c.py
def print_c():
    print("this is c")
```
此时，如果我们分别在主目录和`module`目录测试`a.py`文件时会出现如下错误
```shell
Traceback (most recent call last):
  File "main.py", line 1, in <module>
    from src.module import a
  File "/home/hakureisk/Workspace/tmp/src/module/a.py", line 1, in <module>
    from thirdparty.codes.modules import b
  File "/home/hakureisk/Workspace/tmp/thirdparty/codes/modules/b.py", line 1, in <module>
    from codes.modules.submodule import c
ModuleNotFoundError: No module named 'codes'

Traceback (most recent call last):
  File "a.py", line 1, in <module>
    from thirdparty.codes.modules import b
ModuleNotFoundError: No module named 'thirdparty'
```
这种错误在遇到使用`git submodule`的方式引入第三方项目时非常常见，解决方式呢也有很多种

对于`main.py`这种处于最外层的入口文件，解决方式非常简单，只需要将`thirdparty`目录加入到当前执行的系统路径中即可
```python
import sys
sys.path.append('thirdparty')
sys.path.append('thirdparty/codes')

from src.module import a

if __name__ == '__main__':
    a.print_a()

# output
this is a
this is b
this is c
```
而对于`a.py`这种单独测试文件呢就比较麻烦了，你需要将第三方包的绝对路径添加到系统路径之中
```python
# src/module/a.py
if __name__ == '__main__':
    import sys
    sys.path.append('/path/of/the/project/tmp')
    sys.path.append('/path/of/the/project/tmp/thirdparty')

from thirdparty.codes.modules import b

def print_a():
    print("this is a")
    b.print_b()

if __name__ == '__main__':
    print_a()
```

需要注意的是，这里的demo非常简单，并不能覆盖到全部情况，例如将`b.py`文件中引入部分改为`from modules.submodule import c`则上述方法就失效了（`main.py`除外）

为了解决这种情况，最简单的方式就是将第三方包以`pth`文件的形式导入到系统路径之中，使其变为`pip`安装包的引入方式，具体操作如下
1. 找到目前环境下`site-packages`文件目录

我们需要将`pth`文件保存在和第三方包一样的路径之中，通过下面的脚本即可获取当前环境的第三方包路径
    
    python -m site

2. 创建`pth`文件

在对应的目录中，创建对应的`pth`文件，文件名自定义取一个有辨识度但和已经安装的包不冲突的名字，例如我这使用的`thirdparty.pth`

    cd /path/to/site-packages
    touch thirdparty.pth

创建完成后，将需要引入的第三方包路径添加至文件之中

    # /path/to/site-packages/thirdparty.pth
    /path/to/third/party/module

    # example
    /home/xxx/tmp/
    /home/xxx/tmp/thirdparty

注意，需要将对应的包明确写入文件之中，这是为了方便寻找其子包的存在

最后运行脚本就可以成功引入了: )

