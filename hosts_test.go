package hosts_test

import (
	"bytes"
	"errors"
	"io"
	"net"
	"reflect"
	"testing"
	"unsafe"

	"github.com/ksinica/hosts"
)

func isSameIp(a, b []byte) bool {
	return bytes.Equal(a, b)
}

func TestSuccessfulHostsParse(t *testing.T) {
	const input = `
	127.0.0.1 domain.com	domain2.com    żółta.GĘŚ.pl     
	::1       domain.com       # comment
	1.1.1.1   one.one.one.one`

	hosts, err := hosts.Parse(bytes.NewBufferString(input))
	if err != nil {
		t.Error(err)
	}

	if len(hosts) != 4 {
		t.Fail()
	}

	if !isSameIp(hosts["domain.com"][0], net.ParseIP("127.0.0.1")) {
		t.Fail()
	}

	if !isSameIp(hosts["domain2.com"][0], net.ParseIP("127.0.0.1")) {
		t.Fail()
	}

	if !isSameIp(hosts["żółta.gęś.pl"][0], net.ParseIP("127.0.0.1")) {
		t.Fail()
	}

	if !isSameIp(hosts["domain.com"][1], net.ParseIP("::1")) {
		t.Fail()
	}

	if !isSameIp(hosts["one.one.one.one"][0], net.ParseIP("1.1.1.1")) {
		t.Fail()
	}
}

func TestSuccessfulDomainListParse(t *testing.T) {
	const input = `
	one.one.one.one
	golang.org
	Play.Golang.Org`

	hosts, err := hosts.Parse(bytes.NewBufferString(input))
	if err != nil {
		t.Error(err)
	}

	if len(hosts) != 3 {
		t.Fail()
	}

	if hosts["one.one.one.one"] != nil ||
		hosts["golang.org"] != nil ||
		hosts["play.golang.org"] != nil {
		t.Fail()
	}
}

func TestUnsuccessfulParse(t *testing.T) {
	var errTest = errors.New("test")

	pr, pw := io.Pipe()
	pw.CloseWithError(errTest)

	_, err := hosts.Parse(pr)
	if err != errTest {
		t.Fail()
	}
}

func dataPtr(s []net.IP) uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(&s)).Data
}

func TestParseAddressCache(t *testing.T) {
	const inputA = `
	127.0.0.1 hosta hostb
	127.0.0.1 hostc`

	const inputB = `
	127.0.0.1 hostd`

	hosts, err := hosts.Parse(
		bytes.NewBufferString(inputA),
		bytes.NewBufferString(inputB),
	)
	if err != nil {
		t.Fail()
	}

	if dataPtr(hosts["hosta"]) != dataPtr(hosts["hostb"]) ||
		dataPtr(hosts["hosta"]) != dataPtr(hosts["hostc"]) ||
		dataPtr(hosts["hosta"]) != dataPtr(hosts["hostd"]) {
		t.Fail()
	}
}

func TestParseMultipleReaders(t *testing.T) {
	const inputA = `
	127.0.0.1 hosta
	127.0.0.2 hostb`

	const inputB = `
	127.0.0.3 hostc
	127.0.0.4 hostc
	127.0.0.5 hosta`

	const inputC = `
	hostd`

	hosts, err := hosts.Parse(
		bytes.NewBufferString(inputA),
		bytes.NewBufferString(inputB),
		bytes.NewBufferString(inputC),
	)
	if err != nil {
		t.Fail()
	}

	if len(hosts) != 4 {
		t.Fail()
	}

	if !isSameIp(hosts["hosta"][0], net.ParseIP("127.0.0.1")) &&
		!isSameIp(hosts["hosta"][1], net.ParseIP("127.0.0.5")) {
		t.Fail()
	}

	if !isSameIp(hosts["hostb"][0], net.ParseIP("127.0.0.2")) {
		t.Fail()
	}

	if !isSameIp(hosts["hostc"][0], net.ParseIP("127.0.0.3")) &&
		!isSameIp(hosts["hostc"][1], net.ParseIP("127.0.0.4")) {
		t.Fail()
	}

	if hosts["hostd"] != nil {
		t.Fail()
	}
}

