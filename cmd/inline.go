package cmd

import (
	"context"

	"github.com/luevano/mangal/anilist"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/inline"
	"github.com/luevano/mangal/provider/manager"
	"github.com/luevano/libmangal"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var inlineArgs = inline.InlineArgs{}

func init() {
	subcommands = append(subcommands, inlineCmd)

	inlineCmd.Flags().StringVarP(&inlineArgs.Query, "query", "q", "", "Query to search")
	inlineCmd.Flags().StringVarP(&inlineArgs.Provider, "provider", "p", "", "Load provider by tag")
	inlineCmd.Flags().BoolVarP(&inlineArgs.Download, "download", "d", false, "Load provider by tag")
	inlineCmd.Flags().BoolVarP(&inlineArgs.JSON, "json", "j", false, "Load provider by tag")

	inlineCmd.MarkFlagRequired("provider")
	inlineCmd.MarkFlagRequired("query")
	inlineCmd.MarkFlagsOneRequired("download", "json")
	inlineCmd.MarkFlagsMutuallyExclusive("download", "json")
	inlineCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)
}

var inlineCmd = &cobra.Command{
	Use:     "inline",
	Short:   "Run mangal in inline mode",
	GroupID: groupMode,
	Args:    cobra.NoArgs,
	// TODO: probably move all this logic to inline/inline.go itself, not sure why this pattern was used (and I followed)
	Run: func(cmd *cobra.Command, args []string) {
		loaders, err := manager.Loaders()
		if err != nil {
			errorf(cmd, err.Error())
		}

		loader, ok := lo.Find(loaders, func(loader libmangal.ProviderLoader) bool {
			return loader.Info().ID == inlineArgs.Provider
		})

		if !ok {
			errorf(cmd, "provider with ID %q not found", inlineArgs.Provider)
		}

		client, err := client.NewClient(context.Background(), loader)
		if err != nil {
			errorf(cmd, err.Error())
		}

		var options inline.Options

		options.InlineArgs = inlineArgs
		options.Client = client
		options.Anilist = anilist.Client

		if err := inline.Run(context.Background(), options); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
