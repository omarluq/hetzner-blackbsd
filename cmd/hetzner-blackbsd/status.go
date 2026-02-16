package main

import (
	"fmt"
	"io"
	"log/slog"
	"text/tabwriter"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"

	"github.com/omarluq/hetzner-blackbsd/internal/config"
	"github.com/omarluq/hetzner-blackbsd/internal/hcloud"
)

func newStatusCmd() *cobra.Command {
	var cmd cobra.Command
	cmd.Use = "status"
	cmd.Short = "Show BlackBSD build servers"
	cmd.RunE = runStatus
	return &cmd
}

func runStatus(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return err
	}

	client := hcloud.NewClient(cfg.HCloudToken)
	servers, err := client.ListServers(cmd.Context())
	if err != nil {
		return fmt.Errorf("list servers: %w", err)
	}

	if len(servers) == 0 {
		slog.Info("No BlackBSD servers found.")
		return nil
	}

	return printServers(cmd.OutOrStdout(), servers)
}

func printServers(output io.Writer, servers []*hcloudsdk.Server) error {
	tabWriter := tabwriter.NewWriter(output, 0, 0, 3, ' ', 0)

	if _, err := fmt.Fprintln(tabWriter, "ID\tNAME\tSTATUS\tIPv4\tRESCUE"); err != nil {
		return err
	}

	for _, server := range servers {
		rescue := "no"
		if server.RescueEnabled {
			rescue = "yes"
		}

		ipv4 := ""
		if server.PublicNet.IPv4.IP != nil {
			ipv4 = server.PublicNet.IPv4.IP.String()
		}

		if _, err := fmt.Fprintf(tabWriter, "%d\t%s\t%s\t%s\t%s\n",
			server.ID, server.Name, server.Status, ipv4, rescue); err != nil {
			return err
		}
	}

	if err := tabWriter.Flush(); err != nil {
		return err
	}

	_, err := fmt.Fprintf(output, "\nFound %d BlackBSD server(s).\n", len(servers))
	return err
}
