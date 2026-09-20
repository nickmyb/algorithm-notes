package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "algorithm-notes",
	Short: "A simple command line client for algorithm-notes.",
}

func execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
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
