package custom_utils

import "time"

func CreateRandomDate(minTimestampMs int64, maxTimestampMs int64) string {
	if minTimestampMs == 0 {
		minTimestampMs = 172000000983
	}
	if maxTimestampMs == 0 {
		maxTimestampMs = 1727500009835
	}

	minTimestampSecs := minTimestampMs / 1000
	maxTimestampSecs := maxTimestampMs / 1000

	randomTimestamp := CreateRandomNumber(minTimestampSecs, maxTimestampSecs)
	unixTimeUTC := time.Unix(randomTimestamp, 0)
	formatted := unixTimeUTC.Format("2006-01-02")

	return formatted
}
