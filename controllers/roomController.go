package controllers

import (
	"airbnb-api/database"
	"airbnb-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"fmt"
	"time"
)

// GetRoomMetrics retrieves occupancy and rate metrics for a specific room.
func GetRoomMetrics(c *gin.Context) {
	roomIDStr := c.Param("room_id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room_id"})
		return
	}

	var room models.Room

	// Fetch room from the database. Adjust Preload if AvailableDates is a related table.
	if err := database.DB.First(&room, "room_id = ?", roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// Debug: Print room details after fetching from DB
	fmt.Printf("Room Details: %+v\n", room) // Print full room details
	fmt.Printf("Available Dates: %v\n", room.AvailableDates)

	// Calculate Occupancy Percentage (Next 5 Months)
	now := time.Now()
	occupancy := make(map[string]float64)

	for i := 0; i < 5; i++ {
		targetDate := now.AddDate(0, i, 0) // Add months incrementally
		year := targetDate.Year()
		month := targetDate.Month()

		fmt.Printf("Calculating for Year: %d, Month: %s\n", year, month)

		totalDays := daysInMonth(month, year)
		availableDays := countAvailableDays(room.AvailableDates, month, year)

		// Occupancy is (Total Days - Available Days) / Total Days * 100
		occupiedDays := totalDays - availableDays
		var occupancyPercentage float64
		if totalDays > 0 {
			occupancyPercentage = (float64(occupiedDays) / float64(totalDays)) * 100
		} else {
			occupancyPercentage = 0
		}
		// Use YYYY-MM format to ensure uniqueness
		key := fmt.Sprintf("%d-%02d", year, month)
		occupancy[key] = occupancyPercentage

		fmt.Printf("Month: %s | Total Days: %d | Available Days: %d | Occupancy: %.2f%%\n", key, totalDays, availableDays, occupancyPercentage)
	}

	// Calculate Night Rates (Next 30 Days)
	nowPlus30 := now.AddDate(0, 0, 30)
	var rates []float64

	for _, dateStr := range room.AvailableDates {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Printf("Error parsing date: %s\n", dateStr)
			continue
		}
		fmt.Printf("Checking date: %v\n", date) // Debug date
		if date.After(now) && date.Before(nowPlus30) {
			fmt.Printf("Adding rate for date: %v\n", date)
			rates = append(rates, room.RatePerNight)
		}
	}

	var averageRate, maxRate, minRate *float64
	if len(rates) > 0 {
		sumRates := sum(rates)
		avg := sumRates / float64(len(rates))
		averageRate = &avg
		maxVal := max(rates)
		maxRate = &maxVal
		minVal := min(rates)
		minRate = &minVal
	}

	// Response
	response := gin.H{
		"occupancy_percentage": occupancy,
		"rates_next_30_days": gin.H{
			"average_rate": averageRate,
			"highest_rate": maxRate,
			"lowest_rate":  minRate,
		},
	}

	c.JSON(http.StatusOK, response)
}

// daysInMonth returns the number of days in a given month and year.
func daysInMonth(month time.Month, year int) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// countAvailableDays counts the number of available days in a specific month and year.
func countAvailableDays(dates []string, month time.Month, year int) int {
	count := 0
	for _, dateStr := range dates {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Printf("Error parsing date: %s\n", dateStr)
			continue
		}
		fmt.Printf("Parsed Date: %v\n", date)

		// Compare the year and month
		if date.Year() == year && date.Month() == month {
			count++
		}
	}
	return count
}

// sum calculates the total of a slice of float64 numbers.
func sum(array []float64) float64 {
	var total float64
	for _, num := range array {
		total += num
	}
	return total
}

// max returns the maximum value in a slice of float64 numbers.
func max(array []float64) float64 {
	if len(array) == 0 {
		return 0
	}
	maxVal := array[0]
	for _, num := range array[1:] {
		if num > maxVal {
			maxVal = num
		}
	}
	return maxVal
}

// min returns the minimum value in a slice of float64 numbers.
func min(array []float64) float64 {
	if len(array) == 0 {
		return 0
	}
	minVal := array[0]
	for _, num := range array[1:] {
		if num < minVal {
			minVal = num
		}
	}
	return minVal
}
