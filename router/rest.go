package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/max-gui/logagent/internal/pkg/logagent"
	// nethttp "net/http"
)

func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()
		if raw != "" {
			path = path + "?" + raw
		}
		timestamp := time.Since(start)
		var infomsg string
		if c.Errors.String() == "" {
			infomsg = fmt.Sprintf("%s %s %s from %s cost %s;bodysize is %s;",
				strconv.Itoa(c.Writer.Status()), c.Request.Method, path, c.ClientIP(), timestamp, strconv.Itoa(c.Writer.Size()))
		} else {
			infomsg = fmt.Sprintf("%s %s %s from %s cost %s;bodysize is %s;errormsg: %s",
				strconv.Itoa(c.Writer.Status()), c.Request.Method, path, c.ClientIP(), timestamp, strconv.Itoa(c.Writer.Size()), c.Errors.String())
		}
		logagent.Inst(c).
			WithField("timestamp", timestamp).
			WithField("clientip", c.ClientIP()).
			WithField("method", c.Request.Method).
			WithField("statuscode", c.Writer.Status()).
			WithField("error", c.Errors.String()).
			WithField("bodysize", c.Writer.Size()).
			WithField("path", path).Infof(infomsg)
		// Log only when path is not being skipped

	}
}

func GinHeaderMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// // "trace": "%X{X-B3-TraceId:-}",？
		// // "span": "%X{X-B3-SpanId:-}",？
		// // "parent": "%X{X-B3-ParentSpanId:-}",？
		// // "x-baggage-AF-env": "%X{x-baggage-AF-env:-}",？
		// // "x-baggage-AF-region": "%X{x-baggage-AF-region:-}",？
		trace := c.Request.Header.Get("X-B3-TraceId")
		if trace == "" {
			trace = strings.ReplaceAll(uuid.NewString(), "-", "")
		}
		span := strings.ReplaceAll(uuid.NewString(), "-", "")[0:16]
		// c.Set("trace", c.Request.Header.Get("X-B3-TraceId"))
		// c.Set("span", c.Request.Header.Get("X-B3-SpanId"))
		c.Set("trace", trace)
		c.Set("span", span)
		// c.Set("parentspanid", c.Request.Header.Get("X-B3-ParentSpanId"))
		c.Set("env", c.Request.Header.Get("x-baggage-AF-env"))
		c.Set("region", c.Request.Header.Get("x-baggage-AF-region"))
		// logger := logagent.Inst(c)
		// logger.Print(c.Request.Header)
		// c.Set("log", logagent.Inst(c))
		// log.Print(c.Get("trace"))
		// log.Print(c.Get("span"))
		// log.Print(c.Get("X-B3-ParentSpanId"))
		// log.Print(c.Value("env"))
		// log.Print(c.Get("env"))
		// log.Print(c.Get("region"))
		c.Next()
		// host := c.Request.Host
		// fmt.Printf("Before: %s\n", host)
		// c.Next()
		// fmt.Println("Next: ...")
	}
}

func GinErrorMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": fmt.Sprint(e),
				})
				logger := logagent.Inst(c)
				logger.Panic(e)
			}
		}()

		c.Next()
		// host := c.Request.Host
		// fmt.Printf("Before: %s\n", host)
		// c.Next()
		// fmt.Println("Next: ...")
	}
}
