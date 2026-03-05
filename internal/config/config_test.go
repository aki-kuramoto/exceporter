package config

import "testing"

func TestExtractSpreadsheetID(t *testing.T) {
	const wantID = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	cases := []struct {
		name    string
		input   string
		wantID  string
		wantErr bool
	}{
		{
			name:   "docs URL with edit and gid",
			input:  "https://docs.google.com/spreadsheets/d/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/edit?gid=1582528009#gid=1582528009",
			wantID: wantID,
		},
		{
			name:   "docs URL without scheme",
			input:  "docs.google.com/spreadsheets/d/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/edit?gid=1582528009#gid=1582528009",
			wantID: wantID,
		},
		{
			name:   "googleapis files URL with scheme",
			input:  "https://www.googleapis.com/drive/v3/files/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/",
			wantID: wantID,
		},
		{
			name:   "googleapis files URL without scheme",
			input:  "www.googleapis.com/drive/v3/files/XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/",
			wantID: wantID,
		},
		{
			name:   "raw ID",
			input:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			wantID: wantID,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "unrecognized URL",
			input:   "https://example.com/something/unknown",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractSpreadsheetID(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("エラーを期待しましたが got=%q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if got != tc.wantID {
				t.Errorf("got=%q, want=%q", got, tc.wantID)
			}
		})
	}
}

func TestExtractFolderOrDriveID(t *testing.T) {
	const wantID = "111111111111111111111111111111111"

	cases := []struct {
		name    string
		input   string
		wantID  string
		wantErr bool
	}{
		{
			name:   "drive URL with scheme",
			input:  "https://drive.google.com/drive/u/0/folders/111111111111111111111111111111111",
			wantID: wantID,
		},
		{
			name:   "drive URL without scheme",
			input:  "drive.google.com/drive/u/0/folders/111111111111111111111111111111111",
			wantID: wantID,
		},
		{
			name:   "raw ID",
			input:  "111111111111111111111111111111111",
			wantID: wantID,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "unrecognized URL",
			input:   "https://example.com/something/unknown",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractFolderOrDriveID(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("エラーを期待しましたが got=%q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if got != tc.wantID {
				t.Errorf("got=%q, want=%q", got, tc.wantID)
			}
		})
	}
}
