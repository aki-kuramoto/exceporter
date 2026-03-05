package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Regular expressions for extracting IDs from various URL formats.
//
// Spreadsheet ID:
//   - docs.google.com/spreadsheets/d/<ID>/...
//   - googleapis.com/drive/v3/files/<ID>/...
//   - <ID> (bare ID)
//
// Drive / Folder ID:
//   - drive.google.com/drive/.../folders/<ID>
//   - <ID> (bare ID)

var (
	// /d/<ID> — Google Docs URL pattern
	reSpreadsheetPath = regexp.MustCompile(`/d/([a-zA-Z0-9-_]+)`)
	// /files/<ID> — Drive REST API URL pattern
	reDriveFilePath = regexp.MustCompile(`/files/([a-zA-Z0-9-_]+)`)
	// /folders/<ID> — Google Drive folder URL pattern
	reFolderPath = regexp.MustCompile(`/folders/([a-zA-Z0-9-_]+)`)
)

// ExtractSpreadsheetID extracts a spreadsheet ID from various URL formats or a bare ID.
//
// Supported formats:
//   - https://docs.google.com/spreadsheets/d/<ID>/edit?gid=...#gid=...
//   - docs.google.com/spreadsheets/d/<ID>/edit?gid=...
//   - https://www.googleapis.com/drive/v3/files/<ID>/
//   - www.googleapis.com/drive/v3/files/<ID>/
//   - <ID> (bare spreadsheet ID)
func ExtractSpreadsheetID(urlOrID string) (string, error) {
	// /d/<ID> pattern (docs.google.com)
	if m := reSpreadsheetPath.FindStringSubmatch(urlOrID); len(m) >= 2 {
		return m[1], nil
	}
	// /files/<ID> pattern (googleapis.com)
	if m := reDriveFilePath.FindStringSubmatch(urlOrID); len(m) >= 2 {
		return m[1], nil
	}
	// Treat as a bare ID if no slashes are present
	if urlOrID != "" && !strings.Contains(urlOrID, "/") {
		return urlOrID, nil
	}
	return "", fmt.Errorf("cannot extract spreadsheet ID from: %s", urlOrID)
}

// ExtractFolderOrDriveID extracts a drive/folder ID from a URL or a bare ID.
//
// Supported formats:
//   - https://drive.google.com/drive/u/0/folders/<ID>
//   - drive.google.com/drive/u/0/folders/<ID>
//   - <ID> (bare ID)
func ExtractFolderOrDriveID(urlOrID string) (string, error) {
	// /folders/<ID> pattern (drive.google.com)
	if m := reFolderPath.FindStringSubmatch(urlOrID); len(m) >= 2 {
		return m[1], nil
	}
	// Treat as a bare ID if no slashes are present
	if urlOrID != "" && !strings.Contains(urlOrID, "/") {
		return urlOrID, nil
	}
	return "", fmt.Errorf("cannot extract drive/folder ID from: %s", urlOrID)
}

// SheetsConfig represents the sheets: section and supports both map and list formats.
//
//	Map format (explicit filename):
//	  sheets:
//	    chr: https://docs.google.com/spreadsheets/d/XXX/
//
//	List format (filename derived from spreadsheet title):
//	  sheets:
//	    - https://docs.google.com/spreadsheets/d/XXX/
type SheetsConfig struct {
	Named  map[string]string // name → URL (map format)
	Listed []string          // URL list (list format)
}

// UnmarshalYAML implements a custom decoder that accepts both map and list formats.
func (s *SheetsConfig) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.MappingNode:
		return value.Decode(&s.Named)
	case yaml.SequenceNode:
		return value.Decode(&s.Listed)
	default:
		return fmt.Errorf("sheets must be a map or a list")
	}
}

// IsEmpty reports whether no sheets are configured.
func (s *SheetsConfig) IsEmpty() bool {
	return len(s.Named) == 0 && len(s.Listed) == 0
}

// Config is the top-level structure of the YAML configuration file.
type Config struct {
	Sheets SheetsConfig  `yaml:"sheets"`
	Folder *FolderConfig `yaml:"folder"`
}

// FolderConfig holds the configuration for bulk-exporting all spreadsheets in a folder.
// drive-id and folder-id accept either a bare ID or a Google Drive URL.
type FolderConfig struct {
	DriveID  string `yaml:"drive-id"`
	FolderID string `yaml:"folder-id"`
}

// normalize resolves URL-form drive-id / folder-id values to bare IDs.
// Called automatically after loading.
func (f *FolderConfig) normalize() error {
	driveID, err := ExtractFolderOrDriveID(f.DriveID)
	if err != nil {
		return fmt.Errorf("folder.drive-id: %w", err)
	}
	f.DriveID = driveID

	folderID, err := ExtractFolderOrDriveID(f.FolderID)
	if err != nil {
		return fmt.Errorf("folder.folder-id: %w", err)
	}
	f.FolderID = folderID
	return nil
}

// Load reads the YAML file at path and returns a validated Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validate checks the config and normalizes URL-form IDs.
func (c *Config) validate() error {
	if c.Sheets.IsEmpty() && c.Folder == nil {
		return fmt.Errorf("at least one of 'sheets' or 'folder' must be specified")
	}
	if c.Folder != nil {
		if c.Folder.DriveID == "" {
			return fmt.Errorf("folder.drive-id is empty")
		}
		if c.Folder.FolderID == "" {
			return fmt.Errorf("folder.folder-id is empty")
		}
		if err := c.Folder.normalize(); err != nil {
			return err
		}
	}
	return nil
}
