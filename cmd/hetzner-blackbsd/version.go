package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/omarluq/hetzner-blackbsd/internal/vinfo"
)

func newVersionCmd() *cobra.Command {
	var cmd cobra.Command
	cmd.Use = "version"
	cmd.Short = "Show version information"
	cmd.RunE = func(c *cobra.Command, _ []string) error {
		_, err := fmt.Fprintln(c.OutOrStdout(), vinfo.String())
		return err
	}
	return &cmd
}
