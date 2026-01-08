# Critical Missing Tests - Design & Implementation Plan

**Date:** January 8, 2026  
**Focus:** HIGH IMPACT, QUICK WINS  
**Goal:** Maximum coverage improvement with minimal complexity

---

## Executive Summary

Based on coverage analysis, this document provides detailed designs for the **TOP 9 CRITICAL TESTS** across:

1. **3 Most-Used Language Parsers** (Python, JavaScript, TypeScript)
2. **3 Most-Critical MCP Tools** (Brain/Store, Search/Explore, Briefing)
3. **3 Most-Important HTTP Handlers** (Brain/Remember, Search, Context)

**Expected Total Coverage Improvement:** +25-30%  
**Implementation Effort:** 3-5 days  
**Priority:** P0 - Critical for production readiness

---

## Section 1: Language Parser Tests

### 1.1 Python Parser Tests

**File:** `apps/cli/internal/analysis/parser_python_test.go`

**Current Coverage:** 0% (382 lines uncovered)  
**Target Coverage:** 75%+  
**Impact:** CRITICAL - Python is top 3 most used language

#### Test Structure

```go
package analysis

import (
	"strings"
	"testing"
)

// Table-driven test for Python parsing
func TestPythonParser(t *testing.T) {
	parser := NewPythonParser()

	tests := []struct {
		name          string
		code          string
		wantSymbols   int
		wantKinds     []SymbolKind
		wantNames     []string
		checkRelations bool
	}{
		// Test cases detailed below
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse([]byte(tt.code), "test.py")
			// Assertions
		})
	}
}
```

#### Test Cases (Minimum 15)

1. **Function Definition - Simple**

   ```python
   def hello():
       pass
   ```

   - Expected: 1 symbol, kind=Function, name="hello"

2. **Function Definition - With Parameters**

   ```python
   def greet(name: str, age: int = 0) -> str:
       return f"Hello {name}"
   ```

   - Expected: 1 symbol, signature includes params and return type

3. **Function Definition - With Docstring**

   ```python
   def process():
       """Process data and return results."""
       pass
   ```

   - Expected: Symbol has doc_comment populated

4. **Class Definition - Empty**

   ```python
   class MyClass:
       pass
   ```

   - Expected: 1 symbol, kind=Class, name="MyClass"

5. **Class Definition - With Methods**

   ```python
   class Calculator:
       def add(self, a, b):
           return a + b

       def subtract(self, a, b):
           return a - b
   ```

   - Expected: 3 symbols (1 class + 2 methods), methods are children

6. **Class Definition - With **init****

   ```python
   class Person:
       def __init__(self, name):
           self.name = name
   ```

   - Expected: Class + constructor method

7. **Class Definition - With Inheritance**

   ```python
   class Child(Parent):
       pass
   ```

   - Expected: Relationship of kind "inherits" detected

8. **Decorated Function**

   ```python
   @property
   def value(self):
       return self._value
   ```

   - Expected: Decorator metadata captured

9. **Decorated Class**

   ```python
   @dataclass
   class Config:
       name: str
   ```

   - Expected: Class symbol with decorator info

10. **Async Function**

    ```python
    async def fetch_data():
        pass
    ```

    - Expected: Function marked as async/method

11. **Global Variable Assignment**

    ```python
    VERSION = "1.0.0"
    CONFIG = {"key": "value"}
    ```

    - Expected: 2 symbols, kind=Variable

12. **Import Statements**

    ```python
    import os
    from typing import List, Dict
    import numpy as np
    ```

    - Expected: Relationships of kind "import"

13. **Method Calls (Relationships)**

    ```python
    def caller():
        helper()
        obj.method()
    ```

    - Expected: Call relationships detected

14. **Nested Functions**

    ```python
    def outer():
        def inner():
            pass
    ```

    - Expected: Both functions captured, proper nesting

