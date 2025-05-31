package benchmark

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// GenerateGraphs creates performance graphs from the recorded metrics
func (r *Recorder) GenerateGraphs(outputDir string) error {
	metrics := r.GetMetrics()
	if len(metrics) == 0 {
		return fmt.Errorf("no metrics recorded")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate pages crawled vs time graph
	if err := r.generatePagesGraph(outputDir, metrics); err != nil {
		return fmt.Errorf("failed to generate pages graph: %w", err)
	}

	// Generate crawled/queued ratio vs time graph
	if err := r.generateRatioGraph(outputDir, metrics); err != nil {
		return fmt.Errorf("failed to generate ratio graph: %w", err)
	}

	return nil
}

func (r *Recorder) generatePagesGraph(outputDir string, metrics []Metric) error {
	p := plot.New()
	p.Title.Text = "Pages Crawled vs Time"
	p.X.Label.Text = "Time (seconds)"
	p.Y.Label.Text = "Pages Crawled"

	pts := make(plotter.XYs, len(metrics))
	for i, m := range metrics {
		pts[i].X = m.Timestamp.Sub(r.start).Seconds()
		pts[i].Y = float64(m.PagesCount)
	}

	line, points, err := plotter.NewLinePoints(pts)
	if err != nil {
		return err
	}

	line.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	points.Shape = draw.CircleGlyph{}
	points.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}

	p.Add(line, points)
	p.Legend.Add("Pages", line)

	filename := filepath.Join(outputDir, "pages_vs_time.png")
	if err := p.Save(8*vg.Inch, 6*vg.Inch, filename); err != nil {
		return fmt.Errorf("failed to save pages graph: %w", err)
	}

	return nil
}

func (r *Recorder) generateRatioGraph(outputDir string, metrics []Metric) error {
	p := plot.New()
	p.Title.Text = "Crawled/Queued Ratio vs Time"
	p.X.Label.Text = "Time (seconds)"
	p.Y.Label.Text = "Ratio (Crawled/Queued)"

	pts := make(plotter.XYs, len(metrics))
	for i, m := range metrics {
		pts[i].X = m.Timestamp.Sub(r.start).Seconds()
		if m.QueuedCount > 0 {
			pts[i].Y = float64(m.PagesCount) / float64(m.QueuedCount)
		} else {
			pts[i].Y = 0
		}
	}

	line, points, err := plotter.NewLinePoints(pts)
	if err != nil {
		return err
	}

	line.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	points.Shape = draw.CircleGlyph{}
	points.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

	p.Add(line, points)
	p.Legend.Add("Ratio", line)

	filename := filepath.Join(outputDir, "ratio_vs_time.png")
	if err := p.Save(8*vg.Inch, 6*vg.Inch, filename); err != nil {
		return fmt.Errorf("failed to save ratio graph: %w", err)
	}

	return nil
}
