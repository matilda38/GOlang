package LogAgentInGo

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"math"
	"net"
	"log"
	"crypto/rand"
	"github.com/glycerine/rbuf"
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

type Gelf struct{
	Config Config
	RingBuffer rbuf.FixedSizeRingBuf
	Buffer bytes.Buffer
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

	rb := rbuf.NewFixedSizeRingBuf(1420000)
	var queue bytes.Buffer

	g:= &Gelf{
		Config : config,
		RingBuffer : *rb,
		Buffer : queue,
	}
	return g
}

func (g *Gelf) Log(message string){
	if g.Config.Protocol == "UDP"{
		g.LogUDP(message)
	} else{
		message = fmt.Sprint(message, "\000")
		g.SendTCP([]byte(message))
	}
}

//When message comes
//1. compress
//2. if the length of compressed data is larger than the chuncksize,
// slice it to each chunk and send each chunk.

func (g *Gelf) LogUDP(message string){
	//step 1 : compress
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
			g.SendUDP(packet.Bytes())
		}
	} else{
		g.SendUDP(compressed.Bytes())
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

func (g *Gelf) IntToBytes (i int) []byte {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int8(i))
	if err != nil{
		log.Print("%s", err)
	}
	return buf.Bytes()
}

//For UDP Gelf
//1. receive data
//2. write it in ringbuffer
//3. if the data exceeds the size of ringbuffer size(1420000),
//	write it in normal buffer(unlimited size)
//4. setup for connection
//5. send data in ringbuffer
func (g *Gelf) SendUDP (b []byte){
	var n int
	var err error

	n, err = g.RingBuffer.Write(b)

	if err != nil{
		g.Buffer.Write(b[n:])
	}

	var addrU = "0.0.0.0:12201"
	udpAddr, err := net.ResolveUDPAddr("udp", addrU)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err!=nil{
		log.Printf("%s", err)
		return
	}
	conn.Write(g.RingBuffer.Bytes())
}

// For TCP GELF : no chunking no compressing
func (g *Gelf) SendTCP(b []byte) {
	var addr = "0.0.0.0:5555"
	conn, err := net.Dial("tcp", addr)
	if err != nil{
		log.Print("%s", err)
		return
	}
	defer conn.Close()
	conn.Write(b)
}