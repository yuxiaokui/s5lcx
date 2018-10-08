package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

var lock sync.Mutex
var trueList []string
var ip string

var list string

var help = func() {
	fmt.Println("Usage for lcx tools, coded by xi4okv")
	fmt.Println("==========================================================")
	fmt.Println("Lcx LocalPort RemoteIp:RemotePort        # Port Forwarding")
	fmt.Println("Lcx ListenPort LocalPort                 # Listen ")
	fmt.Println("Lcx HostIp:HostPort TargetIp:TargetPort  # Slave")
	fmt.Println("==========================================================")
}

func main() {
	args := os.Args

	if len(args) != 3 || args == nil {
		help()
		os.Exit(0)
	}

	if !strings.ContainsAny(args[1], ":") && strings.ContainsAny(args[2], ":") {
		ip = "0.0.0.0:" + args[1]
		server(args[2])
	}

	if !strings.ContainsAny(args[1], ":") && !strings.ContainsAny(args[2], ":") {
		ip = "0.0.0.0:" + args[1]
		server("0.0.0.0:" + args[2])
	}

	if strings.ContainsAny(args[1], ":") && strings.ContainsAny(args[2], ":") {
		ip = args[1]
		server(args[2])
	}
}

func server(target string) {
	fmt.Printf("Listening %s\n", ip)
	lis, err := net.Listen("tcp", ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("建立连接错误:%v\n", err)
			continue
		}
		fmt.Println("working....\n")
		//fmt.Println(conn.RemoteAddr(), conn.LocalAddr())
		go handle(conn, target)
	}
}

func handle(sconn net.Conn, target string) {
	defer sconn.Close()
	ip := target
	dconn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Printf("连接%v失败:%v\n", ip, err)
		return
	}
	ExitChan := make(chan bool, 1)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		fmt.Println(dconn.Read)
		io.Copy(dconn, sconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)

	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(sconn, dconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	<-ExitChan
	dconn.Close()
}
