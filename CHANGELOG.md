# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Export Google Spreadsheets to `.xlsx` via the Drive API using Application Default Credentials (ADC)
- YAML configuration supporting both map format (explicit filename) and list format (title auto-fetched from Drive API)
- Folder bulk export: export all spreadsheets in a specified shared drive folder
- Flexible URL parsing for spreadsheet entries:
  - `https://docs.google.com/spreadsheets/d/<ID>/edit?gid=...`
  - `https://www.googleapis.com/drive/v3/files/<ID>/`
  - Bare spreadsheet ID
- Flexible URL parsing for `drive-id` / `folder-id`:
  - `https://drive.google.com/drive/u/0/folders/<ID>`
  - Bare ID
- CLI options: `-o` (output directory), `-v` (verbose), `-h` (help)
- Flags accepted in any position (not just before positional arguments)
- GitHub Actions: CI workflow (`ci.yml`) and release workflow (`release.yml`) for 6 platforms
- Setup guide: `docs/setup-gcloud-auth.md` (English + Japanese)
