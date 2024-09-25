package tundrive

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"tunnel/internal/config"

	"golang.org/x/sys/unix"
)

const tun_drive_path = "/dev/net/tun"

const TUN_TYPE_TAP = unix.IFF_TAP
const TUN_TYPE_TUN = unix.IFF_TUN

type TUNInfo struct {
	FD      int
	Type    uint16
	Name    string
	Mtu     uint16
	LocalIP net.IP
	Mask    net.IP
	Flags   uint16
}

func NewTunInfo(c config.Config) (r TUNInfo) {
	return TUNInfo{
		Name: c.GetString("tun_name", ""),
		Mtu:  uint16(c.GetInt64("mtu", 1430)),
	}
}

func (t *TUNInfo) Check() error {
	if len(t.Name) == 0 {
		return fmt.Errorf("tun name is empty")
	}
	if t.Mtu <= 0 {
		return fmt.Errorf("tun mtu < 0")
	}
	if t.LocalIP.IsUnspecified() {
		return fmt.Errorf("tun local IP is unspecified")
	}
	// if t.PeerIP.IsUnspecified() {
	// 	return fmt.Errorf("tun peerip is unspecified")
	// }
	if t.Mask.IsUnspecified() {
		return fmt.Errorf("tun mask is unspecified")
	}
	return nil
}

func (t TUNInfo) CreateTunDrive() (err error) {
	var ifreq *unix.Ifreq
	ifreq, err = unix.NewIfreq(t.Name)
	var addr [4]byte
	var socketfd int

	t.FD, err = syscall.Open(tun_drive_path, os.O_RDWR, 0)
	if err != nil {
		return
	}

	/* Flags: IFF_TUN   - TUN device (no Ethernet headers)
	 *        IFF_TAP   - TAP device
	 *
	 *        IFF_NO_PI - Do not provide packet information
	 *        IFF_MULTI_QUEUE - Create a queue of multiqueue device
	 */
	ifreq.SetUint16(t.Type)
	err = unix.IoctlIfreq(t.FD, unix.TUNSETIFF, ifreq)
	if err != nil {
		return err
	}

	// 获取设置的socketfd
	if socketfd, err = unix.Socket(
		unix.AF_INET,
		unix.SOCK_DGRAM,
		unix.IPPROTO_IP,
	); err != nil {
		return
	}
	defer unix.Close(socketfd)

	// 设置IP地址
	copy(addr[:], t.LocalIP.To4())
	ifreq.SetInet4Addr(addr[:])
	err = unix.IoctlIfreq(socketfd, unix.SIOCSIFADDR, ifreq)
	if err != nil {
		return
	}

	// 设置子网掩码
	copy(addr[:], t.Mask.To4())
	ifreq.SetInet4Addr(addr[:])
	err = unix.IoctlIfreq(socketfd, unix.SIOCSIFNETMASK, ifreq)
	if err != nil {
		return
	}

	// 设置MTU
	ifreq.SetUint32(uint32(t.Mtu))
	err = unix.IoctlIfreq(socketfd, unix.SIOCSIFMTU, ifreq)
	if err != nil {
		return
	}

	/*  persistent status */
	err = unix.IoctlSetPointerInt(t.FD, syscall.TUNSETPERSIST, 1)
	if err != nil {
		return
	}

	unix.SetNonblock(t.FD, true)

	return
}

func (t TUNInfo) RunTunDrive() (err error) {
	var socketfd int
	var ifreq *unix.Ifreq
	// 获取设置的socketfd
	if socketfd, err = unix.Socket(
		unix.AF_INET,
		unix.SOCK_DGRAM,
		unix.IPPROTO_IP,
	); err != nil {
		return
	}
	defer unix.Close(socketfd)

	if ifreq, err = unix.NewIfreq(t.Name); err != nil {
		return
	}

	err = unix.IoctlIfreq(socketfd, unix.SIOCGIFFLAGS, ifreq)
	ifreq.SetUint32(ifreq.Uint32() | unix.IFF_UP)
	err = unix.IoctlIfreq(socketfd, unix.SIOCSIFFLAGS, ifreq)
	return
}
