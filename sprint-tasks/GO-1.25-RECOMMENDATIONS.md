# Go 1.25+ Development Best Practices - Authoritative Recommendations

**Last Updated:** January 2026  
**Go Version:** 1.25+ (released August 2025)  
**Sources:** Official go.dev/doc/go1.25, go.dev/blog, go.dev/doc/diagnostics

---

## 1. GO 1.25 FEATURES & IMPROVEMENTS

### Runtime Enhancements

#### Container-Aware GOMAXPROCS (Major Change)
- **New Behavior**: On Linux, Go 1.25 now respects cgroup CPU bandwidth limits
- **Impact**: If CPU bandwidth limit is lower than logical CPUs, `GOMAXPROCS` defaults to the lower limit
- **Automatic Updates**: Runtime periodically updates `GOMAXPROCS` if cgroup limits change
- **When Disabled**: Manually setting `GOMAXPROCS` env var or calling `runtime.GOMAXPROCS()` disables auto-adjustment
- **Debug Flags**: `GODEBUG=containermaxprocs=0` or `GODEBUG=updatemaxprocs=0` to disable
- **Action**: Test CLI apps in containerized environments - may see improved resource utilization without code changes

#### New Experimental Garbage Collector (GreenteaGC)
- **Performance**: Expected 10-40% reduction in GC overhead for programs with heavy GC use
- **Enable**: Set `GOEXPERIMENT=greenteagc` at build time
- **Status**: Experimental - design may evolve; feedback welcome on GitHub issue #73581
- **Recommendation**: Benchmark production workloads with this enabled

#### Trace Flight Recorder API
- **Purpose**: Continuously record runtime traces into in-memory ring buffer
- **Usage**: Call `runtime/trace.FlightRecorder.WriteTo()` to snapshot last few seconds to file
- **Benefits**: Captures rare events with minimal overhead (vs. traditional expensive execution traces)
- **Recommendation**: Implement for daemon applications to diagnose rare production issues

#### Memory Leak Detection (ASAN)
- **New Default**: `go build -asan` now enables leak detection at program exit
- **Reports**: Allocations by C code not freed and not referenced by other memory
- **Disable**: Set `ASAN_OPTIONS=detect_leaks=0` when running
- **Use Case**: Critical for CLI tools with C bindings

### Compiler Improvements

#### Nil Pointer Bug Fix (Breaking Change)
```go
// This code will NOW CORRECTLY PANIC in Go 1.25 (was incorrectly allowed in 1.21-1.24)
f, err := os.Open("nonExistentFile")
name := f.Name()  // Uses f BEFORE checking err
if err != nil {
    return
}
```
- **Action Required**: Review error handling - ensure error checks happen BEFORE using results
- **Impact**: May expose hidden bugs in existing code

#### DWARF5 Support (Optimization)
- **Benefit**: Reduced debug info size, faster linking for large binaries
- **Default**: Now generates DWARF5 debug symbols
- **Disable if Needed**: Set `GOEXPERIMENT=nodwarf5` at build time (may be removed in future)
- **Tooling Note**: Ensure debuggers support DWARF5

#### Faster Stack-Allocated Slices
- **Optimization**: Compiler allocates slice backing stores on stack in more situations
- **Trade-off**: Can amplify effects of incorrect `unsafe.Pointer` usage
- **Debug**: Use `golang.org/x/tools/cmd/bisect` with `-compile=variablemake` to find issues
- **Disable**: Use `-gcflags=all=-d=variablemakehash=n` if needed

### Standard Library Additions

#### testing/synctest Package (Now GA)
```go
import "testing/synctest"

// Run tests with virtualized time
synctest.Test(t, func(t *testing.T) {
    // time package functions operate on fake clock
    // Clock moves instantly when all goroutines block
    synctest.Wait(t)  // Wait for goroutines to block
})
```
- **Purpose**: Test concurrent code without flaky timing issues
- **Benefit**: Deterministic concurrent testing - no sleeps or timeouts needed
- **Recommendation**: Use for testing daemon/concurrent components

#### WaitGroup.Go() Method
```go
var wg sync.WaitGroup
wg.Go(func() {
    // do work
})
// Replaces: wg.Add(1); go func() { ... }()
```
- **Convenience**: Reduces boilerplate for goroutine patterns
- **Safety**: Cleaner API

