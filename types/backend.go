package types


// Backend is the interface that wraps the basic methods for the backend.
type Backend interface {
	// SendToBackend sends the data from the tun interface to the backend.
	SendToBackend(iface Iface) error
	// ReceiveFromBackend receives the data from the backend, forwards data to interface
	ReceiveFromBackend(iface Iface) error
	// Setup sets up the backend
	Setup() error
}

type BackendFactory func(backendType string) Backend
