# Hungarian Lottery System

## To run the project

```
go run .
```

## Asymptotic Runtime Analysis

The lottery system consists of two main operations: processing player picks and matching winning picks. Here's the runtime analysis:

### 1. ProcessPlayerPicks
- Time Complexity: O(N * K)
  - N = number of players (lines in input file)
  - K = numberOfPicks (constant, typically 5)
- Space Complexity: O(N)
  - Stores N LotteryPick structs, each containing two uint64 values

Key operations:
- Reading file line by line: O(N)
- For each line:
  - Splitting string: O(K)
  - Converting K numbers to integers: O(K)
  - Setting K bits in LotteryPick: O(K)
  - Total per line: O(K)

### 2. MatchPicks
- Time Complexity: O(N / P) where:
  - N = number of player picks
  - P = number of CPU cores
- Space Complexity: O(M) where:
  - M = numberOfPicks - minMatches + 1 (size of winners of each category)

Key operations:
- Parsing winning entry: O(K) 
- Parallel matching:
  - Split into P chunks: O(1)
  - For each chunk (~N/P picks):
    - Bitwise AND: O(1) - operates on 128 bits
    - Population count: O(1) - hardware-optimized bits.OnesCount64
  - Aggregating results: O(P * M) where M is typically small (e.g., 4)

Overall runtime is O(N)

## Optimization Ideas

### 1. Performance Improvements
- **Caching**
  - Cache validated pick combinations
  - Faster processing of repeated inputs
  - Trade-off - Overhead memory usage if every pick is unique

- **Pre-computation**
  - Pre-compute bit masks for common pick combinations. For e.g. players pick birthdays etc.
  - Faster matching for frequent patterns
  - Trade-off - Increased memory usage

- **Batch Processing**
  - Process batch of picks using Single Instruction Multiple Data (SIMD) instructions. For. e.g. using 256 bit registers to perform AND operation on batch of two picks
  
- **Bit Packing Optimization**
  - Use single uint64 if maxPick ≤ 64
  - Halves memory usage

### 2. Handling More Players
- **Memory Efficiency**
  - Implement streaming processing instead of loading all picks
  - Process file in chunks, maintaining winners count for each chunk
  - Constant memory usage regardless of number of players
  - Trade-off: Can't reprocess same dataset multiple times

- **Distributed Processing**
  - Split player picks across multiple machines
  - Use a map-reduce framework
  - Scales horizontally with number of players

- **Database Integration**
  - Store player picks in a database with indexing
  - Query the database to find winners, rather than doing all the matching in memory
  - Trade-off: Initial loading time and storage overhead

These optimizations could be implemented based on specific use case requirements - memory constraints, number of players, frequency of execution.