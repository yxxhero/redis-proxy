package proxy

import (
	"io"
	"log"
	"net"
	"time"
)

func Handle(sconn net.Conn, ip string, buf []byte) {
	dconn, err := net.DialTimeout("tcp", ip, 2*time.Second)
	if err != nil {
		log.Printf("连接%v失败:%v\n", ip, err)
		return
	}
	func(sconn net.Conn, dconn net.Conn, buf []byte) {
		_, write_err := dconn.Write(buf)
		if write_err != nil {
			log.Println("缓冲写入错误:%v",write_err)
		}
		_, copy_err := io.Copy(sconn, dconn)
		if copy_err != nil {
			log.Println("io复制错误:%v",copy_err)
		}
	}(sconn, dconn, buf)
}