15. **Property Methods**
    ```python
    class X:
        @property
        def value(self):
            return 42

        @value.setter
        def value(self, v):
            pass
    ```
    - Expected: Both property accessor methods

#### Mock/Stub Requirements

- **None** - Pure parsing, no external dependencies
- Uses tree-sitter which is already mocked via test input

#### Expected Coverage Improvement

- **Before:** 0/382 lines = 0%
- **After:** ~285/382 lines = 75%
- **Net Gain:** +75% coverage on critical parser

---

### 1.2 JavaScript Parser Tests

**File:** `apps/cli/internal/analysis/parser_javascript_test.go`

**Current Coverage:** 0% (352 lines uncovered)  
**Target Coverage:** 75%+  
**Impact:** CRITICAL - JavaScript is top 2 most used

#### Test Structure

```go
func TestJavaScriptParser(t *testing.T) {
	parser := NewJavaScriptParser()

	tests := []struct {
		name        string
		code        string
		wantSymbols []wantSymbol
		wantExports int
	}{
		// Test cases
	}
}

type wantSymbol struct {
	name     string
	kind     SymbolKind
	exported bool
}
```

#### Test Cases (Minimum 15)

1. **Function Declaration**

   ```javascript
   function add(a, b) {
     return a + b;
   }
   ```

2. **Arrow Function - Const**

   ```javascript
   const multiply = (a, b) => a * b;
   ```

3. **Arrow Function - Let**

   ```javascript
   let divide = (a, b) => {
     return a / b;
   };
   ```

4. **Class Declaration**

   ```javascript
   class Person {
     constructor(name) {
       this.name = name;
     }
   }
   ```

5. **Class with Methods**

   ```javascript
   class Calculator {
     add(a, b) {
       return a + b;
     }
     subtract(a, b) {
       return a - b;
     }
   }
   ```

6. **Class with Static Methods**

   ```javascript
   class Utils {
     static format(s) {
       return s.trim();
     }
   }
   ```

7. **Class Inheritance**

   ```javascript
   class Child extends Parent {
     method() {}
   }
   ```

8. **Named Export - Function**

   ```javascript
   export function api() {}
   ```

   - Expected: exported=true

9. **Named Export - Class**

   ```javascript
   export class Service {}
   ```

10. **Default Export**

    ```javascript
    export default class App {}
    ```

11. **Variable Declaration - Const**

    ```javascript
    const API_URL = "https://api.example.com";
    ```

12. **Variable Declaration - Let/Var**

    ```javascript
    let counter = 0;
    var legacy = true;
    ```

13. **Object Method Shorthand**

    ```javascript
    const obj = {
      method() {
        return 42;
      },
    };
    ```

14. **Import Statements**

    ```javascript
    import React from "react";
    import { useState, useEffect } from "react";
    ```

    - Expected: Import relationships

15. **Async Function**
    ```javascript
    async function fetchData() {
      await fetch("/api");
    }
    ```

#### Mock/Stub Requirements

- **None** - Pure parsing

#### Expected Coverage Improvement

- **Before:** 0/352 lines = 0%
- **After:** ~264/352 lines = 75%
- **Net Gain:** +75%

---

### 1.3 TypeScript Parser Tests

**File:** `apps/cli/internal/analysis/parser_typescript_test.go`

**Current Coverage:** 0% (471 lines uncovered)  
**Target Coverage:** 80%+  
**Impact:** CRITICAL - TypeScript is most used in this project

#### Test Structure

```go
func TestTypeScriptParser(t *testing.T) {
	parser := NewTypeScriptParser()

	tests := []struct {
		name            string
		code            string
		wantInterfaces  int
		wantTypes       int
		wantEnums       int
		wantClasses     int
	}{
		// Test cases
	}
}
```

#### Test Cases (Minimum 18 - TypeScript has more constructs)

1. **Function with Type Annotations**

   ```typescript
   function add(a: number, b: number): number {
     return a + b;
   }
   ```