#### Test Attributes
```go
func TestExample(t *testing.T) {
    t.Attr("database", "postgres")
    t.Attr("region", "us-west")
}
// Output with -json flag shows in "attr" action
```
- **Use Case**: Tag tests with metadata for filtering/reporting

#### New Test Output() Method
- Returns `io.Writer` for test output (like `TB.Log` but without file/line numbers)
- Useful for structured output from tests

#### JSON v2 (Experimental)
- **Enable**: `GOEXPERIMENT=jsonv2` at build time
- **Status**: Substantially faster decoding (parity or faster encoding)
- **Action**: Test with your codebase, provide feedback on issue #71497

### Other Notable Changes

#### go doc -http
- Serves documentation locally in browser
- Useful for development workflow

#### net/http.CrossOriginProtection
- CSRF protection using Fetch metadata headers
- No tokens required
- Supports origin-based and pattern-based bypasses

#### io/fs.ReadLinkFS Interface
- New interface for reading symbolic links
- Implemented by `os.DirFS`, `os.Root.FS`, `testing/fstest.MapFS`

#### runtime.SetDefaultGOMAXPROCS()
- Sets `GOMAXPROCS` to default (as if env var not set)
- Useful for re-enabling auto-adjustment after manual override

---

## 2. SQLITE INTEGRATION - PURE-GO WITH WAL MODE

### Library Recommendation

**Best Choice for Pure-Go SQLite: `github.com/mattn/go-sqlite3` with Pure-Go Drivers**

Alternative (fully pure-Go, no CGO):
- `zombieland/go-sqlite3` - Pure Go implementation
- `github.com/ncruces/go-sqlite3` - Modern pure-Go SQLite binding

### WAL (Write-Ahead Logging) Mode Best Practices

#### Enable WAL Mode
```go
db, err := sql.Open("sqlite3", "file:mydb.db")
if err != nil {
    return err
}
defer db.Close()

// Enable WAL mode
if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
    return err
}

// Recommended WAL settings
settings := []string{
    "PRAGMA synchronous=NORMAL",        // Balance safety/performance
    "PRAGMA wal_autocheckpoint=1000",   // Checkpoint every 1000 pages
    "PRAGMA cache_size=10000",          // ~40MB cache
    "PRAGMA mmap_size=30000000",        // 30MB memory-mapped I/O
    "PRAGMA foreign_keys=ON",
}

for _, setting := range settings {
    if _, err := db.Exec(setting); err != nil {
        return err
    }
}
```

#### WAL Benefits for CLI/Daemon Apps
- **Concurrent Readers**: Multiple readers while writes happen
- **Performance**: Batched writes (checkpoints) instead of individual commits
- **Durability**: Write-ahead logging provides crash safety
- **Best for**: Event logging, analytics, state persistence

#### WAL Considerations
- **Checkpoint Management**: Manual checkpoints may be needed
  ```go
  // Force checkpoint
  if _, err := db.Exec("PRAGMA wal_checkpoint(RESTART)"); err != nil {
      return err
  }
  ```
- **Storage**: Creates `-wal` and `-shm` files alongside main DB
- **Network Filesystems**: Don't use WAL on network shares (NFS, SMB)
- **Multiple Processes**: WAL works but limit to one writer process

#### Connection Pool Settings for CLI/Daemon
```go
db.SetMaxOpenConns(25)      // Max concurrent connections
db.SetMaxIdleConns(5)       // Idle connections to keep
db.SetConnMaxLifetime(0)    // No lifetime limit
db.SetConnMaxIdleTime(time.Hour)  // Close idle after 1 hour
```

### Memory Profiling with SQLite
```go
// Track memory usage
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("Allocated: %v MB\n", m.Alloc/1024/1024)
fmt.Printf("TotalAlloc: %v MB\n", m.TotalAlloc/1024/1024)
fmt.Printf("HeapAlloc: %v MB\n", m.HeapAlloc/1024/1024)

// With pprof
import _ "net/http/pprof"
// GET http://localhost:6060/debug/pprof/heap
```

---

## 3. CLI PERFORMANCE OPTIMIZATION PATTERNS

