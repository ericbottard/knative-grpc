/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
	"strconv"
	"syscall"
)

type RepeaterService struct {
}

func (*RepeaterService) Repeat(server Repeater_RepeatServer) error {
	for true {
		request, err := server.Recv()
		if err == io.EOF {
			fmt.Println("Done: EOF")
			return nil
		}
		if err != nil {
			fmt.Printf("Aborting on Recv: error %v\n", err)
			return err
		}
		fmt.Printf("Got %v\n", request)
		for i := int64(0); i < request.Quantity; i++ {
			response := RepeatResponse{Content: request.Content, Padding: make([]byte, request.ResponsePaddingSize)}
			err := server.Send(&response)
			if err != nil {
				fmt.Printf("Aborting on Send: error %v\n", err)
				return err
			}
		}
	}
	return nil
}

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	serverPort, err := strconv.Atoi(port)
	if err != nil {
		log.Panicf("error: %v", err)
	}
	fmt.Printf("starting at port %d\n", serverPort)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Panicf("error: %v", err)
	}

	grpcServer := grpc.NewServer()
	RegisterRepeaterServer(grpcServer, &RepeaterService{})

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Panicf("error: %v", err)
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals
		grpcServer.GracefulStop()
	}()
}
