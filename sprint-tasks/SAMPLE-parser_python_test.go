package analysis

import (
	"strings"
	"testing"
)

// TestPythonParser tests Python code parsing with comprehensive test cases
func TestPythonParser(t *testing.T) {
	parser := NewPythonParser()

	tests := []struct {
		name          string
		code          string
		wantSymbols   int
		checkSymbol   *wantSymbol
		checkRelation *wantRelation
	}{
		{
			name: "simple function",
			code: `def hello():
    pass`,
			wantSymbols: 1,
			checkSymbol: &wantSymbol{
				name: "hello",
				kind: KindFunction,
			},
		},
		{
			name: "function with parameters and return type",
			code: `def greet(name: str, age: int = 0) -> str:
    return f"Hello {name}"`,
			wantSymbols: 1,
			checkSymbol: &wantSymbol{
				name:              "greet",
				kind:              KindFunction,
				signatureContains: []string{"name", "age"},
			},
		},
		{
			name: "function with docstring",
			code: `def process():
    """Process data and return results."""
    pass`,
			wantSymbols: 1,
			checkSymbol: &wantSymbol{
				name:          "process",
				kind:          KindFunction,
				hasDocComment: true,
			},
		},
		{
			name: "empty class",
			code: `class MyClass:
    pass`,
			wantSymbols: 1,
			checkSymbol: &wantSymbol{
				name: "MyClass",
				kind: KindClass,
			},
		},
		{
			name: "class with methods",
			code: `class Calculator:
    def add(self, a, b):
        return a + b
    
    def subtract(self, a, b):
        return a - b`,
			wantSymbols: 3, // 1 class + 2 methods
			checkSymbol: &wantSymbol{
				name:        "Calculator",
				kind:        KindClass,
				hasChildren: true,
			},
		},
		{
			name: "class with __init__",
			code: `class Person:
    def __init__(self, name):
        self.name = name`,
			wantSymbols: 2, // class + __init__
			checkSymbol: &wantSymbol{
				name: "Person",
				kind: KindClass,
			},
		},
		{
			name: "class inheritance",
			code: `class Child(Parent):
    pass`,
			wantSymbols: 1,
			checkRelation: &wantRelation{
				kind:        "inherits",
				hasRelation: true,
			},
		},
		{
			name: "decorated function",
			code: `@property
def value(self):
    return self._value`,
			wantSymbols: 1,
			checkSymbol: &wantSymbol{
				name: "value",
				kind: KindFunction,
			},
		},
		{
			name: "decorated class",
			code: `@dataclass
class Config:
    name: str`,
			wantSymbols: 1,
			checkSymbol: &wantSymbol{
				name: "Config",
				kind: KindClass,
			},
		},
		{
			name: "async function",
			code: `async def fetch_data():
    pass`,
			wantSymbols: 1,
			checkSymbol: &wantSymbol{
				name: "fetch_data",
				kind: KindFunction,
			},
		},
		{
			name: "global variable assignments",
			code: `VERSION = "1.0.0"
CONFIG = {"key": "value"}`,
			wantSymbols: 2,
			checkSymbol: &wantSymbol{
				name: "VERSION",
				kind: KindVariable,
			},
		},
		{
			name: "import statements",
			code: `import os
from typing import List, Dict
import numpy as np`,
			wantSymbols: 0, // imports don't create symbols
			checkRelation: &wantRelation{
				kind:        "import",
				hasRelation: true,
			},
		},
		{
			name: "method calls creating relationships",
			code: `def caller():
    helper()
    obj.method()`,
			wantSymbols: 1,
			checkRelation: &wantRelation{
				kind:        "call",
				hasRelation: true,
			},
		},
		{
			name: "nested functions",
			code: `def outer():
    def inner():
        pass`,
			wantSymbols: 2, // both functions captured
			checkSymbol: &wantSymbol{
				name: "outer",
				kind: KindFunction,
			},
		},
		{
			name: "property with setter",
			code: `class X:
    @property
    def value(self):
        return 42
    
    @value.setter
    def value(self, v):
        pass`,
			wantSymbols: 3, // class + 2 property methods
		},
		{
			name:        "empty file",
			code:        ``,
			wantSymbols: 0,
		},
		{
			name: "comments only",
			code: `# This is a comment
# Another comment`,
			wantSymbols: 0,
		},
		{
			name:        "lambda assignment",
			code:        `square = lambda x: x * x`,
			wantSymbols: 1,
			checkSymbol: &wantSymbol{
				name: "square",
				kind: KindVariable,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse([]byte(tt.code), "test.py")
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			// Check language
			if result.Language != "python" {
				t.Errorf("Language = %q, want %q", result.Language, "python")
			}

			// Check symbol count
			if len(result.Symbols) != tt.wantSymbols {
				t.Errorf("symbols count = %d, want %d", len(result.Symbols), tt.wantSymbols)
				for i, sym := range result.Symbols {
					t.Logf("  [%d] %s (%s)", i, sym.Name, sym.Kind)
				}
			}

			// Check specific symbol properties
			if tt.checkSymbol != nil && len(result.Symbols) > 0 {
				checkSymbolProperties(t, result.Symbols[0], tt.checkSymbol)
			}

			// Check relationships
			if tt.checkRelation != nil {
				checkRelationshipProperties(t, result.Relationships, tt.checkRelation)
			}
		})
	}
}

