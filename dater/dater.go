package dater

import (
	"strconv"
	"time"
)

var months = map[string]string{
	"January":   "Janvier",
	"February":  "Février",
	"March":     "Mars",
	"April":     "Avril",
	"May":       "Mai",
	"June":      "Juin",
	"Jully":     "Juillet",
	"August":    "Août",
	"September": "Septembre",
	"October":   "Octobre",
	"November":  "Novembre",
	"December":  "Décembre",
}

// Transform date in from a specific layout (Mon, 02 Jan 2006 15:04:05) into another one more friendly for users : 2 Janvier 2014
func FriendlyDater(date string) string {
	parsed, err := time.Parse("Mon, 02 Jan 2006 15:04:05", date)
	if err == nil {
		year, month, day := parsed.Date()
		translatedMonth, _ := months[month.String()]
		return strconv.Itoa(day) + " " + translatedMonth + " " + strconv.Itoa(year)
	} else {
		return date
	}
}
