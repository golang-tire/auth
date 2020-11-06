#!/usr/bin/env bash

PACKAGE_NAME="auth";
OUTPUT_DIR="internal/proto/v1"

for f in $(find -type f -name "*.json"); do
  output_filename=$(basename $f)
  tire swagger-to-go $f --pkg $PACKAGE_NAME --out $OUTPUT_DIR/$output_filename.go
done