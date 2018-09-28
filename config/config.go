package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
	"github.com/vgraveto/pKafka/kafka"
	"github.com/vgraveto/pKafka/pcap"
	"github.com/vgraveto/pKafka/dummy"
	"github.com/vgraveto/syslogProbe/syslogtoavro"
	"github.com/vgraveto/syslogProbe/unixsocket"
	"os"
	"bufio"
	"regexp"
	"github.com/vgraveto/syslogProbe/encoder"
)

type configType struct {
	EncodingAvro bool
	Verbose      bool
	PrintJSON    bool
	Filename     string
	EventType    string
	Kafka        kafka.Config
	Unixsocket   string
	URI          string
	UserAgent    string
	Meaning      string
	ManifestPath string
	URIofType	 string
	Type		 string
	itemsMeaningsRegex []encoder.Item

}

var (
	GlobalData    configType
)



func getRegexValues(path string)([]encoder.Item){
	//receives path to manifest
	var itemsRegex []encoder.Item

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	re0, _ := regexp.Compile(`[[](.*)[]] = [[](.*)[]]`)
	for scanner.Scan() {
		a := scanner.Text() // line
		if strings.HasPrefix(a, "#"){
			continue
		}
		if strings.HasPrefix(a, " "){
			continue
		}
		if strings.HasPrefix(a, "\n"){
			continue
		}
		line := re0.FindStringSubmatch(a)
		//if manifest contains the correct pattern collects items and meanings regex into a 2d array
		if len(line) != 0 {
			content := line[1]
			meaning := line[2]
			item := encoder.Item{content,meaning}
			itemsRegex = append(itemsRegex, item)
		}
	}
	return itemsRegex
}


func init() {

	viper.SetConfigName("syslog")
	viper.AddConfigPath("/etc/gounixsocketlistener")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("config: Config file not found...syslog.toml")
	} else {
		GlobalData.Verbose = viper.GetBool("global.verbose")
		GlobalData.EncodingAvro = viper.GetBool("global.encodingAvro")
		GlobalData.PrintJSON = viper.GetBool("global.printJSON")
		GlobalData.Filename = viper.GetString("pcap.filename")
		GlobalData.EventType = viper.GetString("avro.eventT")
		GlobalData.Kafka.Nbrokers = viper.GetInt("kafka.nbrokers")
		GlobalData.Kafka.Addrs = strings.Split(viper.GetString("kafka.addrs"), " ")
		GlobalData.Unixsocket = viper.GetString("syslogprobe.Unixsocket")
		GlobalData.UserAgent = viper.GetString("syslogprobe.UserAgent")
		GlobalData.URI = viper.GetString("syslogprobe.URI")
		GlobalData.Meaning = viper.GetString("syslogprobe.Meaning")
		GlobalData.ManifestPath = viper.GetString("syslogprobe.Manifest_path")
		GlobalData.URIofType = viper.GetString("syslogprobe.URIofType")
		GlobalData.Type = viper.GetString("syslogprobe.Type")


		if len(GlobalData.Kafka.Addrs) != GlobalData.Kafka.Nbrokers {
			log.Fatal("config: Wrong number of kafka brokers")
		}
		GlobalData.Kafka.Topic = viper.GetString("kafka.topic")
	}

	// initialize config data for kafka package
	kafka.Verbose = GlobalData.Verbose
	kafka.Cfg = GlobalData.Kafka
	// initialize config data for pcap package
	pcap.EventType = GlobalData.EventType
	pcap.Verbose = GlobalData.Verbose
	pcap.Filename = GlobalData.Filename
	// initialize config data for dummy package
	dummy.Verbose = GlobalData.Verbose
	dummy.PrintJSON = GlobalData.PrintJSON
	//initialize syslogtoavro config
	syslogtoavro.URI = GlobalData.URI
	syslogtoavro.UserAgent = GlobalData.UserAgent
	syslogtoavro.Meaning = GlobalData.Meaning
	//initialize manifest path
	syslogtoavro.URIofType = GlobalData.URIofType
	syslogtoavro.Type = GlobalData.Type
	GlobalData.itemsMeaningsRegex = getRegexValues(GlobalData.ManifestPath)
	encoder.ItemsMeaningsRegex = GlobalData.itemsMeaningsRegex
	encoder.Verbose = GlobalData.Verbose
	encoder.PrintJSON = GlobalData.PrintJSON
	//initialize socket
	unixsocket.UnixSocket = GlobalData.Unixsocket


}