2. **Arrow Function with Types**

   ```typescript
   const multiply = (a: number, b: number): number => a * b;
   ```

3. **Interface Declaration**

   ```typescript
   interface User {
     id: string;
     name: string;
     email?: string;
   }
   ```

4. **Interface with Methods**

   ```typescript
   interface Service {
     fetch(): Promise<Data>;
     save(data: Data): void;
   }
   ```

5. **Interface Inheritance**

   ```typescript
   interface Admin extends User {
     role: string;
   }
   ```

6. **Type Alias - Simple**

   ```typescript
   type ID = string | number;
   ```

7. **Type Alias - Object**

   ```typescript
   type Config = {
     apiUrl: string;
     timeout: number;
   };
   ```

8. **Type Alias - Union**

   ```typescript
   type Result = Success | Error;
   ```

9. **Enum Declaration**

   ```typescript
   enum Status {
     Pending,
     Active,
     Done,
   }
   ```

10. **Const Enum**

    ```typescript
    const enum Color {
      Red = "#ff0000",
      Blue = "#0000ff",
    }
    ```

11. **Class with TypeScript Features**

    ```typescript
    class Service {
      private url: string;
      constructor(url: string) {
        this.url = url;
      }
    }
    ```

12. **Class Implementing Interface**

    ```typescript
    class UserService implements Service {
      fetch(): Promise<User> {}
    }
    ```

13. **Generic Function**

    ```typescript
    function identity<T>(arg: T): T {
      return arg;
    }
    ```

14. **Generic Class**

    ```typescript
    class Box<T> {
      value: T;
    }
    ```

15. **Exported Types**

    ```typescript
    export type { User, Config };
    export interface API {}
    ```

16. **Namespace Declaration**

    ```typescript
    namespace Utils {
      export function format(s: string) {}
    }
    ```

17. **Decorators**

    ```typescript
    @Component
    class AppComponent {}
    ```

18. **Async/Await with Types**
    ```typescript
    async function fetch(): Promise<Response> {
      return await api.get();
    }
    ```

#### Mock/Stub Requirements

- **None** - Pure parsing

#### Expected Coverage Improvement

- **Before:** 0/471 lines = 0%
- **After:** ~377/471 lines = 80%
- **Net Gain:** +80%

---

## Section 2: MCP Tool Handler Tests

### 2.1 Brain/Store Tool Tests

**File:** `apps/cli/internal/butler/mcp_tools_brain_test.go`

**Current Coverage:** ~15% (estimated from mcp_test.go)  
**Target Coverage:** 85%+  
**Impact:** CRITICAL - Core knowledge storage functionality

#### Test Structure

```go
package butler

import (
	"testing"
	"github.com/koksalmehmet/mind-palace/apps/cli/internal/memory"
)

func TestToolStore(t *testing.T) {
	tests := []struct {
		name       string
		args       map[string]interface{}
		wantError  bool
		errorMsg   string
		checkKind  memory.RecordKind
		checkScope string
	}{
		// Test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := setupMCPServer(t)
			resp := server.toolStore(1, tt.args)
			// Assertions
		})
	}
}
```

#### Test Cases (Minimum 12)

1. **Store with Explicit Kind - Idea**

   - Args: `{"content": "Use caching", "as": "idea"}`
   - Expected: Success, kind=idea, scope=palace

2. **Store with Explicit Kind - Decision**

   - Args: `{"content": "Use PostgreSQL", "as": "decision"}`
   - Expected: Success, kind=decision

3. **Store with Explicit Kind - Learning**

   - Args: `{"content": "Async improves perf", "as": "learning"}`
   - Expected: Success, kind=learning

4. **Store with Auto-Classification**

   - Args: `{"content": "We should use Redis for caching"}`
   - Expected: Auto-classified (likely idea/decision)

5. **Store with Tags**

   - Args: `{"content": "Test", "tags": ["performance", "backend"]}`
   - Expected: Tags saved

