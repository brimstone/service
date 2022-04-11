package service

import (
	"context"
	"time"

	"github.com/brimstone/logger"
)

type Runner interface {
	Run(ctx context.Context) error
}

type Stater interface {
	Stat() interface{}
}

type manager struct {
	runables map[string]*runable
	stats    map[string]*statable
}

type runable struct {
	run     Runner
	running bool
	ctx     context.Context
	name    string
	cancel  func()
}

type statable struct {
	stat   Stater
	ctx    context.Context
	name   string
	cancel func()
}

var Manager manager

func (m *manager) Add(name string, svc Runner) {
	if m.runables == nil {
		m.runables = make(map[string]*runable)
	}
	r := &runable{
		name: name,
		run:  svc,
	}
	m.runables[name] = r
}

func (m *manager) RunAll(ctx context.Context) {
	log := logger.New()
	for n, r := range m.runables {
		go func(name string, r *runable) {
			r.ctx, r.cancel = context.WithCancel(ctx)
			log.Info("Starting service",
				log.Field("service", name),
			)
			r.running = true
			defer func() { r.running = false }()
			err := r.run.Run(r.ctx)
			if err != nil {
				log.Error("Service failed",
					log.Field("service", name),
					log.Field("err", err),
				)
			}
		}(n, r)
	}
}

func (m *manager) StopAll() {
	log := logger.New()
	running := false
	for {
		for n, r := range m.runables {
			if r.running {
				log.Info("Service still running",
					log.Field("name", n),
				)
				running = true
				r.cancel()
				break
			}
		}
		if !running {
			break
		}
		time.Sleep(time.Second)
		running = false
	}
}

func (m *manager) Stats() map[string]interface{} {
	s := make(map[string]interface{})
	for n, r := range m.runables {
		if st, ok := r.run.(Stater); ok {
			s[n] = st.Stat()
		} else {
			s[n] = struct {
				Running bool `json:"running"`
			}{
				Running: r.running,
			}
		}
	}
	return s
}
