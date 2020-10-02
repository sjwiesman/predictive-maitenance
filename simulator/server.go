package main

import (
	entities "generator/protobuf"
	"log"
	"math/rand"
	"time"
)

type health int

const (
	healthy health = iota
	failing
	failed
)

func (h health) String() string {
	if h == healthy {
		return "healthy"
	}

	if h == failing {
		return "failing"
	}

	return "failed"
}

type server struct {
	// The id of the server
	serverId string

	// The id of the rack
	rackId string

	// The current reports of the machine
	health health

	// The number of turns before the machines
	// health can change
	resiliency int

	// The numbers of turns since the
	// last health change
	counter int
}

func (s *server) tick() *entities.ServerMetricReport {
	if s.health == failed {
		return nil
	}

	s.counter++
	if s.counter == s.resiliency {
		s.counter = 0
		s.health++
	}

	low := 50
	high := 75
	if s.health == failing {
		low = 76
		high = 100
	}

	return &entities.ServerMetricReport{
		ServerId:                                     s.serverId,
		RackId:                                       s.rackId,
		CpuPercentNanosecondsIdle:                    float64(rand.Intn(high-low) + low),
		CpuPercentNanosecondsNice:                    float64(rand.Intn(high-low) + low),
		CpuPercentNanosecondsInterrupt:               float64(rand.Intn(high-low) + low),
		CpuPercentNanosecondsSoftirq:                 float64(rand.Intn(high-low) + low),
		CpuPercentNanosecondsSteal:                   float64(rand.Intn(high-low) + low),
		CpuPercentNanosecondsSystem:                  float64(rand.Intn(high-low) + low),
		CpuPercentNanosecondsUser:                    float64(rand.Intn(high-low) + low),
		CpuPercentNanosecondsWait:                    float64(rand.Intn(high-low) + low),
		InterfaceIfDroppedIn:                         float64(rand.Intn(high-low) + low),
		InterfaceIfErrorsIn:                          float64(rand.Intn(high-low) + low),
		InterfaceIfOctetsIn:                          float64(rand.Intn(high-low) + low),
		InterfaceIfPacketsIn:                         float64(rand.Intn(high-low) + low),
		InterfaceIfDroppedOut:                        float64(rand.Intn(high-low) + low),
		InterfaceIfErrorsOut:                         float64(rand.Intn(high-low) + low),
		InterfaceIfOctetsOut:                         float64(rand.Intn(high-low) + low),
		InterfaceIfPacketsOut:                        float64(rand.Intn(high-low) + low),
		MemoryMemoryBuffered:                         float64(rand.Intn(high-low) + low),
		MemoryMemoryCached:                           float64(rand.Intn(high-low) + low),
		MemoryMemoryFree:                             float64(rand.Intn(high-low) + low),
		MemoryMemorySlabRecl:                         float64(rand.Intn(high-low) + low),
		MemoryMemorySlabUnrecl:                       float64(rand.Intn(high-low) + low),
		MemoryMemoryTotal:                            float64(rand.Intn(high-low) + low),
		MemoryMemoryUsed:                             float64(rand.Intn(high-low) + low),
		DiskDiskIoTimeIoTime:                         float64(rand.Intn(high-low) + low),
		DiskDiskIoTimeWeightedIoTime:                 float64(rand.Intn(high-low) + low),
		DiskDiskMergedRead:                           float64(rand.Intn(high-low) + low),
		DiskDiskMergedWrite:                          float64(rand.Intn(high-low) + low),
		DiskDiskOctectsRead:                          float64(rand.Intn(high-low) + low),
		DiskDiskOctectsWrite:                         float64(rand.Intn(high-low) + low),
		DiskDiskOpsRead:                              float64(rand.Intn(high-low) + low),
		DiskDiskOpsWrite:                             float64(rand.Intn(high-low) + low),
		DiskDiskTimeRead:                             float64(rand.Intn(high-low) + low),
		DiskDiskTimeWrite:                            float64(rand.Intn(high-low) + low),
		ProcessesPsStateBlocked:                      float64(rand.Intn(high-low) + low),
		ProcessesPsStatePaging:                       float64(rand.Intn(high-low) + low),
		ProcessesPsStateRunning:                      float64(rand.Intn(high-low) + low),
		ProcessesPsStateSleeping:                     float64(rand.Intn(high-low) + low),
		ProcessesPsStateStopped:                      float64(rand.Intn(high-low) + low),
		ProcessesPsStateZombies:                      float64(rand.Intn(high-low) + low),
		MemoryAndCacheMemoryBandwidthLocal:           float64(rand.Intn(high-low) + low),
		MemoryAndCacheMemoryBandwidthRemote:          float64(rand.Intn(high-low) + low),
		MemoryAndCacheBytesLlc:                       float64(rand.Intn(high-low) + low),
		MemoryExceptionsErrorsCorrectedMemoryErrors:  float64(rand.Intn(high-low) + low),
		MemoryExceptionsErrorsUncorrectedMemoryError: float64(rand.Intn(high-low) + low),
		MemoryThrottlingSensor:                       float64(rand.Intn(high-low) + low),
		IpmiVoltage:                                  float64(rand.Intn(high-low) + low),
		IpmiTemperatureAgg:                           float64(rand.Intn(high-low) + low),
		IpmiAirflow:                                  float64(rand.Intn(high-low) + low),
		IpmiWattsInputPower:                          float64(rand.Intn(high-low) + low),
	}
}

func (s *server) fix() {
	if s.health == healthy {
		log.Printf("Wasted part on healthy server %s", s.serverId)
	} else {
		log.Printf("Fixing server %s", s.serverId)
	}

	s.counter = 0
	s.health = healthy
}

func (s *server) run(status chan<- *entities.ServerMetricReport, fix <-chan bool) {
	ticker := time.NewTicker(5 * time.Second)

	jitter := time.Duration(rand.Intn(5000)) * time.Millisecond
	time.Sleep(jitter)

	for {
		select {
		case <-ticker.C:
			report := s.tick()
			if report != nil {
				status <- report
			}
		case <-fix:
			s.fix()
		}
	}
}
