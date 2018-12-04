package hooks

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"runtime"
	"strings"
)

func GetCallFields() string {
	depth := 0
	occur := false
	for i := 2; ; i++ {
		pc, _, _, _ := runtime.Caller(i)
		if occur {
			if !strings.Contains(runtime.FuncForPC(pc).Name(), "Sirupsen/logrus.") {
				depth = i
				if strings.Contains(runtime.FuncForPC(pc).Name(), "common/log.") {
					depth++
				}
				break
			}
		} else {
			if strings.Contains(runtime.FuncForPC(pc).Name(), "Sirupsen/logrus.") {
				occur = true
			}
		}
	}

	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = funcName[strings.LastIndex(funcName, "/")+1:]
	file = file[strings.LastIndex(file, "/")+1:]
	return fmt.Sprintf("%v(%v:%v)", funcName, file, line)
}

type CallerHook struct {
}

func (hook *CallerHook) Fire(entry *logrus.Entry) error {
	entry.Data["caller"] = GetCallFields()
	return nil
}

func (hook *CallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
