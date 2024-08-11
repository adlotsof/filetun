package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/adlotsof/filetun/config"
	"github.com/adlotsof/filetun/types"
)

func TestSetupDevice(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test. (which also requires root privilegues)")
	}
	t.Log("testing setup of device")
	if usr, err := user.Current(); usr.Username != "root" || err != nil {
		t.Fatalf("probably you are not root while you should be %s err %v", usr, err)
	}
	config.CLI.OwnCidr = "10.0.22.0/24"
	config.CLI.PeersCidr = "10.0.21.0/24"
	config.CLI.OwnName = "testTunDevice"
	config.CLI.Input = "testInput.gob"
	config.CLI.Output = "testOutput.gob"
	defer func() {
		os.Remove(config.CLI.Input)
		os.Remove(config.CLI.Output)
		teardownDevice(config.CLI.OwnName, config.CLI.OwnCidr)
	}()
	os.Create(config.CLI.Input)
	os.Create(config.CLI.Output)
	iface, err := setupDevice(config.CLI.OwnName, config.CLI.OwnCidr, config.CLI.PeersCidr)
	if err != nil {
		t.Errorf("Error setting up device %v", err)
	}
	if iface.Name() != config.CLI.OwnName {
		t.Error("Wrong interface name")
	}
	cmd := exec.Command("ip", "a", "show", "dev", config.CLI.OwnName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("error running ip a show dev iface %v", err)
	}
	if !strings.Contains(string(output), fmt.Sprintf("inet %s brd 10.0.22.255 scope global %s", config.CLI.OwnCidr, config.CLI.OwnName)) {
		t.Fatalf("something is off with the ip a output %s", string(output))
	}
}

func TestForwardPacketsToIface(t *testing.T) {
	config.CLI.OwnCidr = "10.0.23.0/24"
	config.CLI.PeersCidr = "10.0.24.0/24"
	config.CLI.OwnName = "testTunDevice"
	config.CLI.Input = "testInput.gob"
	config.CLI.Output = "testOutput.gob"
	f, err := os.Create(config.CLI.Input)
	if err != nil {
		t.Fatalf("Error creating file %v", err)
	}
	defer func(f *os.File) {
		f.Close()
		os.Remove(config.CLI.Input)
	}(f)
	testString := "This is a test string"
	_, err = f.WriteString(testString)
	if err != nil {
		t.Fatalf("Error writing to file %v", err)
	}
	iface := types.MockIface{IfaceName: config.CLI.OwnName}
	// TODO: this is extremely unelegant, FIXME
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		go forwardsPacketsToIface(&iface)
		time.Sleep(time.Millisecond * 10)
		iface.Write([]byte("Writing more to exit"))
	}()
	wg.Wait()
	if !strings.Contains(string(iface.Content), testString) {
		t.Fatalf("not written to iface, %s", string(iface.Content))
	}
}

func TestForwardPacketsToFile(t *testing.T) {
	config.CLI.OwnCidr = "10.0.23.0/24"
	config.CLI.PeersCidr = "10.0.24.0/24"
	config.CLI.OwnName = "testTunDevice"
	config.CLI.Input = "testInput.gob"
	config.CLI.Output = "testOutput.gob"
	f, err := os.Create(config.CLI.Output)
	if err != nil {
		t.Fatalf("Error creating file %v", err)
	}
	defer func(f *os.File) {
		f.Close()
		// os.Remove(config.CLI.Output)
	}(f)
	testString := "This is a test string"

	iface := types.MockIface{IfaceName: config.CLI.OwnName}
	// TODO: this is extremely unelegant, FIXME
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		go forwardPacketsToFile(&iface)
		time.Sleep(50 * time.Millisecond)
		wg.Done()
	}()
	wg.Wait()
	f.Close()
	f, err = os.Open(config.CLI.Output)
	reader := bufio.NewReader(f)
	var fc []byte
	fc = make([]byte, 300)
	_, err = reader.Read(fc)
	if err != nil {
		t.Fatalf("error reading file %v", err)
	}
	if !strings.Contains(string(fc), testString) {
		t.Fatalf("not written to file, %s", string(fc))
	}
}
