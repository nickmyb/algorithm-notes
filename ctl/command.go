package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "algorithm-notes",
	Short: "A simple command line client for algorithm-notes.",
	// cobra 默认自己打一遍错误，下面 execute() 又打一遍，同一句话会出现两次。
	// 交给 execute() 统一处理，顺便写到 stderr 而不是 stdout。
	SilenceErrors: true,
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		versionCmd,
		newBuildCommand(),
		newNewCommand(),
		newRefresh(),
		// PDF 电子书生成依赖 website/ 目录，已随 ctl/pdf.go 一起停用。
		// newPDFCommand(),
	)
}