6. **Store with Room Scope**

   - Args: `{"content": "Test", "scope": "room", "scopePath": "api"}`
   - Expected: Scoped to room

7. **Store with File Scope**

   - Args: `{"content": "Test", "scope": "file", "scopePath": "main.go"}`
   - Expected: Scoped to file

8. **Store - Missing Content**

   - Args: `{"as": "idea"}`
   - Expected: Error "content is required"

9. **Store - Empty Content**

   - Args: `{"content": ""}`
   - Expected: Error

10. **Store - Backward Compatibility (kind param)**

    - Args: `{"content": "Test", "kind": "idea"}`
    - Expected: Works (backward compatible)

11. **Store - Invalid Kind**

    - Args: `{"content": "Test", "as": "invalid_kind"}`
    - Expected: Still works (stored as custom kind or validated)

12. **Store - Multiple Tags Array**
    - Args: `{"content": "Test", "tags": ["t1", "t2", "t3"]}`
    - Expected: All tags saved

#### Additional Test Functions

```go
func TestToolRecall(t *testing.T) {
	// Test memory recall with various filters
}

func TestToolReflect(t *testing.T) {
	// Test reflection/analysis of stored memories
}

func TestToolForget(t *testing.T) {
	// Test deletion of memories
}
```

#### Mock/Stub Requirements

- Memory instance (created by setupMCPServer helper)
- Database with test data
- No external API calls

#### Expected Coverage Improvement

- **Before:** ~134/894 lines = 15%
- **After:** ~760/894 lines = 85%
- **Net Gain:** +70%

---

### 2.2 Search/Explore Tool Tests

**File:** `apps/cli/internal/butler/mcp_tools_search_test.go`

**Current Coverage:** ~25% (basic tests in mcp_test.go)  
**Target Coverage:** 90%+  
**Impact:** HIGH - Primary code discovery mechanism

#### Test Structure

```go
func TestToolExplore(t *testing.T) {
	tests := []struct {
		name         string
		args         map[string]interface{}
		seedData     func(*sql.DB) error
		wantError    bool
		wantContains []string
		wantExclude  []string
	}{
		// Test cases
	}
}
```

#### Test Cases (Minimum 10)

1. **Search - Basic Query**

   - Args: `{"query": "DoWork"}`
   - Expected: Results containing "DoWork"

2. **Search - With Limit**

   - Args: `{"query": "func", "limit": 5}`
   - Expected: Max 5 results

3. **Search - Room Filter**

   - Args: `{"query": "test", "room": "core"}`
   - Expected: Only results from core room

4. **Search - Fuzzy Match Enabled**

   - Args: `{"query": "wrk", "fuzzy": true}`
   - Expected: Matches "work", "worker", etc.

5. **Search - Fuzzy Match Disabled**

   - Args: `{"query": "work", "fuzzy": false}`
   - Expected: Exact matches only

6. **Search - Missing Query**

   - Args: `{}`
   - Expected: Error "query is required"

7. **Search - Empty Query**

   - Args: `{"query": ""}`
   - Expected: Error

8. **Search - Limit Capping (>50)**

   - Args: `{"query": "test", "limit": 100}`
   - Expected: Capped at 50 results

9. **Search - No Results**

   - Args: `{"query": "nonexistent_xyz_123"}`
   - Expected: Empty results, no error

10. **Search - Special Characters**
    - Args: `{"query": "test-function_name"}`
    - Expected: Handles special chars

#### Additional Test Functions

```go
func TestToolExploreImpact(t *testing.T) {
	// Test impact analysis (dependents/dependencies)
}

func TestToolExploreSymbols(t *testing.T) {
	// Test symbol listing with filters
}

func TestToolExploreFile(t *testing.T) {
	// Test file analysis
}

func TestToolExploreDeps(t *testing.T) {
	// Test dependency graph
}

func TestToolExploreCallers(t *testing.T) {
	// Test finding callers of a symbol
}
```

