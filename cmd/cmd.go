package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xylonx/config2go/config"
	"github.com/xylonx/config2go/converter"
	"github.com/xylonx/zapx"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = config.Setup(cfgFile, watch, sourceConfig, targetFile, targetPackage, tags)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

var cfgFile string

var (
	watch         string
	sourceConfig  string
	targetFile    string
	targetPackage string
	tags          string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.default.yaml", "specify config file path")

	// the config related flags
	rootCmd.Flags().StringVarP(&watch, "watch", "w", "", "whether listening file change event")
	rootCmd.Flags().StringVarP(&sourceConfig, "source", "s", "", "the source config file to be convert")
	rootCmd.Flags().StringVarP(&targetFile, "target", "t", "config/setting.go", "the generated target go file containing struct")
	rootCmd.Flags().StringVarP(&targetPackage, "package", "p", "", "the go package name of the generated target go file")
	rootCmd.Flags().StringVar(&tags, "tag", "mapstructure,json,yaml", "add the tag into the generated struct")
}

func Execute() error {
	return rootCmd.Execute()
}

func run() (err error) {
	err = converter.ConvertConfigFile(
		config.Config.SourceConfig,
		config.Config.TargetFile,
		config.Config.TargetPackage,
		config.Config.Tags,
	)
	if err != nil {
		zapx.Info("convert config file failed")
		return
	}
	return
}
