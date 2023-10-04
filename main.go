package main

import (
	g "github.com/AllenDang/giu"
)

var (
	leftMenuWidth   float32 = 200
	selectedProject string
)

func loop() {
	g.SingleWindow().Layout(
		g.SplitLayout(g.DirectionVertical, &leftMenuWidth,
			g.Layout{
				g.Label("Projects"),
				NewProjectButton("Project One"),
				NewProjectButton("Project Two"),
				NewProjectButton("Poject Three"),
			},
			g.Layout{
				g.Label("Main Frame"),
				renderItemsForProject(),
			},
		),
	)
}

func NewProjectButton(title string) *g.ButtonWidget {
	return g.Button(title).OnClick(func() {
		selectedProject = title
	})
}

func renderItemsForProject() g.Widget {
	switch selectedProject {
	case "Project One":
		return g.Layout{
			g.Label("Item 1 for Project One"),
			g.Label("Item 2 for Project One"),
		}
	case "Project Two":
		return g.Layout{
			g.Label("Item 1 for Project Two"),
			g.Label("Item 2 for Project Two"),
		}
	case "Poject Three":
		return g.Layout{
			g.Label("Item 1 for Project Three"),
			g.Label("Item 2 for Project Three"),
		}
	default:
		return g.Label("Select a project to view items.")
	}
}

func main() {
	w := g.NewMasterWindow("Flamego", 1000, 700, g.MasterWindowFlagsFloating)
	w.Run(loop)
}
