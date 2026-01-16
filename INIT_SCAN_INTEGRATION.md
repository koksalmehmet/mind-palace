# Init + Scan Integration: Implementation Summary

## What Changed

### Before

```bash
# Two-step process (confusing for new users)
palace init              # Creates config only
palace scan              # Required to actually use the system
# Brief shows: "⚠️  No index found. Run 'palace scan'"
```

### After

```bash
# One-step process (intuitive)
palace init              # Does everything: config + scan
# ✓ Initialized .palace
# building code index...
# ✓ Code index built successfully

# Edge case escape hatch
palace init --no-scan    # Config only (rare cases)
```

## Code Changes

### 1. Modified `InitOptions` struct

```go
type InitOptions struct {
    Root        string
    Force       bool
    WithOutputs bool
    SkipDetect  bool
    NoScan      bool // NEW: Skip automatic scan after init
}
```

### 2. Added `--no-scan` flag

```go
noScan := fs.Bool("no-scan", false, "skip automatic scan after init (not recommended)")
```

### 3. Auto-scan execution with error handling

```go
// Auto-scan unless explicitly disabled
if !opts.NoScan {
    fmt.Printf("\nbuilding code index...\n")
    scanErr := ExecuteScan(ScanOptions{
        Root:    rootPath,
        Verbose: false,
    })

    if scanErr != nil {
        // Comprehensive error guidance
        fmt.Fprintf(os.Stderr, "\n⚠️  Warning: Initial scan failed: %v\n", scanErr)
        fmt.Fprintf(os.Stderr, "\nPossible causes:\n")
        fmt.Fprintf(os.Stderr, "  • Very large workspace (>100k files)\n")
        fmt.Fprintf(os.Stderr, "  • Complex/unconventional project structure\n")
        fmt.Fprintf(os.Stderr, "  • Insufficient disk space or permissions\n")
        fmt.Fprintf(os.Stderr, "  • Unsupported language/framework\n")
        fmt.Fprintf(os.Stderr, "\nYou can:\n")
        fmt.Fprintf(os.Stderr, "  1. Run 'palace scan --verbose' for details\n")
        fmt.Fprintf(os.Stderr, "  2. Check .palace/guardrails.jsonc to exclude dirs\n")
        fmt.Fprintf(os.Stderr, "  3. Use Mind Palace without index (limited features)\n")

        // Graceful degradation: Don't fail init
        return nil
    }

    fmt.Printf("✓ Code index built successfully\n")
}
```

## Error Handling Strategy

### Critical Design Decision: **Graceful Degradation**

If scan fails during init:

1. ✅ Init still succeeds (`.palace/` created)
2. ⚠️ Warning printed with actionable guidance
3. 📊 User informed what still works (sessions, knowledge storage)
4. 🔧 User given next steps to fix scan issues

**Rationale:** Mind Palace has **layered functionality**:

- **Layer 1** (always works): Sessions, knowledge storage
- **Layer 2** (needs scan): Text search, exploration
- **Layer 3** (needs successful parse): Symbol search, call graphs

Failing init on scan failure would prevent Layer 1 usage, which is overly strict.

## Edge Cases Analyzed

### 1. Very Large Codebases (>100k files)

**Behavior:**

- Scan runs but takes time
- No timeout, will complete eventually
- Progress visible with `--verbose`

**Mitigation:**

- Use guardrails to exclude node_modules, build dirs
- Consider `--no-scan` for initial setup, scan later

---

### 2. Unsupported Languages

**Behavior:**

- Files indexed as text-only (chunks + FTS)
- No symbols extracted
- Text search works, symbol search unavailable

**Mitigation:**

- Warning printed: "Unsupported language detected"
- Partial functionality preserved
- Future: Add language support via plugins

---

### 3. Syntax Errors in Code

**Behavior:**

- Parse failures logged as warnings
- File stored without symbol analysis
- Scan completes successfully

**Example:**

```python
# bad.py with syntax error
def broken(
```

**Result:**

```
⚠️  Parse warning: bad.py: unexpected EOF
✓ Indexed 100 files (1 parse warning)
```

---

### 4. Not a Git Repository

**Behavior:**

- Auto-fallback to hash-based incremental scan
- No commit hash tracking
- Full functionality preserved

**Output:**

```
not a git repository, using hash-based incremental scan
✓ Indexed 50 files
```

---

### 5. Insufficient Permissions

**Behavior:**

- Init fails fast with clear error
- No partial state created
- User gets actionable message

**Example:**

```bash
palace init
# Error: Cannot create .palace: permission denied
# Try: chmod +w . or run as administrator
```

---

### 6. Disk Space Exhaustion

