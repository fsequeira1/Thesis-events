package encoder

import (
	"crypto/sha256"
	"github.com/vgraveto/pKafka/AvroCPIDSEvents"
	"time"
	"fmt"
	"encoding/hex"
	"regexp"
	"bytes"
	"strconv"
	"log"
)
var(
	Verbose   bool   = true
	ItemsMeaningsRegex []Item = nil
	PrintJSON bool = false
)
type loginfo struct{
	month string
	severity int
	day string
	hour string
	hostname string
	application string
	processid string
	message string
}

type Item struct {
	Content string
	Meaning string
}


func parseInputToStruct(data []byte)(loginfo,[]Item){
	//Doesn't have severity level cause rsyslog doesn't have it either by default
	//receives data from rsyslog and converts to string
	data = bytes.Trim(data, "\x00")
	receivedstring := string(data)

	log.Println("DEBUG LOG: ", receivedstring)
	var itemsMeanings []Item
	var facility int
	//regex to parse rsyslog logs
	//re, _ := regexp.Compile(`([a-zA-Z]+) ([0-9]{2}) ([0-9]+:[0-9]+:[0-9]+) ([a-zA-Z]+) ([a-zA-Z]+)([^:]+)?:([^\n]+)`)

	//with priority should look like this
	re, _ := regexp.Compile(`<([0-9]{1,3})>([a-zA-Z]+) *([0-9]{1,2}) ([0-9]+:[0-9]+:[0-9]+) ([a-zA-Z]+) ([a-zA-Z]+)([^:]+)?:([^\n]+)`)
	values:=loginfo{}

	// assigns log values to variables
	line := re.FindStringSubmatch(receivedstring)
	//if Verbose{
	//	fmt.Println("DEBUG LOG: ", line)
	//}
	//with priority should look like this
	if len(line)!=0 {
		priority, _ := strconv.Atoi(line[1])
		facility = priority / 8
		severity := priority - (facility * 8)
		values.severity = severity
		month := line[2]
		values.month = month
		day := line[3]
		values.day = day
		hour := line[4]
		values.hour = hour
		hostname := line[5]
		values.hostname = hostname
		application := line[6]
		values.application = application
		if len(line[7]) > 0 {
			processid := line[7]
			values.processid = processid
		}
		message := line[8]
		values.message = message
	}


	//see if string contains more items or meanings
	for index,_ := range ItemsMeaningsRegex {
		re, _ := regexp.Compile(ItemsMeaningsRegex[index].Content)
		re1, _ := regexp.Compile(ItemsMeaningsRegex[index].Meaning)
		content := re.FindString(values.message)
		if len(content)!=0{
			meaning := re1.FindString(values.message)
			if len(meaning)!=0{
				item:=Item{content,meaning}
				itemsMeanings = append(itemsMeanings, item)
			}
		}
	}
	return values, itemsMeanings
}



func NewEncodeAvro(data []byte, UserAgent string,UserAgentVersion string , URI string, Meaning string, URIofType string, Type string) ([]byte, error) {
	var Items[]map[string]interface{}
	values,itemsMeaningArray:=parseInputToStruct(data)

	if len(itemsMeaningArray)==0{
		Items=[]map[string]interface{}{
			{
				"Meaning": Meaning,
				"Content": map[string]interface{}{
					"string": values.message,
				},
			},
		}
	}else {
		for index, _ := range itemsMeaningArray {
			var item map[string]interface{}
			item = map[string]interface{}{
				"Meaning": itemsMeaningArray[index].Meaning,
					"Content": map[string]interface{}{
					"string": itemsMeaningArray[index].Content,
				},
			}
			Items = append(Items, item)
		}
	}

	if Verbose{
		log.Println(values)
		log.Println(itemsMeaningArray)
	}

	//priority wont work without changing main regex and adding priority to rsyslog default string
	var svr string
	svRand := values.severity
	switch svRand {
	case 0:
		svr = "emerg"
	case 1:
		svr = "alert"
	case 2:
		svr = "crit"
	case 3:
		svr = "err"
	case 4:
		svr = "warn"
	case 5:
		svr = "notice"
	case 6:
		svr = "info"
	case 7:
		svr = "debug"
	default:
		svr = "info"
	}
	if Verbose {
		fmt.Println("PRIORIDADE: " + svr)
	}
	// Create Avro Payload message
	Payload := map[string]interface{}{
		"Events": []map[string]interface{}{
			{
				"URIofType": URIofType,
				"Type":      Type,
				"Severity":  svr,
				"Data": map[string]interface{}{
					"pt.uc.dei.atena.datamodel.Other": map[string]interface{}{
						"Code": 0,
						"Items": Items,
					},
				},
			},
		},
	}

	// Compute checksum
	PayloadBytes, err := AvroCPIDSEvents.CodecPayload.BinaryFromNative(nil, Payload)
	AvroCPIDSEvents.CheckError("Payload - BinaryFromNative", err)
	h := sha256.New()
	h.Write(PayloadBytes)
	hh := hex.EncodeToString(h.Sum(nil))
	//hh := string(h.Sum(nil))

	// Create Avro Message
	AvroMessage := map[string]interface{}{
		"Metadata": map[string]interface{}{
			"Origins": []map[string]interface{}{
				{
					//fazer parse do ip da maquina de onde vem
					"URI":       URI,
					"UserAgent": UserAgent+"/"+UserAgentVersion,
					"Timestamp": time.Now().Format(time.RFC3339Nano),
					"PreviousOrigin": map[string]interface{}{
						"null": nil,
					},
				},
			},
			"Checksum": hh,
		},
		"Payload": Payload,
	}

	binary, err := AvroCPIDSEvents.CodecMessage.BinaryFromNative(nil, AvroMessage)
	AvroCPIDSEvents.CheckError("Message - BinaryFromNative", err)

	if PrintJSON {
		// print message in JSON Avro format
		text, _ := AvroCPIDSEvents.CodecMessage.TextualFromNative(nil, AvroMessage)
		fmt.Printf("Other EncodeAvro: %s\n", text)
	}

	if Verbose {
		fmt.Printf("Encoded in avro msg \n")
	}
	return binary, err
}

