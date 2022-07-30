package cmd

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Posts",
	Long:  `List all posts using grpc stream`,
	Run: func(cmd *cobra.Command, args []string) {
		stream, err := ClientConnection.ListPosts(context.Background(), &pb.Void{})
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
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
