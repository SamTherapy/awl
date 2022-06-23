package query

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

// Resolve a DNS-over-HTTPS query
//
// Currently only supports POST requests
func ResolveHTTPS(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
	httpR := &http.Client{}
	buf, err := msg.Pack()
	if err != nil {
		return nil, 0, err
	}
	// query := server + "?dns=" + base64.RawURLEncoding.EncodeToString(buf)
	req, err := http.NewRequest("POST", server, bytes.NewBuffer(buf))
	if err != nil {
		return nil, 0, fmt.Errorf("DoH: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	now := time.Now()
	res, err := httpR.Do(req)
	rtt := time.Since(now)

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
	response := dns.Msg{}
	err = response.Unpack(fullRes)
	if err != nil {
		return nil, 0, fmt.Errorf("DoH dns message unpack error: %s", err.Error())
	}

	return &response, rtt, nil
}
