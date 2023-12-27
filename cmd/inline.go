package cmd

import (
	"context"

	"github.com/luevano/mangal/anilist"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/inline"
	"github.com/spf13/cobra"
)

var inlineArgs = inline.InlineArgs{}

func init() {
	subcommands = append(subcommands, inlineCmd)

	inlineCmd.PersistentFlags().StringVarP(&inlineArgs.Query, "query", "q", "", "Query to search")
	inlineCmd.PersistentFlags().StringVarP(&inlineArgs.Provider, "provider", "p", "", "Load provider by tag")
	inlineCmd.PersistentFlags().StringVarP(&inlineArgs.MangaSelector, "manga-selector", "m", "", "Manga selector")

	inlineCmd.MarkFlagRequired("provider")
	inlineCmd.MarkFlagRequired("query")
	inlineCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)
}

var inlineCmd = &cobra.Command{
	Use:     "inline",
	Short:   "Inline mode",
	GroupID: groupMode,
	Args:    cobra.NoArgs,
}

func init() {
	inlineCmd.AddCommand(inlineJSONCmd)
}

var inlineJSONCmd = &cobra.Command{
	Use:     "json",
	Short:   "Output search results in JSON",
	Args:    cobra.NoArgs,
	// TODO: change to RunE
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClientByID(context.Background(), inlineArgs.Provider)
		if err != nil {
			errorf(cmd, err.Error())
		}

		var options inline.Options

		options.InlineArgs = inlineArgs
		options.Client = client
		options.Anilist = anilist.Client

		if err := inline.RunJSON(context.Background(), options); err != nil {
			errorf(cmd, err.Error())
		}
	},
}

func init() {
	inlineCmd.AddCommand(inlineDownloadCmd)
}

var inlineDownloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Download manga",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		errorf(cmd, "unimplemented")
	},
}
