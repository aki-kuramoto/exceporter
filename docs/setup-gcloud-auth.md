# Setting Up Google Drive Authentication with ADC

This guide explains how to authenticate with Google Drive API using
**Application Default Credentials (ADC)** via the Google Cloud CLI (`gcloud`),
without registering a dedicated OAuth client application.

---

## Overview

Normally, accessing Google APIs requires registering an app and obtaining a
Client ID. As an alternative for local development, you can borrow the
credentials of the official `gcloud` CLI.

Google Cloud manages API usage per **project**. A default project is usually
available even without explicitly creating one. By enabling the Drive API on
any Google Cloud project and logging in locally via `gcloud`, you can export
spreadsheets (as `.xlsx`) without writing any authentication code yourself.

---

## One-Time Setup

### 0. Install Google Cloud CLI

Download and install from the official page:  
<https://cloud.google.com/sdk/docs/install>

During installation you may be prompted to run `gcloud init`, which will also
walk you through logging in and selecting a project.

### 1. Log in to Google Cloud

```sh
gcloud auth login
```

### 2. Set the target project

Replace `<YOUR_PROJECT_ID>` with your actual Google Cloud project ID.  
If you don't have a specific project, use the default one shown during `gcloud init`.

```sh
gcloud config set project <YOUR_PROJECT_ID>
```

You can check available projects with:

```sh
gcloud projects list
```

### 3. Enable the Drive API for the project

```sh
gcloud services enable drive.googleapis.com
```

---

## Periodic Task — Refreshing Application Default Credentials (ADC)

ADC tokens expire. Run the following command to refresh them:

```sh
gcloud auth application-default login \
  --scopes="https://www.googleapis.com/auth/drive.readonly,https://www.googleapis.com/auth/cloud-platform"
```

This saves a credential file locally (e.g. on Windows:
`%APPDATA%\gcloud\application_default_credentials.json`).  
Client libraries such as `google.golang.org/api` discover this file
automatically — no explicit path configuration is needed.

**Why specify scopes?**

- `drive.readonly` — grants read access to Google Drive files.
- `cloud-platform` — required as the base scope when issuing ADC.

---

## What You Can Do After Setup

With valid ADC in place, any library that supports ADC
(e.g. `google.golang.org/api/drive/v3`) can make requests like:

```
GET https://www.googleapis.com/drive/v3/files/{fileId}/export
    ?mimeType=application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
```

---

## Extending to Server Environments (Service Account)

The same Go code works on servers — just switch the credential source:

1. Create a **service account** in the Google Cloud Console and download its JSON key.
2. Set the environment variable:
   ```sh
   export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
   ```
3. No code changes are needed. The library checks `GOOGLE_APPLICATION_CREDENTIALS`
   before falling back to local ADC.

> **Note:** On Google Cloud infrastructure (GCE, Cloud Run, etc.), the instance's
> attached service account is used automatically without any key file.

> **Important:** When using a service account, make sure its email address has
> **Viewer** (or higher) access to the target spreadsheets or shared drives.

---
---

# Google Drive 認証のセットアップ (ADC 利用)

このガイドでは、OAuth クライアントを独自登録せずに、Google Cloud CLI (`gcloud`) の
**Application Default Credentials (ADC)** を借りて Google Drive API を利用する方法を説明します。

---

## 概要

通常、Google API へのアクセスにはアプリ登録 (独自の Client ID) が必要ですが、
ローカル開発では `gcloud` の認証情報を代用できます。

Google Cloud は **プロジェクト** 単位で API 使用量を管理しています。
明示的に作成しなくてもデフォルトのプロジェクトが存在するはずです。
そのプロジェクトで Drive API を有効化し、`gcloud` でローカルログインすることで、
認証コードなしにスプレッドシートを `.xlsx` 形式でエクスポートできるようになります。

---

## 一回だけ行う設定

### 0. Google Cloud CLI をインストールする

公式ページからダウンロードしてインストールしてください:  
<https://cloud.google.com/sdk/docs/install?hl=ja>

インストール時に `gcloud init` が実行されると、ログインとプロジェクト選択を促されます。

### 1. Google Cloud にログインする

```sh
gcloud auth login
```

### 2. 操作対象のプロジェクトを設定する

`<YOUR_PROJECT_ID>` を実際のプロジェクト ID に置き換えてください。  
特定のプロジェクトがない場合は `gcloud init` 時に表示されたデフォルトプロジェクトを使用します。

```sh
gcloud config set project <YOUR_PROJECT_ID>
```

利用可能なプロジェクトは以下で確認できます:

```sh
gcloud projects list
```

### 3. そのプロジェクトで Drive API を有効化する

```sh
gcloud services enable drive.googleapis.com
```

---

## たまに行う作業 — ADC の更新

ADC トークンには有効期限があります。期限が切れたら以下を実行してください:

```sh
gcloud auth application-default login \
  --scopes="https://www.googleapis.com/auth/drive.readonly,https://www.googleapis.com/auth/cloud-platform"
```

認証情報はローカルに保存されます (例: Windows では `%APPDATA%\gcloud\application_default_credentials.json`) 。  
`google.golang.org/api` などのクライアントライブラリはこのファイルを自動で見つけます。
パスを明示的に指定する必要はありません。

**スコープを指定する理由**

- `drive.readonly` — Google Drive ファイルへの読み取り権限
- `cloud-platform` — ADC 発行時にベースとして必要なスコープ

---

## ADC 設定後にできること

有効な ADC があれば、ADC をサポートするライブラリ (`google.golang.org/api/drive/v3` など) から
以下のようなリクエストを実行できます:

```
GET https://www.googleapis.com/drive/v3/files/{fileId}/export
    ?mimeType=application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
```

---

## サーバー環境への応用 (サービスアカウント) 

Go のコードはそのままに、資格情報のソースを切り替えるだけでサーバーでも動作します:

1. Google Cloud Console で **サービスアカウント** を作成し、JSON キーをダウンロードする
2. 環境変数を設定する:
   ```sh
   export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
   ```
3. コードの変更は不要です。ライブラリはローカル ADC よりも `GOOGLE_APPLICATION_CREDENTIALS` を優先して読み込みます。

> **補足:** GCE や Cloud Run などの Google Cloud 上の環境では、インスタンスに付与された
> サービスアカウントの権限が自動的に使用されるため、キーファイルは不要です。

> **注意:** サービスアカウントを使う場合は、そのアカウントのメールアドレスが対象の
> スプレッドシートや共有ドライブに対して **閲覧権限以上** を持っている必要があります。
