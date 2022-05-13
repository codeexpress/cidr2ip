package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
)

const (
	Version = "1.0.0"
	IPRegex = `\b(?:\d{1,3}\.){3}\d{1,3}\b$`
)

var (
	cidrFilePtr = flag.String("f", "",
		"[Optional] Name of file with CIDR blocks")
	printRangesPtrPtr = flag.Bool("r", false,
		"[Optional] Print IP ranges instead of all IPs")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	args := os.Args[1:]

	if *cidrFilePtr != "" {
		file, err := os.Open(*cidrFilePtr)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			displayIPs(scanner.Text())
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
	} else if info.Mode()&os.ModeNamedPipe != 0 { // data is piped in
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			displayIPs(scanner.Text())
		}
	} else if len(args) > 0 { // look for CIDRs on cmd line
		var cidrs []string
		if *printRangesPtrPtr == true {
			cidrs = args[1:]
		} else {
			cidrs = args
		}

		for _, cidr := range cidrs {
			displayIPs(cidr)
		}
	} else { // no piped input, no file provide and no args, display usage
		flag.Usage()
	}
}

func isIPAddr(cidr string) bool {
	match, _ := regexp.MatchString(IPRegex, cidr)
	return match
}

func displayIPs(cidr string) {
	var ips []string

	// if a IP address, display the IP address and return
	if isIPAddr(cidr) {
		fmt.Println(cidr)
		return
	}

	ipAddr, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Print(err)
		return
	}

	for ip := ipAddr.Mask(ipNet.Mask); ipNet.Contains(ip); increment(ip) {
		ips = append(ips, ip.String())
	}

	// CIDR too small eg. /31
	if len(ips) <= 2 {
		return
	}

	if *printRangesPtrPtr == true {
		fmt.Printf("%s-%s\n", ips[1], ips[len(ips)-1])
	} else {
		for _, ip := range ips[1 : len(ips)-1] {
			fmt.Println(ip)
		}
	}
}

// The next IP address of a given ip address
// https://stackoverflow.com/a/33925954
func increment(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "CIDR to IPs version %s\n", Version)
	fmt.Fprintf(os.Stderr, "Usage:   $ cidr2ip [-r] [-f <filename>] <list of cidrs> \n")
	fmt.Fprintf(os.Stderr, "Example: $ cidr2ip -f cidrs.txt\n")
	fmt.Fprintf(os.Stderr, "         $ cidr2ip 10.0.0.0/24\n")
	fmt.Fprintf(os.Stderr, "         $ cidr2ip -r 10.0.0.0/24\n")
	fmt.Fprintf(os.Stderr, "         $ cidr2ip -r -f cidrs.txt\n")
	fmt.Fprintf(os.Stderr, "         $ cat cidrs.txt | cidr2ip \n")
	fmt.Fprintf(os.Stderr, "--------------------------\nFlags:\n")
	flag.PrintDefaults()
}
