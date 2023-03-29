---
title: "Hugo+Github Action自动部署静态页面"
date: 2021-05-22T17:37:27+08:00
draft: false
toc: true
images:
tags: 
  - 教程
  - Github
---

## [Hugo](https://gohugo.io/)
> Hugo is one of the most popular open-source static site generators. With its amazing speed and flexibility, Hugo makes building websites fun again.

### 安装

如果你是Windows用户，可以直接从[Github Release](https://github.com/gohugoio/hugo/releases)页面下载对应的版本即可

如果你是Linux用户则有三种方式来安装Hugo
- 从[Github Release](https://github.com/gohugoio/hugo/releases)下载
- 从包管理器下载
- 自行编译

在这介绍下自行编译，最新的Hugo使用到了`go1.16`的`io/fs`包，因此如果想自行编译最新版本的Hugo则要安装或更新系统`go`版本至`1.16+`。安装完语言编译环境后按以下步骤进行

    mkdir $HOME/src
    cd $HOME/src
    git clone https://github.com/gohugoio/hugo.git
    cd hugo
    go install

如果你想让Hugo支持`sass`的话则需要编译另一个版本

    CGO_ENABLED=1 go install --tags extended

在终端中输入`hugo version`来检查是否编译成功

### 创建站点

安装完Hugo后就可以来创建一个网站站点，假设名称为`demo`

    hugo new site demo
    cd demo
    git init

在初始化站点后，需要为站点选择一个主题，这里选择我使用的[Hermit](https://github.com/Track3/hermit)

    git submodule add https://github.com/Track3/hermit.git themes/hermit
    cp themes/hermit/exampleSite/config.toml .

随后创建一篇新文章

    hugo new posts/Hello-hugo.md

添加完文章后输入以下指令进行预览

    hugo server -D

![预览图](https://i.loli.net/2021/05/22/ULygZGTj8bPilw6.png)

## Github Action自动化发布

当完成一个站点之后，可以发布到github page中，在这介绍如何使用Github Action来自动化部署站点到github page中

Github Page分为用户Page和项目Page，用户Page通常为`username.github.io`而项目Page通常为`username.github.io/project`，但也可以通过添加`CNAME`文件来自定义域名，在这展示项目Page

首先创建一个项目(`Repository`)，假设为`demo`记住要设置为公开，然后在你的站点的根目录创建以下文件并写入

```yml
# .github/workflows/gh-pages.yml
name: github pages

on:
  push:
    branches:
      - master  # 要编译的文件分支
  pull_request:

jobs:
  deploy:
    runs-on: ubuntu-20.04
    steps:
      - uses: actions/checkout@v2
        with:
          submodules: true  # Fetch Hugo themes (true OR recursive)
          fetch-depth: 0    # Fetch all history for .GitInfo and .Lastmod

      - name: Setup Hugo
        uses: peaceiris/actions-hugo@v2
        with:
          hugo-version: 'latest'
          extended: true

      - name: Build
        run: hugo --minify

      - name: Deploy
        uses: peaceiris/actions-gh-pages@v3
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./public
```

在部署到Github之前，还需要更改Hugo的根配置，将`config.toml`中的`baseURL`更改为`https://{username}.github.io/demo`，如果配置了域名则为`https://{CUSTOM_DOMAIN}/demo`

随后，将你的站点所有文件push到Github中

    git remote add origin https://github.com/{your_user_name}/{your_project}.git
    git add .
    git commit -m "Upload"
    git push -u origin master

这时，前往你的项目页面查看`Actions`标签查看编译结果，如果成功则如下所示

![Action](https://i.loli.net/2021/05/22/bwJztQ7qr9k2PnM.png)

完成编译之后，还不能立刻通过网址来访问，需要在项目的`Settings\Pages`中，选择Github Page的来源，在这选择`gh-pages`

设置成功后访问`https://{username}.github.com/{project_name}`即可，如果配置了`CNAME`则通过`https://{CUSTOM_DOMAIN}/{project_name}`即可

![结果](https://i.loli.net/2021/05/22/Z6gRYaLcfUbpByx.png)