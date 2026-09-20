package main

import (
	"github.com/spf13/cobra"
)

func newRefresh() *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "Refresh all document",
		Run: func(cmd *cobra.Command, args []string) {
			refresh()
		},
	}
}

func refresh() {
	// 站点相关的两步（渲染第二章、拷题解进第四章）依赖 website/ 目录，
	// 已随 ctl/render_website.go 一起停用，现在 refresh 只重建 README。
	// copyLackFile()
	// buildBookMenu()
	buildREADME()
	// buildChapterTwo(true)
}
