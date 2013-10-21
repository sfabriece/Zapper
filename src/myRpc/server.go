package myRpc
import (
	//"fmt"
	"sync"
)
var lock sync.Mutex
type Client struct{
	RefreshRate int
	Address string
}
type RpcServer struct{
	Cl map[string]int
}

func (this *RpcServer) add(cl Client){
	lock.Lock()
	this.Cl[cl.Address] = cl.RefreshRate
	lock.Unlock()
}

func (this *RpcServer) Delete(address string){
	lock.Lock()
	delete(this.Cl, address)
	lock.Unlock()
}

func (this *RpcServer) Subscribe(args Client, reply *int) error {
	this.add(args)
//	fmt.Println(args.Address)
//	*reply = 1
	return nil
}