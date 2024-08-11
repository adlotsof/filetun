package cli

import (
	"bufio"
	"net"

	"github.com/adlotsof/filetun/config"
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"

	"github.com/adlotsof/filetun/types"
	"io"
	"log"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
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

func forwardPacketsToFile(iface types.Iface) {
	file, err := os.OpenFile(config.CLI.Output, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create outgoing packet file: %v", err)
	}
	defer file.Close()
	for {
		buffer := make([]byte, 1600)

		n, err := iface.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
				break
			}
		}
		if n > 0 {
			if _, err := file.Write(buffer[:n]); err != nil {
				log.Printf("Encode error listen: %v", err)
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func forwardsPacketsToIface(iface types.Iface) {
	file, err := os.Open(config.CLI.Input)
	if err != nil {
		log.Fatalf("Failed to open outgoing packet file: %v", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	packetData := make([]byte, 16000)
	for {
		n, err := reader.Read(packetData)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if err != nil {
			log.Printf("Decode error deserialize: %v", err)
		}
		if n > 0 {
			_, err = iface.Write(packetData[:n])
			if err != nil {
				log.Printf("Write error: %v", err)
				break
			}
		}
	}
}

func handleClientConnection(iface types.Iface) {
	go forwardPacketsToFile(iface)
	go forwardsPacketsToIface(iface)
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

	wg.Add(1)
	go handleClientConnection(iface)
	wg.Wait()
}
