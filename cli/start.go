package cli

import (
	"net"

	"github.com/adlotsof/filetun/config"
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"

	"github.com/adlotsof/filetun/types"
	"github.com/adlotsof/filetun/backends"
	"log"
	"errors"
	"fmt"
	"sync"
)

func setupDevice(ifName string, ifCidr string, peersSubnet string) (*water.Interface, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = ifName
	iface, err := water.New(config)
	if err != nil {
		return nil, errors.Join(errors.New("unable to create iface"), err)
	}
	log.Printf("Interface created %s", iface.Name())
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return nil, errors.Join(errors.New("iface not found"), err)
	}
	addr, err := netlink.ParseAddr(ifCidr)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("could not parse ip %s", ifCidr), err)
	}
	if err := netlink.AddrAdd(link, addr); err != nil {
		log.Fatalf("Could not add address to link device %v", err)
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return nil, errors.Join(fmt.Errorf("could not bring device %s up, ", ifName), err)
	}
	_, dst, err := net.ParseCIDR(peersSubnet)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Error parsing peers cidr"), err)
	}
	viaIP, _, err := net.ParseCIDR(ifCidr)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error parsing own cidr"), err)
	}

	route := netlink.Route{
		Dst: dst,
		Gw:  viaIP,
	}
	if err := netlink.RouteAdd(&route); err != nil {
		return nil, errors.Join(fmt.Errorf("error adding route"), err)
	}
	log.Printf("%s setup with cidr %s\n", ifName, ifCidr)
	return iface, nil
}

func teardownDevice(ifName string, ifCidr string) {
	link, err := netlink.LinkByName(ifName)
	if err != nil {
		log.Fatalf("Could not find %s for teardown: %v", ifName, err)
		return
	}

	// Remove the IP address from the interface
	addr, err := netlink.ParseAddr(ifCidr)
	if err != nil {
		log.Fatalf("Could not parse IP %s for teardown: %v", ifCidr, err)
		return
	}
	if err := netlink.AddrDel(link, addr); err != nil {
		log.Printf("Could not remove address from device %s: %v", ifName, err)
		// Not fatal at this point, attempt to bring down the interface anyway
	}

	// Bring the interface down
	if err := netlink.LinkSetDown(link); err != nil {
		log.Printf("Could not bring device %s down: %v", ifName, err)
	}

	log.Printf("%s teardown complete", ifName)
}

func handleClientConnection(iface types.Iface, backend types.Backend) {
	// go forwardPacketsToFile(iface)
	go backend.SendToBackend(iface)
    go backend.ReceiveFromBackend(iface)
	// go forwardsPacketsToIface(iface)
}

func Run() {
	var wg sync.WaitGroup
	conf := &config.CLI
	defer func() {
		teardownDevice(conf.OwnName, conf.OwnCidr)

	}()

	iface, err := setupDevice(conf.OwnName, conf.OwnCidr, conf.PeersCidr)
	if err != nil {
		log.Fatalf("coundt setup device, %v", err)
	}
	backend := backends.BackendFactory(conf.BackendType)

	wg.Add(1)
	go handleClientConnection(iface, backend)
	wg.Wait()
}
