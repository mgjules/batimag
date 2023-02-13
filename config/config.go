package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/go-playground/validator"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v2"
)

// Config represents the configuration.
type Config struct {
	InputDir  string    `yaml:"input_dir" validate:"required,dir"`
	OutputDir string    `yaml:"output_dir" validate:"required,dir"`
	Resize    Resize    `yaml:"resize"`
	Transform Transform `yaml:"transform"`
	Effect    Effect    `yaml:"effect"`
	Adjust    Adjust    `yaml:"adjust"`
}

// Resize represents a bunch of opinionated image resize options.
type Resize struct {
	// Desired width in pixels.
	// Zero value will resize image proportionally to height.
	// Default value is 0.
	Width uint16 `yaml:"width" validate:"gte=0,lte=7680"`
	// Desired height in pixels.
	// Zero value will resize image proportionally to width.
	// Default value is 0.
	Height uint16 `yaml:"height" validate:"gte=0,lte=7680"`
	// Resize type.
	// Accepted values are Normal, Thumbnail, Fit and Fill.
	// Default value is Normal.
	Type string `yaml:"type" default:"Normal" validate:"oneof=Normal Thumbnail Fit Fill"`
	// Resize resampling filter.
	// Accepted values are Lanczos, CatmullRom, MitchellNetravali, Linear, Box and NearestNeighbor.
	// Default value is Lanczos.
	Filter string `yaml:"filter" default:"Lanczos" validate:"oneof=Lanczos CatmullRom MitchellNetravali Linear Box NearestNeighbor"`
	// Anchor is used for image alignment when using Resize Type 'Fill'.
	// Accepted values are Center, TopLeft, Top, TopRight, Left, Right, BottomLeft, Bottom and BottomRight.
	// Default value is Center.
	Anchor string `yaml:"anchor" default:"Center" validate:"oneof=Center TopLeft Top TopRight Left Right BottomLeft Bottom BottomRight"`
}

// Transform represents a bunch of opinionated image transform options.
type Transform struct {
	// Automatically fix orientation after processing.
	// Default value is false.
	AutoOrientation bool `yaml:"auto_orientation"`
	// Rotate clockwise by certain angle in degrees.
	// Accepted values are 0 90 180 270.
	// Default value is 0.
	Rotate uint16 `yaml:"rotate" validate:"oneof=0 90 180 270"`
	// Flip on a certain axis.
	// Accepted values are Vertical and Horizontal.
	// Default value is None.
	Flip string `yaml:"flip" default:"None" validate:"oneof=None Vertical Horizontal"`
}

// Effect represents a bunch of opinionated image effect options.
type Effect struct {
	// Sigma parameter allows to control the strength of the blurring effect.
	// Accepted values are 0.0 to 5.0.
	// Default value is 0.
	GaussianBlur float64 `yaml:"gaussian_blur" validate:"gte=0.0,lte=5.0"`
	// Sharpen uses gaussian function internally.
	// Accepted values are 0.0 to 5.0.
	// Default value is 0.
	Sharpen float64 `yaml:"sharpen" validate:"gte=0.0,lte=5.0"`
}

// Adjust represents a bunch of opinionated image adjust options.
type Adjust struct {
	// Gamma correction.
	// Accepted values are 0.0 to 5.0.
	// Values less that 1.0 are darker while values greater than 1.0 are lighter.
	// Default value is 0.
	Gamma float64 `yaml:"gamma" validate:"gte=0.0,lte=5.0"`
	// Percentage value.
	// Accepted values are -100 to 100.
	// Negative values are grayer.
	// Default value is 0.
	Contrast float64 `yaml:"contrast" validate:"gte=-100,lte=100"`
	// Percentage value.
	// Accepted values are -100 to 100.
	// Lower negative values are darker while higher positive values are lighter.
	// Default value is 0.
	Brightness float64 `yaml:"brightness" validate:"gte=-100,lte=100"`
	// Percentage value.
	// Accepted values are -100 to 100.
	// Negative values decrease image saturation while positive values increase image saturation.
	// Default value is 0.
	Saturation float64 `yaml:"saturation" validate:"gte=-100,lte=100"`
	// Grayscale version.
	// Default value is false.
	Grayscale bool `yaml:"grayscale"`
}

// Load loads the configuration for a given path.
func Load(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	defaults.SetDefaults(&cfg)

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &cfg, nil
}
