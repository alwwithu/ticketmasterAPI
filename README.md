# Ticketmaster Ingest API

A Go-based API service that fetches events from Ticketmaster API and provides a simple web interface for testing.

## Features

- **Ingest API**: Fetch events from Ticketmaster for a specific marketplace
- **Events API**: Retrieve stored events with pagination support
- **Web Frontend**: Simple HTML interface for testing the APIs

## Setup

1. **Set your Ticketmaster API Key**:
   ```bash
   export TICKETMASTER_API_KEY="your_api_key_here"
   ```
   Or create a `.env` file:
   ```
   TICKETMASTER_API_KEY=your_api_key_here
   ```

2. **Run the server**:
   ```bash
   go run main.go
   ```

3. **Access the web interface**:
   Open your browser and go to `http://localhost:8080`

## API Endpoints

### POST /ingest/{marketplace}
Fetches events from Ticketmaster API for the specified marketplace.

**Parameters:**
- `marketplace`: Country code (e.g., "US", "CA", "GB")

**Response:**
```json
{
  "marketplace": "US",
  "pagesFetched": 5,
  "eventsIngested": 1000,
  "duration_ms": 2500
}
```

### GET /events/{marketplace}
Retrieves stored events for the specified marketplace.

**Parameters:**
- `marketplace`: Country code
- `limit` (optional): Number of events to return (default: all)
- `offset` (optional): Number of events to skip (default: 0)

**Response:**
```json
{
  "count": 1000,
  "events": [
    {
      "id": "event_id",
      "name": "Event Name",
      "url": "https://...",
      "dates": {
        "start": {
          "localDate": "2024-01-01",
          "dateTime": "2024-01-01T20:00:00Z"
        }
      }
    }
  ]
}
```

## Web Interface

The web interface provides:

- **Ingest Form**: Select a marketplace and fetch events from Ticketmaster
- **Events Form**: View stored events with pagination controls
- **Real-time Responses**: See API responses formatted as JSON
- **Loading Indicators**: Visual feedback during API calls

## Supported Marketplaces

- US (United States)
- CA (Canada)
- GB (United Kingdom)
- AU (Australia)
- DE (Germany)
- FR (France)
- ES (Spain)
- IT (Italy)
- NL (Netherlands)
- SE (Sweden)
- NO (Norway)
- DK (Denmark)
- FI (Finland)

## Rate Limiting

The service includes built-in rate limiting (510ms between requests) to respect Ticketmaster's API limits.
