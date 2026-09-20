package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version    = "v1.0"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version of algorithm-notes ctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("algorithm-notes ctl version:", version)
		},
	}
)
