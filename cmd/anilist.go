package cmd

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/client/anilist"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func init() {
	subcommands = append(subcommands, anilistCmd)
}

var anilistCmd = &cobra.Command{
	Use:   "anilist",
	Short: "Anilist auth commands",
	Args:  cobra.NoArgs,
}

func init() {
	anilistCmd.AddCommand(anilistAuthCmd)
}

// TODO: replace this with a mini TUI
var anilistAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with anilist",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(cmd.InOrStdin())

		cmd.Print("ID\n> ")
		id, err := reader.ReadString('\n')
		if err != nil {
			errorf(cmd, err.Error())
		}
		id = strings.TrimSpace(id)

		cmd.Print("Secret\n> ")
		secret, err := reader.ReadString('\n')
		if err != nil {
			errorf(cmd, err.Error())
		}
		secret = strings.TrimSpace(secret)

		authURL := fmt.Sprint("https://anilist.co/api/v2/oauth/authorize?client_id=", id, "&response_type=code&redirect_uri=https://anilist.co/api/v2/oauth/pin")
		if err := open.Start(authURL); err != nil {
			errorf(cmd, err.Error())
		}

		cmd.Print("Code\n> ")
		code, err := reader.ReadString('\n')
		if err != nil {
			errorf(cmd, err.Error())
		}
		code = strings.TrimSpace(code)

		err = anilist.Anilist.Authorize(context.Background(), libmangal.AnilistLoginCredentials{
			ID:     id,
			Secret: secret,
			Code:   code,
		})

		if err != nil {
			errorf(cmd, err.Error())
		}

		successf(cmd, "Authorized with the Anilist")
	},
}

func init() {
	anilistCmd.AddCommand(anilistLogoutCmd)
}

var anilistLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from anilist",
	Run: func(cmd *cobra.Command, args []string) {
		if err := anilist.Anilist.Logout(); err != nil {
			errorf(cmd, err.Error())
		}

		successf(cmd, "Logged out from Anilist")
	},
}
