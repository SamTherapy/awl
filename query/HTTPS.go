// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
)

// HTTPSResolver is for DNS-over-HTTPS queries.
type HTTPSResolver struct {
	opts util.Options
}

var _ Resolver = (*HTTPSResolver)(nil)

// LookUp performs a DNS query.
func (r *HTTPSResolver) LookUp(msg *dns.Msg) (util.Response, error) {
	var resp util.Response

	httpR := &http.Client{
		Timeout: r.opts.Request.Timeout,
	}

	buf, err := msg.Pack()
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: packing: %w", err)
	}

	r.opts.Logger.Debug("https: sending HTTPS request")

	req, err := http.NewRequest("POST", r.opts.Request.Server, bytes.NewBuffer(buf))
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: request creation: %w", err)
	}

	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	now := time.Now()
	res, err := httpR.Do(req)
	resp.RTT = time.Since(now)

	if err != nil {
		return util.Response{}, fmt.Errorf("doh: HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return util.Response{}, &ErrHTTPStatus{res.StatusCode}
	}

	r.opts.Logger.Debug("https: reading response")

	fullRes, err := io.ReadAll(res.Body)
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: body read: %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return util.Response{}, fmt.Errorf("doh: body close: %w", err)
	}

	r.opts.Logger.Debug("https: unpacking response")

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