#### Mock/Stub Requirements

- Index database with test symbols
- Butler instance with search configured
- Test rooms defined

#### Expected Coverage Improvement

- **Before:** ~26/103 lines = 25%
- **After:** ~93/103 lines = 90%
- **Net Gain:** +65%

---

### 2.3 Briefing Tool Tests

**File:** `apps/cli/internal/butler/mcp_tools_briefing_test.go`

**Current Coverage:** 0% (277 lines uncovered)  
**Target Coverage:** 75%+  
**Impact:** HIGH - LLM-powered feature, needs reliability

#### Test Structure

```go
func TestToolBriefingSmart(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]interface{}
		mockLLM     bool
		llmResponse string
		wantError   bool
		checkFields []string
	}{
		// Test cases
	}
}
```

#### Test Cases (Minimum 10)

1. **Briefing - Workspace Context**

   - Args: `{"context": "workspace"}`
   - Expected: Workspace-level briefing generated

2. **Briefing - File Context**

   - Args: `{"context": "file", "contextPath": "main.go"}`
   - Expected: File-specific briefing

3. **Briefing - Room Context**

   - Args: `{"context": "room", "contextPath": "api"}`
   - Expected: Room-specific briefing

4. **Briefing - Task Context**

   - Args: `{"context": "task", "contextPath": "Implement auth"}`
   - Expected: Task-focused briefing

5. **Briefing - Summary Style**

   - Args: `{"context": "workspace", "style": "summary"}`
   - Expected: Concise summary format

6. **Briefing - Detailed Style**

   - Args: `{"context": "workspace", "style": "detailed"}`
   - Expected: Comprehensive analysis

7. **Briefing - Actionable Style**

   - Args: `{"context": "workspace", "style": "actionable"}`
   - Expected: Action items highlighted

8. **Briefing - Default Context**

   - Args: `{}`
   - Expected: Defaults to workspace context

9. **Briefing - Invalid Context Type**

   - Args: `{"context": "invalid"}`
   - Expected: Error or defaults

10. **Briefing - Missing Context Path (when required)**
    - Args: `{"context": "file"}`
    - Expected: Error "contextPath required for file context"

#### Mock/Stub Requirements

- **LLM Mock** - Mock anthropic.Client

  ```go
  type mockLLMClient struct {
      response string
      err      error
  }

  func (m *mockLLMClient) Generate(ctx context.Context, req *llm.Request) (*llm.Response, error) {
      if m.err != nil {
          return nil, m.err
      }
      return &llm.Response{Content: m.response}, nil
  }
  ```

- Memory with test records
- Butler with rooms configured

#### Expected Coverage Improvement

- **Before:** 0/277 lines = 0%
- **After:** ~208/277 lines = 75%
- **Net Gain:** +75%

---

## Section 3: HTTP Handler Tests

### 3.1 Brain/Remember Handler Tests

**File:** `apps/cli/internal/dashboard/handlers_brain_test.go`

**Current Coverage:** ~20% (basic tests in handlers_test.go)  
**Target Coverage:** 85%+  
**Impact:** CRITICAL - Primary API for knowledge capture

#### Test Structure

```go
func TestHandleRemember(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		setup      func(*Server)
		wantStatus int
		wantFields map[string]interface{}
	}{
		// Test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := setupServer(t)
			if tt.setup != nil {
				tt.setup(s)
			}

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, "/api/remember",
				strings.NewReader(tt.body))

			s.handleRemember(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			// Additional assertions
		})
	}
}
```

#### Test Cases (Minimum 15)

1. **POST - Valid Idea**

   - Body: `{"content": "Use caching", "kind": "idea"}`
   - Status: 200

2. **POST - Valid Decision**

   - Body: `{"content": "PostgreSQL", "kind": "decision", "tags": ["db"]}`
   - Status: 200

