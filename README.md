# netecho
Netecho is a small echo program to reply client the same info as they send.Writed in Golang.

# Install
Just use run the binary or complie with go(go build netecho.go).

# Usage

Usage: 

    netecho [-d] [-t tcpports] [-u udpports]

Options:

    -d	This program can resp some useful infomation[REQINFO,SVRADDR,CLTADDR] to client,not just request info.
    -h	this help
    -t string
    	The tcp ports to listen,Format: 8080,9000-9010 or single port.
    -u string
    	The udp ports to listen,Format: 8080,9000-9010 or single port.