# Electrolux Backend Developer Test (Fit Tracker)

## Tech Stack

1. Go
1. MongoDB
1. Docker

## Directory & File Structure

### Directory

1. `ingestor` - Service that poll data from websocket
1. `database` - Database service for storing ingested data
1. `api` - HTTP server for exposing API service to frontend (second task)
1. `testdata` - Sample data for unit testing

### File

1. `ingestions.json` - Exported data from database (10 minutes)
1. `Dockerfile` - Docker image for Go
1. `compose.yml` - Docker compose setup that includes Go (built from Dockerfile) and MongoDB
1. `compose.override.yml` - Docker compose setup with override settings for development environment
1. `openapi.yml` - OpenAPI specification for the API server (second task)