3. **POST - Valid Learning**

   - Body: `{"content": "Async is faster", "kind": "learning"}`
   - Status: 200

4. **POST - Auto-Classification (no kind)**

   - Body: `{"content": "Consider using Redis"}`
   - Status: 200, check classification result

5. **POST - Room Scope**

   - Body: `{"content": "Test", "scope": "room", "scopePath": "api"}`
   - Status: 200

6. **POST - File Scope**

   - Body: `{"content": "Test", "scope": "file", "scopePath": "main.go"}`
   - Status: 200

7. **POST - With Multiple Tags**

   - Body: `{"content": "Test", "tags": ["t1", "t2", "t3"]}`
   - Status: 200

8. **POST - Missing Content**

   - Body: `{"kind": "idea"}`
   - Status: 400

9. **POST - Empty Content**

   - Body: `{"content": ""}`
   - Status: 400

10. **POST - Invalid JSON**

    - Body: `{invalid json}`
    - Status: 400

11. **GET - Method Not Allowed**

    - Method: GET
    - Status: 405

12. **PUT - Method Not Allowed**

    - Method: PUT
    - Status: 405

13. **POST - No Memory Available**

    - Setup: s.memory = nil
    - Status: 503

14. **POST - Default Scope (palace)**

    - Body: `{"content": "Test"}`
    - Expected: scope="palace"

15. **POST - Response Structure Validation**
    - Body: `{"content": "Test"}`
    - Check: response has id, classification, etc.

#### Additional Test Functions

```go
func TestHandleRecall(t *testing.T) {
	// Test GET /api/recall with various filters
}

func TestHandleReflect(t *testing.T) {
	// Test reflection endpoint
}

func TestHandleForget(t *testing.T) {
	// Test DELETE endpoint for memories
}
```

#### Mock/Stub Requirements

- Memory instance with test database
- httptest for requests/responses
- No external dependencies

#### Expected Coverage Improvement

- **Before:** ~65/323 lines = 20%
- **After:** ~275/323 lines = 85%
- **Net Gain:** +65%

---

### 3.2 Search Handler Tests

**File:** Extend `apps/cli/internal/dashboard/handlers_search_test.go` (or add to handlers_test.go)

**Current Coverage:** ~30%  
**Target Coverage:** 90%+  
**Impact:** HIGH - Critical API endpoint

#### Test Cases (Minimum 12)

1. **GET - Basic Search**

   - URL: `/api/search?q=DoWork`
   - Status: 200, results contain symbols

2. **GET - Search with Limit**

   - URL: `/api/search?q=func&limit=5`
   - Status: 200, max 5 results per category

3. **GET - Empty Query**

   - URL: `/api/search?q=`
   - Status: 400

4. **GET - Missing Query**

   - URL: `/api/search`
   - Status: 400

5. **GET - No Results**

   - URL: `/api/search?q=nonexistent_xyz`
   - Status: 200, empty arrays

6. **GET - Special Characters**

   - URL: `/api/search?q=test-func_name`
   - Status: 200

7. **GET - Unicode Query**

   - URL: `/api/search?q=测试`
   - Status: 200

8. **GET - No Butler (symbols)**

   - Setup: s.butler = nil
   - Status: 200, empty symbols array

9. **GET - No Memory (learnings)**

   - Setup: s.memory = nil
   - Status: 200, empty learnings array

10. **GET - No Corridor**

    - Setup: s.corridor = nil
    - Status: 200, empty corridor array

11. **POST - Method Not Allowed**

    - Method: POST
    - Status: 405

12. **GET - Large Limit**
    - URL: `/api/search?q=test&limit=9999`
    - Status: 200, check limit is reasonable

#### Mock/Stub Requirements

- Butler with indexed data
- Memory with learnings
- Corridor with cross-workspace data

#### Expected Coverage Improvement

