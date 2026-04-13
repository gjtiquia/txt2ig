package processor

type PreProcessor interface {
	Process(text string) (string, error)
	Name() string
}

type PostProcessor interface {
	Process(line string) (string, *TextStyle, error)
	Name() string
}

type StatefulPostProcessor interface {
	PostProcessor
	ProcessLines(lines []StyledLine) ([]StyledLine, error)
}
