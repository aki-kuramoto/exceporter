# exceporter

A CLI tool to export Google Spreadsheets as `.xlsx` files to your local machine.

[![Latest Release](https://img.shields.io/github/v/release/aki-kuramoto/exceporter?style=for-the-badge&label=Download%20Latest%20Release&color=2ea44f)](https://github.com/aki-kuramoto/exceporter/releases/latest)

**[日本語版は下部をご覧ください / Japanese version below](#exceporter-ja)**

---

## Prerequisites

- Go 1.21+
- [Google Cloud CLI](https://cloud.google.com/sdk/docs/install) installed
- Authenticated via Application Default Credentials (ADC):

```sh
gcloud auth application-default login \
  --scopes="https://www.googleapis.com/auth/drive.readonly,https://www.googleapis.com/auth/cloud-platform"
```

> See [docs/setup-gcloud-auth.md](docs/setup-gcloud-auth.md) for a step-by-step
> guide including how to enable the Drive API on your Google Cloud project.

## Installation

```sh
go install github.com/aki-kuramoto/exceporter/cmd/exceporter@latest
```

Or build from source:

```sh
git clone https://github.com/aki-kuramoto/exceporter.git
cd exceporter
go build -o exceporter ./cmd/exceporter/...
```

## Usage

```sh
exceporter [options] <config.yaml>

Options:
  -o <dir>  Output directory (default: ./exceporter-out)
  -v        Verbose output
  -h        Show help
```

### Configuration YAML

```yaml
# Map format: explicit filename  →  saves as <name>.xlsx
sheets:
  report-2024: https://docs.google.com/spreadsheets/d/SPREADSHEET_ID/edit
  budget:      https://docs.google.com/spreadsheets/d/ANOTHER_ID/edit

# List format: no filename  →  saves using the spreadsheet's title
sheets:
  - https://docs.google.com/spreadsheets/d/SPREADSHEET_ID/edit
  - https://docs.google.com/spreadsheets/d/ANOTHER_ID/edit

folder:
  drive-id:  SHARED_DRIVE_ID   # ID or Google Drive URL
  folder-id: FOLDER_ID         # ID or Google Drive URL
```

- `sheets` accepts either map or list format.
- **Map format**: the key is used as the output filename.
- **List format**: the spreadsheet title is fetched via the Drive API and used as the filename.
- `folder` exports all spreadsheets found in the specified shared drive folder.
- `sheets` and `folder` can be used together.

**Supported URL formats for spreadsheet entries:**

| Format | Example |
|---|---|
| Full Sheets URL | `https://docs.google.com/spreadsheets/d/<ID>/edit?gid=...` |
| Sheets URL (no scheme) | `docs.google.com/spreadsheets/d/<ID>/edit` |
| Drive REST API URL | `https://www.googleapis.com/drive/v3/files/<ID>/` |
| Bare ID | `<SPREADSHEET_ID>` |

**Supported URL formats for `drive-id` / `folder-id`:**

| Format | Example |
|---|---|
| Full Drive URL | `https://drive.google.com/drive/u/0/folders/<ID>` |
| Drive URL (no scheme) | `drive.google.com/drive/u/0/folders/<ID>` |
| Bare ID | `<ID>` |

### Examples

```sh
# Export named sheets
exceporter target.yaml

# Custom output directory
exceporter -o ./downloads target.yaml

# Verbose output
exceporter -v target.yaml
```

Output files are saved as `<output-dir>/<name>.xlsx`.

## Authentication

This tool uses **Application Default Credentials (ADC)** — no OAuth client
registration is required. See
[docs/setup-gcloud-auth.md](docs/setup-gcloud-auth.md) for full setup
instructions, including how to use a service account for server environments.

## License

MIT

---

<a name="exceporter-ja"></a>
# exceporter

Google Spreadsheets をローカルに `.xlsx` 形式でエクスポートする CLI ツールです。

[![最新リリース](https://img.shields.io/github/v/release/aki-kuramoto/exceporter?style=for-the-badge&label=%E6%9C%80%E6%96%B0%E3%83%AA%E3%83%AA%E3%83%BC%E3%82%B9%E3%82%92%E3%83%80%E3%82%A6%E3%83%B3%E3%83%AD%E3%83%BC%E3%83%89&color=2ea44f)](https://github.com/aki-kuramoto/exceporter/releases/latest)

## 前提条件

- Go 1.21+
- [Google Cloud CLI](https://cloud.google.com/sdk/docs/install?hl=ja) がインストール済みであること
- Application Default Credentials (ADC) で認証済みであること:

```sh
gcloud auth application-default login \
  --scopes="https://www.googleapis.com/auth/drive.readonly,https://www.googleapis.com/auth/cloud-platform"
```

> セットアップの詳細（Drive API の有効化方法なども含む）は
> [docs/setup-gcloud-auth.md](docs/setup-gcloud-auth.md) を参照してください。

## インストール

```sh
go install github.com/aki-kuramoto/exceporter/cmd/exceporter@latest
```

またはソースからビルド:

```sh
git clone https://github.com/aki-kuramoto/exceporter.git
cd exceporter
go build -o exceporter ./cmd/exceporter/...
```

## 使い方

```sh
exceporter [オプション] <設定YAMLファイル>

オプション:
  -o <ディレクトリ>  出力ディレクトリ (デフォルト: ./exceporter-out)
  -v               詳細ログを出力する
  -h               ヘルプを表示する
```

### 設定 YAML フォーマット

```yaml
# マップ形式: 名前を指定する → <名前>.xlsx で保存
sheets:
  report-2024: https://docs.google.com/spreadsheets/d/SPREADSHEET_ID/edit
  budget:      https://docs.google.com/spreadsheets/d/ANOTHER_ID/edit

# リスト形式: 名前を指定しない → スプレッドシートのタイトルで保存
sheets:
  - https://docs.google.com/spreadsheets/d/SPREADSHEET_ID/edit
  - https://docs.google.com/spreadsheets/d/ANOTHER_ID/edit

folder:
  drive-id:  SHARED_DRIVE_ID   # ID または Google Drive の URL
  folder-id: FOLDER_ID         # ID または Google Drive の URL
```

- `sheets` はマップ形式・リスト形式どちらでも動作します。
- **マップ形式**: キーが保存先ファイル名になります。
- **リスト形式**: Drive API からスプレッドシートのタイトルを取得してファイル名に使用します。
- `folder` を指定した場合、そのフォルダ内の全スプレッドシートを一括エクスポートします。
- `sheets` と `folder` は同時に指定できます。

**スプレッドシート指定の対応フォーマット:**

| 形式 | 例 |
|---|---|
| Sheets URL (https あり) | `https://docs.google.com/spreadsheets/d/<ID>/edit?gid=...` |
| Sheets URL (https なし) | `docs.google.com/spreadsheets/d/<ID>/edit` |
| Drive REST API URL | `https://www.googleapis.com/drive/v3/files/<ID>/` |
| 生 ID | `<SPREADSHEET_ID>` |

**`drive-id` / `folder-id` の対応フォーマット:**

| 形式 | 例 |
|---|---|
| Drive URL (https あり) | `https://drive.google.com/drive/u/0/folders/<ID>` |
| Drive URL (https なし) | `drive.google.com/drive/u/0/folders/<ID>` |
| 生 ID | `<ID>` |

### 実行例

```sh
# 名前指定ありのシートをエクスポート
exceporter target.yaml

# 出力先を指定
exceporter -o ./downloads target.yaml

# 詳細ログあり
exceporter -v target.yaml
```

出力ファイルは `<出力ディレクトリ>/<名前>.xlsx` として保存されます。

## 認証について

このツールは **Application Default Credentials (ADC)** を使用するため、OAuth クライアントの登録は不要です。
サービスアカウントを使ったサーバー環境での利用方法も含め、詳細は
[docs/setup-gcloud-auth.md](docs/setup-gcloud-auth.md) を参照してください。

## ライセンス

MIT
