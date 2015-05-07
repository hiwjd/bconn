package bconn

import (
	"errors"
	"log"
	"net"
)

// A buffered net.Conn
// ...
// bconn := NewBufferedConn(conn)
// for {
//     err := bconn.Read() // must call this at first
//     if bconn.Buffered() < minLen { // determine whether the received data is long enough
//         continue
//     }
//     b := bconn.Peek(0, minLen) // retrieve ref of data between[0...minLen]
//     // processing...
//     bconn.Clear(minLen) // after process clean the data [0...minLen]
// }
type BufferedConn struct {
	conn net.Conn
	pos  int // current data length filled in the buf
	buf  []byte
	size int // buf`s max length
}

func NewBufferedConn(conn net.Conn) *BufferedConn {
	size := 2048 // default size of buf
	buf := make([]byte, size)
	return &BufferedConn{
		conn: conn,
		pos:  0,
		buf:  buf,
		size: size,
	}
}

// size: init buf size
func NewBufferedConnWithSize(conn net.Conn, size int) *BufferedConn {
	buf := make([]byte, size)
	return &BufferedConn{
		conn: conn,
		pos:  0,
		buf:  buf,
		size: size,
	}
}

// read data from connection, this function must be called
func (this *BufferedConn) Read() error {
	n, err := this.conn.Read(this.buf[this.pos:])
	this.pos = this.pos + n
	if err != nil {
		return err
	}
	return nil
}

// return current length of data filled in buf
func (this *BufferedConn) Buffered() int {
	return this.pos
}

// return the buf`s max size
func (this *BufferedConn) Size() int {
	return this.size
}

// return the ref of the data between start and stop(exclued)
func (this *BufferedConn) Peek(start, stop int) []byte {
	if start > this.pos || stop < 0 {
		return nil
	}

	return this.buf[start:stop]
}

// increase the buf`s max size
// extendLen: the length of enlarged part
func (this *BufferedConn) Grow(extendLen int) {
	extendBuf := make([]byte, extendLen)
	this.buf = append(this.buf, extendBuf[:]...)
	this.size = this.size + extendLen
	log.Printf("the buf`s max size has been set: %d \n", this.size)
}

// clean the data between 0 to position
func (this *BufferedConn) Clear(position int) error {
	if position > this.size || position < 0 {
		return errors.New("position wrong")
	}

	newPos := this.pos - position
	if position < this.pos {
		// |____________________________________________|
		// |<---              pos               --->|
		// |<-position->|
		//	        |<-        newPos        ->|
		//dst := this.buf[:position]
		src := this.buf[position:this.pos]
		for i := range src {
			this.buf[i] = src[i]
		}
		needZero := this.buf[newPos:this.pos]
		for i := range needZero {
			needZero[i] = 0
		}
	} else {
		needZero := this.buf[:position]
		for i := range needZero {
			needZero[i] = 0
		}
	}
	this.pos = newPos

	return nil
}