### Startup Time Optimization
```go
// Profile startup time
import "os"
import "runtime/debug"

func init() {
    debug.SetGCPercent(50)  // Aggressive GC for memory-heavy startup
}

func main() {
    defer func() {
        debug.SetGCPercent(100)  // Reset to default
    }()
    
    // Application code
}
```

### Memory-Efficient Large File Processing
```go
// Use buffered readers/writers
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()
    
    // Allocate once, reuse
    buf := make([]byte, 32*1024)  // 32KB buffer
    for {
        n, err := f.Read(buf)
        if err != nil && err != io.EOF {
            return err
        }
        if n == 0 {
            break
        }
        
        // Process buf[:n]
    }
    return nil
}
```

### Concurrent CLI Operations with Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

results := make(chan Result, numWorkers)
for i := 0; i < numWorkers; i++ {
    go func() {
        select {
        case <-ctx.Done():
            return
        default:
            // Do work
        }
    }()
}
```

### Resource Cleanup Pattern
```go
// Use defer chains for guaranteed cleanup
func processWithResources() error {
    f, err := os.Open("file.txt")
    if err != nil {
        return err
    }
    defer f.Close()
    
    db, err := sql.Open("sqlite3", "db.sqlite")
    if err != nil {
        f.Close()  // Manual cleanup if second fails
        return err
    }
    defer db.Close()
    
    // Process
    return nil
}
```

### Exit Code Best Practices
```go
func main() {
    if err := run(context.Background()); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func run(ctx context.Context) error {
    // All business logic here
    return nil
}
```

---

## 4. ERROR HANDLING & WRAPPING BEST PRACTICES

### Sentinel Errors (Go 1.13+)
```go
var (
    ErrNotFound = errors.New("not found")
    ErrTimeout  = errors.New("timeout")
)

// Check with errors.Is()
if err := getValue(); err != nil {
    if errors.Is(err, ErrNotFound) {
        fmt.Println("Value not found")
    } else if errors.Is(err, ErrTimeout) {
        fmt.Println("Operation timed out")
    }
}
```

### Error Wrapping Best Practices
```go
import "fmt"
import "errors"

// Always wrap errors with context
func readConfig(path string) (Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, fmt.Errorf("read config: %w", err)
    }
    
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return Config{}, fmt.Errorf("parse config: %w", err)
    }
    
    return cfg, nil
}

// Check with errors.Is() and errors.As()
var cfg Config
if err := readConfig("app.json"); err != nil {
    var pathErr *os.PathError
    if errors.As(err, &pathErr) {
        fmt.Printf("Path error on %s\n", pathErr.Path)
    }
}
```

### Custom Error Types (When Needed)
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message)
}

// Check with errors.As()
if err := validate(input); err != nil {
    var ve ValidationError
    if errors.As(err, &ve) {
        fmt.Printf("Invalid field: %s\n", ve.Field)
    }
}
```

### Panic/Recover Pattern (Rare)
```go
// Only use internally, convert to errors
func safeExecute(fn func()) error {
    defer func() {
        if r := recover(); r != nil {
            // Log and handle
            fmt.Printf("Recovered: %v\n", r)
        }
    }()
    
    fn()
    return nil
}
```

### Go 1.25 New Panic Output
- Panics recovered and repanicked now show: `panic: <value> [recovered, repanicked]`
- No longer repeats panic text (shorter, cleaner output)

---

## 5. TESTING PATTERNS FOR CLI & DAEMON APPLICATIONS

### Table-Driven Tests (Recommended)
```go
func TestParseCommand(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Command
        wantErr bool
    }{
        {"simple", "list", Command{Op: "list"}, false},
        {"with_args", "delete id=5", Command{Op: "delete", ID: "5"}, false},
        {"invalid", "unknown", Command{}, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := parseCommand(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("parseCommand(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
            }
            if err == nil && got != tt.want {
                t.Errorf("parseCommand(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}
```

### Subtests for Organization
```go
func TestServer(t *testing.T) {
    t.Run("startup", func(t *testing.T) {
        // Test initialization
    })
    
    t.Run("requests", func(t *testing.T) {
        t.Run("valid", func(t *testing.T) {
            // Test valid requests
        })
        t.Run("invalid", func(t *testing.T) {
            // Test invalid requests
        })
    })
    
    t.Run("shutdown", func(t *testing.T) {
        // Test cleanup
    })
}
```

