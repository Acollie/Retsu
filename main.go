package main

import (
	"context"
	"log"
	"queue/awsx"
	"queue/queue"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	screenWidth  int
	screenHeight int
)

type cli struct {
	url      []string
	pageSize uint32
	page     uint32
	handler  queue.Handler
}

func (d *cli) renderBasicTable(table *widgets.Table, selected int) {
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.SetRect(0, 0, 60, 100)

	rows := [][]string{}

	// Calculate start and end indices for pagination
	start := int(d.page) * int(d.pageSize)
	end := start + int(d.pageSize)
	if end > len(d.url) {
		end = len(d.url)
	}
	if selected >= end {
		d.page++
		d.renderBasicTable(table, selected)
		return
	}

	// Pick rows to display
	for i, row := range d.url[start:end] {
		if i == (selected - start) {
			rows = append(rows, []string{"[x]" + row})
		} else {
			rows = append(rows, []string{"[ ]" + row})
		}
	}
	table.Rows = rows
}

func (d cli) renderInfo(queueInfo *queue.Queue) {
	p := widgets.NewParagraph()
	p.Title = "Queue:" + queueInfo.Name
	p.Text = "Messages:"

	p.SetRect(60, 0, screenWidth, 100)
	ui.Render(p)
}

func main() {
	ctx := context.Background()
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	cfg := awsx.LoadConfig(ctx)
	handler := awsx.NewHandler(*cfg)

	screenHeight, screenWidth = ui.TerminalDimensions()
	selected := 0
	moreInfo := false

	basTable := widgets.NewTable()

	urls, err := handler.Scan(ctx)
	if err != nil {
		panic(err)
	}

	basTable.TextStyle = ui.NewStyle(ui.ColorWhite)
	d := cli{
		url:      urls,
		page:     0,
		pageSize: uint32(screenHeight - 2),
	}

	d.renderBasicTable(basTable, selected)
	ui.Render(basTable)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {

		case "i":
			queue, err := d.handler.GetQueue(ctx, d.url[selected])
			if err != nil {
				panic(err)
			}
			if moreInfo {
				d.renderInfo(queue)
				moreInfo = true
			} else {
				moreInfo = false
			}

		case "k":
			if len(basTable.Rows) < 0 {
				selected = 0
			} else {
				selected--
			}
			d.renderBasicTable(basTable, selected)
			ui.Render(basTable)

		case "j":
			selected++
			if selected >= len(basTable.Rows) {
				selected = 0
			}
			d.renderBasicTable(basTable, selected)

			ui.Render(basTable)
		case "q", "<C-c>":
			return
		}
	}
}
