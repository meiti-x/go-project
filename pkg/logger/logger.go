package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	"agentic/commerce/config"
	"github.com/rs/zerolog"
)

type AsyncWriter struct {
	ch     chan []byte
	writer io.Writer
	quit   chan struct{}
	wg     sync.WaitGroup
}

func NewAsyncWriter(w io.Writer, bufferSize int) *AsyncWriter {
	aw := &AsyncWriter{
		ch:     make(chan []byte, bufferSize),
		writer: w,
		quit:   make(chan struct{}),
	}
	aw.wg.Add(1)
	go aw.writeLoop()
	return aw
}

func (aw *AsyncWriter) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))
	copy(buf, p)

	select {
	case aw.ch <- buf:
		return len(p), nil
	case <-aw.quit:
		return 0, io.ErrClosedPipe
	}
}

func (aw *AsyncWriter) writeLoop() {
	defer aw.wg.Done()
	for {
		select {
		case buf := <-aw.ch:
			_, err := aw.writer.Write(buf)

			if err != nil {
				log.Fatal("write error:", err)
			}
		case <-aw.quit:
			close(aw.ch)
			for buf := range aw.ch {
				_, err := aw.writer.Write(buf)

				if err != nil {
					log.Fatal("write error:", err)
				}

			}
			return
		}
	}
}

func (aw *AsyncWriter) Close() error {
	close(aw.quit)
	aw.wg.Wait()
	return nil
}

// AppLogger encapsulates zerolog logger and config
type AppLogger struct {
	cfg    *config.Config
	logger zerolog.Logger
	writer *AsyncWriter
}

// NewAppLogger creates new AppLogger with config
func NewAppLogger(cfg *config.Config) *AppLogger {
	return &AppLogger{cfg: cfg}
}

func (l *AppLogger) InitLogger() {
	if l.cfg.Logger == nil || l.cfg.Logger.Level == "" {
		l.cfg.Logger = &config.Logger{Level: config.LevelInfo}
	}
	level, err := zerolog.ParseLevel(string(l.cfg.Logger.Level))
	if err != nil {
		fmt.Println("invalid log level:", err)
	}

	asyncWriter := NewAsyncWriter(os.Stdout, 1000)
	l.writer = asyncWriter

	logger := zerolog.New(asyncWriter).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(level)

	if l.cfg.IsDev() {
		logger = logger.Output(zerolog.ConsoleWriter{Out: asyncWriter})
	}

	l.logger = logger
}

// Trace logs trace level message with optional fields
func (l *AppLogger) Trace(msg string, fields ...interface{}) {
	if l.logger.GetLevel() > zerolog.TraceLevel {
		return
	}
	l.log(l.logger.Trace(), msg, fields)
}

// Debug logs debug level message with optional fields
func (l *AppLogger) Debug(msg string, fields ...interface{}) {
	if l.logger.GetLevel() > zerolog.DebugLevel {
		return
	}
	l.log(l.logger.Debug(), msg, fields)
}

// Info logs info level message with optional fields
func (l *AppLogger) Info(msg string, fields ...interface{}) {
	if l.logger.GetLevel() > zerolog.InfoLevel {
		return
	}
	l.log(l.logger.Info(), msg, fields)
}

// Warn logs warn level message with optional fields
func (l *AppLogger) Warn(msg string, fields ...interface{}) {
	if l.logger.GetLevel() > zerolog.WarnLevel {
		return
	}
	l.log(l.logger.Warn(), msg, fields)
}

// Error logs error level message with error field
func (l *AppLogger) Error(msg string, err error, fields ...interface{}) {
	if l.logger.GetLevel() > zerolog.ErrorLevel {
		return
	}
	l.error(l.logger.Error(), err, msg, fields)
}

// Panic logs dpanic level message with error field
func (l *AppLogger) Panic(msg string, err error, fields ...interface{}) {
	if l.logger.GetLevel() > zerolog.PanicLevel {
		return
	}
	l.error(l.logger.Panic(), err, msg, fields)
}

// Fatal logs fatal level message with error field
func (l *AppLogger) Fatal(msg string, err error, fields ...interface{}) {
	if l.logger.GetLevel() > zerolog.FatalLevel {
		return
	}
	l.error(l.logger.Fatal(), err, msg, fields)
}

// WithScope adds module name to the logger context
func (l *AppLogger) WithScope(instance any) *AppLogger {
	var moduleName string

	switch v := instance.(type) {
	case string:
		moduleName = v
	default:
		moduleName = reflect.TypeOf(instance).Name()
	}

	newLogger := l.logger.With().Str("module", moduleName).Logger()
	return &AppLogger{
		cfg:    l.cfg,
		logger: newLogger,
		writer: l.writer,
	}
}

// With adds custom field name to the logger context
func (l *AppLogger) With(key string, value interface{}) *AppLogger {
	newLogger := l.logger.With().Interface(key, value).Logger()
	return &AppLogger{
		cfg:    l.cfg,
		logger: newLogger,
		writer: l.writer,
	}
}

// Close closes AsyncWriter to flush logs
func (l *AppLogger) Close() error {
	if l.writer != nil {
		l.logger.Debug().Msg("close logger writer")
		return l.writer.Close()
	}
	return nil
}

// GetLogLevel get current log level
func (l *AppLogger) GetLogLevel() zerolog.Level {
	return l.logger.GetLevel()
}

func (l *AppLogger) log(event *zerolog.Event, msg string, fields []interface{}) {
	event.Msgf(convertBracesToPlaceholder(msg), derefAll(fields)...)
}

func (l *AppLogger) error(event *zerolog.Event, err error, msg string, fields []interface{}) {
	if err != nil {
		event = event.Err(err)
	}
	l.log(event, msg, fields)
}

// deepDeref: follows pointers and expands slices/maps so fmt prints values, not addresses
func deepDeref(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return "<nil>"
	}
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return "<nil>"
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		n := rv.Len()
		out := make([]interface{}, n)
		for i := 0; i < n; i++ {
			out[i] = deepDeref(rv.Index(i).Interface())
		}
		return out
	case reflect.Map:
		out := make(map[interface{}]interface{}, rv.Len())
		for _, k := range rv.MapKeys() {
			out[deepDeref(k.Interface())] = deepDeref(rv.MapIndex(k).Interface())
		}
		return out
	case reflect.Struct:
		out := make(map[string]interface{}, rv.NumField())
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			// skip unexported fields
			if rt.Field(i).PkgPath != "" {
				continue
			}
			name := rt.Field(i).Name
			out[name] = deepDeref(rv.Field(i).Interface())
		}
		return out
	default:
		return rv.Interface()
	}
}

func derefAll(vals []interface{}) []interface{} {
	out := make([]interface{}, len(vals))
	for i, v := range vals {
		out[i] = deepDeref(v)
	}
	return out
}

func convertBracesToPlaceholder(s string) string {
	var b strings.Builder
	// single pass
	for i := 0; i < len(s); {
		// escaped \{} -> literal {}
		if s[i] == '\\' && i+2 < len(s) && s[i+1] == '{' && s[i+2] == '}' {
			b.WriteString("{}")
			i += 3
			continue
		}
		// placeholder {} -> %v
		if s[i] == '{' && i+1 < len(s) && s[i+1] == '}' {
			b.WriteString("%+v")
			i += 2
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
