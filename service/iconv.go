package service

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	zh "golang.org/x/text/encoding/simplifiedchinese"
)

type iconvHTTPBodyWriter struct {
	http.ResponseWriter
	charset string
	encoder io.Writer
}

// Write writes bytes to encoder, which in turn writes bytes to http.ResponseWriter
func (w *iconvHTTPBodyWriter) Write(p []byte) (int, error) {
	if len(w.charset) > 0 {
		w.WriteHeader(http.StatusOK)
		w.charset = ""
	}
	return w.encoder.Write(p)
}

// WriteHeader sets correct charset to the Content-Type of response
func (w *iconvHTTPBodyWriter) WriteHeader(code int) {
	// parse response Content-Type for mime
	mediaType, param, err := mime.ParseMediaType(w.Header().Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	// set charset
	param["charset"] = w.charset
	w.charset = ""
	// overwrite with the new charset
	w.Header().Set("Content-Type", mime.FormatMediaType(mediaType, param))
	// perform real write header
	w.ResponseWriter.WriteHeader(code)
}

// IconvHandler struct contains the ServeHTTP method
type IconvHandler struct{}

// NewIconvHandler returns a handler which will convert the request body from
// the encoding in the request to utf-8, and convert the response body back to
// that encoding in ServeHTTP.
func NewIconvHandler() *IconvHandler {
	return &IconvHandler{}
}

// ServeHTTP wraps the http.ResponseWriter with a gzip.Writer.
func (h *IconvHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// get Content-Type
	contentType := r.Header["Content-Type"]
	// if there is no Content-Type, just skip
	if len(contentType) == 0 {
		next(w, r)
		return
	}

	// parse Content-Type for charset
	_, param, err1 := mime.ParseMediaType(contentType[len(contentType)-1])
	if err1 != nil {
		http.Error(w, "Unrecognized Content Type", http.StatusUnprocessableEntity)
		return
	}
	charset := strings.ToLower(param["charset"])
	if len(charset) == 0 {
		charset = "utf-8"
	}

	// get decoder and encoder by charset
	// raw io.Reader & io.Writer
	var rd io.Reader
	var wt io.Writer
	switch charset {
	case "utf-8":
		rd = r.Body
		wt = w
	case "gb2312":
		fallthrough
	case "gbk":
		rd = zh.GBK.NewDecoder().Reader(r.Body)
		wt = zh.GBK.NewEncoder().Writer(w)
	default:
		http.Error(w, "Unsupported Encoding", http.StatusUnprocessableEntity)
		return
	}

	// wraps raw io.Reader to ReadCloser
	r.Body = ioutil.NopCloser(rd)

	// wraps raw io.Writer to iconvHTTPBodyWriter
	proxyWriter := &iconvHTTPBodyWriter{
		ResponseWriter: w,
		charset:        charset,
		encoder:        wt,
	}
	next(proxyWriter, r)
}
