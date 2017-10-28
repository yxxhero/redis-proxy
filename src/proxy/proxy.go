package proxy

import (
	"log"
	"net"
	"time"
	"util"
)

func Handle(sconn net.Conn, ip string, buf []byte, exitchanobj chan int, max_conn chan net.Conn) {
	select {
	case dconn := <-max_conn:
		func(sconn net.Conn, dconn net.Conn, buf []byte, exitobj chan int, max_conn chan net.Conn) {
			_, write_err := dconn.Write(buf)
			if write_err != nil {
				log.Println("缓冲写入错误:%v", write_err.Error())
				return
			}
			util.Copy(sconn, dconn, max_conn)

			<-exitchanobj
		}(sconn, dconn, buf, exitchanobj, max_conn)
	default:
		dconn, err := net.DialTimeout("tcp", ip, 2*time.Second)
		if err != nil {
			log.Printf("连接%v失败:%v\n", ip, err)
			return
		} 
		func(sconn net.Conn, dconn net.Conn, buf []byte, exitobj chan int, max_conn chan net.Conn) {
			_, write_err := dconn.Write(buf)
			if write_err != nil {
				log.Println("缓冲写入错误:%v", write_err.Error())
				return
			}
			util.Copy(sconn, dconn, max_conn)

			<-exitchanobj
		}(sconn, dconn, buf, exitchanobj, max_conn)

	}
}
