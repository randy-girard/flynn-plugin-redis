package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/randy-girard/flynn/discoverd/client"
)

// Ensure the program can register with discoverd.
func TestMain_Discoverd(t *testing.T) {
	if _, err := exec.LookPath("redis-server"); err != nil {
		t.Skip("redis-server not found in PATH")
	}
	m := newTestMain()
	defer m.Close()

	// Mock heartbeater.
	var hbClosed bool
	hb := NewHeartbeater("127.0.0.1:0")
	hb.CloseFn = func() error { hbClosed = true; return nil }

	lnPick, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	redisPort := strconv.Itoa(lnPick.Addr().(*net.TCPAddr).Port)
	if err := lnPick.Close(); err != nil {
		t.Fatal(err)
	}

	// Validate arguments passed to discoverd.
	m.DiscoverdClient.AddServiceFn = func(name string, config *discoverd.ServiceConfig) error {
		if name != "redis" {
			t.Fatalf("unexpected service name: %s", name)
		}
		return nil
	}
	m.DiscoverdClient.RegisterInstanceFn = func(service string, inst *discoverd.Instance) (discoverd.Heartbeater, error) {
		if service != "redis" {
			t.Fatalf("unexpected service: %s", service)
		} else if !reflect.DeepEqual(inst, &discoverd.Instance{
			Addr: ":" + redisPort,
			Meta: map[string]string{"REDIS_ID": m.Process.ID},
		}) {
			t.Fatalf("unexpected inst: %#v", inst)
		}
		return hb, nil
	}

	// set a password and bind Redis server to an ephemeral port (avoid clashes with a host redis on 6379)
	m.Process.Password = "test"
	m.Process.Port = redisPort

	// Execute program.
	if err := m.Run(); err != nil {
		t.Fatal(err)
	}

	// Close program and validate that the heartbeater was closed.
	if err := m.Close(); err != nil {
		t.Fatal(err)
	} else if !hbClosed {
		t.Fatal("expected heartbeater close")
	}
}

// testMain is a test wrapper for Main.
type testMain struct {
	*Main
	DiscoverdClient *DiscoverdClient

	Stdin  bytes.Buffer
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func newTestMain() *testMain {
	binPath, err := exec.LookPath("redis-server")
	if err != nil {
		panic("redis-server not found in PATH (caller should skip)")
	}

	// Create a temporary data directory.
	dataDir, err := ioutil.TempDir("", "flynn-redis-")
	if err != nil {
		panic(err)
	}

	// Create test wrapper with random port and temporary data directory.
	m := &testMain{
		Main:            NewMain(),
		DiscoverdClient: NewDiscoverdClient(),
	}
	m.Main.Addr = "127.0.0.1:0"
	m.Main.DataDir = dataDir
	m.Main.Process.BinDir = filepath.Dir(binPath)
	m.Main.DiscoverdClient = m.DiscoverdClient

	m.Main.Stdin = &m.Stdin
	m.Main.Stdout = &m.Stdout
	m.Main.Stderr = &m.Stderr

	if testing.Verbose() {
		m.Main.Stdout = io.MultiWriter(os.Stdout, m.Main.Stdout)
		m.Main.Stderr = io.MultiWriter(os.Stderr, m.Main.Stderr)
	}

	return m
}

// Close cleans up temporary paths and closes the program.
func (m *testMain) Close() error {
	defer os.RemoveAll(m.DataDir)
	return m.Main.Close()
}

// DiscoverdClient is a mock implementation of Main.DiscoverdClient.
type DiscoverdClient struct {
	AddServiceFn       func(name string, config *discoverd.ServiceConfig) error
	RegisterInstanceFn func(service string, inst *discoverd.Instance) (discoverd.Heartbeater, error)
}

// NewDiscoverdClient returns a new instance of DiscoverdClient with default mock implementations.
func NewDiscoverdClient() *DiscoverdClient {
	return &DiscoverdClient{
		AddServiceFn: func(name string, config *discoverd.ServiceConfig) error { return nil },
		RegisterInstanceFn: func(service string, inst *discoverd.Instance) (discoverd.Heartbeater, error) {
			return NewHeartbeater(inst.Addr), nil
		},
	}
}

func (c *DiscoverdClient) AddService(name string, config *discoverd.ServiceConfig) error {
	return c.AddServiceFn(name, config)
}
func (c *DiscoverdClient) RegisterInstance(service string, inst *discoverd.Instance) (discoverd.Heartbeater, error) {
	return c.RegisterInstanceFn(service, inst)
}

// Heartbeater is a mock implementation of discoverd.Heartbeater.
type Heartbeater struct {
	SetMetaFn   func(map[string]string) error
	CloseFn     func() error
	AddrFn      func() string
	SetClientFn func(*discoverd.Client)
}

// NewHeartbeater returns a new instance of Heartbeater with default mock implementations.
func NewHeartbeater(addr string) *Heartbeater {
	return &Heartbeater{
		SetMetaFn:   func(map[string]string) error { return nil },
		CloseFn:     func() error { return nil },
		AddrFn:      func() string { return addr },
		SetClientFn: func(*discoverd.Client) {},
	}
}

func (h *Heartbeater) SetMeta(m map[string]string) error { return h.SetMetaFn(m) }
func (h *Heartbeater) Close() error                      { return h.CloseFn() }
func (h *Heartbeater) Addr() string                      { return h.AddrFn() }
func (h *Heartbeater) SetClient(c *discoverd.Client)     { h.SetClientFn(c) }
