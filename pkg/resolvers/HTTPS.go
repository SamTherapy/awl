// SPDX-License-Identifier: BSD-3-Clause

package resolvers

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
)

// HTTPSResolver is for DNS-over-HTTPS queries.
type HTTPSResolver struct {
	client http.Client
	opts   util.Options
}

var _ Resolver = (*HTTPSResolver)(nil)

// LookUp performs a DNS query.
func (resolver *HTTPSResolver) LookUp(msg *dns.Msg) (util.Response, error) {
	var resp util.Response

	resolver.client = http.Client{
		Timeout: resolver.opts.Request.Timeout,
		Transport: &http.Transport{
			MaxConnsPerHost:     1,
			MaxIdleConns:        1,
			MaxIdleConnsPerHost: 1,
			Proxy:               http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				//nolint:gosec // This is intentional if the user requests it
				InsecureSkipVerify: resolver.opts.TLSNoVerify,
				ServerName:         resolver.opts.TLSHost,
			},
		},
	}

	buf, err := msg.Pack()
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: packing: %w", err)
	}

	resolver.opts.Logger.Debug("https: sending HTTPS request")

	req, err := http.NewRequest("POST", resolver.opts.Request.Server, bytes.NewBuffer(buf))
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: request creation: %w", err)
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	now := time.Now()
	res, err := resolver.client.Do(req)
	resp.RTT = time.Since(now)

	if err != nil {
		return util.Response{}, fmt.Errorf("doh: HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return util.Response{}, &ErrHTTPStatus{res.StatusCode}
	}

	resolver.opts.Logger.Debug("https: reading response")

	fullRes, err := io.ReadAll(res.Body)
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: body read: %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: body close: %w", err)
	}

	resolver.opts.Logger.Debug("https: unpacking response")

	resp.DNS = &dns.Msg{}

	err = resp.DNS.Unpack(fullRes)
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: dns message unpack: %w", err)
	}

	return resp, nil
}

// ErrHTTPStatus is returned when DoH returns a bad status code.
type ErrHTTPStatus struct {
	code int
}

func (e *ErrHTTPStatus) Error() string {
	return fmt.Sprintf("doh server responded with HTTP %d", e.code)
}
