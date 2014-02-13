package memstats

import (
	"fmt"
	"runtime"
	"time"
)

type MemStatsSource struct {
	getMstatsRec chan chan MemStatsRecord
}

func NewMemStatsSource() *MemStatsSource {
	getCh := make(chan chan MemStatsRecord)
	source := &MemStatsSource{getCh}
	go source.loadStats()
	return source
}

func (ms *MemStatsSource) Name() string {
	return "mem"
}

func (ms *MemStatsSource) Contents(name string) ([]byte, bool) {
	content, _, found := ms.itemInfo(name)
	return content, found
}

func (ms *MemStatsSource) itemInfo(name string) ([]byte, time.Time, bool) {
	stats, updated := ms.getMemStats()

	var content string
	switch name {
	case "genmem":
		content = fmt.Sprintf("%d %d %d %d %d %d\n", stats.Alloc, stats.TotalAlloc, stats.Sys, stats.Lookups, stats.Mallocs, stats.Frees)
	case "heap":
		content = fmt.Sprintf("%d %d %d %d %d %d\n", stats.HeapAlloc, stats.HeapSys, stats.HeapIdle, stats.HeapInuse, stats.HeapReleased, stats.HeapObjects)
	case "stack":
		content = fmt.Sprintf("%d %d\n", stats.StackInuse, stats.StackSys)
	case "mspan":
		content = fmt.Sprintf("%d %d\n", stats.MSpanInuse, stats.MSpanSys)
	case "mcache":
		content = fmt.Sprintf("%d %d\n", stats.MCacheInuse, stats.MCacheSys)
	case "buckethash":
		content = fmt.Sprintf("%d\n", stats.BuckHashSys)
	case "gc":
		content = fmt.Sprintf("%d %d %d %d %t %t\n", stats.NextGC, stats.LastGC, stats.PauseTotalNs, stats.NumGC, stats.EnableGC, stats.DebugGC)
	default:
		return []byte{}, updated, false
	}

	return []byte(content), updated, true
}

func (ms *MemStatsSource) getMemStats() (*runtime.MemStats, time.Time) {
	ret := make(chan MemStatsRecord)
	ms.getMstatsRec <- ret
	result := <-ret
	return &result.Stats, result.LastUpdated
}

func (ms *MemStatsSource) loadStats() {
	var statLastUpdated time.Time
	var stats runtime.MemStats
	for {
		select {
		case ret := <-ms.getMstatsRec:
			now := time.Now()
			if now.After(statLastUpdated.Add(2 * time.Second)) {
				runtime.ReadMemStats(&stats)
				statLastUpdated = now
			}
			ret <- MemStatsRecord{stats, statLastUpdated}
		}
	}
}
