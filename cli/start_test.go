package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"testing"

	"github.com/adlotsof/filetun/config"
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
	defer func() {
		os.Remove(config.CLI.Input)
		os.Remove(config.CLI.Output)
	}()
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
