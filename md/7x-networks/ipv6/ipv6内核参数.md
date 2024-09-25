net.ipv6.conf.all.disable_ipv6 是否禁用ipv6

 0：不禁用

 1：禁用

net.ipv6.conf.all.forwarding 所有网络接口开启ipv6转发

 0：关闭

 1：开启

net.ipv6.conf.all.accept_dad 0：取消DAD功能

 1：启用DAD功能，但link-local地址冲突时，不关闭ipv6功能

 2：启用DAD功能，但link-local地址冲突时，关闭ipv6功能

net.ipv6.conf.all.accept_ra 接受IPv6路由通告.并且根据得到的信息自动设定.

 0：不接受路由通告

 1：当forwarding禁止时接受路由通告

 2：任何情况下都接受路由通告

net.ipv6.conf.all.accept_ra_defrtr 是否接受ipv6路由器发出的默认路由设置

 0：不接受

 1：接受

net.ipv6.conf.all.accept_ra_pinfo 当accept_ra开启时此选项会自动开启，关闭时则会关闭

net.ipv6.conf.all.accept_ra_rt_info_max_plen 在路由通告中路由信息前缀的最大长度。当

net.ipv6.conf.all.accept_ra_rtr_pref

net.ipv6.conf.all.accept_redirects 是否接受ICMPv6重定向包

 0：拒绝接受ICMPv6，当forwarding=1时，此值会自动设置为0

 1：启动接受ICMPv6，当forwarding=0时，此值会自动设置为1

net.ipv6.conf.all.accept_source_route 接收带有SRR选项的数据报。主机设为0，路由设为1

net.ipv6.conf.all.autoconf 设定本地连结地址使用L2硬件地址. 它依据界面的L2-MAC address自动产生一个地址如:"fe80::201:23ff:fe45:6789"

net.ipv6.conf.all.dad_transmits 接口增加ipv6地址时，发送几次DAD包

net.ipv6.conf.all.force_mld_version

net.ipv6.conf.all.force_tllao

net.ipv6.conf.all.hop_limit 缺省hop限制

net.ipv6.conf.all.max_addresses 所有网络接口自动配置IP地址的数量最大值

 0：不限制

 \>0：最大值

net.ipv6.conf.all.max_desync_factor DESYNC_FACTOR的最大值，DESYNC_FACTOR是一个随机数，用于防止客户机在同一时间生成新的地址

net.ipv6.conf.all.mc_forwarding 是否使用多路广播进行路由选择，需要内核编译时开启了CONFIG_MROUTE选项并且开启了多路广播路由选择的后台daemon

 0：关闭

 1：开启

net.ipv6.conf.all.mldv1_unsolicited_report_interval 每次发送MLDv1的主动报告的时间间隔(ms)

net.ipv6.conf.all.mldv2_unsolicited_report_interval 每次发送MLDv2的主动报告的时间间隔(ms)

net.ipv6.conf.all.mtu ipv6的最大传输单元

net.ipv6.conf.all.ndisc_notify 如何向邻居设备通知地址和设备的改变

 0：不通知

 1：主动向邻居发送广播报告硬件地址或者设备发生了改变

net.ipv6.conf.all.optimistic_dad 是否启用optimistic DAD(乐观地进行重复地址检查)

 0：关闭

 1：开启

net.ipv6.conf.all.proxy_ndp 此功能类似于ipv4的nat，可将内网的包转发到外网，外网不能主动发给内网。

 0：关闭

 1：开启

net.ipv6.conf.all.regen_max_retry 尝试生成临时地址的次数

net.ipv6.conf.all.router_probe_interval 路由器探测间隔(秒)

net.ipv6.conf.all.router_solicitation_delay 在发送路由请求之前的等待时间(秒).

net.ipv6.conf.all.router_solicitation_interval 在每个路由请求之间的等待时间(秒).

net.ipv6.conf.all.router_solicitations 假定没有路由的情况下发送的请求个数

net.ipv6.conf.all.temp_prefered_lft

net.ipv6.conf.all.temp_valid_lft

net.ipv6.conf.all.use_tempaddr