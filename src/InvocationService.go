package src

type Invocation struct {
	message *ClientMessage
	partitionId int32
	address *Address
	connection *ClientConnection
	messages chan *ClientMessage
}

func NewInvocation(request *ClientMessage) {
	invocation := new(Invocation)
	invocation.message = request
}

type InvocationService struct {

}