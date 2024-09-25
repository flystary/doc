Docker仓库部署（支持多CPU架构）

1. ### 基本环境配置

   1. dockerhub.7x-networks.net是解析到公网的。仓库构建和push时直接访问127.0.0.1就可以，不需要访问公网，在/etc/hosts里添加

```Plain%20Text
127.0.0.1 dockerhub.7x-networks.net
```

1. 安装内核5.16.0(必须大于4.8，略)
2. 安装docker-ce(用阿里云的docker-ce.repo)

```Bash
yum install docker-ce docker-ce-cli
```

1. 配置docker，修改/etc/docker/daemon.json文件，使用国内的源和自己，配置好data目录，启用experimental特性

```JSON
{

    "registry-mirrors": ["https://oy9g5j5n.mirror.aliyuncs.com"],       

    "bip": "192.168.117.1/24",

    "iptables":false,

    "data-root": "/data/docker",

    "insecure-registries": ["https://dockerhub.7x-networks.net"],

    "experimental": true

}
```

1. portal服务器上，修改/etc/hosts，通过内网接口访问docker仓库

```Plain%20Text
172.31.5.2 dockerhub.7x-networks.net
```

1. 检查是否启用experimental，`docker version`或者`docker system info`命令
2. 配置dnsmasq和iptables，允许容器访问公网(构建时apt/yum需要访问公网)，把容器内的dns解析，劫持到宿主机的dnsmasq进程

```Shell
# snat

iptables -t nat -A POSTROUTING -s 192.168.117.0/24 -o eth0 -j MASQUERADE

# dns劫持

iptables -t nat -A PREROUTING -p udp -m udp --dport 53 -j REDIRECT --to-ports 53
```

1. ucloud上云主机防火墙放开443端口

1. ### 拉取docker镜像

   1. 拉取镜像：

```Shell
 docker image pull registry
```

1. 安装模拟器

```Shell
 docker run --rm --privileged tonistiigi/binfmt:latest --install all
```

1. 创建新的builder，默认的builder不支持

```Shell
 docker buildx create --use --name=svxbuilder-cn --driver docker-container \ 

 --driver-opt image=dockerpracticesig/buildkit:master
```

1. 查看创建的builder

```Shell
docker buildx ls
```

1. ### 创建7x-registry容器

   1. docker + nginx方式，nginx实现认证和https功能，更换证书只需要重启nginx
      1. 创建docker

```Shell
# 默认5000端口

docker create --restart=always --name 7x-registry \

 -v /data/docker/registry:/var/lib/registry \

 -v /etc/localtime:/etc/localtime:ro \

 -e REGISTRY_STORAGE_DELETE_ENABLED=true \

 -p 5000:5000 registry
```

1. 容器内的dns服务器是8.8.8.8，构建的时候有时会解析不了
   1. 方式一：执行docker exec -it 7x-registry sh进入docker，把/etc/resolv.conf修改为nameserver 192.168.117.1。即宿主机的dnsmasq进程，每次重启容器都会重置为8.8.8.8，所以每次重启容器都需要重新修改
   2. 方式二：通过iptables规则拦截到宿主机dnsmasq
2. 上传7x-networks证书到/etc/ssl_certificate/7x-networks.net/
3. 在/etc/nginx目录下，创建nginx的基本认证的密码文件，然后输入2次密码

```Shell
htpasswd -c /etc/nginx/docker_htpasswd 7xnetwork.docker
```

1. nginx的配置文件

```Nginx
upstream docker_repo {

    server 127.0.0.1:5000;

}

# cpe上pull镜像通过https，443端口需要在防火墙上放开

server {

    listen       443 ssl;

    server_name  dockerhub.7x-networks.net;



    ssl_certificate   /etc/ssl_certificate/7x-networks.net/7x-networks.net.pem;

    ssl_certificate_key  /etc/ssl_certificate/7x-networks.net/7x-networks.net.key;

    ssl_session_timeout 5m;

    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;

    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;

    ssl_prefer_server_ciphers on;



    location / {

        auth_basic "Please enter Password";

        auth_basic_user_file docker_htpasswd;

        proxy_set_header Host $host;

        proxy_set_header X-Forwarded-For $remote_addr;

        proxy_pass http://docker_repo/;

    }

}



# 本地构建完成后通过http协议push到仓库，不打开认证，80端口不能在防火墙上放开

server {

    listen 80;

    server_name  dockerhub.7x-networks.net;

    location / {

#        auth_basic "Please enter Password";

#        auth_basic_user_file docker_htpasswd;

        proxy_set_header Host $host;

        proxy_set_header X-Forwarded-For $remote_addr;

        proxy_pass http://docker_repo/;

    }

}
```

1. docker仓库直接启用认证和https，缺点更换证书需要重启docker，后期不方便扩展(没采用)

```Shell
docker create --restart=always --name 7x-registry \

 -v /data/docker/cert:/certs \

 -v /data/docker/auth:/auth \

 -v /etc/localtime:/etc/localtime:ro \

 -v /data/docker/registry:/var/lib/registry \

 -e REGISTRY_STORAGE_DELETE_ENABLED=true \

 -e REGISTRY_HTTP_ADDR=0.0.0.0:443 \

 -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/cert.pem \

 -e REGISTRY_HTTP_TLS_KEY=/certs/privkey.pem \

 -e "REGISTRY_AUTH=htpasswd" \

 -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \

 -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \

 -p 443:443 registry
```

1. 使用已经弄好的dockerfile，构建镜像并push

```Shell
 docker login https://dockerhub.7x-networks.net

 # 构建arm

 docker buildx build --platform linux/amd64,linux/arm64 -t dockerhub.7x-networks.net/valor/7x-networks/dnsmasq:v1.0 . --push

 # 构建x86

 docker build -t dockerhub.7x-networks.net/valor/7x-networks/dnsmasq:v1.0 .

 docker push dockerhub.7x-networks.net/valor/7x-networks/dnsmasq:v1.0

 
```