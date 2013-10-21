package main

import (
	"fmt"
	"net/rpc"
	"net"
	"os"
	"ztorage"
	"data"
	"time"
	"myRpc"
	"log"
	"strings"
	"strconv"
)

var server *myRpc.RpcServer
var store *ztorage.ViewersList
var c1 chan string
var c2 chan string
var c3 chan error
func main() {
	server = &myRpc.RpcServer{make(map[string]int)}
	store = ztorage.NewZapStore()
	c1 = make(chan string)
	c2 = make(chan string)
	c3 = make(chan error)
	go startRpc()
	go startStatsServer(store)
	go sendStats()
	for{
		select {
            case msg1 := <- c1:
                fmt.Println(msg1)
            case msg2 := <- c2:
                fmt.Println("clients: ", msg2)
			case err := <- c3:
				fmt.Println(err)
        }
	}
}

func startRpc(){
	fmt.Println("startRpc")
	rpc.Register(server)
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:60000")
	checkError(err)
	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("sever error: ", err)
		return
	}
	fmt.Println("connection established!")
	
	go func(ln net.Listener){
		for{
			c, err := ln.Accept()
			if err != nil {
			continue
			}
		
		go rpc.ServeConn(c)
		}
	}(ln)
}

func startStatsServer(store *ztorage.ViewersList){
	fmt.Println("start Server")
	gaddr, err := net.ResolveUDPAddr("udp", "224.0.1.130:10000")
	checkError(err)

	//go nrkViewers(store)
	//go tv2Viewers(store)
	//go top10()
	conn, err := net.ListenMulticastUDP("udp", nil, gaddr)
	checkError(err)
	for {		
		message(conn, store)
	}
}

func sendStats(){
	fmt.Println("Send Stats")
	for{
		time.Sleep(time.Second * 2)
		for i,v := range server.Cl{
			c2<-strings.Join([]string{i, strconv.Itoa(v)}, ", ")
			
			go func (rr int, addr string){
					refresh, err := time.ParseDuration(strings.Join([]string{strconv.Itoa(rr), "s"}, ""))
					if err != nil{
						c3<-err
						return
					}
					conn, err := net.Dial("tcp", addr)
					if err != nil{
						c3<-err
						return
					}
					for{
						time.Sleep(refresh)
						
						_, err = conn.Write([]byte(computeTop10()))
						
						if err != nil {
							c3<-err
							return
						}	
					}
			}(v, i)
			
			server.Delete(i)
		}
		/*for i,v := range server.Cl{
				go func (rr int, addr string){
					refresh,_ := time.ParseDuration(string(rr))
					for{
						time.Sleep(refresh * time.Second)
						
						conn, err := net.Dial("tcp", addr)
						if err != nil{
							return
						}
						top10 := computeTop10()
						_, err = conn.Write([]byte(top10))
						
						if err != nil {
							fmt.Println("an error occured: ", err)
							return
						}
						c <- top10
						conn.Close()		
					}
				}(v, i)
				server.Delete(i)
		}*/
	}
}

func message(conn *net.UDPConn, store *ztorage.ViewersList) {
	var buf [150]byte

	n, _,err := conn.ReadFromUDP(buf[0:])

	if err != nil {
		return
	}
	msg := string(buf[0:n])
	zap, _ := data.NewZap(msg)
	if err != nil {
		fmt.Println("error creating new zap", err)
		return
	}
	store.StoreZap(*zap)

	//fmt.Println(msg)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func nrkViewers(store *ztorage.ViewersList){
	for{
		time.Sleep(time.Second)
		fmt.Println("Nrk1 has: ", store.ComputeViewers("NRK1"), " viewers right now.")
	}
}

func tv2Viewers(store *ztorage.ViewersList){
	for{
		time.Sleep(time.Second)
		fmt.Println("Tv2 Norge has: ", store.ComputeViewers("TV2 Norge"), " viewers right now.")
	}
}

func top10(){
	for{
		msg := computeTop10()
		fmt.Println(msg)
		time.Sleep(time.Second * 1)
	}
	c1<-""
}

func computeTop10() string{
		if len(store.M) == 0{
			return "No watchers!"
		}
		lists := make([]string, 10)
		topList := make(map[string]int)
		for key, value := range store.M{
			topList[key] = value
		}
		//fmt.Println("\ntop 10:")
		for i := 0; i < 10; i++{
			maxKey := ""
			maxValue := 0
			for key, value := range topList{
				if value > maxValue{
					maxValue = value
					maxKey = key
				}
			}
			lists[i] = fmt.Sprint("\t", i + 1, "\t", maxKey, ":  ", maxValue, " viewers")
			//fmt.Println("\t", i + 1, "\t", maxKey, ":  ", maxValue, " viewers")
			delete(topList, maxKey)
		}
		return strings.Join(lists, "\n")
}
