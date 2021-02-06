package widget

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"

	//	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/internal/widget"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with Widget interface.
var _ fyne.Disableable = (*Lentry)(nil)

//var _ fyne.Draggable = (*Lentry)(nil)
//var _ fyne.Focusable = (*Lentry)(nil)

//var _ fyne.Tappable = (*Lentry)(nil)
var _ fyne.Widget = (*Lentry)(nil)

//var _ desktop.Mouseable = (*Lentry)(nil)
//var _ desktop.Keyable = (*Lentry)(nil)

//var _ mobile.Keyboardable = (*Entry)(nil)

// Lentry is a widget that pools lentry items for performance and
// lays the items out in a vertical direction inside of a scroller.
// Lentry requires that all items are the same size.
//
// Since: 1.4
type Lentry struct {
	DisableableWidget
	// list stuff
	OnSelected   func(id LentryItemID)
	OnUnselected func(id LentryItemID)

	scroller      *widget.Scroll
	selected      []LentryItemID
	itemMin       fyne.Size
	offsetY       float32
	offsetUpdated func(fyne.Position)

	// entry stuff
	focused   bool
	Text      string
	provider  *textProvider
	Password  bool
	Alignment fyne.TextAlign // The alignment of the Text
	Wrapping  fyne.TextWrap  // The wrapping of the Text
	TextStyle fyne.TextStyle // The style of the label text
}

// NewLentry creates and returns a lentry widget for displaying items in
// a vertical layout with scrolling and caching for performance.
//
// Since: 1.4
func NewLentry(text string) *Lentry {
	lentry := &Lentry{Text: text, Wrapping: fyne.TextTruncate}
	lentry.ExtendBaseWidget(lentry)
	return lentry
}

// NewLentryWithData creates a new list widget that will display the contents of the provided data.
//
// Since: 2.0
//TODO sort out databinding
// func NewLentryWithData(data binding.DataList) *Lentry {
// 	l := NewLentry(
// 		data.Length,
// 		createItem,
// 		func(i LentryItemID, o fyne.CanvasObject) {
// 			item, err := data.GetItem(i)
// 			if err != nil {
// 				fyne.LogError(fmt.Sprintf("Error getting data item %d", i), err)
// 				return
// 			}
// 			updateItem(item, o)
// 		})

// 	data.AddListener(binding.NewDataListener(l.Refresh))
// 	return l
// }

// CreateRenderer is a private method to Fyne which links this widget to its renderer.
func (l *Lentry) CreateRenderer() fyne.WidgetRenderer {
	l.ExtendBaseWidget(l)
	// initialise
	l.provider = newTextProvider(l.Text, l)
	l.provider.extraPad = fyne.NewSize(0, theme.Padding())
	l.provider.size = l.size

	if l.itemMin.IsZero() {
		l.itemMin = newLentryItem(nil, -1, nil).MinSize()
	}
	layout := fyne.NewContainerWithLayout(newLentryLayout(l))
	layout.Resize(layout.MinSize())
	l.scroller = widget.NewVScroll(layout)
	objects := []fyne.CanvasObject{l.scroller}
	lr := newLentryRenderer(objects, l, l.scroller, layout)
	l.offsetUpdated = lr.offsetUpdated
	return lr
}

// Cursor returns the cursor type of this widget
//
// Implements: desktop.Cursorable
func (l *Lentry) Cursor() desktop.Cursor {
	return desktop.TextCursor
}

// Disable this widget so that it cannot be interacted with, updating any style appropriately.
//
// Implements: fyne.Disableable
func (l *Lentry) Disable() {
	l.DisableableWidget.Disable()
}

// Disabled returns whether the entry is disabled or read-only.
//
// Implements: fyne.Disableable
func (l *Lentry) Disabled() bool {
	return l.DisableableWidget.disabled
}

// Enable this widget, updating any style or features appropriately.
//
// Implements: fyne.Disableable
func (l *Lentry) Enable() {
	l.DisableableWidget.Enable()
}

// FocusGained is called when the Entry has been given focus.
//
// Implements: fyne.Focusable
func (l *Lentry) FocusGained() {
	if l.Disabled() {
		return
	}
	l.setFieldsAndRefresh(func() {
		l.focused = true
	})
	fmt.Println("lentry FocusGained")
}

// FocusLost is called when the Entry has had focus removed.
//
// Implements: fyne.Focusable
func (l *Lentry) FocusLost() {
	l.setFieldsAndRefresh(func() {
		l.focused = false
	})
	fmt.Println("lentry FocusLost")

}

// MinSize returns the size that this widget should not shrink below.
func (l *Lentry) MinSize() fyne.Size {
	l.ExtendBaseWidget(l)

	return l.BaseWidget.MinSize()
}

// Select add the item identified by the given ID to the selection.
func (l *Lentry) Select(id LentryItemID) {
	if len(l.selected) > 0 && id == l.selected[0] {
		return
	}
	length := 0
	if f := l.Length; f != nil {
		length = f()
	}
	if id < 0 || id >= length {
		return
	}
	old := l.selected
	l.selected = []LentryItemID{id}
	defer func() {
		if f := l.OnUnselected; f != nil && len(old) > 0 {
			f(old[0])
		}
		if f := l.OnSelected; f != nil {
			f(id)
		}
	}()
	if l.scroller == nil {
		return
	}
	y := float32(id) * l.itemMin.Height
	if y < l.scroller.Offset.Y {
		l.scroller.Offset.Y = y
	} else if y+l.itemMin.Height > l.scroller.Offset.Y+l.scroller.Size().Height {
		l.scroller.Offset.Y = y + l.itemMin.Height - l.scroller.Size().Height
	}
	l.offsetUpdated(l.scroller.Offset)
	l.Refresh()
}

// keyHandler processes key events.
func (l *Lentry) keyHandler(key *fyne.KeyEvent) {
	if l.Disabled() {
		return
	}
	fmt.Printf("key=%s\n", key.Name)
}

// runeHandler proceses rune events.
func (l *Lentry) runeHandler(r rune) {
	if l.Disabled() {
		return
	}
	fmt.Printf("rune=%s\n", string(r))
}

// Unselect removes the item identified by the given ID from the selection.
func (l *Lentry) Unselect(id LentryItemID) {
	if len(l.selected) == 0 {
		return
	}

	l.selected = nil
	l.Refresh()
	if f := l.OnUnselected; f != nil {
		f(id)
	}
}

// Length returns the number of rows in the testProvider
func (l *Lentry) Length() int {
	return l.provider.rows()
}

// Declare conformity with textPresenter interface.
var _ textPresenter = (*Lentry)(nil)

// textAlign tells the rendering textProvider our alignment
func (l *Lentry) textAlign() fyne.TextAlign {
	return l.Alignment
}

// textWrap tells the rendering textProvider our wrapping
func (l *Lentry) textWrap() fyne.TextWrap {
	return l.Wrapping
}

// textStyle tells the rendering textProvider our style
func (l *Lentry) textStyle() fyne.TextStyle {
	return l.TextStyle
}

// textColor tells the rendering textProvider our color
func (l *Lentry) textColor() color.Color {
	return theme.ForegroundColor()
}

// concealed tells the rendering textProvider if we are a concealed field
func (l *Lentry) concealed() bool {
	return false
}

// object returns the root object of the widget so it can be referenced
func (l *Lentry) object() fyne.Widget {
	return l.super()
}
