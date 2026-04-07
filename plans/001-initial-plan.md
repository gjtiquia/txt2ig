# txt2ig Implementation Plan - Phase 1

## Project Overview

**Goal**: Build a CLI tool that converts text files to images for Instagram postings.

**Philosophy**: Simple, minimal, "grug brain" approach. Easy to unit test. Pure functions where possible.

**Tech Stack**: Golang only.

---

## Project Structure

```
txt2ig/
├── main.go                         # Entry point
├── go.mod
├── go.sum
├── README.md
├── plans/                          # Implementation plans
│   └── 001-initial-plan.md
└── internal/
    ├── config/
    │   ├── config.go               # Config struct types
    │   ├── loader.go               # Load config from file hierarchy
    │   └── jsonc.go                # JSONC (JSON with Comments) parser
    ├── cli/
    │   └── cli.go                  # Kong CLI definition
    ├── processor/
    │   ├── processor.go            # Processor interface
    │   ├── pre_processors.go       # Built-in pre-processors
    │   └── post_processors.go      # Built-in post-processors
    ├── renderer/
    │   ├── renderer.go             # Main renderer interface
    │   ├── text.go                 # Text layout, wrapping, metrics
    │   └── image.go                # Image creation and saving
    └── font/
        ├── font.go                 # Font manager
        └── fallback.go             # Font fallback chain logic
```

**Key Decision**: All code in `internal/`, no `pkg/`. Entry point is root `main.go`, no `/cmd` directory.

---

## Dependencies

### Core Dependencies
```go
require (
    github.com/alecthomas/kong v1.0.1      // CLI parsing
    golang.org/x/image v0.38.0             // Font rendering and drawing
)
```

### Standard Library
- `image`, `image/color`, `image/draw` - Image creation and manipulation
- `image/jpeg`, `image/png` - Image encoding
- `os`, `io` - File I/O
- `path/filepath` - Path handling
- `fmt`, `errors` - Error handling
- `strings`, `regexp` - Text processing
- `time` - DateTime processing

---

## Configuration

### Config Struct

```go
type Config struct {
    // Font
    Font            []string    `json:"font"`            // Fallback chain, supports file paths and font names
    FontSize        int         `json:"fontSize"`        // pixels
    FontColor       string      `json:"fontColor"`       // hex color (#FFFFFF)
    
    // Layout
    BgColor         string      `json:"bgColor"`         // hex color (#000000)
    ScreenSize      []int       `json:"screenSize"`      // [width, height] pixels
    TextBoxJustify  string      `json:"textBoxJustify"`  // "left", "center", "right"
    TextBoxAlign    string      `json:"textBoxAlign"`    // "top", "center", "bottom"
    TextBoxOffset   []int       `json:"textBoxOffset"`   // [x, y] pixels offset from anchor
    TextBoxMaxWidth int         `json:"textBoxMaxWidth"` // max width for text box, 0 = auto (90% of screen)
    
    // Text
    TextJustify     string      `json:"textJustify"`     // "left", "center", "right"
    TextWrap        bool        `json:"textWrap"`        // enable text wrapping
    LineHeight      float64     `json:"lineHeight"`      // line height multiplier (1.0 = font size)
    
    // Processing
    PreProcessors   []interface{} `json:"preProcessors"`
    PostProcessors  []interface{} `json:"postProcessors"`
}
```

### Config Hierarchy

Load order (first found wins):
1. Command-line flag (`--config`/`-c`)
2. Local config: `./.txt2igconfig.jsonc`
3. XDG global: `$XDG_CONFIG_HOME/txt2ig/config.jsonc`
4. Home global: `~/.txt2ig/config.jsonc`
5. Built-in defaults

### Default Values

```jsonc
{
    "font": ["GoMono"],
    "fontSize": 18,
    "fontColor": "#FFFFFF",
    "bgColor": "#000000",
    "textJustify": "left",
    "textBoxJustify": "center",
    "textBoxAlign": "center",
    "textBoxOffset": [0, 0],
    "textBoxMaxWidth": 972,  // 90% of 1080
    "screenSize": [1080, 1920],
    "textWrap": true,
    "lineHeight": 1.4,
    "preProcessors": [],
    "postProcessors": []
}
```

