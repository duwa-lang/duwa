package main

import (
	"context"
	"log"

	"github.com/sevenreup/duwa/src/tools/lsp"
	"github.com/sevenreup/duwa/src/tools/lsp/transport"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/tliron/kutil/util"
	versionpkg "github.com/tliron/kutil/version"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

var (
	address     string
	protocolSVR string

	version  bool
	verbose  int
	logFile  string
	logLevel int
)

func init() {
	lspCommand.Flags().StringVarP(&logFile, "logFile", "f", "", "log to file (defaults to stderr)")
	lspCommand.Flags().IntVarP(&logLevel, "logLevel", "l", int(zap.DebugLevel), "log level")
	lspCommand.Flags().BoolVar(&version, "version", false, "print version")

	lspCommand.Flags().StringVarP(&protocolSVR, "protocol", "p", "stdio", "protocol (\"stdio\", \"tcp\", \"websocket\", or \"nodejs\"")
	lspCommand.Flags().StringVarP(&address, "address", "a", ":4389", "address (for \"tcp\" and \"websocket\"")

}

var lspCommand = &cobra.Command{
	Use:   "lsp",
	Short: "Duwa LSP",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if version {
			versionpkg.Print()
			util.Exit(0)
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		startSever()
	},
}

func startSever() {
	lsp.InitLog(logLevel, logFile)
	server := lsp.NewServer("duwa-lsp", verbose > 0)
	server.Logger.Sugar().Info("version %s", versionpkg.GitVersion)

	var conn jsonrpc2.Conn

	switch protocolSVR {
	case "stdio":
		server.Logger.Sugar().Info("Starting with stdio")
		conn = transport.CreateTransport()
	case "tcp":
	case "websocket":
	case "nodejs":
	default:
		server.Logger.Sugar().Error("unsupported protocol: %s", protocolSVR)
		util.Exit(1)
	}
	ctx := protocol.WithLogger(context.Background(), server.Logger)
	conn.Go(ctx, protocol.ServerHandler(server, nil))
	<-ctx.Done()
	log.Println("exit")
}
