package core

import (
	"sync"

	"github.com/fatih/color"
	"github.com/zer0ne-hub/z0ne/internal/recon"
)

// Task:Single pipeline step with its
// dependencies and execution logic
type Task struct {
	Name         string
	Dependencies []string
	Execute      func(results map[string]interface{}) error
}

// RunPipeline: Executes tasks in parallel
// while respecting dependencies
func RunPipeline(tasks []Task, maxWorkers int) {
	results := make(map[string]interface{})
	var mu sync.Mutex
	var wg sync.WaitGroup

	taskCompleted := make(map[string]bool)
	taskInProgress := make(map[string]bool)
	taskQueue := make(chan Task)

	// Workers
	for i := 0; i < maxWorkers; i++ {
		go func() {
			for task := range taskQueue {
				err := task.Execute(results)
				mu.Lock()
				if err != nil {
					color.Red("[!] %s failed: %v", task.Name, err)
				} else {
					color.Green("[✓] %s completed", task.Name)
					taskCompleted[task.Name] = true
				}
				mu.Unlock()
				wg.Done()
			}
		}()
	}

	// Scheduler
	for {
		startedAnyTask := false

		for _, t := range tasks {
			mu.Lock()

			if taskCompleted[t.Name] || taskInProgress[t.Name] {
				mu.Unlock()
				continue
			}

			// Check if dependencies are all completed
			depsReady := true
			for _, dep := range t.Dependencies {
				if !taskCompleted[dep] {
					depsReady = false
					break
				}
			}

			if depsReady {
				taskInProgress[t.Name] = true
				wg.Add(1)
				taskQueue <- t
				startedAnyTask = true
			}

			mu.Unlock()
		}

		if !startedAnyTask {
			mu.Lock()
			if len(taskCompleted) == len(tasks) {
				mu.Unlock()
				break
			}
			mu.Unlock()
		}
	}

	wg.Wait()
	close(taskQueue)
}

// buildTasks returns the ordered tasks
// for the given mode
func buildTasks(mode string, target string) []Task {
	switch mode {
	case "scan":
		return []Task{
			{Name: "naabu", Execute: func(r map[string]interface{}) error {
				out, err := recon.RunNaabu(target, "", "")
				if err != nil {
					return err
				}
				return SaveResultToJSON("naabu", out)
			}},
			{Name: "subfinder", Execute: func(r map[string]interface{}) error {
				out, err := recon.RunSubfinder(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("subfinder", out)
			}},

			{Name: "dnsx", Dependencies: []string{"subfinder"}, Execute: func(r map[string]interface{}) error {
				out, err := recon.RunDnsX(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("dnsx", out)
			}},

			{Name: "httpx", Dependencies: []string{"dnsx", "naabu"}, Execute: func(r map[string]interface{}) error {
				out, err := recon.RunHttpX(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("httpx", out)
			}},

			{Name: "katana", Dependencies: []string{"httpx"}, Execute: func(r map[string]interface{}) error {
				out, err := recon.RunKatana(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("katana", out)
			}},
		}

	case "probe":
		return []Task{
			{Name: "naabu", Execute: func(r map[string]interface{}) error {
				out, err := recon.RunNaabu(target, "", "")
				if err != nil {
					return err
				}
				return SaveResultToJSON("naabu", out)
			}},
			{Name: "subfinder", Execute: func(r map[string]interface{}) error {
				out, err := recon.RunSubfinder(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("subfinder", out)
			}},
			{Name: "dnsx", Dependencies: []string{"subfinder"}, Execute: func(r map[string]interface{}) error {
				out, err := recon.RunDnsX(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("dnsx", out)
			}},

			{Name: "httpx", Dependencies: []string{"dnsx", "naabu"}, Execute: func(r map[string]interface{}) error {
				out, err := recon.RunHttpX(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("httpx", out)
			}},
			{Name: "nuclei", Dependencies: []string{"httpx"}, Execute: func(r map[string]interface{}) error {
				out, err := recon.RunNuclei(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("nuclei", out)
			}},

			{Name: "uncover", Dependencies: []string{"httpx"}, Execute: func(r map[string]interface{}) error {
				out, err := recon.RunUncover(target)
				if err != nil {
					return err
				}
				return SaveResultToJSON("uncover", out)
			}},
		}
	}
	return nil
}

// RunRecon executes a broad scanning sequence
func RunRecon(target string) {
	if targetType := detectTargetType(target); targetType == IP || targetType == DOMAIN {
		color.Cyan("🎯 Target: %s", target)
		RunPipeline(buildTasks("scan", target), 3)
	} else {
		color.Red("Unknown target type: %s", target)
	}
}

// RunProbe executes a precise probing sequence
func RunProbe(target string) {
	if targetType := detectTargetType(target); targetType == IP || targetType == DOMAIN {
		color.Cyan("🎯 Target: %s", target)
		RunPipeline(buildTasks("probe", target), 3)
	} else {
		color.Red("Unknown target type: %s", target)
	}
}
