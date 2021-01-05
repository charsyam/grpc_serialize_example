package main

import (
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"errors"
    "io"
    "os"
    "log"
    pb "grpc_serialize_example/proto"
)

var errInvalidVarint = errors.New("invalid varint32 encountered")

func ReadDelimited(r io.Reader, m proto.Message) (n int, err error) {
	// Per AbstractParser#parsePartialDelimitedFrom with
	// CodedInputStream#readRawVarint32.
	var headerBuf [binary.MaxVarintLen32]byte
	var bytesRead, varIntBytes int
	var messageLength uint64
	for varIntBytes == 0 { // i.e. no varint has been decoded yet.
		if bytesRead >= len(headerBuf) {
			return bytesRead, errInvalidVarint
		}
		// We have to read byte by byte here to avoid reading more bytes
		// than required. Each read byte is appended to what we have
		// read before.
		newBytesRead, err := r.Read(headerBuf[bytesRead : bytesRead+1])
		if newBytesRead == 0 {
			if err != nil {
				return bytesRead, err
			}
			// A Reader should not return (0, nil), but if it does,
			// it should be treated as no-op (according to the
			// Reader contract). So let's go on...
			continue
		}
		bytesRead += newBytesRead
		// Now present everything read so far to the varint decoder and
		// see if a varint can be decoded already.
		messageLength, varIntBytes = proto.DecodeVarint(headerBuf[:bytesRead])
	}
	messageBuf := make([]byte, messageLength)
	newBytesRead, err := io.ReadFull(r, messageBuf)
	bytesRead += newBytesRead
	if err != nil {
		return bytesRead, err
	}
	return bytesRead, proto.Unmarshal(messageBuf, m)
}

func main() {
	in, err := os.Open("out.bin")
    if err != nil {
        log.Fatalln("Error reading file:", err)
    }

    for {
        metric := &pb.Metric{}
        n, err := ReadDelimited(in, metric)
        if n == 1 {
            break
        }
        if err != nil {
            log.Printf("err %s", err)
        }
        log.Printf("metric: %d %s %s %s", n, metric.Name, metric.Type, metric.Tags)
    }
}
