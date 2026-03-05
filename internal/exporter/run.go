package exporter

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aki-kuramoto/exceporter/internal/config"
)

// Run exports spreadsheets according to cfg.
// Processing continues even if individual exports fail; the total failure count
// is returned as an error at the end.
func Run(ctx context.Context, cfg *config.Config, outDir string, verbose bool) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory (%s): %w", outDir, err)
	}

	exp, err := New(ctx, outDir, verbose)
	if err != nil {
		return err
	}

	var failCount int

	// Map-format sheets: user-specified filename.
	for name, urlOrID := range cfg.Sheets.Named {
		id, err := config.ExtractSpreadsheetID(urlOrID)
		if err != nil {
			log.Printf("[%s] skipped: %v", name, err)
			failCount++
			continue
		}
		if err := exp.ExportByID(name, id); err != nil {
			log.Printf("[%s] error: %v", name, err)
			failCount++
		}
	}

	// List-format sheets: filename derived from spreadsheet title via Drive API.
	for _, urlOrID := range cfg.Sheets.Listed {
		id, err := config.ExtractSpreadsheetID(urlOrID)
		if err != nil {
			log.Printf("[skip] cannot extract ID (%s): %v", urlOrID, err)
			failCount++
			continue
		}
		name, err := exp.GetFileName(id)
		if err != nil {
			// Fall back to the ID as filename and log a warning.
			log.Printf("[%s] warning: could not fetch title, using ID as filename: %v", id, err)
			name = id
		}
		if err := exp.ExportByID(name, id); err != nil {
			log.Printf("[%s] error: %v", name, err)
			failCount++
		}
	}

	// Folder bulk export.
	if cfg.Folder != nil {
		if err := exp.ExportFolder(ctx, cfg.Folder.DriveID, cfg.Folder.FolderID); err != nil {
			log.Printf("[folder] error: %v", err)
			failCount++
		}
	}

	if failCount > 0 {
		return fmt.Errorf("%d item(s) failed to export", failCount)
	}
	return nil
}
