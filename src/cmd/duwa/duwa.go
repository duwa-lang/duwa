package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/duwa-lang/duwa/src/duwa"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/runtime"
	"github.com/duwa-lang/duwa/src/runtime/native"

	"github.com/spf13/cobra"
)

var (
	// Version information
	version = " v0.1.1"
	commit  = "unknown"
	date    = "unknown"

	// Flags
	verbose bool
	debug   bool
	trace   bool
)

var rootCmd = &cobra.Command{
	Use:   "duwa",
	Short: "Duwa CLI tool",
	Long:  `A CLI tool for running duwa programs`,
}

var runCmd = &cobra.Command{
	Use:   "run [filepath]",
	Short: "Run a duwa file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]

		if filepath == "" {
			log.Fatal("Please provide a file to run")
		}

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		console := native.NewConsole()
		env := object.New(logger, console)

		// Enable debugging if requested
		if debug || trace {
			if debug {
				debugObserver := runtime.NewDebuggerObserver(logger, verbose)
				env.ObserverManager.Register(debugObserver)
			}
			if trace {
				traceObserver := runtime.NewTraceObserver(logger)
				env.ObserverManager.Register(traceObserver)
			}
		}

		duwa := duwa.New(env)
		value := duwa.RunFile(filepath)
		if value != nil {
			if object.IsError(value) {
				log.Fatal(value.String())
			}
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Duwa CLI Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", date)
	},
}

func init() {
	// Add flags to run command
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	runCmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debug mode")
	runCmd.Flags().BoolVarP(&trace, "trace", "t", false, "enable function call tracing")

	// Add subcommands
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
