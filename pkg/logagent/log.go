package logagent

import (

	// rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/google/uuid"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	Logagent *logrus.Logger
	once     sync.Once
)

func GetRootContextWithTrace() context.Context {
	trace := strings.ReplaceAll(uuid.NewString(), "-", "")
	span := strings.ReplaceAll(uuid.NewString(), "-", "")

	d := context.WithValue(context.Background(), "trace", trace)
	c := context.WithValue(d, "span", span)

	return c
}

func Inst(c context.Context) *logrus.Entry {
	once.Do(func() {

		logfile := fmt.Sprintf("logs/json/%s.json", *logsets.Appname)
		pathMap := lfshook.PathMap{
			logrus.InfoLevel:  logfile,
			logrus.ErrorLevel: logfile,
			logrus.PanicLevel: logfile,
			logrus.FatalLevel: logfile,
		}

		Logagent = logrus.New()

		// Logagent.Hooks.Add(lfshook.NewHook(
		// 	os.Stdout,
		// 	&logrus.TextFormatter{
		// 		DisableColors:    false,
		// 		FullTimestamp:    false,
		// 		DisableTimestamp: true,
		// 	},
		// ))

		if *logsets.Jsonlog {
			// Logagent.Out = nil
			Logagent.Hooks.Add(lfshook.NewHook(
				pathMap,
				&logrus.JSONFormatter{
					FieldMap: logrus.FieldMap{
						logrus.FieldKeyTime: "@timestamp",
						logrus.FieldKeyMsg:  "message",
					},
					TimestampFormat: time.RFC3339Nano},
			))
		} else {

			Logagent.Out = os.Stdout
		}
	})
	// logrus.SetFormatter(&logrus.JSONFormatter{
	// 	FieldMap: logrus.FieldMap{
	// 		logrus.FieldKeyTime: "@timestamp",
	// 		logrus.FieldKeyMsg:  "message",
	// 	},
	// })

	Logagent.SetLevel(logrus.InfoLevel)
	Logagent.SetReportCaller(true)

	// logrus.SetOutput(os.Stdout)
	// file, err := os.OpenFile("out.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err == nil {
	// 	logrus.SetOutput(file)
	// }
	// defer file.Close()
	var logger *logrus.Entry

	if *logsets.Jsonlog {
		c.Value("trace")
		c.Value("span")
		c.Value("env")
		c.Value("region")
		logger = Logagent.WithFields(logrus.Fields{
			"service":             logsets.Appname,
			"app_name":            logsets.Appname,
			"app-env":             logsets.Appenv,
			"app-region":          "default",
			"trace":               c.Value("trace"),
			"span":                c.Value("span"),
			"x-baggage-AF-env":    c.Value("env"),
			"x-baggage-AF-region": c.Value("region"),
		})
	} else {
		Logagent.SetFormatter(&nested.Formatter{
			HideKeys: true,
			// FieldsOrder: []string{"component", "category"},
		})

		logger = Logagent.WithContext(c)
	}
	// log.Print()
	return logger
	// // "trace": "%X{X-B3-TraceId:-}",？
	// // "span": "%X{X-B3-SpanId:-}",？
	// // "parent": "%X{X-B3-ParentSpanId:-}",？
	// // "logger_name": "%logger{40}",？
	// // "x-baggage-AF-env": "%X{x-baggage-AF-env:-}",？
	// // "x-baggage-AF-region": "%X{x-baggage-AF-region:-}",？
	// fields := log.Fields{"userId": 12}
	// log.WithFields(fields).Info("User logged in!")

	// fields = log.Fields{"userId": 12}
	// log.WithFields(fields).Info("Sent a message!")

	// fields = log.Fields{"userId": 12}
	// log.WithFields(fields).Info("Failed to get a message!")

	// fields = log.Fields{"userId": 12}
	// log.WithFields(fields).Info("User logged out!")
}

// func init() {
// 	// if Logagent != nil {
// 	// 	return Logagent
// 	// }
// 	logfile := fmt.Sprintf("logs/json/%s.json", *logsets.Appname)
// 	pathMap := lfshook.PathMap{
// 		logrus.InfoLevel:  logfile,
// 		logrus.ErrorLevel: logfile,
// 		logrus.PanicLevel: logfile,
// 		logrus.FatalLevel: logfile,
// 	}

// 	Logagent = logrus.New()

// 	// Logagent.Hooks.Add(lfshook.NewHook(
// 	// 	os.Stdout,
// 	// 	&logrus.TextFormatter{
// 	// 		DisableColors:    false,
// 	// 		FullTimestamp:    false,
// 	// 		DisableTimestamp: true,
// 	// 	},
// 	// ))

// 	if *logsets.Jsonlog {
// 		// Logagent.Out = nil
// 		Logagent.Hooks.Add(lfshook.NewHook(
// 			pathMap,
// 			&logrus.JSONFormatter{
// 				FieldMap: logrus.FieldMap{
// 					logrus.FieldKeyTime: "@timestamp",
// 					logrus.FieldKeyMsg:  "message",
// 				},
// 				TimestampFormat: time.RFC3339Nano,
// 			},
// 		))
// 	} else {

// 		Logagent.Out = os.Stdout
// 	}
// 	// return Logagent
// }
