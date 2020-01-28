package hosts

import (
	"bufio"
	"bytes"
	"io"
	"net"
)

func parseIP(c map[string]List, p []byte) List {
	if x, ok := c[string(p)]; ok {
		return x
	}

	ip := net.ParseIP(string(p))
	if ip != nil {
		x := List{ip}
		c[string(p)] = x
		return x
	}
	return nil
}

// ParseSize parses multiple host files contents provided by rs and returns
// map of domain to IP address associations, or error in case of first
// unsuccessful read operation.
// The size parameter is a capacity hint for the returned map.
//
// Domains are treated as case-insensitive and converted to lower case.
// IP addresses occur in the same order as in hosts file.
//
// Parse can parse normal host file, as well as domain lists. In the second
// case, the domain is associated with nil IP address list.
//
// There is an optimization for a space efficency: different domains with
// the same singular address entry will share the same address list.
// It is useful for domain blocklists, where domains are mapped to localhost
// address.
func ParseSize(size int, rs ...io.Reader) (map[string]List, error) {
	cache := map[string]List{
		"0.0.0.0":   {net.IPv4(0, 0, 0, 0)},
		"127.0.0.1": {net.IPv4(127, 0, 0, 1)},
		"::":        {net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		"::1":       {net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
	}

	ret := make(map[string]List, size)
	for _, r := range rs {
		s := bufio.NewScanner(r)
		for s.Scan() {
			l := bytes.TrimSpace(s.Bytes())
			if len(l) == 0 || l[0] == '#' {
				continue
			}

			var ip []net.IP
			fields(l, func(i int, p []byte) bool {
				p = bytes.TrimSpace(p)

				if len(p) > 0 {
					if i == 0 {
						ip = parseIP(cache, p)
						if ip != nil {
							return true
						}
					}

					if p[0] == '#' {
						return false
					}

					p = toLower(p)

					if xs, ok := ret[string(p)]; ok {
						if len(ip) > 0 && !xs.ContainsIP(ip[0]) {
							if len(xs) == 1 {
								xs = xs.Clone()
							}
							xs = append(xs, ip...)
							ret[string(p)] = xs
						}
					} else {
						ret[string(p)] = ip
					}
				}
				return true
			})
		}
		if err := s.Err(); err != nil {
			return ret, err
		}
	}
	return ret, nil
}

// Parse parses multiple host files contents provided by rs and returns
// map of domain to IP address associations, or error in case of first
// unsuccessful read operation.
//
// See ParseSize for more information.
func Parse(rs ...io.Reader) (map[string]List, error) {
	return ParseSize(0, rs...)
}
