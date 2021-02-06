package widget

import (
	"fmt"
	"math"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"

	//	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// LentryItemID uniquely identifies an item within a lentry.
type LentryItemID = int

// Declare conformity with interfaces.
var _ fyne.Widget = (*lentryItem)(nil)
var _ fyne.Tappable = (*lentryItem)(nil)
var _ desktop.Hoverable = (*lentryItem)(nil)
var _ fyne.Focusable = (*lentryItem)(nil)
var _ fyne.Disableable = (*lentryItem)(nil)

type lentryItem struct {
	BaseWidget
	lentry            *Lentry
	onTapped          func()
	statusIndicator   *canvas.Rectangle
	textCanvas        *canvas.Text
	hovered, selected bool
	focused           bool
	RowId             int
}

func newLentryItem(lentry *Lentry, rowId int, tapped func()) *lentryItem {
	li := &lentryItem{
		lentry:     lentry,
		textCanvas: canvas.NewText("Place Holder", theme.ForegroundColor()),
		onTapped:   tapped,
		RowId:      rowId,
	}

	li.ExtendBaseWidget(li)
	return li
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (li *lentryItem) CreateRenderer() fyne.WidgetRenderer {
	li.ExtendBaseWidget(li)

	li.statusIndicator = canvas.NewRectangle(theme.BackgroundColor())
	li.statusIndicator.Hide()

	objects := []fyne.CanvasObject{li.statusIndicator, li.textCanvas}

	return &lentryItemRenderer{widget.NewBaseRenderer(objects), li}
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
//
// Implements: fyne.Disableable
func (li *lentryItem) Disable() {
	li.lentry.DisableableWidget.Disable()
}

// Disabled returns whether the entry is disabled or read-only.
//
// Implements: fyne.Disableable
func (li *lentryItem) Disabled() bool {
	return li.lentry.DisableableWidget.disabled
}

// Enable this widget, updating any style or features appropriately.
//
// Implements: fyne.Disableable
func (li *lentryItem) Enable() {
	li.lentry.DisableableWidget.Enable()
}

// FocusGained is called when the Entry has been given focus.
//
// Implements: fyne.Focusable
func (li *lentryItem) FocusGained() {
	if li.Disabled() {
		return
	}
	li.setFieldsAndRefresh(func() {
		li.focused = true
		li.selected = true
	})
	fmt.Printf("FocusGained %d\n", li.RowId)
}

// FocusLost is called when the Entry has had focus removed.
//
// Implements: fyne.Focusable
func (li *lentryItem) FocusLost() {
	li.setFieldsAndRefresh(func() {
		li.focused = false
		li.selected = false
	})
	fmt.Printf("FocusLost %d\n", li.RowId)
}

// MinSize returns the size that this widget should not shrink below.
func (li *lentryItem) MinSize() fyne.Size {
	li.ExtendBaseWidget(li)
	return li.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget.
func (li *lentryItem) MouseIn(*desktop.MouseEvent) {
	li.hovered = true
	li.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget.
func (li *lentryItem) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget.
func (li *lentryItem) MouseOut() {
	li.hovered = false
	li.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler.
func (li *lentryItem) Tapped(*fyne.PointEvent) {
	if li.onTapped != nil {
		li.selected = true
		li.Refresh()
		li.onTapped()
	}
}

// TypedKey receives key input events when the Entry widget is focused.
//
// Implements: fyne.Focusable
func (li *lentryItem) TypedKey(key *fyne.KeyEvent) {
	li.lentry.keyHandler(key)
}

// TypedRune receives text input events when the Entry widget is focused.
//
// Implements: fyne.Focusable
func (li *lentryItem) TypedRune(r rune) {
	li.lentry.runeHandler(r)
}

//UpdateItem updates the given LentryItem with the row text specified by the id
func (li *lentryItem) UpdateItem(id LentryItemID, provider *textProvider) {
	row := provider.row(id)
	var line string
	if provider.presenter.concealed() {
		line = strings.Repeat(passwordChar, len(row))
	} else {
		line = string(row)
	}

	li.RowId = id
	li.textCanvas.Text = line
	li.textCanvas.Alignment = provider.presenter.textAlign()
	li.textCanvas.TextStyle = provider.presenter.textStyle()
	li.textCanvas.Refresh()
}

// Declare conformity with the WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*lentryItemRenderer)(nil)

type lentryItemRenderer struct {
	widget.BaseRenderer

	item *lentryItem
}

// MinSize calculates the minimum size of a lentryItem.
// This is based on the size of the status indicator and the size of the child object.
func (li *lentryItemRenderer) MinSize() (size fyne.Size) {
	itemSize := li.item.textCanvas.MinSize()
	size = fyne.NewSize(itemSize.Width+theme.Padding()*3, itemSize.Height) //+theme.Padding()*2)
	return
}

// Layout the components of the lentryItem widget.
func (li *lentryItemRenderer) Layout(size fyne.Size) {
	li.item.statusIndicator.Move(fyne.NewPos(0, 0))
	s := fyne.NewSize(theme.Padding(), size.Height)
	li.item.statusIndicator.SetMinSize(s)
	li.item.statusIndicator.Resize(s)

	li.item.textCanvas.Move(fyne.NewPos(theme.Padding()*2, 0))                         //theme.Padding()))
	li.item.textCanvas.Resize(fyne.NewSize(size.Width-theme.Padding()*3, size.Height)) //-theme.Padding()*2))
}

func (li *lentryItemRenderer) Refresh() {
	if li.item.selected {
		li.item.statusIndicator.FillColor = theme.PrimaryColor()
		li.item.statusIndicator.Show()
	} else if li.item.hovered {
		li.item.statusIndicator.FillColor = theme.HoverColor()
		li.item.statusIndicator.Show()
	} else {
		li.item.statusIndicator.Hide()
	}
	li.item.statusIndicator.Refresh()
	canvas.Refresh(li.item.super())
}

// Declare conformity with Layout interface.
var _ fyne.Layout = (*lentryLayout)(nil)

type lentryLayout struct {
	lentry     *Lentry
	dividers   []fyne.CanvasObject
	children   []fyne.CanvasObject
	layoutEndY float32
}

func newLentryLayout(lentry *Lentry) fyne.Layout {
	return &lentryLayout{lentry: lentry}
}

func (l *lentryLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if l.lentry.offsetY != 0 {
		return
	}
	y := float32(0)
	for _, child := range l.children {
		child.Move(fyne.NewPos(0, y))
		y += l.lentry.itemMin.Height
		child.Resize(fyne.NewSize(l.lentry.size.Width, l.lentry.itemMin.Height))
	}
	l.layoutEndY = y
}

func (l *lentryLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if f := l.lentry.Length; f != nil {
		return fyne.NewSize(l.lentry.itemMin.Width, l.lentry.itemMin.Height*float32(f()))
	}
	return fyne.NewSize(0, 0)
}

func (l *lentryLayout) getObjects() []fyne.CanvasObject {
	objects := l.children
	objects = append(objects, l.dividers...)
	return objects
}

func (l *lentryLayout) appendedItem(objects []fyne.CanvasObject) {
	if len(objects) > 1 {
		objects[len(objects)-1].Move(fyne.NewPos(0, objects[len(objects)-2].Position().Y+l.lentry.itemMin.Height))
	} else {
		objects[len(objects)-1].Move(fyne.NewPos(0, 0))
	}
	objects[len(objects)-1].Resize(fyne.NewSize(l.lentry.size.Width, l.lentry.itemMin.Height))
}

func (l *lentryLayout) prependedItem(objects []fyne.CanvasObject) {
	objects[0].Move(fyne.NewPos(0, objects[1].Position().Y-l.lentry.itemMin.Height))
	objects[0].Resize(fyne.NewSize(l.lentry.size.Width, l.lentry.itemMin.Height))
}

// Declare conformity with WidgetRenderer interface.
var _ fyne.WidgetRenderer = (*lentryRenderer)(nil)

type lentryRenderer struct {
	widget.BaseRenderer

	lentry           *Lentry
	scroller         *widget.Scroll
	layout           *fyne.Container
	itemPool         *syncPool
	children         []fyne.CanvasObject
	size             fyne.Size
	visibleItemCount int
	firstItemIndex   LentryItemID
	lastItemIndex    LentryItemID
	previousOffsetY  float32
}

func newLentryRenderer(objects []fyne.CanvasObject, l *Lentry, scroller *widget.Scroll, layout *fyne.Container) *lentryRenderer {
	lr := &lentryRenderer{BaseRenderer: widget.NewBaseRenderer(objects), lentry: l, scroller: scroller, layout: layout}
	lr.scroller.OnScrolled = lr.offsetUpdated
	return lr
}

func (l *lentryRenderer) Layout(size fyne.Size) {
	length := 0
	if f := l.lentry.Length; f != nil {
		length = f()
	}
	if length <= 0 {
		if len(l.children) > 0 {
			for _, child := range l.children {
				l.itemPool.Release(child)
			}
			l.previousOffsetY = 0
			l.firstItemIndex = 0
			l.lastItemIndex = 0
			l.visibleItemCount = 0
			l.lentry.offsetY = 0
			l.layout.Layout.(*lentryLayout).layoutEndY = 0
			l.children = nil
			l.layout.Objects = nil
			l.lentry.Refresh()
		}
		return
	}
	if size != l.size {
		if size.Width != l.size.Width {
			for _, child := range l.children {
				child.Resize(fyne.NewSize(size.Width, l.lentry.itemMin.Height))
			}
		}
		l.scroller.Resize(size)
		l.size = size
	}
	if l.itemPool == nil {
		l.itemPool = &syncPool{}
	}

	// Relayout What Is Visible - no scroll change - initial layout or possibly from a resize.
	l.visibleItemCount = int(math.Ceil(float64(l.scroller.Size().Height) / float64(l.lentry.itemMin.Height)))
	if l.visibleItemCount <= 0 {
		return
	}
	min := int(fyne.Min(float32(length), float32(l.visibleItemCount)))
	if len(l.children) > min {
		for i := len(l.children); i >= min; i-- {
			l.itemPool.Release(l.children[i-1])
		}
		l.children = l.children[:min-1]
	}
	for i := len(l.children) + l.firstItemIndex; len(l.children) <= l.visibleItemCount && i < length; i++ {
		l.appendItem(i)
	}
	l.layout.Layout.(*lentryLayout).children = l.children
	l.layout.Layout.Layout(l.children, l.lentry.itemMin)
	l.layout.Objects = l.layout.Layout.(*lentryLayout).getObjects()
	l.lastItemIndex = l.firstItemIndex + len(l.children) - 1

	i := l.firstItemIndex
	for _, child := range l.children {
		child.(*lentryItem).UpdateItem(i, l.lentry.provider)
		l.setupLentryItem(child, i)
		i++
	}
}

func (l *lentryRenderer) MinSize() fyne.Size {
	return l.scroller.MinSize().Max(l.lentry.itemMin)
}

func (l *lentryRenderer) Refresh() {
	if l.lentry.Text != string(l.lentry.provider.buffer) {
		l.lentry.provider.setText(l.lentry.Text)
	} else {
		l.lentry.provider.updateRowBounds() // if truncate/wrap has changed
	}
	l.lentry.itemMin = newLentryItem(nil, -1, nil).MinSize()

	l.Layout(l.lentry.Size())
	l.scroller.Refresh()

	l.lentry.provider.propertyLock.Lock()
	l.lentry.provider.updateRowBounds()
	l.lentry.provider.propertyLock.Unlock()

	canvas.Refresh(l.lentry.super())
}

// Resize sets a new size for the lentry.
// Note this should not be used if the widget is being managed by a Layout within a Container.
func (l *Lentry) Resize(size fyne.Size) {
	l.BaseWidget.Resize(size)
	if l.provider == nil { // not created until visible
		return
	}
	l.provider.Resize(size)
}

func (l *lentryRenderer) appendItem(id LentryItemID) {
	item := l.getItem()
	l.children = append(l.children, item)
	l.setupLentryItem(item, id)
	l.layout.Layout.(*lentryLayout).children = l.children
	l.layout.Layout.(*lentryLayout).appendedItem(l.children)
	l.layout.Objects = l.layout.Layout.(*lentryLayout).getObjects()
}

func (l *lentryRenderer) getItem() fyne.CanvasObject {
	item := l.itemPool.Obtain()
	if item == nil {
		item = newLentryItem(l.lentry, -1, nil)
	}
	return item
}

func (l *lentryRenderer) offsetChanged() {
	offsetChange := float32(math.Abs(float64(l.previousOffsetY - l.lentry.offsetY)))

	if l.previousOffsetY < l.lentry.offsetY {
		// Scrolling Down.
		l.scrollDown(offsetChange)
	} else if l.previousOffsetY > l.lentry.offsetY {
		// Scrolling Up.
		l.scrollUp(offsetChange)
	}
}

func (l *lentryRenderer) prependItem(id LentryItemID) {
	item := l.getItem()
	l.children = append([]fyne.CanvasObject{item}, l.children...)
	l.setupLentryItem(item, id)
	l.layout.Layout.(*lentryLayout).children = l.children
	l.layout.Layout.(*lentryLayout).prependedItem(l.children)
	l.layout.Objects = l.layout.Layout.(*lentryLayout).getObjects()
}

func (l *lentryRenderer) scrollDown(offsetChange float32) {
	itemChange := 0
	layoutEndY := l.children[len(l.children)-1].Position().Y + l.lentry.itemMin.Height
	scrollerEndY := l.scroller.Offset.Y + l.scroller.Size().Height
	if layoutEndY < scrollerEndY {
		itemChange = int(math.Ceil(float64(scrollerEndY-layoutEndY) / float64(l.lentry.itemMin.Height)))
	} else if offsetChange < l.lentry.itemMin.Height {
		return
	} else {
		itemChange = int(math.Floor(float64(offsetChange) / float64(l.lentry.itemMin.Height)))
	}
	l.previousOffsetY = l.lentry.offsetY
	length := 0
	if f := l.lentry.Length; f != nil {
		length = f()
	}
	if length == 0 {
		return
	}
	for i := 0; i < itemChange && l.lastItemIndex != length-1; i++ {
		l.itemPool.Release(l.children[0])
		l.children = l.children[1:]
		l.firstItemIndex++
		l.lastItemIndex++
		l.appendItem(l.lastItemIndex)
	}
}

func (l *lentryRenderer) scrollUp(offsetChange float32) {
	itemChange := 0
	layoutStartY := l.children[0].Position().Y
	if layoutStartY > l.scroller.Offset.Y {
		itemChange = int(math.Ceil(float64(layoutStartY-l.scroller.Offset.Y) / float64(l.lentry.itemMin.Height)))
	} else if offsetChange < l.lentry.itemMin.Height {
		return
	} else {
		itemChange = int(math.Floor(float64(offsetChange) / float64(l.lentry.itemMin.Height)))
	}
	l.previousOffsetY = l.lentry.offsetY
	for i := 0; i < itemChange && l.firstItemIndex != 0; i++ {
		l.itemPool.Release(l.children[len(l.children)-1])
		l.children = l.children[:len(l.children)-1]
		l.firstItemIndex--
		l.lastItemIndex--
		l.prependItem(l.firstItemIndex)
	}
}

func (l *lentryRenderer) setupLentryItem(item fyne.CanvasObject, id LentryItemID) {
	li := item.(*lentryItem)
	previousIndicator := li.selected
	li.selected = false
	for _, s := range l.lentry.selected {
		if id == s {
			li.selected = true
		}
	}
	if previousIndicator != li.selected {
		item.Refresh()
	}
	li.UpdateItem(id, l.lentry.provider)

	li.onTapped = func() {
		l.lentry.Select(id)
	}
}

func (l *lentryRenderer) offsetUpdated(pos fyne.Position) {
	if l.lentry.offsetY == pos.Y {
		return
	}
	l.lentry.offsetY = pos.Y
	l.offsetChanged()
}
