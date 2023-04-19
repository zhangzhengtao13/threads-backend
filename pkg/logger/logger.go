package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type Level int8

type Fields map[string]interface{} // 日志公共字段

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

// 返回日志等级的名称
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	level     Level
	fields    Fields
	callers   []string
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	l := log.New(w, prefix, flag)
	return &Logger{newLogger: l}
}

/*-------------------日志标准化---------------------------*/

func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

func (l *Logger) withLevel(lvl Level) *Logger {
	ll := l.clone()
	ll.level = lvl
	return ll
}

func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}
	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

// 设置当前某一层调用栈的信息（程序计数器、文件信息和行号）
func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		ll.callers = []string{fmt.Sprintf("%s:%d %s", file, line, f.Name())}
	}
	return ll
}

// 设置当前的整个调用栈信息
func (l *Logger) WithCallesFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames((pcs[:depth]))
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callers = append(callers, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	ll := l.clone()
	ll.callers = callers
	return ll
}

func (l *Logger) WithTrace() *Logger {
	ginCtx, ok := l.ctx.(*gin.Context)
	if ok {
		return l.WithFields(Fields{
			"trace_id": ginCtx.MustGet("X-Trace-id"),
			"span_id":  ginCtx.MustGet("X-Span-ID"),
		})
	}
	return l
}

/*-----------------日志格式化及输出----------------*/

func (l *Logger) jsonFormat(message string) map[string]interface{} {
	data := make(Fields, len(l.fields)+4)
	data["level"] = l.level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["callers"] = l.callers
	data["message"] = message
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

func (l *Logger) output(message string) {
	body, err := json.Marshal(l.jsonFormat(message))
	if err != nil {
		l.newLogger.Println(err)
	}
	content := string(body)
	switch l.level {
	case LevelDebug:
		l.newLogger.Println(content)
	case LevelInfo:
		l.newLogger.Println(content)
	case LevelWarn:
		l.newLogger.Println(content)
	case LevelError:
		l.newLogger.Println(content)
	case LevelFatal:
		l.newLogger.Fatalln(content)
	case LevelPanic:
		l.newLogger.Panicln(content)
	}
}

func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.withLevel(LevelDebug).WithContext(ctx).WithTrace().output(fmt.Sprint(v...))
}

// func (l *Logger) Debug(v ...interface{}) {
// 	l.withLevel(LevelDebug).output(fmt.Sprint(v...))
// }

func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.withLevel(LevelDebug).WithContext(ctx).WithTrace().output(fmt.Sprintf(format, v...))
}

// func (l *Logger) Debugf(format string, v ...interface{}) {
// 	l.withLevel(LevelDebug).output(fmt.Sprintf(format, v...))
// }

func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.withLevel(LevelInfo).WithContext(ctx).WithTrace().output(fmt.Sprint(v...))
}

// func (l *Logger) Info(v ...interface{}) {
// 	l.withLevel(LevelInfo).output(fmt.Sprint(v...))
// }

func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.withLevel(LevelInfo).WithContext(ctx).WithTrace().output(fmt.Sprintf(format, v...))
}

// func (l *Logger) Infof(format string, v ...interface{}) {
// 	l.withLevel(LevelInfo).output(fmt.Sprintf(format, v...))
// }

func (l *Logger) Warn(ctx context.Context, v ...interface{}) {
	l.withLevel(LevelWarn).WithContext(ctx).WithTrace().output(fmt.Sprint(v...))
}

// func (l *Logger) Warn(v ...interface{}) {
// 	l.withLevel(LevelWarn).output(fmt.Sprint(v...))
// }

func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.withLevel(LevelWarn).WithContext(ctx).WithTrace().output(fmt.Sprintf(format, v...))
}

// func (l *Logger) Warnf(format string, v ...interface{}) {
// 	l.withLevel(LevelWarn).output(fmt.Sprintf(format, v...))
// }

func (l *Logger) Err(ctx context.Context, v ...interface{}) {
	l.withLevel(LevelError).WithContext(ctx).WithTrace().output(fmt.Sprint(v...))
}

// func (l *Logger) Err(v ...interface{}) {
// 	l.withLevel(LevelError).output(fmt.Sprint(v...))
// }

func (l *Logger) Errf(ctx context.Context, format string, v ...interface{}) {
	l.withLevel(LevelError).WithContext(ctx).WithTrace().output(fmt.Sprintf(format, v...))
}

// func (l *Logger) Errf(format string, v ...interface{}) {
// 	l.withLevel(LevelError).output(fmt.Sprintf(format, v...))
// }

func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
		l.withLevel(LevelFatal).WithContext(ctx).WithTrace().output(fmt.Sprint(v...))
	}

// func (l *Logger) Fatal(v ...interface{}) {
// 	l.withLevel(LevelFatal).output(fmt.Sprint(v...))
// }

func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.withLevel(LevelFatal).WithContext(ctx).WithTrace().output(fmt.Sprintf(format, v...))
}

// func (l *Logger) Fatalf(format string, v ...interface{}) {
// 	l.withLevel(LevelFatal).output(fmt.Sprintf(format, v...))
// }

func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.withLevel(LevelPanic).WithContext(ctx).WithTrace().output(fmt.Sprint(v...))
}

// func (l *Logger) Panic(v ...interface{}) {
// 	l.withLevel(LevelPanic).output(fmt.Sprint(v...))
// }

func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.withLevel(LevelPanic).WithContext(ctx).WithTrace().output(fmt.Sprintf(format, v...))
}


// func (l *Logger) Panicf(format string, v ...interface{}) {
// 	l.withLevel(LevelPanic).output(fmt.Sprintf(format, v...))
// }
