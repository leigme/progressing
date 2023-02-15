package progressing

import (
	"fmt"
	"strings"
	"time"
)

type ProcessBar interface {
	Start()
	Refresh(i int)
	Stop()
}

func (pb *processBar) Start() {
	go func() {
		for {
			select {
			case <-pb.stopC:
				return
			case i := <-pb.refreshC:
				fmt.Printf("%d ", i)
			default:
				fmt.Printf("%s ", pb.symbol)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (pb *processBar) Refresh(i int) {
	pb.refreshC <- i
}

func (pb *processBar) Stop() {
	defer close(pb.stopC)
	time.Sleep(1 * time.Second)
}

type processBar struct {
	options
}

type options struct {
	symbol   string
	stopC    chan struct{}
	refreshC chan int
}

type Option func(o *options)

func New(ops ...Option) (pb ProcessBar) {
	options := newDefault()
	for _, apply := range ops {
		apply(&options)
	}
	pb = &processBar{
		options: options,
	}
	return
}

func newDefault() (o options) {
	o = options{
		symbol:   ".",
		stopC:    make(chan struct{}),
		refreshC: make(chan int),
	}
	return
}

func WithSymbol(symbol string) Option {
	return func(o *options) {
		if !strings.EqualFold(symbol, "") {
			o.symbol = symbol
		}
	}
}

func WithStopC(stopC chan struct{}) Option {
	return func(o *options) {
		if stopC != nil {
			o.stopC = stopC
		}
	}
}

func WithRefreshC(refreshC chan int) Option {
	return func(o *options) {
		if refreshC != nil {
			o.refreshC = refreshC
		}
	}
}
