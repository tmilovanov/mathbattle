package mstd

import "time"

func DurationToDayHourMinute(d time.Duration) (int, int, int) {
	d = d.Round(time.Minute)
	day := time.Hour * 24
	days := d / day
	d = d - days*day
	hours := d / time.Hour
	d = d - hours*time.Hour
	minutes := d / time.Minute

	return int(days), int(hours), int(minutes)
}
