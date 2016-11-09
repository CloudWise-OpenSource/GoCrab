package main

import (
	"github.com/CloudWise-OpenSource/GoCrab/Libs/goavro"
	"bytes"
	//"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	//"time"
	//"fmt"
	//"io/ioutil"
	"log"
)

var host = flag.String("host", "localhost", "host")
var port = flag.String("port", "26799", "port")

//go run timeclient.go -host time.nist.gov
func main() {
	flag.Parse()

	recordSchemaJSON := `
{
  "type": "record",
  "name": "comments",
  "doc:": "A basic schema for storing blog comments",

  "fields": [
    {
      "doc": "Name of user",
      "type": "string",
      "name": "username"
    },
    {
      "doc": "The content of the user's message",
      "type": "string",
      "name": "comment"
    },
    {
      "doc": "Unix epoch time in milliseconds",
      "type": "long",
      "name": "timestamp"
    }
  ]
}
`
	someRecord, err := goavro.NewRecord(goavro.RecordSchema(recordSchemaJSON))
	if err != nil {
		log.Fatal(err)
	}
	// identify field name to set datum for
	someRecord.Set("username", " neeke ")
	someRecord.Set("comment", "The Atlantic is oddly cold this morning!")
	someRecord.Set("timestamp", int64(1082196484))

	codec, err := goavro.NewCodec(recordSchemaJSON)
	if err != nil {
		log.Fatal(err)
	}

	bb := new(bytes.Buffer)
	if err = codec.Encode(bb, someRecord); err != nil {
		log.Fatal(err)
	}

	//actual := bb.Bytes()

	bf := bytes.NewBuffer(bb.Bytes())
	//body := ioutil.NopCloser(bf)

	addr, err := net.ResolveUDPAddr("udp", *host+":"+*port)
	if err != nil {
		fmt.Println("Can't resolve address: ", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Can't dial: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(bf.String()))
	if err != nil {
		fmt.Println("failed:", err)
		os.Exit(1)
	}

	//data := make([]byte, 4)
	//_, err = conn.Read(data)
	//if err != nil {
	//	fmt.Println("failed to read UDP msg because of ", err)
	//	os.Exit(1)
	//}

	//t := binary.BigEndian.Uint32(data)
	//fmt.Println(time.Unix(int64(t), 0).String())

	os.Exit(0)
}
