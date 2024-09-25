# NAT64服务软件的选择：

1. **tayga**
   1. 用户态程序，yum/apt直接安装，安装方便，可信度高
   2. tayga不停地从tun字符设备读取IPv6或者IPv4报文，经过6-4/4-6转换后再发回该tun字符设备
   3. 需要配合tun和路由、snat
   4. tun编程有积累，后期有修改源码的可能
   5. 无状态NAT64
2. Jool
   1. 内核模块，编译、加载和卸载比较麻烦
   2. 配合iptables使用比较麻烦
   3. 无状态NAT64

# 使用tayga实现NAT64

1. 安装tayga
   1. Centos
      1. ```Bash
         yum install epel-release
         yum install tayga
         ```
   2. Ubuntu
   3. ```Bash
      apt install tayga
      ```
2. 配置tayga
   1. 配置/etc/tayga.conf，查看tayga配置：`grep '^\w' /etc/tayga.conf` 
   2. ```Plain
      tun-device nat64
      ipv4-addr 192.168.255.1
      ipv6-addr fdaa:bb:1::1
      prefix 64:ff9b::/96
      dynamic-pool 192.168.255.0/24
      data-dir /var/spool/tayga
      ```

   3. /etc/default/tayga
   4. ```SQL
      CONFIGURE_IFACE="yes"
      CONFIGURE_NAT44="yes"
      DAEMON_OPTS=""
      IPV4_TUN_ADDR="192.168.255.1"
      IPV6_TUN_ADDR="fdaa:bb:1::1"
      ```
3. 启动tayga
   1. ```SQL
      systemctl restart tayga
      systemctl enable tayga
      ```
4. 启动tayga后，tayga对系统做了一些网络相关的配置修改：
   1. 创建了一个nat64的tun网卡，并配置了一个ipv4和ipv6的地址（跟配置文件里一样）
   2. 添加了2条路由：
      1. 192.168.255.0/24 dev nat64 scope link
      2. 64:ff9b::/96 dev nat64
   3. 添加了一条iptables的nat规则
      1. tptables -t nat -A POSTROUTING -s 192.168.255.0/24 -j MASQUERADE
5. 测试：
   1. 下挂一台设备ping 64:ff9b::223.5.5.5，测试ping正常
   2. 64:ff9b::223.5.5.5，这个ip是怎么来的，就是配置文件里配置的64:ff9b::/96前缀，加上ipv4的32，刚刚好是128位，凑成一个ipv6的ip地址。
   3. 对不提供ipv6服务的业务，如果通过NAT64去访问，需要跟DNS64结合

# DNS64

1. dnsmasq不支持DNS64，暂时使用bind9实现，增加一个nat64的配置：（64:ff9b::/96要跟NAT64的配置对应）

   1. ```Python
      dns64 64:ff9b::/96 {
          clients { any; };
      };
      ```

2. ping测试

   ![](/Users/lichenlu/Desktop/md/7x-networks/img/18af0f6a.png)

ipv4 ping测试正常，解析到的ip是106.75.254.202

ipv6 ping测试正常，解析到的ip是64:ff9b::6a4b:feca，6a4b:feca对应ipv4的106.75.254.202

64:ff9b::对应的NAT64和DNS64里的ipv6前缀配置