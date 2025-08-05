package models

import (
	"time"
)


type TrackerScan struct {
    ID           int       `json:"id"`
    ScannedURL   string    `json:"scanned_url"`
    TrackerName  string    `json:"tracker_name"`
    TrackerDomain string   `json:"tracker_domain"`
    RiskLevel    string    `json:"risk_level"`
    PageURL      string    `json:"page_url"`
    CreatedAt    time.Time `json:"created_at"`
}
