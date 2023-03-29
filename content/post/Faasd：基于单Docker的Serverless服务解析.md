---
title: "Faasd：基于单Docker的Serverless服务解析-序"
date: 2023-03-09T17:45:54+08:00
draft: true
toc: true
images:
tags: 
  - Docker
  - Serverless
  - OpenFaas
---
> [github链接](https://github.com/openfaas/faasd)

![faasd!](https://github.com/openfaas/faasd/raw/master/docs/media/social.png)

Serverless 是云计算的一种模型。以平台即服务（PaaS）为基础，无服务器运算提供一个微型的架构，终端客户不需要部署、配置或管理服务器服务，代码运行所需要的服务器服务皆由云端平台来提供。根据CNCF的定义，Serverless是指构建和运行不需要服务器管理的应用程序的概念

目前个人所知的开源的Serverless平台包括
- Knative，谷歌官方所推崇的Serverless平台范式，旨在统一Serverless的接口定义相关的规范，和k8s高度绑定
- OpenFunction，国内青云开源产品，目前文档比较简陋，功能较为简单
- OpenFaas，一种可以运行在单机Docker和K8s集群的Serverless平台，也是本文Faasd的父产品

为了更清晰的了解一个Serverless平台的工作流程，我选择OpenFaas的子项目Faasd作为切入点，是因为Faasd是一个简单的可随时部署的单机Serverless平台，但兼容集群操作的API

本文关注于OpenFaas是如何从集群内布创建、调用Function容器的，针对如何通过Gateway入口外部调用Function则将会在后续进行解析

- [工作流]()
- [创建函数]()
- [调用函数]()