package main

import (
	pb "216/proto"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}

	defer conn.Close()

	grpcClient := pb.NewExpressionClient(conn)

	expression, err := grpcClient.Do(context.TODO(), &pb.Request{Messgae: "123321"})

	if err != nil {
		log.Println("failed invoked Expression: ", err)
	}

	fmt.Println("Expression: ", expression.Message)
}
