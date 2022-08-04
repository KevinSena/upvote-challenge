package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Find a post",
	Long:  "Find a post with the provided id",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatal(err)
		}
		res, err := ClientConnection.GetPost(context.Background(), &pb.GetPostRequest{
			Id: id,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	},
}

func init() {
	getCmd.Flags().StringP("id", "i", "", "Post identifier")
	getCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(getCmd)
}
