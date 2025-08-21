package cmd

import (
	"os"

	"github.com/sdcio/config-diff/pkg/configdiff/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var optsP = config.ConfigPersistentOpts{}
var jsonOutput bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "config-diff",
	Short: "A CLI tool to interact with NOS configs based on YANG schemas",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd.AddCommand(makeCompletionCmd(rootCmd))
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func AddTargetPersistentFlag(c *cobra.Command) error {
	c.PersistentFlags().StringVarP(&targetName, "target", "t", "", "the target to use")
	err := c.MarkPersistentFlagRequired("target")
	if err != nil {
		return err
	}

	// Register autocompletion for the flag
	err = c.RegisterFlagCompletionFunc("target", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		opts := config.ConfigOpts{}
		c, err := config.NewConfigPersistent(opts, optsP)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		entries, err := os.ReadDir(c.TargetBasePath())
		if err != nil {
			log.Fatal(err)
		}

		result := make([]string, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				result = append(result, entry.Name())
			}
		}

		return result, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		return err
	}

	return nil
}

func EnableFlagAndDisableFileCompletion(cmd *cobra.Command) {
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			// Skip if already set
			if !f.Changed {
				completions = append(completions, "--"+f.Name)
			}
		})
		// skip just --help
		if len(completions) == 1 && completions[0] == "--help" {
			completions = nil
		}

		return completions, cobra.ShellCompDirectiveNoFileComp
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Return output in JSON instead of text")
}