- **Before:** ~64/214 lines = 30%
- **After:** ~193/214 lines = 90%
- **Net Gain:** +60%

---

### 3.3 Context Handler Tests

**File:** `apps/cli/internal/dashboard/handlers_context_test.go`

**Current Coverage:** 0% (172 lines uncovered)  
**Target Coverage:** 80%+  
**Impact:** HIGH - Auto-injection feature

#### Test Structure

```go
func TestHandleContextPreview(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       interface{}
		setup      func(*Server)
		wantStatus int
		checkFunc  func(*testing.T, map[string]interface{})
	}{
		// Test cases
	}
}
```

#### Test Cases (Minimum 12)

1. **POST - Valid File Path**

   - Body: `{"filePath": "main.go"}`
   - Status: 200

2. **POST - With Max Tokens**

   - Body: `{"filePath": "main.go", "maxTokens": 1000}`
   - Status: 200

3. **POST - Include Learnings**

   - Body: `{"filePath": "main.go", "includeLearnings": true}`
   - Status: 200

4. **POST - Include Decisions**

   - Body: `{"filePath": "main.go", "includeDecisions": true}`
   - Status: 200

5. **POST - Include Failures**

   - Body: `{"filePath": "main.go", "includeFailures": true}`
   - Status: 200

6. **POST - Min Confidence Filter**

   - Body: `{"filePath": "main.go", "minConfidence": 0.8}`
   - Status: 200

7. **POST - All Options Combined**

   - Body: `{"filePath": "main.go", "maxTokens": 2000, "includeLearnings": true, "minConfidence": 0.7}`
   - Status: 200

8. **POST - Missing File Path**

   - Body: `{}`
   - Status: 400

9. **POST - Empty File Path**

   - Body: `{"filePath": ""}`
   - Status: 400

10. **POST - Invalid JSON**

    - Body: `{invalid}`
    - Status: 400

11. **GET - Method Not Allowed**

    - Method: GET
    - Status: 405

12. **POST - No Butler Available**
    - Setup: s.butler = nil
    - Status: 503

#### Additional Test Functions

```go
func TestHandleContextPack(t *testing.T) {
	// Test context pack generation
}

func TestHandleContextValidate(t *testing.T) {
	// Test context validation
}
```

#### Mock/Stub Requirements

- Butler with context building capabilities
- Test files in workspace
- Memory with learnings/decisions

#### Expected Coverage Improvement

- **Before:** 0/172 lines = 0%
- **After:** ~138/172 lines = 80%
- **Net Gain:** +80%

---

## Implementation Priority & Timeline

### Phase 1: Language Parsers (Days 1-2)

1. **Day 1 Morning:** Python parser tests
2. **Day 1 Afternoon:** JavaScript parser tests
3. **Day 2:** TypeScript parser tests

**Deliverable:** +230% coverage on 3 critical parsers

### Phase 2: MCP Tools (Days 3-4)

4. **Day 3:** Brain/Store tool tests
5. **Day 4 Morning:** Search/Explore tool tests
6. **Day 4 Afternoon:** Briefing tool tests

**Deliverable:** +210% coverage on core MCP functionality

### Phase 3: HTTP Handlers (Day 5)

7. **Day 5 Morning:** Brain/Remember handler tests
8. **Day 5 Early Afternoon:** Search handler tests
9. **Day 5 Late Afternoon:** Context handler tests

**Deliverable:** +205% coverage on critical HTTP APIs

---

## Coverage Impact Summary

