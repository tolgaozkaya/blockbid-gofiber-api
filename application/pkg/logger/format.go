package logger

import (
	"fmt"
	"path"
	"strings"
	"time"

	l "github.com/sirupsen/logrus"
)

func (f *Formatter) Format(entry *l.Entry) ([]byte, error) {
	var res []string
	res = append(res, entry.Time.Format(time.RFC3339))
	res = append(res, fmt.Sprintf("[%s]", strings.ToUpper(entry.Level.String())))
	if entry.HasCaller() && f.MinLevel <= entry.Logger.GetLevel() {
		res = append(res, fmt.Sprintf("[%s()]", path.Base(entry.Caller.Function)))
		res = append(res, fmt.Sprintf("[%s:%d]", path.Base(entry.Caller.File), entry.Caller.Line))
	}
	res = append(res, strings.TrimSpace(entry.Message))
	return []byte(strings.Join(res, " ") + "\n"), nil
}
