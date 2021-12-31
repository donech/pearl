//package ginzap provides log handling using ginzap package.
// Code structure based on ginrus package.
// Reference: github.com/xgin-contrib/zap

package middleware

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/donech/pearl/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var connectionNum int64

type bodyLogWriter struct {
	gin.ResponseWriter
	body      *bytes.Buffer
	mode      string
	writeBody bool
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if w.needWriteBody() {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	if w.needWriteBody() {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w bodyLogWriter) needWriteBody() bool {
	if w.writeBody {
		return w.writeBody
	}
	if w.mode == gin.DebugMode {
		w.writeBody = true
	}
	contentType := w.Header().Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		w.writeBody = true
	}
	return w.writeBody
}

// GinZap returns a xgin.HandlerFunc (middleware) that logs requests using uber-go/ginzap.
//
// Requests with errors are logged using ginzap.Error().
// Requests without errors are logged using ginzap.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
//   3. A string stating whether to log response body
func GinZap(timeFormat string, utc bool, mod string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer, mode: mod}
		c.Writer = bodyLogWriter
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		header := c.Request.Header

		// 注入 trace-id
		traceID := log.GetTraceIDFromHTTPHeader(header)
		var ctx context.Context
		if traceID == "" {
			ctx = log.NewCtxWithTraceID(c.Request.Context())
			traceID = log.GetTraceIDFromContext(ctx)
		} else {
			ctx = context.WithValue(c.Request.Context(), log.TraceKeyName, traceID)
		}
		c.Request = c.Request.WithContext(ctx)
		c.Header(string(log.TraceKeyName), traceID)
		log.L(ctx).Info("Request receive:",
			zap.String(string(log.TraceKeyName), traceID),
			zap.String("path", path),
			zap.String("method", c.Request.Method),
			zap.Reflect("header", header),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", start.Format(timeFormat)),
		)

		//不再记录请求body, 接收参数由内部方法自行记录
		//if c.Request.ContentLength < 1024 {
		//	body, err := ioutil.ReadAll(c.Request.Body)
		//	if err != nil {
		//		logger.Error("读取请求 body 失败", zap.String("err", err.Error()))
		//	}
		//	//decoded, err := base64.StdEncoding.DecodeString(string(body))
		//	//if err != nil {
		//	//	logger.Error("decode body 失败", zap.String("err", err.Error()))
		//	//}
		//	logger.Info("Request body:",
		//		zap.String(log.TraceKeyName, traceID),
		//		zap.String("body", string(body)),
		//	)
		//}
		connectionNum++
		c.Next()
		connectionNum--
		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				log.L(ctx).Error(e)
			}
		} else {
			responseBody := bodyLogWriter.body.String()
			header = bodyLogWriter.Header()
			log.L(ctx).Info("Request response:：",
				zap.String(string(log.TraceKeyName), traceID),
				zap.String("path", path),
				zap.String("method", c.Request.Method),
				zap.Reflect("header", header),
				zap.String("query", query),
				zap.Int("status", c.Writer.Status()),
				zap.String("time", end.Format(timeFormat)),
				zap.Duration("latency", latency),
				zap.String("body", responseBody),
			)
		}
	}
}

// RecoveryWithZap returns a ginzap.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/ginzap.
// All errors are logged using ginzap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func RecoveryWithZap(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			traceID := log.GetTraceIDFromHTTPHeader(c.Writer.Header())

			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack log.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					log.L(c.Request.Context()).Error(c.Request.URL.Path,
						zap.String(string(log.TraceKeyName), traceID),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					log.L(c.Request.Context()).Error("[Recovery from panic]",
						zap.String(string(log.TraceKeyName), traceID),
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					log.L(c.Request.Context()).Error("[Recovery from panic]",
						zap.String(string(log.TraceKeyName), traceID),
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}

//GetConnectionNum http connectionNum
func GetConnectionNum() int64 {
	return connectionNum
}
