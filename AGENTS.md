# AGENTS.md

Context and conventions for AI agents working on txt2ig.

## Project Philosophy

**Grug Brain Development**: Keep things simple, avoid over-engineering, prefer boring solutions over clever ones.

### Core Principles

1. **KISS (Keep It Simple, Stupid)**: Every feature should be the simplest possible solution that works
2. **YAGNI (You Ain't Gonna Need It)**: Don't build features until you actually need them
3. **Boring is Good**: Use standard libraries, well-known patterns, minimal dependencies
4. **No Premature Abstractions**: Don't generalize until you've seen the pattern 3+ times

### Grug Brain Rules

- When in doubt, pick the simpler solution
- Avoid clever code (premature optimization is root of all evil)
- Prefer explicit over implicit
- Dependencies are a cost, not a benefit
- Testing is good, but don't over-test
- Single-purpose functions are better than multi-purpose ones
- Avoid abstraction layers (keep code flat, not nested)

## Architecture Decisions

### Current State (v1.0)

#### Tech Stack
- **Language**: Go 1.25 (Golang only, no exceptions)
- **CLI Framework**: Kong (minimal, declarative, boring)
- **Config Format**: JSON5/JSONC (JSON with comments and trailing commas)
- **Font Rendering**: golang.org/x/image/font/opentype
- **Syntax Highlighting**: Chroma (70+ styles, bash support)
- **Embedded Fonts**: GoMono variants (GoMono, GoMonoBold, GoMonoItalic, GoMonoBoldItalic)

#### Core Components

```
internal/
├── cli/           # Kong CLI definitions
├── config/        # Config loading and parsing (JSONC)
├── font/          # Font loading with fallback chain
├── processor/     # Pre and post processors
├── renderer/      # Image rendering and text drawing
└── (web/)         # Future: Web server (stateless)
```

#### Key Patterns

**Config Hierarchy** (fallback chain):
```
custom > local > XDG > home > defaults
```

**Font Fallback Chain**:
```
User Font → System Font → Embedded GoMono
```

**Processors**:
- **Pre-processors**: Text transformation (search/replace, datetime)
- **Post-processors**: Styling (markdown headers, bash syntax)
- **Stateless vs Stateful**: Post-processors can be stateless (line-by-line) or stateful (multi-line)

**Styled Segments** (per-token styling):
```go
type StyledSegment struct {
    Text  string      // Token text
    Style *TextStyle  // Color, Bold, Italic, etc.
}

type StyledLine struct {
    Segments []StyledSegment  // Multiple styles per line
}
```

## Code Conventions

### Go Style

1. **No comments** unless absolutely necessary (code should be self-documenting)
2. **Single-purpose functions**: Each function does one thing well
3. **Early returns**: Guard clauses over nested if-else
4. **Explicit is better than implicit**: Avoid magic, prefer explicit code
5. **Avoid globals**: Pass dependencies explicitly
6. **NoInterfaces until needed**: Concrete types first, interfaces when needed
7. **NoGenerics until needed**: Keep it simple, avoid premature abstraction

### Naming Conventions

```go
// Good: Clear, descriptive, no stutter
func LoadConfig() (*Config, error)  // Good
func LoadFont() (*Font, error)      // Good

// Bad: Stuttering
func ConfigLoad() (*Config, error)  // Bad

// Good: Package provides context
// font.Load() is clear
// font.FontLoad() is stuttering
```

### Error Handling

- **Log warnings, not errors**: `log.Printf("WARN: ...")`
- **Return errors to caller**: Don't swallow errors
- **Wrap errors with context**: `fmt.Errorf("load font: %w", err)`
- **Document error conditions**: Only if not obvious

### Testing Philosophy

- **TDD when appropriate**: Write tests first for complex logic
- **Coverage is not a goal**: Tests should catch bugs, not hit numbers
- **Test behavior, not implementation**: Black-box testing preferred
- **Keep tests simple**: No test frameworks beyond stdlib testing
- **CI/CD tests**: Run `go test ./...` to check all

### File Organization

```
internal/
├── package/
│   ├── file.go          # Main implementation
│   ├── file_test.go     # Tests (same package)
│   └── types.go         # If types are large
```

## Git Workflow

### Branches

- **main/master**: Stable, release-ready code
- **feat/***: Feature branches (e.g., feat/initial-implementation)
- **opencode/***: OpenCode-specific branches

### Commit Messages

Follow conventional commits (boring style):

```
feat: add font family support with bold/italic variants
fix: recognize embedded font names to avoid warnings
docs: update README with installation instructions
test: add tests for styled segments
refactor: simplify config loading logic
style: format code with gofmt
```

### Commit Frequency

- **Commit often**: Small, logical changes
- **One concept per commit**: Don't mix features and fixes
- **Clear messages**: Future you (or another AI) should understand

### Push Policy

- **NO remote pushes**: Unless explicitly asked
- **NO force pushes**: If explicitly asked, warn user first
- **NO rebase on shared branches**: Only rebase local branches

## Working with AI Agents

### When an AI Agent Works on This Project

1. **Read AGENTS.md first**: Understand philosophy and conventions
2. **Read README.md**: Understand user-facing features
3. **Check plans/**: Look for implementation plans
4. **Run `go test ./...`**: Ensure tests pass before starting
5. **Check git status**: See what's been changed recently
6. **Ask clarifying questions**: Don't assume, ask

### Safety Rules

- **Don't run dangerous commands**: Avoid `rm -rf`, `sudo`, or any destructive operations
- **User trust is earned**: Don't assume user will accept all changes
- **When in doubt, ask**: Better to ask than to execute something risky

### Communication Style

- **Be concise**: Short sentences, no preamble
- **Be direct**: Get to the point quickly
- **Show, don't tell**: Code examples over explanations
- **Ask "Why?" questions**: Understand rationale before implementing
- **Challenge complexity**: If a solution seems complex, propose simpler alternative

### Technical Discussions

When discussing technical decisions with user:

1. **Present options clearly**: List A, B, C approaches with pros/cons
2. **State your recommendation**: With rationale
3. **Wait for user decision**: Don't assume
4. **Document decision**: In commit message or code comment (if needed)

### Code Review Checklist

When reviewing AI-generated code:

1. **Does it follow Grug Brain?**: Is it simple?
2. **Is it tested?**: At least basic tests for complex logic
3. **Are errors handled?**: No swallowed errors
4. **Are warnings logged?**: User should know about fallbacks
5. **Is it documented?**: Only if not obvious from code
6. **Are dependencies minimal?**: Each dep should be justified
7. **Is it stateless?**: Prefer stateless over stateful

## Implementation Patterns

### Pattern: Fallback Chain

```go
// Good: Try multiple sources, fallback to default
func (m *Manager) LoadFont(name string) (*Font, error) {
    // Try embedded font
    if embedded := m.loadEmbedded(name); embedded != nil {
        return embedded, nil
    }
    
    // Try system font
    if system := m.loadSystem(name); system != nil {
        return system, nil
    }
    
    // Fallback to default
    log.Printf("WARN: font not found: %s, using embedded default", name)
    return m.loadDefault(), nil
}
```

### Pattern: Processor Pipeline

```go
// Good: Chain processors, each transforms input
func ApplyProcessors(lines []string, configs []interface{}) ([]string, error) {
    for _, cfg := range configs {
        processor, err := ParseProcessor(cfg)
        if err != nil {
            return nil, err
        }
        lines = processor.Process(lines)
    }
    return lines, nil
}
```

### Pattern: Functional Options (Avoid)

```go
// Bad: Over-engineered, not needed for this project
func NewRenderer(opts ...Option) *Renderer

// Good: Simple, explicit, boring
func NewRenderer(config *Config) *Renderer
```

### Pattern: Embedded Assets

```go
// Good: No external dependencies, guaranteed availability
//go:embed fonts/gomono.ttf
var gomonoTTF []byte

func (m *Manager) loadEmbedded(name string) *Font {
    data := getEmbeddedData(name)
    font, _ := opentype.Parse(data)
    return font
}
```

## Common Gotchas

### Go Gotchas

1. **Embedded fonts**: Must use `//go:embed` directive, not runtime file reads
2. **Font face caching**: Create face once, reuse for performance
3. **Color parsing**: Hex colors must start with `#`, not `#` + hex
4. **JSONC**: Use `github.com/titanous/json5` for trailing commas

### Architecture Gotchas

1. **Processor order**: Pre-processors run before text wrapping, post-processors after
2. **Stateful processors**: Must check if line is already styled (check `Style != nil`)
3. **Styled segments**: Each token can have its own style, not one style per line
4. **Font variants**: Bold text requires bold font face, not just style flag

### Web Server Gotchas

1. **Templ generation**: Must run `templ generate` before build
2. **Base64 images**: Encode as data URL, not raw base64
3. **HTMX responses**: Must set `Content-Type: text/html` even for partials
4. **LocalStorage**: Use debouncing (1 second delay) to avoid excessive writes

## Testing Strategy

### Unit Tests

```go
// Good: Test behavior, not implementation
func TestMarkdownBoldHeaders(t *testing.T) {
    p := &MarkdownBoldHeaders{Bold: true}
    line, style, err := p.Process("# My Header")
    
    assert.NoError(t, err)
    assert.Equal(t, "# My Header", line)
    assert.True(t, style.Bold)
}
```

### Integration Tests

```bash
# Good: Test end-to-end with real file
txt2ig test.md -o output.jpg
# Open output.jpg and verify visually
```

### Performance Tests

- **Not needed**: Premature optimization is root of all evil
- **Measure first**: Only optimize when there's a measured problem
- **Keep it simple**: Prefer readable code over fast code

## Questions to Ask User

When in doubt, ask these questions:

1. **Is this the simplest solution that works?**
2. **Do we need this abstraction now?**
3. **Is this stateless? Can it be?**
4. **Does this follow Grug Brain?**
5. **Is this dependency necessary?**
6. **Can this be tested easily?**
7. **Will this be hard to change later?**
8. **Is this boring?** (Boring is good!)

## References

### Project Documents
- `README.md` - User-facing documentation
- `plans/*.md` - Implementation plans and summaries
- `AGENTS.md` - This file (AI agent context)

### External Resources
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Practical Go](https://dave.cheney.net/practical-go)
- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Grug Brain](https://grugbrain.dev/) (Inspiration for philosophy)

### Key Go Packages
- [Kong CLI](https://github.com/alecthomas/kong) - Declarative CLI parsing
- [Chroma](https://github.com/alecthomas/chroma) - Syntax highlighting
- [JSON5](https://github.com/titanous/json5) - JSON with comments
- [golang.org/x/image/font](https://pkg.go.dev/golang.org/x/image/font) - Font rendering

---

**Remember**: Keep it simple, keep it boring, keep it Grug Brain. 🦍
