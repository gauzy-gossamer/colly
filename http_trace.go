package colly

import (
	"net/http"
	"net/http/httptrace"
	"crypto/tls"
	"time"
)

// HTTPTrace provides a datastructure for storing an http trace.
type HTTPTrace struct {
	start, connect    time.Time
	ConnectDuration   time.Duration
	FirstByteDuration time.Duration
	DNSInfo           httptrace.DNSDoneInfo
	TLSConn           tls.ConnectionState
}

// trace returns a httptrace.ClientTrace object to be used with an http
// request via httptrace.WithClientTrace() that fills in the HttpTrace.
func (ht *HTTPTrace) trace() *httptrace.ClientTrace {
	trace := &httptrace.ClientTrace{
		ConnectStart: func(network, addr string) { ht.connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			ht.ConnectDuration = time.Since(ht.connect)
		},

		DNSDone: func(dns_info httptrace.DNSDoneInfo) {
			ht.DNSInfo = dns_info
		},
		TLSHandshakeDone: func(tls_conn tls.ConnectionState, err error) {
			ht.TLSConn = tls_conn
		},
		GetConn: func(hostPort string) { ht.start = time.Now() },
		GotFirstResponseByte: func() {
			ht.FirstByteDuration = time.Since(ht.start)
		},
	}
	return trace
}

// WithTrace returns the given HTTP Request with this HTTPTrace added to its
// context.
func (ht *HTTPTrace) WithTrace(req *http.Request) *http.Request {
	return req.WithContext(httptrace.WithClientTrace(req.Context(), ht.trace()))
}
