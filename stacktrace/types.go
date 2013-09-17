package stacktrace

import (
    "time"
)

type StacktraceRecord struct {
    Stacktrace []byte
    LastUpdated time.Time
}
