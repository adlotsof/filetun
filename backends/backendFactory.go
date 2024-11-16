package backends

import (
	"github.com/adlotsof/filetun/types"
	"github.com/adlotsof/filetun/backends/fileBackend"
)

func BackendFactory (backendType string) types.Backend {
	var backend types.Backend
	switch backendType {
	case "file":
		backend = &fileBackend.FileBackend{}
	default:
		panic("unknown backend")
	}
	backend.Setup()
	return backend
}
