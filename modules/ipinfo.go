package modules

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// ipTrigger Specifies the trigger.
const ipTrigger string = "ip"

// IpModule is a module for ip info
type IpModule struct {
	colorBg string
	colorFg string

	iconWif string
	iconEth string

	Alignment string
	Interval  time.Duration
}

// Initialize from config file. TODO
func (ip *IpModule) Initialize() {
	ip.colorBg = "#ff282828"
	ip.colorFg = "#ffebdbb2"

	ip.iconWif = ""
	ip.iconEth = ""

	if ip.Interval == 0 {
		ip.Interval = time.Second
	}
	if ip.Alignment == "" {
		ip.Alignment = "r"
	}

}

// IsTriggered specifies whether test is one of its trigger keywords.
func (ip *IpModule) IsTriggered(test string) bool {
	return test == ipTrigger
}

// Info is the main wrapper function for getting Information.
func (ip *IpModule) Info(output chan string, trigger chan string) {
	output <- ip.process()
	ticker := time.NewTicker(ip.GetInterval())
	for {
		select {
		case <-ticker.C:
			output <- ip.process()
		case <-trigger:
			output <- ip.process()
		}
	}
}

// GetAlignment allows us to get alignment.
func (ip *IpModule) GetAlignment() string {
	return ip.Alignment
}

// GetInterval allows us to get interval duration.
func (ip *IpModule) GetInterval() time.Duration {
	return ip.Interval
}

func (ip *IpModule) process() string {
	format := "%%{F%s B%s A:%s:} %s %v %%{A F- B-}"
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		return fmt.Sprintf(format, ip.colorFg, ip.colorBg, ipTrigger, "", "down")
	}
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					if strings.HasPrefix(i.Name, "wl") {
						return fmt.Sprintf(format, ip.colorFg, ip.colorBg, ipTrigger, ip.iconWif, ipnet.IP)
					} else if strings.HasPrefix(i.Name, "en") {
						return fmt.Sprintf(format, ip.colorFg, ip.colorBg, ipTrigger, ip.iconEth, ipnet.IP)
					} else {
						return fmt.Sprintf(format, ip.colorFg, ip.colorBg, ipTrigger, "?", ipnet.IP)
					}
				}
			}
		}
	}
	return fmt.Sprintf(format, ip.colorFg, ip.colorBg, ipTrigger, "", "down")

}

// Ensure Interface is Implemented
var _ Module = (*IpModule)(nil)
