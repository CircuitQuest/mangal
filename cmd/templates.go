package cmd

import (
	"encoding/json"
	"strings"
	"text/template"

	"github.com/luevano/mangal/template/funcs"
	"github.com/luevano/mangal/theme/style"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func templatesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "templates",
		Short: "Name templates commands",
		Args:  cobra.NoArgs,
	}

	c.AddCommand(templatesFuncsCmd())
	c.AddCommand(templatesExecCmd())

	return c
}

func templatesFuncsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "funcs",
		Short: "Show available functions",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			for k, v := range funcs.Funcs {
				cmd.Println(style.Bold.Accent.Render(k))
				cmd.Println(style.Normal.Secondary.Render(v.Description))
				cmd.Println()
			}
		},
	}

	return c
}

func templatesExecCmd() *cobra.Command {
	templatesExecArgs := struct {
		Value string
	}{}

	c := &cobra.Command{
		Use:   "exec template...",
		Short: "Execute template",
		Args:  cobra.MinimumNArgs(1),
		// TODO: fix issue when using spaces in -v
		Run: func(cmd *cobra.Command, args []string) {
			tmpl, err := template.
				New("exec").
				Funcs(funcs.FuncMap).
				Parse(strings.Join(args, " "))
			if err != nil {
				errorf(cmd, err.Error())
			}

			var value map[string]any

			if err := json.Unmarshal([]byte(templatesExecArgs.Value), &value); err != nil {
				errorf(cmd, err.Error())
			}

			if err := tmpl.Execute(cmd.OutOrStdout(), value); err != nil {
				errorf(cmd, err.Error())
			}

			cmd.Println()
		},
	}

	exampleValue := struct {
		Title  string
		Number float64
	}{
		Title:  "Example Title",
		Number: 32.5,
	}
	marshalled := lo.Must(json.Marshal(exampleValue))
	c.Flags().StringVarP(&templatesExecArgs.Value, "value", "v", string(marshalled), "JSON object to use as value")

	return c
}
