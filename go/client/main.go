package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Fail to make client connection: %v", err)
	}
	defer conn.Close()
	client := pb.NewPostServiceClient(conn)
	stream, err := client.ListPosts(context.Background(), &pb.Void{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}
}
