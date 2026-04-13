package processor

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
