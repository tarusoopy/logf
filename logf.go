package logf

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// These flags define which text to prefix to each log entry generated by the Logger (compatible with log package).
const (
	Ldate         = 1 << iota              // the date in the local time zone: 2009/01/23
	Ltime                                  // the time in the local time zone: 01:23:23
	Lmicroseconds                          // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                              // full file name and line number: /a/b/c/d.go:23
	Lshortfile                             // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                                   // if Ldate or Ltime is set, use UTC rather than the local time zone
	Llevel                                 // log level of message
	LstdFlags     = Ldate | Ltime | Llevel // initial values for the standard logger
)

const maskStdLogFlags = Ldate | Ltime | Lmicroseconds | Llongfile | Lshortfile | LUTC

//Logger is logger class
type Logger struct {
	lg   *log.Logger // logger
	mu   sync.Mutex  // ensures atomic writes; protects the following fields
	flag int         // properties
	min  Level       // minimum level for filtering
}

//OptFunc is self-referential function for functional options pattern
type OptFunc func(*Logger)

// New creates a new Logger.
func New(opts ...OptFunc) *Logger {
	l := &Logger{lg: log.New(os.Stderr, "", LstdFlags&maskStdLogFlags), flag: LstdFlags, min: TRACE}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

//WithWriter returns function for setting Writer
func WithWriter(w io.Writer) OptFunc {
	return func(l *Logger) {
		if w != nil {
			l.SetOutput(w)
		}
	}
}

//WithFlags returns function for setting flags
func WithFlags(flag int) OptFunc {
	return func(l *Logger) {
		l.SetFlags(flag)
	}
}

//WithPrefix returns function for setting prefix string
func WithPrefix(prefix string) OptFunc {
	return func(l *Logger) {
		l.SetPrefix(prefix)
	}
}

//WithMinLevel returns function for setting minimum level
func WithMinLevel(lv Level) OptFunc {
	return func(l *Logger) {
		l.SetMinLevel(lv)
	}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lg.SetOutput(w)
}

// SetFlags sets the output flags for the logger.
func (l *Logger) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
	l.lg.SetFlags(flag & maskStdLogFlags)
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.lg.SetPrefix(prefix)
}

// SetMinLevel sets the minimum level for the logger.
func (l *Logger) SetMinLevel(lv Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.min = lv
}

//GetLogger returns log.Logger instance
func (l *Logger) GetLogger() *log.Logger {
	return l.lg
}

//Output writes the output for a logging event.
func (l *Logger) Output(lv Level, calldepth int, s string) error {
	if lv >= l.min {
		if (l.flag & Llevel) != 0 {
			return l.lg.Output(calldepth, fmt.Sprintf("[%v] %s", lv, s))
		}
		return l.lg.Output(calldepth, s)
	}
	return nil
}

//lprintf calls l.Output() to print to the logger.
//Arguments are handled in the manner of fmt.Printf.
func (l *Logger) lprintf(lv Level, format string, v ...interface{}) {
	l.Output(lv, 4, fmt.Sprintf(format, v...))
}

//lprint calls l.Output() to print to the logger.
//Arguments are handled in the manner of fmt.Print.
func (l *Logger) lprint(lv Level, v ...interface{}) { l.Output(lv, 4, fmt.Sprint(v...)) }

//lprintln calls l.Output() to print to the logger.
//Arguments are handled in the manner of fmt.Println.
func (l *Logger) lprintln(lv Level, v ...interface{}) { l.Output(lv, 4, fmt.Sprintln(v...)) }

//Tracef calls l.lprintf() to print to the logger.
func (l *Logger) Tracef(format string, v ...interface{}) { l.lprintf(TRACE, format, v...) }

//Trace calls l.lprint() to print to the logger.
func (l *Logger) Trace(v ...interface{}) { l.lprint(TRACE, v...) }

//Traceln calls l.lprintln() to print to the logger.
func (l *Logger) Traceln(v ...interface{}) { l.lprintln(TRACE, v...) }

//Debugf calls l.lprintf() to print to the logger.
func (l *Logger) Debugf(format string, v ...interface{}) { l.lprintf(DEBUG, format, v...) }

//Debug calls l.lprint() to print to the logger.
func (l *Logger) Debug(v ...interface{}) { l.lprint(DEBUG, v...) }

//Debugln calls l.lprintln() to print to the logger.
func (l *Logger) Debugln(v ...interface{}) { l.lprintln(DEBUG, v...) }

//Printf calls l.lprintf() to print to the logger.
func (l *Logger) Printf(format string, v ...interface{}) { l.lprintf(INFO, format, v...) }

//Print calls l.lprint() to print to the logger.
func (l *Logger) Print(v ...interface{}) { l.lprint(INFO, v...) }

//Println calls l.lprintln() to print to the logger.
func (l *Logger) Println(v ...interface{}) { l.lprintln(INFO, v...) }

