package main

import (
	"github.com/anhnguyentb/grpc-implement/global"
	"fmt"
	"github.com/anhnguyentb/grpc-implement/server"
)

func main() {

	//Load config
	if err := global.LoadConfig(); err != nil {
		panic(fmt.Errorf("Fatal due load config %s \n", err))
	}

	//Load logger
	if err := global.LoadLogger(false); err != nil {
		panic(fmt.Errorf("Fatal due load logger %s \n", err))
	}

	fmt.Println("Server starting ....")
	server.InitServer()
}
