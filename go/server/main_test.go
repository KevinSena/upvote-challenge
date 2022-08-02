package main_test

import (
	"context"
	"encoding/hex"
	"io"
	"log"
	"testing"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var connection pb.PostServiceClient
var postMongo *mongo.Collection
var mongoCtx context.Context

func init() {
	conn, err := grpc.Dial("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Fail to make client connection, verify if server is running: %v", err)
	}
	connection = pb.NewPostServiceClient(conn)

	uri := "mongodb://localhost:27017/"
	mongoCtx = context.Background()
	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Panic(err)
	}
	postMongo = db.Database("upvote").Collection("posts")
}

func TestList(t *testing.T) {
	stream, err := connection.ListPosts(context.Background(), &pb.Void{})
	if err != nil {
		t.Errorf("Occur an unexpected error on connect to server or db: %v", err)
	}
	var got int64 = 0
	want, _ := postMongo.CountDocuments(mongoCtx, bson.M{})
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
		}
		got++
		_, err = hex.DecodeString(res.XId)
		if err != nil {
			t.Errorf("the string %v isn't a Hex", res.XId)
		}
	}
	if got != want {
		t.Errorf("Expected to get %v elements and got %v", want, got)
	}
}

func TestGet(t *testing.T) {
	data := bson.M{
		"Title":  "Get test",
		"Desc":   "test if the get function is correct",
		"Author": "TestGet",
	}
	created, _ := postMongo.InsertOne(context.Background(), data)

	id := created.InsertedID.(primitive.ObjectID).Hex()

	res, err := connection.GetPost(context.Background(), &pb.GetPostRequest{
		Id: id,
	})
	if err != nil {
		t.Errorf("Occur an unexpected error on connect to server or db: %v", err)
	}

	if res.XId != id {
		t.Errorf("Expecte to get %v but got %v", created, res)
	}
	if res.Title != "Get test" {
		t.Errorf("Expecte to get %v but got %v", created, res)
	}
	if res.Author != "TestGet" {
		t.Errorf("Expecte to get %v but got %v", created, res)
	}

	wrongRes, _ := connection.GetPost(context.Background(), &pb.GetPostRequest{
		Id: "fwe984fsd",
	})
	if wrongRes != nil {
		t.Error("Expect to return nil with false id")
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	postMongo.DeleteOne(context.Background(), bson.M{"_id": oid})
}

func TestVote(t *testing.T) {
	data := bson.M{
		"Title":  "Vote test",
		"Desc":   "test if the vote function is correct",
		"Author": "TestVote",
	}
	created, _ := postMongo.InsertOne(context.Background(), data)
	id := created.InsertedID.(primitive.ObjectID).Hex()

	want := 10
	var count int64
	for i := 0; i < want; i++ {
		res, err := connection.Vote(context.Background(), &pb.VoteRequest{
			Id: id,
		})
		if err != nil {
			t.Errorf("Occur an unexpected error on connect to server or db: %v", err)
		}
		count = res.Votes
	}

	if count != int64(want) {
		t.Errorf("Excpect to receive %v votes but got %v", want, count)
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	postMongo.DeleteOne(context.Background(), bson.M{"_id": oid})
}

func TestCreate(t *testing.T) {
	want := &pb.Post{
		Title:  "Create test",
		Desc:   "test if the Create function is correct",
		Author: "TestCreate",
		Votes:  0,
	}
	created, err := connection.CreatePost(context.Background(), &pb.Post{
		Title:  "Create test",
		Desc:   "test if the Create function is correct",
		Author: "TestCreate",
	})
	if err != nil {
		t.Errorf("Unexpected error on create post: %v", err)
	}
	if created.Votes != want.Votes {
		t.Errorf("Excpect to receive %v votes but got %v", want.Votes, created.Votes)
	}
	if created.Author != want.Author {
		t.Errorf("Excpect to receive %v but got %v", want.Author, created.Author)
	}
	if created.Desc != want.Desc {
		t.Errorf("Excpect to receive %v but got %v", want.Desc, created.Desc)
	}
	if created.Title != want.Title {
		t.Errorf("Excpect to receive %v but got %v", want.Title, created.Title)
	}
	_, err = hex.DecodeString(created.XId)
	if err != nil {
		t.Errorf("The id created isn't a hexadecimal")
	}

	oid, _ := primitive.ObjectIDFromHex(created.XId)
	postMongo.DeleteOne(context.Background(), bson.M{"_id": oid})

	_, err = connection.CreatePost(context.Background(), &pb.Post{
		Desc:   "test if the vote function is correct",
		Author: "TestVote",
	})
	if err == nil {
		t.Error("Expect to return an error if title is not passed")
	}
	_, err = connection.CreatePost(context.Background(), &pb.Post{
		Title:  "Create test",
		Author: "TestVote",
	})
	if err == nil {
		t.Error("Expect to return an error if desc is not passed")
	}
	_, err = connection.CreatePost(context.Background(), &pb.Post{
		Title: "Create test",
		Desc:  "test if the Create function is correct",
	})
	if err == nil {
		t.Error("Expect to return an error if author is not passed")
	}
}

func TestUpdate(t *testing.T) {
	data := bson.M{
		"Title":  "test",
		"Desc":   "test if the function is correct",
		"Author": "Test",
	}
	created, _ := postMongo.InsertOne(context.Background(), data)
	id := created.InsertedID.(primitive.ObjectID).Hex()
	newPost := &pb.PostDB{
		XId:    id,
		Title:  "Update test",
		Desc:   "test if the update function is correct",
		Author: "TestUpdate",
	}

	updated, err := connection.UpdatePost(context.Background(), newPost)
	if err != nil {
		t.Errorf("unexpected error on update post: %v", err)
	}
	if updated.Votes != 0 {
		t.Error("Votes cant be changed by update")
	}
	if updated.Author != newPost.Author {
		t.Errorf("Expected to get %v but got %v", updated.Author, updated.Author)
	}
	if updated.Title != newPost.Title {
		t.Errorf("Expected to get %v but got %v", updated.Title, updated.Title)
	}
	if updated.Desc != newPost.Desc {
		t.Errorf("Expected to get %v but got %v", updated.Desc, updated.Desc)
	}

	oid, _ := primitive.ObjectIDFromHex(id)
	postMongo.DeleteOne(context.Background(), bson.M{"_id": oid})
}

func TestDelete(t *testing.T) {
	data := bson.M{
		"Title":  "delete test",
		"Desc":   "test if the Delete function is correct",
		"Author": "TestDelete",
	}
	created, _ := postMongo.InsertOne(context.Background(), data)
	id := created.InsertedID.(primitive.ObjectID).Hex()

	deleted, err := connection.DeletePost(context.Background(), &pb.DeletePostRequest{Id: id})
	if err != nil {
		t.Errorf("Unexpected error on delete method: %v", err)
	}
	if deleted.Msg == "" {
		t.Error("Expect a deletion message")
	}

	_, err = connection.DeletePost(context.Background(), &pb.DeletePostRequest{Id: id})
	if err == nil {
		t.Error("Expected an error when delete post already deleted")
	}
}
