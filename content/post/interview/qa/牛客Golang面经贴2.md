---
title: "牛客Golang面经(二)：居学科技"
date: 2023-04-18T11:50:27+08:00
draft: false
toc: false
images:
tags: 
  - 面试
  - 面经
---
> [牛客地址](https://www.nowcoder.com/feed/main/detail/451716c717ce47a9989e4ab043b282a3?sourceSSR=post)


## 1. GRPC有几种通信方式
- 简单RPC调用，类似HTTP
- 服务端流式RPC调用，类似监听一个消息队列信息
- 客户端流式RPC调用，类似持续往消息队列发送信息
- 双向流式RPC，类似Websocket

## 2. GRPC使用什么进行端口暴露
原贴描述不够清晰，个人猜测是通过什么方式将grpc的服务暴露给外部HTTP？

如果是上述这个问题，则可以使用grpc-gateway插件在编译protoc文件时自动生成对应RPC函数的HTTP接口

或者我们可以自行写一个网关程序，专门用于内部GRPC接口的请求。另外GRPC使用的还是HTTP，不过是2.0协议，在数据头中协议的版本号与普通HTTP接口所使用的的1.0协议不一致，可以通过请求Header中的协议号进行区分从而调用GRPC函数

## 3. MySQL索引优化
- 最左前缀匹配
- 尽量使用in，避免使用<,>,bettwen
- 建立索引的字段大小尽可能小

## 4. Channel使用场景
- 协程之间的通信
- 生产-消费者模式

## 5. Channel并发安全吗
channel是并发安全的，在channel内部数据结构中包含了一个互斥锁用于控制并发

## 6. 直接使用互斥锁控制并发不行吗
直接使用互斥锁控制并发也是可以的，但是人工维护锁的成本比直接使用channel来说要高

## 7. 除了Channel还可以用什么进行通信
context

## 8. context有几种类型
- Background，空
- TODO, 空
- WithTimeout, 携带超时定时器
- WithCancel, 携带退出触发函数
- WithValues, 携带数据
- WithDeadline, 携带超时定时器以及父Context

## 9. Redis数据类型
List, Set, Hash, String

## 10. Linux如何查询CPU占用率
```shell
cat /proc/stat | grep cpu | awk -F ' ' '{usage = ($2+$4)*100/($2+$4+$5)} END {print usage "%"}'
```

## 11. TCP为什么使用三次握手，两次行不行
两次握手后，服务端会要求收到确认请求再接受数据，如果不发送第三次握手确认报文，服务端会持续监听确认报文从而浪费服务端资源