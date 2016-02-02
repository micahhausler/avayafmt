package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	flag "github.com/spf13/pflag"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "0102061504"
)

var portP = flag.IntP("port", "p", 9000, "Default port to use")
var addressP = flag.StringP("address", "a", "127.0.0.1", "Default listen address to use")

type Header struct {
	Data [16]byte
}

type Record struct {
	Date        [6]byte
	_           byte
	Time        [4]byte
	_           byte
	SecDur      [5]byte
	_           byte
	CodCode     [1]byte
	_           byte
	CodeDial    [4]byte
	_           byte
	CodeUsed    [4]byte
	_           byte
	DialedNum   [18]byte
	_           byte
	ClgNum      [10]byte
	_           byte
	AuthCode    [7]byte
	_           byte
	InCrtId     [3]byte
	_           byte
	OutCrtId    [3]byte
	_           byte
	IsDnCc      [11]byte
	_           byte
	Ppm         [5]byte
	_           byte
	AcctCode    [15]byte
	_           byte
	InTrkCode   [4]byte
	_           byte
	AttdConsole [2]byte
	_           byte
	Vdn         [5]byte
	_           [2]byte
}

func (r Record) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	var date interface{}
	date, err := time.Parse(dateFormat, string(append(r.Date[:], r.Time[:]...)))
	if err != nil {
		date = string(append(r.Date[:], r.Time[:]...))
	}

	dur, err := strconv.Atoi(string(r.SecDur[:]))
	if err != nil {
		panic(err)
	}

	data["Date"] = date
	data["Duration"] = dur
	data["CodCode"] = strings.Trim(string(r.CodCode[:]), " ")
	data["CodeDial"] = strings.Trim(string(r.CodeDial[:]), " ")
	data["CodeUsed"] = strings.Trim(string(r.CodeUsed[:]), " ")
	data["DialedNum"] = strings.Trim(string(r.DialedNum[:]), " ")
	data["ClgNum"] = strings.Trim(string(r.ClgNum[:]), " ")
	data["AuthCode"] = strings.Trim(string(r.AuthCode[:]), " ")
	data["InCrtId"] = strings.Trim(string(r.InCrtId[:]), " ")
	data["OutCrtId"] = strings.Trim(string(r.OutCrtId[:]), " ")
	data["IsDnCc"] = strings.Trim(string(r.IsDnCc[:]), " ")
	data["Ppm"] = strings.Trim(string(r.Ppm[:]), " ")
	data["AcctCode"] = strings.Trim(string(r.AcctCode[:]), " ")
	data["InTrkCode"] = strings.Trim(string(r.InTrkCode[:]), " ")
	data["AttdConsole"] = strings.Trim(string(r.AttdConsole[:]), " ")
	data["Vdn"] = strings.Trim(string(r.Vdn[:]), " ")

	return json.Marshal(data)
}

func main() {
	flag.Parse()

	runserver(*addressP, *portP)

}

// runserver is the main event loop for the server
func runserver(host string, port int) {

	hp := fmt.Sprintf("%s:%s", host, strconv.Itoa(port))
	fmt.Printf("Listening on %s\n", hp)
	listener, err := net.Listen("tcp", hp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		c, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Connection from %v established.\n", c.RemoteAddr())
		go Handler(c)
	}
}

// Handler handels a network connection
func Handler(c net.Conn) {
	defer c.Close()
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in Handler () for remote %s: %s", c.RemoteAddr(), r)
		}
	}()

	// Read Header
	head := new(Header)
	err := binary.Read(c, binary.LittleEndian, head)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(head)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	for {
		rec := new(Record)
		err := binary.Read(c, binary.LittleEndian, rec)
		if err != nil && err != io.EOF {
			// pass
		}
		if err == io.EOF {
			break
		}
		//data, err := json.MarshalIndent(rec, "", "  ")
		data, err := json.Marshal(rec)
		if err != nil {
			fmt.Printf("Error Marshaling: %s\n", err)
		}
		fmt.Println(string(data))
	}

}
