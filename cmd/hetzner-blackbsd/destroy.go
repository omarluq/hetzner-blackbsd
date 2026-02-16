package main

import (
	"fmt"
	"io"
	"log/slog"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"

	"github.com/omarluq/hetzner-blackbsd/internal/config"
	"github.com/omarluq/hetzner-blackbsd/internal/hcloud"
)

func newDestroyCmd() *cobra.Command {
	var cmd cobra.Command
	cmd.Use = "destroy"
	cmd.Short = "Destroy BlackBSD build servers"
	cmd.RunE = runDestroy
	return &cmd
}

func runDestroy(cmd *cobra.Command, _ []string) error {
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
		slog.Info("No BlackBSD servers to destroy.")
		return nil
	}

	return destroyServers(cmd.OutOrStdout(), client, cmd, servers)
}

func destroyServers(
	output io.Writer,
	client *hcloud.Client,
	cmd *cobra.Command,
	servers []*hcloudsdk.Server,
) error {
	if _, writeErr := fmt.Fprintf(output, "Destroying %d BlackBSD server(s)...\n", len(servers)); writeErr != nil {
		return writeErr
	}

	for _, server := range servers {
		if writeErr := destroySingleServer(output, client, cmd, server); writeErr != nil {
			return writeErr
		}
	}

	_, writeErr := fmt.Fprintln(output, "\nDone.")
	return writeErr
}

func destroySingleServer(
	output io.Writer,
	client *hcloud.Client,
	cmd *cobra.Command,
	server *hcloudsdk.Server,
) error {
	if _, writeErr := fmt.Fprintf(output, "  %s (%d)... ", server.Name, server.ID); writeErr != nil {
		return writeErr
	}

	deleted, delErr := client.DeleteServer(cmd.Context(), server)
	if delErr != nil {
		_, writeErr := fmt.Fprintf(output, "error: %v\n", delErr)
		return writeErr
	}

	msg := "destroyed"
	if !deleted {
		msg = "not found (skipped)"
	}

	_, writeErr := fmt.Fprintln(output, msg)
	return writeErr
}
