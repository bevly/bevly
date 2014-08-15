package throttle

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Throttle struct {
	Name           string
	MinMillis      int64
	MaxMillis      int64
	LastInvocation time.Time
	Mutex          sync.Mutex
}

func Default(name string) *Throttle {
	return New(name, 4500, 18000)
}

func New(name string, minMillis, maxMillis int64) *Throttle {
	return &Throttle{Name: name, MinMillis: minMillis, MaxMillis: maxMillis}
}

func (t *Throttle) Throttle(action func()) {
	t.DelayInvocation()
	action()
}

func (t *Throttle) DelayInvocation() {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	if !t.LastInvocation.IsZero() {
		sleepDur := t.RandomSleepDuration(time.Now().Sub(t.LastInvocation))
		if sleepDur > 0 {
			log.Printf("Throttle(%s): Sleeping %.2fs", t.Name,
				float64(sleepDur)/1e9)
			time.Sleep(sleepDur)
		}
	}
	t.LastInvocation = time.Now()
}

func clampedDuration(millis int64, base time.Duration) time.Duration {
	remainingTime := millis*1000000 - base.Nanoseconds()
	if remainingTime < 0 {
		remainingTime = 0
	}
	return time.Duration(remainingTime)
}

func (t *Throttle) RandomSleepDuration(durSinceLastInvoke time.Duration) time.Duration {
	minDur := clampedDuration(t.MinMillis, durSinceLastInvoke)
	maxDur := clampedDuration(t.MaxMillis, durSinceLastInvoke)
	if maxDur == 0 {
		return 0
	}
	if minDur == maxDur {
		return minDur
	}
	return time.Duration(
		rand.Int63n(maxDur.Nanoseconds()-minDur.Nanoseconds()) +
			minDur.Nanoseconds())
}
