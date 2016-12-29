package main

import (
	"time"

	tr "github.com/aeden/traceroute"
)

// Hop represents an individual hop in the traceroute
type Hop struct {
	TTL     int           `json:"ttl"`
	Host    string        `json:"host"`
	Address string        `json:"address"`
	RTT     time.Duration `json:"rtt"`
}

func traceroute(dest string) ([]Hop, error) {
	hops := []Hop{}

	out, err := tr.Traceroute(dest, &tr.TracerouteOptions{})
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
