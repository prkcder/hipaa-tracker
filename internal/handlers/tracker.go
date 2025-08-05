package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"time"

	// "io"
	"net/http"
	"net/url"

	// "os/exec"

	// "bytes"

	dbrepo "github.com/freshpaint/hipaa-tracker/internal/db"
	"github.com/freshpaint/hipaa-tracker/internal/models"
)

// func HandleTrackerScan(w http.ResponseWriter, r *http.Request) {
//     var req struct {
//         URL string `json:"url"`
//     }
//     json.NewDecoder(r.Body).Decode(&req)

//     cmd := exec.Command("node", "crawler/scan.js", req.URL)
//     output, err := cmd.Output()
//     if err != nil {
//         http.Error(w, "Scan failed", http.StatusInternalServerError)
//         return
//     }

//     var trackers []models.TrackerScan
//     json.Unmarshal(output, &trackers)

//     for _, tracker := range trackers {
//         tracker.ScannedURL = req.URL
//         db.SaveTracker(trackers) // write to DB
//     }

//     json.NewEncoder(w).Encode(trackers)
// }

type ScanRequest struct {
	URL string `json:"url"`
}

// func HandleTrackerScan(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		slog.Info("Incoming /api/scan request", "method", r.Method, "remote", r.RemoteAddr)

// 		if r.Method != http.MethodPost {
// 			slog.Warn("Invalid method for /api/scan", "method", r.Method)
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		var req ScanRequest

// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
// 			http.Error(w, "Invalid or missing URL", http.StatusBadRequest)
// 			return
// 		}

// 		// Call the Node-based crawler
// 		crawlerURL := "http://crawler:4000/scan"
// 		payload, _ := json.Marshal(map[string]string{"url": req.URL})
// 		resp, err := http.Post(crawlerURL, "application/json", bytes.NewBuffer(payload))

// 		if err != nil {
// 			http.Error(w, "Failed to contact crawler", http.StatusInternalServerError)
// 			slog.Error("Crawler POST failed", "error", err)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		if resp.StatusCode != http.StatusOK {
// 			body, _ := io.ReadAll(resp.Body)
// 			slog.Error("Crawler error response", "status", resp.StatusCode, "body", string(body))
// 			http.Error(w, "Crawler failed to return valid data", http.StatusInternalServerError)
// 			return
// 		}

// 		// Parse crawler response
// 		var scanResults []models.TrackerScan
// 		if err := json.NewDecoder(resp.Body).Decode(&scanResults); err != nil {
// 			http.Error(w, "Invalid crawler JSON", http.StatusInternalServerError)
// 			slog.Error("Failed to decode crawler response", "error", err)
// 			return
// 		}

// 		// Save to DB
// 		for _, scan := range scanResults {
// 			if err := dbrepo.SaveTracker(db, scan); err != nil {
// 				slog.Error("DB save failed", "scan", scan, "error", err)
// 			}
// 		}

// 		// Return the results
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(scanResults)

// 	}
// }

// Struct to match what your crawler returns
type CrawlerResponse struct {
	URL       string        `json:"url"`
	Trackers  []string      `json:"trackers"`
	Cookies   []interface{} `json:"cookies"`
	Requests  []string      `json:"requests"`
	Timestamp string        `json:"timestamp"`
}

// Helper function to extract domain from URL
func extractDomain(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr // Return original if parsing fails
	}
	return parsedURL.Host
}

// Helper function to determine risk level based on tracker domain
func getRiskLevel(domain string) string {
	highRiskDomains := []string{
		"facebook.com",
		"doubleclick.net",
		"google-analytics.com",
		"googletagmanager.com",
	}

	for _, highRisk := range highRiskDomains {
		if strings.Contains(domain, highRisk) {
			return "High"
		}
	}
	return "Medium"
}

// func HandleTrackerScan(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		slog.Info("Incoming /api/scan request", "method", r.Method, "remote", r.RemoteAddr)

// 		if r.Method != http.MethodPost {
// 			slog.Warn("Invalid method for /api/scan", "method", r.Method)
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		var req ScanRequest

// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
// 			http.Error(w, "Invalid or missing URL", http.StatusBadRequest)
// 			return
// 		}

// 		// Call the Node-based crawler
// 		crawlerURL := "http://crawler:4000/scan"
// 		payload, _ := json.Marshal(map[string]string{"url": req.URL})
// 		resp, err := http.Post(crawlerURL, "application/json", bytes.NewBuffer(payload))

