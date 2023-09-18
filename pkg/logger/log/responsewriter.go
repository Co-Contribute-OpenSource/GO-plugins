package log

import "net/http"

// ResponseWriterInterceptor is used to decorate http.ResponseWriter,
// so we can intercept the Write process
type ResponseWriterInterceptor struct {
	http.ResponseWriter
	data []byte
}

func (w *ResponseWriterInterceptor) Write(b []byte) (int, error) {
	w.data = b
	return w.ResponseWriter.Write(b)
}
