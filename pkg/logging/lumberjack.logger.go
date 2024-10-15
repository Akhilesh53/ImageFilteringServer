package logging

import (
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

// FlushingLumberjackLogger wraps lumberjack.Logger and adds sync before rotating.
type FlushingLumberjackLogger struct {
	*lumberjack.Logger
	sync.Mutex
}

// Write ensures thread-safe writes.
func (f *FlushingLumberjackLogger) Write(p []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	return f.Logger.Write(p)
}

// Rotate calls zap.Sync before rotating the log file.
func (f *FlushingLumberjackLogger) Rotate() error {
	f.Lock()
	defer f.Unlock()
	if err := logInstance.Sync(); err != nil {
		return err
	}
	return f.Logger.Rotate()
}
