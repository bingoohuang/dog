package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bingoohuang/gou/htt"
	"github.com/bingoohuang/gou/lo"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	parseFlags()

	CreateAgApp().GoStart()
	select {}
}

func parseFlags() {
	conf := pflag.StringP("conf", "c", "./conf.toml", "config file path")
	help := pflag.BoolP("help", "h", false, "help")
	pflag.StringP("addr", "", ":9910", "listen address, eg :9910")
	pprofAddr := htt.PprofAddrPflag()
	pflag.Parse()
	args := pflag.Args()
	if len(args) > 0 {
		fmt.Printf("Unknown args %s\n", strings.Join(args, " "))
		pflag.PrintDefaults()
		os.Exit(0)
	}
	if *help {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if *conf != "" {
		cf, _ := homedir.Expand(*conf)
		if _, err := os.Stat(cf); err == nil {
			viper.SetConfigFile(*conf)
		}
	}

	htt.StartPprof(*pprofAddr)
	viper.SetEnvPrefix("DOG")
	viper.AutomaticEnv()

	// 绑定命令行参数，（优先级比配置文件高）
	lo.Err(viper.BindPFlags(pflag.CommandLine))

	logrus.SetLevel(logrus.InfoLevel)
}
