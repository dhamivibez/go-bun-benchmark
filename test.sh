#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
RUNS=5
WARMUP_RUNS=2

echo -e "${BLUE}=== Go vs JavaScript Benchmark ===${NC}\n"

# Check if files exist
if [ ! -f "main.go" ]; then
    echo -e "${RED}Error: main.go not found${NC}"
    exit 1
fi

if [ ! -f "script.js" ]; then
    echo -e "${RED}Error: script.js not found${NC}"
    exit 1
fi

# Build Go program
echo -e "${YELLOW}Building Go program...${NC}"
go build -ldflags="-s -w" -o benchmark main.go

if [ $? -ne 0 ]; then
    echo -e "${RED}Go build failed${NC}"
    exit 1
fi

echo -e "${GREEN}Build successful!${NC}\n"

# Function to run benchmark
run_benchmark() {
    local cmd="$1"
    local name="$2"
    local runs="$3"
    
    echo -e "${BLUE}$name Benchmark (${runs} runs):${NC}"
    
    local total_ms=0
    
    for ((i=1; i<=runs; i++)); do
        echo -n "  Run $i: "
        
        # Use milliseconds for better precision
        local start_time=$(date +%s%3N)
        eval "$cmd" >/dev/null 2>&1
        local end_time=$(date +%s%3N)
        
        local duration_ms=$((end_time - start_time))
        total_ms=$((total_ms + duration_ms))
        
        # Convert to seconds for display
        local duration_sec=$((duration_ms / 1000))
        local duration_frac=$((duration_ms % 1000))
        printf "%d.%03ds\n" $duration_sec $duration_frac
    done
    
    local avg_ms=$((total_ms / runs))
    local avg_sec=$((avg_ms / 1000))
    local avg_frac=$((avg_ms % 1000))
    printf "  ${GREEN}Average: %d.%03ds${NC}\n\n" $avg_sec $avg_frac
    
    # Return average time in ms
    return $avg_ms
}

# Warmup runs (don't count these)
echo -e "${YELLOW}Warming up...${NC}"
for ((i=1; i<=WARMUP_RUNS; i++)); do
    ./benchmark >/dev/null 2>&1
    bun script.js >/dev/null 2>&1
done
echo -e "${GREEN}Warmup complete!${NC}\n"

# Run benchmarks and capture results differently
echo -e "${BLUE}Go Benchmark (${RUNS} runs):${NC}"
go_total=0
for ((i=1; i<=RUNS; i++)); do
    echo -n "  Run $i: "
    start_time=$(date +%s%3N)
    ./benchmark >/dev/null 2>&1
    end_time=$(date +%s%3N)
    duration_ms=$((end_time - start_time))
    go_total=$((go_total + duration_ms))
    duration_sec=$((duration_ms / 1000))
    duration_frac=$((duration_ms % 1000))
    printf "%d.%03ds\n" $duration_sec $duration_frac
done
go_avg=$((go_total / RUNS))
go_sec=$((go_avg / 1000))
go_frac=$((go_avg % 1000))
printf "  ${GREEN}Average: %d.%03ds${NC}\n\n" $go_sec $go_frac

echo -e "${BLUE}JavaScript Benchmark (${RUNS} runs):${NC}"
js_total=0
for ((i=1; i<=RUNS; i++)); do
    echo -n "  Run $i: "
    start_time=$(date +%s%3N)
    bun script.js >/dev/null 2>&1
    end_time=$(date +%s%3N)
    duration_ms=$((end_time - start_time))
    js_total=$((js_total + duration_ms))
    duration_sec=$((duration_ms / 1000))
    duration_frac=$((duration_ms % 1000))
    printf "%d.%03ds\n" $duration_sec $duration_frac
done
js_avg=$((js_total / RUNS))
js_sec=$((js_avg / 1000))
js_frac=$((js_avg % 1000))
printf "  ${GREEN}Average: %d.%03ds${NC}\n\n" $js_sec $js_frac

# Compare results
echo -e "${BLUE}=== Results ===${NC}"
printf "Go average:         %d.%03ds\n" $go_sec $go_frac
printf "JavaScript average: %d.%03ds\n" $js_sec $js_frac

# Calculate which is faster (avoid division by zero)
if [ $go_avg -lt $js_avg ] && [ $go_avg -gt 0 ]; then
    # Go is faster
    speedup=$((js_avg * 100 / go_avg))
    speedup_int=$((speedup / 100))
    speedup_frac=$((speedup % 100))
    printf "${GREEN}Go is %d.%02dx faster${NC}\n" $speedup_int $speedup_frac
elif [ $js_avg -lt $go_avg ] && [ $js_avg -gt 0 ]; then
    # JavaScript is faster
    speedup=$((go_avg * 100 / js_avg))
    speedup_int=$((speedup / 100))
    speedup_frac=$((speedup % 100))
    printf "${GREEN}JavaScript is %d.%02dx faster${NC}\n" $speedup_int $speedup_frac
else
    echo -e "${YELLOW}Performance is roughly equal${NC}"
fi

# Cleanup
rm -f benchmark

echo -e "\n${BLUE}Benchmark complete!${NC}"