package cli

import "github.com/alecthomas/kong"

type CLI struct {
	Init    InitCmd    `cmd:"" help:"Create default config file"`
	Convert ConvertCmd `cmd:"" default:"withargs" help:"Convert text file to image"`
	Web     WebCmd     `cmd:"" help:"Start web server"`
}

type InitCmd struct {
	Output string `short:"o" long:"output" help:"Output file path" type:"path" default:".txt2igconfig.jsonc"`
	Force  bool   `short:"f" long:"force" help:"Overwrite existing file"`
}

type ConvertCmd struct {
	InputFile string `arg:"" name:"file" help:"Text file to convert" type:"existingfile" optional:""`
	Output    string `short:"o" help:"Output file name (.jpg or .png)" type:"path"`
	Config    string `short:"c" long:"config" help:"Custom config file" type:"existingfile"`
	Debug     bool   `short:"d" long:"debug" help:"Print config info and exit"`
	Watch     bool   `short:"w" long:"watch" help:"Watch file and regenerate on save"`
	Port      int    `short:"p" long:"port" help:"Port for web preview (requires --watch)"`
}

type WebCmd struct {
	Port int `short:"p" long:"port" default:"3000" help:"Port to run server on"`
}

func Parse(args []string) (*CLI, error) {
	var cli CLI
	parser, err := kong.New(&cli, kong.Name("txt2ig"), kong.Description("Convert text files to images for Instagram"))
	if err != nil {
		return nil, err
	}

	_, err = parser.Parse(args)
	if err != nil {
		return nil, err
	}

	return &cli, nil
}
