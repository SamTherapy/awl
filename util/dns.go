package util

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/miekg/dns"
)

func ResolveHTTPS(msg *dns.Msg, server string) (*dns.Msg, error) {
	httpR := &http.Client{}
	buf, err := msg.Pack()
	if err != nil {
		return nil, err
	}
	query := server + "?dns=" + base64.RawURLEncoding.EncodeToString(buf)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/dns-message")

	res, err := httpR.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad HTTP Request: %d", res.StatusCode)
	}

	fullRes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := dns.Msg{}
	err = response.Unpack(fullRes)
	if err != nil {
		return nil, err
	}

	return &response, nil

}
