package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

var ip= "192.168.1.2"

func main(){
	activeThreads := 0
	doneChannel := make(chan bool)
	for port := 0; port <=1024; port++{
		go grabBanner(ipToScan, port, doneChannel)
		activeThreads++
	}

	for activeThreads > 0 {
		<-doneChannel
		activeThreads--
	}
}

func grabBanner(ip string, port int, doneChannel chan bool){
	connection, err := net.DialTimeout(
		"tcp",
		ip+":"+strconv.Itoa(port),
		time.Second*10)
	if err != nil {
		doneChannel <- true
		return
	}
	//See if server offers anything to read
	//quizás hay que subir el número de bytes que recibimos idk
	buffer := make([]byte,4096)
	connection.SetReadDeadline(time.Now().Add(time.Second*5))

	//Set timeOut
	numBytesRead, err := connection.Read(buffer)
	if err != nil {
		doneChannel <- true
		return
	}
	log.Printf("Banner port  %d:\n%s\n ", port, buffer[0:numBytesRead])
	doneChannel <- true
}