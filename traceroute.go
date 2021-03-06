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
func NewTraceConfig(opts ...TraceOpt) (*TraceConfig, error) {
	config := &TraceConfig{
		Port:    33434,
		Hops:    64,
		Timeout: 500,
		Retries: 3,
		Size:    52,
	}
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return config, err
		}
	}
	return config, nil
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
		if size < 1 || size > 1400 {
			return fmt.Errorf("packet size must be between 1 and 1400 inclusive")
		}
		t.Size = size
		return nil
	}
}

// MakeTracerouteOptions makes a TracerouteOptions for the github.com/aeden/traceroute package
func MakeTracerouteOptions(config *TraceConfig) *tr.TracerouteOptions {
	t := &tr.TracerouteOptions{}
	t.SetMaxHops(config.Hops)
	t.SetPacketSize(config.Size)
	t.SetPort(config.Port)
	t.SetRetries(config.Retries)
	t.SetTimeoutMs(config.Timeout)
	return t
}

// TraceHop represents an individual hop in the traceroute
type TraceHop struct {
	TTL     int           `json:"ttl"`
	Host    string        `json:"host"`
	Address string        `json:"address"`
	RTT     time.Duration `json:"rtt"`
}

func traceroute(dest string, opts ...TraceOpt) ([]TraceHop, error) {
	var err error
	hops := []TraceHop{}

	config, err := NewTraceConfig(opts...)
	if err != nil {
		return hops, err
	}
	traceopts := MakeTracerouteOptions(config)
	out, err := tr.Traceroute(dest, traceopts)
	if err != nil {
		return hops, err
	}
	for _, hop := range out.Hops {
		hops = append(hops, TraceHop{
			TTL:     hop.TTL,
			Host:    hop.HostOrAddressString(),
			Address: hop.AddressString(),
			RTT:     hop.ElapsedTime,
		})
	}
	return hops, nil
}

func liveTraceroute(dest string, ch chan TraceHop, done chan bool, opts ...TraceOpt) error {
	var err error
	fmt.Println("wtf livetraceroute")
	config, err := NewTraceConfig(opts...)
	if err != nil {
		return err
	}
	traceopts := MakeTracerouteOptions(config)
	trch := make(chan tr.TracerouteHop, 0)
	go tr.Traceroute(dest, traceopts, trch)
	for {
		hop, ok := <-trch
		if !ok {
			done <- true
			return nil
		}
		ch <- TraceHop{
			TTL:     hop.TTL,
			Host:    hop.HostOrAddressString(),
			Address: hop.AddressString(),
			RTT:     hop.ElapsedTime,
		}
	}
}
