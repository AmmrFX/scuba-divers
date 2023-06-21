package handlers

import (
        "database/sql"
        "net/http"

        "github.com/gin-gonic/gin"
)

type Report struct {
        DiverID     int    `json:"diver_id"`
        TotalDives  int    `json:"total_dives"`
        MaxDepth    int    `json:"max_depth"`
        LatestDive  string `json:"latest_dive"`
        EarliestDive string `json:"earliest_dive"`
}

func GenerateDiveReports(c *gin.Context) {
        db, err := connectToDB()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
                return
        }
        defer db.Close()

        rows, err := db.Query("SELECT diver_id, COUNT(*) as total_dives, MAX(depth) as max_depth, MAX(date_time) as latest_dive, MIN(date_time) as earliest_dive FROM divelogs GROUP BY diver_id")
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate dive reports"})
                return
        }
        defer rows.Close()

        reports := []Report{}
        for rows.Next() {
                var report Report
                err := rows.Scan(&report.DiverID, &report.TotalDives, &report.MaxDepth, &report.LatestDive, &report.EarliestDive)
                if err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate dive reports"})
                        return
                }
                reports = append(reports, report)
        }

        c.JSON(http.StatusOK, gin.H{"dive_reports": reports})
}
