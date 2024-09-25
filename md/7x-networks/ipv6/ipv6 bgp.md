1. 使用BGP配置
   1. CPE A的frr里增加一下ipv6配置，绿色标记的部分，ipv4的配置是原来就支持的ipv4配置：
   2. ```Bash
      router bgp 65001
       bgp router-id 192.168.200.81
       no bgp ebgp-requires-policy
       no bgp default ipv4-unicast
       no bgp network import-check
       timers bgp 3 9
       neighbor 192.168.200.82 remote-as 65002
       neighbor fe80::2b67:8a6d:67d1:c43f remote-as 65002
       neighbor fe80::2b67:8a6d:67d1:c43f interface ens37
       !
       address-family ipv4 unicast
        network 192.168.10.0/24
        neighbor 192.168.200.82 activate
        neighbor 192.168.200.82 soft-reconfiguration inbound
       exit-address-family
       !
       address-family ipv6 unicast
        network fd00:6:6:1::/64
        neighbor fe80::2b67:8a6d:67d1:c43f activate
        neighbor fe80::2b67:8a6d:67d1:c43f soft-reconfiguration inbound
       exit-address-family
      ```

   3. CPE B的frr里增加一下ipv6配置，绿色标记的部分，ipv4的配置是原来就支持的ipv4配置：
   4. ```Bash
      router bgp 65002
       bgp router-id 192.168.200.82
       no bgp ebgp-requires-policy
       no bgp default ipv4-unicast
       no bgp network import-check
       timers bgp 3 9
       neighbor 192.168.200.81 remote-as 65001
       neighbor 192.168.200.81 ebgp-multihop 2
       neighbor fe80::292e:b940:3804:e5d0 remote-as 65001
       neighbor fe80::292e:b940:3804:e5d0 interface ens37
       !
       address-family ipv4 unicast
        network 192.168.20.0/24
        neighbor 192.168.200.81 activate
        neighbor 192.168.200.81 soft-reconfiguration inbound
       exit-address-family
       !
       address-family ipv6 unicast
        network fd00:6:6:2::/64
        neighbor fe80::292e:b940:3804:e5d0 activate
        neighbor fe80::292e:b940:3804:e5d0 soft-reconfiguration inbound
       exit-address-family
      ```
2. 状态查看：
   1. show ip bgp summary和show ip bgp ipv6 summary
   2. ```Plain
      net1# show ip bgp summary 
      
      IPv4 Unicast Summary (VRF default):
      BGP router identifier 192.168.200.81, local AS number 65001 vrf-id 0
      BGP table version 4
      RIB entries 3, using 576 bytes of memory
      Peers 1, using 20 KiB of memory
      
      Neighbor        V         AS   MsgRcvd   MsgSent   TblVer  InQ OutQ  Up/Down State/PfxRcd   PfxSnt Desc
      192.168.200.82  4      65002       527       528        4    0    0 00:24:13            1        2 N/A
      
      Total number of neighbors 1
      net1# 
      net1# show ip bgp ipv6 summary 
      
      IPv6 Unicast Summary (VRF default):
      BGP router identifier 192.168.200.81, local AS number 65001 vrf-id 0
      BGP table version 5
      RIB entries 3, using 576 bytes of memory
      Peers 1, using 20 KiB of memory
      
      Neighbor                  V         AS   MsgRcvd   MsgSent   TblVer  InQ OutQ  Up/Down State/PfxRcd   PfxSnt Desc
      fe80::2b67:8a6d:67d1:c43f 4      65002       527       528        5    0    0 00:24:15            1        2 N/A
      
      Total number of neighbors 1
      ```

   3. show ip bgp和show ip bgp ipv6
   4. ```Plain
      net1# show ip bgp 
      BGP table version is 4, local router ID is 192.168.200.81, vrf id 0
      Default local pref 100, local AS 65001
      Status codes:  s suppressed, d damped, h history, * valid, > best, = multipath,
                     i internal, r RIB-failure, S Stale, R Removed
      Nexthop codes: @NNN nexthop's vrf id, < announce-nh-self
      Origin codes:  i - IGP, e - EGP, ? - incomplete
      RPKI validation codes: V valid, I invalid, N Not found
      
          Network          Next Hop            Metric LocPrf Weight Path
       *> 192.168.10.0/24  0.0.0.0                  0         32768 i
       *> 192.168.20.0/24  192.168.200.82           0             0 65002 i
      
      Displayed  2 routes and 2 total paths
      net1# show ip bgp ipv6 
      BGP table version is 5, local router ID is 192.168.200.81, vrf id 0
      Default local pref 100, local AS 65001
      Status codes:  s suppressed, d damped, h history, * valid, > best, = multipath,
                     i internal, r RIB-failure, S Stale, R Removed
      Nexthop codes: @NNN nexthop's vrf id, < announce-nh-self
      Origin codes:  i - IGP, e - EGP, ? - incomplete
      RPKI validation codes: V valid, I invalid, N Not found
      
          Network          Next Hop            Metric LocPrf Weight Path
       *> fd00:6:6:1::/64  ::                       0         32768 i
       *> fd00:6:6:2::/64  fe80::2b67:8a6d:67d1:c43f
                                                   0             0 65002 i
      
      Displayed  2 routes and 2 total paths
      ```
3. 配置说明：
   1. no bgp default ipv4-unicast
      1. 默认frr只会向邻居宣告ipv4的地址，使用"no bgp default ipv4-unicast"覆盖此默认配置，因此需要显式启用所有地址族(address-family ipv4/ipv6)
   2. neighbor fe80::292e:b940:3804:e5d0 interface ens37
      1. 如果使用的是link-local地址，必选指定接口名称
      2. 如果使用不是link-local地址，这个不需要配置
   3. network fd00:6:6:2::/64
      1. 必须在"address-family ipv6 unicast"里配置
   4. no bgp network import-check和no bgp ebgp-requires-policy，原来ipv4就有的配置，参考frr说明文档：
   5. ![img](/Users/lichenlu/Desktop/md/7x-networks/img/jjj.png)

   6. ![img](/Users/lichenlu/Desktop/md/7x-networks/img/kkkkk.png)