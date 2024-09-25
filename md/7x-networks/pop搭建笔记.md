1. 拓扑图

   ![](/Users/lichenlu/Desktop/md/7x-networks/img/e75c54a1.svg)

   

2. 抛开路由协议的实现
   1. OS：ubuntu-server-20.4-LTS
   2. 内核模块
   3. ```undefined
      modprobe mpls_router
      modprobe mpls_gso
      modprobe mpls_iptunnel
      
      sysctl -w net.mpls.conf.lo.input=1
      sysctl -w net.mpls.platform_labels=1048575
      ```

   4. 内核参数
      1. R1
      2. ```undefined
         sysctl -w net.mpls.conf.ens34.input=1
         ```

      3. R2
      4. ```Plain
         sysctl -w net.mpls.conf.ens34.input=1
         sysctl -w net.mpls.conf.ens35.input=1
         sysctl -w net.mpls.conf.ens39.input=1
         ```

      5. R3
      6. ```Plain
         sysctl -w net.mpls.conf.ens34.input=1
         sysctl -w net.mpls.conf.ens35.input=1
         ```

      7. R4
      8. ```Plain
         sysctl -w net.mpls.conf.ens34.input=1
         sysctl -w net.mpls.conf.ens35.input=1
         sysctl -w net.mpls.conf.ens39.input=1
         ```

      9. R5
      10. ```Plain
          sysctl -w net.mpls.conf.ens34.input=1
          ```
   5. VRF
      1. R1
      2. ```undefined
         ip link add blue type vrf table 100
         ip link set blue up
         ip route add vrf blue unreachable default metric 4278198272
         ip -6 route add vrf blue unreachable default metric 4278198272
         
         ip link add blueveth0 type veth peer name blueveth1
         ip link set blueveth0 vrf blue 
         ip link set blueveth0 up
         ip addr add 10.1.1.1/24 dev blueveth0
         
         ip netns add nsblue
         ip link set blueveth1 netns nsblue
         ip netns exec nsblue ip addr add 10.1.1.2/24 dev blueveth1
         ip netns exec nsblue ip link set blueveth1 up
         ip netns exec nsblue ip route add default via 10.1.1.1
         ```

      3. R5
      4. ```Plain
         ip link add blue type vrf table 100
         ip link set blue up
         ip route add vrf blue unreachable default metric 4278198272
         ip -6 route add vrf blue unreachable default metric 4278198272
         
         ip link add blueveth0 type veth peer name blueveth1
         ip link set blueveth0 vrf blue 
         ip link set blueveth0 up
         ip addr add 10.1.5.1/24 dev blueveth0
         
         ip netns add nsblue
         ip link set blueveth1 netns nsblue
         ip netns exec nsblue ip addr add 10.1.5.2/24 dev blueveth1
         ip netns exec nsblue ip link set blueveth1 up
         ip netns exec nsblue ip route add default via 10.1.5.1
         ```
   6. P的转发配置
      1. R2
      2. ```undefined
         ip -f mpls route add 1102 dev lo
         ip -f mpls route add 1202 as 1202 via inet 172.16.12.2
         
         ip -f mpls route add 1201 dev lo
         ip -f mpls route add 1101 as 1101 via inet 172.16.11.1
         ```

      3. R3
      4. R4
      5. ```undefined
         ip -f mpls route add 1202 dev lo
         ip -f mpls route add 1301 as 1301 via inet 172.16.13.1
         
         ip -f mpls route add 1302 dev lo
         ip -f mpls route add 1201 as 1201 via inet 172.16.12.1
         ```
   7. PE的转发配置：
      1. R1
      2. ```undefined
         ip route add 10.1.5.0/24 encap mpls 1102/1202/1301/100 via inet 172.16.11.2 vrf blue
         
         ip -f mpls route add 1101 dev lo
         ip -f mpls route add 100 dev blue
         ```

      3. R5
      4. ```undefined
         ip route add 10.1.1.0/24 encap mpls 1302/1201/1101/100 via inet 172.16.13.2 vrf blue
         
         ip -f mpls route add 1301 dev lo
         ip -f mpls route add 100 dev blue
         ```

3. 基于路由协议的实现
   1. PE间三层网络互通
      1. PE初始化（路由模式、MPLS模式等内核参数）
      2. PE间点到点互联建立（点到点IP地址互通）
         1. 二层:  纯二层
         2. 三层：GRE
      3. PE上增加loopback地址
      4. PE上启用OSPF路由协议（loopback地址之间能相互ping通）
         1. ```undefined
            router ospf
             ospf router-id 100.64.240.10
             network 100.64.240.10/32 area 0   //loopback地址
             network 100.64.241.0/24 area 0
             capability opaque
             segment-routing on                //启用SR-MPLS
             segment-routing prefix 100.64.240.10/32 index 10
             router-info area
            exit
            ```

         2. OSPF router-id  lo地址，计算index值
            1. (0-65535)  Index value inside SRGB
         3. OSPF network 是所有链路的地址段
   2. 调通SR/MPLS网络（BGP）
      1. ```undefined
         router bgp 65000
          bgp router-id 100.64.240.9
          no bgp     ebgp    -requires-policy
          no bgp suppress-duplicates
          no bgp default ipv4-unicast
          no bgp network import-check
          neighbor 100.64.240.3 remote-as 65000
          neighbor 100.64.240.3 update-source 100.64.240.9
          neighbor 100.64.240.5 remote-as 65000
          neighbor 100.64.240.5 update-source 100.64.240.9
          !
          address-family ipv4 unicast
           neighbor 100.64.240.3 activate
           neighbor 100.64.240.3 soft-reconfiguration inbound
           neighbor 100.64.240.5 activate
           neighbor 100.64.240.5 soft-reconfiguration inbound
          exit-address-family
          !
          address-family ipv4 vpn
           neighbor 100.64.240.3 activate
           neighbor 100.64.240.5 activate
          exit-address-family
         exit
         ```

      2. BGP router-id  lo地址
      3. BGP要建立 full-mesh（RR）
      4. BGP 邻居用对端lo地址进行配置
   3. CE接入PE打通CE间的网络