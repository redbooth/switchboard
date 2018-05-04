package header

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
)

type Header interface {
	Id() []byte
	String() string
}

type Conf interface{}

type Constructor (func(Conf) Header)

/*
// Note: struct fields _must_ be exported (i.e. capitalized)
type Example struct {
	First, Second, Third byte
}

h := new(Example)
r := bytes.NewReader([]byte{0x01, 0x02, 0x03})
Extract(h, r)
fmt.Printf("%#v\n", h)
>> &main.Example{First:0x1, Second:0x2, Third:0x3}
*/
func Extract(header Header, reader io.Reader) {
	err := binary.Read(reader, binary.LittleEndian, header)
	if err != nil {
		log.Panicf("Unable to extract header from stream: %v\n", err)
	}
}

// Recreate the original input stream
func Inject(header Header, reader io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		err := binary.Write(pw, binary.LittleEndian, header)
		if err != nil {
			log.Panicf("Unable to inject header into stream: %v\n", err)
		}
	}()
	return io.MultiReader(pr, reader)
}

func Bytes(header Header) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, header)
	return buf.Bytes(), err
}
