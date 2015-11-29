package hlc

import (
	"fmt"
	"time"
)

type TimeSupplier func() time.Time

type HLCTime struct {
	wallTime    time.Time
	logicalTime int16
}

type HLCClock struct {
	t            HLCTime
	timeSupplier TimeSupplier // Wall time supplier
}

func NewHLCTime(wallTime time.Time, logical int16) HLCTime {
	return HLCTime{
		wallTime:    wallTime,
		logicalTime: logical,
	}
}

func NewHLCClock(timeSupplier TimeSupplier) *HLCClock {
	return &HLCClock{
		t:            NewHLCTime(time.Unix(0, 0), 0),
		timeSupplier: timeSupplier,
	}
}

func (t HLCTime) String() string {
	return fmt.Sprintf("%d+%d", t.wallTime.UnixNano(), t.logicalTime)
}

func (c *HLCClock) GetTime() HLCTime {
	wallTime := c.timeSupplier()

	if c.t.wallTime.Before(wallTime) {
		c.t.wallTime = wallTime
		c.t.logicalTime = 0
	} else {
		c.t.logicalTime++
	}

	return c.t
}

func (c *HLCClock) UpdateTime(event HLCTime) HLCTime {
	wallTime := c.timeSupplier()

	if wallTime.After(event.wallTime) && wallTime.After(c.t.wallTime) {
		c.t.wallTime = wallTime
		c.t.logicalTime = 0
	} else if event.wallTime.After(c.t.wallTime) {
		c.t.wallTime = event.wallTime
		c.t.logicalTime++
	} else if c.t.wallTime.After(event.wallTime) {
		c.t.logicalTime++
	} else {
		if event.logicalTime > c.t.logicalTime {
			c.t.logicalTime = event.logicalTime
		}

		c.t.logicalTime += 1
	}

	return c.t
}
