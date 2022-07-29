package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedPostServiceServer
}

type Post struct {
	Id     primitive.ObjectID `bson:"_id, omitempty"`
	Title  string             `bson:"title"`
	Desc   string             `bson:"desc"`
	Votes  int64              `bson:"votes"`
	Author string             `bson:"author"`
}

func (s *Server) ListPosts(nothing *pb.Void, stream pb.PostService_ListPostsServer) error {
	data := &Post{}

	cursor, err := postColl.Find(context.Background(), bson.D{{}})
	if err != nil {
		return status.Errorf(codes.Internal, "Internal Server Error: %v", err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Unavailable, "Failed to decode: %v", err)
		}

		stream.Send(&pb.PostDB{
			XId:    data.Id.Hex(),
			Title:  data.Title,
			Desc:   data.Desc,
			Votes:  data.Votes,
			Author: data.Author,
		})
	}

	return nil
}

var db *mongo.Client
var postColl *mongo.Collection

func main() {
	port := ":3001"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Fail to listen: %v", err)
	}

	uri := "mongodb://localhost:27017/"
	db, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Panic(err)
	}
	postColl = db.Database("upvote").Collection("posts")

	grpcServer := grpc.NewServer()
	pb.RegisterPostServiceServer(grpcServer, &Server{})

	go func() {
		if grpcErr := grpcServer.Serve(listener); grpcErr != nil {
			log.Fatal(grpcErr)
		}
	}()

	fmt.Printf("Server succesfully started on port %v", port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("\nStopping the server and killing db connection")
	grpcServer.Stop()
	listener.Close()
	db.Disconnect(context.Background())
	fmt.Println("\nDone")
}
