package ui

import (
	"fmt"
	"sync"
	"time"
)

var spinnerFrames = []string{"|", "/", "-", "\\"}

type Spinner struct {
	msg  string
	done chan struct{}
	wg   sync.WaitGroup
}

func NewSpinner(msg string) *Spinner {
	s := &Spinner{msg: msg, done: make(chan struct{})}
	s.wg.Add(1)
	go s.run()
	return s
}

func (s *Spinner) run() {
	defer s.wg.Done()
	if !colorsEnabled {
		fmt.Printf("%s... ", s.msg)
		<-s.done
		return
	}
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-s.done:
			fmt.Print("\r\033[K")
			return
		case <-ticker.C:
			fmt.Printf("\r%s %s", Cyan(spinnerFrames[i%len(spinnerFrames)]), s.msg)
			i++
		}
	}
}

func (s *Spinner) Stop(msg string) {
	close(s.done)
	s.wg.Wait()
	Success(msg)
}

func (s *Spinner) Fail(msg string) {
	close(s.done)
	s.wg.Wait()
	Fail(msg)
}