### Concurrent Testing (Go 1.25)
```go
import "testing/synctest"

func TestConcurrentOperations(t *testing.T) {
    synctest.Test(t, func(t *testing.T) {
        var done int32
        
        go func() {
            // Concurrent work
            atomic.StoreInt32(&done, 1)
        }()
        
        // Time progresses instantly when goroutines block
        synctest.Wait(t)
        
        if atomic.LoadInt32(&done) != 1 {
            t.Fatal("work not completed")
        }
    })
}
```

### Setup/Teardown Pattern
```go
func TestDatabase(t *testing.T) {
    // Setup
    db := setupDB(t)
    defer db.Close()
    
    t.Run("insert", func(t *testing.T) {
        // Use db
    })
    
    t.Run("query", func(t *testing.T) {
        // Use db
    })
    
    // Teardown happens when all subtests complete
}
```

### Parallel Tests
```go
func TestProcessing(t *testing.T) {
    tests := []struct {
        name string
        data string
    }{
        {"case1", "data1"},
        {"case2", "data2"},
    }
    
    for _, tt := range tests {
        tt := tt  // Capture for closure
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // Run in parallel with others
            result := process(tt.data)
            if result == "" {
                t.Error("unexpected empty result")
            }
        })
    }
}
```

### Daemon Application Testing
```go
func TestDaemon(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Start daemon
    done := make(chan error, 1)
    go func() {
        done <- runDaemon(ctx)
    }()
    
    // Give it time to initialize
    time.Sleep(100 * time.Millisecond)
    
    // Test operations
    if err := callDaemon(); err != nil {
        t.Fatalf("daemon call failed: %v", err)
    }
    
    // Cleanup
    cancel()
    if err := <-done; err != nil && err != context.Canceled {
        t.Errorf("daemon error: %v", err)
    }
}
```

### Benchmark Patterns
```go
func BenchmarkOperation(b *testing.B) {
    cases := []struct {
        name string
        size int
    }{
        {"small", 100},
        {"medium", 1000},
        {"large", 10000},
    }
    
    for _, bc := range cases {
        b.Run(bc.name, func(b *testing.B) {
            data := prepareData(bc.size)
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                process(data)
            }
        })
    }
}
```

---

## 6. CONTEXT & DEADLINE MANAGEMENT IN CONCURRENT APPS

### Basic Context Patterns
```go
// Create root context
ctx := context.Background()

// Add timeout
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()  // IMPORTANT: Always call cancel

// Add deadline
deadline := time.Now().Add(5*time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// Add cancellation
ctx, cancel := context.WithCancel(context.Background())
// Cancel when needed
cancel()

// Add values (use custom types for keys)
type key string
ctx = context.WithValue(ctx, key("user"), "alice")
```

### Propagating Context Through Functions
```go
// ALWAYS accept context as first parameter
func fetchData(ctx context.Context, id string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return io.ReadAll(resp.Body)
}
```

### Goroutine Lifecycle with Context
```go
func worker(ctx context.Context, jobs <-chan Job, results chan<- Result) {
    for {
        select {
        case <-ctx.Done():
            // Context canceled/deadline exceeded
            return
        case job, ok := <-jobs:
            if !ok {
                return
            }
            // Process job
            results <- process(job)
        }
    }
}

// Usage
ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()

for i := 0; i < numWorkers; i++ {
    go worker(ctx, jobs, results)
}
```

### Timeout Pattern for Operations
```go
func doWithTimeout(ctx context.Context, duration time.Duration, fn func() error) error {
    ctx, cancel := context.WithTimeout(ctx, duration)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        done <- fn()
    }()
    
    select {
    case err := <-done:
        return err
    case <-ctx.Done():
        return ctx.Err()  // Returns context.DeadlineExceeded
    }
}
```

### Deadline Checking Before Starting Work
```go
func maybeDoWork(ctx context.Context) error {
    deadline, ok := ctx.Deadline()
    if ok && time.Until(deadline) < 5*time.Second {
        // Not enough time, bail out early
        return context.DeadlineExceeded
    }
    
    // Do work
    return nil
}
```

