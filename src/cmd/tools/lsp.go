package main

import (
	"fmt"
	"github.com/sevenreup/duwa/src/tools/lsp"

	"github.com/spf13/cobra"
	"github.com/tliron/commonlog"
	serverpkg "github.com/tliron/glsp/server"
	"github.com/tliron/kutil/util"
	versionpkg "github.com/tliron/kutil/version"
)

var (
	address  string
	protocol string

	version bool
	verbose int
	logTo   string
)

func init() {
	lspCommand.Flags().StringVarP(&logTo, "log", "l", "", "log to file (defaults to stderr)")
	lspCommand.Flags().CountVarP(&verbose, "verbose", "v", "add a log verbosity level (can be used twice)")
	lspCommand.Flags().BoolVar(&version, "version", false, "print version")

	lspCommand.Flags().StringVarP(&protocol, "protocol", "p", "stdio", "protocol (\"stdio\", \"tcp\", \"websocket\", or \"nodejs\"")
	lspCommand.Flags().StringVarP(&address, "address", "a", ":4389", "address (for \"tcp\" and \"websocket\"")

}

var lspCommand = &cobra.Command{
	Use:   "lsp",
	Short: "Duwa LSP",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if logTo == "" {
			commonlog.Configure(verbose, nil)
		} else {
			commonlog.Configure(verbose, &logTo)
		}

		if verbose > 0 {
			// Reduce Puccini logging even in verbose mode
			commonlog.SetMaxLevel(commonlog.Warning, "duwa")
		}

		if version {
			versionpkg.Print()
			util.Exit(0)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := startSever(); err != nil {
			log.Errorf("failed to start LSP server: %s", err)
			util.Exit(1)
		}
	},
}

func startSever() error {
	log.Infof("version %s", versionpkg.GitVersion)

	server := serverpkg.NewServer(&lsp.Handler, toolName, verbose > 0)

	switch protocol {
	case "stdio":
		return server.RunStdio()

	case "tcp":
		return server.RunTCP(address)

	case "websocket":
		return server.RunWebSocket(address)

	case "nodejs":
		return server.RunNodeJs()

	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}
}
