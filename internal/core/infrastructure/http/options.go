package http

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

const (
	Timeout               = 2 * time.Second
	DialContextTimeOut    = 30 * time.Second
	DialContextKeepAlive  = 30 * time.Second
	MaxIdleConnections    = 100
	IdleConnTimeout       = 90 * time.Second
	TLSHandshakeTimeout   = 10 * time.Second
	ExpectContinueTimeout = 1 * time.Second
)

// Option is a functional option type that allows us to configure the Client.
type Option func(*PageParser)

// WithHTTPClient is a functional option to set the HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(provider *PageParser) {
		provider.httpClient = httpClient
	}
}

func WithHTTPTimeOut(timeOut int) Option {
	return func(provider *PageParser) {
		provider.httpClient.Timeout = time.Duration(timeOut) * time.Second
	}
}

func getDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: Timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   DialContextTimeOut,
				KeepAlive: DialContextKeepAlive,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          MaxIdleConnections,
			IdleConnTimeout:       IdleConnTimeout,
			TLSHandshakeTimeout:   TLSHandshakeTimeout,
			ExpectContinueTimeout: ExpectContinueTimeout,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
}
