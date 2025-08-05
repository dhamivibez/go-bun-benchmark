const { Worker, isMainThread, parentPort, workerData } = require('worker_threads');
const os = require('os');

// CPU-intensive calculation: Prime number generation with factorization
function isPrime(n) {
    if (n < 2) return false;
    if (n === 2) return true;
    if (n % 2 === 0) return false;
    
    const sqrt = Math.sqrt(n);
    for (let i = 3; i <= sqrt; i += 2) {
        if (n % i === 0) return false;
    }
    return true;
}

function findPrimesInRange(start, end) {
    const primes = [];
    for (let i = start; i <= end; i++) {
        if (isPrime(i)) {
            primes.push(i);
        }
    }
    return primes;
}

// Additional CPU-intensive work: Matrix multiplication
function matrixMultiply(size) {
    const a = Array(size).fill().map(() => Array(size).fill().map(() => Math.random()));
    const b = Array(size).fill().map(() => Array(size).fill().map(() => Math.random()));
    const result = Array(size).fill().map(() => Array(size).fill(0));
    
    for (let i = 0; i < size; i++) {
        for (let j = 0; j < size; j++) {
            for (let k = 0; k < size; k++) {
                result[i][j] += a[i][k] * b[k][j];
            }
        }
    }
    
    return result;
}

// Combined CPU-intensive task
function complexCalculation(taskId, iterations) {
    const start = taskId * 10000;
    const end = start + iterations;
    
    // Find primes in range
    const primes = findPrimesInRange(start, end);
    
    // Perform matrix operations
    const matrices = [];
    for (let i = 0; i < 3; i++) {
        matrices.push(matrixMultiply(50));
    }
    
    // Additional computational work: Fibonacci with memoization stress
    function fibonacci(n) {
        if (n <= 1) return n;
        let a = 0, b = 1;
        for (let i = 2; i <= n; i++) {
            [a, b] = [b, a + b];
        }
        return b;
    }
    
    const fibs = [];
    for (let i = 0; i < 100; i++) {
        fibs.push(fibonacci(1000 + i));
    }
    
    return {
        taskId,
        primeCount: primes.length,
        matrixCount: matrices.length,
        fibCount: fibs.length,
        lastPrime: primes[primes.length - 1] || -1
    };
}

if (isMainThread) {
    // Main thread - benchmark orchestrator
    async function runBenchmark() {
        const numCPUs = os.cpus().length;
        const numWorkers = numCPUs;
        const tasksPerWorker = 5;
        const iterationsPerTask = 2000;
        
        console.log(`JavaScript Worker Threads Benchmark`);
        console.log(`CPUs: ${numCPUs}, Workers: ${numWorkers}`);
        console.log(`Tasks per worker: ${tasksPerWorker}, Iterations per task: ${iterationsPerTask}`);
        console.log('Starting benchmark...\n');
        
        const startTime = Date.now();
        
        // Create workers
        const workers = [];
        const promises = [];
        
        for (let i = 0; i < numWorkers; i++) {
            const worker = new Worker(__filename, {
                workerData: { workerId: i, tasksPerWorker, iterationsPerTask }
            });
            
            workers.push(worker);
            promises.push(new Promise((resolve, reject) => {
                worker.on('message', resolve);
                worker.on('error', reject);
            }));
        }
        
        // Wait for all workers to complete
        const results = await Promise.all(promises);
        
        // Cleanup workers
        workers.forEach(worker => worker.terminate());
        
        const endTime = Date.now();
        const totalTime = endTime - startTime;
        
        // Calculate totals
        let totalPrimes = 0;
        let totalMatrices = 0;
        let totalFibs = 0;
        
        results.forEach(workerResult => {
            workerResult.tasks.forEach(task => {
                totalPrimes += task.primeCount;
                totalMatrices += task.matrixCount;
                totalFibs += task.fibCount;
            });
        });
        
        console.log('\n=== BENCHMARK RESULTS ===');
        console.log(`Total execution time: ${totalTime}ms`);
        console.log(`Total primes found: ${totalPrimes}`);
        console.log(`Total matrices computed: ${totalMatrices}`);
        console.log(`Total fibonacci numbers: ${totalFibs}`);
        console.log(`Operations per second: ${((totalPrimes + totalMatrices + totalFibs) / totalTime * 1000).toFixed(2)}`);
        console.log(`Parallel efficiency: ${numWorkers} workers utilized`);
    }
    
    runBenchmark().catch(console.error);
    
} else {
    // Worker thread
    const { workerId, tasksPerWorker, iterationsPerTask } = workerData;
    
    const workerStart = Date.now();
    const tasks = [];
    
    for (let i = 0; i < tasksPerWorker; i++) {
        const taskId = workerId * tasksPerWorker + i;
        const result = complexCalculation(taskId, iterationsPerTask);
        tasks.push(result);
    }
    
    const workerEnd = Date.now();
    
    parentPort.postMessage({
        workerId,
        executionTime: workerEnd - workerStart,
        tasks
    });
}