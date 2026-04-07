# txt2ig

a cli tool to convert text files to images with plain text, ready for posting to instagram

meant to be dead simple, minimal config

## usage

```bash
# basic usage
# searches for my-post.md on the same directory
# if exists, creates my-post.jpg (we only support .jpg for now)
txt2ig my-post.md

# output flag
# set the output file name
# (will error if output does not end in .jpg)
txt2ig my-post.md -o another-name.jpg

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

```jsonc
{
    // target font and fallback fonts
    "font": ["FiraMono Nerd Font Mono", "mono"],
    // unit: px
    "fontSize": 18,
    "fontColor": "#FFFFFF"
    "bgColor": "#000000"
    // text x-axis postion within the text box
    "textJustify": "left", // "center", "right"
    // text box (bounding box of text) x-axis position on screen
    "textBoxJustify": "center", // "left", "right"
    // text box (bounding box of text) y-axis position on screen
    "textBoxAlign": "center", // "top", "bottom"
    // text box (bounding box of text) (x, y) offset. unit: px
    "textBoxOffset": [0, 0],
    // (x, y). unit: px
    "screenSize": [1080. 1920],

    // processors will run in sequence, 
    // you may chain several processors of the same name to get different results if you so desire

    // typically pre-process text
    "preProcessors": [
        // "exactSearchAndReplace": {
        //     "searchString": "apple", // exact match only
        //     "replaceString": "bananas",
        // },

        // "grepSearchAndReplace": {
        //     "pattern": "^@foo",
        //     "replaceString": "bar",
        // }

        // "exactSearchAndReplaceWithDateTimeNow": {
        //     "searchString": "@date",
        //     "replaceFormat": "yyyy-mm-dd", // supports yyyy, mm, dd, hh, mm, ss
        // },
    ],

    // typically post-process styling
    "postProcessors": [
        // lines starting with # will be bold
        "markdown-bold-headers": {
            "bold": true,
            // set to empty to use the same default color
            // "fontColor": "#EC9006" // orange
        },

        // lines starting with # will be italic and different color
        // "bash-comments": {
        //     "italic": true,
        //     "fontColor": "#CCCCCC" // gray
        // },
    ],
}
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
