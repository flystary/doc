nxp编译

NXP目录结构

（1）packages/apps                 //app 源码目录

（2）packages/firmware         //uboot atf rcw 源码目录

（3）packages/linux                 //linux 内核目录

（4）build/apps                         //apps编译后生成文件

（5）build/firmware                 // uboot atf rcw编译后生成文件

（6）build/linux                         // mkbootpartition命令会生成的文件，boot分区的文件

（7）build/rfs                         // rfs编译后生成文件，包含所有app的文件系统

（8）build/images                 // 合成镜像文件

编译方法：

A、自动编译命令：flex-builder -m ls1043ardb -a arm64  //全部编译，一步到位

B、分步独立编译：                                                            // -a arm64可以不写，默认为arm64

$ flex-builder -a arm64 -m ls1043ardb -i mkfw  -b sd              //下载uboot,rcw,atf编译生成 firmware_ls1043ardb_uboot_sdboot.img

$ flex-builder -a arm64 -m ls1043ardb -i mkrfs                                     //生成rootfs目录

$ flex-builder -a arm64 -m ls1043ardb -c apps                                      //下载编译apps，这步可以略去

$ flex-builder -a arm64 -m ls1043ardb -i merge-component  // merge app components and kernel modules into main userland

$ flex-builder -a arm64 -m ls1043ardb -i packrfs                          // pack target userland as rootfs_lsdk1909_LS_arm64_main.tgz,~700M

$ flex-builder -a arm64 -m ls1043ardb -i mkbootpartition      //下载linux内核编译,编译出bootpartiton和linux kernel

烧录到SD卡（nxp板子在20.04这个包不支持）：

$ cd build/images

$ flex-installer -b bootpartition_LS_arm64_lts_5.4.tgz -r rootfs_lsdk2004_ubuntu_main_arm64.tgz -f firmware_ls1043ardb_uboot_sdboot.img -d /dev/sdx  //烧录到SD卡

boot partition: The boot partition includes kernel image, DTB, distro boot script, secure boot headers, tiny initrd etc

mkbootpartition命令会生成两个文件：

​        bootpartition_LS_arm64_lts_5.4.tgz // boot分区的内容

​        linux_4.19_LS_arm64.tgz                                //包含最新的boot分区内容以及最新的kernel modules

自定义内核过程：

方式一：

flex-builder -c linux:custom -m ls1043ardb -a arm64          //自定义配置生成config文件

flex-builder -c linux -m ls1043ardb -a arm64                         //内核编译

flex-builder -i mkbootpartition -m ls1043ardb -a arm64  //生成内核文件

方式二：

新建文件packages/linux/linux/arch/arm64/configs/svxnetworks.config，然后把要新增的模块，加到这个配置文件里

builder -c linux -B fragment:svxnetworks.config -a arm64 -m ls1043ardb -b sd

登录板子上下载内核文件

root@localhost:/# tar xfmv bootpartition_LS_arm64_lts_5.4.tgz -C /boot

root@localhost:/# reboot