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
	formattedAverage := fmt.Sprintf("%.2f", averageDuration)

	println("HTTP Request Statistics:")
	println("Total Requests:", len(durations))
	println("Successful Requests:", successCount)
	println("Failed Requests:", len(durations)-successCount)
	println("Average Duration (seconds):", formattedAverage, "ms")
}

func PrintTotalDuration(totalDuration float64) {
	if totalDuration <= 0 {
		return
	}

	formattedTotal := fmt.Sprintf("%.2f", totalDuration*1000) // Convert to milliseconds
	println("Total time taken for all requests:", formattedTotal, "ms")
}
