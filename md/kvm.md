# KVM

## 实现

### 名词解释

​	KVM是Linux开源社区大力支持的虚拟化技术，基于Intel和AMD的硬件虚拟化技术。KVM（Kernel-bashdVirtual Machine，即基于内核的虚拟机），它是用于Linux内核中的虚拟化环境设施，是Linux内核中的一个功能模块，在Linux内核中默认被安装，可以将Linux内核转化为一个 Hypervisor。

### 何为qemu-kvm

​	所谓kvm技术中，应用到两个技术，分别是：qemu、kvm，其中kvm负责CPU虚拟化+内存虚拟化，实现了CPU和内存的虚拟化，但kvm不能模拟其他设备；qemu是模拟IO设备（网卡、磁盘等），kvm加上qemu之后就能实现真正的服务器虚拟化，故称之为qemu-kvm。由于kvm技术已经相当成熟，并且对很多方面都进行了隔离，但是像网卡、磁盘等设备依然无法虚拟出真是的机器。qemu-kvm补充了kvm技术的不足，而且在性能上对kvm进行了优化。

### 检查是否支持

#### 判断是否支持kvm，数量需要大于0

```bash
grep -Eoc '(vmx|svm)' /proc/cpuinfo //数字大于0，则代表CPU支持硬件虚拟化，反之则不支
```

#### 查询本机是否支持虚拟化技术

```bash
cat /proc/cpuinfo | grep -E 'vmx|svm'
```

#### 检查VT是否在BIOS开启

```bash
apt install cpu-checker //检查 VT 是否在 BIOS 中启用
kvm-ok //如果处理器虚拟化能力没有在 BIOS 中被禁用会输出如下信息
		INFO: /dev/kvm exists
		KVM acceleration can be used
```

### 网络模式

支持多种网络模式，包括Bridge（桥接模式）、NAT（NAT模式）、Router（路由模式）、host-only（隔离模式）

#### 桥接模式

在这种模式下，所有虚拟机都好像与主机物理机器在同一个子网内。同一物理网络中的所有其他物理机器都知道这些虚拟机，并可以访问这些虚拟机。桥接操作在OSI网络模型的第2层进行。

#### NAT模式

默认情况下，虚拟网络交换机以NAT模式运行。使用IP伪装技术，连接的guest虚拟机可以使用主机物理机器IP地址与任何外部网络进行通信。默认情况下，虚拟网络交换机在NAT模式下运行时，放置在主机物理机外部的计算机无法与其中的guest虚拟机进行通信。

#### 路由模式

当使用路由模式时，虚拟交换机连接到连接到主机物理机器的物理LAN，在不使用NAT的情况下来回传输流量。所有虚拟机都位于其自己的子网中，通过虚拟交换机进行路由。这种情况并不总是理想的，因为物理网络上的其他主机物理机器不通过手工配置的路由信息是没法发现这些虚拟机，并且不能访问虚拟机。

#### 隔离模式

在这种模式下，连接到虚拟交换机的虚拟机可以相互通信，也可以与主机物理机通信，但其通信不会传到主机物理机外，也不能从主机物理机外部接收通信。这种模式下使用dnsmasq对于诸如DHCP的基本功能是必需的。

## 安装

#### 所需软件

qemu-kvm：完整的虚拟化平台

libvirtd：用于硬件虚拟化的开源API、守护进程与管理工具。

virt-manager：虚拟机管理器

#### 使用apt安装

```bash
apt install -y qemu-kvm libvirt-daemon-system libvirt-clients bridge-utils virtinst virt-manager
```

#### 判断程序是否运行
```bash
systemctl is-active libvirtd
```

#### 检查安装结果

```bash
lsmod | grep kvm
```

#### 配置用户组

将当前用户添加到`libvirt`和`kvm`用户组，以便在不使用`sudo`的情况下管理虚拟机

```bash
sudo usermod -aG libvirt $USER
sudo usermod -aG kvm $USER
```

#### 启用并设置开机自启

```bash
systemctl start libvirtd 
systemctl enable libvirtd
```

#### 打印启动虚拟化和设置开机自启情况

```bash
systemctl list-unit-files |grep libvirtd.service
```

## 使用

### 网络

#### 桥接网络

##### Netplan

使用此方式可以直接创建网桥以及配置IP地址

###### 创建配置文件

```bash
touch /etc/netplan/01-network-manager-all.yaml
```

###### 文件写入以下内容

```yaml
# Let NetworkManager manage all devices on this system
network:
  version: 2
  #renderer: NetworkManager
  ethernets:
        ens33:(网卡名称)
            dhcp4: false
            dhcp6: false
  bridges:
        br0:
            addresses: [172.16.10.10/24]（本机IP）
            gateway4: 172.16.10.1（网关）
            nameservers:
                addresses: [172.16.10.1, 114.114.114.114]（DNS）
                search: [msnode]
            interfaces: [ens33]（网卡名称）
```