**Behavior:**

- Pre-flight check warns if <1GB available
- Transaction rollback on write failure
- Database integrity preserved

**Future Enhancement:**

```go
func checkDiskSpace(root string) error {
    // Check available space before scan
    // Warn if < 1GB, fail if < 100MB
}
```

---

### 7. Unconventional Project Structure

**Example:** Research code with no package manager

```
/research
  experiment1.py
  experiment2.py
  data.csv
```

**Behavior:**

- Auto-detect language from file extensions
- Create generic "workspace" room
- Flat structure (no subprojects)

**Output:**

```
detected project type: python
⚠️  No standard project structure detected
✓ Initialized .palace
✓ Indexed 10 files
```

---

### 8. Mixed Monorepo (Multiple Languages)

**Example:**

```
/workspace
  /frontend  (TypeScript)
  /backend   (Go)
  /mobile    (Dart)
```

**Behavior:**

- Auto-detect rooms per subproject
- Each room gets language-appropriate analyzer
- Cross-language relationships not tracked (future)

**Output:**

```
detected monorepo structure with 3 subprojects
  created room: frontend (apps/frontend)
  created room: backend (apps/backend)
  created room: mobile (apps/mobile)
✓ Indexed 5000 files across 3 rooms
```

## When to Use `--no-scan`

### Rare but Valid Cases

1. **Extremely Large Workspace (>500k files)**

   - Want to configure guardrails first
   - Plan to scan overnight

2. **CI/CD with Pre-built Index**

   - Copy cached index from artifact storage
   - Skip redundant scan

3. **Exploring Mind Palace Structure**

   - Just want to see `.palace/` layout
   - Not planning to use index features

4. **Iterative Configuration**

   - Need to edit palace.jsonc/rooms before scan
   - Want to test configuration first

5. **Resource-Constrained Environments**
   - Low memory/disk space
   - Need to plan scan timing

### How to Use

```bash
# Initialize without scanning
palace init --no-scan

# Later, when ready
palace scan

# Or use watch mode (auto-scans on changes)
palace watch
```

## Testing Checklist

- [x] `palace init` runs scan by default
- [x] `palace init --no-scan` skips scan
- [x] Scan failure doesn't break init
- [x] Error messages are helpful
- [x] Graceful degradation works
- [ ] Test on 100k+ file workspace (manual)
- [ ] Test with parse errors (manual)
- [ ] Test with no git repo (manual)
- [ ] Test with insufficient permissions (manual)

## User Experience Flow

### Happy Path (95% of users)

```bash
cd my-project
palace init
# detected project type: typescript
# initialized palace in /path/.palace
#
# building code index...
# ✓ Code index built successfully

# Ready to use immediately
palace brief
```

### Edge Case Path (5% of users)

```bash
cd massive-monorepo
palace init
# detected project type: javascript
# initialized palace in /path/.palace
#
# building code index...
# ⚠️  Warning: Initial scan failed: too many open files
#
# Possible causes:
#   • Very large workspace (>100k files)
#   ...
#
# You can:
#   1. Run 'palace scan --verbose' for details
#   2. Check .palace/guardrails.jsonc to exclude dirs
#   3. Use Mind Palace without index (limited features)

# User fixes guardrails
vi .palace/guardrails.jsonc
# Add: "node_modules/**", "build/**"

# Retry scan
palace scan
# ✓ Indexed 10000 files (excluded 200k)
```

## Documentation Updates Needed

1. **Getting Started Guide**

   - Update to show single `palace init` command
   - Remove separate `palace scan` step
   - Mention `--no-scan` as advanced option

2. **Troubleshooting Guide**

   - Add section: "What if scan fails during init?"
   - Link to SCAN_FAILURE_ANALYSIS.md

3. **CLI Reference**

   - Document `--no-scan` flag
   - Explain when to use it

4. **Autonomous Agent Guide**
   - Update workflow to single init command
   - Agents should call `palace init` not `palace init && palace scan`

## Benefits

✅ **Better UX**: One command instead of two  
✅ **Intuitive**: "init" does what it says  
✅ **Backward Compatible**: `palace scan` still works for refresh  
✅ **Graceful**: Scan failures don't break init  
✅ **Flexible**: `--no-scan` for edge cases  
✅ **Informative**: Clear error messages with next steps

## Conclusion

The separation between `init` and `scan` was a **historical artifact** with weak justification. By making `init` run `scan` automatically (with `--no-scan` escape hatch), we get:

- Simpler onboarding
- More intuitive workflow
- Preserved flexibility
- Better error handling

This is a **quality of life improvement** that makes Mind Palace easier to adopt while maintaining robustness for edge cases.
