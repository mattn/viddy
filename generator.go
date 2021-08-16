package main

import "time"

func ClockSnapshot(name string, args []string, interval time.Duration) <-chan *Snapshot {
	c := make(chan *Snapshot)

	go func() {
		var s *Snapshot
		t := time.Tick(interval)
		begin := time.Now().UnixNano()

		for {
			select {
			case now := <-t:
				finish := make(chan struct{})
				id := (now.UnixNano() - begin) / int64(time.Millisecond)
				s = NewSnapshot(id, name, args, s, finish)
				c <- s
			}
		}
	}()

	return c
}

func PreciseSnapshot(name string, args []string, interval time.Duration) <-chan *Snapshot {
	c := make(chan *Snapshot)

	go func() {
		var s *Snapshot
		begin := time.Now().UnixNano()

		for {
			finish := make(chan struct{})
			start := time.Now()
			id := (start.UnixNano() - begin) / int64(time.Millisecond)
			ns := NewSnapshot(id, name, args, s, finish)
			s = ns
			c <- ns
			<-finish
			pTime := time.Since(start)

			if pTime > interval {
				continue
			} else {
				time.Sleep(interval - pTime)
			}
		}
	}()

	return c
}

func SequentialSnapshot(name string, args []string, interval time.Duration) <-chan *Snapshot {
	c := make(chan *Snapshot)

	go func() {
		var s *Snapshot
		begin := time.Now().UnixNano()

		for {
			finish := make(chan struct{})
			id := (time.Now().UnixNano() - begin) / int64(time.Millisecond)
			s = NewSnapshot(id, name, args, s, finish)
			c <- s
			<-finish

			time.Sleep(interval)
		}
	}()

	return c
}
