package main

import (
	"fmt"
	"github.com/MattSPalmer/svgraphs"
	"github.com/ajstarks/svgo"
	"os"
	"time"
)

var colors = []string{"red", "orange", "green", "teal", "blue", "purple", "maroon"}

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

	rh := int((float64(h) - (10*float64(len(m)) + 60)) / float64(len(m)))
	rh = rh - (rh % 30)
	nextColor := colorer()

	canvas.Grid(0, 0, w, h, 30, "stroke-width:3;fill:black;opacity:10%")

	var i int
	for _, blocks := range m {
		style := fmt.Sprintf("fill:%v", nextColor())
		for _, b := range blocks {
			start := w * (timeToMins(b.start) - 540) / 660
			pixelW := w * int(b.Duration().Minutes()) / 660
			canvas.Rect(start, i, pixelW, rh, style)
		}
		i += rh + 30
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

func colorer() func() string {
	i := -1
	return func() string {
		i += 1
		return colors[i%len(colors)]
	}
}

func callDensity(calls []CallRecord) []int {
	density := make([]int, 1440)
	for _, call := range calls {
		a := timeToMins(call.Interval().start)
		b := timeToMins(call.Interval().end)
		fmt.Printf("%v, %v\n", a, b)
		for i := a; i < b; i++ {
			density[i] += 1
		}
	}
	return density[539:1200]
}

func timeToMins(t time.Time) int { return 60*t.Hour() + t.Minute() }
func minsToTime(i int) time.Time {
	return time.Date(0, time.January, 1, (i%60)/60, i-(i%60), 0, 0, time.UTC)
}
