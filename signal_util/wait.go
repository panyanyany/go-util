package signal_util

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitForBreak() {
	// wait Ctrl + C
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	<-signals // wait for signal
	go func() {
		<-signals // hard exit on second signal (in case shutdown gets stuck)
		os.Exit(1)
	}()
}
