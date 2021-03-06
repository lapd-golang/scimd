package cmd

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabbricadigitale/scimd/config"
	"github.com/fabbricadigitale/scimd/schemas/core"
	"github.com/fabbricadigitale/scimd/validation"
	"github.com/spf13/cobra"
)

func init() {
	scimd.AddCommand(printConfig)
	// (fixme) > should have the same flags of the root command to completely mimic its functioning and print the current config
}

var printConfig = &cobra.Command{
	Use:   "print-config",
	Short: "Print the current configuration",
	Long:  `Explicitly dump the current configuration objects for debugging purposes.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// spew.Dump(scimd.Flags())

		// Printing errors
		if len(config.Errors) > 0 {
			fmt.Fprintln(os.Stderr, validation.Errors(config.Errors))
			os.Exit(1)
		}

		// Ensure eventual custom config is ok and load it
		if err := config.Custom(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		dump := spew.NewDefaultConfig()
		dump.DisablePointerAddresses = true
		dump.SortKeys = true

		fmt.Fprintln(os.Stdout, "SERVICE PROVIDER CONFIG")
		dump.Dump(config.ServiceProviderConfig())

		fmt.Fprintln(os.Stdout, "\nSCHEMAS")
		dump.Dump(core.GetSchemaRepository().List())

		fmt.Fprintln(os.Stdout, "\nRESOURCE TYPES")
		dump.Dump(core.GetResourceTypeRepository().List())

		fmt.Fprintln(os.Stdout, "\nCONFIGURATION VALUES")
		dump.Dump(config.Values)
	},
	DisableAutoGenTag: true,
}
