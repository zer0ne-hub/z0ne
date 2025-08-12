// Package core: Handles core engine logic
package core

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
	"github.com/zer0ne-hub/z0ne/internal/recon"
)

// Task represents a single pipeline step with dependencies and execution logic
type Task struct {
	Name         string
	Dependencies []string
	Execute      func(results map[string]interface{}) error
}

// RunPipeline executes tasks in parallel while respecting dependencies.
func RunPipeline(tasks []Task, maxWorkers int) {
	results := make(map[string]interface{})
	var mu sync.Mutex
	var wg sync.WaitGroup

	taskCompleted := make(map[string]bool)
	taskInProgress := make(map[string]bool)
	taskQueue := make(chan Task)

	startWorkers(maxWorkers, taskQueue, &wg, &mu, taskCompleted, results)

	scheduleTasks(tasks, taskQueue, &wg, &mu, taskCompleted, taskInProgress)

	wg.Wait()
	close(taskQueue)
}

// startWorkers launches goroutines to process tasks
func startWorkers(maxWorkers int, taskQueue chan Task, wg *sync.WaitGroup, mu *sync.Mutex,
	taskCompleted map[string]bool, results map[string]interface{}) {

	for i := 0; i < maxWorkers; i++ {
		go func() {
			for task := range taskQueue {
				err := task.Execute(results)

				mu.Lock()
				if err != nil {
					color.Red("[!] %s failed: %v \n\n", task.Name, err)
				} else {
					color.Green("[✓] %s completed\n\n", task.Name)
					taskCompleted[task.Name] = true
				}
				mu.Unlock()

				wg.Done()
			}
		}()
	}
}

// scheduleTasks manages dependencies and sends ready tasks to the queue
func scheduleTasks(tasks []Task, taskQueue chan Task, wg *sync.WaitGroup, mu *sync.Mutex,
	taskCompleted, taskInProgress map[string]bool) {

	for {
		startedAnyTask := false

		for _, t := range tasks {
			mu.Lock()

			if taskCompleted[t.Name] || taskInProgress[t.Name] {
				mu.Unlock()
				continue
			}

			if dependenciesMet(t.Dependencies, taskCompleted) {
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
}

// dependenciesMet checks if all dependencies are completed
func dependenciesMet(deps []string, completed map[string]bool) bool {
	for _, dep := range deps {
		if !completed[dep] {
			return false
		}
	}
	return true
}

// buildTasks returns the ordered tasks for a given mode and target
func buildTasks(mode string, target string) []Task {
	switch mode {
	case "scan":
		return scanTasks(target)
	case "probe":
		return probeTasks(target)
	}
	return nil
}

// scanTasks defines the scanning task sequence
func scanTasks(target string) []Task {
	return []Task{
		{Name: "naabu", Execute: func(r map[string]interface{}) error {
			result, err := recon.RunNaabu(target, "", "")
			return runAndSave("naabu", result, err)
		}},
		{Name: "subfinder", Execute: func(r map[string]interface{}) error {
			if detectTargetType(target) != DOMAIN {
				color.Yellow("[!] Skipping subfinder: target is not a domain\n")
				return nil
			}
			result, err := recon.RunSubfinder(target)
			return runAndSave("subfinder", result, err)
		}},
		{Name: "dnsx", Dependencies: []string{"subfinder"}, Execute: func(r map[string]interface{}) error {
			if detectTargetType(target) != DOMAIN {
				color.Yellow("[!] Skipping dnsx: target is not a domain\n")
				return nil
			}
			result, err := recon.RunDnsX(target)
			return runAndSave("dnsx", result, err)
		}},
		{Name: "httpx", Dependencies: []string{"dnsx", "naabu"}, Execute: func(r map[string]interface{}) error {
			result, err := recon.RunHttpX(target)
			return runAndSave("httpx", result, err)
		}},
		{Name: "katana", Dependencies: []string{"httpx"}, Execute: func(r map[string]interface{}) error {
			result, err := recon.RunKatana(target)
			return runAndSave("katana", result, err)
		}},
	}
}

// probeTasks defines the probing task sequence
func probeTasks(target string) []Task {
	return []Task{
		{Name: "naabu", Execute: func(r map[string]interface{}) error {
			result, err := recon.RunNaabu(target, "", "")
			return runAndSave("naabu", result, err)
		}},
		{Name: "subfinder", Execute: func(r map[string]interface{}) error {
			result, err := recon.RunSubfinder(target)
			return runAndSave("subfinder", result, err)
		}},
		{Name: "dnsx", Dependencies: []string{"subfinder"}, Execute: func(r map[string]interface{}) error {
			result, err := recon.RunDnsX(target)
			return runAndSave("dnsx", result, err)
		}},
		{Name: "httpx", Dependencies: []string{"dnsx", "naabu"}, Execute: func(r map[string]interface{}) error {
			result, err := recon.RunHttpX(target)
			return runAndSave("httpx", result, err)
		}},
		{Name: "nuclei", Dependencies: []string{"httpx"}, Execute: func(r map[string]interface{}) error {
			result, err := recon.RunNuclei(target)
			return runAndSave("nuclei", result, err)
		}},
	}
}

// runAndSave wraps tool execution with result saving
func runAndSave(toolName string, out interface{}, err error) error {
	if err != nil {
		return err
	}
	return SaveResultToJSON(toolName, out)
}

// RunRecon executes a broad scanning task sequence
func RunRecon(target string) error {
	if targetType := detectTargetType(target); targetType == IP || targetType == DOMAIN {
		RunPipeline(buildTasks("scan", target), 3)
	} else {
		color.Red("Unknown target type: %s", target)
	}
	return nil
}

// RunProbe executes a precise probing task sequence
func RunProbe(target string, keys ProbeKeys) error {
	if targetType := detectTargetType(target); targetType == IP || targetType == DOMAIN {
		tasks := buildTasks("probe", target)
		if keys.ShodanKey != "" {
			fmt.Println("Shodan API key found, Running Uncover...")
			tasks = append(tasks, Task{
				Name:         "uncover",
				Dependencies: []string{"httpx"},
				Execute: func(r map[string]interface{}) error {
					result, err := recon.RunUncover(target, keys.ShodanKey)
					return runAndSave("uncover", result, err)
				},
			})
		}
		RunPipeline(tasks, 3)
	} else {
		color.Red("Unknown target type: %s", target)
	}
	return nil
}
