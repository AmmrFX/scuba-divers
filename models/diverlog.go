package models

import "time"

type DiveLog struct {
        ID       int       `json:"id"`
        DiverID  int       `json:"diver_id"`
        Depth    int       `json:"depth"`
        DateTime time.Time `json:"date_time"`
}
