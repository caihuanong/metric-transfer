package log

import (
	"bytes"
	"fmt"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mgutz/ansi"
)

const reset = ansi.Reset

var (
	baseTimestamp time.Time
	isTerminal    bool
)

func init() {
	baseTimestamp = time.Now()
	isTerminal = logrus.IsTerminal()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging of just the time passed since beginning of execution.
	ShortTimestamp bool

	// Timestamp format to use for display when a full timestamp is printed.
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// Pad msg field with spaces on the right for display.
	// The value for this parameter will be the size of padding.
	// Its default value is zero, which means no padding will be applied for msg.
	SpacePadding int
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var keys []string = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		if k != "prefix" && k != "caller" {
			keys = append(keys, k)
		}
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}

	b := &bytes.Buffer{}

	prefixFieldClashes(entry.Data)

	isColorTerminal := isTerminal && (runtime.GOOS != "windows")
	isColored := (f.ForceColors || isColorTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.Stamp
	}
	if isColored {
		f.printColored(b, entry, keys, timestampFormat)
	} else {
		f.printNoColored(b, entry, keys, timestampFormat)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string) {
	var levelColor string
	var levelText string
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = ansi.Green
	case logrus.WarnLevel:
		levelColor = ansi.Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = ansi.Red
	default:
		levelColor = ansi.Blue
	}

	if entry.Level != logrus.WarnLevel {
		levelText = strings.ToUpper(entry.Level.String())
	} else {
		levelText = "WARN"
	}

	prefix := " "
	message := entry.Message

	if prefixValue, ok := entry.Data["prefix"]; ok {
		prefix = fmt.Sprint(ansi.Cyan, prefixValue, ":", reset, " ")
	} else {
		prefixValue, trimmedMsg := extractPrefix(entry.Message)
		if len(prefixValue) > 0 {
			prefix = fmt.Sprint(ansi.Cyan, prefixValue, ":", reset, " ")
			message = trimmedMsg
		}
	}

	caller, _ := entry.Data["caller"]
	messageFormat := "%s"
	if f.SpacePadding != 0 {
		messageFormat = fmt.Sprintf("%%-%ds", f.SpacePadding)
	}

	if f.ShortTimestamp {
		fmt.Fprintf(b, "%s[%s %04d %s]%s%s"+messageFormat, levelColor, levelText[:1], miniTS(), caller, reset, prefix, message)
	} else {
		fmt.Fprintf(b, "%s[%s %s %s]%s%s"+messageFormat, levelColor, levelText[:1], entry.Time.Format(timestampFormat), caller, reset, prefix, message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " %s%s%s=%+v", levelColor, k, reset, v)
	}
}

func (f *TextFormatter) printNoColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string) {
	var levelText string

	if entry.Level != logrus.WarnLevel {
		levelText = strings.ToUpper(entry.Level.String())
	} else {
		levelText = "WARN"
	}

	prefix := " "
	message := entry.Message

	if prefixValue, ok := entry.Data["prefix"]; ok {
		prefix = fmt.Sprint(prefixValue, ":", " ")
	} else {
		prefixValue, trimmedMsg := extractPrefix(entry.Message)
		if len(prefixValue) > 0 {
			prefix = fmt.Sprint(prefixValue, ":", " ")
			message = trimmedMsg
		}
	}

	caller, _ := entry.Data["caller"]
	messageFormat := "%s"
	if f.SpacePadding != 0 {
		messageFormat = fmt.Sprintf("%%-%ds", f.SpacePadding)
	}

	if f.ShortTimestamp {
		fmt.Fprintf(b, "[%s %04d %s]%s"+messageFormat, levelText[:1], miniTS(), caller, prefix, message)
	} else {
		fmt.Fprintf(b, "[%s %s %s]%s"+messageFormat, levelText[:1], entry.Time.Format(timestampFormat), caller, prefix, message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " %s=%+v", k, v)
	}
}

func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return false
		}
	}
	return true
}

func extractPrefix(msg string) (string, string) {
	prefix := ""
	regex := regexp.MustCompile("^\\[(.*?)\\]")
	if regex.MatchString(msg) {
		match := regex.FindString(msg)
		prefix, msg = match[1:len(match)-1], strings.TrimSpace(msg[len(match):])
	}
	return prefix, msg
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	b.WriteString(key)
	b.WriteByte('=')

	switch value := value.(type) {
	case string:
		if needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	case error:
		errmsg := value.Error()
		if needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	default:
		fmt.Fprint(b, value)
	}

	b.WriteByte(' ')
}

func prefixFieldClashes(data logrus.Fields) {
	_, ok := data["time"]
	if ok {
		data["fields.time"] = data["time"]
	}
	_, ok = data["msg"]
	if ok {
		data["fields.msg"] = data["msg"]
	}
	_, ok = data["level"]
	if ok {
		data["fields.level"] = data["level"]
	}
}
