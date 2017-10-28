package main

import (
	"config"
	"fmt"
	"log"
	"loghandle"
	"net"
	"proxy"
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
	log.SetPrefix("[info] ")
	var confobj, err = config.Load("../conf/default.json")
	if err != nil {
		log.SetPrefix("[error] ")
		fmt.Println(err)
		log.Println("配置文件读取错误")
		return
	}
	var max_cointegration chan int 
	var max_rconnection chan net.Conn 
	var max_wconnection chan net.Conn 
	max_cointegration = make(chan int, confobj.MaxCointegration)
        max_rconnection = make(chan net.Conn,confobj.MaxConnection)
        max_wconnection = make(chan net.Conn,confobj.MaxConnection)

	for _, value := range confobj.Backends {
		ipelement := value.Url()
		trueList = append(trueList, ipelement)
	}
	log.Println("Start tcp server....")
	listener, err := net.Listen("tcp", confobj.Host+":"+strconv.Itoa(int(confobj.Port)))
	if err != nil {
		log.SetPrefix("[error] ")
		log.Println("Error listening", err.Error())
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.SetPrefix("[error] ")
			log.Println("Error accepting", err.Error())
			return
		}
		//log.Println("waiting")
		go doServerStuff(conn, confobj, max_cointegration,max_rconnection,max_wconnection)

	}

}

func doServerStuff(conn net.Conn, confobj *config.Config, chanobj chan int,max_rconnchan chan net.Conn,max_wconnchan chan net.Conn) {
	for {
		totalinfo := []byte{}
		for {

			buf := make([]byte, 1024)
			len, err := conn.Read(buf)
			if len == 1024 {
				if err != nil {
					if err.Error() == "EOF" {
						log.Println("客户端断开连接")
						return
					} else {
						log.Println("Error accepting", err.Error())
						return

					}
				}
				totalinfo = util.BytesCombine(totalinfo, buf[:len])
			} else {
				if err != nil {
					if err.Error() == "EOF" {
						//log.Println("客户端断开连接")
						return
					} else {
						log.Println("Error accepting", err.Error())
						return

					}
				}
				totalinfo = util.BytesCombine(totalinfo, buf[:len])

				break
			}
		}
		cmd := strings.Split(string(totalinfo[:]), "\r\n")[2]
		if cmd == "COMMAND" {
			sayok(conn)
		} else {
			result := util.Cmdanalysis(strings.ToUpper(cmd))
			if result == 1 {
				ip := confobj.MasterHost + ":" + strconv.Itoa(int(confobj.MasterPort))
				chanobj <- 1
				go proxy.Handle(conn, ip, totalinfo[:], chanobj,max_wconnchan)
			} else {
				ip, _ := getIP()
				chanobj <- 1
				go proxy.Handle(conn, ip, totalinfo[:], chanobj,max_rconnchan)
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
