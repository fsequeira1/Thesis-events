package syslogtoavro

import (
	//"time"
	//"github.com/vgraveto/pKafka/dummy"
	//"github.com/vgraveto/pKafka/CPIDSNetwork"
	"log"
	"github.com/vgraveto/syslogProbe/encoder"
)

var (
	EventType string = "other"
	Verbose   bool   = false
	URI		  string = "resource://172.27.100.127/IDS"
	UserAgent string = "IntrusionDetectionSystem"
	URIofType string = "AvroCPIDSEvents.json"
	Type 	  string = "Other"
	UserAgentVersion string = "0_4"
	Meaning   string = "nidsReport"
)

func ToAvro(inSyslog <-chan []byte, outAvro chan<- []byte, stopChan <-chan struct{}) <-chan struct{} {
	stoppedchan := make (chan struct{},1)
	go func() {
		defer func() {
			stoppedchan <- struct{}{}
		}()


	ToAvroLoop:
		for {
			select {
			case <-stopChan:
				log.Printf("syslog - ToAvro: Ctrl+C\n")
				break ToAvroLoop
			case msg, ok := <-inSyslog:
				if !ok {
					log.Println("syslog - ToAvro: Broken pipe inSyslog")
					//close(outAvro)
					break ToAvroLoop
				}

				// convert message to Avro

				if EventType == "network" {
					// create network event type
					//msg, _ = CPIDSNetwork.EncodeAvro(msg, 0)
				} else if EventType == "other" {
					// create other event type
					//msg, _ = encoder.EncodeAvro(msg, UserAgent, UserAgentVersion, URI, Meaning)
					msg, _ = encoder.NewEncodeAvro(msg, UserAgent, UserAgentVersion, URI, Meaning, URIofType, Type)
				} else { // dummy event
					//tt := time.Now()
					//msg, _ = dummy.EncodeAvro(tt.String(), 0)
					//				fmt.Printf("Metadata: %s - Payload: %d\n", tt, nmsg)
				}

				// Send message out after encoding
				select {
				case <-stopChan:
					log.Printf("syslog - ToAvro: Message not sent\n")
					break ToAvroLoop
				case outAvro <- msg:
				}

			}
		}
		log.Println("syslog - ToAvro: Closing toAvro ...")
	}()
	return stoppedchan

}