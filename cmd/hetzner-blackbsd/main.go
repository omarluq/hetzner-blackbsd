// Package main is the entry point for blackbsd.
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"github.com/omarluq/hetzner-blackbsd/internal/vinfo"
)

const defaultConfigFile = "blackbsd.yml"

var cfgFile string

func newRootCmd() *cobra.Command {
	var cmd cobra.Command
	cmd.Use = "hetzner-blackbsd"
	cmd.Short = "BlackBSD ISO build pipeline on Hetzner Cloud"
	cmd.Long = `BlackBSD builds NetBSD-based security ISO images on Hetzner Cloud ephemeral servers.`
	cmd.Example = `  # List build servers
  hetzner-blackbsd status

  # Destroy orphaned build servers
  hetzner-blackbsd destroy

  # Show version
  hetzner-blackbsd version`
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", defaultConfigFile, "config file path")
	return &cmd
}

var rootCmd = newRootCmd()

func init() {
	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newDestroyCmd())
	rootCmd.AddCommand(newVersionCmd())
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rootCmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")

	fangOpts := []fang.Option{
		fang.WithVersion(vinfo.String()),
	}

	cobra.CheckErr(fang.Execute(ctx, rootCmd, fangOpts...))
}
