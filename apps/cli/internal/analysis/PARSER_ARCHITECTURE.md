# Parser Architecture

## Overview

Mind Palace uses a **three-tier fallback strategy** for code parsing, prioritizing accuracy while ensuring broad compatibility.

## Priority Chain

### 1. LSP (Language Server Protocol) - **Preferred**
**Best for**: Semantic analysis, production accuracy

- **Most accurate** - Uses official language tools with full type information
- **Semantic understanding** - Not just syntax, but types, references, implementations
- **Auto-updated** - Maintained by language communities
- **Call hierarchy** - Accurate caller/callee relationships

**Available Language Servers**:
- **Go**: `gopls` (official)
- **TypeScript/JavaScript**: `typescript-language-server`
- **Python**: `pyright`, `pylsp`
- **Rust**: `rust-analyzer`
- **Java**: `jdtls`
- **C/C++**: `clangd`
- **C#**: `omnisharp`
- **Ruby**: `solargraph`
- **PHP**: `intelephense`
- [Full list](https://langserver.org)

**Implementation**: See `dart_lsp.go` for the LSP client pattern

**Status**: ⏳ **Planned** - Dart LSP client exists as reference implementation

---

### 2. Tree-Sitter - **Current Default**
**Best for**: Fast AST parsing when LSP unavailable

- **Good accuracy** - Proper AST parsing with node types
- **Fast** - Compiled C parsers via CGO
- **30+ languages** - Wide language coverage
- **Requires CGO** - Needs C compiler (gcc/clang/MinGW)

**CI Support**:
- ✅ macOS: Clang/LLVM (built-in)
- ✅ Linux: GCC (built-in)  
- ✅ Windows: MinGW (installed via msys2 in CI)

**Local Development**:
- macOS/Linux: Works out of box
- Windows: Requires MinGW/TDM-GCC installation
  ```powershell
  choco install mingw
  # or download from: https://jmeubank.github.io/tdm-gcc/
  ```

**Status**: ✅ **Active** - All 31 parsers available in CI builds

**Supported Languages**:
- Tree-sitter parsers: Go, TypeScript, JavaScript, Python, Rust, Java, C, C++, C#
- Backend: Ruby, PHP, Kotlin, Scala, Swift
- Infrastructure: Bash, SQL, Dockerfile, HCL
- Config/Web: HTML, CSS, YAML, TOML, JSON, Markdown
- Other: Elixir, Lua, Groovy, Svelte, OCaml, Elm, Protobuf

---

### 3. Regex - **Fallback**
**Best for**: Simple symbol extraction, guaranteed availability

- **Always works** - Pure Go, no dependencies
- **Basic extraction** - Classes, functions, imports
- **Simple patterns** - No AST, just regex matching
- **Limited accuracy** - Can miss nested/complex structures

**Status**: ✅ **Active** - Dart, CUE

**Use Cases**:
- Languages without LSP/tree-sitter support
- Environments without C compiler
- Quick symbol extraction where full AST not needed

---

## Implementation Roadmap

### Phase 1: LSP Infrastructure (Next)
1. Extract generic LSP client from `dart_lsp.go`
2. Create `lsp_client.go` with reusable protocol handling
3. Add configuration for LSP binary paths

### Phase 2: Priority Parsers
Implement LSP parsers for most common languages:
1. **Go** - `gopls` (highest priority)
2. **TypeScript/JavaScript** - `typescript-language-server`
3. **Python** - `pyright`
4. **Rust** - `rust-analyzer`

### Phase 3: Fallback Logic
```go
func (r *ParserRegistry) GetParser(lang Language) Parser {
    // 1. Try LSP if available
    if lspParser, ok := r.lspParsers[lang]; ok {
        if lspParser.IsAvailable() {
            return lspParser
        }
    }
    
    // 2. Fall back to tree-sitter
    if tsParser, ok := r.treeSitterParsers[lang]; ok {
        return tsParser
    }
    
    // 3. Fall back to regex
    return r.regexParsers[lang]
}
```

### Phase 4: Additional Languages
Expand LSP coverage based on usage:
- Java (jdtls)
- C/C++ (clangd)
- C# (omnisharp)
- Ruby (solargraph)
- PHP (intelephense)

---

## Current Status

| Parser Type | Languages | Status | Requires |
|-------------|-----------|--------|----------|
| **LSP** | Dart (example) | ✅ Working | Language server binary |
| **Tree-Sitter** | 31 languages | ✅ CI builds | gcc/clang/MinGW |
| **Regex** | Dart, CUE | ✅ Working | Nothing |

---

## Developer Notes

### Local Development (Windows)
If you want to build locally on Windows with tree-sitter support:

1. **Install MinGW**:
   ```powershell
   choco install mingw
   ```

2. **Or download TDM-GCC**: https://jmeubank.github.io/tdm-gcc/

3. **Verify**:
   ```powershell
   gcc --version
   ```

4. **Build**:
   ```powershell
   $env:CGO_ENABLED="1"
   go build ./apps/cli
   ```

Without MinGW, builds will fail with:
```
cgo: C compiler "gcc" not found
```

### CI Builds
All platforms supported via GitHub Actions:
- MinGW automatically installed on Windows runners
- CGO enabled for all builds
- Full tree-sitter support across platforms

### Adding New LSP Parser
See `dart_lsp.go` and `dart_analyzer.go` for reference implementation:

1. Create `lsp_client_<lang>.go` with language-specific setup
2. Implement `Parser` interface using LSP calls
3. Add to parser registry with priority check
4. Document required language server in README

---

## Benefits of LSP-First Approach

1. **Most Accurate**: Official language tools, semantic analysis
2. **Auto-Updated**: Language servers updated by communities
3. **Rich Features**: Call hierarchy, type information, cross-references
4. **Developer-Friendly**: Tools developers already have installed
5. **Graceful Degradation**: Falls back to tree-sitter if LSP unavailable

---

## References

- **LSP Specification**: https://microsoft.github.io/language-server-protocol/
- **Tree-Sitter**: https://tree-sitter.github.io/tree-sitter/
- **Language Servers**: https://langserver.org/
- **Dart LSP Example**: `dart_lsp.go`, `dart_analyzer.go`
