package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xylonx/config2go/config"
	"github.com/xylonx/config2go/converter"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "config2go",
	Short: "convert yaml, json, toml, etc. to go struct source code",
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
	rootCmd.Flags().StringVarP(&watch, "watch", "w", "true", "whether listening file change event")
	rootCmd.Flags().StringVarP(&sourceConfig, "source", "s", "", "the source config file to be convert")
	rootCmd.Flags().StringVarP(&targetFile, "target", "t", "config/setting.go", "the generated target go file containing struct")
	rootCmd.Flags().StringVarP(&targetPackage, "package", "p", "", "the go package name of the generated target go file")
	rootCmd.Flags().StringVar(&tags, "tag", "mapstructure,json,yaml", "add the tag into the generated struct")
}

func Execute() error {
	return rootCmd.Execute()
}

func run() (err error) {
	v := viper.New()
	v.SetConfigFile(config.Config.SourceConfig)
	if err := v.ReadInConfig(); err != nil {
		zapx.Error("read source config file failed", zap.Error(err))
		return err
	}

	if config.Config.Watch != "true" {
		if err = generateGoStruct(v); err != nil {
			return err
		}
		return nil
	}

	v.OnConfigChange(func(in fsnotify.Event) {
		if err = generateGoStruct(v); err != nil {
			return
		}
	})
	v.WatchConfig()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM)
	<-sig

	zapx.Info("stop config2go.")

	return nil
}

func generateGoStruct(v *viper.Viper) error {
	data := make(map[string]interface{})
	if err := v.Unmarshal(&data); err != nil {
		zapx.Error("unmarshal source config file failed", zap.Error(err))
		return err
	}

	parser := converter.NewMapParser(data)
	converter := converter.NewConverter(parser, converter.AppendAllFields(config.Config.Tags))

	if err := converter.Convert(config.Config.TargetPackage, config.Config.TargetFile); err != nil {
		return err
	}

	zapx.Info("genereate go struct successfully")
	return nil
}
