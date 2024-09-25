## 完整规则库：

https://rules.emergingthreats.net/open/suricata-5.0/rules/

### 规则库说明：

1. botcc.portgrouped.rules botcc.rules

这些是已知和确认的活动僵尸网络和其C&C(command and control)服务器。由一些组织生成，每日更新。

1. ciarmy.rules

封锁被ciArmy.com标记出来的Top Attackers

1. compromised.rules

这是一个已知的受影响的主机列表，每天更新。

1. drop.rules

每天更新的Spamhaus DROP(Don't Route or Peer)列表。列出了著名的、专业的垃圾邮件发送者。

1. dshield.rules

每天更新的DShield top attackers

1. emerging-activex.rules

主要用来检测与ActiveX控件有关的攻击

1. emerging-attack_response.rules

这些规则是为了捕获成功攻击的结果，诸如“id=root”之类的东西，或者表示可能发生妥协的错误消息（即虽然产生了错误消息，但是攻击已经成功）。

1. emerging-chat.rules

主要检测聊天软件、即时通讯软件的攻击，大部分是国外的一些软件，比如facebook，雅虎，msn

1. emerging-current_events.rules

这些规则是不打算在规则集中长期保存的，或者是在被包含之前进行测试。大多数情况下，这些都是针对当天的大量二进制URL的简单sigs，用来捕获CLSID新发现的易受攻击的应用程序，我们没有这些漏洞的任何细节。这些sigs很有用，却不是长期有效的。

1. emerging-dns.rules

检测dns协议相关的攻击

1. emerging-dos.rules

目的是捕获入站的DOS（拒绝服务）活动和出站指示。

1. emerging-ftp.rules

检测ftp协议相关的攻击

1. emerging-games.rules

魔兽世界、星际争霸和其他流行的在线游戏都在这里。我们不打算把这些东西贴上邪恶的标签，只是它们不适合所有的攻击环境，所以将它们放在了这里。

1. emerging-icmp_info.rules emerging-icmp.rules

检测与icmp协议相关的攻击

1. emerging-exploit.rules

直接检测exploits（漏洞）的规则。如果你在寻找windows的漏洞等，他们会在这里被列出。就像sql注入一样，exploits有着自己的体系。总之就是用来检测exploits漏洞的。

1. emerging-imap.rules

检测与imap相关的攻击

1. emerging-inappropriate.rules

色情、儿童色情，你不应该在工作中访问的网站等等。WARNING：这些都大量使用了正则表达式，因此存在高负荷和频繁的误报问题。只有当你真正对这些规则感兴趣时才去运行这些规则。

1. emerging-info.rules

检测与信息泄露、信息盗取等事件的规则，里面会检测后门、特洛伊木马等与info相关的攻击

1. emerging-malware.rules

这一套最初只是间谍软件。间谍软件和恶意软件之间的界限已经很模糊了

1. emerging-misc.rules

检测混杂的攻击，这种攻击一般没有确切的分类，或者使用了多种技术

1. emerging-mobile_malware.rules

检测移动设备上的恶意软件

1. emerging-netbios.rules

检测与netbios有关的攻击

1. emerging-p2p.rules

P2P(Peer to Peer)之类的。我们并不想将它定义为有害的，只是不适合出现在IPS/IDS的网络环境中。

1. emerging-policy.rules

对于经常被公司或组织政策禁止的事务的规则。Myspace、Ebay之类的东西。

1. emerging-pop3.rules

检测与pop3协议有关的攻击

1. emerging-rpc.rules

检测与rpc（远程过程调用协议)有关的攻击

1. emerging-scada.rules

检测与SCADA（数据采集与监控系统）相关的攻击

1. emerging-scan.rules==

检测探测行为。Nessus，Nikto，端口扫描等这样的活动。这是攻击前准备时期的警告。

1. emerging-shellcode.rules

检测shellcode，shellcode是一段用于利用软件漏洞而执行的代码，以其经常让攻击者获得shell而得名。

1. emerging-smtp.rules

检测与smtp协议相关的攻击

1. emerging-snmp.rules

检测与snmp协议相关的攻击

1. emerging-sql.rules

这是一个巨大的规则集，用于捕获在特殊应用程序上的特殊攻击。这里面有一些普遍的SQL注入攻击规则，效果很好，可以捕获大多数攻击。

但是这些规则根据不同的app和不同的web服务器，有很大的差别。如果你需要运行非常严格的web服务或者很重视信息的安全性，请使用这个规则集。

1. emerging-telnet.rules

检测与telnet协议相关的攻击

1. emerging-tftp.rules

检测与tftp协议相关的攻击

1. emerging-user_agents.rules

检测异常的user-agents

1. emerging-voip.rules

检测voip相关的异常

1. emerging-web_client.rules

检测web客户端的攻击

1. emerging-web_server.rules

检测web服务端的攻击

1. emerging-web_specific_apps.rules

检测特殊的web应用程序的异常

1. emerging-worm.rules

检测蠕虫

1. tor.rules

检测使用tor进行匿名通信的流量，tor本身没有威胁，但却是很可以的行为