const (
	benchData = `
0.0.0.0 livingthenourishedlife.com
0.0.0.0 shanty-2-chic.com
0.0.0.0 dietnavi.com
0.0.0.0 lifeline.de
0.0.0.0 sc-s.com
0.0.0.0 paginagospel.com.br
0.0.0.0 webhostingtalk.lk
0.0.0.0 prodieta.ro
0.0.0.0 datapremiery.pl
0.0.0.0 donaperfeitinha.com
0.0.0.0 mountainhardwear.com
0.0.0.0 jinxykids.com
0.0.0.0 a-to-z-of-manners-and-etiquette.com
0.0.0.0 epubdump.com
0.0.0.0 songbirdgarden.com
0.0.0.0 beliefnet.com
0.0.0.0 bling99.com
0.0.0.0 playsport.cc
0.0.0.0 thetaylor-house.com
0.0.0.0 kyxaodienanh.com
0.0.0.0 2dishingdivas.com
0.0.0.0 tamilsonglyrics.org
0.0.0.0 dailymasalla.com
0.0.0.0 securitykiller.org
0.0.0.0 prostoblog.com.ua
0.0.0.0 doguhakimiyet.com
0.0.0.0 chocolateamais.com
0.0.0.0 maslatip.com
0.0.0.0 woman.com.au
0.0.0.0 jornaldascaldas.com
0.0.0.0 poradum.com.ua
0.0.0.0 acselkita.com
0.0.0.0 currentincarmel.com
0.0.0.0 paranormal360.co.uk
0.0.0.0 ioannina24.gr
0.0.0.0 myboredtoddler.com
0.0.0.0 spotry.me
0.0.0.0 novaperspectiva.com
0.0.0.0 escritosdederecho.com
0.0.0.0 diybastelideen.com
0.0.0.0 napolimagazine.info
0.0.0.0 olabloga.pl
0.0.0.0 doctoradriancormillot.com
0.0.0.0 mjna50.net
0.0.0.0 modasemcensura.com
0.0.0.0 orchidsforum.com
0.0.0.0 aquamaniya.ru
0.0.0.0 teckknow.com
0.0.0.0 hnbmg.com
0.0.0.0 acleanplate.com
0.0.0.0 ladymilonguera.fr
0.0.0.0 wdwinfo.com
0.0.0.0 guesstheemoji-answers.com
0.0.0.0 mariara.info
0.0.0.0 cactushugs.com
0.0.0.0 wyborcza.pl
0.0.0.0 sobesednik.ru
0.0.0.0 banker.bg
0.0.0.0 funcheaporfree.com
0.0.0.0 halosheaven.com
0.0.0.0 hitmovie.co
0.0.0.0 seniorhousingnews.com
0.0.0.0 bergamonews.it
0.0.0.0 tigerland.com
0.0.0.0 freebeacon.com
0.0.0.0 xatakafoto.com
0.0.0.0 chipandco.com
0.0.0.0 spankbang.com
0.0.0.0 villagesclubsdusoleil.com
0.0.0.0 kangolstore.com
0.0.0.0 laptopbatteryexpress.com
0.0.0.0 sectornolimits.com
0.0.0.0 deinterieurcollectie.nl
0.0.0.0 kimonomodern.com
0.0.0.0 marionnaud.it
0.0.0.0 skin1.com
0.0.0.0 elcorteingles.eu
0.0.0.0 kiwoko.com
0.0.0.0 cisalfasport.it
0.0.0.0 medicbatteries.com
0.0.0.0 heels.com
0.0.0.0 1-800homeopathy.com
0.0.0.0 badkamermarkt.nl
0.0.0.0 badkamerconcurrent.nl
0.0.0.0 utahskis.com
0.0.0.0 bigs.jp
0.0.0.0 coccinelle.com
0.0.0.0 mitiendadearte.com
0.0.0.0 metroshoes.net
0.0.0.0 surlatable.com
0.0.0.0 pixibeauty.com
0.0.0.0 trutechtools.com
0.0.0.0 dvd.it
0.0.0.0 wysada.com
0.0.0.0 promofarma.com
0.0.0.0 laredoute.es
0.0.0.0 oldpueblotraders.com
0.0.0.0 thomascook.in
0.0.0.0 thewarehouse.co.nz
0.0.0.0 fairyglen.com`
)

func BenchmarkParse100(b *testing.B) {
	var buf bytes.Buffer
	for n := 0; n < b.N; n++ {
		buf.Reset()
		buf.WriteString(benchData)
		hosts.ParseSize(100, &buf)
	}
}
