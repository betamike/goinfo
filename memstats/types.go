package memstats

import (
    "runtime"
    "time"
)

type MemStatsRecord struct {
    Stats runtime.MemStats
    LastUpdated time.Time
}
