## Overview

This project retrieves webcam information from the [Windy API](https://api.windy.com/webcams) and generates embedding vectors using the URLs of thumbnail images. The generated embedding vectors and webcam information are stored in MongoDB, while the downloaded thumbnail images are uploaded to Amazon S3.

## Features

- Retrieve webcam information from the Windy API
- Generate embedding vectors from thumbnail image URLs
- Store webcam information and embedding vectors in MongoDB
- Upload downloaded thumbnail images to Amazon S3

## Implementation Steps

1. Use the Windy API to fetch webcam information.
2. Generate embedding vectors based on the thumbnail image URLs.
3. Register the webcam information and embedding vectors in MongoDB.
4. Download the thumbnail images and save them to Amazon S3.
