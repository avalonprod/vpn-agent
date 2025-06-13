package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func isPeerExists(pubKey string) bool {
	cmd := exec.Command("wg", "show", ServerConfig.WGInterface, "peers")
	output, err := cmd.Output()
	if err != nil {
		log.Println("wg show error:", err)
		return false
	}

	peers := strings.Split(string(output), "\n")
	for _, peer := range peers {
		if peer == pubKey {
			return true
		}
	}

	return false
}

func getNextIP() string {
	for i := nextIP; i < 255; i++ {
		ip := fmt.Sprintf("%s%d", ServerConfig.BaseIP, i)
		if !usedIPs[ip] {
			usedIPs[ip] = true
			nextIP = i + 1
			return ip
		}
	}
	return ""
}

func restoreUsedIPs() {
	cmd := exec.Command("wg", "show", ServerConfig.WGInterface, "allowed-ips")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("failed to get IPs list: %v", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		ipCIDR := fields[1]
		ip := strings.Split(ipCIDR, "/")[0]
		usedIPs[ip] = true
	}
}
