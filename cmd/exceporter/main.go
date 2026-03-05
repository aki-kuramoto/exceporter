package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aki-kuramoto/exceporter/internal/config"
	"github.com/aki-kuramoto/exceporter/internal/exporter"
)

const defaultOutDir = "./exceporter-out"

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: exceporter [options] <config.yaml>

Options:
  -o <dir>  Output directory (default: %s)
  -v        Verbose output
  -h        Show this help
`, defaultOutDir)
}

// parseArgs parses command-line arguments allowing flags in any position.
// The standard flag package stops parsing at the first non-flag argument,
// so we implement our own parser here.
func parseArgs(args []string) (yamlPath, outDir string, verbose bool, err error) {
	outDir = defaultOutDir

	var positionals []string
	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-v" || arg == "--v":
			verbose = true

		case arg == "-o" || arg == "--o":
			if i+1 >= len(args) {
				err = fmt.Errorf("-o requires an argument")
				return
			}
			i++
			outDir = args[i]

		case strings.HasPrefix(arg, "-o="):
			outDir = arg[len("-o="):]

		case arg == "-h" || arg == "--help":
			printUsage()
			os.Exit(0)

		case strings.HasPrefix(arg, "-"):
			err = fmt.Errorf("unknown option: %s", arg)
			return

		default:
			positionals = append(positionals, arg)
		}
	}

	if len(positionals) == 0 {
		err = fmt.Errorf("config YAML file is required")
		return
	}
	yamlPath = positionals[0]
	return
}

func main() {
	yamlPath, outDir, verbose, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.Load(yamlPath)
	if err != nil {
		log.Fatalf("error: failed to load config: %v", err)
	}

	ctx := context.Background()
	if err := exporter.Run(ctx, cfg, outDir, verbose); err != nil {
		log.Fatalf("error: %v", err)
	}

	if verbose {
		log.Println("all exports completed successfully.")
	}
}
