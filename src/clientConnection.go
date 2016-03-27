package src
import (
	"net"
	"fmt"
)

type ClientConnection struct {
	address Address //TCPAddress
	socket net.Conn //TCPConnection

	lastRead int
	readBuffer []byte
}


func (connection *ClientConnection) Connect() {
	buffer := make([]byte,3)
	buffer = []byte(CLIENT_BINARY_NEW)

	connection.socket, _ = net.Dial("tcp", "127.0.0.1:5701")
	fmt.Println(connection.socket.RemoteAddr())
//	connection.socket, _ = net.Dial("tcp",connection.address.Host + ":" + connection.address.Port)
	connection.socket.Write(buffer)
	request := EncodeRequest("dev", "dev-pass", nil, nil, true, "GO", 1)
	request.SetCorrelationId(0)
	request.SetPartitionId(-1)
	request.SetFlags(BEGIN_END_FLAG)
	connection.socket.Write(request.Buffer)

	rBuffer := make([]byte, 360)

	readBytes , _ := connection.socket.Read(rBuffer)
//	readBytes, err := bufio.NewReader(connection.socket).Read(connection.readBuffer)
	fmt.Println(readBytes)
	message := CreateForDecode(rBuffer[0:readBytes])

	response := DecodeResponse(message)

	fmt.Println(response.Address)
	fmt.Println(response.Uuid)
	fmt.Println(response.OwnerUuid)
}
