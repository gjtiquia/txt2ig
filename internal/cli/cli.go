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
}

type WebCmd struct {
	Port  int    `short:"p" long:"port" help:"Port to run server on" default:"3000"`
	Watch string `short:"w" long:"watch" help:"File to watch for changes" type:"existingfile"`
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
