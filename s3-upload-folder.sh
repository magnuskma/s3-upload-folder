#!/bin/bash

# Function to check if the binary exists and download it if not
function check_and_download_binary() {
    if [ ! -f "./s3-upload-folder" ]; then
        echo "Binary not found. Downloading the binary using curl..."

        # Download the binary using curl
        curl -L -o s3-upload-folder https://github.com/magnuskma/s3-upload-folder/raw/main/s3-upload-folder

        # Make the binary executable
        chmod +x s3-upload-folder

        echo "Binary downloaded and made executable."
    else
        echo "Binary found. Proceeding with the upload..."
    fi
}

# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Check if the binary exists, if not, download it
check_and_download_binary

# Execute the binary with the provided parameters
./s3-upload-folder \
    --accessKeyId "$AWS_ACCESS_KEY_ID" \
    --secretAccessKey "$AWS_SECRET_ACCESS_KEY" \
    --region "$AWS_REGION" \
    --endpoint "$AWS_ENDPOINT_URL_S3" \
    --bucket "$S3_BUCKET_NAME" \
    --folder "$FOLDER_TO_UPLOAD" \
    --prefix "$S3_KEY_PREFIX" \
    --workers 10
