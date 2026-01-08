package butler

import (
	"testing"

	"github.com/koksalmehmet/mind-palace/apps/cli/internal/memory"
)

// TestToolStore tests the brain store/remember functionality
func TestToolStore(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]interface{}
		wantError     bool
		errorContains string
		checkResult   func(*testing.T, jsonRPCResponse)
	}{
		{
			name: "store idea with explicit kind",
			args: map[string]interface{}{
				"content": "Use Redis for caching to improve performance",
				"as":      "idea",
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				if !containsText(result, "stored") || !containsText(result, "idea") {
					t.Errorf("result should confirm idea was stored: %s", result)
				}
			},
		},
		{
			name: "store decision with tags",
			args: map[string]interface{}{
				"content": "Decided to use PostgreSQL for primary database",
				"as":      "decision",
				"tags":    []interface{}{"database", "architecture"},
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				if !containsText(result, "decision") {
					t.Errorf("result should mention decision: %s", result)
				}
			},
		},
		{
			name: "store learning",
			args: map[string]interface{}{
				"content": "Async operations significantly improve throughput",
				"as":      "learning",
			},
			wantError: false,
		},
		{
			name: "store with auto-classification",
			args: map[string]interface{}{
				"content": "Maybe we should consider using GraphQL instead of REST",
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				// Should be auto-classified as idea or decision
				if !containsText(result, "stored") {
					t.Errorf("result should confirm storage: %s", result)
				}
			},
		},
		{
			name: "store with room scope",
			args: map[string]interface{}{
				"content":   "API endpoint optimization ideas",
				"as":        "idea",
				"scope":     "room",
				"scopePath": "api",
			},
			wantError: false,
		},
		{
			name: "store with file scope",
			args: map[string]interface{}{
				"content":   "Refactor needed in this file",
				"as":        "idea",
				"scope":     "file",
				"scopePath": "main.go",
			},
			wantError: false,
		},
		{
			name: "store with multiple tags",
			args: map[string]interface{}{
				"content": "Performance optimization strategy",
				"as":      "idea",
				"tags":    []interface{}{"performance", "optimization", "backend"},
			},
			wantError: false,
		},
		{
			name: "store defaults to palace scope",
			args: map[string]interface{}{
				"content": "General project idea",
				"as":      "idea",
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				// Should use default palace scope
				if resp.Error != nil {
					t.Errorf("should succeed with default scope")
				}
			},
		},
		{
			name: "backward compatibility - kind parameter",
			args: map[string]interface{}{
				"content": "Old style storage",
				"kind":    "idea",
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				// 'kind' should work same as 'as'
				if resp.Error != nil {
					t.Errorf("'kind' parameter should work for backward compatibility")
				}
			},
		},
		{
			name: "missing content",
			args: map[string]interface{}{
				"as": "idea",
			},
			wantError:     true,
			errorContains: "content is required",
		},
		{
			name: "empty content",
			args: map[string]interface{}{
				"content": "",
				"as":      "idea",
			},
			wantError:     true,
			errorContains: "content is required",
		},
		{
			name: "empty tags in array",
			args: map[string]interface{}{
				"content": "Test content",
				"as":      "idea",
				"tags":    []interface{}{"", "valid", ""},
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				// Should filter out empty tags
				if resp.Error != nil {
					t.Errorf("should handle empty tags gracefully")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := setupMCPServer(t)
			resp := server.toolStore(1, tt.args)

			if tt.wantError {
				if resp.Error == nil {
					t.Fatalf("expected error, got none")
				}
				if tt.errorContains != "" && !containsText(resp.Error.Message, tt.errorContains) {
					t.Errorf("error message = %q, should contain %q", resp.Error.Message, tt.errorContains)
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %v", resp.Error)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, resp)
			}
		})
	}
}

// TestToolRecall tests memory recall functionality
func TestToolRecall(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*testing.T, *Butler)
		args          map[string]interface{}
		wantError     bool
		errorContains string
		checkResult   func(*testing.T, jsonRPCResponse)
	}{
		{
			name: "recall all memories",
			setup: func(t *testing.T, b *Butler) {
				// Store some test memories
				storeTestMemory(t, b, "Test idea", memory.KindIdea)
				storeTestMemory(t, b, "Test decision", memory.KindDecision)
			},
			args:      map[string]interface{}{},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				if !containsText(result, "Test idea") || !containsText(result, "Test decision") {
					t.Errorf("should contain both memories: %s", result)
				}
			},
		},
		{
			name: "recall filtered by kind",
			setup: func(t *testing.T, b *Butler) {
				storeTestMemory(t, b, "Idea content", memory.KindIdea)
				storeTestMemory(t, b, "Decision content", memory.KindDecision)
			},
			args: map[string]interface{}{
				"kind": "idea",
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				if !containsText(result, "Idea content") {
					t.Errorf("should contain idea: %s", result)
				}
				if containsText(result, "Decision content") {
					t.Errorf("should not contain decision: %s", result)
				}
			},
		},
		{
			name: "recall with limit",
			setup: func(t *testing.T, b *Butler) {
				for i := 0; i < 10; i++ {
					storeTestMemory(t, b, "Memory content", memory.KindIdea)
				}
			},
			args: map[string]interface{}{
				"limit": float64(5),
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				// Should have limited results
				if resp.Error != nil {
					t.Errorf("should succeed with limit")
				}
			},
		},
		{
			name: "recall with query filter",
			setup: func(t *testing.T, b *Butler) {
				storeTestMemory(t, b, "Redis caching strategy", memory.KindIdea)
				storeTestMemory(t, b, "PostgreSQL database choice", memory.KindDecision)
			},
			args: map[string]interface{}{
				"query": "Redis",
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				if !containsText(result, "Redis") {
					t.Errorf("should contain Redis memory: %s", result)
				}
			},
		},
		{
			name: "recall with tag filter",
			setup: func(t *testing.T, b *Butler) {
				storeTestMemoryWithTags(t, b, "Backend idea", memory.KindIdea, []string{"backend"})
				storeTestMemoryWithTags(t, b, "Frontend idea", memory.KindIdea, []string{"frontend"})
			},
			args: map[string]interface{}{
				"tags": []interface{}{"backend"},
			},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				if !containsText(result, "Backend") {
					t.Errorf("should contain backend memory: %s", result)
				}
			},
		},
		{
			name:      "recall with no memories",
			setup:     func(t *testing.T, b *Butler) {},
			args:      map[string]interface{}{},
			wantError: false,
			checkResult: func(t *testing.T, resp jsonRPCResponse) {
				result := extractToolResult(t, resp)
				if !containsText(result, "No") || !containsText(result, "found") {
					t.Errorf("should indicate no memories found: %s", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, butler := setupMCPServer(t)

			if tt.setup != nil {
				tt.setup(t, butler)
			}

			resp := server.toolRecall(1, tt.args)

			if tt.wantError {
				if resp.Error == nil {
					t.Fatalf("expected error, got none")
				}
				if tt.errorContains != "" && !containsText(resp.Error.Message, tt.errorContains) {
					t.Errorf("error message = %q, should contain %q", resp.Error.Message, tt.errorContains)
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %v", resp.Error)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, resp)
			}
		})
	}
}

// TestToolReflect tests memory reflection/analysis
func TestToolReflect(t *testing.T) {
	server, butler := setupMCPServer(t)

	// Store some contradictory decisions
	storeTestMemory(t, butler, "Use MySQL for database", memory.KindDecision)
	storeTestMemory(t, butler, "Use PostgreSQL for database", memory.KindDecision)

	resp := server.toolReflect(1, map[string]interface{}{
		"context": "workspace",
	})

	if resp.Error != nil {
		t.Fatalf("toolReflect() error = %v", resp.Error)
	}

	result := extractToolResult(t, resp)
	// Should provide some analysis
	if result == "" {
		t.Error("reflection should provide analysis")
	}
}

// TestToolForget tests memory deletion
func TestToolForget(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*testing.T, *Butler) string
		args          map[string]interface{}
		wantError     bool
		errorContains string
	}{
		{
			name: "forget by id",
			setup: func(t *testing.T, b *Butler) string {
				return storeTestMemory(t, b, "To be deleted", memory.KindIdea)
			},
			wantError: false,
		},
		{
			name: "forget nonexistent id",
			setup: func(t *testing.T, b *Butler) string {
				return "nonexistent-id-12345"
			},
			wantError:     true,
			errorContains: "not found",
		},
		{
			name: "forget missing id",
			setup: func(t *testing.T, b *Butler) string {
				return ""
			},
			wantError:     true,
			errorContains: "id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, butler := setupMCPServer(t)

			id := ""
			if tt.setup != nil {
				id = tt.setup(t, butler)
			}

			args := map[string]interface{}{}
			if id != "" {
				args["id"] = id
			}

			resp := server.toolForget(1, args)

			if tt.wantError {
				if resp.Error == nil {
					t.Fatalf("expected error, got none")
				}
				if tt.errorContains != "" && !containsText(resp.Error.Message, tt.errorContains) {
					t.Errorf("error message = %q, should contain %q", resp.Error.Message, tt.errorContains)
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %v", resp.Error)
			}
		})
	}
}

