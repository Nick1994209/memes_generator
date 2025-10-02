# Meme Generator Service

A web application for generating and viewing memes with a React frontend and Golang backend.

## Features

- Web UI for viewing created memes with thumbnail images
- Page for creating new memes
- Detailed view for individual memes with full-size images
- REST API for meme management
- Command-line tool for batch meme generation
- Docker containerization for easy deployment
- Image serving for generated memes
- Placeholder images for memes without generated images

## Architecture

The application follows a clean architecture pattern with the following components:

- **Frontend**: React 18.0 application for the user interface
- **Backend**: Golang 1.24 + Gin 1.11.0 web server
- **Storage**: File-based storage in `./data/memes/<meme_id>` directories
- **API**: RESTful endpoints for meme management

## Prerequisites

- Docker (for containerized deployment)
- Node.js 16+ (for local development)
- Go 1.24+ (for local development)

## Getting Started

### Using Docker (Recommended)

1. Build the Docker image:
   ```bash
   docker build -t memes-generator .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 -v $(pwd)/data:/app/data memes-generator
   ```

3. Access the application at http://localhost:8080

### Local Development

1. Install frontend dependencies:
   ```bash
   cd web
   npm install
   ```

2. Build the frontend:
   ```bash
   npm run build
   ```

3. Build the backend:
   ```bash
   go build -o memes-generator cmd/web/main.go
   ```

4. Run the application:
   ```bash
   ./memes-generator
   ```

5. Access the application at http://localhost:8080

## Command-Line Tool

The application includes a command-line tool for regenerating memes from existing metadata:

```bash
# Build the CLI tool
go build -o generate-meme cmd/generate/main.go

# Regenerate memes using metadata from an existing meme directory
./generate-meme --meme-path data/memes/meme_1759442111813095000
```

Note: The CLI tool only works with memes that were originally created with source images. Memes created through the web interface don't have source images saved, so they cannot be regenerated using this tool. The web interface generates memes on-the-fly and saves only the metadata.

## API Endpoints

- `GET /api/memes` - List all memes
- `POST /api/memes` - Create a new meme
- `GET /api/memes/:id` - Get a specific meme
- `DELETE /api/memes/:id` - Delete a meme
- `GET /memes/:id/image` - Get the image for a specific meme (returns actual image or placeholder)

## Project Structure

```
.
├── cmd/
│   ├── web/          # Web server main package
│   └── generate/     # CLI tool main package
├── internal/
│   ├── delivery/     # HTTP handlers
│   ├── usecase/      # Business logic
│   ├── repository/   # Data access layer
│   ├── meme/         # Meme generation logic
│   └── domain/       # Core domain models
├── web/              # React frontend
├── data/             # Meme storage directory
├── Dockerfile        # Docker configuration
└── README.md         # This file
```

## Data Storage

Memes are stored in the `./data/memes` directory with the following structure:

```
data/memes/
└── meme_<unique_id>/
    ├── metadata.json    # Meme metadata
    └── images/          # Generated meme images (for CLI-generated memes)
```

Each meme has a unique ID and is stored in its own directory with metadata. Images are stored in the `images` subdirectory for memes generated via the CLI tool. Memes created through the web interface only have metadata.