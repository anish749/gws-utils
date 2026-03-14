package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/anish749/gws_utils/internal/converter"
	"github.com/anish749/gws_utils/internal/exporter"
	"github.com/anish749/gws_utils/internal/gws"
	"github.com/spf13/cobra"
)

var (
	outputDir string
	account   string
)

func init() {
	downloadCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "output directory")
	downloadCmd.Flags().StringVarP(&account, "account", "a", "", "gws account email to use")
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download [document-url-or-id]",
	Short: "Download a Google Doc as markdown, preserving tabs",
	Args:  cobra.ExactArgs(1),
	RunE:  runDownload,
}

func runDownload(cmd *cobra.Command, args []string) error {
	docID := extractDocID(args[0])
	if docID == "" {
		return fmt.Errorf("could not extract document ID from %q", args[0])
	}

	client := gws.NewClient(account)
	conv := converter.NewMarkdownConverter()
	fs := exporter.NewFilesystemExporter()

	fmt.Printf("Fetching document %s...\n", docID)
	doc, err := client.GetDocument(docID)
	if err != nil {
		return fmt.Errorf("fetching document: %w", err)
	}

	tabs := doc.AllTabs()
	fmt.Printf("Found %d tab(s) in %q\n", len(tabs), doc.Title)

	if len(tabs) == 1 {
		md := conv.Convert(tabs[0])
		filename := sanitizeFilename(doc.Title) + ".md"
		outPath := filepath.Join(outputDir, filename)
		if err := fs.WriteFile(outPath, md); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Printf("Saved: %s\n", outPath)
	} else {
		dirName := sanitizeFilename(doc.Title)
		dirPath := filepath.Join(outputDir, dirName)
		if err := fs.EnsureDir(dirPath); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
		for i, tab := range tabs {
			md := conv.Convert(tab)
			filename := fmt.Sprintf("%02d_%s.md", i, sanitizeFilename(tab.Title))
			outPath := filepath.Join(dirPath, filename)
			if err := fs.WriteFile(outPath, md); err != nil {
				return fmt.Errorf("writing tab %q: %w", tab.Title, err)
			}
			fmt.Printf("Saved: %s\n", outPath)
		}
	}

	return nil
}

// extractDocID pulls the document ID from a Google Docs URL or returns the input as-is.
func extractDocID(input string) string {
	// Handle full URLs like https://docs.google.com/document/d/DOC_ID/edit...
	if strings.Contains(input, "docs.google.com/document/d/") {
		parts := strings.Split(input, "/d/")
		if len(parts) < 2 {
			return ""
		}
		id := parts[1]
		if idx := strings.Index(id, "/"); idx != -1 {
			id = id[:idx]
		}
		return id
	}
	return input
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	name = replacer.Replace(name)
	name = strings.TrimSpace(name)
	if len(name) > 200 {
		name = name[:200]
	}
	return name
}
