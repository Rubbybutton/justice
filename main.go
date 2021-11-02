package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	// encrypt
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const CTRkey = "1443flfsaWfdasds"

// CTR encrypt and decrypt
func aesCtrCrypt(plainText []byte, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)

	return dst, nil
}

func system() (string, []string) {
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
	return hostsPath, flushDnsCmd
}

func tryCmd(cmds []string) error {
	for _, cmd := range cmds {
		fmt.Println("[Exetuting] " + cmd)
		splitCmd := strings.Split(cmd, " ")
		execCmd := exec.Command(splitCmd[0], splitCmd[1:]...)
		_, err := execCmd.CombinedOutput()
		if err != nil {
			fmt.Println("[Execute command (", cmd, ") error]")
		} else {
			// fmt.Println("ouput: ", string(out))
			return nil // when successfully execute a command, than the program exit
		}
	}
	return nil
}

func openCompetitionPage(os string, url string) {
	var cmds []string
	switch os {
	case "windows":
		cmds = append(cmds, "cmd /c start "+url)
	case "linux":
		cmds = append(cmds, "xdg-open "+url)
	case "darwin":
		cmds = append(cmds, "open "+url)
	}
	tryCmd(cmds)
}

func genHosts(hosts []string) string {
	res := ""
	for _, host := range hosts {
		res += "127.0.0.1 " + host + "\n"
	}
	return res
}

func isInCompetitionTime(begin string, end string) bool {
	// though users can change time to pass the program

	timeLocal, _ := time.LoadLocation("Asia/Shanghai")
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", begin, timeLocal)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", end, timeLocal)
	nowTime := time.Now()
	if nowTime.After(startTime) && nowTime.Before(endTime) {
		return true
	}
	return false
}

func add2Hosts(bannedHosts []string, hostsPath string, flushDnsCmd []string) {

	hostsStr := genHosts(bannedHosts)

	// write to hosts
	filePath := hostsPath
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	check(err)
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(hostsStr)
	write.Flush()

	//flush DNS
	tryCmd(flushDnsCmd)
}

func main() {
	var bannedHosts = []string{"www.google.com", "www.google.cn", "www.google.hk", "www.google.com.hk",
		"www.baidu.com",
		"www.bing.com", "cn.bing.com",
		"duckduckgo.com", "www.duckduck.com"}
	hostsPath, flushDnsCmd := system()

	curFolderPath, _ := os.Getwd()
	tmpFilePath := path.Join(curFolderPath, "tmpFile")
	// FIXME change the startTime and endTime into proper time
	// NOTE startTime should be 30m before the competition start
	if isInCompetitionTime("2021-11-02 18:00:00", "2021-11-02 23:00:00") {
		// write encrypted original content to tmpFile
		originHosts, _ := ioutil.ReadFile(hostsPath)
		originHosts, _ = aesCtrCrypt(originHosts, []byte(CTRkey))
		err := ioutil.WriteFile(tmpFilePath, originHosts, 0666)
		check(err)
		// write to hosts
		add2Hosts(bannedHosts, hostsPath, flushDnsCmd)

		fmt.Println("All Done!")
		competitionUrl := "https://dict.youdao.com/search?q=is%20opening&keyfrom=new-fanyi.smartResult"
		// exec.Command("xdg-open", competitionUrl).Run()
		// openCompetitionPage(runtime.GOOS, competitionUrl)
		fmt.Println("The competiton page is opening...")
		fmt.Println("if it goes wrong, please visit the website " + competitionUrl)
		fmt.Println("[Please don't delete tmpfile in the current directory]")
		fmt.Println("[Please reexecute this program after the competition to ensure your internet connection is normal]")
	} else {
		// read and
		encrypted, e := ioutil.ReadFile(tmpFilePath)
		check(e)
		originHosts, _ := aesCtrCrypt(encrypted, []byte(CTRkey))
		err := ioutil.WriteFile(hostsPath, originHosts, 0666)
		check(err)
		fmt.Println("[ALL have been restored to original state]")
		fmt.Println("Appreciation for yout attendance and Congratulations for your completion!")
	}
	fmt.Println("[Press ENTER key to exit...]")
	fmt.Scanln()
}
