package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/c-robinson/iplib"
	"github.com/miekg/dns"
)

type Response struct {
	Answers Answers `json:"Response"` //
}

// The Answers struct is the basic structure of a DNS request
// to be returned to the user upon making a request
type Answers struct {
	Server  string        `json:"Server"`  // The server to make the DNS request from
	Request uint16        `json:"Request"` // The type of request
	Name    string        `json:"Name"`    // The domain name to make a DNS request for
	RTT     time.Duration `json:"RTT"`     // When AWL was ran
}

func ReverseDNS(dom string, q string) (string, error) {
	if q == "PTR" {
		if strings.Contains(dom, ".") {
			// It's an IPv4 address
			ip := net.ParseIP(dom)
			if ip != nil {
				return iplib.IP4ToARPA(ip), nil
			} else {
				return "", errors.New("error: Could not parse IPv4 address")
			}

		} else if strings.Contains(dom, ":") {
			// It's an IPv6 address
			ip := net.ParseIP(dom)
			if ip != nil {
				return iplib.IP6ToARPA(ip), nil
			} else {
				return "", errors.New("error: Could not parse IPv6 address")
			}
		}
	}
	return "", errors.New("error: -x flag given but no IP found")
}

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
