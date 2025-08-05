// Go benchmark with goroutines
// Run with: go run benchmark.go

package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// CPU-intensive calculation: Prime number generation with factorization
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	sqrt := int(math.Sqrt(float64(n)))
	for i := 3; i <= sqrt; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func findPrimesInRange(start, end int) []int {
	var primes []int
	for i := start; i <= end; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

// Additional CPU-intensive work: Matrix multiplication
func matrixMultiply(size int) [][]float64 {
	// Create matrices
	a := make([][]float64, size)
	b := make([][]float64, size)
	result := make([][]float64, size)

	for i := 0; i < size; i++ {
		a[i] = make([]float64, size)
		b[i] = make([]float64, size)
		result[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			a[i][j] = rand.Float64()
			b[i][j] = rand.Float64()
		}
	}

	// Multiply matrices
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return result
}

// Task result structure
type TaskResult struct {
	TaskID      int
	PrimeCount  int
	MatrixCount int
	FibCount    int
	LastPrime   int
}

// Worker result structure
type WorkerResult struct {
	WorkerID      int
	ExecutionTime time.Duration
	Tasks         []TaskResult
}

// Combined CPU-intensive task
func complexCalculation(taskID, iterations int) TaskResult {
	start := taskID * 10000
	end := start + iterations

	// Find primes in range
	primes := findPrimesInRange(start, end)

	// Perform matrix operations
	matrices := make([][][]float64, 3)
	for i := 0; i < 3; i++ {
		matrices[i] = matrixMultiply(50)
	}

	// Additional computational work: Fibonacci
	fibonacci := func(n int) int {
		if n <= 1 {
			return n
		}
		a, b := 0, 1
		for i := 2; i <= n; i++ {
			a, b = b, a+b
		}
		return b
	}

	fibs := make([]int, 100)
	for i := 0; i < 100; i++ {
		fibs[i] = fibonacci(1000 + i)
	}

	lastPrime := -1
	if len(primes) > 0 {
		lastPrime = primes[len(primes)-1]
	}

	return TaskResult{
		TaskID:      taskID,
		PrimeCount:  len(primes),
		MatrixCount: len(matrices),
		FibCount:    len(fibs),
		LastPrime:   lastPrime,
	}
}

// Worker function
func worker(workerID, tasksPerWorker, iterationsPerTask int, resultChan chan<- WorkerResult, wg *sync.WaitGroup) {
	defer wg.Done()

	workerStart := time.Now()
	tasks := make([]TaskResult, tasksPerWorker)

	for i := 0; i < tasksPerWorker; i++ {
		taskID := workerID*tasksPerWorker + i
		tasks[i] = complexCalculation(taskID, iterationsPerTask)
	}

	workerEnd := time.Now()

	resultChan <- WorkerResult{
		WorkerID:      workerID,
		ExecutionTime: workerEnd.Sub(workerStart),
		Tasks:         tasks,
	}
}

func main() {
	numCPUs := runtime.NumCPU()
	numWorkers := numCPUs
	tasksPerWorker := 5
	iterationsPerTask := 2000

	// Set GOMAXPROCS to utilize all CPUs
	runtime.GOMAXPROCS(numCPUs)

	fmt.Printf("Go Goroutines Benchmark\n")
	fmt.Printf("CPUs: %d, Goroutines: %d\n", numCPUs, numWorkers)
	fmt.Printf("Tasks per goroutine: %d, Iterations per task: %d\n", tasksPerWorker, iterationsPerTask)

	startTime := time.Now()

	// Create channels and wait group
	resultChan := make(chan WorkerResult, numWorkers)
	var wg sync.WaitGroup

	// Launch goroutines
	for i := range numWorkers {
		wg.Add(1)
		go worker(i, tasksPerWorker, iterationsPerTask, resultChan, &wg)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var results []WorkerResult
	for result := range resultChan {
		results = append(results, result)
	}

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	// Calculate totals
	totalPrimes := 0
	totalMatrices := 0
	totalFibs := 0

	for _, workerResult := range results {
		for _, task := range workerResult.Tasks {
			totalPrimes += task.PrimeCount
			totalMatrices += task.MatrixCount
			totalFibs += task.FibCount
		}
	}

	fmt.Println("\n=== BENCHMARK RESULTS ===")
	fmt.Printf("Total execution time: %v\n", totalTime)
	fmt.Printf("Total primes found: %d\n", totalPrimes)
	fmt.Printf("Total matrices computed: %d\n", totalMatrices)
	fmt.Printf("Total fibonacci numbers: %d\n", totalFibs)
	totalOps := float64(totalPrimes + totalMatrices + totalFibs)
	opsPerSec := totalOps / totalTime.Seconds()
	fmt.Printf("Operations per second: %.2f\n", opsPerSec)
	fmt.Printf("Parallel efficiency: %d goroutines utilized\n", numWorkers)
}
