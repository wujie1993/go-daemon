package daemon

import (
	"context"
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
	cancels                 map[string]context.CancelFunc
	cancelsMutex            sync.RWMutex
	log                     Logger
}

type DaemonConfig struct {
	Ctx                     context.Context
	GracefulExitWaitSeconds time.Duration
	Log                     Logger
}

func (d *Daemon) Run(name string, run DaemonRunFunc) {
	ctx, cancel := context.WithCancel(d.ctx)

	d.cancelsMutex.Lock()
	d.cancels[name] = cancel
	d.cancelsMutex.Unlock()

	d.wg.Add(1)

	go func() {
		defer func() {
			d.cancelsMutex.Lock()
			delete(d.cancels, name)
			d.cancelsMutex.Unlock()

			d.wg.Done()
		}()

		run(ctx)
	}()
}

func (d *Daemon) Kill(name string) {
	d.cancelsMutex.RLock()
	cancel, ok := d.cancels[name]
	d.cancelsMutex.RUnlock()
	if !ok {
		return
	}
	cancel()
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
			d.log.Warn("There is no running routines.")
			return
		case killSignal := <-d.interrupt:
			d.log.Warnf("Got signal: %s.\n", killSignal)
			if killSignal == os.Interrupt {
				d.log.Warnf("Waiting %d seconds for graceful terminate. You can pass interrupt signal again to forced exit daemon.\n", d.gracefulExitWaitSeconds)
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
			d.log.Info("Daemon exit.")
			return
		case killSignal := <-d.interrupt:
			if killSignal == os.Interrupt {
				d.log.Warn("Manually forced exit.")
				return
			}
		case <-time.After(time.Duration(d.gracefulExitWaitSeconds * time.Second)):
			d.log.Error("Timeout waiting for graceful exit, perform forced exit.")
			return
		}
	}
}

func NewDaemon(config DaemonConfig) *Daemon {
	if config.Ctx == nil {
		config.Ctx = context.Background()
	}
	if config.GracefulExitWaitSeconds <= 0 {
		config.GracefulExitWaitSeconds = 5
	}
	if config.Log == nil {
		config.Log = DefaultLog{}
	}
	ctx, cancel := context.WithCancel(config.Ctx)
	return &Daemon{
		ctx:                     ctx,
		cancel:                  cancel,
		gracefulExitWaitSeconds: config.GracefulExitWaitSeconds,
		cancels:                 make(map[string]context.CancelFunc),
		log:                     config.Log,
	}
}
