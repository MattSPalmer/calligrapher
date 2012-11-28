package main

import (
	"github.com/MattSPalmer/svgraphs"
	"github.com/ajstarks/svgo"
	"os"
)

func WriteToSVG(cg CallGraph, fp string) error {
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
