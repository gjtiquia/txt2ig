# example workflow

```bash
# set up environment
txt2ig init
touch example.md

# set up watch mode and live preview at port 3000
txt2ig example.md -w -p 3000

# edit post for content and config for styling
vim example.md
vim .txt2igconfig.jsonc

# generate example.jpg
txt2ig example.md
```
