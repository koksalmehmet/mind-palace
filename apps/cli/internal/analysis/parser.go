package analysis

type Parser interface {
	Parse(content []byte, filePath string) (*FileAnalysis, error)
	Language() Language
}

type ParserRegistry struct {
	parsers map[Language]Parser
}

func NewParserRegistry() *ParserRegistry {
	reg := &ParserRegistry{
		parsers: make(map[Language]Parser),
	}
	reg.registerDefaults()
	return reg
}

func (r *ParserRegistry) registerDefaults() {
	// NOTE: Tree-sitter based parsers require CGO and are excluded on Windows.
	// LSP-based parsers will be added to replace them (see dart_lsp.go for pattern).
	// For now, only regex-based parsers are available without CGO.

	// Regex-based parsers (no CGO required, work on all platforms)
	r.Register(NewDartParser())
	r.Register(NewCUEParser())

	// TODO: Add LSP-based parsers for major languages:
	// - Go (gopls)
	// - TypeScript/JavaScript (typescript-language-server)
	// - Python (pyright or pylsp)
	// - Rust (rust-analyzer)
	// - etc.
}

func (r *ParserRegistry) Register(p Parser) {
	r.parsers[p.Language()] = p
}

func (r *ParserRegistry) GetParser(lang Language) (Parser, bool) {
	p, ok := r.parsers[lang]
	return p, ok
}

func (r *ParserRegistry) Parse(content []byte, filePath string) (*FileAnalysis, error) {
	lang := DetectLanguage(filePath)
	if lang == LangUnknown {
		return &FileAnalysis{
			Path:     filePath,
			Language: string(LangUnknown),
		}, nil
	}

	parser, ok := r.GetParser(lang)
	if !ok {
		return &FileAnalysis{
			Path:     filePath,
			Language: string(lang),
		}, nil
	}

	return parser.Parse(content, filePath)
}

var defaultRegistry *ParserRegistry

func init() {
	defaultRegistry = NewParserRegistry()
}

func Analyze(content []byte, filePath string) (*FileAnalysis, error) {
	return defaultRegistry.Parse(content, filePath)
}
