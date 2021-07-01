package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}

}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	HandleErr(gob.NewDecoder(bytes.NewReader(data)).Decode(i))
}
