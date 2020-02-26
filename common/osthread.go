package common

import "runtime"

func init() {
	runtime.LockOSThread()
}

// queue of work to run in main thread.
var Mainfunc chan func()

func MakeMainThread() {
	Mainfunc = make(chan func(), 1)
}

func CloseMainThread() {
	close(Mainfunc)
}

// do runs f on the main thread.
func Do(f func()) {
	done := make(chan bool, 1)
	Mainfunc <- func() {
		f()
		done <- true
	}
	<-done
}
