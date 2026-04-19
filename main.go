package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AmadlaOrg/weaver-mustache/input"
	"github.com/cbroglie/mustache"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

var (
	infoOutputFlag string
	infoHeryFlag   bool

	infoCmd = &cobra.Command{
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
			if err := writeInfoOutput(os.Stdout, infoOutputFlag, infoHeryFlag, metadata); err != nil {
				fmt.Fprintf(os.Stderr, "Error encoding metadata: %v\n", err)
				os.Exit(1)
			}
		},
	}
)

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
	infoCmd.Flags().StringVarP(&infoOutputFlag, "output", "o", "table", "Output format: table, json, yaml")
	infoCmd.Flags().BoolVar(&infoHeryFlag, "hery", false, "Wrap output in HERY envelope (_type, _body)")

	renderCmd.Flags().StringVarP(&templatePath, "template", "t", "", "Path to the Mustache template file")
	renderCmd.Flags().StringVarP(&filePath, "file", "f", "-", "Input data file path (- for stdin)")
	renderCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: stdout)")
	_ = renderCmd.MarkFlagRequired("template")

	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(renderCmd)
}

type heryEnvelope struct {
	Type string `json:"_type" yaml:"_type"`
	Body any    `json:"_body" yaml:"_body"`
}

func writeInfoOutput(w io.Writer, format string, hery bool, data map[string]any) error {
	if hery {
		return writeHeryOutput(w, format, data)
	}

	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	case "yaml":
		bytes, err := yaml.Marshal(data)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, string(bytes))
		return err
	default:
		table := tablewriter.NewWriter(w)
		table.Header("Field", "Value")
		table.Append("Name", fmt.Sprint(data["name"]))
		table.Append("Version", fmt.Sprint(data["version"]))
		if eng, ok := data["engine"]; ok {
			table.Append("Engine", fmt.Sprint(eng))
		}
		table.Append("Description", fmt.Sprint(data["description"]))
		if supports, ok := data["supports"].([]string); ok {
			table.Append("Supports", strings.Join(supports, "\n"))
		}
		if exts, ok := data["file_extensions"].([]string); ok {
			table.Append("File Extensions", strings.Join(exts, ", "))
		}
		table.Render()
		return nil
	}
}

func writeHeryOutput(w io.Writer, format string, data map[string]any) error {
	envelope := heryEnvelope{
		Type: "amadla.org/entity/tools/info@v1.0.0",
		Body: data,
	}

	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(envelope)
	case "table":
		fmt.Fprintf(w, "_type: %s\n\n", envelope.Type)
		table := tablewriter.NewWriter(w)
		table.Header("Field", "Value")
		table.Append("Name", fmt.Sprint(data["name"]))
		table.Append("Version", fmt.Sprint(data["version"]))
		if eng, ok := data["engine"]; ok {
			table.Append("Engine", fmt.Sprint(eng))
		}
		table.Append("Description", fmt.Sprint(data["description"]))
		if supports, ok := data["supports"].([]string); ok {
			table.Append("Supports", strings.Join(supports, "\n"))
		}
		if exts, ok := data["file_extensions"].([]string); ok {
			table.Append("File Extensions", strings.Join(exts, ", "))
		}
		table.Render()
		return nil
	default:
		bytes, err := yaml.Marshal(envelope)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(w, string(bytes))
		return err
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(2)
	}
}
