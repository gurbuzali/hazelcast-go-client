package src
import (
	"net"
	"fmt"
	"bufio"
	"errors"
	"strconv"
	"encoding/binary"
)

type ClientConnection struct {
	address Address  //TCPAddress
	socket  net.Conn //TCPConnection
	readBuffer []byte
}

func NewClientConnection(address Address) *ClientConnection {
	connection := new(ClientConnection)
	connection.address = address
	connection.readBuffer = make([]byte,0) //todo
	return connection
}

func (this *ClientConnection) Connect(address Address) *Promise {
	result := new(Promise)

	result.SuccessChannel = make(chan interface{}, 1)
	result.FailureChannel = make(chan error, 1)

	go func() {
		socket, err := net.Dial("tcp", this.address.Host + ":" + strconv.Itoa(this.address.Port))
		this.socket = socket
		if err == nil {
			this.socket.Write([]byte(CLIENT_BINARY_NEW))
			result.SuccessChannel <- this
		} else {
			result.FailureChannel <- errors.New("Could not connect to address" + this.address.String())
			fmt.Println(err)
		}
	}()

	return result
}

func (this *ClientConnection) registerResponseCallback(callback func([]byte)) {
	rBuffer := make([]byte, 1024) //todo dynamic
	this.socket.Read(rBuffer)
	this.readBuffer = append(this.readBuffer,rBuffer)
	for ;this.readBuffer >= INT_SIZE_IN_BYTES; {
		frameLength := binary.LittleEndian.Uint32(this.readBuffer[0:])
		if frameLength > len(this.readBuffer) {
			return
		}
		message := make([]byte,frameLength)
		copy(message,this.readBuffer[0:frameLength])
		this.readBuffer = this.readBuffer[frameLength:]
		callback(message)
	}
}
