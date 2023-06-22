/*
Restful API Scuba Divers 
===================================================
Title:         Scuba Divers
Release date:  21/06/2023
Author:        Amr Effat
Copyright:     (c) Amr Effat
===================================================
*/

package main

import (
        "database/sql"
        "errors"
        "log"
        "net/http"
        "strconv"
        "strings"
        "time"

        "github.com/gin-gonic/gin"
        _ "github.com/go-sql-driver/mysql"
)

const (
        MaxDepth             = 40
        MaxAllowedDives      = 10
        MinAllowedTimeInterval = 5
)

type Diver struct {
        ID       int    `json:"id"`
        Name     string `json:"name"`
        DiverEqp string `json:"diverEqp"`
}

type DiveLog struct {
        ID       int       `json:"id"`
        DiverID  int       `json:"diver_id"`
        Depth    int       `json:"depth"`
        DateTime time.Time `json:"date_time"`
}

func main() {
        router := gin.Default()

        router.POST("/divers", createDiverProfile)
        router.POST("/divelogs", logNewDive)
        router.GET("/divelogs", getDiverDiveLogs)
        router.GET("/reports", generateDiveReports)
        router.GET("/divers", getDiversInformation)

        db, err := connectToDB()
        if err != nil {
                log.Fatalf("Failed to connect to database: %v", err)
        }
        defer db.Close()

        _, err = db.Exec("INSERT INTO divers (name, diverEqp) VALUES (?, ?)", "Debug Diver", "Debug Equipment")
        if err != nil {
                log.Fatalf("Failed to write debug data to database: %v", err)
        }

        if err := router.Run(":8080"); err != nil {
                log.Fatalf("Failed to start server: %v", err)
        }
}


func createDiverProfile(c *gin.Context) {
        var diver Diver

        if err := c.ShouldBindJSON(&diver); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        db, err := connectToDB()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
                return
        }
        defer db.Close()

        result, err := db.Exec("INSERT INTO divers (name, diverEqp) VALUES (?, ?)", diver.Name, diver.DiverEqp)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create diver profile"})
                return
        }

        diverID, err := result.LastInsertId()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve diver ID"})
                return
        }

        c.JSON(http.StatusOK, gin.H{"diver_id": diverID})
}

func logNewDive(c *gin.Context) {
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

func getDiverDiveLogs(c *gin.Context) {
        diverID, err := strconv.Atoi(c.Query("diver_id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diver ID"})
                return
        }

        db, err := connectToDB()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
                return
        }
        defer db.Close()

        rows, err := db.Query("SELECT id, diver_id, depth, date_time FROM divelogs WHERE diver_id = ?", diverID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dive logs"})
                return
        }
        defer rows.Close()

        diveLogs := []DiveLog{}
        for rows.Next() {
                var diveLog DiveLog
                err := rows.Scan(&diveLog.ID, &diveLog.DiverID, &diveLog.Depth, &diveLog.DateTime)
                if err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dive logs"})
                        return
                }
                diveLogs = append(diveLogs, diveLog)
        }

        c.JSON(http.StatusOK, gin.H{"dive_logs": diveLogs})
}

func generateDiveReports(c *gin.Context) {
        db, err := connectToDB()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
                return
        }
        defer db.Close()

        rows, err := db.Query("SELECT diver_id, COUNT(*) as total_dives FROM divelogs GROUP BY diver_id")
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate dive reports"})
                return
        }
        defer rows.Close()

        reports := make(map[int]int)
        for rows.Next() {
                var diverID, totalDives int
                err := rows.Scan(&diverID, &totalDives)
                if err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate dive reports"})
                        return
                }
                reports[diverID] = totalDives
        }

        c.JSON(http.StatusOK, gin.H{"dive_reports": reports})
}

func getDiversInformation(c *gin.Context) {
        diverIDs := c.Query("diver_ids")

        ids, err := parseDiverIDs(diverIDs)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diver IDs"})
                return
        }

        db, err := connectToDB()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
                return
        }
        defer db.Close()

        query := "SELECT id, name, diverEqp FROM divers WHERE id IN (?)"
        query = strings.Replace(query, "(?)", "(?)"+strings.Repeat(", (?)", len(ids)-1), 1)
        rows, err := db.Query(query, ids...)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve divers information"})
                return
        }
        defer rows.Close()

        divers := []Diver{}
        for rows.Next() {
                var diver Diver
                err := rows.Scan(&diver.ID, &diver.Name, &diver.DiverEqp)
                if err != nil {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve divers information"})
                        return
                }
                divers = append(divers, diver)
        }

        c.JSON(http.StatusOK, gin.H{"divers": divers})
}

func connectToDB() (*sql.DB, error) {
        db, err := sql.Open("mysql", "root:ihackstuff@tcp(db:3306)/effat")
        if err != nil {
                return nil, err
        }
        return db, nil
}

func calculateMinAllowedTimeInterval(depth, diveCount int) time.Duration {
        minTimeInterval := time.Duration(MinAllowedTimeInterval) * time.Minute
        if diveCount > 0 {
                minTimeInterval += time.Duration(depth*intPow(2, diveCount)) * time.Minute
        }
        return minTimeInterval
}

func intPow(base, exp int) int {
        result := 1
        for i := 0; i < exp; i++ {
                result *= base
        }
        return result
}

func parseDiverIDs(diverIDs string) ([]interface{}, error) {
        ids := []interface{}{}
        idRanges := strings.Split(diverIDs, ",")
        for _, idRange := range idRanges {
                if strings.Contains(idRange, "-") {
                        rangeLimits := strings.Split(idRange, "-")
                        if len(rangeLimits) != 2 {
                                return nil, errors.New("Invalid diver IDs")
                        }
                        start, err := strconv.Atoi(rangeLimits[0])
                        if err != nil {
                                return nil, err
                        }
                        end, err := strconv.Atoi(rangeLimits[1])
                        if err != nil {
                                return nil, err
                        }
                        for i := start; i <= end; i++ {
                                ids = append(ids, i)
                        }
                } else {
                        id, err := strconv.Atoi(idRange)
                        if err != nil {
                                return nil, err
                        }
                        ids = append(ids, id)
                }
        }
        return ids, nil
}
