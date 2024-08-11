package types

import "fmt"

type Iface interface {
	Close() error
	//  IsTAP() bool
	//  IsTUN() bool
	Name() string
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

type MockIface struct {
	IfaceName string
	Content   []byte
	Reads    int
}

func (m *MockIface) Close() error {
	return nil
}

func (m *MockIface) Name() string {
	return m.IfaceName
}

func (m *MockIface) Read(p []byte) (n int, err error) {
	if m.Reads > 1 {
		return 1, fmt.Errorf("erroring out on second read")
	}
	m.Reads++
	n = copy(p, []byte(" This is a test string "))
	return n, nil
}

func (m *MockIface) Write(p []byte) (n int, err error) {
	if len(m.Content) > len("This is a test string") {
		return 1, fmt.Errorf("Erroring out on second write")
	}
	m.Content = append(m.Content, p...)
	return len(m.Content), nil
}
