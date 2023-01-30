package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	Timeout time.Duration
	Thread  int
)
var rootCmd = &cobra.Command{
	Use:   "root",       //使用命令的别称
	Short: "short desc", //简短的提示
	Long:  "long desc",  //详细的提示
}

func Execute() {
	rootCmd.Execute()
}
