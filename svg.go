package main

import (
	"github.com/MattSPalmer/svgraphs"
	"github.com/ajstarks/svgo"
	"os"
)

func barGraph(cg CallGraph, fp string) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	dist, err := cg.Distribution()
	if err != nil {
		return err
	}
	output := make(map[string]int)

	for k, v := range dist {
		key := cg.Labels(k)
		output[key] = len(v)
	}
	canvas := svg.New(f)
	g, err := graphs.NewBarGraph(output)
	if err != nil {
		return err
	}
	canvas.Start(700, 500)
	defer canvas.End()

	g.Draw(canvas, 30, 30, 600, 400)

	return nil
}

func durationStackGraph(w, h int, m map[string][]timeBlock, fp string) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	canvas := svg.New(f)
	canvas.Start(w, h)
	defer canvas.End()

	rh := (h - (10*len(m) + 60)) / len(m)
	i := 30
	for _, blocks := range m {
		for _, b := range blocks {
			start := w * (60*b.start.Hour() + b.start.Minute() - 540) / 660
			pixelW := w * int(b.Duration().Minutes()) / 660
			canvas.Rect(start, i-rh/2, pixelW, rh)
		}
		i += rh + 10
	}
	return nil
}

func mapCallTime(calls []CallRecord) (map[string][]timeBlock, error) {
	m := make(map[string][]timeBlock)
	dist, err := GraphByAgent(calls).Distribution()
	if err != nil {
		return nil, err
	}
	for k, v := range dist {
		if k == -1 {
			continue
		}
		b := make([]timeBlock, 0)
		for _, call := range v {
			b = append(b, *call.Interval())
		}
		name := agentsByID[k]
		m[name] = b
	}
	return m, nil
}
