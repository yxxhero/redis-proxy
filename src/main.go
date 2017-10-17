package main

import (
	"config"
	"fmt"
	"log"
	"loghandle"
	"net"
	"proxy"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"util"
)

var trueList []string
var ip string
var lock sync.Mutex

func main() {
	loghandle.Init("redisproxy.log")
	log.SetPrefix("[info]")
	var confobj, err = config.Load("conf/default.json")
	if err != nil {
		log.SetPrefix("[error]")
		fmt.Println(err)
		log.Println("配置文件读取错误")
		return
	}
	for _, value := range confobj.Backends {
		ipelement := value.Url()
		trueList = append(trueList, ipelement)
	}
	log.Println("Start tcp server....")
	listener, err := net.Listen("tcp", confobj.Host+":"+strconv.Itoa(int(confobj.Port)))
	if err != nil {
		log.SetPrefix("[error]")
		log.Println("Error listening", err.Error())
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.SetPrefix("[error]")
			log.Println("Error accepting", err.Error())
			return
		}
		go doServerStuff(conn, confobj)

	}

}

func doServerStuff(conn net.Conn, confobj *config.Config) {
	for {

		buf := make([]byte, 1024)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return
		}
		log.Println(string(buf[:len]))
		log.Println(reflect.TypeOf(buf[:len]))
		cmd := strings.Split(string(buf[:len]), "\r\n")[2]
		if cmd == "COMMAND" {
			sayok(conn)
		} else {
			result := util.Cmdanalysis(strings.ToLower(cmd))
			log.Println("cmd is ", result)
			if result == 1 {
				ip := confobj.MasterHost + ":" + strconv.Itoa(int(confobj.MasterPort))
				log.Println(ip)
				go proxy.Handle(conn, ip, buf[:len])
			} else {
				ip, _ := getIP()
				log.Println(ip)
				go proxy.Handle(conn, ip, buf[:len])
			}
		}
	}

}

func sayok(to net.Conn) {
	obuf := []byte{'+', 'o', 'k', '\r', '\n'}
	_, err := to.Write(obuf)

	util.CheckError(err)
}
func getIP() (string, bool) {
	lock.Lock()
	defer lock.Unlock()

	if len(trueList) < 1 {
		return "", false
	}
	ip := trueList[0]
	trueList = append(trueList[1:], ip)
	return ip, true
}
