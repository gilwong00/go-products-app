package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type Compression struct {
}

func (c *Compression) CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// gzip response
			writer := NewGzipResponseWriter(w)
			writer.Header().Set("Content-Encoding", "gzip")
			/*
				call next with our gzip writer so upstream func will automatically
				apply gzip compression when it writes its data
			*/
			next.ServeHTTP(writer, r)
			defer writer.Flush()
			return
		}
		// normal processing
		next.ServeHTTP(w, r)
	})
}

type GZipResponseWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

func NewGzipResponseWriter(w http.ResponseWriter) *GZipResponseWriter {
	/*
		gzip.NewWriter takes in an io.Writer
		http.ResponseWriter implements an io.Writer so we can use the passed in writer
	*/
	gw := gzip.NewWriter(w)
	return &GZipResponseWriter{w: w, gw: gw}
}

func (grw *GZipResponseWriter) Header() http.Header {
	return grw.w.Header()
}

// when we call Write, we wanna use the gzip writer to compress
func (grw *GZipResponseWriter) Write(d []byte) (int, error) {
	return grw.gw.Write(d)
}

func (grw *GZipResponseWriter) WriteHeader(statuscode int) {
	grw.w.WriteHeader(statuscode)
}

/*
Flushes anything that has not been sent in the stream
*/
func (grw *GZipResponseWriter) Flush() {
	grw.gw.Flush()
	grw.gw.Close()
}
