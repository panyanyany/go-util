package wall_clock_util

import (
	"time"
)

type TimeRange [2]Clock

func (r TimeRange) ClockInRange(cur Clock) (ret bool) {
	if r[0].LessThan(r[1]) {
		ret = r[0].Seconds() <= cur.Seconds() && cur.Seconds() <= r[1].Seconds()
	} else {
		ret = r[1].Seconds() <= cur.Seconds() || cur.Seconds() <= r[0].Seconds()
	}
	return
}
func (r TimeRange) TimeInRange(cur time.Time) bool {
	c := NewClockFromTime(cur)
	return r.ClockInRange(*c)
}
func (r TimeRange) NowInRange() bool {
	return r.TimeInRange(time.Now())
}
