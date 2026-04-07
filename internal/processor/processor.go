package processor

type PreProcessor interface {
	Process(text string) (string, error)
	Name() string
}

type PostProcessor interface {
	Process(line string) (string, *TextStyle, error)
	Name() string
}

type TextStyle struct {
	Bold      bool
	Italic    bool
	Underline bool
	FontColor string
	Size      *int
}

type ProcessorRegistry struct {
	preProcessors  map[string]PreProcessor
	postProcessors map[string]PostProcessor
}

func NewProcessorRegistry() *ProcessorRegistry {
	return &ProcessorRegistry{
		preProcessors:  make(map[string]PreProcessor),
		postProcessors: make(map[string]PostProcessor),
	}
}

func (r *ProcessorRegistry) RegisterPreProcessor(name string, processor PreProcessor) {
	r.preProcessors[name] = processor
}

func (r *ProcessorRegistry) RegisterPostProcessor(name string, processor PostProcessor) {
	r.postProcessors[name] = processor
}

func (r *ProcessorRegistry) GetPreProcessor(name string) (PreProcessor, bool) {
	p, ok := r.preProcessors[name]
	return p, ok
}

func (r *ProcessorRegistry) GetPostProcessor(name string) (PostProcessor, bool) {
	p, ok := r.postProcessors[name]
	return p, ok
}

func ApplyPreProcessors(text string, processors []interface{}, registry *ProcessorRegistry) (string, error) {
	result := text
	for range processors {
		// TODO: Parse processor config and apply
		// For now, return the text as-is
	}
	return result, nil
}

func ApplyPostProcessors(lines []string, processors []interface{}, registry *ProcessorRegistry) ([]string, error) {
	// TODO: Implement post-processor application
	return lines, nil
}
