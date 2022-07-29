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
	Id     primitive.ObjectID `bson:"_id,omitempty"`
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

func (s *Server) GetPost(_ context.Context, req *pb.GetPostRequest) (*pb.PostDB, error) {
	data := &Post{}
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Fail to convert hex to OID: %v", err)
	}
	res := postColl.FindOne(context.Background(), bson.M{"_id": id})

	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, "Post not found: %v", err)
	}

	result := &pb.PostDB{
		XId:    data.Id.Hex(),
		Title:  data.Title,
		Desc:   data.Desc,
		Votes:  data.Votes,
		Author: data.Author,
	}
	return result, nil
}

func (s *Server) Vote(_ context.Context, req *pb.VoteRequest) (*pb.PostDB, error) {
	data := &Post{}
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Fail to convert hex to OID: %v", err)
	}

	res := postColl.FindOneAndUpdate(context.Background(), bson.M{"_id": id}, bson.M{"$inc": bson.M{"votes": 1}}, options.FindOneAndUpdate().SetReturnDocument(1))
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, "Post not found: %v", err)
	}
	return &pb.PostDB{
		XId:    data.Id.Hex(),
		Title:  data.Title,
		Desc:   data.Desc,
		Votes:  data.Votes,
		Author: data.Author,
	}, nil
}

func (s *Server) CreatePost(_ context.Context, req *pb.Post) (*pb.PostDB, error) {
	data := &Post{
		Title:  req.Title,
		Desc:   req.Desc,
		Votes:  0,
		Author: req.Author,
	}
	fmt.Println(data)
	res, err := postColl.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Fail to create post: %v", err)
	}

	id := res.InsertedID.(primitive.ObjectID)

	return &pb.PostDB{
		XId:    id.Hex(),
		Title:  data.Title,
		Desc:   data.Desc,
		Votes:  data.Votes,
		Author: data.Author,
	}, nil
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
