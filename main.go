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

func split(ipbytes [16]byte, plen int) {
	if plen >= 48 {
		return
	}

	i := plen / 8
	j := 7 - (plen % 8)

	// set to 0
	ipbL := ipbytes
	ipbL[i] &= ^(1 << j)
	prefixL, err := netip.AddrFrom16(ipbL).Prefix(plen+1)
	if err != nil {
		log.Fatal(err)
	} else if opt_L || !opt_R {
		fmt.Println(prefixL)
	}

	// set to 1
	ipbR := ipbytes
	ipbR[i] |= (1 << j)
	prefixR, err := netip.AddrFrom16(ipbR).Prefix(plen+1)
	if err != nil {
		log.Fatal(err)
	} else if !opt_L || opt_R {
		fmt.Println(prefixR)
	}

	// recurse
	split(ipbL, plen+1)
	split(ipbR, plen+1)
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
	split(ipbytes, plen)
}
