package log

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
)

//TraceKeyName http stander header style
const TraceKeyName TraceKey = "X-Request-ID"

type TraceKey string

var randomFuc func() uint64

//ID traceID
type ID struct {
	High, Low uint64
}

//String String
func (o ID) String() string {
	if o.High == 0 {
		return fmt.Sprintf("%x", o.Low)
	}
	return fmt.Sprintf("%x%016x", o.High, o.Low)
}

// IsValid checks if the ID is valid, i.e. not zero.
func (o ID) IsValid() bool {
	return o.High != 0 || o.Low != 0
}

func init() {
	rand.Seed(time.Now().UnixNano())
	randomFuc = func() uint64 {
		return rand.Uint64()
	}
}

//NewTraceID return a random string
func NewTraceID() string {
	return ID{
		High: randomFuc(),
		Low:  randomFuc(),
	}.String()
}

//GetTraceIDFromHTTPHeader GetTraceIDFromHTTPHeader
func GetTraceIDFromHTTPHeader(header http.Header) string {
	traceID := header.Get(string(TraceKeyName))
	return traceID
}

//GetTraceIDFromContext GetTraceIDFromContext
func GetTraceIDFromContext(ctx context.Context) string {
	traceID := ctx.Value(TraceKeyName)
	if traceID == nil {
		return ""
	}
	return traceID.(string)
}

//GetTraceIDFromGrpcMetadata GetTraceIDFromGrpcMetadata case insensitive
func GetTraceIDFromGrpcMetadata(md metadata.MD) string {
	if t, ok := md[string(TraceKeyName)]; ok {
		if len(t) > 0 {
			return t[0]
		}
	}

	if t, ok := md[strings.ToLower(string(TraceKeyName))]; ok {
		if len(t) > 0 {
			return t[0]
		}
	}
	return ""
}

//NewCtxWithTraceID NewCtxWithTraceID
func NewCtxWithTraceID(ctx context.Context) context.Context {
	traceID := GetTraceIDFromContext(ctx)

	if traceID != "" {
		return ctx
	}

	traceID = NewTraceID()
	return context.WithValue(ctx, TraceKeyName, traceID)
}
