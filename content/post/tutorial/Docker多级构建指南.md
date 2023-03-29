---
title: "Docker多级构建指南"
date: 2021-08-25T09:01:12+08:00
draft: false
toc: true
images:
tags: 
  - 教程
  - 容器
---

> 多级构建让你的容器更小巧~

在构建Docker镜像的时候经常发现构建出来的容器大小非常的大，而我本地编译出来的二进制文件也不过26MB左右而Docker容器居然有120MB！
问题出现在哪里了呢？我们一步一步从最起始的地方开始看起

## 构建方式
Docker容器的构建方法我个人常用以下几种
1. 本地打包二进制后直接放入Docker中运行
2. 容器内打包直接运行

第一种的方式就是在本地编译好要执行的文件之后放入容器之中运行，这种方式无疑是最简单的方法也是最容易理解的方法。
但这种方法仔细一想就会发现有很多问题，假设我们需要在Windows机器编译Linux版本的容器内运行要怎么办呢？
除此之外，假设我们换了一台机器也需要编译一份Docker镜像时，本地没有对应的编译工具要如何解决呢？
为了应对上述问题，也就有了第二种方法。

第二种的方式就是将代码放在容器内部编译，这样既解决了编译环境问题也解决了编译平台的问题。
但第二种方法带来的问题就是我今天要说的，那就是打包后的镜像体积过大，完全不如第一种打包出来的镜像。
这是因为，第一种只是将二进制文件放入容器内直接运行，而第二种还添加了项目的源代码文件，除此之外，还有一堆用于
编译二进制文件的工具在镜像之中，因此打包出来的镜像体积会非常巨大。
那有没有一种方法既能随意构建又不产生大体积的镜像呢？那就要谈谈Docker的多级构建了。

## 多级构建
Docker的多级构建也就是常说的`multi-stage build`，你可以通过指定多个stage分别完成不同的任务最后合在一起完成最终的构建。

例如有一个Go的项目，我们可以在stage1阶段进行编译操作，而在stage2阶段进行运行操作，这样就相当于构建1和2相结合，完成最后的构建。

多级构建其实非常简单，其精髓就在于一个命令那就是`COPY`，这个命令不仅仅可以将本地文件拷贝至Docker的build进程上下文中，还可以在多个stage中
进行文件的复制，而`ADD`命令则只能用于前者，因此我通常在第一级构建时使用`ADD`命令将所有的源代码复制到Docker上下文中后，使用`COPY`应对其他层级的构建。`COPY`的参数非常简单

    COPY --from=stage source dest

其中`from`就是用于表明从哪个stage中复制文件，不添加该参数则默认从宿主机中复制文件。stage可以用数字表示从第几级复制(从0开始)，也可以用字符串来指定层级复制但需要对层级进行命名

    # builder stage
    FROM x as builder

    # final stage
    FROM x as final
    COPY --from=builder /abc /abc
    COPY --from=0 /abc/abc

## 实战
以我的[Yuki](https://github.com/miRemid/yuki)为例，该项目由`React`和`Go`组成，其中`Go`将会提供服务器提供前端接口。

首先分析阶层，我们需要编译两个项目，其中是`React`另一个是`Go`，并且`Go`的编译需要依赖于`React`。因此很容易得出，以下顺序
1. 编译React
2. 编译Go
3. 运行

为了让编译出来的容器尽可能的小，我们在选取构建容器时也尽量选择小的容器来进行编译例如我最喜欢的`alpine`。因此在[Docker Hub](https://hub.docker.com)中寻找关于`nodejs`和`golang`的`alpine`版本，由于原版`golang`的`alpine`版本不附带`gcc`因此我选择了第三方的容器`tetafro/golang-gcc`。而提供运行环境的容器我也选择了`alpine:3.14`

- NodeJS: node:14.17.5-alpine(117MB)
- Golang: tetafro/golang-gcc:1.16-alpine(425MB)
- Runtime: alpine:3.14(5.6MB)

在选择完构建容器后就可以正式构建自己的程序了，首先是React的打包，注意的是如果将`node_modules`文件一并Add的话将不会再拉取一遍不符合我们的要求，因此需要提请编写一个类似于`.gitignore`文件的`.dockerignore`文件，例如我使用的
```shell
cat ./dockerignore

release
data
web/node_modules
web/dist
.git
.cache
```

这样就可以安心写`Dockerfile`了

```Dockerfile
# build react
# 使用下载的node容器，并将该阶段命名为node
FROM node:14.17.5-alpine as node
# 将所有源代码放入容器的/yuki文件夹中
ADD . /yuki
# 切换当前工作路径到/yuki/web中，该目录是React工程根目录
WORKDIR /yuki/web
# 拉取依赖并打包，/yuki/web/dist
RUN yarn && yarn build
```

在构建完前端后就可以来编译后端项目了，这时就要使用多级构建中的`COPY`命令了
```Dockerfile
# build golang
# 使用下载的go容器，并将该阶段命名为golang
FROM tetafro/golang-gcc:1.16-alpine as golang
# 从node层中复制源代码到本层的/yuki中，此时文件夹中已经包含了node层打包好的前端项目
COPY --from=node /yuki /yuki
# 切换工作路径
WORKDIR /yuki
# 静态编译
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -ldflags \
	' -extldflags "-static"' \
	-o yuki_linux_amd64
```

构建完二进制文件后，就需要将二进制文件和一些必要的静态文件放入到最终的运行容器中运行了
```Dockerfile
FROM alpine:3.14

# 从golang层复制所需要的文件
COPY --from=golang /yuki/yuki_linux_amd64 /yuki_linux_amd64
COPY --from=golang /yuki/docs /docs
# 切换工作路径
WORKDIR /
# 运行
CMD [ "sh", "-c", "/yuki_linux_amd64" ]
```
可以看到，我们最后只复制了必要的二进制文件和所需要的静态文件到最后的容器之中，并没有其他任何东西，因此最终的容器大小理论上是默认大小加上二进制文件大小，事实是这样吗。事实上也确实如此，查看本地静态编译文件大小和打包后的镜像大小发现两者相差几乎等于镜像的大小
```shell
zsh> cd release && ll
-rwxrwxr-x 1 kamir kamir  26M Aug 25 09:39 yuki_linux_amd64
zsh> docker images | grep yuki
yuki                  latest                  eb2e260d4701   5 hours ago    32.8MB
zsh> docker images | grep alpine
alpine                3.14                    021b3423115f   2 weeks ago    5.6MB
```
这样一个完美的小容器就诞生了~