package middle

import (
	"bytes"
	"io"
	"net/http"
)

// RespSniffer implements http.ResponseWriter while also allowing
// introspection into what's being returned to the client. Useful
// in middleware applications such as logging or inline modification
// of Write calls.
type RespSniffer struct {
	orig   http.ResponseWriter
	out    io.Writer
	Status int
	Size   int
}

// NewRespSniffer returns a new RespSniffer middleware. It takes the original
// http.ResponseWriter object along with a new io.Writer that will be used
// instead of the original. This can be used in conjunction with io.MultiWriter
// and bytes.Buffer to make a simple response data copier or with a gzip.Writer
// to make an inline data compressor
func NewRespSniffer(orig http.ResponseWriter, out io.Writer) *RespSniffer {
	return &RespSniffer{orig: orig, Status: 200, out: out}
}

// Header simply returns the original ResponseWriter's Header
func (l *RespSniffer) Header() http.Header {
	return l.orig.Header()
}

// WriteHeader takes note of i in the RespSniffer.Status field while also
// passing it along to the original ResponseWriter
func (l *RespSniffer) WriteHeader(i int) {
	l.Status = i
	l.orig.WriteHeader(i)
}

// Write redirects data to the io.Writer used in the RespSniffer initialization.
// It also keeps track of the number of bytes written in the RespSniffer.Size
// field.
// It does not write to the original ResponseWriter
func (l *RespSniffer) Write(b []byte) (int, error) {
	n, err := l.out.Write(b)
	l.Size += n
	return n, err
}

// Used for copying the body of an http.Request. Must have a paired reqSniffer
type reqSniffer struct {
	body io.Closer
	r    io.Reader
	buf  *bytes.Buffer
}

// Read from either the buffer if it contains data or from the reader,
// which should copy to the paired reqSniffer's buffer
func (rc reqSniffer) Read(b []byte) (int, error) {
	if rc.buf.Len() > 0 {
		return rc.buf.Read(b)
	}
	return rc.r.Read(b)
}

// For properly implementing the ReadCloser interface
func (rc reqSniffer) Close() error {
	return rc.body.Close()
}

// make a pair of ReadClosers that fill the opposite buffer when they read,
// allowing both to read the same data
func makeReqSniffPair(body io.ReadCloser) (reqSniffer, reqSniffer) {
	var rightB, leftB bytes.Buffer
	rightR, leftR := io.TeeReader(body, &leftB), io.TeeReader(body, &rightB)
	return reqSniffer{body: body, r: leftR, buf: &leftB},
		reqSniffer{body: body, r: rightR, buf: &rightB}
}

// NewReqSniffer installes a sniffer in the http.Request and returns it
// as an io.Reader. Reads from this sniffer will not interfere with reads
// from the original Request Body.
func NewReqSniffer(orig *http.Request) io.Reader {
	pass, ret := makeReqSniffPair(orig.Body)
	orig.Body = pass
	return ret
}
