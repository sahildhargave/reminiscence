package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sahildhargave/memories/account/model/apperrors"
)

// Timeout wraps the request with a timeout.
func Timeout(timeout time.Duration, errTimeout *apperrors.Error) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Setting Gin's writer as custom writer
		tw := &timeoutWriter{ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

		// Wrap the request with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Update Gin request context
		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})
		panicChan := make(chan interface{})

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next() // Call subsequent middleware(s) and handler
			finished <- struct{}{}
		}()

		select {
		case <-panicChan:
			// If cannot recover from panic, send internal server error
			e := apperrors.NewInternal()
			tw.ResponseWriter.WriteHeader(e.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": e,
			})
			tw.ResponseWriter.Write(eResp)
		case <-finished:
			// If finished, set headers and write response
			tw.mu.Lock()
			defer tw.mu.Unlock()
			// Copy headers from tw.Header() (written to by Gin)
			// to tw.ResponseWriter for response
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.h {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			tw.ResponseWriter.Write(tw.wbuf.Bytes())
		case <-ctx.Done():
			// Timeout has occurred, send errTimeout and write headers
			tw.mu.Lock()
			defer tw.mu.Unlock()
			// Set response headers
			tw.ResponseWriter.Header().Set("Content-Type", "application/json")
			tw.ResponseWriter.WriteHeader(errTimeout.Status())
			eResp, _ := json.Marshal(gin.H{
				"error": errTimeout,
			})
			tw.ResponseWriter.Write(eResp)
			c.Abort()
			tw.SetTimedOut()
		}
	}
}

// timeoutWriter implements http.ResponseWriter and tracks if the writer has timed out.
type timeoutWriter struct {
	gin.ResponseWriter
	h          http.Header
	wbuf       bytes.Buffer
	mu         sync.Mutex
	timedOut   bool
	wroteHeader bool
	code       int
}

// Write writes the response, but only if the writer hasn't already timed out.
func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, nil
	}
	return tw.wbuf.Write(b)
}

// WriteHeader writes the response header, but only if the writer hasn't already timed out.
func (tw *timeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

// writeHeader sets the response header.
func (tw *timeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

// Header returns the header map.
func (tw *timeoutWriter) Header() http.Header {
	return tw.h
}

// SetTimedOut sets the timedOut flag to true.
func (tw *timeoutWriter) SetTimedOut() {
	tw.timedOut = true
}

// checkWriteHeaderCode checks if the provided status code is valid.
func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
