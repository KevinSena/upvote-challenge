package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/KevinSena/upvote-challenge/go/pb"
	"github.com/spf13/cobra"
)

var voteCmd = &cobra.Command{
	Use:   "vote",
	Short: "Vote in some post",
	Long:  "Vote in some post with the provided id",
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatal(err)
		}
		res, err := ClientConnection.Vote(context.Background(), &pb.VoteRequest{
			Id: id,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	},
}

func init() {
	voteCmd.Flags().StringP("id", "i", "", "Post identifier")
	voteCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(voteCmd)
}
