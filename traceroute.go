package main

import (
	"fmt"
	"time"

	tr "github.com/aeden/traceroute"
)

// TraceConfig contains configuration details for a traceroute
type TraceConfig struct {
	Port    int
	Hops    int
	Timeout int
	Retries int
	Size    int
}

// NewTraceConfig creates a new TraceConfig with default values
func NewTraceConfig() *TraceConfig {
	config := &TraceConfig{
		Port:    33434,
		Hops:    64,
		Timeout: 500,
		Retries: 3,
		Size:    52,
	}
	return config
}

// TraceOpt is passed into traceroute to set optional configuration
type TraceOpt func(*TraceConfig) error

// PortOpt sets the port option on the config
func PortOpt(port int) func(*TraceConfig) error {
	return func(t *TraceConfig) error {
		if port < 1 || port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535 inclusive")
		}
		t.Port = port
		return nil
	}
}

// HopsOpt sets the maximum number of hops in the config
func HopsOpt(hops int) func(*TraceConfig) error {
	return func(t *TraceConfig) error {
		if hops < 1 || hops > 255 {
			return fmt.Errorf("hops must be between 1 and 255 inclusive")
		}
		t.Hops = hops
		return nil
	}
}

// TimeoutOpt sets the timeout in the config
func TimeoutOpt(timeout int) func(*TraceConfig) error {
	return func(t *TraceConfig) error {
		if timeout < 1 || timeout > 5000 {
			return fmt.Errorf("timeout must be between 1 and 5000 inclusive")
		}
		t.Timeout = timeout
		return nil
	}
}

// RetriesOpt sets the maximum number of retries in the config
func RetriesOpt(retries int) func(*TraceConfig) error {
	return func(t *TraceConfig) error {
		if retries < 1 || retries > 5 {
			return fmt.Errorf("retries must be between 1 and 5 inclusive")
		}
		t.Retries = retries
		return nil
	}
}

// SizeOpt sets the maximum packet size in the config
func SizeOpt(size int) func(*TraceConfig) error {
	return func(t *TraceConfig) error {
		fmt.Println("Processing size opt", size)
		if size < 1 || size > 1400 {
			return fmt.Errorf("packet size must be between 1 and 1400 inclusive")
		}
		t.Size = size
		return nil
	}
}

// Hop represents an individual hop in the traceroute
type Hop struct {
	TTL     int           `json:"ttl"`
	Host    string        `json:"host"`
	Address string        `json:"address"`
	RTT     time.Duration `json:"rtt"`
}

func traceroute(dest string, opts ...TraceOpt) ([]Hop, error) {
	var err error
	hops := []Hop{}

	config := NewTraceConfig()
	for _, opt := range opts {
		fmt.Println("Processing opt")
		if err = opt(config); err != nil {
			return hops, err
		}
	}

	traceopts := &tr.TracerouteOptions{}
	traceopts.SetMaxHops(config.Hops)
	traceopts.SetPacketSize(config.Size)
	traceopts.SetPort(config.Port)
	traceopts.SetRetries(config.Retries)
	traceopts.SetTimeoutMs(config.Timeout)
	fmt.Printf("%+v\n", traceopts)
	out, err := tr.Traceroute(dest, traceopts)
	if err != nil {
		return hops, err
	}
	for _, hop := range out.Hops {
		hops = append(hops, Hop{
			TTL:     hop.TTL,
			Host:    hop.HostOrAddressString(),
			Address: hop.AddressString(),
			RTT:     hop.ElapsedTime,
		})
	}
	return hops, nil
}
