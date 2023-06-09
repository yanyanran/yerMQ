package main

import (
	"Y-MQ/yerMQ"
	"context"
	"github.com/judwhite/go-svc"
	"log"
	"os"
	"sync"
	"syscall"
)

type program struct {
	once  sync.Once
	yerMq *yerMQ.YERMQ
}

func main() {
	prg := &program{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatalf("[yerMQ]%s", err)
	}
}

func (p *program) Init(environment svc.Environment) error {
	// 日志输出文件
	file, err := os.OpenFile("sys.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("[FAILED] failed to open error logger file:", err)
	}
	// 自定义日志格式
	log.SetOutput(file)

	opts := yerMQ.NewOptions()
	yermq, err := yerMQ.New(opts)
	if err != nil {
		log.Fatalf("[FAILED] failed to instantiate yermq - %s\n", err)
	}
	p.yerMq = yermq
	return nil
}

func (p *program) Start() error { // in Run
	err := p.yerMq.LoadMetadata()
	if err != nil {
		log.Fatalf("[FAILED] failed to LOAD metadata - %s", err)
	}
	err = p.yerMq.PersistMetadata()
	if err != nil {
		log.Fatalf("[FAILED] failed to PERSIST metadata - %s", err)
	}

	go func() {
		err := p.yerMq.Main() // go Main
		if err != nil {
			p.Stop()
			os.Exit(1)
		}
	}()

	return nil
}

func (p *program) Stop() error {
	p.once.Do(func() {
		p.yerMq.Exit()
	})
	return nil
}

func (p *program) Context() context.Context {
	return p.yerMq.Context()
}
