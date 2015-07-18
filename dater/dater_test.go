package dater

import (
	"testing"
)

func TestFriendlyDates(t *testing.T) {
	var tests = []struct {
		original string
		expected string
	}{
		{"Tue, 07 Oct 2014 20:52:39", "7 Octobre 2014"},
		{"Sun, 05 Oct 2014 21:59:34", "5 Octobre 2014"},
		{"Sat, 25 Oct 2014 18:46:26", "25 Octobre 2014"},
		{"Mon, 03 Nov 2014 20:52:19", "3 Novembre 2014"},
		{"Tue, 16 Dec 2014 21:09:03", "16 Décembre 2014"},
		{"Fri, 14 Nov 2014 12:06:25", "14 Novembre 2014"},
	}

	for _, date := range tests {
		actual := FriendlyDater(date.original)
		if actual != date.expected {
			t.Errorf("%s should lead to friendly date %s but got %s", date.original, date.expected, actual)
		}
	}
}
