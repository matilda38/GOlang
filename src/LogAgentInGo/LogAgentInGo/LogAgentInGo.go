package LogAgentInGo

import (
	"bytes"
	"compress/gzip"
	"math"
	"net"
	"log"
	"crypto/rand"
	"fmt"
)

const (
	defaultGraylogPort = 12201
	defaultGraylogHostname = "0.0.0.0"
	defaultConnection = "wan"
	defaultMaxChunkSizeWan = 1420
	defaultMaxChunkSizeLan = 8154
)
type Config struct{
	GraylogPort int
	GraylogHostname string
	Connection string
	MaxChunkSizeWan int
	MaxChunkSizeLan int
	Protocol string
}

type OutputChannel chan []byte
type InputChannel chan []byte

type RingBuffer struct{
	Input InputChannel
	Output OutputChannel
}

type Gelf struct{
	Config Config
	RingBuffer RingBuffer
}

func New(config Config) *Gelf{
	if config.GraylogPort == 0{
		config.GraylogPort = defaultGraylogPort
	}
	if config.GraylogHostname == ""{
		config.GraylogHostname = defaultGraylogHostname
	}
	if config.Connection == ""{
		config.Connection = defaultConnection
	}
	if config.MaxChunkSizeWan == 0{
		config.MaxChunkSizeWan = defaultMaxChunkSizeWan
	}
	if config.MaxChunkSizeLan == 0{
		config.MaxChunkSizeLan = defaultMaxChunkSizeLan
	}
	if config.Protocol == ""{
		config.Protocol = "UDP"
	}

	oc := make(chan []byte, 1460)
	ic := make(chan []byte)
	rb := RingBuffer{ic,oc}

	g:= &Gelf{
		Config : config,
		RingBuffer: rb,
	}

	return g
}

func (g *Gelf) InputMessage(message chan<- string) {
	sampleDoc := "{ \"version\" : \"1.1\", \"host\": \"example.org\", \"short_message\": \"A short message that helps you identify what is going on\", \"full_message\": \"Backtrace here\n\nmore stuff\"}"
	message <- sampleDoc
}

func (g *Gelf) ProcessMessage(message <-chan string){
	messageToProcess := <- message

	if g.Config.Protocol == "UDP"{
		g.LogUDP(messageToProcess)
		close(g.RingBuffer.Input)
	} else {
		messageToProcess = fmt.Sprint(messageToProcess, "\000")
		g.RingBuffer.Input <- ([]byte)(messageToProcess)
	}
}

//When message comes
//1. compress
//2. if the length of compressed data is larger than the chuncksize,
// slice it to each chunk and send each chunk.
func (g *Gelf) LogUDP(message string) {
	compressed := g.Compress([]byte(message))
	chunksize := g.Config.MaxChunkSizeWan

	length :=compressed.Len()

	if length > chunksize {
		chunkCountInt := int(math.Ceil(float64(length) / float64(chunksize)))

		messageId := make([]byte, 8)
		rand.Read(messageId)

		//until length of the message
		for i, index := 0,0; i < length; i, index = i + chunksize , index+1{
			packet := g.CreateChunkedMessage(index, chunkCountInt, messageId, &compressed)
			g.RingBuffer.Input <- packet.Bytes()
		}
	} else{
		g.RingBuffer.Input <- compressed.Bytes()
	}
}


func (g *Gelf) CreateChunkedMessage(index int, chunkCountInt int, messageId []byte, compressed *bytes.Buffer) bytes.Buffer{
	var packet bytes.Buffer

	//temporary default
	chunksize := 1420

	packet.WriteByte(byte(30))
	packet.WriteByte(byte(15))
	packet.Write(messageId)

	packet.WriteByte(byte(index))
	packet.WriteByte(byte(chunkCountInt))

	packet.Write(compressed.Next(chunksize))

	return packet
}

func (g *Gelf) Compress(b []byte) bytes.Buffer{
	var buf bytes.Buffer

	//zlip not working
	comp := gzip.NewWriter(&buf)
	comp.Write(b)
	comp.Close()

	return buf
}

//For UDP Gelf
//1. receive data
//2. write it in ringbuffer
//3. setup for connection
//4. send data in ringbuffer
//Consumer
func (g *Gelf) Send (){
	if g.Config.Protocol == "UDP"{
		g.ConnectUDP()
	} else {
		g.ConnectTCP()
	}
}

func (g *Gelf) Consume(){
	for input := range g.RingBuffer.Input{
		select {
		case g.RingBuffer.Output <- input:
		default:
			<-g.RingBuffer.Output
			g.RingBuffer.Output <- input
		}
	}
	close(g.RingBuffer.Output)
}

func (g *Gelf) ConnectUDP(){
	var addr = "0.0.0.0:12201"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err!=nil{
		log.Printf("%s", err)
		return
	}

	for output := range g.RingBuffer.Output{
		conn.Write(output)
	}
}

// For TCP GELF : no chunking no compressing
func (g *Gelf) ConnectTCP() {
	var addr = "0.0.0.0:5555"
	conn, err := net.Dial("tcp", addr)
	if err != nil{
		log.Print("%s", err)
		return
	}
	defer conn.Close()

	for output := range g.RingBuffer.Output{
		conn.Write(output)
	}
}