package fileBackend
import (
	"github.com/adlotsof/filetun/types"
	"github.com/adlotsof/filetun/config"
	"testing"
	"os"
	"time"
	"bufio"
	"strings"
	"sync"
)


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
	backend := FileBackend{}
	// TODO: this is extremely unelegant, FIXME
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		go backend.ReceiveFromBackend(&iface)
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
		os.Remove(config.CLI.Output)
	}(f)
	testString := "This is a test string"

	iface := types.MockIface{IfaceName: config.CLI.OwnName}
	backend := FileBackend{}
	// TODO: this is extremely unelegant, FIXME
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		go backend.SendToBackend(&iface)
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

func TestSetup(t *testing.T) {
	config.CLI.OwnCidr = "10.0.23.0/24"
	config.CLI.PeersCidr = "10.0.24.0/24"
	config.CLI.OwnName = "testTunDevice"
	config.CLI.Input = "testInput.gob"
	config.CLI.Output = "testOutput.gob"
	defer func() {
		os.Remove(config.CLI.Input)
		os.Remove(config.CLI.Output)
	}()
	backend := FileBackend{}
	err := backend.Setup()
	if err != nil {
		t.Fatalf("error setting up backend %v", err)
	}
	// check if the input / output file exist
	if _, err := os.Stat(config.CLI.Input); os.IsNotExist(err) {
		t.Fatalf("input file does not exist %v", err)
	}
	if _, err := os.Stat(config.CLI.Output); os.IsNotExist(err) {
		t.Fatalf("output file does not exist %v", err)
	}
}
