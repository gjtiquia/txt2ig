package processor

type PreProcessor interface {
	Process(text string) (string, error)
	Name() string
}

type PostProcessor interface {
	ProcessLines(lines []StyledLine) ([]StyledLine, error)
	Name() string
}
