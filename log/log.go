package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	logger     Logger
	bufferPool *sync.Pool
)

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	logger = NewStdLogger()
}

// SetLogger 设置自定义日志对象
func SetLogger(l Logger) {
	logger = l
}

// Logger 日志输出interface
type Logger interface {
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
	Errorf(string, ...interface{})
	Errorln(...interface{})
	Warnf(string, ...interface{})
	Warnln(...interface{})
	Infof(string, ...interface{})
	Infoln(...interface{})
	Tracef(string, ...interface{})
	Traceln(...interface{})
}

// Level 日志级别类型
type Level int

func (l Level) String() string {
	switch l {
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	case ErrorLevel:
		return "ERROR"
	case WarningLevel:
		return "WARNING"
	case InfoLevel:
		return "INFO"
	case TraceLevel:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel 根据提供的日志级别字符串转换为类型
func ParseLevel(level string) (Level, error) {
	switch strings.ToUpper(level) {
	case "PANIC":
		return PanicLevel, nil
	case "FATAL":
		return FatalLevel, nil
	case "WARNING":
		return WarningLevel, nil
	case "INFO":
		return InfoLevel, nil
	case "TRACE":
		return TraceLevel, nil
	default:
		var l Level
		return l, fmt.Errorf("invalid log level: %s", level)
	}
}

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarningLevel
	InfoLevel
	TraceLevel
)

func NewStdLogger() *StdLogger {
	return &StdLogger{
		out:   os.Stdout,
		level: TraceLevel,
		mu:    sync.Mutex{},
	}
}

type StdLogger struct {
	out   io.Writer
	level Level
	mu    sync.Mutex
}

func (logger *StdLogger) SetOutput(output io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.out = output
}

func (logger *StdLogger) SetLevel(l Level) {
	logger.level = l
}

func (logger *StdLogger) GetLevel() Level {
	return logger.level
}

func (logger *StdLogger) isLevelEnable(level Level) bool {
	return logger.GetLevel() >= level
}

func (logger *StdLogger) output(l Level, msg string) {
	if !logger.isLevelEnable(l) {
		return
	}
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)

	buffer.WriteString(l.String())
	buffer.WriteString(" ")
	buffer.WriteString(time.Now().Format(time.RFC3339Nano))
	buffer.WriteString(" ")
	buffer.WriteString(msg)
	buffer.WriteString("\n")

	logger.out.Write(buffer.Bytes())
}

func (logger *StdLogger) Fatalf(format string, args ...interface{}) {
	logger.output(FatalLevel, fmt.Sprintf(format, args...))
}

func (logger *StdLogger) Fatalln(args ...interface{}) {
	logger.output(WarningLevel, fmt.Sprintln(args...))
}

func (logger *StdLogger) Errorf(format string, args ...interface{}) {
	logger.output(ErrorLevel, fmt.Sprintf(format, args...))
}

func (logger *StdLogger) Errorln(args ...interface{}) {
	logger.output(ErrorLevel, fmt.Sprintln(args...))
}

func (logger *StdLogger) Warnf(format string, args ...interface{}) {
	logger.output(WarningLevel, fmt.Sprintf(format, args...))
}

func (logger *StdLogger) Warnln(args ...interface{}) {
	logger.output(WarningLevel, fmt.Sprintln(args...))
}

func (logger *StdLogger) Infof(format string, args ...interface{}) {
	logger.output(InfoLevel, fmt.Sprintf(format, args...))
}

func (logger *StdLogger) Infoln(args ...interface{}) {
	logger.output(InfoLevel, fmt.Sprintln(args...))
}

func (logger *StdLogger) Tracef(format string, args ...interface{}) {
	logger.output(TraceLevel, fmt.Sprintf(format, args...))
}

func (logger *StdLogger) Traceln(args ...interface{}) {
	logger.output(TraceLevel, fmt.Sprintln(args...))
}
