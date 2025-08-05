package db

import (
	"database/sql"
	// "encoding/json"

	"github.com/freshpaint/hipaa-tracker/internal/models"
)

var SaveTracker = func(db *sql.DB, scan *models.TrackerScan) error {

	query := `
		INSERT INTO web_tracker_scans 
			(scanned_url, tracker_name, tracker_domain, risk_level, page_url, created_at)
		VALUES 
			($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`

	err := db.QueryRow(
		query,
		scan.ScannedURL,
		scan.TrackerName,
		scan.TrackerDomain,
		scan.RiskLevel,
		scan.PageURL,
		scan.CreatedAt,
	).Scan(&scan.ID)

	return err
}

func GetAllScans(db *sql.DB) ([]models.TrackerScan, error) {

	rows, err := db.Query(`
		SELECT id, scanned_url, tracker_name, tracker_domain, risk_level, page_url, created_at 
		FROM web_tracker_scans 
		ORDER BY created_at DESC`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var scans []models.TrackerScan

	for rows.Next() {
		var s models.TrackerScan
		err := rows.Scan(
			&s.ID,
			&s.ScannedURL,
			&s.TrackerName,
			&s.TrackerDomain,
			&s.RiskLevel,
			&s.PageURL,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		scans = append(scans, s)
	}

	return scans, nil
}
