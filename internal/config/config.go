package config

type FontFamily struct {
	Regular    []string `json:"regular"`
	Bold       []string `json:"bold"`
	Italic     []string `json:"italic"`
	BoldItalic []string `json:"boldItalic"`
}

type Config struct {
	FontFamily      FontFamily    `json:"fontFamily"`
	FontSize        int           `json:"fontSize"`
	FontColor       string        `json:"fontColor"`
	BgColor         string        `json:"bgColor"`
	TextJustify     string        `json:"textJustify"`
	TextBoxJustify  string        `json:"textBoxJustify"`
	TextBoxAlign    string        `json:"textBoxAlign"`
	TextBoxOffset   []int         `json:"textBoxOffset"`
	TextBoxMaxWidth int           `json:"textBoxMaxWidth"`
	ScreenSize      []int         `json:"screenSize"`
	TextWrap        bool          `json:"textWrap"`
	LineHeight      float64       `json:"lineHeight"`
	PreProcessors   []interface{} `json:"preProcessors"`
	PostProcessors  []interface{} `json:"postProcessors"`
}

func DefaultConfig() *Config {
	return &Config{
		FontFamily: FontFamily{
			Regular:    []string{"GoMono"},
			Bold:       []string{"GoMonoBold"},
			Italic:     []string{"GoMonoItalic"},
			BoldItalic: []string{"GoMonoBoldItalic"},
		},
		FontSize:        18,
		FontColor:       "#FFFFFF",
		BgColor:         "#000000",
		TextJustify:     "left",
		TextBoxJustify:  "center",
		TextBoxAlign:    "center",
		TextBoxOffset:   []int{0, 0},
		TextBoxMaxWidth: 972,
		ScreenSize:      []int{1080, 1920},
		TextWrap:        true,
		LineHeight:      1.4,
		PreProcessors:   []interface{}{},
		PostProcessors:  []interface{}{},
	}
}
