package main

import (
	"os"
	"log"
	"os/signal"
	"syscall"
	"github.com/vgraveto/syslogProbe/config"
	"github.com/vgraveto/pKafka/kafka"
	"github.com/vgraveto/syslogProbe/unixsocket"
	"github.com/vgraveto/syslogProbe/syslogtoavro"
	//"time"
	"sync"
)

var syslogVersion string = "0_4"
func main(){
	log.Printf("main: started syslog_%s\n", syslogVersion)
	syslogtoavro.UserAgentVersion = syslogVersion

	//remove unix socket
	//depois de isto ser removido o go já não cria o socket... dunno y

	/*
	//for some reason this deletes the unix socket, but after that the socket isn't starting
	if _, err := os.Stat("/tmp/snort"); !os.IsNotExist(err) {
		unixsocket.DeleteFile("/tmp/snort")
	}

	time.Sleep(10 * time.Second)
	*/


	toEncode := make(chan []byte)
	toKafka:=  make(chan []byte)
	stopChan:=   make(chan struct{},3)

	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		log.Println("main: Stopping ...")
		stopChan <- struct{}{}
		stopChan <- struct{}{}
		stopChan <- struct{}{}

		close(stopChan)
	}()
	signal.Notify(signalChan,syscall.SIGINT,syscall.SIGTERM)
	var mu sync.Mutex
	writeStoppedChan := kafka.Write(toKafka, stopChan,mu)
	if config.GlobalData.EncodingAvro {
		toAvroStoppedChan := syslogtoavro.ToAvro(toEncode, toKafka, stopChan)
		readStoppedChan := unixsocket.Unixreceiver(toEncode, stopChan)
		<-readStoppedChan
		<-toAvroStoppedChan
	} else {
		readStoppedChan := unixsocket.Unixreceiver(toKafka, stopChan)
		<-readStoppedChan
	}


	<-writeStoppedChan
	log.Printf("main: ended %s\n", syslogVersion)
}
