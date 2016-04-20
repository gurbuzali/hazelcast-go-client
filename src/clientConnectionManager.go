package src
import (
	"errors"
	"fmt"
)

type ClientConnectionManager struct {

}

func (manager *ClientConnectionManager) GetOrConnect(address Address) *Promise {
	connection := NewClientConnection(address)

	promise := connection.Connect(address)
	promise2 := promise.Then(func(obj interface{}) (interface{}, error) {
		connection := obj.(*ClientConnection)
		fmt.Println("Connection:" + connection.address.String())
		return connection,nil
	}, func(err error) error{
		fmt.Println("Connection Failed")
		return err
	})

	promise3 := promise2.ThenPromise(func(obj interface{}) *Promise {
		connection := obj.(*ClientConnection)
		return authenticate(connection)
	}, func(err error) error{
		fmt.Println("Authentication problem")
		return err
	})

//	promise4 := promise3.ThenSuccessReturnPromise(func(obj interface{}) *Promise {
//		isAuthenticated := obj.(*bool)
//		if(isAuthenticated){
//			fmt.Println("Hello")
//		}
//		return nil
//	}, func(err error) {
//		fmt.Println("Not Hello")
//	})

	return promise3
}

func authenticate(connection *ClientConnection) *Promise {
	result := new(Promise)

	result.SuccessChannel = make(chan interface{}, 1)
	result.FailureChannel = make(chan error, 1)

	//////////////////////////////////////////
	request := EncodeRequest("dev", "dev-pass", nil, nil, true, "GO", 1) //config
	request.SetCorrelationId(0)
	request.SetPartitionId(-1)
	request.SetFlags(BEGIN_END_FLAG)

	connection.socket.Write(request.Buffer)

	rBuffer := make([]byte, 1024)
	readBytes, _ := connection.socket.Read(rBuffer)
	fmt.Println(readBytes)
	response := CreateForDecode(rBuffer[:readBytes])
	///////////////////////////////////////////////

	go func() {
		var isAuthenticated bool
		authResponse := DecodeResponse(response)
		if authResponse.Status == 0 {
			connection.address.Host = authResponse.Address.Host
			connection.address.Port = authResponse.Address.Port
			isAuthenticated = true
			result.SuccessChannel <- isAuthenticated
		} else {
			result.FailureChannel <- errors.New("Connection is not authenticated" + connection.address.String())
		}
	}()

	return result
}