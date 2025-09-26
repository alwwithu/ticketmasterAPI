package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// POST /ingest/{marketplace}
func ingestHandler(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		market := c.Param("marketplace")
		if market == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing marketplace"})
			return
		}

		// Configurable parameters
		size := 200      // page size (<= 500 per API)
		page := 0        // start page
		maxItems := 1000 // deep-paging safety (per TM docs)
		client := &http.Client{Timeout: 15 * time.Second}
		totalIngested := 0
		pagesFetched := 0

		start := time.Now()

		// Clear previous data for this marketplace (optional)
		store.Lock()
		store.data[market] = []Event{}
		store.Unlock()

		for {
			// Build request URL (we treat marketplace as countryCode by default)
			url := fmt.Sprintf("https://app.ticketmaster.com/discovery/v2/events.json?countryCode=%s&size=%d&page=%d&apikey=%s",
				market, size, page, apiKey)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			}
			// If you prefer header auth (some partner endpoints support x-api-key), use:
			// req.Header.Set("x-api-key", apiKey)

			resp, err := client.Do(req)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "error calling ticketmaster: " + err.Error()})
				return
			}
			if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "ticketmaster auth error"})
				resp.Body.Close()
				return
			}
			if resp.StatusCode >= 500 {
				c.JSON(http.StatusBadGateway, gin.H{"error": "ticketmaster server error"})
				resp.Body.Close()
				return
			}

			var tm TMResponse
			if err := json.NewDecoder(resp.Body).Decode(&tm); err != nil {
				resp.Body.Close()
				c.JSON(http.StatusBadGateway, gin.H{"error": "failed to decode response: " + err.Error()})
				return
			}
			resp.Body.Close()

			pagesFetched++

			if tm.Embedded != nil && len(tm.Embedded.Events) > 0 {
				store.Lock()
				store.data[market] = append(store.data[market], tm.Embedded.Events...)
				store.Unlock()
				totalIngested += len(tm.Embedded.Events)
			}

			// Stop conditions
			if tm.Page == nil {
				// no paging info — just stop
				break
			}
			// Stop if we've fetched all pages
			if page >= tm.Page.TotalPages-1 {
				break
			}
			// Respect deep paging limits (size * page < maxItems)
			if (page+1)*size >= maxItems {
				// reached deep-paging cap (avoid endless paging)
				break
			}

			page++

			// Basic pacing: sleep to avoid hitting default quotas (safe default ~2 req/sec)
			time.Sleep(510 * time.Millisecond)
		}

		duration := time.Since(start)
		out := map[string]interface{}{
			"marketplace":    market,
			"pagesFetched":   pagesFetched,
			"eventsIngested": totalIngested,
			"duration_ms":    duration.Milliseconds(),
		}
		c.JSON(http.StatusOK, out)
	}
}

// GET /events/{marketplace}
func getEventsHandler(c *gin.Context) {
	market := c.Param("marketplace")

	store.RLock()
	events := store.data[market]
	store.RUnlock()

	if events == nil {
		events = []Event{}
	}

	// Optional: support ?limit= and ?offset=
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = len(events)
	}
	end := offset + limit
	if end > len(events) {
		end = len(events)
	}
	if offset < 0 {
		offset = 0
	}

	resp := struct {
		Count  int     `json:"count"`
		Events []Event `json:"events"`
	}{
		Count:  len(events),
		Events: events[offset:end],
	}

	c.JSON(http.StatusOK, resp)
}
