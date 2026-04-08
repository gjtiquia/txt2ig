package config

import "fmt"

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
	cfg := &Config{}
	if err := ParseJSONC(defaultConfig, cfg); err != nil {
		panic(fmt.Errorf("parse embedded default config: %w", err))
	}
	return cfg
}
