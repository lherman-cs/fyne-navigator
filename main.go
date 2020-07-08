package main

import (
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

type Page2 struct {
	fyne.Widget
}

type Page2Context struct {
	Label string
}

func NewPage2(navigator Navigator, ctx Page2Context) fyne.Widget {
	var content *widget.Box
	label := widget.NewLabel(ctx.Label)
	button := widget.NewButton("Pop", func() {
		err := navigator.Pop()
		if err != nil {
			panic(err)
		}
	})
	content = widget.NewVBox(
		label,
		button,
	)

	return &Page2{
		Widget: content,
	}
}

func (page *Page2) BeforeDestroy() {
	fmt.Println("Page2: BeforeDestroy")
}

func main() {
	var err error

	a := app.New()
	w := a.NewWindow("Hello")

	page1 := NewNavigationItem("/page1", func(navigator Navigator, ctx interface{}) (fyne.Widget, error) {
		fmt.Println("Pushing /page1")
		return widget.NewVBox(
			widget.NewButton("Push", func() {
				navigator.Push("/page2", Page2Context{"From Page 1"})
			}),
		), nil
	})

	page2 := NewNavigationItem("/page2", func(navigator Navigator, ctx interface{}) (fyne.Widget, error) {
		fmt.Println("Pushing /page2")
		return NewPage2(navigator, ctx.(Page2Context)), nil
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
