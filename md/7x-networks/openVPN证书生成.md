OSS部署VPN证书生成脚本部署(easyrsa3) 

1. 下载easyrsa3(或者成其他OSS服务器上复制)，[GitHub - OpenVPN/easy-rsa: easy-rsa - Simple shell based CA utility](https://github.com/OpenVPN/easy-rsa/tree/master)
2. 下载后，复制2份，/opt/easyrsa3-client和/opt/easyrsa3-server
3. `cd /opt/easyrsa3-server;cp vars.example vars`，然后修改vars里的参数

```bash
set_var EASYRSA_REQ_COUNTRY     "CN"

set_var EASYRSA_REQ_PROVINCE    "ShangHai"

set_var EASYRSA_REQ_CITY        "ShangHai"

set_var EASYRSA_REQ_ORG         "7x-networks"

set_var EASYRSA_REQ_EMAIL       "service@7x-networks.com"

set_var EASYRSA_REQ_OU          "7x-networks"

set_var EASYRSA_CA_EXPIRE       36500

set_var EASYRSA_CERT_EXPIRE     36500
```

1. 执行下面的脚本，根据提示输入"yes"，其他默认

```Bash
./easyrsa init-pki # 初始化PKI

./easyrsa build-ca nopass # 无密码方式创建ca

./easyrsa gen-req server nopass # 创建服务端key文件

./easyrsa sign server server # 注册服务端CN名，生产服务端crt文件

./easyrsa gen-dh dh.pem #文件生产
```

1. server证书路径：

```Bash
# server端相关证书放到/etc/openvpn：

mkdir -p /etc/openvpn

cp /opt/easyrsa3-server/pki/ca.crt /etc/openvpn

cp /opt/easyrsa3-server/pki/issued/server.crt /etc/openvpn

cp /opt/easyrsa3-server/pki/private/server.key /etc/openvpn

cp /opt/easyrsa3-server/pki/dh.pem /etc/openvpn

# 给文件增加可读权限

chown +r /etc/openvpn/*
```

1. 创建/usr/local/svxnetworks目录，把7xcli_gen_openvpn_client_certificate脚本复制打/usr/local/svxnetworks目录下，OSS会调用这个脚本自动生成盒子的证书
2. 安装expect ，yum -y install expect
3. 执行/usr/local/svxnetworks/7xcli_gen_openvpn_client_certificate --sn test看看能不能正常生成test对应的证书。