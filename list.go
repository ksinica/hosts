package hosts

import (
	"bytes"
	"net"
)

type List []net.IP

// Clone returns deep copy of the list.
func (l List) Clone() List {
	return append(l[:0:0], l...)
}

func (l List) filterIP(f func(net.IP) bool, g func(net.IP) bool) {
	for _, x := range l {
		if x != nil {
			if f(x) {
				if g(x) == false {
					return
				}
			}
		}
	}
}

// FilterIPv4 takes IP address list ips and calls f for each list
// element, that is an IPv4 address.
//
// If f returns false, the loop is terminated.
func (l List) FilterIPv4(f func(net.IP) bool) {
	l.filterIP(isIPv4, f)
}

// FilterIPv6 takes IP address list ips and calls f for each list
// element, that is an IPv6 address.
//
// If f returns false, the loop is terminated.
func (l List) FilterIPv6(f func(net.IP) bool) {
	l.filterIP(isIPv6, f)
}

// ContainsIP returns true when IP address provided in ip is present
// in the list, false otherwise.
func (l List) ContainsIP(ip net.IP) bool {
	for _, x := range l {
		if bytes.Equal(x, ip) {
			return true
		}
	}
	return false
}
