#!/bin/bash

# Default values
DEFAULT_DOMAIN="https://peachystudio.com"
DEFAULT_CONFIG="configs/default.yaml"
MONGO_URI=""  # Leave empty by default

# Help function
show_help() {
    echo "Usage: ./run.sh [domain] [-m mongo_uri] [-c config_file]"
    echo
    echo "Options:"
    echo "  domain      The domain to crawl (default: $DEFAULT_DOMAIN)"
    echo "  -m, --mongo MongoDB connection string (optional)"
    echo "  -c, --config Configuration file path (default: $DEFAULT_CONFIG)"
    echo "  -h, --help  Show this help message"
    echo
    echo "Example:"
    echo "  ./run.sh https://golang.org"
    echo "  ./run.sh -m \"mongodb+srv://user:pass@cluster.mongodb.net\""
    echo "  ./run.sh -c \"configs/custom.yaml\""
    echo "  ./run.sh https://golang.org -m \"mongodb+srv://user:pass@cluster.mongodb.net\" -c \"configs/custom.yaml\""
}

# Parse command line arguments
DOMAIN="$DEFAULT_DOMAIN"
CONFIG="$DEFAULT_CONFIG"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -m|--mongo)
            MONGO_URI="$2"
            shift 2
            ;;
        -c|--config)
            CONFIG="$2"
            shift 2
            ;;
        -*)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
        *)
            # First non-option argument is the domain
            if [[ "$1" =~ ^https?:// ]]; then
                DOMAIN="$1"
                shift
            else
                echo "Error: Domain must start with http:// or https://"
                exit 1
            fi
            ;;
    esac
done

# Build the crawler
echo "Building crawler..."
go build -o crawler ./cmd/crawler

# Construct the command
CMD="./crawler -seed \"$DOMAIN\" -config \"$CONFIG\""
if [ ! -z "$MONGO_URI" ]; then
    CMD="$CMD -mongo \"$MONGO_URI\""
fi

# Run the crawler
echo "Starting crawler with domain: $DOMAIN"
echo "Using config file: $CONFIG"
eval $CMD 