###### 启用配置并重启网络

```bash
netplan apply
```



##### NetworkManager

###### 安装network-manager

```bash
apt update
apt install network-manager -y
```

###### 启用并设置开机自启

```bash
systemctl enable NetworkManager
systemctl start NetworkManager
```

###### 创建10-globally-managed-devices.conf文件

```bash
touch /etc/NetworkManager/conf.d/10-globally-managed-devices.conf
```

###### 重启network-manager

```bash
systemctl restart NetworkManager
```

需要创建上述文件后NetworkManager才能完全管理网卡

###### 创建网桥并将网卡放入其中

使用brctl创建br1网桥并将网卡放入网桥中

```bash
brctl addbr br1
brctl addif br1 ens2f3
```

使用nmcli创建br1网桥并将ens2f3网卡放入br1网桥中

```bash
nmcli connection add type bridge con-name br1 ifname br1
nmcli connection modify br1 bridge.stp off
nmcli connection modify br1 bridge.forward-delay 4
nmcli connection modify ens2f3 master br1
nmcli connection modify br1 ipv4.addresses 172.16.30.80/24
nmcli connection modify br1 ipv4.gateway 172.16.30.1
nmcli connection modify br1 ipv4.dns 223.5.5.5
nmcli connection modify br1 ipv4.method manual
nmcli connection up br1
nmcli connection up ens2f3
```

###### 创建kvm桥接网络配置文件

```bash
cat > /tmp/kbr0.xml << EOF
<network>
  <name>kbr0</name>
  <forward mode='bridge'/>
  <bridge name='br1'/>
</network>
EOF
```

###### 创建桥接网络

```bash
virsh net-define --file  /tmp/kbr0.xml
virsh net-start --network kbr0
virsh net-autostart --network kbr0
```

#### 其他网络

vethMac

##### 创建直连网卡配置文件

```bash
cat > /tmp/enp0s8.xml << EOF
<network>
  <name>enp0s8</name>
  <uuid>01a45bc6-577a-4e8a-82be-785e8a35c9ba</uuid>
  <forward dev='enp0s8' mode='bridge'>
    <interface dev='enp0s8'/>
  </forward>
</network>
EOF
```

##### 创建网络

```bash
virsh net-undefine --file /tmp/enp0s8.xml
virsh net-define --file /tmp/enp0s8.xml
virsh net-start enp0s8
virsh net-autostart enp0s8
```

### 磁盘

#####  创建一个空的 qcow2磁盘

```bash
qemu-img create -f qcow2 /data/kvm/volume/ubuntu.qcow2 20G
```

### 虚机

##### 创建镜像虚机

```bash
virt-install --name=u1 --vcpus=1 --memory=2048 \
--boot hd --os-type=generic \
--disk path=/data/kvm/volume/ubuntu.qcow2,format=qcow2,size=50,bus=virtio \
--cdrom=/data/kvm/iso/ubuntu-20.04.6-live-server-amd64.iso \
--network network=kbr0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5901 \
--noautoconsole
```

停止虚机

```bash
virsh shutdown u2
```

##### 虚机空间清理

```bash
virt-syprep -d u1
```

##### 压缩虚机镜像

```bash
virt-sparsify --compress  /data/kvm/volume/ubuntu.qcow2   /data/kvm/volume/ubuntu.base.qcow2
```

##### 复制镜像

```bash
cp ubuntu.base.qcow2  ubuntu2.qcow2
cp ubuntu.base.qcow2  ubuntu3.qcow2
cp ubuntu.base.qcow2  ubuntu4.qcow2
cp ubuntu.base.qcow2  ubuntu5.qcow2
```

##### 使用镜像创建虚机

```bash
virt-install --name=u2 --vcpus=1 --memory=2048 \
--boot hd --os-type=generic \
--disk path=/data/kvm/volume/ubuntu2.qcow2,format=qcow2,size=50,bus=virtio \
--network network=kbr0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5902 \
--noautoconsole
```

```bash
virt-install --name=u3 --vcpus=1 --memory=2048 \
--boot hd --os-type=generic \
--disk path=/data/kvm/volume/ubuntu3.qcow2,format=qcow2,size=50,bus=virtio \
--network network=kbr0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5903 \
--noautoconsole
```

```bash
virt-install --name=u4 --vcpus=1 --memory=2048 \
--boot hd --os-type=generic \
--disk path=/data/kvm/volume/ubuntu4.qcow2,format=qcow2,size=50,bus=virtio \
--network network=kbr0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5904 \
--noautoconsole
```

```bash
virt-install --name=u5 --vcpus=1 --memory=2048 \
--boot hd --os-type=generic \
--disk path=/data/kvm/volume/ubuntu5.qcow2,format=qcow2,size=50,bus=virtio \
--network network=kbr0,model=virtio \
--graphics vnc,listen=0.0.0.0,port=5905 \
--noautoconsole
```

##### 设置虚机开机自启

```bash
virsh autostart u2
```