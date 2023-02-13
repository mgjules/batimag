# batimag

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/mgjules/batimag)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](LICENSE)

Batimag is short for Batch Imager which applies basic image processing functions to a batch of images in a specific directory recursively.

## Installation

### From Source (recommended)

```shell
go install github.com/mgjules/batimag@v0.1.0
```

### From Binary

Download the corresponding binary for your operating system and architecture from the [releases](https://github.com/mgjules/batimag/releases) page and place it in a `$PATH` directory.

## Usage

```shell
❯ batimag
NAME:
   batimag - Batch Image Processor!

USAGE:
   batimag [global options] command [command options] [arguments...]

DESCRIPTION:
   Batimag applies a set of image processing functions to images in a given directory recursively.

AUTHOR:
   Michaël Giovanni Jules <julesmichaelgiovanni@gmail.com>

COMMANDS:
   process     Process some images!
   version, v  Shows the version
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

COPYRIGHT:
   (c) 2023 Michaël Giovanni Jules
```

```shell
❯ batimag help process
NAME:
   batimag process - Process some images!

USAGE:
   batimag process [command options] [arguments...]

OPTIONS:
   --debug              run in debug mode (default: false)
   --clean              remove images in output directory (default: false)
   --config-path value  path to config file (default: config.yml)
   --help, -h           show help
```

### Example usage

Run batimag using the default configuration file `config.yml`.

```shell
batimag process --clean --debug
```

## Configuration

Example configuration (see [config.example.yml](config.example.yml)):

```yaml
input_dir: "./samples/input"
output_dir: "./samples/output"
resize:
  # Desired width in pixels.
  # Zero value will resize image proportionally to height.
  # Default value is 0.
  width: 100
  # Desired height in pixels.
  # Zero value will resize image proportionally to width.
  # Default value is 0.
  height: 50
  # Resize type.
  # Accepted values are Normal, Thumbnail, Fit and Fill.
  # Default value is Normal.
  type: Normal
  # Resize resampling filter.
  # Accepted values are Lanczos, CatmullRom, MitchellNetravali, Linear, Box and NearestNeighbor.
  # Default value is Lanczos.
  filter: Lanczos
  # Anchor is used for image alignment when using Resize Type 'Fill'.
  # Accepted values are Center, TopLeft, Top, TopRight, Left, Right, BottomLeft, Bottom and BottomRight.
  # Default value is Center.
  anchor: Center
transform:
  # Automatically fix orientation after processing.# Default value is false.
  auto_orientation: true
  # Rotate clockwise by certain angle in degrees.
  # Accepted values are 0 90 180 270.
  # Default value is 0.
  rotate: 90
  # Flip on a certain axis.
  # Accepted values are Vertical and Horizontal.
  # Default value is None.
  flip: Vertical
effect:
  # Strength of the gaussian blur effect.
  # Accepted values are 0.0 to 5.0.
  # Default value is 0.
  gaussian_blur: 2.0
  # Strength of the sharpening effect.
  # Accepted values are 0.0 to 5.0.
  # Default value is 0.
  sharpen: 1.0
adjust:
  # Gamma correction.
  # Accepted values are 0.0 to 5.0.
  # Values less that 1.0 are darker while values greater than 1.0 are lighter.
  # Default value is 0.
  gamma: 0.5
  # Percentage value.
  # Accepted values are -100 to 100.
  # Negative values are grayer.
  # Default value is 0.
  contrast: 50
  # Percentage value.
  # Accepted values are -100 to 100.
  # Lower negative values are darker while higher positive values are lighter.
  # Default value is 0.
  brightness: 50
  # Percentage value.
  # Accepted values are -100 to 100.
  # Negative values decrease image saturation while positive values increase image saturation.
  # Default value is 0.
  saturation: 50
  # Grayscale version.
  # Default value is false.
  grayscale: true
```

## Credits

- [github.com/disintegration/imaging](https://github.com/disintegration/imaging) - Image processing functions.
- [github.com/chai2010/webp](https://github.com/chai2010/webp) - Webp support.
- [github.com/sourcegraph/conc](https://github.com/sourcegraph/conc) - Concurrency made easy.
- [github.com/mcuadros/go-defaults](https://github.com/mcuadros/go-defaults) - Struct default values.
- [github.com/go-playground/validator/v10](https://github.com/go-playground/validator) - Struct validation.
- [github.com/urfave/cli/v2](https://github.com/urfave/cli) - CLI.
