package lager

import (
	"io"
	"sync"
)

// A Sink represents a write destination for a Logger. It provides
// a thread-safe interface for writing logs
type Sink interface {
	//Log to the sink.  Best effort -- no need to worry about errors.
	Log(LogFormat)
}

type writerSink struct {
	writer      io.Writer
	minLogLevel LogLevel
	writeL      *sync.Mutex
}

func NewWriterSink(writer io.Writer, minLogLevel LogLevel) Sink {
	return &writerSink{
		writer:      writer,
		minLogLevel: minLogLevel,
		writeL:      new(sync.Mutex),
	}
}

func (sink *writerSink) Log(log LogFormat) {
	if log.LogLevel < sink.minLogLevel {
		return
	}

	// Convert to json outside of critical section to minimize time spent holding lock
	message := append(log.ToJSON(), '\n')

	sink.writeL.Lock()
	sink.writer.Write(message) //nolint:errcheck
	sink.writeL.Unlock()
}

type prettySink struct {
	writer      io.Writer
	minLogLevel LogLevel
	writeL      sync.Mutex
}

func NewPrettySink(writer io.Writer, minLogLevel LogLevel) Sink {
	return &prettySink{
		writer:      writer,
		minLogLevel: minLogLevel,
	}
}

func (sink *prettySink) Log(log LogFormat) {
	if log.LogLevel < sink.minLogLevel {
		return
	}

	// Convert to json outside of critical section to minimize time spent holding lock
	message := append(log.toPrettyJSON(), '\n')

	sink.writeL.Lock()
	sink.writer.Write(message) //nolint:errcheck
	sink.writeL.Unlock()
}
