package handlers

import (
        "database/sql"
        "encoding/json"
        "net/http"
        "strconv"
        "strings"
        "time"

        "github.com/gin-gonic/gin"
)

type DiveLog struct {
        ID       int       `json:"id"`
        DiverID  int       `json:"diver_id"`
        Depth    int       `json:"depth"`
        DateTime time.Time `json:"date_time"`
}

func LogNewDive(c *gin.Context) {
        var diveLog DiveLog

        if err := c.ShouldBindJSON(&diveLog); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        db, err := connectToDB()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
                return
        }
        defer db.Close()

        var diverID int
        err = db.QueryRow("SELECT id FROM divers WHERE id = ?", diveLog.DiverID).Scan(&diverID)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diver ID"})
                return
        }

        if diveLog.Depth > MaxDepth {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Exceeded maximum allowed depth"})
                return
        }

        var totalDives int
        err = db.QueryRow("SELECT COUNT(*) FROM divelogs WHERE diver_id = ?", diveLog.DiverID).Scan(&totalDives)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dive logs"})
                return
        }
        if totalDives >= MaxAllowedDives {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Exceeded maximum allowed dives"})
                return
        }

        var lastDiveDepth int
        err = db.QueryRow("SELECT depth FROM divelogs WHERE diver_id = ? ORDER BY id DESC LIMIT 1", diveLog.DiverID).Scan(&lastDiveDepth)
        if err != nil {
                if err == sql.ErrNoRows {
                        // No previous dive logs for the diver, skip the depth check
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last dive depth"})
                        return
                }
        } else {
                if diveLog.Depth >= lastDiveDepth {
                        c.JSON(http.StatusBadRequest, gin.H{"error": "New dive depth must be less than the previous dive depth"})
                        return
                }
        }

        minTimeInterval := calculateMinAllowedTimeInterval(diveLog.Depth, totalDives)

        var lastDiveTime time.Time
        err = db.QueryRow("SELECT date_time FROM divelogs WHERE diver_id = ? ORDER BY id DESC LIMIT 1", diveLog.DiverID).Scan(&lastDiveTime)
        if err != nil {
                if err != sql.ErrNoRows {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve last dive time"})
                        return
                }
        } else {
                elapsedTime := time.Since(lastDiveTime)
                if elapsedTime < minTimeInterval {
                        c.JSON(http.StatusBadRequest, gin.H{"error": "Minimum allowed time interval between dives violated", "minimum_interval": minTimeInterval.String()})
                        return
                }
        }

        _, err = db.Exec("INSERT INTO divelogs (diver_id, depth, date_time) VALUES (?, ?, NOW())", diveLog.DiverID, diveLog.Depth)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log new dive"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"message": "New dive logged successfully"})
}
