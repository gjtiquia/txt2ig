# txt2ig

a cli tool to convert text files to images with plain text, ready for posting to instagram

meant to be dead simple, minimal config

## installation

### from release

download the latest binary from [releases](https://github.com/gjtiquia/txt2ig/releases)

```bash
# make it executable
chmod +x txt2ig

# move to PATH (optional)
sudo mv txt2ig /usr/local/bin/
```

### from source

```bash
# clone the repository
git clone https://github.com/gjtiquia/txt2ig.git
cd txt2ig

# build
go build -o txt2ig

# install to $GOPATH/bin (optional)
go install
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

# and of cuz, a helpful help menu
txt2ig -help
txt2ig -h
```

## config

config is a simple jsonc (JSON with Comments) file

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
    "fontSize": 18,
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
        }

        // {
        //      // lines starting with # will be italic and different color
        //     "bash-comments": {
        //         "italic": true,
        //         "fontColor": "#CCCCCC" // gray
        //     }
        // }
    ],
}
```

## examples

### basic example

create a config file `.txt2igconfig.jsonc`:

```jsonc
{
    "fontFamily": {
        "regular": ["GoMono"],
        "bold": ["GoMonoBold"],
        "italic": ["GoMonoItalic"],
        "boldItalic": ["GoMonoBoldItalic"]
    },
    "fontSize": 24,
    "bgColor": "#000000",
    "fontColor": "#FFFFFF",
    "screenSize": [1080, 1920]
}
```

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

### hot-reload file watching preview support

- spins up a web server which watches for file changes
- previews image with canvas, as much as possible matching the exact image output without html drift

```bash
txt2ig -w post.md --port 3000
# -w or --watch
# -p or --port
```

### plugin system

a way to build 3rd party plugins in addition to the official plugins available
- pre-processors
- post-processors

perhaps look into go plugins...? see how gonotify does it

## tech stack

- golang

## license

MIT
