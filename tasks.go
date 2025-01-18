package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
	MaxRetries   uint
	MaxProcesses uint8
}

type Worker struct {
	Ticker   *Ticker
	Context  *Context
	Logger   *Logger
	Settings Settings
}

func NewTicker(i int) *Ticker {
	var interval time.Duration = time.Duration(i)
	return &Ticker{interval: interval, done: make(chan bool), t: time.NewTicker(interval * time.Second)}
}

func NewWorker(retries uint, processes uint8, l *Logger) Worker {
	if l == nil {
		return Worker{
			Logger:   logger,
			Settings: Settings{retries, processes},
		}
	}

	return Worker{
		Logger:   l,
		Settings: Settings{retries, processes},
	}
}

// function DoWork dispatches/spawns goroutines to execute tasks
func (w *Worker) DoWork() error {
	return nil
}

func (w *Worker) StartHeartbeat() {
	go func() {
		for {
			select {
			case <-w.Ticker.done:
				return
			case t := <-w.Ticker.t.C:
				msg := fmt.Sprintf("heartbeat at %v", t.Format(time.DateTime))
				w.Logger.Info(msg)
			}
		}
	}()
}

func (w *Worker) ListenForSignals() {
	sigChannel := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChannel
		done <- true
		w.Logger.Info("received signal " + sig.String())
	}()

	<-done
	w.Ticker.t.Stop()
}

func (w *Worker) CreateContext() {
	w.Context = &Context{}
	w.Context.ctx, w.Context.cancel = context.WithCancel(context.Background())
}

func (w *Worker) CreateTicker(sec int) {
	w.Ticker = NewTicker(sec)
}

func (w *Worker) Run() {
	w.Logger.Info("starting worker...")
	w.CreateContext()
	w.CreateTicker(5)
	w.StartHeartbeat()
	w.ListenForSignals()
	w.DoWork()
}
