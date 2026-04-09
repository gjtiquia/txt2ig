package cli

import "github.com/alecthomas/kong"

type CLI struct {
	Convert ConvertCmd `cmd:"" name:"convert" help:"Convert text file to image" default:"withargs"`
	Web     WebCmd     `cmd:"" name:"web" help:"Start web server"`
}

type ConvertCmd struct {
	InputFile string `arg:"" name:"file" help:"Text file to convert" type:"existingfile"`
	Output    string `short:"o" help:"Output file name (.jpg or .png)" type:"path"`
	Config    string `short:"c" long:"config" help:"Custom config file" type:"existingfile"`
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
