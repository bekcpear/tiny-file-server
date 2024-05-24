package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// LoggingResponseWriter is a custom http.ResponseWriter that captures the status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewLoggingResponseWriter creates a new LoggingResponseWriter
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

// WriteHeader captures the status code and writes the header
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Logger wraps an http.Handler and logs each request
func Logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := NewLoggingResponseWriter(w)
		inner.ServeHTTP(lrw, r)

		log.Printf(
			"[%d] %s %s %s %s",
			lrw.statusCode,
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

func main() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	addr := flag.String("addr", "", "address to serve on")
	port := flag.Int("port", 8000, "port to serve on")
	dir := flag.String("dir", wd, "the directory to serve on")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatal(err)
	}

	// Serve the current directory
	fs := http.FileServer(http.Dir(absDir))
	http.Handle("/", Logger(fs))

	// Start the server
	log.Printf("Serving %s on HTTP %s:%d\n", absDir, *addr, *port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", *addr, *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
