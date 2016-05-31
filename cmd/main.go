package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codeskyblue/go-sh"
	ui "github.com/gizak/termui"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	dataBuffer  = kingpin.Flag("data-buffer", "Length of data.").Short('b').Default("370").Int()
	processName = kingpin.Flag("process-name", "Name of process").Short('p').Required().String()
)

func main() {
	kingpin.Parse()
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	rss := make([]float64, *dataBuffer)

	lc1 := ui.NewLineChart()
	lc1.BorderLabel = "braille-mode Line Chart"
	lc1.Data = rss
	lc1.Width = 100
	lc1.Height = 20
	lc1.X = 0
	lc1.Y = 0
	lc1.AxesColor = ui.ColorWhite
	lc1.LineColor = ui.ColorYellow | ui.AttrBold

	par := ui.NewPar(":PRESS q TO QUIT DEMO")
	par.Height = 3
	par.Width = 100
	par.TextFgColor = ui.ColorWhite
	par.BorderLabel = "Current tracking process"
	par.BorderFg = ui.ColorCyan

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, par),
		),
		ui.NewRow(
			ui.NewCol(12, 0, lc1),
		),
	)

	draw := func(t int) {
		newValue, processInfo := getRSS("Sublime")
		lc1.Data = updateRSS(rss, newValue)
		par.Text = processInfo

		ui.Clear()
		ui.Body.Align()
		ui.Render(ui.Body)
	}
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		draw(int(t.Count))
	})
	ui.Loop()
}

func updateRSS(dataArr []float64, value *float64) []float64 {
	if value == nil {
		return dataArr
	}

	n := len(dataArr)
	for i := 0; i < n-1; i++ {
		(dataArr)[i] = (dataArr)[i+1]
	}
	(dataArr)[n-1] = *value / 1000
	return dataArr
}

func getRSS(process string) (*float64, string) {
	raw, err := sh.Command("ps", "aux").
		Command("grep", *processName).
		Command("grep", "-v", "grep").
		Command("grep", "-v", "ptrack").
		Command("head", "-n1").
		Output()
	if err != nil {
		// fmt.Println(err.Error())
		return nil, ""
	}

	processInfo := string(raw)

	raw2, err := sh.Command("echo", processInfo).
		Command("awk", "{print $6}").
		Output()

	if err != nil {
		fmt.Println(err.Error())
		return nil, ""
	}

	data := strings.Replace(string(raw2), "\n", "", -1)

	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		// fmt.Println(err.Error())
		return nil, ""
	}
	return &f, processInfo
}
