package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var h = flag.Bool("h", false, "this help")
var tcpPorts = flag.String("t", "", "The tcp ports to listen,Format: 8080,9000-9010 or single port.")
var udpPorts = flag.String("u", "", "The udp ports to listen,Format: 8080,9000-9010 or single port.")
var detail = flag.Bool("d", false, `This program can resp some useful infomation[REQINFO,SVRADDR,CLTADDR] to client,not just request info.`)

func useage() {
	fmt.Fprintf(os.Stderr, `Version: netecho/1.0.0 
Author: walker@walkerfeng.com	
A small echo program to reply client the same info as they send.
Usage: netecho [-d] [-t tcpports] [-u udpports]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = useage
	flag.Parse()
	if *h {
		flag.Usage()
		os.Exit(0)
	}

	tcpPortSlice := getPort(*tcpPorts)
	udpPortSlice := getPort(*udpPorts)

	if len(tcpPortSlice)+len(udpPortSlice) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var wait = make(chan struct{})
	for _, tport := range tcpPortSlice {
		tp := tport
		go tcpListen(tp)
	}

	for _, uport := range udpPortSlice {
		up := uport
		go udpListen(up)
	}
	<-wait
}

func getPort(ports string) []string {
	// Used to get port slice from args string.
	var portSlice []string
	for _, v := range strings.Split(ports, ",") {
		if v != "" {
			if strings.Contains(v, "-") {
				if strings.Count(v, "-") > 1 {
					log.Fatalf("Port Format Error,have more than 1 \"-\":%s", v)
					os.Exit(2)
				} else {
					startS := strings.Split(v, "-")[0]
					endS := strings.Split(v, "-")[1]
					start, _ := strconv.Atoi(startS)
					end, _ := strconv.Atoi(endS)
					if start > end {
						log.Fatalf("Port Format Error,start bigger than end \"-\":%s", v)
						os.Exit(3)
					} else {
						for i := start; i <= end; i++ {
							portSlice = append(portSlice, strconv.Itoa(i))
						}
					}
				}
			} else {
				portSlice = append(portSlice, v)
			}
		}
	}
	return portSlice
}

func tcpListen(port string) {
	//tcp listener.
	log.Printf("Begin to Listen TCP: %s", port)
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("%s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("%s", err)
		}

		go echo(conn)
	}
}

func echo(conn net.Conn) {
	// tcp echoer
	log.Printf("Get TCP connect from %s\n", conn.RemoteAddr())
	defer conn.Close()
	input := bufio.NewScanner(conn)

	for input.Scan() {

		log.Printf("Get TCP message from %s:%s\n", conn.RemoteAddr(), input.Text())
		if *detail {
			fmt.Fprintf(conn, "%s,%s,%s\n", input.Text(), conn.LocalAddr(), conn.RemoteAddr())
		} else {
			fmt.Fprintf(conn, "%s\n", input.Text())
		}
	}
	//io.Copy(c,c)
}

func udpListen(port string) {
	log.Printf("Begin to Listen UDP: %s\n", port)
	portInt, _ := strconv.Atoi(port)
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: portInt})
	if err != nil {
		log.Fatalf("%s", err)
	}
	for {
		buf := make([]byte, 1024)
		var tmpAddr *net.UDPAddr
		n, tmpAddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalf("%s", err)
		}
		if buf[n-1] == '\n' {
			log.Printf("Get UDP Message from %v: %s\n", tmpAddr, string(buf[:n-1]))
		} else {
			log.Printf("Get UDP Message from %v: %s\n", tmpAddr, string(buf[:n]))
		}

		if *detail {
			var tmpStr string
			if buf[n-1] == '\n' {
				tmpStr = fmt.Sprintf("%s,%s,%s\n", string(buf[:n-1]), udpConn.LocalAddr(), tmpAddr)
			} else {
				tmpStr = fmt.Sprintf("%s,%s,%s\n", string(buf[:n]), udpConn.LocalAddr(), tmpAddr)
			}
			udpConn.WriteToUDP([]byte(tmpStr), tmpAddr)
		} else {
			udpConn.WriteToUDP(buf, tmpAddr)
		}
	}
}
