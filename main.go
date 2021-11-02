package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func system() (string, string) {
	var hostsPath string
	var flushDnsCmd []string
	switch runtime.GOOS {
	case "windows":
		hostsPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
		flushDnsCmd = append(flushDnsCmd, "ipconfig /flushdns")
	case "darwin":
		hostsPath = "/etc/hosts"
		// modern version
		flushDnsCmd = append(flushDnsCmd, "sudo dscacheutil -flushcache;sudo killall -HUP mDNSResponder")
		// old version
		flushDnsCmd = append(flushDnsCmd, "sudo discoveryutil udnsflushcaches;sudo discoveryutil mdnsflushcaches")
	case "linux":
		hostsPath = "/etc/hosts"
		// don't do anything to dns, for some distributions they don't have dns services
		// hopefully, the number of linux users is small
	}

}

func genHosts(hosts []string) string {
	res := ""
	for _, host := range hosts {
		res += "127.0.0.1 " + host + "\n"
	}
	return res
}

func main() {

	var bannedHost = []string{"www.google.com", "www.google.cn", "www.google.hk", "www.google.com.hk",
		"www.baidu.com",
		"www.bing.com", "cn.bing.com",
		"duckduckgo.com", "www.duckduck.com"}
	// bannedHost := []string{"www.baidu.com"}

	hostsStr := genHosts(bannedHost)

	// write to hosts
	filePath := hostsPath
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	check(err)
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(hostsStr)
	write.Flush()
	fmt.Println("Done")
}
