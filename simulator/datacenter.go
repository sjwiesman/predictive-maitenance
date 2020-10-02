package main

import (
	"fmt"
	entities "generator/protobuf"
	"log"
	"math/rand"
	"time"
)

type datacenter struct {
	servers map[string]*server

	fixRequests chan string

	reports chan *entities.ServerMetricReport

	// Pending requests for fixes
	// that are waiting on parts
	pendingFixes []string

	// The number of available parts
	parts int32
}

func newDataCenter(numServers int) datacenter {
	log.Printf("starting %d servers", numServers)
	dc := datacenter{
		servers:      make(map[string]*server, numServers),
		fixRequests:  make(chan string),
		reports:      make(chan *entities.ServerMetricReport),
		pendingFixes: make([]string, numServers),
		parts:        int32(numServers),
	}
	for serverId := 0; serverId < numServers; serverId++ {
		fmt.Printf("creating server %d", serverId)
		server := &server{
			serverId:   fmt.Sprintf("server-%d", serverId),
			rackId:     fmt.Sprintf("rack-%d", len(dc.servers)/40),
			health:     healthy,
			resiliency: rand.Intn(100) + 1,
			counter:    0,
		}
		dc.servers[server.serverId] = server
	}

	return dc
}

func (d *datacenter) run() {
	fmt.Printf("starting simulation")
	runServers := time.NewTicker(5 * time.Second)
	newParts := time.NewTicker(1 * time.Second)

	capacity := d.parts

	for {
		select {
		case <-runServers.C:
			for _, s := range d.servers {
				report := s.tick()
				if report != nil {
					d.reports <- report
				}
			}
		case <-newParts.C:
			if d.parts < capacity {
				d.parts++
			}
		case id := <-d.fixRequests:
			s, ok := d.servers[id]
			if !ok {
				log.Printf("Unknown server %s", s)
				continue
			}

			if d.parts > 0 {
				d.parts--
				s.fix()
			} else {
				d.pendingFixes = append(d.pendingFixes, id)
				log.Printf("Request to fix server %s but no parts available", id)
			}
		}
	}
}
