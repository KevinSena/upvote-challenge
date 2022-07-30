package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete post",
	Long:  "Delete post specified by id",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatal(err)
		}

		res, err := ClientConnection.DeletePost(context.Background(), &pb.DeletePostRequest{
			Id: id,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.Msg)
	},
}

func init() {
	deleteCmd.Flags().StringP("id", "i", "", "Provide post identifier")
	deleteCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(deleteCmd)
}
