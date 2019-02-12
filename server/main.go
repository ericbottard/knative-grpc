package main

import (
	"fmt"
	. "github.com/ericbottard/knative-grpc/repeater"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type RepeaterService struct {
}

func (*RepeaterService) Repeat(server Repeater_RepeatServer) error {
	for true {
		input, err := server.Recv()
		if err == io.EOF {
			fmt.Println("Aborting on Recv: EOF")
			return nil
		}
		if err != nil {
			fmt.Printf("Aborting on Recv: error %v\n", err)
			return err
		}
		fmt.Printf("Got %v\n", input)
		for i := int64(0); i < input.Quantity; i++ {
			err := server.Send(&RepeatResponse{Content: input.Content})
			if err != nil {
				fmt.Printf("Aborting on Send: error %v\n", err)
				return err
			}
		}
	}
	return nil
}

func main() {
	port := "9999"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	fmt.Printf("starting at port %s", port)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panicf("error: %v", err)
	}

	server := grpc.NewServer()
	service := RepeaterService{}
	RegisterRepeaterServer(server, &service)

	err = server.Serve(listener)
	if err != nil {
		log.Panicf("error: %v", err)
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals
		server.GracefulStop()
	}()
}