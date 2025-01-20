package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Ticker struct {
	interval time.Duration
	done     chan bool
	t        *time.Ticker
}

type Context struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type Settings struct {
	MaxRetries   int
	MaxProcesses int
}

type Worker struct {
	Ticker   *Ticker
	Context  *Context
	Logger   *Logger
	Settings Settings
}

func NewTicker(i int) *Ticker {
	interval := time.Duration(i)
	return &Ticker{
		done:     make(chan bool),
		interval: interval,
		t:        time.NewTicker(interval * time.Second),
	}
}

func NewWorker(retries int, processes int, hr int) Worker {
	c := Context{}
	c.ctx, c.cancel = context.WithCancel(context.Background())

	return Worker{
		Logger:   logger,
		Settings: Settings{retries, processes},
		Context:  &c,
		Ticker:   NewTicker(hr),
	}
}

// function DoWork dispatches/spawns goroutines to execute tasks
func (w *Worker) DoWork() error {
	return nil
}

func (w *Worker) StartListener() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)
	messenger := make(chan string)

	go func() {
		for {
			select {
			case <-w.Ticker.done:
				return
			case t := <-w.Ticker.t.C:
				messenger <- fmt.Sprintf("heartbeat at %v", t.Format(time.DateTime))
			}
		}
	}()

	go func() {
		for {
			w.Logger.Info(<-messenger)
		}
	}()

	go func() {
		sig := <-sigChannel
		w.Logger.Info("received signal: " + sig.String())
		w.Ticker.done <- true
	}()

	<-w.Ticker.done
	w.Ticker.t.Stop()
}

// ParseArgs is a part of the [Commander] interface implementation
func ParseWorkerArgs(args []string) map[string]int {
	parsed := make(map[string]int, 3)
	parsed["heartRate"] = 2
	parsed["retries"] = 3
	parsed["processes"] = 1

	for i, arg := range args {
		var val string

		if i != len(args)-1 {
			val = args[i+1]
		} else {
			break
		}

		value, err := strconv.Atoi(val)
		if err != nil {
			logger.Error(fmt.Printf("unable to parse %v value %v", arg, err.Error()))
		}

		switch arg {
		case "--heartrate", "--hr", "-heartrate", "-hr":
			parsed["heartRate"] = value
		case "--retries", "--r", "-retries", "-r":
			parsed["retries"] = value
		case "--processes", "--p", "-processes", "-p":
			parsed["processes"] = value
		}
	}

	return parsed
}

// func Run "turns on the bot," i.e. starts the worker and parses the command-line
// argument slice from [ParseArgs]
func Run(args []string) error {
	parsed := ParseWorkerArgs(args)
	w := NewWorker(parsed["retries"], parsed["processes"], parsed["heartRate"])

	w.StartListener()

	if err := w.DoWork(); err != nil {
		return err
	} else {
		return nil
	}
}
