# txt2ig

![test-bash](https://github.com/user-attachments/assets/fc598510-7b30-4fd4-a9ab-827815d34a2c)

a cli tool to convert text files to images with plain text, ready for posting to instagram

meant to be dead simple, intuitive config

> not to be confused with `txt2img` for ai image gen

## installation

### go install (recommended)

```bash
go install github.com/gjtiquia/txt2ig@latest
```

<details>
<summary>other useful commands:</summary>

```bash
# checks what is the latest available version on go proxy cache
go list -m github.com/gjtiquia/txt2ig@latest

# checks what is the latest version directly from GitHub
GOPROXY=direct go list -m github.com/gjtiquia/txt2ig@latest

# installs latest version directly from GitHub
GOPROXY=direct go install github.com/gjtiquia/txt2ig@latest

# installs binary at current directory instead of a global install
GOBIN=$(pwd) go install github.com/gjtiquia/txt2ig@latest
```

## features

- ✅ **text to image**: convert plain text files to images
- ✅ **multiple formats**: output as JPG or PNG (auto-detected from file extension)
- ✅ **configurable**: customize appearance with JSONC config files
- ✅ **font fallback**: supports fonts with fallback chain, embedded GoMono as final fallback
- ✅ **text wrapping**: automatic text wrapping with configurable width
- ✅ **newline preservation**: respects newlines in input, preserves paragraph structure
- ✅ **processors**: pre-processors for text transformation, post-processors for styling
- ✅ **simple CLI**: minimal flags, helpful defaults
- ✅ **watch mode**: live preview with automatic regeneration on file save

## usage

```bash
# basic usage
# searches for my-post.md on the same directory
# if exists, creates my-post.jpg (supports .jpg and .png)
txt2ig my-post.md

# output flag
# set the output file name
# (supports .jpg and .png extensions)
txt2ig my-post.md -o another-name.png

# supports any plain text file
# supports relative and absolute paths
txt2ig ./src/post.txt -o ~/Downloads/img.jpg

# supports custom config
txt2ig post.md --config ./custom-config.jsonc
txt2ig post.md -c ./custom-config.jsonc

# debug currently used config
txt2ig --debug

# create a default config file for quick customization
txt2ig init
# creates .txt2igconfig.jsonc in current directory

# init with custom output path
txt2ig init -o my-config.jsonc

# init with force overwrite
txt2ig init --force

# and of cuz, a helpful help menu
txt2ig -help
txt2ig -h

# watch mode: regenerate image on file save
# watches post.md and regenerates post.jpg every time you save
txt2ig post.md -w
txt2ig post.md --watch

# watch mode with web preview
# watches post.md, regenerates post.jpg, and shows live preview at localhost:3000
txt2ig post.md -w --port 3000
txt2ig post.md --watch -p 3000

# watch mode with custom output and config
txt2ig post.md -w -o output.png -c custom.jsonc
```

## web server

Start a web server for converting text to images:

```bash
# start on port 3000 (default)
txt2ig web

# start on custom port
txt2ig web --port 8080

# short form
txt2ig web -p 8080
```

Features:
- Mobile-friendly interface
- Auto-save config and text to LocalStorage
- Preview image before download
- Real-time validation
- Download image as PNG

Note: If the port is already in use, the server will fail with an error message.

## watch mode

Watch a file for changes and automatically regenerate the image:

```bash
# basic watch mode
# watches post.md and regenerates post.jpg on every save
txt2ig post.md -w

# watch mode with live web preview
# shows live preview in browser at http://localhost:3000
txt2ig post.md -w --port 3000

# custom output file
txt2ig post.md -w -o custom.png

# custom config
txt2ig post.md -w -c custom.jsonc
```

Behavior:
- **Default format**: JPG (e.g., `post.md` → `post.jpg`)
- **Custom format**: Specify `-o custom.png` to use PNG
- **Live preview**: With `--port`, opens web preview that updates automatically
- **File regeneration**: Always regenerates image file on save, even with web preview
- **Port conflicts**: Fails with error if port is already in use
- **Stop**: Press `Ctrl+C` to stop watching

Web preview features:
- Live-updating image
- Config file name display
- Connection status indicator
- Download button
- Manual browser open (no auto-launch)

## config

config is a simple jsonc (JSON with Comments) file

### quick start

```bash
# create a default config file to customize
txt2ig init
# creates .txt2igconfig.jsonc in current directory

# or specify a custom path
txt2ig init -o my-config.jsonc
```

### config location

txt2ig will look for config file in the following order
- custom config (`-c`/`--config`)
- local: `./.txt2igconfig.jsonc`
- global: `XDGCONFIG/txt2ig/config.jsonc`
- global: `~/.txt2ig/config.jsonc`
- use defaults

### config params

the following are the default config params,
feel free to copy and paste this to your own config and override what you need

**note**: newlines in your input file are always preserved. text wrapping happens within each line/paragraph separately. empty lines create spacing between paragraphs.

```jsonc
{
    // font family configuration with bold/italic variants
    // supports font names and file paths with fallback chain
    // embedded GoMono variants available: GoMono, GoMonoBold, GoMonoItalic, GoMonoBoldItalic
    "fontFamily": {
        "regular": ["GoMono"],
        "bold": ["GoMonoBold"],
        "italic": ["GoMonoItalic"],
        "boldItalic": ["GoMonoBoldItalic"]
    },
    // unit: px
    "fontSize": 32,
    "fontColor": "#FFFFFF",
    "bgColor": "#000000",
    // text x-axis position within the text box
    "textJustify": "left", // "center", "right"
    // text box (bounding box of text) x-axis position on screen
    "textBoxJustify": "center", // "left", "right"
    // text box (bounding box of text) y-axis position on screen
    "textBoxAlign": "center", // "top", "bottom"
    // text box (bounding box of text) (x, y) offset. unit: px
    "textBoxOffset": [0, 0],
    // maximum width for text box. unit: px. 0 = auto (90% of screen width)
    "textBoxMaxWidth": 972,
    // screen/canvas size. [width, height]. unit: px
    "screenSize": [1080, 1920],
    // enable automatic text wrapping
    "textWrap": true,
    // line height multiplier (1.4 = 1.4x font size)
    "lineHeight": 1.4,

    // processors will run in sequence, 
    // you may chain several processors of the same name to get different results if you so desire

    // typically pre-process text
    "preProcessors": [
        // {
        //     "exactSearchAndReplace": {
        //         "searchString": "apple", // exact match only
        //         "replaceString": "bananas",
        //     }
        // }

        // {
        //     "grepSearchAndReplace": {
        //         "pattern": "^@foo",
        //         "replaceString": "bar",
        //     }
        // }

        // {
        //     "exactSearchAndReplaceWithDateTimeNow": {
        //         "searchString": "@date",
        //         "replaceFormat": "yyyy-mm-dd", // supports yyyy, mm, dd, hh, mm, ss
        //     }
        // }
    ],

    // typically post-process styling
    "postProcessors": [
        {
            // lines starting with # will be bold
            "markdown-bold-headers": {
                "bold": true,
                // set to empty to use the same default color
                // "fontColor": "#EC9006" // orange
            }
        },

        // {
        //      // lines starting with # will be italic and different color
        //     "bash-comments": {
        //         "italic": true,
        //         "fontColor": "#CCCCCC" // gray
        //     }
        // }

        {
            // highlights bash code blocks with Chroma syntax highlighting
            "bash-code-highlighting": {
                "style": "monokai", // Chroma style (monokai, dracula, github, etc.)
                "defaultColor": "#FFFFFF" // fallback color when Chroma doesn't provide one
            }
        },
    ],
}
```

## examples

### basic example

create a text file `post.md`:

```md
# My First Post

This is my first post!

Features:
- Simple text
- Clean output
- Ready for Instagram
```

run the tool:

```bash
txt2ig post.md
# creates post.jpg in the same directory
```

### with custom colors

```jsonc
{
    "fontFamily": {
        "regular": ["GoMono"],
        "bold": ["GoMonoBold"],
        "italic": ["GoMonoItalic"],
        "boldItalic": ["GoMonoBoldItalic"]
    },
    "fontSize": 20,
    "bgColor": "#1E1E1E",
    "fontColor": "#00FF00",
    "screenSize": [1080, 1920]
}
```

### with text wrapping disabled

```jsonc
{
    "textWrap": false,
    "fontSize": 16
}
```

## development

### run tests

```bash
# run all tests
go test ./...

# run with coverage
go test -cover ./...

# run specific package
go test ./internal/renderer -v
```

### build

```bash
# build for current platform
go build -o txt2ig

# build for specific platform
GOOS=linux GOARCH=amd64 go build -o txt2ig
GOOS=darwin GOARCH=amd64 go build -o txt2ig
GOOS=windows GOARCH=amd64 go build -o txt2ig.exe
```

## future roadmap

### plugin system

a way to build 3rd party plugins in addition to the official plugins available
- pre-processors
- post-processors

perhaps look into go plugins...? see how gonotify does it

## tech stack

- golang

## ai-usage disclosure

ai is heavily used for generating code in this project

continuing the ai workflow exploration after [ifg](https://github.com/gjtiquia/ifg)

some thoughts
- golang is still a really great language to vibecode with!
- Opencode Go GLM 5 is apparently pretty good at doing long running tasks and thorough research!
- a really great workflow i find to work well so far is
  - create a solid README.md first, as if the tool already exists
  - this helps you scope out the project and hv a vision for the final product first
  - also provides an anchor for the LLM to develop on (kinda like GDD for game devs)

(inspired by [ghostty's ai usage policy](https://github.com/ghostty-org/ghostty/blob/main/AI_POLICY.md))

## license

MIT
