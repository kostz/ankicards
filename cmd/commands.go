package cmd

import (
	"github.com/spf13/cobra"
	"os"

	"github.com/kostz/ankicards/internal"
)

var rootCmd = &cobra.Command{
	Use:   "ankicards",
	Short: "Ankicards generator",
	Long: ` 
A bunch of scripts generating Ankicards from different sources for studying German
`,
}

var extractVerbsFromImagesCmd = &cobra.Command{
	Use: "extractVerbsFromImages",
	Run: func(cmd *cobra.Command, args []string) {
		a := internal.NewApplication(
			internal.WithLLM(),
		)
		a.ExtractVerbsFromImages()
		a.WriteResult()
	},
}

var addVerbExamplesCmd = &cobra.Command{
	Use: "addVerbExamples",
	Run: func(cmd *cobra.Command, args []string) {
		a := internal.NewApplication(
			internal.WithLLM(),
		)
		a.LoadResult()
		a.AddVerbExamples()
		a.WriteResult()
	},
}

var makeAnkicardsCmd = &cobra.Command{
	Use: "makeAnkicards",
	Run: func(cmd *cobra.Command, args []string) {
		a := internal.NewApplication()
		a.LoadResult()
		a.MakeAnkicards()
	},
}

func init() {
	rootCmd.AddCommand(extractVerbsFromImagesCmd)
	rootCmd.AddCommand(addVerbExamplesCmd)
	rootCmd.AddCommand(makeAnkicardsCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