// TestPythonParserLanguage tests the Language() method
func TestPythonParserLanguage(t *testing.T) {
	parser := NewPythonParser()
	if got := parser.Language(); got != LangPython {
		t.Errorf("Language() = %v, want %v", got, LangPython)
	}
}

// TestPythonParserComplexClass tests parsing of a complex class
func TestPythonParserComplexClass(t *testing.T) {
	parser := NewPythonParser()

	code := `class ComplexService:
    """A complex service class with multiple methods."""
    
    def __init__(self, config):
        """Initialize the service."""
        self.config = config
    
    def fetch(self, id: str) -> dict:
        """Fetch data by ID."""
        return {}
    
    async def save(self, data: dict) -> bool:
        """Save data asynchronously."""
        return True
    
    @property
    def status(self):
        """Get current status."""
        return "active"
    
    @staticmethod
    def validate(data):
        """Validate data."""
        return True
    
    @classmethod
    def from_config(cls, config):
        """Create instance from config."""
        return cls(config)`

	result, err := parser.Parse([]byte(code), "service.py")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should have 1 class + 6 methods
	if len(result.Symbols) < 6 {
		t.Errorf("symbols count = %d, want at least 6", len(result.Symbols))
	}

	// Check class symbol
	var classSymbol *Symbol
	for i := range result.Symbols {
		if result.Symbols[i].Kind == KindClass {
			classSymbol = &result.Symbols[i]
			break
		}
	}

	if classSymbol == nil {
		t.Fatal("class symbol not found")
	}

	if classSymbol.Name != "ComplexService" {
		t.Errorf("class name = %q, want %q", classSymbol.Name, "ComplexService")
	}

	if !strings.Contains(classSymbol.DocComment, "complex service") {
		t.Errorf("class docstring = %q, should contain 'complex service'", classSymbol.DocComment)
	}

	// Check for at least some methods
	methodCount := 0
	for _, sym := range result.Symbols {
		if sym.Kind == KindFunction || sym.Kind == KindMethod {
			methodCount++
		}
	}

	if methodCount < 5 {
		t.Errorf("method count = %d, want at least 5", methodCount)
	}
}

// TestPythonParserErrorHandling tests error cases
func TestPythonParserErrorHandling(t *testing.T) {
	parser := NewPythonParser()

	tests := []struct {
		name     string
		code     string
		filePath string
	}{
		{
			name:     "malformed syntax",
			code:     `def (broken`,
			filePath: "broken.py",
		},
		{
			name:     "incomplete class",
			code:     `class Incomplete`,
			filePath: "incomplete.py",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse([]byte(tt.code), tt.filePath)

			// Parser may or may not error on malformed code
			// If it doesn't error, result should still be valid
			if err != nil {
				t.Logf("Parse() error = %v (expected for malformed code)", err)
				return
			}

			if result == nil {
				t.Fatal("result should not be nil when no error")
			}

			if result.Path != tt.filePath {
				t.Errorf("path = %q, want %q", result.Path, tt.filePath)
			}
		})
	}
}

// Helper types and functions

type wantSymbol struct {
	name              string
	kind              SymbolKind
	signatureContains []string
	hasDocComment     bool
	hasChildren       bool
	exported          bool
}

type wantRelation struct {
	kind        string
	hasRelation bool
}

func checkSymbolProperties(t *testing.T, got Symbol, want *wantSymbol) {
	t.Helper()

	if want.name != "" && got.Name != want.name {
		t.Errorf("symbol name = %q, want %q", got.Name, want.name)
	}

	if want.kind != "" && got.Kind != want.kind {
		t.Errorf("symbol kind = %v, want %v", got.Kind, want.kind)
	}

	for _, substr := range want.signatureContains {
		if !strings.Contains(got.Signature, substr) {
			t.Errorf("signature %q should contain %q", got.Signature, substr)
		}
	}

	if want.hasDocComment && got.DocComment == "" {
		t.Error("symbol should have doc comment")
	}

	if want.hasChildren && len(got.Children) == 0 {
		t.Error("symbol should have children")
	}

	if want.exported && !got.Exported {
		t.Error("symbol should be exported")
	}
}

func checkRelationshipProperties(t *testing.T, rels []Relationship, want *wantRelation) {
	t.Helper()

	if want.hasRelation && len(rels) == 0 {
		t.Error("should have at least one relationship")
		return
	}

	if want.kind != "" {
		found := false
		for _, rel := range rels {
			if rel.Kind == want.kind {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("no relationship of kind %q found", want.kind)
		}
	}
}
