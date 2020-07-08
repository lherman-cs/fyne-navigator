package main

import (
	"fmt"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

type Page1 struct {
	fyne.Widget
	cntr int
}

func NewPage1(navigator Navigator) fyne.Widget {
	var content *widget.Box
	cntr := 0
	button := widget.NewButton("Push", func() {
		err := navigator.Push("/page2", nil)
		if err != nil {
			panic(err)
		}

		content.Append(widget.NewButton("Pada", func() {}))
		cntr++
		content.Refresh()
	})
	content = widget.NewVBox(
		button,
	)

	return &Page1{
		Widget: content,
	}
}

type Page2 struct {
	fyne.Widget
}

func (p *Page2) BeforeDestroy() {
	fmt.Println("Page 2: BeforeDestroy")
}

func main() {
	var err error

	a := app.New()
	w := a.NewWindow("Hello")

	page1 := NewNavigationItem("/page1", func(navigator Navigator, ctx interface{}) (fyne.Widget, error) {
		return NewPage1(navigator), nil
	})

	page2 := NewNavigationItem("/page2", func(navigator Navigator, ctx interface{}) (fyne.Widget, error) {
		log.Println("Building page 2")
		return &Page2{
			Widget: widget.NewVBox(
				widget.NewButton("Pop", func() {
					navigator.Pop()
				}),
			),
		}, nil
	})

	router, err := NewNavigationContainer("/page1", page1, page2)
	if err != nil {
		panic(err)
	}

	router.BeforeEnter = func(to string) error {
		fmt.Println(to)
		return nil
	}
	w.SetContent(router)
	w.ShowAndRun()
}
