package hosts_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/ksinica/hosts"
)

func TestListFilter(t *testing.T) {
	var xs hosts.List

	xs.FilterIPv4(func(_ net.IP) bool {
		t.Fail()
		return true
	})

	xs.FilterIPv6(func(_ net.IP) bool {
		t.Fail()
		return true
	})

	var ips = []string{
		"1.1.1.1",
		"1.0.0.1",
		"2606:4700:4700::1111",
		"2606:4700:4700::1001",
	}

	for _, x := range ips {
		xs = append(xs, net.ParseIP(x))
	}

	var ipv4 hosts.List
	xs.FilterIPv4(func(x net.IP) bool {
		ipv4 = append(ipv4, x)
		return true
	})

	if !reflect.DeepEqual(ipv4, xs[:2]) {
		t.Fail()
	}

	var ipv6 hosts.List
	xs.FilterIPv6(func(x net.IP) bool {
		ipv6 = append(ipv6, x)
		return true
	})

	if !reflect.DeepEqual(ipv6, xs[2:4]) {
		t.Fail()
	}
}
