package daemon

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type DaemonRunFunc func(ctx context.Context)

type Daemon struct {
	ctx                     context.Context
	cancel                  context.CancelFunc
	wg                      sync.WaitGroup
	gracefulExitWaitSeconds time.Duration
	interrupt               chan os.Signal
}

func (d *Daemon) Run(run DaemonRunFunc) {
	go func() {
		d.wg.Add(1)
		defer d.wg.Done()

		run(d.ctx)
	}()
}

func (d *Daemon) WaitSignal() {
	d.interrupt = make(chan os.Signal, 1)
	signal.Notify(d.interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		defer d.cancel()
		d.wg.Wait()
	}()

	defer d.waitExit()

	for {
		select {
		case <-d.ctx.Done():
			log.Println("There is no running routines.")
			return
		case killSignal := <-d.interrupt:
			log.Printf("Got signal: %s.\n", killSignal)
			if killSignal == os.Interrupt {
				log.Printf("Waiting %d seconds for graceful terminate. You can pass interrupt signal again to forced exit daemon.\n", d.gracefulExitWaitSeconds)
				d.cancel()
				return
			}
		}
	}
}

func (d *Daemon) waitExit() {
	c := make(chan struct{})
	go func() {
		defer close(c)
		d.wg.Wait()
	}()
	for {
		select {
		case <-c:
			log.Println("Daemon has graceful exit.")
			return
		case killSignal := <-d.interrupt:
			if killSignal == os.Interrupt {
				log.Println("Manually forced exit.")
				return
			}
		case <-time.After(time.Duration(d.gracefulExitWaitSeconds * time.Second)):
			log.Println("Timeout waiting for graceful exit, perform forced exit.")
			return
		}
	}
}

func NewDaemon(ctx context.Context) *Daemon {
	ctx, cancel := context.WithCancel(context.Background())
	return &Daemon{
		ctx:                     ctx,
		cancel:                  cancel,
		gracefulExitWaitSeconds: 5,
	}
}
