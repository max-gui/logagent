package logagent

import (

	// rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/google/uuid"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/max-gui/logagent/pkg/logsets"
	"github.com/sirupsen/logrus"
)

var (
	Logagent *logrus.Logger
	logger   *logrus.Entry
	once     sync.Once
)

func GetRootContextWithTrace() context.Context {
	trace := strings.ReplaceAll(uuid.NewString(), "-", "")
	span := strings.ReplaceAll(uuid.NewString(), "-", "")

	d := context.WithValue(context.Background(), "trace", trace)
	c := context.WithValue(d, "span", span)

	return c
}

// func zip(source string) {
// 	target := source + ".tmp.gz"
// 	fw, err := os.Create(target) //"demo.gzip") // 创建gzip包文件，返回*io.Writer
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer fw.Close()

// 	// 实例化心得gzip.Writer
// 	gw := gzip.NewWriter(fw)
// 	defer gw.Close()

// 	// 获取要打包的文件信息
// 	fr, err := os.Open(source) //"demo.txt")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer fr.Close()

// 	// 获取文件头信息
// 	fi, err := fr.Stat()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// 创建gzip.Header
// 	gw.Header.Name = fi.Name()

// 	// 读取文件数据
// 	buf := make([]byte, fi.Size())
// 	_, err = fr.Read(buf)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	// 写入数据到zip包
// 	_, err = gw.Write(buf)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	err = os.Remove(source)
// 	if err != nil {
// 		log.Fatalln(filepath.Abs(source))
// 		log.Fatalln(err)
// 	}
// 	err = os.Rename(target, source)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }

func InstArch(c context.Context) *logrus.Entry {
	logger = inst(c).WithField("log_type", "arch")
	once.Do(func() {
		logger = inst(c).WithField("log_type", "arch")
	})

	return logger
}

func InstNomal(c context.Context) *logrus.Entry {
	logger = inst(c).WithField("", "")

	return logger
}

func inst(c context.Context) *logrus.Entry {
	once.Do(func() {

		logfilelink := fmt.Sprintf("logs/json/%s.json", *logsets.Appname)
		logfilepath := "logs/json/" + *logsets.Appname + ".%Y%m%d%H%M.json.gz" // fmt.Sprintf("logs/json/%s%Y%m%d%H%M.json", *logsets.Appname)

		writter, err := rotatelogs.New(
			logfilepath,
			// "/path/to/access_log.%Y%m%d%H%M",
			rotatelogs.WithLinkName(logfilelink),
			rotatelogs.WithRotationTime(30*time.Minute),
			rotatelogs.WithMaxAge(2*time.Hour),
			// rotatelogs.WithHandler(rotatelogs.Handler(rotatelogs.HandlerFunc(func(e rotatelogs.Event) {
			// 	if e.Type() != rotatelogs.FileRotatedEventType {
			// 		return
			// 	}
			// 	source := e.(*rotatelogs.FileRotatedEvent).PreviousFile()
			// 	if source != "" {
			// 		// target := source + ".json.gz"
			// 		// zip(source, target)
			// 		zip(source)
			// 	}
			// 	// Do what you want with the data. This is just an idea:
			// 	// storeLogFileToRemoteStorage(e.(*FileRotatedEvent).PreviousFile())
			// }))),
		)

		if err != nil {
			log.Panic(err)
		}
		// writerMap := lfshook.WriterMap{
		// 	logrus.InfoLevel:  writter,
		// 	logrus.ErrorLevel: writter,
		// 	logrus.PanicLevel: writter,
		// 	logrus.FatalLevel: writter,
		// }

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
			// asyncwritter := NewHook(8*1024, lfshook.NewHook(

			// 	writerMap,
			// 	&logrus.JSONFormatter{
			// 		FieldMap: logrus.FieldMap{
			// 			logrus.FieldKeyTime: "@timestamp",
			// 			logrus.FieldKeyMsg:  "message",
			// 		},
			// 		TimestampFormat: time.RFC3339Nano},
			// ))
			asyncwritter := NewHook(
				8*1024, writter,
				&logrus.JSONFormatter{
					FieldMap: logrus.FieldMap{
						logrus.FieldKeyTime: "@timestamp",
						logrus.FieldKeyMsg:  "message",
					},
					TimestampFormat: time.RFC3339Nano},
			)

			// Logagent.Out = nil
			Logagent.SetNoLock()
			Logagent.SetOutput(io.Discard)
			Logagent.SetFormatter(&NullFormatter{})

			// Logagent.Hooks.Add(lfshook.NewHook(
			// Logagent.AddHook(lfshook.NewHook(

			// 	writerMap,
			// 	&logrus.JSONFormatter{
			// 		FieldMap: logrus.FieldMap{
			// 			logrus.FieldKeyTime: "@timestamp",
			// 			logrus.FieldKeyMsg:  "message",
			// 		},
			// 		TimestampFormat: time.RFC3339Nano},
			// ))
			Logagent.AddHook(asyncwritter)
		} else {

			Logagent.Out = os.Stdout
		}
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
		// var logger *logrus.Entry

		if *logsets.Jsonlog {
			// c.Value("trace")
			// c.Value("span")
			// c.Value("env")
			// c.Value("region")
			logger = Logagent.WithFields(logrus.Fields{
				"service":             *logsets.Appname,
				"app_name":            *logsets.Appname,
				"app-env":             *logsets.Appenv,
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
	})
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
