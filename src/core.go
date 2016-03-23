package src
import "fmt"


/*
	Hazelcast Core Objects
 */


type Address struct {
	Host string
	Port int
}
func (address Address) String() string {
	return fmt.Sprintf("Address: host=%s, port=%d", address.Host, address.Port)
}