### HTTP Server Graceful Shutdown
```go
func main() {
    server := &http.Server{Addr: ":8080"}
    
    // Start server in goroutine
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("server error: %v", err)
        }
    }()
    
    // Wait for interrupt
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("shutdown error: %v", err)
    }
}
```

### Request-Scoped Values (Advanced)
```go
type contextKey string

const userKey contextKey = "user"

func setUserInContext(ctx context.Context, user string) context.Context {
    return context.WithValue(ctx, userKey, user)
}

func getUserFromContext(ctx context.Context) (string, bool) {
    user, ok := ctx.Value(userKey).(string)
    return user, ok
}

// In handler
func handleRequest(ctx context.Context) {
    if user, ok := getUserFromContext(ctx); ok {
        fmt.Printf("User: %s\n", user)
    }
}
```

---

## 7. MEMORY PROFILING & LEAK DETECTION

### Enable Profiling in CLI
```go
import (
    "flag"
    "log"
    "os"
    "runtime/pprof"
)

func main() {
    cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
    memProfile := flag.String("memprofile", "", "write memory profile to file")
    flag.Parse()
    
    if *cpuProfile != "" {
        f, err := os.Create(*cpuProfile)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()
        
        if err := pprof.StartCPUProfile(f); err != nil {
            log.Fatal(err)
        }
        defer pprof.StopCPUProfile()
    }
    
    // Application logic
    
    if *memProfile != "" {
        f, err := os.Create(*memProfile)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()
        
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal(err)
        }
    }
}
```

### Usage
```bash
# Run with CPU profiling
./myapp -cpuprofile=cpu.prof

# Run with memory profiling
./myapp -memprofile=mem.prof

# Analyze
go tool pprof cpu.prof
go tool pprof mem.prof
```

### Runtime Memory Statistics
```go
import "runtime"

func reportMemory() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Alloc = %v MB\n", m.Alloc/1024/1024)
    fmt.Printf("TotalAlloc = %v MB\n", m.TotalAlloc/1024/1024)
    fmt.Printf("Sys = %v MB\n", m.Sys/1024/1024)
    fmt.Printf("NumGC = %v\n", m.NumGC)
    fmt.Printf("NumGoroutine = %v\n", runtime.NumGoroutine())
}
```

### Goroutine Leak Detection
```go
func testForLeaks(t *testing.T, fn func()) {
    // Count goroutines before
    before := runtime.NumGoroutine()
    
    // Run test
    fn()
    
    // Give time for cleanup
    time.Sleep(100 * time.Millisecond)
    
    // Count after
    after := runtime.NumGoroutine()
    
    if after > before {
        t.Fatalf("goroutine leak: before=%d, after=%d", before, after)
    }
}
```

### Heap Profiling in HTTP Daemon
```go
import _ "net/http/pprof"

func main() {
    // Automatically registers /debug/pprof endpoints
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Your daemon runs here
}

// Access profiles:
// curl http://localhost:6060/debug/pprof/heap > heap.prof
// go tool pprof heap.prof
```

### Go 1.25 Memory Leak Detection
```bash
# Enable ASAN leak detection
go build -asan

# Run with leak detection enabled (default in 1.25)
./myapp

# Disable if needed
ASAN_OPTIONS=detect_leaks=0 ./myapp
```

### GODEBUG Settings for Memory Analysis
```bash
# Print GC events
GODEBUG=gctrace=1 ./myapp

# Print init timing
GODEBUG=inittrace=1 ./myapp

# New in Go 1.25: Finalizer diagnostics
GODEBUG=checkfinalizers=1 ./myapp
```

### Identifying Memory Leaks
```go
// Check for stuck goroutines
func checkForBlockedGoroutines() {
    stack := runtime.Stack(make([]byte, 1<<20), true)
    fmt.Println(string(stack))
}

// Dump heap
func dumpHeap(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()
    
    if err := debug.WriteHeapDump(f.Fd()); err != nil {
        return err
    }
    
    fmt.Printf("Heap dumped to %s\n", filename)
    return nil
}
```

### Profiling Daemon Applications
```bash
# Profile running daemon (periodically)
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap

# CPU profile (30 seconds)
curl 'http://localhost:6060/debug/pprof/profile?seconds=30' > cpu.prof
go tool pprof -http=:8080 cpu.prof

# Goroutine profile
curl http://localhost:6060/debug/pprof/goroutine > goroutines.prof
go tool pprof goroutines.prof
```

