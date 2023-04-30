package main

import (
	"github.com/snowmerak/rmq/gen/proto/message"
	msgsvr "github.com/snowmerak/rmq/lib/server/message"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("invalid args")
	}

	addr := os.Args[1]

	server := grpc.NewServer()
	message.RegisterMessageQueueServer(server, &msgsvr.Server{})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
