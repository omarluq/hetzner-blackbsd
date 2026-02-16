package main_test

import (
	"bytes"
	"net"
	"testing"

	blackbsd "github.com/omarluq/hetzner-blackbsd/cmd/hetzner-blackbsd"

	hcloudsdk "github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	t.Parallel()

	cmd := blackbsd.NewVersionCmdForTest()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.RunE(cmd, []string{})

	require.NoError(t, err)
	output := buf.String()
	assert.NotEmpty(t, output, "version output should not be empty")
}

func TestPrintServers(t *testing.T) {
	t.Parallel()

	t.Run("prints server table", func(t *testing.T) {
		t.Parallel()

		var ipv4_1 hcloudsdk.ServerPublicNetIPv4
		ipv4_1.IP = net.ParseIP("1.2.3.4")

		var publicNet1 hcloudsdk.ServerPublicNet
		publicNet1.IPv4 = ipv4_1

		var srv1 hcloudsdk.Server
		srv1.ID = 1
		srv1.Name = "blackbsd-builder-123"
		srv1.Status = hcloudsdk.ServerStatusRunning
		srv1.RescueEnabled = false
		srv1.PublicNet = publicNet1

		var ipv4_2 hcloudsdk.ServerPublicNetIPv4
		ipv4_2.IP = net.ParseIP("5.6.7.8")

		var publicNet2 hcloudsdk.ServerPublicNet
		publicNet2.IPv4 = ipv4_2

		var srv2 hcloudsdk.Server
		srv2.ID = 2
		srv2.Name = "blackbsd-builder-456"
		srv2.Status = hcloudsdk.ServerStatusRunning
		srv2.RescueEnabled = true
		srv2.PublicNet = publicNet2

		servers := []*hcloudsdk.Server{&srv1, &srv2}

		buf := new(bytes.Buffer)
		err := blackbsd.PrintServersForTest(buf, servers)

		require.NoError(t, err)
		output := buf.String()

		assert.Contains(t, output, "blackbsd-builder-123")
		assert.Contains(t, output, "blackbsd-builder-456")
		assert.Contains(t, output, "1.2.3.4")
		assert.Contains(t, output, "5.6.7.8")
		assert.Contains(t, output, "yes")
		assert.Contains(t, output, "no")
		assert.Contains(t, output, "Found 2 BlackBSD server(s)")
	})

	t.Run("prints with nil IPv4", func(t *testing.T) {
		t.Parallel()

		var srv hcloudsdk.Server
		srv.ID = 1
		srv.Name = "test-server"
		srv.Status = hcloudsdk.ServerStatusInitializing
		srv.RescueEnabled = false

		servers := []*hcloudsdk.Server{&srv}

		buf := new(bytes.Buffer)
		err := blackbsd.PrintServersForTest(buf, servers)

		require.NoError(t, err)
		output := buf.String()

		assert.Contains(t, output, "test-server")
		assert.Contains(t, output, "initializing")
		assert.Contains(t, output, "Found 1 BlackBSD server(s)")
	})
}

func TestRootCommandSetup(t *testing.T) {
	t.Parallel()

	cmd := blackbsd.NewRootCmdForTest()

	assert.Equal(t, "hetzner-blackbsd", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestDestroyCommandSetup(t *testing.T) {
	t.Parallel()

	cmd := blackbsd.NewDestroyCmdForTest()

	assert.Equal(t, "destroy", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotNil(t, cmd.RunE)
}

func TestStatusCommandSetup(t *testing.T) {
	t.Parallel()

	cmd := blackbsd.NewStatusCmdForTest()

	assert.Equal(t, "status", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotNil(t, cmd.RunE)
}
