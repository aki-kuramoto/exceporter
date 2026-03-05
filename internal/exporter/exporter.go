package exporter

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const xlsxMimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

// Exporter downloads Google Spreadsheets via the Drive API.
type Exporter struct {
	srv     *drive.Service
	outDir  string
	verbose bool
}

// New initializes an Exporter using Application Default Credentials.
func New(ctx context.Context, outDir string, verbose bool) (*Exporter, error) {
	srv, err := drive.NewService(ctx, option.WithScopes(drive.DriveReadonlyScope))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Google Drive service: %w", err)
	}
	return &Exporter{srv: srv, outDir: outDir, verbose: verbose}, nil
}

// ExportByID exports a spreadsheet to outDir/<name>.xlsx.
// Progress and completion messages are printed only when verbose is true.
func (e *Exporter) ExportByID(name, spreadsheetID string) error {
	if e.verbose {
		log.Printf("[%s] downloading... (ID: %s)", name, spreadsheetID)
	}

	resp, err := e.srv.Files.Export(spreadsheetID, xlsxMimeType).Download()
	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}
	defer resp.Body.Close()

	destPath := filepath.Join(e.outDir, name+".xlsx")
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file (%s): %w", destPath, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file (%s): %w", destPath, err)
	}

	if e.verbose {
		log.Printf("[%s] done -> %s", name, destPath)
	}
	return nil
}

// ExportFolder exports all spreadsheets found in the specified shared drive folder.
// Progress messages are printed only when verbose is true; errors are always printed.
func (e *Exporter) ExportFolder(ctx context.Context, driveID, folderID string) error {
	if e.verbose {
		log.Printf("[folder] searching spreadsheets... (driveID: %s, folderID: %s)", driveID, folderID)
	}

	query := fmt.Sprintf(
		"'%s' in parents and mimeType = 'application/vnd.google-apps.spreadsheet' and trashed = false",
		folderID,
	)

	var errs []error
	err := e.srv.Files.List().
		Q(query).
		Corpora("drive").
		DriveId(driveID).
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Fields("files(id, name)").
		Pages(ctx, func(fl *drive.FileList) error {
			for _, f := range fl.Files {
				if err := e.ExportByID(f.Name, f.Id); err != nil {
					log.Printf("[%s] error: %v", f.Name, err)
					errs = append(errs, fmt.Errorf("%s: %w", f.Name, err))
				}
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("folder listing failed: %w", err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d export(s) failed", len(errs))
	}
	return nil
}

// GetFileName fetches the title of a spreadsheet from the Drive API.
// Used to determine the output filename for list-format sheets entries.
func (e *Exporter) GetFileName(spreadsheetID string) (string, error) {
	f, err := e.srv.Files.Get(spreadsheetID).
		SupportsAllDrives(true).
		Fields("name").
		Do()
	if err != nil {
		return "", fmt.Errorf("failed to fetch filename: %w", err)
	}
	return f.Name, nil
}