//Warnf calls l.lprintf() to print to the logger.
func (l *Logger) Warnf(format string, v ...interface{}) { l.lprintf(WARN, format, v...) }

//Warn calls l.lprint() to print to the logger.
func (l *Logger) Warn(v ...interface{}) { l.lprint(WARN, v...) }

//Warnln calls l.lprintln() to print to the logger.
func (l *Logger) Warnln(v ...interface{}) { l.lprintln(WARN, v...) }

//Errorf calls l.lprintf() to print to the logger.
func (l *Logger) Errorf(format string, v ...interface{}) { l.lprintf(ERROR, format, v...) }

//Error calls l.lprint() to print to the logger.
func (l *Logger) Error(v ...interface{}) { l.lprint(ERROR, v...) }

//Errorln calls l.lprintln() to print to the logger.
func (l *Logger) Errorln(v ...interface{}) { l.lprintln(ERROR, v...) }

//Fatalf calls l.lprintf() to print to the logger.
func (l *Logger) Fatalf(format string, v ...interface{}) { l.lprintf(FATAL, format, v...) }

//Fatal calls l.lprint() to print to the logger.
func (l *Logger) Fatal(v ...interface{}) { l.lprint(FATAL, v...) }

//Fatalln calls l.lprintln() to print to the logger.
func (l *Logger) Fatalln(v ...interface{}) { l.lprintln(FATAL, v...) }

//Panicf is equivalent() to l.Output() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.Output(FATAL, 4, s)
	panic(s)
}

//Panic is equivalent() to l.Output() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.Output(FATAL, 4, s)
	panic(s)
}

//Panicln is equivalent() to l.Output() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.Output(FATAL, 4, s)
	panic(s)
}

var std = New()

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) { std.SetOutput(w) }

// SetFlags sets the output flags for the logger.
func SetFlags(flag int) { std.SetFlags(flag) }

// SetPrefix sets the output prefix for the logger.
func SetPrefix(prefix string) { std.SetPrefix(prefix) }

// SetMinLevel sets the minimum level for the logger.
func SetMinLevel(lv Level) { std.SetMinLevel(lv) }

//Output writes the output for a logging event.
func Output(lv Level, calldepth int, s string) error {
	return std.Output(lv, calldepth, s)
}

//Tracef calls std.Tracef() to print to the logger.
func Tracef(format string, v ...interface{}) { std.lprintf(TRACE, format, v...) }

//Trace calls std.Trace() to print to the logger.
func Trace(v ...interface{}) { std.lprint(TRACE, v...) }

//Traceln calls std.Traceln() to print to the logger.
func Traceln(v ...interface{}) { std.lprintln(TRACE, v...) }

//Debugf calls std.Debugf() to print to the logger.
func Debugf(format string, v ...interface{}) { std.lprintf(DEBUG, format, v...) }

//Debug calls std.Debug() to print to the logger.
func Debug(v ...interface{}) { std.lprint(DEBUG, v...) }

//Debugln calls std.Debugln() to print to the logger.
func Debugln(v ...interface{}) { std.lprintln(DEBUG, v...) }

//Printf calls std.Printf() to print to the logger.
func Printf(format string, v ...interface{}) { std.lprintf(INFO, format, v...) }

//Print calls std.Print() to print to the logger.
func Print(v ...interface{}) { std.lprint(INFO, v...) }

//Println calls std.Println() to print to the logger.
func Println(v ...interface{}) { std.lprintln(INFO, v...) }

//Warnf calls std.Warnf() to print to the logger.
func Warnf(format string, v ...interface{}) { std.lprintf(WARN, format, v...) }

//Warn calls std.Warn() to print to the logger.
func Warn(v ...interface{}) { std.lprint(WARN, v...) }

//Warnln calls std.Warnln() to print to the logger.
func Warnln(v ...interface{}) { std.lprintln(WARN, v...) }

//Errorf calls std.Errorf() to print to the logger.
func Errorf(format string, v ...interface{}) { std.lprintf(ERROR, format, v...) }

//Error calls std.Error() to print to the logger.
func Error(v ...interface{}) { std.lprint(ERROR, v...) }

//Errorln calls std.Errorln() to print to the logger.
func Errorln(v ...interface{}) { std.lprintln(ERROR, v...) }

//Fatalf calls std.Fatalf() to print to the logger.
func Fatalf(format string, v ...interface{}) { std.lprintf(FATAL, format, v...) }

//Fatal calls std.Fatal() to print to the logger.
func Fatal(v ...interface{}) { std.lprint(FATAL, v...) }

//Fatalln calls std.Fatalln() to print to the logger.
func Fatalln(v ...interface{}) { std.lprintln(FATAL, v...) }

//Panicf is equivalent() to std.Output() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.Output(FATAL, 4, s)
	panic(s)
}

//Panic is equivalent() to std.Output() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.Output(FATAL, 4, s)
	panic(s)
}

//Panicln is equivalent() to std.Output() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	std.Output(FATAL, 4, s)
	panic(s)
}
