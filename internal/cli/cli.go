package cli

import "github.com/alecthomas/kong"

type CLI struct {
	InputFile string `arg:"" name:"file" help:"Text file to convert" type:"existingfile"`
	Output    string `short:"o" help:"Output file name (.jpg or .png)" type:"path"`
	Config    string `short:"c" long:"config" help:"Custom config file" type:"existingfile"`
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
