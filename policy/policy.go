package policy

import "time"

const BeverageResyncIntervalDays = 3
const BeverageDiscardThresholdDays = 35

type Clock interface {
	Now() time.Time
}

var TimeProvider Clock = nowProvider{}

func BeverageResyncThresholdTime() time.Time {
	return TimeAgoDays(BeverageResyncIntervalDays)
}

func BeverageDiscardThresholdTime() time.Time {
	return TimeAgoDays(BeverageDiscardThresholdDays)
}

func TimeAgoDays(days int) time.Time {
	return TimeProvider.Now().Add(-24 * time.Hour * time.Duration(days))
}

type nowProvider struct{}

func (nowProvider) Now() time.Time {
	return time.Now()
}
