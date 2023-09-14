# go-daemon

[![GoDoc](https://godoc.org/github.com/wujie1993/go-daemon?status.svg)](https://godoc.org/github.com/wujie1993/go-daemon)
[![Go Report Card](https://goreportcard.com/badge/github.com/wujie1993/go-daemon)](https://goreportcard.com/report/github.com/wujie1993/go-daemon)

## Installation

```
go get github.com/wujie1993/go-daemon
```

## Usage

```
package main

import (
	"context"
	"log"
	"time"

	"github.com/wujie1993/go-daemon"
)

var (
	TermTimeA = 3
	TermTimeB = 4
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

	d.WaitSignal()
}
```

### Example 1: Graceful Terminate With Timeout

set `TermTimeA = 3` and `TermTimeB = 7`

send terminate signal once when running process

**outputs**

```
2023/09/14 15:01:55 routine A loop
2023/09/14 15:01:55 routine B loop
^C2023/09/14 15:01:56 [warn] Got signal: interrupt.
2023/09/14 15:01:56 [warn] Waiting 5 seconds for graceful terminate. You can pass interrupt signal again to forced exit daemon.
2023/09/14 15:01:56 routine A will stop in 3 seconds
2023/09/14 15:01:56 routine B will stop in 7 seconds
2023/09/14 15:01:59 routine A stopped
2023/09/14 15:02:01 [error] Timeout waiting for graceful exit, perform forced exit.

shell returned 1
```

### Example 2: Manually Forced Terminate

set `TermTimeA = 3` and `TermTimeB = 7`

send terminate signal twice when running process

**outputs**

```
2023/09/14 15:03:02 routine A loop
2023/09/14 15:03:02 routine B loop
^C2023/09/14 15:03:03 [warn] Got signal: interrupt.
2023/09/14 15:03:03 [warn] Waiting 5 seconds for graceful terminate. You can pass interrupt signal again to forced exit daemon.
2023/09/14 15:03:03 routine A will stop in 3 seconds
2023/09/14 15:03:03 routine B will stop in 7 seconds
^C2023/09/14 15:03:04 [warn] Manually forced exit.

shell returned 1
```
