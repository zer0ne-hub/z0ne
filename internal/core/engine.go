package core

import (
	"sync"

	"github.com/fatih/color"
	"github.com/zer0ne-hub/z0ne/internal/recon"
)

// Pipeline step
type Task struct {
	Name string
	Deps []string
	Run  func(results map[string]interface{}) error
}

// Parallel running pipeline
func RunPipeline(tasks []Task, maxWorkers int) {
	results := make(map[string]interface{})
	var mu sync.Mutex
	var wg sync.WaitGroup
	done := make(map[string]bool)
	inProgress := make(map[string]bool)
	taskChan := make(chan Task)

	// Workers
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for task := range taskChan {
				err := task.Run(results)
				if err != nil {
					color.Red("[!] %s failed: %v", task.Name, err)
				} else {
					color.Green("[✓] %s completed", task.Name)
					mu.Lock()
					done[task.Name] = true
					mu.Unlock()
				}
				wg.Done()
			}
		}()
	}

	// Scheduler
	for {
		started := false
		for _, t := range tasks {
			mu.Lock()
			if done[t.Name] || inProgress[t.Name] {
				mu.Unlock()
				continue
			}
			depsReady := true
			for _, dep := range t.Deps {
				if !done[dep] {
					depsReady = false
					break
				}
			}
			if depsReady {
				inProgress[t.Name] = true
				wg.Add(1)
				taskChan <- t
				started = true
			}
			mu.Unlock()
		}
		if !started {
			mu.Lock()
			if len(done) == len(tasks) {
				mu.Unlock()
				break
			}
			mu.Unlock()
		}
	}

	wg.Wait()
	close(taskChan)
}

func RunRecon(target string) {
	TargetType := detectTargetType(target)
	switch TargetType {
	case IP, DOMAIN:
		color.Cyan("🎯 Target: %s", target)
		tasks := []Task{
			{
				Name: "naabu",
				Deps: []string{},
				Run: func(results map[string]interface{}) error {
					return recon.RunNaabu(target, "", "")
				},
			},
			{
				Name: "subfinder",
				Deps: []string{},
				Run: func(results map[string]interface{}) error {
					return recon.RunSubfinder(target)
				},
			},
			{
				Name: "dnsx",
				Deps: []string{"subfinder"},
				Run: func(results map[string]interface{}) error {
					return recon.RunDnsX(target)
				},
			},
			{
				Name: "httpx",
				Deps: []string{"dnsx", "naabu"},
				Run: func(results map[string]interface{}) error {
					return recon.RunHttpX(target)
				},
			},
			{
				Name: "katana",
				Deps: []string{"httpx"},
				Run: func(results map[string]interface{}) error {
					return recon.RunKatana(target)
				},
			},
			{
				Name: "nuclei",
				Deps: []string{"httpx"},
				Run: func(results map[string]interface{}) error {
					return recon.RunNuclei(target)
				},
			},
		}
		RunPipeline(tasks, 3)
	default:
		color.Red("Unknown target type: %s", target)
	}
}

func RunProbe(target string) {
	TargetType := detectTargetType(target)
	switch TargetType {
	case IP, DOMAIN:
		color.Cyan("🧿 Target: %s", target)

		tasks := []Task{
			{
				Name: "naabu",
				Deps: []string{},
				Run: func(results map[string]interface{}) error {
					return recon.RunNaabu(target, "", "")
				},
			},
			{
				Name: "subfinder",
				Deps: []string{},
				Run: func(results map[string]interface{}) error {
					return recon.RunSubfinder(target)
				},
			},
			{
				Name: "dnsx",
				Deps: []string{"subfinder"},
				Run: func(results map[string]interface{}) error {
					return recon.RunDnsX(target)
				},
			},
			{
				Name: "httpx",
				Deps: []string{"dnsx", "naabu"},
				Run: func(results map[string]interface{}) error {
					return recon.RunHttpX(target)
				},
			},
			{
				Name: "katana",
				Deps: []string{"httpx"},
				Run: func(results map[string]interface{}) error {
					return recon.RunKatana(target)
				},
			},
			{
				Name: "nuclei",
				Deps: []string{"httpx"},
				Run: func(results map[string]interface{}) error {
					return recon.RunNuclei(target)
				},
			},
		}
		RunPipeline(tasks, 3)
	default:
		color.Red("Unknown target type: %s", target)
	}
}
