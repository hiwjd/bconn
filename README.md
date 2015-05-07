# bconn
A go package focus on help parse tcp binary protocol packet

## example

### here is a protocol:

|2 bytes version|4 bytes body length|8 bytes packet id|the body|

```go
type Packet struct {
    Version string
    BodyLen uint32
    Pid     uint64
    Body    []byte
}

func parsePacketFromConn(conn net.Conn, cb func(*Packet)) {
	headLen := 14
	bc := bconn.NewBufferedConn(conn)
	for {
		bc.Read()
		if bc.Buffered() < headLen {
			continue
		}
		
		b := bc.Peek(2, 6)
		bodyLen := binary.BigEndian.Uint32(b)
		totLen := headLen + int(bodyLen)
		if bc.Buffered() < totLen {
			if bc.Size() < totLen {
				bc.Grow(totLen - bc.Size())
			}
			continue
		}

		// totLen is a packet`s complete length
		packet := someFuncMakeBytesToPacket(bc.Peek(0, totLen))
		bc.Clear(totLen)
		cb(packet)
	}
}

func main() {
	addr := "127.0.0.1:8080"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Printf("start server at %s\n", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			parsePacketFromConn(conn, func(p *context.Packet) {
				fmt.Printf("got packet %#v \n", p)
			})
		}(conn)
	}
}
```
