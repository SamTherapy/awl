// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

type HTTPSResolver struct {
	server string
	opts   Options
}

func (r *HTTPSResolver) LookUp(msg *dns.Msg) (*dns.Msg, time.Duration, error) {
	var resp Response
	httpR := &http.Client{}
	buf, err := msg.Pack()
	if err != nil {
		return nil, 0, err
	}
	r.opts.Logger.Debug("making DoH request")
	// query := server + "?dns=" + base64.RawURLEncoding.EncodeToString(buf)
	req, err := http.NewRequest("POST", r.server, bytes.NewBuffer(buf))
	if err != nil {
		return nil, 0, fmt.Errorf("DoH: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	now := time.Now()
	res, err := httpR.Do(req)
	resp.Answers.RTT = time.Since(now)

	if err != nil {
		return nil, 0, fmt.Errorf("DoH HTTP request error: %s", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("DoH server responded with HTTP %d", res.StatusCode)
	}

	fullRes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("DoH body read error: %s", err.Error())
	}
	resp.DNS = dns.Msg{}
	r.opts.Logger.Debug("unpacking response")
	err = resp.DNS.Unpack(fullRes)
	if err != nil {
		return nil, 0, fmt.Errorf("DoH dns message unpack error: %s", err.Error())
	}

	return &resp.DNS, resp.Answers.RTT, nil
}
