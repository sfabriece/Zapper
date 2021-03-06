package main

import (
	"fmt"
	"net"
	"os"
	"ztorage"
	"data"
	"time"
)

func main() {
	gaddr, err := net.ResolveUDPAddr("udp", "224.0.1.130:10000")
	checkError(err)
	store := ztorage.NewZapStore()
	//go nrkViewers(store)
	go tv2Viewers(store)
	//go top10(store)
	conn, err := net.ListenMulticastUDP("udp", nil, gaddr)
	checkError(err)
	for {		
		message(conn, store)
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

func top10(store *ztorage.ViewersList){
	for{
		time.Sleep(time.Second)
		topList := make(map[string]int)
		for key, value := range store.M{
			topList[key] = value
		}
		fmt.Println("\ntop 10:")
		for i := 0; i < 10; i++{
			maxKey := ""
			maxValue := 0
			for key, value := range topList{
				if value > maxValue{
					maxValue = value
					maxKey = key
				}
			}
			fmt.Println("\t", i + 1, "\t", maxKey, ":  ", maxValue, " viewers")
			delete(topList, maxKey)
		}
	}
}
