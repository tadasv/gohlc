package hlc_test

import (
	gohlc "github.com/tadasv/gohlc"
	"testing"
	"time"
)

func syntheticNow(times []time.Time) gohlc.TimeSupplier {
	ch := make(chan time.Time, len(times))
	for _, v := range times {
		ch <- v
	}

	return func() time.Time {
		tt := <-ch
		return tt
	}
}

func checkTime(t *testing.T, res gohlc.HLCTime, expected gohlc.HLCTime) {
	if res != expected {
		t.Fatalf("Expected %v got %v", expected, res)
	}
}

func TestHLCClock(t *testing.T) {
	times := []time.Time{
		time.Unix(1, 0),
		time.Unix(1, 0),
		time.Unix(0, 9),
		time.Unix(2, 0),
		time.Unix(3, 0),
		time.Unix(3, 0),
		time.Unix(3, 0),
		time.Unix(3, 0),
		time.Unix(3, 5),
		time.Unix(5, 0),
		time.Unix(4, 9),
		time.Unix(0, 0),
	}

	expected := []gohlc.HLCTime{
		gohlc.NewHLCTime(time.Unix(1, 0), 0),
		gohlc.NewHLCTime(time.Unix(1, 0), 1),
		gohlc.NewHLCTime(time.Unix(1, 0), 2),
		gohlc.NewHLCTime(time.Unix(2, 0), 0),
		gohlc.NewHLCTime(time.Unix(3, 0), 0),
		gohlc.NewHLCTime(time.Unix(3, 0), 1),
		gohlc.NewHLCTime(time.Unix(3, 0), 2),
		gohlc.NewHLCTime(time.Unix(3, 0), 100),
		gohlc.NewHLCTime(time.Unix(4, 4), 101),
		gohlc.NewHLCTime(time.Unix(5, 0), 0),
		gohlc.NewHLCTime(time.Unix(5, 0), 100),
		gohlc.NewHLCTime(time.Unix(5, 0), 101),
	}

	// Wall clock does not move
	clock := gohlc.NewHLCClock(syntheticNow(times))

	tt := clock.GetTime()
	checkTime(t, tt, expected[0])

	tt = clock.GetTime()
	checkTime(t, tt, expected[1])

	tt = clock.GetTime()
	checkTime(t, tt, expected[2])

	// Clocked moved back
	tt = clock.GetTime()
	checkTime(t, tt, expected[3])

	tt = clock.UpdateTime(gohlc.NewHLCTime(time.Unix(1, 2), 3))
	checkTime(t, tt, expected[4])

	tt = clock.UpdateTime(gohlc.NewHLCTime(time.Unix(1, 2), 3))
	checkTime(t, tt, expected[5])

	tt = clock.UpdateTime(gohlc.NewHLCTime(time.Unix(3, 0), 1))
	checkTime(t, tt, expected[6])

	tt = clock.UpdateTime(gohlc.NewHLCTime(time.Unix(3, 0), 99))
	checkTime(t, tt, expected[7])

	// Event with greated wall time
	tt = clock.UpdateTime(gohlc.NewHLCTime(time.Unix(4, 4), 100))
	checkTime(t, tt, expected[8])

	tt = clock.UpdateTime(gohlc.NewHLCTime(time.Unix(4, 5), 0))
	checkTime(t, tt, expected[9])

	tt = clock.UpdateTime(gohlc.NewHLCTime(time.Unix(5, 0), 99))
	checkTime(t, tt, expected[10])

	tt = clock.UpdateTime(gohlc.NewHLCTime(time.Unix(5, 0), 50))
	checkTime(t, tt, expected[11])

}
