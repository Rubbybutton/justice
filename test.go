package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func tryCmd(cmds []string) error {
	for _, cmd := range cmds {
		splitCmd := strings.Split(cmd, " ")
		execCmd := exec.Command(splitCmd[0], splitCmd[1:]...)
		out, err := execCmd.CombinedOutput()
		if err != nil {
			fmt.Println("[Execute command (", cmd, ") error]")
		} else {
			fmt.Println("ouput: ", string(out))
		}
	}
	return nil
}

func main() {
	cmds := []string{
		"toucha hello.txt",
		"pwd"}
	print(tryCmd(cmds))
}
