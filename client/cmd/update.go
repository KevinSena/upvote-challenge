package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Edit some post",
	Long:  "Edit some post by id. You can change the title, description and author",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatal(err)
		}
		title, _ := cmd.Flags().GetString("title")
		desc, _ := cmd.Flags().GetString("desc")
		author, _ := cmd.Flags().GetString("author")

		res, err := ClientConnection.UpdatePost(context.Background(), &pb.PostDB{
			XId:    id,
			Title:  title,
			Desc:   desc,
			Author: author,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	},
}

func init() {
	updateCmd.Flags().StringP("id", "i", "", "Provide a id to find post")
	updateCmd.Flags().StringP("title", "t", "", "The title of your post")
	updateCmd.Flags().StringP("desc", "d", "", "The description of your post")
	updateCmd.Flags().StringP("author", "a", "", "Your name")
	updateCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(updateCmd)
}
