package data

import (
	"net"
	"strings"
	"time"
	//"fmt"
)

const timeLayout = "2006/01/02, 15:04:05"
var counter int = 0
type Zap struct{
	time time.Time
	IP net.IP
	FromCH string
	ToCH string
}
type ZapError struct {
	msg string
}
func (e *ZapError) Error() string {
	return e.msg
}
func NewZap(msg string) (*Zap, error){
	if msg == ""{
		return new(Zap), &ZapError{"input is empty or nil"}
	}
	var params []string = strings.Split(msg, ",")

	if strings.HasPrefix(params[3], "Mute_Status") || len(params) < 5{
		return new(Zap), &ZapError{"not a channel change"}
	}
	//fmt.Println(params)
	for  i := 0; i < len(params); i++ {
		params[i] = strings.TrimSpace(params[i])
	}
	t, _ := time.Parse(timeLayout, strings.Join([]string{params[0], ", ", params[1]}, ""))
	IP := net.ParseIP(params[2])
	FromCH := params[3]
	ToCH := params[4]
	return &Zap{t, IP, FromCH, ToCH}, nil
}

func (z *Zap) String() string{
	return strings.Join([]string{z.time.String()," \n", z.IP.String(), "\nFrom: ", z.FromCH, "\nTO: ", z.ToCH}, "")
}
func (z *Zap) Duration(provided Zap) time.Duration{
	return z.time.Sub(provided.time)
}
