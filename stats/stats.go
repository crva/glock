package stats

import "fmt"

func PrintHttpReport(durations []float64, successCount int) {
	if len(durations) == 0 {
		return
	}

	var totalDuration float64
	for _, duration := range durations {
		totalDuration += duration
	}

	averageDuration := (totalDuration / float64(len(durations))) * 1000 // Convert to milliseconds
	formattedAverage := fmt.Sprintf("%.2f", averageDuration)            // Format to 2 decimal places

	println("HTTP Request Statistics:")
	println("Total Requests:", len(durations))
	println("Successful Requests:", successCount)
	println("Failed Requests:", len(durations)-successCount)
	println("Average Duration (seconds):", formattedAverage, "ms")
}
