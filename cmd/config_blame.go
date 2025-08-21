package cmd

import (
	"context"

	"github.com/sdcio/config-diff/pkg/configdiff"
	"github.com/sdcio/config-diff/pkg/configdiff/config"
	"github.com/sdcio/config-diff/pkg/types"
	sdcpb "github.com/sdcio/sdc-protos/sdcpb"
	"github.com/spf13/cobra"
)

var includeDefaults bool

// cconfigValidateCmd represents the validate command
var configBlameCmd = &cobra.Command{
	Use:          "blame",
	Short:        "blame config",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		ctx := context.Background()

		opts := config.ConfigOpts{}
		c, err := config.NewConfigPersistent(opts, optsP)
		if err != nil {
			return err
		}

		cd, err := configdiff.NewConfigDiffPersistence(ctx, c)
		if err != nil {
			return err
		}
		err = cd.InitTargetFolder(ctx)
		if err != nil {
			return err
		}

		sdcpbPath, err := sdcpb.ParsePath(path)
		if err != nil {
			return err
		}

		blameresult, err := cd.TreeBlame(ctx, includeDefaults, sdcpbPath)
		if err != nil {
			return err
		}

		bro := types.NewBlameResultOutput(blameresult)
		WriteOutput(bro)

		return nil
	},
}

func init() {
	configCmd.AddCommand(configBlameCmd)
	configBlameCmd.Flags().BoolVar(&includeDefaults, "include-defaults", false, "include the schema based default values in the output")
	AddPathPersistentFlag(configBlameCmd)
	EnableFlagAndDisableFileCompletion(configBlameCmd)
}
