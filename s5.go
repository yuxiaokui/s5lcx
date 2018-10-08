package main

import socks5 "github.com/armon/go-socks5"

func main() {

	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on port 10080
	if err := server.ListenAndServe("tcp", "0.0.0.0:10080"); err != nil {
		panic(err)
	}

}
