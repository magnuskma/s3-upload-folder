# S3 Folder Uploader

This project provides a binary for uploading files from a local folder to an AWS S3 bucket (or compatible service) with support for custom endpoints, content type detection, folder structure preservation, and parallel uploads.

## Features

- Stream files from a local folder for efficient memory usage.
- Detect and set the correct `Content-Type` for each file based on its extension.
- Preserve the folder structure during upload to S3.
- Support for custom S3-compatible endpoints (e.g., Fly.io).
- Parallel file uploads with configurable worker count.

## Requirements

- AWS S3 credentials or compatible storage service credentials
- A Unix-like environment for running the Bash script (Linux, macOS, WSL, etc.)

## Usage

### Using the Bash Script

1. Make the Bash script executable:

   ```bash
   curl -L -o upload_to_s3.sh https://raw.githubusercontent.com/magnuskma/s3-upload-folder/main/upload_to_s3.sh && chmod +x upload_to_s3.sh && curl -L -o .env.example https://raw.githubusercontent.com/magnuskma/s3-upload-folder/main/.env.example
   ```

2. Run the script:

   ```bash
   ./s3-upload-folder.sh
   ```

   The script will load the environment variables from the `.env` file and execute the Go binary with the provided parameters.

### Directly Using Binary

You can also run binary directly, passing in the required parameters:

```bash
./s3-upload-folder \
    --accessKeyId "your_access_key_id" \
    --secretAccessKey "your_secret_access_key" \
    --region "auto" \
    --endpoint "https://fly.storage.tigris.dev" \
    --bucket "your_bucket_name" \
    --folder "/path/to/your/folder" \
    --prefix "optional-prefix" \
    --workers 10
```

## Parameters

- `--accessKeyId` (required): Your AWS Access Key ID.
- `--secretAccessKey` (required): Your AWS Secret Access Key.
- `--region`: The AWS region (default: "auto").
- `--endpoint`: The custom S3-compatible endpoint (default: "https://fly.storage.tigris.dev").
- `--bucket` (required): The name of your S3 bucket.
- `--folder` (required): The path to the local folder containing the files to upload.
- `--prefix`: An optional prefix for the keys in S3.
- `--workers`: Number of parallel workers for uploading files (default: 10).

## Contributing

Feel free to fork the repository and submit pull requests. Any contributions are welcome!

## Contact

If you have any questions or issues, please create an issue in the GitHub repository or contact [o.tsyhankov@tontilagunamobile.com](mailto:o.tsyhankov@tontilagunamobile.com).



### Instructions for Setting Up the Project

1. **Clone the repository**:
   ```bash
   git clone git@github.com:magnuskma/s3-upload-folder.git
   cd s3-upload-folder
   ```

2. **Build the Go binary**:
   ```bash
   go build -o s3-upload-folder main.go
   ```

3. **Create a `.env` file**:
   ```bash
   cp .env.example .env
   ```

4. **Edit the `.env` file**:
   ```bash
   nano .env
   ```

5. **Run the upload script**:
   ```bash
   chmod +x s3-upload-folder.sh
   ./s3-upload-folder.sh
   ```