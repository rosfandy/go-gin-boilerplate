package config

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type BracketFormatter struct{}

func (f *BracketFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	timestamp := entry.Time.Format("2006/01/02 - 15:04:05")
	levelText := strings.ToUpper(entry.Level.String())

	b.WriteString(colorizeLevel(levelText, entry.Level))

	b.WriteString(" ")
	b.WriteString(timestamp)

	b.WriteString(": ")
	b.WriteString(entry.Message)

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		if k == "category" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		b.WriteString(" ")
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(fmt.Sprint(entry.Data[k]))
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func colorizeLevel(levelText string, level logrus.Level) string {
	const (
		reset  = "\x1b[0m"
		red    = "\x1b[31m"
		yellow = "\x1b[33m"
		green  = "\x1b[32m"
		blue   = "\x1b[34m"
		gray   = "\x1b[37m"
	)

	color := gray
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		color = red
	case logrus.WarnLevel:
		color = yellow
	case logrus.InfoLevel:
		color = green
	case logrus.DebugLevel, logrus.TraceLevel:
		color = blue
	}

	return color + "[" + levelText + "]" + reset
}

func NewLogrus() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&BracketFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

func NewLogrusWithCategory(category string) *logrus.Entry {
	logger := NewLogrus()

	return logger.WithField("category", strings.ToUpper(category))
}