// 		if err != nil {
// 			http.Error(w, "Failed to contact crawler", http.StatusInternalServerError)
// 			slog.Error("Crawler POST failed", "error", err)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		if resp.StatusCode != http.StatusOK {
// 			body, _ := io.ReadAll(resp.Body)
// 			slog.Error("Crawler error response", "status", resp.StatusCode, "body", string(body))
// 			http.Error(w, "Crawler failed to return valid data", http.StatusInternalServerError)
// 			return
// 		}

// 		// Parse crawler response with the correct structure
// 		var crawlerResponse CrawlerResponse
// 		if err := json.NewDecoder(resp.Body).Decode(&crawlerResponse); err != nil {
// 			http.Error(w, "Invalid crawler JSON", http.StatusInternalServerError)
// 			slog.Error("Failed to decode crawler response", "error", err)
// 			return
// 		}

// 		// Transform crawler response into TrackerScan objects
// 		var scanResults []models.TrackerScan

// 		for _, trackerURL := range crawlerResponse.Trackers {
// 			domain := extractDomain(trackerURL)
// 			riskLevel := getRiskLevel(domain)

// 			scan := models.TrackerScan{
// 				ScannedURL:    crawlerResponse.URL,
// 				TrackerName:   domain,
// 				TrackerDomain: domain,
// 				RiskLevel:     riskLevel,
// 				PageURL:       crawlerResponse.URL,
// 				CreatedAt:     time.Now(),
// 			}
// 			scanResults = append(scanResults, scan)
// 		}

// 		// Save to DB
// 		for _, scan := range scanResults {
// 			if err := dbrepo.SaveTracker(db, &scan); err != nil {
// 				slog.Error("DB save failed", "scan", scan, "error", err)
// 			}
// 		}

// 		// Log what we're returning
// 		slog.Info("Returning scan results", "count", len(scanResults))

// 		// Return the results
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(scanResults)
// 	}
// }

func HandleTrackerScan(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("=== DEBUG: Handler started ===", "method", r.Method, "remote", r.RemoteAddr)

		if r.Method != http.MethodPost {
			slog.Warn("=== DEBUG: Invalid method ===", "method", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ScanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			slog.Error("=== DEBUG: Invalid request body ===", "error", err, "url", req.URL)
			http.Error(w, "Invalid or missing URL", http.StatusBadRequest)
			return
		}

		slog.Info("=== DEBUG: Calling crawler ===", "url", req.URL)
		crawlerURL := "http://crawler:4000/scan"
		payload, _ := json.Marshal(map[string]string{"url": req.URL})
		resp, err := http.Post(crawlerURL, "application/json", bytes.NewBuffer(payload))

		if err != nil {
			slog.Error("=== DEBUG: Crawler POST failed ===", "error", err)
			http.Error(w, "Failed to contact crawler", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			slog.Error("=== DEBUG: Crawler error response ===", "status", resp.StatusCode, "body", string(body))
			http.Error(w, "Crawler failed to return valid data", http.StatusInternalServerError)
			return
		}

		var crawlerResponse CrawlerResponse
		if err := json.NewDecoder(resp.Body).Decode(&crawlerResponse); err != nil {
			slog.Error("=== DEBUG: Crawler JSON decode failed ===", "error", err)
			http.Error(w, "Invalid crawler JSON", http.StatusInternalServerError)
			return
		}

		slog.Info("=== DEBUG: Crawler response ===", "trackers_count", len(crawlerResponse.Trackers), "url", crawlerResponse.URL)

		var scanResults []models.TrackerScan
		for i, trackerURL := range crawlerResponse.Trackers {
			domain := extractDomain(trackerURL)
			riskLevel := getRiskLevel(domain)

			scan := models.TrackerScan{
				ScannedURL:    crawlerResponse.URL,
				TrackerName:   domain,
				TrackerDomain: domain,
				RiskLevel:     riskLevel,
				PageURL:       crawlerResponse.URL,
				CreatedAt:     time.Now(),
			}
			slog.Info("=== DEBUG: Creating scan entry ===", "index", i, "domain", domain)
			scanResults = append(scanResults, scan)
		}

		slog.Info("=== DEBUG: About to save to DB ===", "scanResults_count", len(scanResults))
		for i, scan := range scanResults {
			slog.Info("=== DEBUG: Saving scan ===", "index", i, "tracker", scan.TrackerName)
			if err := dbrepo.SaveTracker(db, &scan); err != nil {
				slog.Error("=== DEBUG: DB save failed ===", "error", err)
			} else {
				slog.Info("=== DEBUG: DB save success ===", "id", scan.ID)
			}
		}

		slog.Info("=== DEBUG: About to return response ===", "count", len(scanResults))
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(scanResults); err != nil {
			slog.Error("=== DEBUG: JSON encode failed ===", "error", err)
		} else {
			slog.Info("=== DEBUG: JSON encode success ===")
		}
	}
}
