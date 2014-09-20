package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func proxyConnection(conn net.Conn, destinations stringSlice) {

	var upstreamConns []io.Writer

	for _, addr := range destinations {
		forwardConn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Println(err)
			continue
		}
		upstreamConns = append(upstreamConns, forwardConn)
	}

	if len(upstreamConns) == 0 {
		log.Println("Could not connect to any destinations")
		return
	}

	upstreamWriter := io.MultiWriter(upstreamConns...)
	firstUpstreamCon := upstreamConns[0].(io.Reader)

	go func() {
		n, err := io.Copy(upstreamWriter, conn)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%d bytes copied from %s", n, conn.RemoteAddr().String())
	}()

	go func() {
		n, err := io.Copy(conn, firstUpstreamCon)
		defer conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("%d bytes copied to %s", n, conn.RemoteAddr().String())
	}()
}

func main() {
	var destinations stringSlice
	var source string

	flag.StringVar(&source, "s", "", "Host to listen on")
	flag.Var(&destinations, "d", "List of destination hosts")
	flag.Parse()

	if flag.NFlag() < 2 {
		panic("At least source and destination must be provided")
	}

	if source == "" {
		panic("You must provide a source host")
	}

	log.Printf("Source host: %s\n", source)
	log.Printf("Destination hosts: %v\n", destinations)

	listner, err := net.Listen("tcp", source)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			panic(err)
		}
		go proxyConnection(conn, destinations)
	}
}
