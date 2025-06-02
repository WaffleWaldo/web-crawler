#!/bin/bash

# Browse Crawled Content Script
# This script helps you explore the saved page content from the web crawler

echo "🕷️  Ultra-High-Performance Web Crawler - Content Browser"
echo "========================================================"
echo

# Check if crawled_content directory exists
if [ ! -d "crawled_content" ]; then
    echo "❌ No crawled_content directory found. Run the crawler first!"
    exit 1
fi

# Show statistics
echo "📊 Content Statistics:"
echo "----------------------"
total_files=$(find crawled_content -name "*.html" | wc -l | tr -d ' ')
total_size=$(du -sh crawled_content | cut -f1)
domains=$(ls crawled_content | wc -l | tr -d ' ')

echo "📄 Total HTML files: $total_files"
echo "💾 Total size: $total_size"
echo "🌐 Domains crawled: $domains"
echo

# List domains
echo "🌐 Crawled Domains:"
echo "-------------------"
for domain in crawled_content/*/; do
    if [ -d "$domain" ]; then
        domain_name=$(basename "$domain")
        file_count=$(find "$domain" -name "*.html" | wc -l | tr -d ' ')
        domain_size=$(du -sh "$domain" | cut -f1)
        echo "  📁 $domain_name ($file_count files, $domain_size)"
    fi
done
echo

# Show recent files
echo "📝 Recent Files (last 10):"
echo "---------------------------"
find crawled_content -name "*.html" -exec ls -lt {} + | head -10 | while read line; do
    filename=$(echo "$line" | awk '{print $NF}')
    size=$(echo "$line" | awk '{print $5}')
    echo "  📄 $(basename "$filename") ($(numfmt --to=iec $size))"
done
echo

# Interactive options
echo "🔍 Browse Options:"
echo "------------------"
echo "1. List all files in a domain"
echo "2. View a specific file"
echo "3. Search for files containing text"
echo "4. Show file metadata"
echo "5. Exit"
echo

read -p "Choose an option (1-5): " choice

case $choice in
    1)
        echo
        echo "Available domains:"
        ls crawled_content/
        echo
        read -p "Enter domain name: " domain
        if [ -d "crawled_content/$domain" ]; then
            echo
            echo "Files in $domain:"
            ls -la "crawled_content/$domain/"
        else
            echo "❌ Domain not found!"
        fi
        ;;
    2)
        echo
        read -p "Enter file path (e.g., crawled_content/domain/file.html): " filepath
        if [ -f "$filepath" ]; then
            echo
            echo "📄 Content of $filepath:"
            echo "========================"
            head -30 "$filepath"
            echo
            echo "... (showing first 30 lines)"
            echo
            read -p "View full file? (y/n): " view_full
            if [ "$view_full" = "y" ]; then
                less "$filepath"
            fi
        else
            echo "❌ File not found!"
        fi
        ;;
    3)
        echo
        read -p "Enter search text: " search_text
        echo
        echo "🔍 Files containing '$search_text':"
        grep -r -l "$search_text" crawled_content/ 2>/dev/null | head -10
        ;;
    4)
        echo
        read -p "Enter file path: " filepath
        if [ -f "$filepath" ]; then
            echo
            echo "📋 Metadata for $filepath:"
            echo "=========================="
            head -15 "$filepath" | grep -A 10 "CRAWLED PAGE METADATA"
        else
            echo "❌ File not found!"
        fi
        ;;
    5)
        echo "👋 Goodbye!"
        exit 0
        ;;
    *)
        echo "❌ Invalid option!"
        ;;
esac 