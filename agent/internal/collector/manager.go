package collector

import (
	"context"
	"log"
	"sync"
	"time"
)

// Runner is the common interface implemented by all agent data collectors.
type Runner interface {
	Name() string
	Run(ctx context.Context)
}

// Manager orchestrates and runs multiple registered collectors concurrently.
type Manager struct {
	collectors []Runner
	interval   time.Duration
}

func NewManager(interval time.Duration) *Manager {
	return &Manager{
		interval: interval,
	}
}

func (m *Manager) Register(c Runner) {
	m.collectors = append(m.collectors, c)
}

func (m *Manager) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for _, c := range m.collectors {
		wg.Add(1)

		go func(col Runner) {
			defer wg.Done()

			ticker := time.NewTicker(m.interval)
			defer ticker.Stop()

			// Initial collection on start
			log.Printf("[%s] collector initialized", col.Name())
			col.Run(ctx)

			for {
				select {
				case <-ctx.Done():
					log.Printf("[%s] collector stopping", col.Name())
					return

				case <-ticker.C:
					log.Printf("[%s] collecting...", col.Name())
					col.Run(ctx)
				}
			}
		}(c)
	}

	wg.Wait()
}
