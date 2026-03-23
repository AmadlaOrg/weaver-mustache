package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/AmadlaOrg/weaver-mustache/input"
	"github.com/cbroglie/mustache"
	"github.com/spf13/cobra"
)

const (
	appName = "weaver-mustache"
	version = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:     appName,
	Short:   "Mustache template engine plugin for Amadla Weaver",
	Version: version,
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show plugin metadata",
	Run: func(cmd *cobra.Command, args []string) {
		metadata := map[string]any{
			"name":            appName,
			"version":         version,
			"engine":          "mustache",
			"supports":        []string{"amadla.org/entity/template@^v1.0.0"},
			"file_extensions": []string{".mustache"},
			"description":     "Mustache template engine plugin for Weaver",
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(metadata); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding metadata: %v\n", err)
			os.Exit(1)
		}
	},
}

var (
	templatePath string
	filePath     string
	outputPath   string
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a Mustache template",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load input data
		data, err := input.Load(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading input: %v\n", err)
			os.Exit(1)
		}

		// Render template
		result, err := mustache.RenderFile(templatePath, data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering template: %v\n", err)
			os.Exit(1)
		}

		// Write output
		var out io.Writer = os.Stdout
		if outputPath != "" {
			f, err := os.Create(outputPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()
			out = f
		}

		if _, err := fmt.Fprint(out, result); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	renderCmd.Flags().StringVarP(&templatePath, "template", "t", "", "Path to the Mustache template file")
	renderCmd.Flags().StringVarP(&filePath, "file", "f", "-", "Input data file path (- for stdin)")
	renderCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: stdout)")
	_ = renderCmd.MarkFlagRequired("template")

	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(renderCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
