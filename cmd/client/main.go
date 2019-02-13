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
	"context"
	"fmt"
	"github.com/ericbottard/knative-grpc/repeater"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s [address] [host/authority] [response-padding-size]\n", os.Args[0])
		os.Exit(2)
	}

	address := os.Args[1]
	authority := os.Args[2]
	padding, _ := strconv.Atoi(os.Args[3])

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithAuthority(authority))
	if err != nil {
		log.Panicf("error: %v", err)
	}
	defer unsafeClose(conn)
	client := repeater.NewRepeaterClient(conn)

	repeaterClient, err := client.Repeat(context.Background())
	if err != nil {
		log.Panicf("error: %v", err)
	}


	go func() {
		for {
			response, err := repeaterClient.Recv()
			fmt.Printf("response: %v, error: %v\n", response, err)
		}
	}()


	request := repeater.RepeatRequest{
		Quantity:            3,
		Content:             "hello",
		ResponsePaddingSize: int32(padding),
	}

	err = repeaterClient.Send(&request)
	if err != nil {
		log.Panicf("error: %v", err)
	}

	fmt.Println("Press Ctrl-C to complete...")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Terminating now...")
}

func unsafeClose(listener closeable) {
	err := listener.Close()
	if err != nil {
		panic(err)
	}
}

type closeable interface {
	Close() error
}
