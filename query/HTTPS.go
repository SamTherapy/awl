// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"github.com/miekg/dns"
)

type HTTPSResolver struct {
	opts cli.Options
}

func (r *HTTPSResolver) LookUp(msg *dns.Msg) (helpers.Response, error) {
	var resp helpers.Response
	httpR := &http.Client{
		Timeout: r.opts.Request.Timeout,
	}
	buf, err := msg.Pack()
	if err != nil {
		return helpers.Response{}, err
	}
	r.opts.Logger.Debug("making DoH request")
	// query := server + "?dns=" + base64.RawURLEncoding.EncodeToString(buf)
	req, err := http.NewRequest("POST", r.opts.Request.Server, bytes.NewBuffer(buf))
	if err != nil {
		return helpers.Response{}, fmt.Errorf("DoH: %w", err)
	}
	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	now := time.Now()
	res, err := httpR.Do(req)
	resp.RTT = time.Since(now)

	if err != nil {
		return helpers.Response{}, fmt.Errorf("DoH HTTP request error: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return helpers.Response{}, fmt.Errorf("DoH server responded with HTTP %d", res.StatusCode)
	}

	fullRes, err := io.ReadAll(res.Body)
	if err != nil {
		return helpers.Response{}, fmt.Errorf("DoH body read error: %w", err)
	}
	resp.DNS = &dns.Msg{}
	r.opts.Logger.Debug("unpacking response")
	err = resp.DNS.Unpack(fullRes)
	if err != nil {
		return helpers.Response{}, fmt.Errorf("DoH dns message unpack error: %w", err)
	}

	return resp, nil
}
