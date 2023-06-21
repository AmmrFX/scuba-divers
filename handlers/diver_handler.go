package handlers

import (
        "database/sql"
        "encoding/json"
        "net/http"

        "github.com/gin-gonic/gin"
)

type Diver struct {
        ID       int    `json:"id"`
        Name     string `json:"name"`
        DiverEqp string `json:"diverEqp"`
}

func CreateDiverProfile(c *gin.Context) {
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

func GetDiverByID(c *gin.Context) {
        diverID := c.Param("id")

        db, err := connectToDB()
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
                return
        }
        defer db.Close()

        var diver Diver
        err = db.QueryRow("SELECT id, name, diverEqp FROM divers WHERE id = ?", diverID).Scan(&diver.ID, &diver.Name, &diver.DiverEqp)
        if err != nil {
                if err == sql.ErrNoRows {
                        c.JSON(http.StatusNotFound, gin.H{"error": "Diver not found"})
                } else {
                        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve diver profile"})
                }
                return
        }

        c.JSON(http.StatusOK, gin.H{"diver": diver})
}

func connectToDB() (*sql.DB, error) {
        db, err := sql.Open("mysql", "root:ihackstuff@tcp(localhost:3306)/effat")
        if err != nil {
                return nil, err
        }
        return db, nil
}