// Helper functions

func extractToolResult(t *testing.T, resp jsonRPCResponse) string {
	t.Helper()

	if resp.Error != nil {
		t.Fatalf("response has error: %v", resp.Error)
	}

	result, ok := resp.Result.(mcpToolResult)
	if !ok {
		t.Fatalf("result type = %T, want mcpToolResult", resp.Result)
	}

	if len(result.Content) == 0 {
		return ""
	}

	return result.Content[0].Text
}

func containsText(text, substr string) bool {
	return text != "" && substr != "" &&
		(text == substr || len(text) > 0 && len(substr) > 0 &&
			(text[0:min(len(text), len(substr))] == substr ||
				findInString(text, substr)))
}

func findInString(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func storeTestMemory(t *testing.T, b *Butler, content string, kind memory.RecordKind) string {
	t.Helper()

	// Mock implementation - actual implementation will use butler's memory
	// This is a placeholder showing the pattern
	rec := memory.Record{
		Content: content,
		Kind:    kind,
		Scope:   memory.ScopePalace,
	}

	// Store via butler's memory system
	// id, err := b.memory.Store(rec)
	// For now, return mock ID
	return "test-id-" + content[:min(10, len(content))]
}

func storeTestMemoryWithTags(t *testing.T, b *Butler, content string, kind memory.RecordKind, tags []string) string {
	t.Helper()

	rec := memory.Record{
		Content: content,
		Kind:    kind,
		Scope:   memory.ScopePalace,
		Tags:    tags,
	}

	// Store via butler's memory system
	return "test-id-" + content[:min(10, len(content))]
}
