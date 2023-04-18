---
title: "牛客Golang面经(一): 默安科技"
date: 2023-04-18T11:50:27+08:00
draft: false
toc: false
images:
tags: 
  - 面试
  - 面经
---
> [牛客地址](https://www.nowcoder.com/discuss/476117375096578048?sourceSSR=home)

#### 1. string的底层结构
在string中包含两个字段，分别是保存数据的byte数组，表示字符串长度的len

#### 2. string底层是一个指针，函数传参时修改string会影响吗？
首先在golang中，string类型的数据规定是不可修改的，另外在go中的所有函数参数中只存在值传递

#### 3. 如果传参是一个slice会影响原数组吗
分情况讨论，如果通过下标形式对原数组进行修改，则会影响到外层slice中，如果对slice进行添加元素或删除元素等操作，则不会影响到外层的slice

具体来说，slice的数据结构中有三个字段，包括指向数组数据的指针、切片长度以及容量。当slice作为参数传递时，由于go不存在引用传递，会对slice进行复制后传入到函数中。该slice的地址和最外层的slice地址并不相同，但指向的数据地址是相同的（因为是复制的），因此，如果在函数中对该slice进行添加或删除操作时，该slice的长度变化以及可能触发的扩容操作导致的容量变化并不能影响到函数外的slice

#### 4. map有了解吗，怎么保证并发安全
map是一个引用类型的Hash结构，可以进行高效的查询和插入操作。保证并发安全可以添加一个互斥锁对map进行保护，或使用sync包的Map

#### 5. sync.map的底层是什么有了解吗
在sync.Map中也包含一个互斥锁，但包含着两个存放数据的结构read和dirty以及一个计数器用于计算diry的数量。Read是一个原子类型数据，其中包含着一个map结构和标记dirty是否存在当前不存在的数据的字段

![sync.Map关系图](https://segmentfault.com/img/remote/1460000020946992)

__参考资料__
- [sync.Map的拓扑关系图](https://segmentfault.com/a/1190000020946989)
- [通过实例深入理解sync.Map的工作原理](https://tonybai.com/2020/11/10/understand-sync-map-inside-through-examples/)

#### 6. context的用处
context是golang中的上下文结构，可以用于函数或协程之间的数据传递，也可以用于控制协程的退出

#### 7. 对一个关闭的channel进行读操作会发生什么，写呢？
如果对一个已经关闭的channel进行读操作并不会报错，channel的读取操作会默认返回定义channel时数据结构的空值。如果对一个已经关闭的channel进行写操作则会直接panic

#### 8. GRPC用什么协议
protobuf

#### 9. protobuf用什么序列化
protobuf使用`varint`和`zigzag`编码器进行序列化和反序列化。前者是一个可变长编码，使用1个或多个字节对证书进行编码，后者是一个整数压缩编码算法。

__参考资料__
- [Google Protobuf编码原理](https://sunyunqiang.com/blog/protobuf_encode/)
- [整数压缩编码ZigZag](https://www.cnblogs.com/en-heng/p/5570609.html)

#### 10. mysql索引的优化
- 联合索引注意最左前缀索引
- 建立索引的字段尽可能小： 由于mysql索引使用B+树，索引的查询收到树高影响，每个叶子所占用的空间是固定的，所占空间越小，数据项数据越多，高度越低，查询更快速

__参考资料__
- [MySQL索引原理及慢查询优化](https://tech.meituan.com/2014/06/30/mysql-index.html)

#### 11. 你一般设置的varchar一般都是变长的吗，utf8 varchar最长是多少
varchar我一般设置定长，长文本我使用text类型存放。utf8的行长度为65532字节，字符占3字节，因此字符长度不超过21845，超过会强行转换为text并warning

实际上，每一行的数据从第2个字节开始，并且每一个varchar会占用2个字节用于长度，因此实际上varchar最大长度为21844

#### 12. MySQL中的原子性
MySQL的原子性实现依靠undolog，原子性保证每个事物被视为一个单独的单元执行只有成功和失败两种状态。在失败的时候会根据undolog进行回滚，undolog中记录的是sql执行的相关信息，回滚时，innodb会根据undolog的内容进行相反的操作

__参考资料__
- [MySQL教程(十)---MySQL ACID 实现原理](https://www.lixueduan.com/posts/mysql/10-acid/)

#### 13. 有了解过Linux的虚拟内存吗
Linux的虚拟内存是为了解决多进程直接使用物理内存造成的bug，每个进行在执行时都会分配独立的虚拟内存用于隔离其他进程的数据。每个PID维护一张虚拟内存和物理内存之间相互映射的页表，程序在虚拟内存中的数据是线性存放的而实际上在物理内存中可能是离散的。

__参考资料__
- [深入理解 Linux 虚拟内存管理](https://xiaolincoding.com/os/3_memory/linux_mem.html#_5-%E8%BF%9B%E7%A8%8B%E8%99%9A%E6%8B%9F%E5%86%85%E5%AD%98%E7%A9%BA%E9%97%B4%E7%9A%84%E7%AE%A1%E7%90%86)

#### 14. 我现在需要查看一个文件最新更新的100行怎么做Linux你一般用什么指令

```shell
watch tail -n 100 xxx.file
```

#### 15. 进程和进程之间的通信是怎么实现的
Linux的进程通信有多种形式
- 信号量机制：如kill
- 管道：通过管道可以将上一个进程的输出直接输入到下一个进程的输入中如常见的查询`cat xxx | grep xxx`
- 消息队列
- 共享内存
- 套接字

#### 16. 两次握手之后一直发送数据和三次握手后再发送数据结果一样吗
如果持续发送数据结果是一样的，三次握手是为了确认客户端已经收到服务端允许连接的通知，在第三次握手时同样可以携带数据发送到服务端

#### 17. TCP断开时的TIME_OUT有什么讲究吗
没找到合适的资料，正在找

#### 18. 对一个排好序的数组进行快排时间复杂度是多少
如果是正序则是O(n)，如果是逆序则是O(n^2)

#### TopK
堆，选前K个，分治法