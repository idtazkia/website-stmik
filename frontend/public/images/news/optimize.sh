#!/bin/bash
# Script untuk optimize dan resize foto berita
# Run dengan: bash optimize.sh

cd "$(dirname "$0")"

echo "Optimizing photos for web..."

# Featured image (DSC_8095) - max 1200px, quality 85%
if [ -f "DSC_8095.JPG" ]; then
    echo "Processing DSC_8095.JPG (featured image)..."
    convert DSC_8095.JPG -resize 1200x -quality 85 DSC_8095_opt.JPG && mv DSC_8095_opt.JPG DSC_8095.JPG
fi

# Other photos - max 800px, quality 85%
for file in DSC_8066.JPG DSC_8079.JPG DSC_8102.JPG DSC_8108.JPG DSC_8112.JPG DSC_8113.JPG; do
    if [ -f "$file" ]; then
        echo "Processing $file..."
        convert "$file" -resize 800x -quality 85 "${file%.JPG}_opt.JPG" && mv "${file%.JPG}_opt.JPG" "$file"
    fi
done

echo "Done! Showing file sizes:"
ls -lh DSC_*.JPG | awk '{print $5, $9}'
