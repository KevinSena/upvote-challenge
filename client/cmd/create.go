package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a post",
	Long:  "Add a post with a tittle, description and author. The post start with 0 votes",
	Run: func(cmd *cobra.Command, args []string) {
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			log.Fatal(err)
		}
		desc, err := cmd.Flags().GetString("desc")
		if err != nil {
			log.Fatal(err)
		}
		author, err := cmd.Flags().GetString("author")
		if err != nil {
			log.Fatal(err)
		}
		res, err := ClientConnection.CreatePost(context.Background(), &pb.Post{
			Title:  title,
			Desc:   desc,
			Author: author,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(res)
	},
}

func init() {
	createCmd.Flags().StringP("title", "t", "", "The title of your post")
	createCmd.Flags().StringP("desc", "d", "", "The description of your post")
	createCmd.Flags().StringP("author", "a", "", "Your name")
	createCmd.MarkFlagRequired("title")
	createCmd.MarkFlagRequired("desc")
	createCmd.MarkFlagRequired("author")
	rootCmd.AddCommand(createCmd)
}
