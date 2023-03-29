---
title: "Websocket关闭错误代码含义"
date: 2021-08-07T15:00:19+08:00
draft: false
toc: false
images:
tags: 
  - 技术
  - Web
---

## 前情
在处理B站直播的Websocket源时，经常发生连接关闭的错误，最常见的就是`close 1006 (abnormal closure): unexpected EOF`错误，你说http的状态码还知道，这websocet的状态码还真不知道，于是去查了查记录一下

## RFC 6455
根据[RFC 6455](https://tools.ietf.org/pdf/rfc6455.pdf)定义的内容，Websocket在处理关闭时设置了一系列的代码提示。当对一个已经建立的连接进行关闭时，在终端处可能提供一个关闭的解释说明，客户端可以根据这个代码来推测终端关闭连接的原因从而更新客户端连接的代码。当然，终端在关闭时也可以忽视代码直接进行关闭

- 1000 Normal Closure

1000表明这是一个正常的关闭，表明要传输的数据已经全部完成可以退出

- 1001 Going Away

1001说明终端可能已经找不到该连接，例如服务可能宕机或浏览器重定向至其他页面

- 1002 Protocal error

1002说明连接被终端由于消息协议错误而进行强制性退出

- 1003 Unsupported Data

1003说明终端接受到一个无法处理的数据而进行强制性退出(例如，服务器可能只能够处理文本数据但接受到了二进制数据)

- 1004 ---Reserved---

1004字段保留，未来可能用得到

- 1005 No Status Rcvd

1005是一个保留数据，绝对不允许终端将其设置为关闭时的状态码。It is designated for use in applications expecting a status code to indicate the no stats code was actually present.

- 1006 Abnormal Closure

1006是一个保留数据，绝对不允许终端将其设置为关闭时的状态码。它是用来指定需要状态码标志连接异常关闭的程序，例如没有发送或接受控制数据

- 1007 Invalid frame payload data

1007说明终端收到了一个不符合规定的格式数据而关闭连接(例如，non-UTF-8格式的数据包括了一段text消息)

- 1008 Policy Violation

1008说明终端收到了一个违反服务规则的消息而关闭连接，这是一个多功能状态码，它可以用于没有其他合适的状态码时发送出去例如(1003或1009)或者在有必要隐藏自己的规则信息时发送出去

- 1009 Message Too Big

1009说明终端收到了一个过大的数据而关闭连接

- 1010 Mandatory Ext

- 1011 Internal Server Error

1011说明服务端在处理请求时遇到了一个意外情况

- 1015 TLS handshake

1005是一个保留数据，绝对不允许终端将其设置为关闭时的状态码。它用于处理在TLS握手时发生的错误。

## 错误处理
了解1006错误之后猜测是客户端在退出时没有发送退出信息导致异常关闭，图省事直接处理err忽略完事

```go
_, body, err := cli.conn.ReadMessage()
if err != nil {
  if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
    DPrintf("Unexpected Close Error")
    closeConn <- struct{}{}
    continue
  }
}
```