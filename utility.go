package hosts

import (
	"bytes"
	"net"
	"unicode"
	"unicode/utf8"
)

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func isAsciiSpace(x byte) bool {
	return asciiSpace[x] == 1
}

func indexAsciiSpace(p []byte) int {
	for i, x := range p {
		// Fallback to bytes.IndexFunc, when input slice contains
		// multibyte characters.
		if x >= utf8.RuneSelf {
			return bytes.IndexFunc(p[i:], unicode.IsSpace)
		}
		if isAsciiSpace(x) {
			return i
		}
	}
	return -1
}

func fields(p []byte, f func(i int, p []byte) bool) {
	for i := 0; ; i++ {
		idx := indexAsciiSpace(p)
		if idx == -1 {
			f(i, p)
			return
		}

		if idx > 0 {
			if !f(i, p[:idx]) {
				return
			}
		}

		p = p[idx+1:]
	}
}

func toLower(xs []byte) []byte {
	for i, x := range xs {
		// Fallback to bytes.Map, when input slice contains
		// multibyte characters.
		if xs[i] >= utf8.RuneSelf {
			return bytes.Map(unicode.ToLower, xs)
		}
		if x >= 'A' && x <= 'Z' {
			xs[i] = 'a' + (x - 'A')
		}
	}
	return xs
}

func isZeros(xs []byte) bool {
	for _, x := range xs {
		if x != 0 {
			return false
		}
	}
	return true
}

func isIPv4(ip net.IP) bool {
	return len(ip) == net.IPv4len || isZeros(ip[0:10]) && ip[10] == 0xff && ip[11] == 0xff
}

func isIPv6(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return !isIPv4(ip)
}
