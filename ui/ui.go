package ui

import (
	"github.com/Stroby241/UntisQuerry/event"
	"github.com/Stroby241/UntisQuerry/state"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

var pages []widget.PreferredSizeLocateableWidget

func Init() (*ebitenui.UI, func(), error) {
	res, err := newUIResources()
	if err != nil {
		return nil, nil, err
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true}),
			widget.GridLayoutOpts.Spacing(0, 20))),
	)

	drag := &dragContents{
		res: res,
	}

	toolTips := toolTipContents{
		tips: map[widget.HasWidget]string{},
		res:  res,
	}

	toolTip := widget.NewToolTip(
		widget.ToolTipOpts.Container(rootContainer),
		widget.ToolTipOpts.ContentsCreater(&toolTips),
	)

	dnd := widget.NewDragAndDrop(
		widget.DragAndDropOpts.Container(rootContainer),
		widget.DragAndDropOpts.ContentsCreater(drag),
	)

	ui := &ebitenui.UI{
		Container:   rootContainer,
		ToolTip:     toolTip,
		DragAndDrop: dnd,
	}
	createPages(res, ui)

	flipBook := widget.NewFlipBook(
		widget.FlipBookOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		}))),
	)
	rootContainer.AddChild(flipBook)

	event.On(event.EventUpdate, func(data interface{}) {
		ui.Update()
	})

	event.On(event.EventDraw, func(data interface{}) {
		ui.Draw(data.(*ebiten.Image))
	})

	event.On(event.EventSetPage, func(data interface{}) {
		setPage(data.(int), flipBook)
	})

	return ui, func() {
		res.close()
	}, nil
}

func createPages(res *uiResources, ui *ebitenui.UI) {

	pages = make([]widget.PreferredSizeLocateableWidget, state.PageMax)

	/*
		uiFunc := func() *ebitenui.UI {
			return ui
		}
	*/

	pages[state.PageStart] = createStartPage(res)
	pages[state.PageLogin] = createLoginPage(res)
	pages[state.PageQuerry] = createQuerryPage(res)
	pages[state.PageAddTeacher] = createAddTeacherPage(res)
	pages[state.PageLoading] = createLoadingPage(res)

	event.On(event.EventUpdateQuerryPanel, func(data interface{}) {
		pages[state.PageQuerry] = createQuerryPage(res)
	})
}

func setPage(pageId int, flipBook *widget.FlipBook) {
	if len(pages) >= pageId && pages[pageId] == nil {
		return
	}

	page := pages[pageId]
	flipBook.SetPage(page)
}
