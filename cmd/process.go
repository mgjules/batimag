package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/mgjules/batimag/config"
	"github.com/mgjules/batimag/logger"
	"github.com/sourcegraph/conc/pool"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	// Add webp support.
	_ "golang.org/x/image/webp"
)

var process = &cli.Command{
	Name:  "process",
	Usage: "Process some images!",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Value: false,
			Usage: "run in debug mode",
		},
		&cli.BoolFlag{
			Name:  "clean",
			Value: false,
			Usage: "remove images in output directory",
		},
		&cli.StringFlag{
			Name:        "config-path",
			Value:       "config.yml",
			DefaultText: "config.yml",
			Usage:       "path to config file",
		},
	},
	Action: func(c *cli.Context) error {
		cfg, err := config.Load(c.String("config-path"))
		if err != nil {
			return fmt.Errorf("config load: %w", err)
		}

		debug := c.Bool("debug")

		logger, err := logger.New(debug)
		if err != nil {
			return fmt.Errorf("new logger: %w", err)
		}

		logger.Debugw("Config loaded", "config", cfg)

		if c.Bool("clean") {
			if err = os.RemoveAll(cfg.OutputDir); err != nil {
				return fmt.Errorf("remove output dir: %w", err)
			}

			if err = os.MkdirAll(cfg.OutputDir, os.ModePerm); err != nil {
				return fmt.Errorf("create output dir: %w", err)
			}
		}

		logger.Debug("Processing images...")
		defer timer(logger.SugaredLogger, "Finished processing images")()

		pool := pool.New().WithMaxGoroutines(500)
		err = fs.WalkDir(os.DirFS(cfg.InputDir), ".", func(path string, d fs.DirEntry, err error) error {
			inputFilepath := filepath.Join(cfg.InputDir, path)

			log := logger.With("input path", inputFilepath)

			if err != nil {
				log.Errorw("Internal", "error", err)
				return nil
			}

			// Skip directories.
			if d.IsDir() {
				if err := os.MkdirAll(filepath.Join(cfg.OutputDir, path), os.ModePerm); err != nil {
					log.Errorw("Create output nested directory", "error", err)
					return fs.SkipDir
				}

				log.Debug("Skipping directory")
				return nil
			}

			pool.Go(func() {
				log.Debug("Processing image...")
				defer timer(log, "Finished processing image")()

				outputFilepath := filepath.Join(cfg.OutputDir, path)
				if fileExists(outputFilepath) {
					log.Debugw("Skipping already processed image", "output path", outputFilepath)
					return
				}

				img, err := imaging.Open(
					inputFilepath,
					imaging.AutoOrientation(cfg.Transform.AutoOrientation),
				)
				if err != nil {
					log.Errorw("Open image", "error", err)
					return
				}

				if cfg.Resize.Width > 0 || cfg.Resize.Height > 0 {
					var filter imaging.ResampleFilter
					switch cfg.Resize.Filter {
					case "Lanczos":
						filter = imaging.Lanczos
					case "CatmullRom":
						filter = imaging.CatmullRom
					case "MitchellNetravali":
						filter = imaging.MitchellNetravali
					case "Linear":
						filter = imaging.Linear
					case "Box":
						filter = imaging.Box
					case "NearestNeighbor":
						filter = imaging.NearestNeighbor
					default:
						log.Errorw("Unsupported resize filter", "filter", cfg.Resize.Filter)
						return
					}

					switch cfg.Resize.Type {
					case "Normal":
						img = imaging.Resize(img, int(cfg.Resize.Width), int(cfg.Resize.Height), filter)
					case "Thumbnail":
						img = imaging.Thumbnail(img, int(cfg.Resize.Width), int(cfg.Resize.Height), filter)
					case "Fit":
						img = imaging.Fit(img, int(cfg.Resize.Width), int(cfg.Resize.Height), filter)
					case "Fill":
						var anchor imaging.Anchor
						switch cfg.Resize.Anchor {
						case "Center":
							anchor = imaging.Center
						case "TopLeft":
							anchor = imaging.TopLeft
						case "Top":
							anchor = imaging.Top
						case "TopRight":
							anchor = imaging.TopRight
						case "Left":
							anchor = imaging.Left
						case "Right":
							anchor = imaging.Right
						case "BottomLeft":
							anchor = imaging.BottomLeft
						case "Bottom":
							anchor = imaging.Bottom
						case "BottomRight":
							anchor = imaging.BottomRight
						default:
							log.Errorw("Unsupported resize anchor", "anchor", cfg.Resize.Anchor)
							return
						}

						img = imaging.Fill(img, int(cfg.Resize.Width), int(cfg.Resize.Height), anchor, filter)
					default:
						log.Errorw("Unsupported resize type", "type", cfg.Resize.Type)
						return
					}
				}

				switch cfg.Transform.Rotate {
				case 0:
					img = imaging.Clone(img)
				case 90:
					img = imaging.Rotate90(img)
				case 180:
					img = imaging.Rotate180(img)
				case 270:
					img = imaging.Rotate270(img)
				default:
					log.Errorw("Unsupported rotate angle", "angle", cfg.Transform.Rotate)
					return
				}

				switch cfg.Transform.Flip {
				case "Vertical":
					img = imaging.FlipV(img)
				case "Horizontal":
					img = imaging.FlipH(img)
				case "None":
					// do nothing.
				default:
					log.Errorw("Unsupported flip axis", "axis", cfg.Transform.Flip)
					return
				}

				if cfg.Effect.GaussianBlur > 0 {
					img = imaging.Blur(img, cfg.Effect.GaussianBlur)
				}

				if cfg.Effect.Sharpen > 0 {
					img = imaging.Sharpen(img, cfg.Effect.Sharpen)
				}

				if cfg.Adjust.Gamma > 0 {
					img = imaging.AdjustGamma(img, cfg.Adjust.Gamma)
				}

				if cfg.Adjust.Contrast > 0 {
					img = imaging.AdjustContrast(img, cfg.Adjust.Contrast)
				}

				if cfg.Adjust.Brightness > 0 {
					img = imaging.AdjustBrightness(img, cfg.Adjust.Brightness)
				}

				if cfg.Adjust.Saturation > 0 {
					img = imaging.AdjustSaturation(img, cfg.Adjust.Saturation)
				}

				if cfg.Adjust.Grayscale {
					img = imaging.Grayscale(img)
				}

				ext := filepath.Ext(outputFilepath)
				if ext == ".webp" {
					f, err := os.Create(outputFilepath)
					if err != nil {
						log.Errorw("Create webp image", "error", err)
						return
					}
					defer f.Close()

					err = webp.Encode(f, img, &webp.Options{Lossless: true})
					if err != nil {
						log.Errorw("Encode webp image", "error", err)
						return
					}
				} else {
					err = imaging.Save(img, filepath.Join(cfg.OutputDir, path))
					if err != nil {
						log.Errorw("Save image", "error", err)
						return
					}
				}
			})

			return nil
		})
		if err != nil {
			return fmt.Errorf("transform images: %w", err)
		}

		pool.Wait()

		return nil
	},
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func timer(logger *zap.SugaredLogger, msg string) func() {
	start := time.Now()
	return func() {
		logger.Debugw(msg, "elasped", time.Since(start))
	}
}
