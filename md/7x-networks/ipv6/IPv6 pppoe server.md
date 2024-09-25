1. 安装：
   1. 参考安装说明文档：https://docs.accel-ppp.org/installation/centos.html
   2. 不需要支持IPOE和VLAN，cmake的时候，指定参数，不使用他建议的：
   3. ```Bash
      # DBUILD_IPOE_DRIVER=FALSE
      # DBUILD_VLAN_MON_DRIVER=FALSE
      # centos
      cmake -DBUILD_IPOE_DRIVER=FALSE -DBUILD_VLAN_MON_DRIVER=FALSE -DCMAKE_INSTALL_PREFIX=/usr -DLUA=TRUE -DCPACK_TYPE=Centos7 ..
      # ubuntu 20.04
      cmake -DBUILD_IPOE_DRIVER=FALSE -DBUILD_VLAN_MON_DRIVER=FALSE -DCMAKE_INSTALL_PREFIX=/usr -DLUA=TRUE -DCPACK_TYPE=Ubuntu20 ..
      ```
2. 配置：
   1. https://docs.accel-ppp.org/configuration/ipv6-nd.html
   2. ipv6-nd的参数，参考文档：[IPv6 自动寻址(DHCPv6/SLAAC/DHCPv6-PD)](https://7x-networks.feishu.cn/docx/O8pwdvywsoHRlZxxq9WcEsdSnke?theme=LIGHT&contentTheme=DARK) 
   3. 账号管理配置文件：/etc/ppp/chap-secrets
   4. 办公室配置文件：/etc/accel-ppp.conf

```Bash
[modules]
log_file
#log_syslog
#log_tcp
#log_pgsql
#pptp
#l2tp
#sstp
pppoe
#ipoe
auth_mschap_v2
auth_mschap_v1
auth_chap_md5
auth_pap
#radius
chap-secrets
# IPv4 address assigning module.
ippool
#pppd_compat
#shaper
#net-snmp
#logwtmp
#connlimit
# IPv6 Neighbor Discovery module.
ipv6_nd
# IPv6 DHCP module.
ipv6_dhcp
# IPv6 address assigning module.
ipv6pool

[core]
log-error=/var/log/accel-ppp/core.log
thread-count=4

[common]
#single-session=replace
#sid-case=upper
#sid-source=seq
#max-sessions=1000

[ppp]
verbose=1
min-mtu=1280
mtu=1400
mru=1400
#accomp=deny
#pcomp=deny
#ccp=0
#check-ip=0
#mppe=require
ipv4=require
ipv6=allow
ipv6-intf-id=0:0:0:1
ipv6-peer-intf-id=0:0:0:2
ipv6-accept-peer-intf-id=1
lcp-echo-interval=20
#lcp-echo-failure=3
lcp-echo-timeout=120
unit-cache=1
#unit-preallocate=1

[auth]
#any-login=0
#noauth=0

[chap-secrets]
chap-secrets=/etc/ppp/chap-secrets
#gw-ip-address=192.168.100.1
#encrypted=0
username-hash=md5

[pppoe]
verbose=1
#ac-name=xxx
#service-name=yyy
#pado-delay=0
#pado-delay=0,100:100,200:200,-1:500
called-sid=mac
#tr101=1
#padi-limit=0
ip-pool=pool_v4
ipv6-pool=pool_v6
ipv6-pool-delegate=pool_v6_pd
ifname=pppoe%d
#sid-uppercase=0
#vlan-mon=eth0,10-200
#vlan-timeout=60
#vlan-name=%I.%N
#interface=eth1,padi-limit=1000
interface=re:enp[23]s0

[dns]
dns1=223.5.5.5
dns2=114.114.114.114

[wins]
#wins1=172.16.0.1
#wins2=172.16.1.1

[client-ip-range]
#10.0.0.0/8
#100.0.0.0/8

[ip-pool]
gw-ip-address=100.70.0.1
#vendor=Cisco
#attr=Cisco-AVPair
attr=Framed-Pool
100.70.0.10-255,name=pool_v4

[ipv6-dns]
2400:3200::1

[ipv6-dhcp]
verbose=1
pref-lifetime=604800
valid-lifetime=2592000
route-via-gw=1

[ipv6-pool]
fc00:0:1::/48,64,name=pool_v6
delegate=fc00:0:7777::/48,60,name=pool_v6_pd

[ipv6-nd]
verbose=1
#AdvManagedFlag不能设置成1，否则networkmanager拨号后会自动退出
AdvManagedFlag=0
AdvOtherConfigFlag=1
#AdvPrefixAutonomousFlag=1
#AdvLinkMTU=1
AdvAutonomousFlag=1

[log]
log-file=/var/log/accel-ppp/accel-ppp.log
log-emerg=/var/log/accel-ppp/emerg.log
log-fail-file=/var/log/accel-ppp/auth-fail.log
log-debug=/var/log/accel-ppp/debug.log
#syslog=accel-pppd,daemon
#log-tcp=127.0.0.1:3000
copy=1
#color=1
#per-user-dir=per_user
#per-session-dir=per_session
#per-session=1
level=3
```