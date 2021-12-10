package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

var Config *Setting = new(Setting)

// Setup the project config
// it will read the config from:
// 1. CLI flags
// 2. Special config file
//
// the priority of them is shown below:
// CLI flags > config file
func Setup(cfgFile, watch, sourceConfig, targetFile, targetPackage string, tags string) (err error) {
	v := viper.New()

	v.SetConfigType("yaml")

	if cfgFile != "" {
		var file *os.File
		file, err = os.Open(cfgFile)
		if err != nil {
			zapx.Error("open config file failed: ", zap.Error(err))
			return err
		}
		err = v.ReadConfig(file)
		if err != nil {
			zapx.Error("read config file failed", zap.Error(err))
			return
		}
	}

	// merge config
	cliConfigMap := make(map[string]interface{})
	switch watch {
	case "true":
		cliConfigMap["watch"] = "true"
	case "false":
		cliConfigMap["watch"] = "false"
	}
	if sourceConfig != "" {
		cliConfigMap["source_config"] = sourceConfig
	}
	if targetFile != "" {
		cliConfigMap["target_file"] = targetFile
	}
	if targetPackage != "" {
		cliConfigMap["target_package"] = targetPackage
	}

	cliConfigMap["tags"] = mergeTags(tags, v.GetString("tag"))

	err = v.MergeConfigMap(cliConfigMap)
	if err != nil {
		zapx.Error("merge cli config failed", zap.Error(err))
		return err
	}

	err = v.Unmarshal(Config)
	if err != nil {
		zapx.Error("unmarshal config failed")
	}

	return nil
}

func mergeTags(tags1, tags2 string) []string {
	t1 := strings.Split(tags1, ",")
	t2 := strings.Split(tags2, ",")

	t3 := make([]string, 0, len(t1)+len(t2))
	tagMaps := make(map[string]bool)
	for i := range t1 {
		if _, ok := tagMaps[t1[i]]; !ok {
			tagMaps[t1[i]] = true
			t3 = append(t3, t1[i])
		}
	}
	for i := range t2 {
		if _, ok := tagMaps[t2[i]]; !ok {
			tagMaps[t2[i]] = true
			t3 = append(t3, t2[i])
		}
	}
	return t3
}
