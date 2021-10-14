package main

import (
	"flag"
	"fmt"
	"log"

	"inet.af/netaddr"
)

func split(ipbytes [16]byte, plen uint8) {
	// print
	ipaddr := netaddr.IPFrom16(ipbytes)
	ipprefix, err := ipaddr.Prefix(plen)
	if err != nil { log.Fatal(err) }
	fmt.Println(ipprefix)

	if plen < 48 {
		i := plen / 8
		j := 7 - (plen % 8)

		// set to 0
		ipbytes[i] &= ^(1 << j)
		split(ipbytes, plen+1)

		// set to 1
		ipbytes[i] |= (1 << j)
		split(ipbytes, plen+1)
	}
}

func main() {
	flag.Parse()

	arg1 := flag.Arg(0)
	if len(arg1) == 0 { log.Fatal("pass IPv6 prefix as the 1st arg") }

	ipprefix, err := netaddr.ParseIPPrefix(arg1)
	if err != nil { log.Fatalf("parsing IP prefix in the 1st arg: %s", err) }

	ipprefix = ipprefix.Masked()
	ipaddr := ipprefix.IP().Unmap()
	if !ipaddr.Is6() { log.Fatalf("%s: not an IPv6 prefix", ipprefix) }

	ipbytes := ipaddr.As16()
	plen := ipprefix.Bits()
	// fmt.Printf("parsed: %s\nraw: %x / %d\n", ipprefix, ipbytes, plen)
	split(ipbytes, plen)
}
