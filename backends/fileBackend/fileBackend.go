package fileBackend

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/adlotsof/filetun/config"
	"github.com/adlotsof/filetun/types"
)

type FileBackend struct {
	// contains filtered or unexported fields
}

func (f *FileBackend) SendToBackend(iface types.Iface) error {
	file, err := os.OpenFile(config.CLI.Output, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to create outgoing packet file: %v", err)
		return err
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
	return nil
}

func (f *FileBackend) ReceiveFromBackend(iface types.Iface) error {
	file, err := os.Open(config.CLI.Input)
	if err != nil {
		log.Fatalf("Failed to open outgoing packet file: %v", err)
		return err
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
	return nil
}

func (f *FileBackend) Setup() error {
	// create the files if they do not exist
	if _, err := os.Stat(config.CLI.Input); os.IsNotExist(err) {
		os.Create(config.CLI.Input)
	} else {
		log.Printf("input file already exists")
	}
	if _, err := os.Stat(config.CLI.Output); os.IsNotExist(err) {
		os.Create(config.CLI.Output)
	} else {
		log.Printf("output file already exists")
	}
  return nil
}
