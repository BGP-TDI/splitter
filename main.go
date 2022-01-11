package main

import (
	"flag"
	"fmt"
	"log"
	"net/netip"
)

var (
	opt_L bool
	opt_R bool
)

func split(ipbytes [16]byte, plen int, left bool) {
	// print
	ipaddr := netip.AddrFrom16(ipbytes)
	ipprefix, err := ipaddr.Prefix(plen)
	if err != nil { log.Fatal(err) }

	switch {
	case opt_L:
		if left  { fmt.Println(ipprefix) }
	case opt_R:
		if !left { fmt.Println(ipprefix) }
	default:
		fmt.Println(ipprefix)
	}

	if plen < 48 {
		i := plen / 8
		j := 7 - (plen % 8)

		// set to 0
		ipbytes[i] &= ^(1 << j)
		split(ipbytes, plen+1, true)

		// set to 1
		ipbytes[i] |= (1 << j)
		split(ipbytes, plen+1, false)
	}
}

func main() {
	flag.BoolVar(&opt_L, "L", false, "print only the prefixes after a split on 0 (left-hand side)")
	flag.BoolVar(&opt_R, "R", false, "print only the prefixes after a split on 1 (right-hand side)")
	flag.Parse()

	arg1 := flag.Arg(0)
	if len(arg1) == 0 { log.Fatal("pass IPv6 prefix as the 1st arg") }

	ipprefix, err := netip.ParsePrefix(arg1)
	if err != nil { log.Fatalf("parsing IP prefix in the 1st arg: %s", err) }

	ipprefix = ipprefix.Masked()
	ipaddr := ipprefix.Addr().Unmap()
	if !ipaddr.Is6() { log.Fatalf("%s: not an IPv6 prefix", ipprefix) }

	ipbytes := ipaddr.As16()
	plen := ipprefix.Bits()
	// fmt.Printf("parsed: %s\nraw: %x / %d\n", ipprefix, ipbytes, plen)
	split(ipbytes, plen, true)
}
