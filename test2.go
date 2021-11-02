package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func tryCmd(cmds []string) error {
	for _, cmd := range cmds {
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

func main() {
	var cmds []string
	url := "https://www.jianshu.com"
	switch runtime.GOOS {
	case "windows":
		cmds = append(cmds, "cmd /c start "+url)
	case "linux":
		cmds = append(cmds, "xdg-open "+url)
	case "darwin":
		cmds = append(cmds, "open "+url)
	}
	// tryCmd(cmds)
	exec.Command("xdg-open", url).Start()
	//
	// fmt.Println(cmds)
	// exec.Command(`xdg-open`, `https://www.jianshu.com`).Start()
}
