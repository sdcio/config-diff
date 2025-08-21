package cmd

import (
	"context"

	"github.com/sdcio/config-diff/pkg/configdiff"
	"github.com/sdcio/config-diff/pkg/configdiff/config"
	"github.com/sdcio/config-diff/pkg/types"
	"github.com/spf13/cobra"
)

// SchemaListCmd represents the list command
var SchemaListCmd = &cobra.Command{
	Use:          "list",
	Short:        "list available schemas",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var c *config.Config
		var cd *configdiff.ConfigDiff

		ctx := context.Background()

		opts := config.ConfigOpts{}
		c, err = config.NewConfig(opts)
		if err != nil {
			return err
		}

		cd, err = configdiff.NewConfigDiff(ctx, c)
		if err != nil {
			return err
		}

		schemas, err := cd.SchemasList(ctx)
		if err != nil {
			return err
		}

		outputSchemas := types.NewSchemaOutputList(schemas)

		err = WriteOutput(outputSchemas)
		if err != nil {
			return err
		}

		return nil

	},
}

func init() {
	schemaCmd.AddCommand(SchemaListCmd)
	EnableFlagAndDisableFileCompletion(SchemaListCmd)
}