---

## Technical Decisions

### 1. CLI Framework: `alecthomas/kong`

**Why Kong over Cobra?**
- Simpler, less boilerplate
- Declarative via struct tags
- Produces clean, testable code
- Perfect "grug brain" alignment

**Example:**
```go
var cli struct {
    File   string `arg:"" name:"file" help:"Text file to convert" type:"existingfile"`
    Output string `short:"o" help:"Output file name (.jpg or .png)" type:"path"`
    Config string `short:"c" long:"config" help:"Custom config file" type:"existingfile"`
}

func main() {
    ctx := kong.Parse(&cli)
    // ... process
}
```

### 2. Font Loading: Embedded GoMono as Fallback

**Why GoMono?**
- Embedded in `golang.org/x/image/font/gofont/gomono`
- ~173KB embedded in binary
- BSD-3-Clause license (very permissive)
- Zero external dependencies
- Clean monospace font designed for code/text

**Font Resolution Order:**
```
For each font in user's font array (e.g., ["FiraMono", "/custom.ttf", "mono"]):
  1. Try loading as file path (if exists on filesystem)
  2. Try loading as system font name (lookup in system font directories)
  3. If successful → use this font
  4. If failed → continue to next font in array

If ALL fonts in array fail:
  5. Final fallback: embedded GoMono
```

**Implementation:**
```go
import (
    "golang.org/x/image/font/gofont/gomono"
    "golang.org/x/image/font/opentype"
)

func LoadFont(fontSpec string) (*opentype.Font, error) {
    // 1. Try as file path
    if _, err := os.Stat(fontSpec); err == nil {
        data, err := os.ReadFile(fontSpec)
        if err != nil {
            return nil, fmt.Errorf("read font file %s: %w", fontSpec, err)
        }
        return opentype.Parse(data)
    }
    
    // 2. Try as system font (implementation needed)
    // ...
    
    return nil, fmt.Errorf("font not found: %s", fontSpec)
}

func LoadFontWithFallback(fonts []string) (*opentype.Font, error) {
    for _, font := range fonts {
        if f, err := LoadFont(font); err == nil {
            return f, nil
        }
    }
    
    // Final fallback: GoMono
    return opentype.Parse(gomono.TTF)
}
```

### 3. JSONC Parsing: Strip Comments, Then Parse JSON

