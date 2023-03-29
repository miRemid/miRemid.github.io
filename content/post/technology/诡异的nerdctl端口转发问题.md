---
title: "诡异的nerdctl端口转发问题"
date: 2023-03-16T10:40:57+08:00
draft: false
toc: false
images:
  - https://images.unsplash.com/photo-1615796153287-98eacf0abb13?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=774&q=80
tags: 
  - 技术
  - 容器
summary: 没有什么是重启解决不了的，如果有那就重启两次吧
---


由于某些原因，需要从Docker转到containerd上面来，为了保持和Docker相同的体验选择Nerdctl作为cli工具

但是在安装之后运行简单的nginx测试时出现了网络问题

```shell
nerdctl pull nginx:alpine

nerdctl run -d --name nginx -p 80:80 nginx:alpine

➜  ~ nerdctl ps
CONTAINER ID    IMAGE                             COMMAND                   CREATED          STATUS    PORTS                 NAMES
25ec8773aa63    docker.io/library/nginx:alpine    "/docker-entrypoint.…"    5 minutes ago    Up        0.0.0.0:80->80/tcp    nginx
```
从打印的日志上来看，当前Nerdctl已经将nginx容器的80端口转发到宿主机的80端口，在宿主机的shell中使用curl访问80端口也能正常访问
```shell
➜  ~ curl http://localhost:80
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```
但是从宿主机外部访问时，却怎么也访问不到
```shell
curl http://10.0.8.9:80
```
但是我换Docker则一切正常
```shell
docker pull nginx:alpine

docker run -d --name nginx -p 80:80 nginx:alpine
curl http://10.0.8.9:80
```

检查iptables可以发现，nerdctl已经帮我们把nerdctl创建的网桥转发出去了
```shell
# 这是还没创建容器时的iptables
➜  ~ iptables --list-rules
-P INPUT ACCEPT
-P FORWARD DROP
-P OUTPUT ACCEPT
-N CNI-ADMIN
-N CNI-FORWARD
-N CNI-ISOLATION-STAGE-1
-N CNI-ISOLATION-STAGE-2

-A FORWARD -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j CNI-ISOLATION-STAGE-1
-A FORWARD -m comment --comment "CNI firewall plugin rules" -j CNI-FORWARD

-A CNI-FORWARD -m comment --comment "CNI firewall plugin admin overrides" -j CNI-ADMIN
-A CNI-ISOLATION-STAGE-1 -i nerdctl0 ! -o nerdctl0 -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j CNI-ISOLATION-STAGE-2
-A CNI-ISOLATION-STAGE-1 -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j RETURN
-A CNI-ISOLATION-STAGE-2 -o nerdctl0 -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j DROP
-A CNI-ISOLATION-STAGE-2 -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j RETURN

# 当创建nginx容器之后的iptables
-P INPUT ACCEPT
-P FORWARD DROP
-P OUTPUT ACCEPT
-N CNI-ADMIN
-N CNI-FORWARD
-N CNI-ISOLATION-STAGE-1
-N CNI-ISOLATION-STAGE-2

-A FORWARD -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j CNI-ISOLATION-STAGE-1
-A FORWARD -m comment --comment "CNI firewall plugin rules" -j CNI-FORWARD
-A FORWARD -j DOCKER-USER
-A FORWARD -j DOCKER-ISOLATION-STAGE-1

-A CNI-FORWARD -m comment --comment "CNI firewall plugin admin overrides" -j CNI-ADMIN
-A CNI-FORWARD -d 10.4.0.8/32 -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
-A CNI-FORWARD -s 10.4.0.8/32 -j ACCEPT
-A CNI-ISOLATION-STAGE-1 -i nerdctl0 ! -o nerdctl0 -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j CNI-ISOLATION-STAGE-2
-A CNI-ISOLATION-STAGE-1 -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j RETURN
-A CNI-ISOLATION-STAGE-2 -o nerdctl0 -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j DROP
-A CNI-ISOLATION-STAGE-2 -m comment --comment "CNI firewall plugin rules (ingressPolicy: same-bridge)" -j RETURN
```
可以发现，nerdctl帮助我们为nginx容器创建了两条iptables规则，检查nginx容器的ip地址可以发现与上述规则是对应的
```shell
➜  ~ nerdctl inspect nginx
[
    {
        "Id": "029f32038d2ab668cae206620b83f143fa799638d6fd623808aa55175825b251",
        "Created": "2023-03-16T02:58:34.799160852Z",
        "Path": "/docker-entrypoint.sh",
        "Args": [
            "nginx",
            "-g",
            "daemon off;"
        ],
        "State": {
            "Status": "running",
            "Running": true,
            "Paused": false,
            "Restarting": false,
            "Pid": 1552,
            "ExitCode": 0,
            "FinishedAt": "0001-01-01T00:00:00Z"
        },
        "Image": "docker.io/library/nginx:alpine",
        "ResolvConfPath": "/var/lib/nerdctl/1935db59/containers/default/029f32038d2ab668cae206620b83f143fa799638d6fd623808aa55175825b251/resolv.conf",
        "HostnamePath": "/var/lib/nerdctl/1935db59/containers/default/029f32038d2ab668cae206620b83f143fa799638d6fd623808aa55175825b251/hostname",
        "LogPath": "/var/lib/nerdctl/1935db59/containers/default/029f32038d2ab668cae206620b83f143fa799638d6fd623808aa55175825b251/029f32038d2ab668cae206620b83f143fa799638d6fd623808aa55175825b251-json.log",
        "Name": "nginx",
        "RestartCount": 0,
        "Driver": "overlayfs",
        "Platform": "linux",
        "AppArmorProfile": "",
        "Mounts": null,
        "Config": {
          ...
        },
        "NetworkSettings": {
            "Ports": {
                "80/tcp": [
                    {
                        "HostIp": "0.0.0.0",
                        "HostPort": "80"
                    }
                ]
            },
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "IPAddress": "10.4.0.8",
            "IPPrefixLen": 24,
            "MacAddress": "02:52:2f:d9:d2:52",
            "Networks": {
                "unknown-eth0": {
                    "IPAddress": "10.4.0.8", # IP地址
                    "IPPrefixLen": 24,
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "MacAddress": "02:52:2f:d9:d2:52"
                }
            }
        }
    }
]
```
按道理按照规则`-A CNI-FORWARD -s 10.4.0.8/32 -j ACCEPT`已经可以从外部访问80端口的服务了，但是怎么都行不通，所以干脆手动添加一条iptable规则吧，直接指定转发80端口出去
```shell
➜  ~ iptables -A CNI-FORWARD -p tcp -m tcp --dport 80 -j ACCEPT
```
这时，访问对应地址的80端口就可以访问服务了，但这明显只是一个临时解决方案，那如何彻底解决呢？直到我找到这个[issue](https://github.com/containerd/nerdctl/issues/1048#issuecomment-1325991013)...

没有什么是重启解决不了的，如果有那就重启两次！重启之后，删除刚才临时添加的iptable之后测试直接启动nginx服务，发现就可以从外部访问了🤣