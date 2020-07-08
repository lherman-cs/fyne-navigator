package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type Navigator interface {
	Push(path string, ctx interface{}) error
	Pop() error
	Reset()
}

type Page interface {
	fyne.Widget
	BeforeDestroy()
}

type navigationRenderer struct {
	navigation *NavigationContainer
}

func (r *navigationRenderer) Layout(size fyne.Size) {
	r.navigation.currentPage.Resize(size)
}

func (r *navigationRenderer) MinSize() fyne.Size {
	return r.navigation.currentPage.MinSize()
}

func (r *navigationRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (r *navigationRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.navigation.currentPage}
}

func (r *navigationRenderer) Destroy() {}

func (r *navigationRenderer) Refresh() {
	r.navigation.currentPage.Refresh()
}

type NavigationHandler func(navigator Navigator, ctx interface{}) (fyne.Widget, error)

type NavigationItem struct {
	Path    string
	Handler NavigationHandler
}

func NewNavigationItem(path string, handle func(Navigator, interface{}) (fyne.Widget, error)) *NavigationItem {
	return &NavigationItem{path, NavigationHandler(handle)}
}

type NavigationContainer struct {
	widget.BaseWidget
	BeforeEnter func(to string) error
	routes      map[string]NavigationHandler
	history     []fyne.Widget
	currentPage fyne.Widget
}

func NewNavigationContainer(initialPath string, items ...*NavigationItem) (*NavigationContainer, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("routes can't be empty")
	}

	routes := make(map[string]NavigationHandler)
	for _, item := range items {
		routes[item.Path] = item.Handler
	}

	buildInitialPage, ok := routes[initialPath]
	if !ok {
		return nil, fmt.Errorf("%s doesn't exist in routes: %v", initialPath, routes)
	}

	navigation := NavigationContainer{
		routes: routes,
	}
	initialPage, err := buildInitialPage(&navigation, nil)
	if err != nil {
		return nil, err
	}
	navigation.currentPage = initialPage
	navigation.ExtendBaseWidget(&navigation)
	return &navigation, nil
}

func (navigation *NavigationContainer) CreateRenderer() fyne.WidgetRenderer {
	return &navigationRenderer{navigation: navigation}
}

func (navigation *NavigationContainer) Push(path string, ctx interface{}) error {
	buildNextPage, ok := navigation.routes[path]
	if !ok {
		return fmt.Errorf("%s doesn't exist in routes: %v", path, navigation.routes)
	}

	if navigation.BeforeEnter != nil {
		err := navigation.BeforeEnter(path)
		if err != nil {
			return err
		}
	}
	navigation.history = append(navigation.history, navigation.currentPage)
	nextPage, err := buildNextPage(navigation, ctx)
	if err != nil {
		return err
	}
	navigation.currentPage = nextPage
	navigation.Refresh()
	return nil
}

func (navigation *NavigationContainer) Pop() error {
	if len(navigation.history) == 0 {
		return fmt.Errorf("there's no more pages in history")
	}

	if page, ok := navigation.currentPage.(Page); ok {
		page.BeforeDestroy()
	}
	navigation.currentPage = navigation.history[len(navigation.history)-1]
	navigation.history = navigation.history[:len(navigation.history)-1]
	navigation.Refresh()
	return nil
}

func (navigation *NavigationContainer) Reset() {
	for navigation.Pop() == nil {
	}
}
