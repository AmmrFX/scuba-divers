package utils

import "time"

func CalculateMinAllowedTimeInterval(depth, diveCount int) time.Duration {
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