---

## 8. BREAKING CHANGES GO 1.24 → 1.25

### Critical

1. **Nil Pointer Bug Fix**
   - Code accessing results BEFORE error checks now panics (was incorrectly allowed)
   - Action: Reorder error checks to come first

2. **Container-Aware GOMAXPROCS**
   - Default behavior changes in containerized environments
   - May see performance changes without code modification
   - Test in target environment

### Minor

3. **DWARF5 Debug Symbols**
   - Default format changed
   - Requires DWARF5-capable debuggers
   - Can disable with `GOEXPERIMENT=nodwarf5`

4. **Panic Output Format**
   - Recovered/repanicked panics show `[recovered, repanicked]` instead of repeating value
   - Affects log parsing if it depends on exact format

5. **Stack Allocations for Slices**
   - Slice allocations on stack vs heap may change
   - Can expose `unsafe.Pointer` bugs
   - Use bisect tool to diagnose

6. **Deprecated APIs**
   - `go/ast.FilterPackage`, `PackageExports`, `MergePackageFiles` (deprecated)
   - `go/parser.ParseDir` (deprecated)
   - Migrate to alternatives

### Not Breaking But Behavioral Changes

7. **Crypto/TLS**
   - SHA-1 now disallowed in TLS 1.2 (can re-enable with `GODEBUG=tlssha1=1`)
   - FIPS mode stricter about key exchange methods

8. **ASN.1 Parsing**
   - More strict parsing of T61String and BMPString types
   - May reject previously accepted malformed encodings

9. **Testing**
   - `AllocsPerRun` panics if parallel tests running (to catch flaky tests)

---

## 9. RECOMMENDED PRACTICES FOR YOUR CODEBASE

### For CLI Applications
1. **Startup**: Use profiling to identify slow init
2. **Resource Cleanup**: Use defer chains, context timeouts
3. **Error Handling**: Wrap all errors with context, use errors.Is/As
4. **Exit Codes**: Follow convention (0=success, 1=error)

### For Daemon Applications
1. **Graceful Shutdown**: Use context with timeout
2. **Resource Management**: Monitor with pprof
3. **Concurrent Testing**: Use testing/synctest
4. **Health Checks**: Expose /debug/pprof endpoints

### For SQLite Usage
1. **Enable WAL**: For concurrent readers
2. **Connection Pooling**: Set reasonable pool size
3. **Pragmas**: Configure for your workload
4. **Monitoring**: Track memory with pprof

### For Performance-Critical Code
1. **Profile First**: Use CPU/heap profiles to find bottlenecks
2. **Benchmark**: Use table-driven benchmarks
3. **Memory**: Pre-allocate where possible, track with ReadMemStats
4. **Concurrency**: Use context for lifecycle management

### For Production Deployment
1. **Health Monitoring**: Enable /debug/pprof endpoints (restricted!)
2. **Memory Tracking**: Periodic heap dumps or memory stats
3. **Container Awareness**: Go 1.25 GOMAXPROCS should handle automatically
4. **Graceful Shutdown**: Implement signal handling with context deadline

---

## 10. QUICK REFERENCE: GO 1.25 BUILD FLAGS

```bash
# Standard build
go build

# With profiling support
go build

# No optimizations (for debugging)
go build -gcflags=all="-N -l"

# With DWARF location lists (optimized + debug)
go build -gcflags="-dwarflocationlists=true"

# ASAN with leak detection (Go 1.25 default)
go build -asan

# Experimental greenteagc
go build -tags=greenteagc

# Experimental jsonv2
go build -tags=jsonv2

# Experimental synctest
go build -tags=synctest
```

---

## Summary & Action Items

✅ **Verify**: Nil pointer error handling in existing code  
✅ **Test**: Container-aware GOMAXPROCS in target deployment environment  
✅ **Implement**: Graceful shutdown with context in daemon apps  
✅ **Add**: Table-driven tests with subtests pattern  
✅ **Enable**: SQLite WAL mode for concurrent access  
✅ **Setup**: Memory profiling endpoints for production diagnostics  
✅ **Update**: Error handling to use errors.Is/As consistently  
✅ **Benchmark**: CLI startup and critical paths with Go 1.25 compiler optimizations  
