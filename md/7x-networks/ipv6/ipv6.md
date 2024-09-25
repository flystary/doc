1. ### 功能验证

   1. #### 内核参数

      1. linux上打开ipv6功能、转发功能，以及自动配置(只有WAN角色的网卡accept_ra=2，其他角色的网卡要设置成0)：
         1. 

         2. ```Plain
            net.ipv6.conf.all.disable_ipv6 = 0
            net.ipv6.conf.default.disable_ipv6 = 0
            net.ipv6.conf.all.forwarding = 1
            net.ipv6.conf.default.forwarding = 1
            net.ipv6.conf.all.autoconf = 1
            net.ipv6.conf.all.accept_ra = 2
            net.ipv6.conf.default.autoconf = 1
            net.ipv6.conf.default.accept_ra = 2
            ```

   1. ### WAN区域

      1. WAN口ipv6配置信息
         1. 配置方式：
            1. centos通过ifcfg文件配置
            2. ubuntu通过nmcli配置
         2. 支持的方式：
            1. 通过PPPOE获取
               1. 暂时没环境验证
            2. 通过DHCP/SLAAC获取
               1. 参考：[IPv6 自动寻址(DHCPv6/SLAAC)](https://7x-networks.feishu.cn/docx/O8pwdvywsoHRlZxxq9WcEsdSnke?theme=LIGHT&contentTheme=DARK) 
            3. 通过STATIC配置，类似IPv4
      2. wan探测策略路由
         1. 类似IPv4通过IPv6策略路由方式来探测，参考(基于IPv6的策略路由)

   1. ### LAN区域

      1. DHCPv6 SERVER
         1. 参考：[IPv6 自动寻址(DHCPv6/SLAAC)](https://7x-networks.feishu.cn/docx/O8pwdvywsoHRlZxxq9WcEsdSnke?theme=LIGHT&contentTheme=DARK) 
      2. DNS server
         1. 参考：[IPv6 自动寻址(DHCPv6/SLAAC)](https://7x-networks.feishu.cn/docx/O8pwdvywsoHRlZxxq9WcEsdSnke?theme=LIGHT&contentTheme=DARK) 
         2. dnsmasq不支持DNS64功能：
            1. 修改DNSMASQ来支持
            2. 使用其他工具替代，改动比较大
      3. 发送RA的工具：
         1. dnsmasq
         2. radvd

   1. ### NAT

      1. IPv6->IPv6
         1. 类似ipv4，通过ip6tables添加IPv6的snat规则：支持MASQUERADE和SNAT
         2. ```Bash
            ip addr add fd00:1:2:3::22/64 dev vlan7
            ip6tables -t nat -A POSTROUTING -s fd00:1:2:3::/64 -o enp2s0 -j MASQUERADE
            
            ping -6 64:ff9b::223.5.5.5 -I fd00:1:2:3::22
            PING 64:ff9b::223.5.5.5(64:ff9b::df05:505) from fd00:1:2:3::22 : 56 data bytes
            64 bytes from 64:ff9b::df05:505: icmp_seq=1 ttl=114 time=7.20 ms
            64 bytes from 64:ff9b::df05:505: icmp_seq=2 ttl=114 time=5.76 ms
            64 bytes from 64:ff9b::df05:505: icmp_seq=3 ttl=114 time=5.95 ms
            ```
      2. IPv6->IPv4
         1. 需要配置DNS64，参考: [IPv6 NAT64](https://7x-networks.feishu.cn/docx/IcNIdamlIomchzxWlROcd5tinBQ?theme=LIGHT&contentTheme=DARK) 

   1. ### 基于IPv6的ACL

      1. 类似于IPv4：通过ip6tables实现

      ```Bash
      # 添加规则后可以不能ping通，删除规则后可以ping通
      ip6tables -A FORWARD -d 64:ff9b::223.5.5.5/128 -j DROP
      ```

      1. ipset支持IPv6

      ```Bash
      ipset create test hash:ip family inet6
      ip6tables -A OUTPUT -m set --match-set test dst -j DROP
      # 添加ipv6地址到ipset集合，添加成功后，ping不通
      ipset add test 240e:390:340:86b0::10
      # 删除ipset集合里的ipv6地址，删除后，可以ping通
      ipset del test 240e:390:340:86b0::10
      ```

   1. ### 基于IPv6的QOS

      1. 类似于IPv4：通过ip6tables实现

      ```Bash
      # 添加tc规则
      interface=fm1-mac5
      tc qdisc del dev $interface root 2>/dev/null
      tc qdisc add dev $interface root handle 1: htb default 255 2>/dev/null
      tc class add dev $interface parent 1: classid 1:1 htb rate 1000mbit ceil 1000mbit
      tc class add dev $interface parent 1:1 classid 1:255 htb rate 50kbit ceil 1000mbit prio 7
      tc class replace dev $interface parent 1:1 classid 1:11 htb rate 2048kbit ceil 2048kbit prio 3
      tc qdisc replace dev $interface parent 1:11 handle 11: sfq perturb 10
      tc filter add dev $interface parent 1:0 protocol ipv6 handle 0xb00000/0x0ff00000 fw classid 1:11
      
      # 通过iptables, 打mark（0xb00000）
      ip6tables -t mangle -A POSTROUTING -d 240e:390:340:86b1::10/128 -p tcp --dport 5201 -j MARK --set-mark 0xb00000/0x0ff00000
      ```

      

   1. ### 基于IPv6的策略路由

      1. 类似于IPv4：
         1. 通过ip -6 rule添加规则
         2. 通过ip -6 route add table xxx 添加路由
         3. ```Bash
            ip -6 rule add from fd00:1:2:3::22 table 201
            # table 201加路由走wan2，抓包确定是走enp2s0
            ip -6 route add ::/0 via fe80::fbfe:8ecb:396c:4700 dev enp2s0 table 20
            # table 201加路由走vlan10，抓包确定是走vlan10
            ip -6 route add ::/0 via fe80::64f4:5ff:fece:1 dev vlan10 table 201
            ```

   1. ### POP支持IPv6

   1. ### 基于IPv6的加速业务

      1. 暂时没验证
         1. 类似于IPv4：通过策略路由和ip6tables SNAT实现
         2. 需要POP支持

   1. ### 基于IPv6的经典互联业务

      1. 遗留问题： arm没编译IPv6隧道相关内核模块，ipip/gre都不支持
      2. 通过ipv4建立GRE隧道，内层ip头部可以是ipv4也可以是ipv6，优点：不需要PE支持IPv6，ipv4和ipv6可以共用一个隧道

      ```Bash
      ip tunnel add test64 mode gre local 172.18.1.1 remote 172.18.1.191 ttl 255
      ip link set test64 mtu 1410
      ip link set test64 up
      ip addr add dev test64 fd00:2:2:2::1/64
      # 可以配置ipv4地址
      ip addr add dev test64 100.100.100.1/30
      ```

      1. 通过ipv6建立GRE隧道，内层ip头部可以是ipv4也可以是ipv6，缺点：需要PE支持IPv6

      ```Bash
      ip tunnel add test66 mode ip6gre local 240e:390:340:86b0::1 remote 240e:390:340:86b0::10
      ip link set test66 mtu 1410
      ip link set test66 up
      ip addr add dev test66 fd00:2:2:2::1/64
      # 可以配置ipv4地址
      ip addr add dev test66 100.100.100.1/30
      ```

      1. 通过ipv6建立ip6tnl隧道(mode any)，内层ip头部可以是ipv4也可以是ipv6，缺点：需要PE支持IPv6

      ```Bash
      ip link add test66 type ip6tnl mode any local 240e:390:340:86b0::1 remote 240e:390:340:86b0::10
      ip link set test66 mtu 1410
      ip link set test66 up
      ip addr add dev test66 fd00:2:2:2::1/64
      # 可以配置ipv4地址
      ip addr add dev test66 100.100.100.1/30
      ```

      1. 通过ipv6建立ip6ip6隧道，内层ip头部只能是ipv6（ipip6只支持IPv6封装IPv4），缺点：需要PE支持IPv6，而且内层ip头不支持ipv4

      ```Bash
      ip tunnel add test66 mode ip6ip6 local 240e:390:340:86b0::1 remote 240e:390:340:86b0::10
      ip link set test66 mtu 1410
      ip link set test66 up
      # 配置ipv4不通
      ip addr add dev test66 fd00:2:2:2::1/64
      ```

      1. Tunnel参考文档: https://developers.redhat.com/blog/2019/05/17/an-introduction-to-linux-virtual-interfaces-tunnels 

   1. ### 动态路由OSPF/BGP

      1. [ipv6 BGP](https://7x-networks.feishu.cn/docx/ZRLGdmyhjoNGyGx0NSBcyETknid) 
      2. OSPF需要修改/etc/frr/daemons，启动osfp6d进程

   1. ### HA: keepalived

      1. 仅IPv6模式，正常切换

      ```Bash
      vrrp_instance enp1s0 {
          state MASTER
          interface enp1s0
          virtual_router_id 139
          priority 10
          advert_int 0.3
          version 3
          authentication {
              auth_type PASS
              auth_pass 7network
          }
          virtual_ipaddress {
              240e:390:340:86b0::20/64
          }
          unicast_src_ip 240e:390:340:86b0::10
          unicast_peer {
              240e:390:340:86b0::1
          }
      }
      ```

      1. ipv4和ipv6混合模式，需要同时配置一个ipv4和一个ipv6的vip，正常切换

      ```Bash
      vrrp_instance test2 {
          state MASTER
          interface enp1s0
          virtual_router_id 140
          priority 10
          advert_int 0.3
          version 3
          authentication {
              auth_type PASS
              auth_pass 7network
          }
          virtual_ipaddress {
              10.168.1.62/24
          }
          virtual_ipaddress_excluded {
              240e:390:340:86b0::21/64
          }
          unicast_src_ip 10.168.1.200
          unicast_peer {
              10.168.1.114
          }
      }
      ```

   1. ## portal和em都支持ipv6

   uc上申请ipv6地址

   阿里云上，portal和em配置AAAA记录

   ### 问题

   1. pop不支持ipv6
   2. dnsmasq不支持dns64

   ### cpe wan配置：

   1. centos要添加配置：
      1. /etc/sysconfig/network文件里添加：IPV6FORWARDING=yes
      2. 静态配置ipv6：
         1. ```Bash
            IPV6INIT=yes
            IPV6_AUTOCONF=yes
            IPV6_FAILURE_FATAL=no
            IPV6FORWARDING=yes
            IPV6ADDR=2001:db8:20::2/64
            IPV6_DEFAULTGW=2001:db8:20::1/64
            ```
      3. DHCP
         1. ```Bash
             IPV6INIT=yes
             DHCPV6C=yes
            ```
   2. ubuntu：通过nmcli直接可以配置，类似于ipv4