package main

import (
	"context"
	"log"
	"time"

	"github.com/wujie1993/go-daemon"
)

var (
	TermTimeA = 3
	TermTimeB = 7
)

func routineA(ctx context.Context) {
	defer log.Println("routine A stopped")
	for {
		select {
		case <-ctx.Done():
			log.Printf("routine A will stop in %d seconds\n", TermTimeA)
			time.Sleep(time.Second * time.Duration(TermTimeA))
			return
		default:
			log.Println("routine A loop")
		}
		time.Sleep(time.Second)
	}
}

func routineB(ctx context.Context) {
	defer log.Println("routine B stopped")
	for {
		select {
		case <-ctx.Done():
			log.Printf("routine B will stop in %d seconds\n", TermTimeB)
			time.Sleep(time.Second * time.Duration(TermTimeB))
			return
		default:
			log.Println("routine B loop")
		}
		time.Sleep(time.Second)
	}
}

func main() {
	d := daemon.NewDaemon(daemon.DaemonConfig{
		Ctx:                     context.Background(),
		GracefulExitWaitSeconds: 5,
	})
	d.Run("routine A", routineA)
	d.Run("routine B", routineB)

	//go func() {
	//	// kill routine A after 1s
	//	time.Sleep(time.Second)
	//	d.Kill("routine A")
	//}()

	//go func() {
	//	// kill routine B after 2s
	//	time.Sleep(time.Second)
	//	d.Kill("routine B")
	//}()

	// wait until all routines stop
	d.WaitSignal()
}
