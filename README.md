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
	TermTimeA int
	TermTimeB int
)

func routineA(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			time.Sleep(time.Second * time.Duration(TermTimeA))
			log.Printf("routine A terminated in %d seconds\n", TermTimeA)
			return
		default:
			log.Println("routine A loop")
		}
		time.Sleep(time.Second)
	}
}

func routineB(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			time.Sleep(time.Second * time.Duration(TermTimeB))
			log.Printf("routine B terminated in %d seconds\n", TermTimeB)
			return
		default:
			log.Println("routine B loop")
		}
		time.Sleep(time.Second)
	}
}

func main() {
	d := daemon.NewDaemon(context.Background())
	d.Run(routineA)
	d.Run(routineB)
	d.WaitSignal()
}
```

### Example 1: Graceful Terminate

set `TermTimeA = 3` and `TermTimeB = 4`

send terminate signal once

**outputs**

```
2021/10/08 22:13:30 routine B loop
2021/10/08 22:13:30 routine A loop
^C2021/10/08 22:13:30 Got signal: interrupt.
2021/10/08 22:13:30 Waiting 5 seconds for graceful terminate. You can pass interrupt signal again to forced exit daemon.
2021/10/08 22:13:34 routine A terminated in 3 seconds
2021/10/08 22:13:35 routine B terminated in 4 seconds
2021/10/08 22:13:35 Daemon exit.

shell returned 1
```

### Example 2: Automatically Forced Terminate

set `TermTimeA = 3` and `TermTimeB = 7`

send terminate signal once

**outputs**

```
2021/10/08 22:17:52 routine B loop
2021/10/08 22:17:52 routine A loop
^C2021/10/08 22:17:53 Got signal: interrupt.
2021/10/08 22:17:53 Waiting 5 seconds for graceful terminate. You can pass interrupt signal again to forced exit daemon.
2021/10/08 22:17:56 routine A terminated in 3 seconds
2021/10/08 22:17:58 Timeout waiting for graceful exit, perform forced exit.

shell returned 1
```

### Example 3: Manually Forced Terminate

set `TermTimeA = 3` and `TermTimeB = 7`

send terminate signal twice

**outputs**

```
2021/10/08 22:25:09 routine A loop
2021/10/08 22:25:09 routine B loop
^C2021/10/08 22:25:09 Got signal: interrupt.
2021/10/08 22:25:09 Waiting 5 seconds for graceful terminate. You can pass interrupt signal again to forced exit daemon.
^C2021/10/08 22:25:09 Manually forced exit.

shell returned 1
```
