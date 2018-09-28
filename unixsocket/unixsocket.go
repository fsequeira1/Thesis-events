package unixsocket

import (
	"fmt"
	"net"
	"log"
	"os"
	_"strings"
	_"bufio"
	"bufio"
)

var(
	listener net.Listener
	err error
	UnixSocket = "/tmp/snort"
)

func DeleteFile(filePath string) {
	var err = os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Socket is now available")
}


func Unixreceiver(outavro chan <- []byte, stopchan <- chan struct{}) <-chan struct{} {

	listener, err = net.Listen("unix", UnixSocket)
	if err != nil {
		log.Fatal(err)
	}
	stoppedchan := make (chan struct{},1)
	go func() {
		defer func() {
			stoppedchan <- struct{}{}
		}()
	UnixsocketLoop:
		for {
			select {
			case <-stopchan:
				log.Printf("syslog - ReadFile: Ctrl+C\n")
				break UnixsocketLoop
			default:
				fd, err := listener.Accept()
				if err != nil {
					log.Fatal(err)
					return
				}

				buff:=bufio.NewScanner(fd)
				for buff.Scan(){
					//fmt.Println("TEST: ",buff.Text())
					select {
					case <-stopchan:
						log.Printf("syslogProbe: Event dropped from unixsocket\n")
						fd.Close()
						break UnixsocketLoop
					case outavro <- []byte(buff.Text()):
					}
				}
				fd.Close()
			}
		}
		listener.Close()
		log.Println("syslog - unixReceiver: Closing toEncode ...")
	}()

	return stoppedchan
}