| Component            | File                       | Before | After | Gain | Lines Added |
| -------------------- | -------------------------- | ------ | ----- | ---- | ----------- |
| Python Parser        | parser_python_test.go      | 0%     | 75%   | +75% | ~400        |
| JavaScript Parser    | parser_javascript_test.go  | 0%     | 75%   | +75% | ~350        |
| TypeScript Parser    | parser_typescript_test.go  | 0%     | 80%   | +80% | ~450        |
| Brain MCP Tool       | mcp_tools_brain_test.go    | 15%    | 85%   | +70% | ~500        |
| Search MCP Tool      | mcp_tools_search_test.go   | 25%    | 90%   | +65% | ~400        |
| Briefing MCP Tool    | mcp_tools_briefing_test.go | 0%     | 75%   | +75% | ~350        |
| Brain HTTP Handler   | handlers_brain_test.go     | 20%    | 85%   | +65% | ~450        |
| Search HTTP Handler  | handlers_search_test.go    | 30%    | 90%   | +60% | ~350        |
| Context HTTP Handler | handlers_context_test.go   | 0%     | 80%   | +80% | ~400        |

**Total Estimated Coverage Improvement: +25-30% overall project coverage**  
**Total Test Code: ~3,650 lines**  
**Implementation Effort: 3-5 days**

---

## Common Patterns & Utilities

### Helper Functions (Reusable)

```go
// For all test files
func setupTestServer(t *testing.T) *Server {
	t.Helper()
	root := t.TempDir()
	// Setup logic
	return server
}

func mustMarshalJSON(t *testing.T, v interface{}) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	return string(data)
}

func assertSymbol(t *testing.T, got Symbol, wantName string, wantKind SymbolKind) {
	t.Helper()
	if got.Name != wantName {
		t.Errorf("name = %q, want %q", got.Name, wantName)
	}
	if got.Kind != wantKind {
		t.Errorf("kind = %q, want %q", got.Kind, wantKind)
	}
}

func assertHTTPStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Errorf("status = %d, want %d", rec.Code, want)
	}
}

func seedTestDatabase(t *testing.T, db *sql.DB, fixtures []Fixture) {
	t.Helper()
	// Common database seeding logic
}
```

---

## Success Metrics

### Quantitative

- [ ] Overall test coverage: 55% → 80%+
- [ ] Language parsers: 0% → 75%+
- [ ] MCP tools: 20% → 85%+
- [ ] HTTP handlers: 25% → 85%+
- [ ] All tests pass in CI/CD
- [ ] Test execution time: <30 seconds

### Qualitative

- [ ] Tests are maintainable (table-driven)
- [ ] No flaky tests
- [ ] Clear error messages
- [ ] Tests document expected behavior
- [ ] Easy to add new test cases

---

## Next Steps

After completing these 9 critical test files:

1. **Run full coverage report:**

   ```bash
   make test-coverage
   ```

2. **Identify remaining gaps** in:

   - Less common language parsers (Rust, Scala, etc.)
   - Secondary MCP tools (Oracle, Postmortem, etc.)
   - Utility functions

3. **Implement integration tests** for:

   - End-to-end MCP workflows
   - Multi-component interactions
   - Real workspace scenarios

4. **Performance benchmarks** for:
   - Parser performance on large files
   - Search query performance
   - Memory query performance

---

## Appendix: Test Template

```go
package analysis

import (
	"testing"
)

func TestXYZParser(t *testing.T) {
	parser := NewXYZParser()

	tests := []struct {
		name        string
		code        string
		wantSymbols int
		wantKind    SymbolKind
		wantError   bool
	}{
		{
			name: "simple function",
			code: `function test() {}`,
			wantSymbols: 1,
			wantKind: KindFunction,
		},
		// More test cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse([]byte(tt.code), "test.ext")

			if (err != nil) != tt.wantError {
				t.Errorf("Parse() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if err != nil {
				return
			}

			if len(result.Symbols) != tt.wantSymbols {
				t.Errorf("symbols count = %d, want %d", len(result.Symbols), tt.wantSymbols)
			}

			if tt.wantSymbols > 0 && result.Symbols[0].Kind != tt.wantKind {
				t.Errorf("symbol kind = %v, want %v", result.Symbols[0].Kind, tt.wantKind)
			}
		})
	}
}
```

---

**End of Test Design Document**