**Algorithm:**
1. Remove `//` single-line comments
2. Remove `/* */` multi-line comments
3. Handle strings properly (don't strip comments inside strings)
4. Parse remaining JSON

**Implementation:**
```go
func ParseJSONC(data []byte, v interface{}) error {
    // Strip comments
    cleaned := stripComments(data)
    
    // Parse as standard JSON
    return json.Unmarshal(cleaned, v)
}

func stripComments(data []byte) []byte {
    // Implementation handles:
    // - Single-line comments: //
    // - Multi-line comments: /* */
    // - String literals (preserves content)
    // ...
}
```

### 4. Text Wrapping Algorithm

**Inputs:**
- `text`: Input text (may contain newlines)
- `maxWidth`: Maximum width in pixels
- `face`: Font face for measuring
- `lineHeight`: Multiplier (default 1.4)

**Algorithm:**
```
For each line in input text:
  If textWrap is disabled:
    - Render line as-is
  Else:
    - Split line into words
    - Initialize current line as empty
    - For each word:
      - If word fits in current line (measure width):
        - Add word to current line
      - Else:
        - Embark current line (if not empty)
        - Start new line with this word
    - Emit final line
  After each line:
    - Apply lineHeight spacing
```

**Implementation Note:**
- Use `font.Drawer.MeasureString()` to measure text width
- Respect existing newlines in input (don't merge paragraphs)

### 5. Image Output: PNG and JPEG Support

**Detection:**
```go
func SaveImage(img *image.RGBA, outputPath string) error {
    ext := strings.ToLower(filepath.Ext(outputPath))
    
    f, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("create output file %s: %w", outputPath, err)
    }
    defer f.Close()
    
    switch ext {
    case ".png":
        return png.Encode(f, img)
    case ".jpg", ".jpeg":
        return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
    default:
        return fmt.Errorf("unsupported output format: %s (use .jpg or .png)", ext)
    }
}
```

---

## Rendering Pipeline

### Stage 1: Load Input
```
1. Parse CLI args (Kong)
2. Read input text file
3. Load config (hierarchical: custom > local > global > defaults)
4. Validate config
```

### Stage 2: Pre-Process Text
```
1. Run pre-processors in sequence
2. Each processor transforms text:
   Input: "Hello @date"
   Output: "Hello 2026-04-08"
```

### Stage 3: Load Fonts
```
1. Parse font array from config
2. Try each font in order (file path → system font)
3. Fall back to GoMono if all fail
4. Create font face with specified size
```

### Stage 4: Calculate Layout
```
1. Determine text box dimensions
   - maxWidth = textBoxMaxWidth (or 90% of screen width)
   - Calculate wrapped text
2. Determine text box position
   - Apply textBoxJustify (horizontal)
   - Apply textBoxAlign (vertical)
   - Apply textBoxOffset
3. Calculate text anchor point
```

### Stage 5: Create Image
```
1. Create RGBA image with screenSize
2. Fill background with bgColor
3. Draw text at calculated position
   - Apply textJustify within text box
   - Apply lineHeight for spacing
```

### Stage 6: Post-Process Styling
```
1. Run post-processors in sequence
2. Each processor modifies styling:
   Input: "#" headers → Bold style
```

### Stage 7: Save Output
```
1. Detect output format from extension
2. Encode as JPG or PNG
3. Write to file
```

---

## Processor System

### Interface

```go
type PreProcessor interface {
    Process(text string) (string, error)
    Name() string
}

type PostProcessor interface {
    Process(line string) (string, *TextStyle, error)
    Name() string
}
```

### Built-In Pre-Processors

#### 1. `exactSearchAndReplace`
```go
type ExactSearchAndReplace struct {
    SearchString  string `json:"searchString"`
    ReplaceString string `json:"replaceString"`
}

func (p *ExactSearchAndReplace) Process(text string) (string, error) {
    return strings.ReplaceAll(text, p.SearchString, p.ReplaceString), nil
}
```

#### 2. `exactSearchAndReplaceWithDateTimeNow`
```go
type ExactSearchAndReplaceWithDateTimeNow struct {
    SearchString string `json:"searchString"`
    ReplaceFormat string `json:"replaceFormat"` // yyyy-mm-dd, etc.
}

func (p *ExactSearchAndReplaceWithDateTimeNow) Process(text string) (string, error) {
    now := time.Now()
    formatted := formatDateTime(now, p.ReplaceFormat)
    return strings.ReplaceAll(text, p.SearchString, formatted), nil
}
```

### Built-In Post-Processors

#### 1. `markdown-bold-headers`
```go
type MarkdownBoldHeaders struct {
    Bold bool `json:"bold"`        // true
    FontColor string `json:"fontColor"` // empty = use default, or hex color
}

func (p *MarkdownBoldHeaders) Process(line string) (string, *TextStyle, error) {
    if strings.HasPrefix(line, "#") {
        return strings.TrimLeft(line, "# "), &TextStyle{
            Bold: p.Bold,
            FontColor: p.FontColor, // empty string means default
        }, nil
    }
    return line, nil, nil
}
```

---

## Error Handling

### Style: Context-Aware, Helpful

**Good Examples:**
```go
// Config loading
return nil, fmt.Errorf("load config from %s: parse JSONC: %w", path, err)

// Font loading
return nil, fmt.Errorf("load font %s: file not found or not in system fonts", fontSpec)

// Text rendering
return nil, fmt.Errorf("render text: measure string width: %w", err)
```

**Principle:**
- Always include context about what operation failed
- Wrap underlying errors with `%w`
- Provide helpful hints in error messages

---

## Testing Strategy

### Unit Tests

**Pure Functions (Easy to Test):**
- JSONC parsing
- Text wrapping logic
- Processor transformations
- Layout calculations
- Color parsing

**Example:**
```go
func TestTextWrapping(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        maxWidth int
        expected []string
    }{
        {
            name: "simple wrap",
            input: "hello world",
            maxWidth: 60, // assume "hello world" is > 60px
            expected: []string{"hello", "world"},
        },
        // ...more cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := wrapText(tt.input, tt.maxWidth, mockFace)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

**Dependency Injection:**
```go
type FileSystem interface {
    ReadFile(path string) ([]byte, error)
    Stat(path string) (os.FileInfo, error)
}

type ConfigLoader struct {
    fs FileSystem
}

func NewConfigLoader() *ConfigLoader {
    return &ConfigLoader{fs: osFS{}}
}

// For testing:
func TestConfigLoading(t *testing.T) {
    loader := &ConfigLoader{fs: mockFS{
        files: map[string][]byte{
            "./.txt2igconfig.jsonc": []byte(`{"fontSize": 20}`),
        },
    }}
    // ...
}
```

### Integration Tests

**End-to-End:**
```go
func TestEndToEnd(t *testing.T) {
    input := "Hello World\n# Header"
    cfg := defaultConfig()
    
    img, err := Render(input, cfg)
    assert.NoError(t, err)
    assert.NotNil(t, img)
    
    // Save to temp file
    tmp := t.TempDir()
    outputPath := filepath.Join(tmp, "output.jpg")
    err = SaveImage(img, outputPath)
    assert.NoError(t, err)
    
    // Verify file exists
    _, err = os.Stat(outputPath)
    assert.NoError(t, err)
}
```

---

## Implementation Phases

### Phase 1: Foundation (Days 1-2)
**Tasks:**
1. Initialize Go module (`go mod init github.com/yourusername/txt2ig`)
2. Create directory structure (`internal/...`)
3. Create root `main.go` with basic Kong CLI
4. Implement JSONC parser (`internal/config/jsonc.go`)
5. Define config types (`internal/config/config.go`)
6. Implement config loader with hierarchy (`internal/config/loader.go`)

**Deliverables:**
- Project scaffolding complete
- Can parse JSONC config files
- Can load config from hierarchical locations
- Unit tests for JSONC parsing

---

### Phase 2: Core Rendering (Days 3-5)
**Tasks:**
1. Implement font manager (`internal/font/`)
   - Load fonts from file path
   - Load fonts from system (basic lookup)
   - Fallback chain logic
   - Embedded GoMono
2. Create basic image generation (`internal/renderer/image.go`)
   - Create canvas
   - Fill background color
   - Save as PNG/JPEG
3. Implement text rendering (`internal/renderer/text.go`)
   - Basic text drawing
   - Text justification (left/center/right)
   - Text box alignment (top/center/bottom)
4. Implement text wrapping (`internal/renderer/text.go`)
   - Word-level wrapping
   - Line height calculation
5. Coordinate system (`internal/renderer/renderer.go`)
   - Map text box position to screen coordinates

**Deliverables:**
- Can render basic text to image
- Can save as JPG or PNG
- Font loading works with fallback chain
- Text wrapping functional

---

### Phase 3: Styling & Layout (Days 6-7)
**Tasks:**
1. Refine text box calculations
   - Max width constraints
   - Offset handling
2. Color parsing (hex string to color.RGBA)
3. Font color and background color
4. Test different screen sizes (1080x1920, etc.)

**Deliverables:**
- Text positioned correctly
- Colors working
- Layout matches config

---

### Phase 4: Processors (Days 8-9)
**Tasks:**
1. Define processor interfaces (`internal/processor/processor.go`)
2. Implement processor registration system
3. Implement pre-processors:
   - `exactSearchAndReplace`
   - `exactSearchAndReplaceWithDateTimeNow`
4. Implement post-processors:
   - `markdown-bold-headers`
5. Integration: Run processors before/after rendering

**Deliverables:**
- Processors working end-to-end
- Can chain multiple processors
- Unit tests for each processor

---

### Phase 5: Polish & Testing (Days 10-12)
**Tasks:**
1. Comprehensive unit tests
   - Aim for >80% coverage
   - Table-driven tests
   - Mock filesystem for tests
2. Integration tests
   - End-to-end: text file → config → image
   - Visual verification of sample outputs
3. Add example test fixtures (`internal/testdata/`)
4. Update README
   - Add installation instructions
   - Add usage examples
   - Document all config fields
   - Mention PNG support
5. Add to README: Build/install instructions

**Deliverables:**
- Test coverage >80%
- README complete
- Example configurations provided
- Ready for use

---

## Future Enhancements (Post v1)

### Hot-Reload Preview
- Watch file for changes
- Spin up local web server
- Preview in browser with canvas
- Real-time updates

### Plugin System
**Why postpone?**
- Go's native `plugin` package has limitations:
  - Only on Linux/macOS/FreeBSD
  - Requires exact Go version match
  - Complex deployment

**Alternative approaches:**
- RPC-based plugins (gRPC)
- Lua scripting
- WASM
- Simple Go interface-based extensions compiled into binary

---

## Questions & Decisions Log

### Q1: Font File Paths vs Font Names?
**Decision**: Support both.
- Font names: Try to find in system font directories
- File paths: Load directly from filesystem
- Order: Try file path first, then system lookup

### Q2: Default Font Behavior?
**Decision**: Embed GoMono as fallback.
- Embedded in binary (zero external deps)
- Used only after exhausting user's fallback chain
- Keeps binary portable

### Q3: Output Format Support?
**Decision**: Support both PNG and JPEG.
- Detect from output file extension
- jpeg: quality 95
- png: best compression

### Q4: Text Wrapping Default?
**Decision**: Enable by default (`textWrap: true`).
- Most users want this
- Can be disabled if needed
- Configurable max width (default 90% of screen)

### Q5: Line Height Default?
**Decision**: `lineHeight: 1.4` default.
- Good balance for monospace fonts
- Configurable per user preference
- Multiplier of font size

### Q6: Error Handling Style?
**Decision**: Fail-fast with context.
- Clear error messages
- Wrap errors with `%w`
- Include operation context
- Exit with non-zero status

### Q7: Padding vs TextBoxOffset?
**Decision**: No separate padding config.
- `textBoxOffset` provides positioning control
- `textBoxMaxWidth` controls wrap width
- Users can achieve padding effect with offset + max width

---

## Development Commands

```bash
# Initialize project
cd /path/to/txt2ig
go mod init github.com/yourusername/txt2ig

# Add dependencies
go get github.com/alecthomas/kong@latest
go get golang.org/x/image@latest

# Run
go run main.go input.md -o output.jpg

# Build
go build -o txt2ig

# Test
go test ./...

# Test with coverage
go test -cover ./...

# Install
go install
```

---

## Notes for Future Sessions

### Context forFresh Implementation
This plan contains all decisions, trade-offs, and implementation details needed to implement txt2ig from scratch. Key points:

1. **Architecture**: All code in `internal/`, root `main.go` entry point, no `/cmd`
2. **CLI**: Use `alecthomas/kong` - declarative, simple
3. **Font Strategy**: Fallback chain → GoMono as last resort
4. **Config**: JSONC parsing, hierarchical loading
5. **Image**: Use Go standard library + golang.org/x/image for fonts
6. **Testing**: Pure functions + dependency injection
7. **Philosophy**: Simple, minimal, "grug brain"

### Implementation Order
1. Project structure + CLI
2. Config loading
3. Font management
4. Basic rendering
5. Text wrapping + layout
6. Processors
7. Tests + polish

### Key Files to Start With
- `main.go` - Entry point
- `internal/cli/cli.go` - Kong CLI definition
- `internal/config/config.go` - Config struct
- `internal/config/jsonc.go` - JSONC parser
- `internal/font/font.go` - Font loading

### Don't Forget
- Update README after implementation
- Fix JSONC typos in README (missing commas)
- Test with both JPG and PNG outputs
- Test font fallback chain
- Test config hierarchy
- Add example configs in README

---

**Last Updated**: 2026-04-08
**Status**: Ready for implementation
**Estimated Time**: 3-4 hours total work