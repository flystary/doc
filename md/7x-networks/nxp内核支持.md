nxp需要添加的内核支持

### **内核参数修改**

#### 

1. Networking support:

- networking options
  - TCP/IP networking
    - IP: advanced router
      - FIB TRIE statisctiscs
      - IP : policy routing
      - IP: equal cost multipath
      - IP: verbose route monitoring
  - IP: tunneling
  - IP: GRE demultiplerxer
  - IP: GRE tunnels over ip（IP: broadcast GRE over ip）
  - IP: TCP syncookie support
  - Virtual (secure) IP: tunneling
  - Network packet filtering framwork (Netfilter)
  - MultiProtocol Label Switching（以及2个子选项）

1. Netfliter:

- IP set support
- core Netfilter Configuration
  - FTP protocol support
  - TFTP protocol support
  - Connection tracking netlink interface
  - Netfliter flow table module
  - nfmakr/ctmakr target and support
  - set target and match support
  - DSCP and TOS target support
  - LOG target support
  - MARK target support
  - NFLOG target support
  - NFQUEUE target support
  - REDIRECT target support
  - TCPMSS target support
  - commet match support
  - connlimit match support
  - connlmark connection mark match support
  - conntrack connection tracking match support
  - dccp protocol match support
  - dscp and tos match support
  - esp match support
  - iprange address range match support
  - length match support
  - hashlimit match support（M）
  - mac address match support
  - mark match support
  - multiport match support
  - IPsec policy match support
  - state match support
  - statistics match support
  - string match support
  - tcpmss match support
  - time match support
  - u32 match support
  - Netfilter nf_tables support
    - conntrack
    - connlimit 
    - log
    - limit
    - masquerade
    - redirect
    - nat
    - quata
    - reject
    - hash
    - xfrm/IPSec
    - tporxy
- IP Netfilter Configuration
  - Netfilter flow table IPv4 module
  - IP tables support
    - 除CLUSTERIP target和Security table外全选



ppp over Ethernet

bond

team