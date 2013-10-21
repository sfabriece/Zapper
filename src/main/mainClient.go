package main

import (
	"fmt"
	"net/rpc"
	"os"
	"log"
	"net"
	"myRpc"
)

type Stats struct{}

func (this *Stats) Receive(msg string, reply *int) error{
	fmt.Println(msg)
	*reply = 1
	return nil
}

var myAdress = "127.0.0.1:61000"
var c chan string
func main (){
	c = make(chan string)
	go startServer()
	go startClient()
	for{
		select{
			case msg := <-c:
				fmt.Println(msg)
		}
	}
}

func startServer(){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", myAdress)
	if err != nil{
		log.Fatal("client server error: ", err)
		checkError(err)
	}
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil{
		log.Fatal("client server error: ", err)
		checkError(err)
	}
	
	for{
		conn, err := ln.Accept()
		if err != nil{
			continue
		}
		handleConn(conn)
	}
}

func handleConn(conn net.Conn){
	var buf [512]byte
	for{
		n, err := conn.Read(buf[0:])
		if err != nil {
			continue
		}
		c<- string(buf[0:n])
	}
}

func startClient(){
	client, err := rpc.Dial("tcp", "localhost:60000")
	if err != nil{
		log.Fatal("dialing: ", err)
		checkError(err)
	}
	
	rr := 5
	var reply int
	err = client.Call("RpcServer.Subscribe", myRpc.Client{rr, myAdress}, &reply)
	if err != nil {
		fmt.Println("an error occured: ", err)
		checkError(err)
	}
	c<-"recieved connection!"
	client.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}