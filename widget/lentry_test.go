package widget

import (
	"fmt"
	"image/color"
	"math"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"

	"github.com/stretchr/testify/assert"
)

func TestNewLentry(t *testing.T) {
	lentry := createLentry(1000)

	template := newLentryItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)
	firstItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(lentry).(*lentryRenderer).children)

	assert.Equal(t, 1000, lentry.Length())
	assert.GreaterOrEqual(t, lentry.MinSize().Width, template.MinSize().Width)
	assert.Equal(t, lentry.MinSize(), template.MinSize().Max(test.WidgetRenderer(lentry).(*lentryRenderer).scroller.MinSize()))
	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)
}

func TestLentry_MinSize(t *testing.T) {
	for name, tt := range map[string]struct {
		cellSize        fyne.Size
		expectedMinSize fyne.Size
	}{
		"small": {
			fyne.NewSize(1, 1),
			fyne.NewSize(float32(32), float32(32)),
		},
		"large": {
			fyne.NewSize(100, 100),
			fyne.NewSize(100+3*theme.Padding(), 100+2*theme.Padding()),
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expectedMinSize, NewLentry(
				func() int { return 5 },
				func() fyne.CanvasObject {
					r := canvas.NewRectangle(color.Black)
					r.SetMinSize(tt.cellSize)
					r.Resize(tt.cellSize)
					return r
				},
				func(LentryItemID, fyne.CanvasObject) {}).MinSize())
		})
	}
}

func TestLentry_Resize(t *testing.T) {
	defer test.NewApp()
	lentry, w := setupLentry(t)
	template := newLentryItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)

	firstItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(lentry).(*lentryRenderer).children)
	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)

	w.Resize(fyne.NewSize(200, 600))

	newFirstItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex
	newLastItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).lastItemIndex
	newVisibleCount := len(test.WidgetRenderer(lentry).(*lentryRenderer).children)
	expectedNewLastItemIndex := int(math.Ceil(float64(600)/float64(template.MinSize().Height))) + newFirstItemIndex

	assert.Equal(t, firstItemIndex, newFirstItemIndex)
	assert.NotEqual(t, lastItemIndex, newLastItemIndex)
	assert.Equal(t, expectedNewLastItemIndex, newLastItemIndex)
	assert.NotEqual(t, visibleCount, newVisibleCount)
	assert.Equal(t, newVisibleCount, newLastItemIndex-newFirstItemIndex+1)
	test.AssertRendersToMarkup(t, "lentry/resized.xml", w.Canvas())
}

func TestLentry_OffsetChange(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	lentry := createLentry(1000)
	w := test.NewWindow(lentry)
	w.Resize(fyne.NewSize(200, 400))
	template := newLentryItem(fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object")), nil)

	firstItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).lastItemIndex
	visibleCount := test.WidgetRenderer(lentry).(*lentryRenderer).visibleItemCount

	assert.Equal(t, 0, firstItemIndex)
	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex)

	scroll := test.WidgetRenderer(lentry).(*lentryRenderer).scroller
	scroll.Scrolled(&fyne.ScrollEvent{Scrolled: fyne.NewDelta(0, -280)})

	indexChange := int(math.Floor(float64(300) / float64(template.MinSize().Height)))

	newFirstItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex
	newLastItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).lastItemIndex
	newVisibleCount := test.WidgetRenderer(lentry).(*lentryRenderer).visibleItemCount

	assert.NotEqual(t, firstItemIndex, newFirstItemIndex)
	assert.Equal(t, newFirstItemIndex, firstItemIndex+indexChange-1)
	assert.NotEqual(t, lastItemIndex, newLastItemIndex)
	assert.Equal(t, newLastItemIndex, lastItemIndex+indexChange-1)
	assert.Equal(t, visibleCount, newVisibleCount)
	assert.Equal(t, newVisibleCount, newLastItemIndex-newFirstItemIndex)
	test.AssertRendersToMarkup(t, "lentry/offset_changed.xml", w.Canvas())
}

func TestLentry_Hover(t *testing.T) {
	lentry := createLentry(1000)
	children := test.WidgetRenderer(lentry).(*lentryRenderer).children

	for i := 0; i < 2; i++ {
		assert.False(t, children[i].(*lentryItem).statusIndicator.Visible())
		children[i].(*lentryItem).MouseIn(&desktop.MouseEvent{})
		assert.Equal(t, children[i].(*lentryItem).statusIndicator.FillColor, theme.HoverColor())
		children[i].(*lentryItem).MouseOut()
		assert.False(t, children[i].(*lentryItem).statusIndicator.Visible())
	}
}

func TestLentry_Selection(t *testing.T) {
	lentry := createLentry(1000)
	children := test.WidgetRenderer(lentry).(*lentryRenderer).children

	assert.False(t, children[0].(*lentryItem).statusIndicator.Visible())
	children[0].(*lentryItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[0].(*lentryItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[0].(*lentryItem).statusIndicator.Visible())
	assert.Equal(t, 1, len(lentry.selected))
	assert.Equal(t, 0, lentry.selected[0])
	children[1].(*lentryItem).Tapped(&fyne.PointEvent{})
	assert.Equal(t, children[1].(*lentryItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[1].(*lentryItem).statusIndicator.Visible())
	assert.Equal(t, 1, len(lentry.selected))
	assert.Equal(t, 1, lentry.selected[0])
	assert.False(t, children[0].(*lentryItem).statusIndicator.Visible())
}

func TestLentry_Select(t *testing.T) {
	lentry := createLentry(1000)

	assert.Equal(t, test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex, 0)
	lentry.Select(50)
	assert.Equal(t, test.WidgetRenderer(lentry).(*lentryRenderer).lastItemIndex, 50)
	children := test.WidgetRenderer(lentry).(*lentryRenderer).children
	assert.Equal(t, children[len(children)-1].(*lentryItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[len(children)-1].(*lentryItem).statusIndicator.Visible())

	lentry.Select(5)
	assert.Equal(t, test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex, 5)
	children = test.WidgetRenderer(lentry).(*lentryRenderer).children
	assert.Equal(t, children[0].(*lentryItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[0].(*lentryItem).statusIndicator.Visible())

	lentry.Select(6)
	assert.Equal(t, test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex, 5)
	children = test.WidgetRenderer(lentry).(*lentryRenderer).children
	assert.False(t, children[0].(*lentryItem).statusIndicator.Visible())
	assert.Equal(t, children[1].(*lentryItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[1].(*lentryItem).statusIndicator.Visible())
}

func TestLentry_Unselect(t *testing.T) {
	lentry := createLentry(1000)

	lentry.Select(10)
	children := test.WidgetRenderer(lentry).(*lentryRenderer).children
	assert.Equal(t, children[10].(*lentryItem).statusIndicator.FillColor, theme.PrimaryColor())
	assert.True(t, children[10].(*lentryItem).statusIndicator.Visible())

	lentry.Unselect(10)
	children = test.WidgetRenderer(lentry).(*lentryRenderer).children
	assert.False(t, children[10].(*lentryItem).statusIndicator.Visible())
	assert.Nil(t, lentry.selected)
}

func TestLentry_DataChange(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	lentry, w := setupLentry(t)
	children := test.WidgetRenderer(lentry).(*lentryRenderer).children

	assert.Equal(t, children[0].(*lentryItem).child.(*fyne.Container).Objects[1].(*Label).Text, "Test Item 0")
	changeLentryData(lentry)
	lentry.Refresh()
	children = test.WidgetRenderer(lentry).(*lentryRenderer).children
	assert.Equal(t, children[0].(*lentryItem).child.(*fyne.Container).Objects[1].(*Label).Text, "a")
	test.AssertRendersToMarkup(t, "lentry/new_data.xml", w.Canvas())
}

func TestLentry_ThemeChange(t *testing.T) {
	defer test.NewApp()
	lentry, w := setupLentry(t)

	test.AssertImageMatches(t, "lentry/lentry_initial.png", w.Canvas().Capture())

	test.WithTestTheme(t, func() {
		time.Sleep(100 * time.Millisecond)
		lentry.Refresh()
		test.AssertImageMatches(t, "lentry/lentry_theme_changed.png", w.Canvas().Capture())
	})
}

func TestLentry_SmallLentry(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	var data []string
	data = append(data, "Test Item 0")

	lentry := NewLentry(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(id LentryItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	w := test.NewWindow(lentry)
	w.Resize(fyne.NewSize(200, 400))

	visibleCount := len(test.WidgetRenderer(lentry).(*lentryRenderer).children)
	assert.Equal(t, visibleCount, 1)

	data = append(data, "Test Item 1")
	lentry.Refresh()

	visibleCount = len(test.WidgetRenderer(lentry).(*lentryRenderer).children)
	assert.Equal(t, visibleCount, 2)

	test.AssertRendersToMarkup(t, "lentry/small.xml", w.Canvas())
}

func TestLentry_ClearLentry(t *testing.T) {
	defer test.NewApp()
	lentry, w := setupLentry(t)
	assert.Equal(t, 1000, lentry.Length())

	firstItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).firstItemIndex
	lastItemIndex := test.WidgetRenderer(lentry).(*lentryRenderer).lastItemIndex
	visibleCount := len(test.WidgetRenderer(lentry).(*lentryRenderer).children)

	assert.Equal(t, visibleCount, lastItemIndex-firstItemIndex+1)

	lentry.Length = func() int {
		return 0
	}
	lentry.Refresh()

	visibleCount = len(test.WidgetRenderer(lentry).(*lentryRenderer).children)

	assert.Equal(t, visibleCount, 0)

	test.AssertRendersToMarkup(t, "lentry/cleared.xml", w.Canvas())
}

func TestLentry_RemoveItem(t *testing.T) {
	test.NewApp()
	defer test.NewApp()

	var data []string
	data = append(data, "Test Item 0")
	data = append(data, "Test Item 1")
	data = append(data, "Test Item 2")

	lentry := NewLentry(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(id LentryItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	w := test.NewWindow(lentry)
	w.Resize(fyne.NewSize(200, 400))

	visibleCount := len(test.WidgetRenderer(lentry).(*lentryRenderer).children)
	assert.Equal(t, visibleCount, 3)

	data = data[:len(data)-1]
	lentry.Refresh()

	visibleCount = len(test.WidgetRenderer(lentry).(*lentryRenderer).children)
	assert.Equal(t, visibleCount, 2)
	test.AssertRendersToMarkup(t, "lentry/item_removed.xml", w.Canvas())
}

func TestLentry_NoFunctionsSet(t *testing.T) {
	lentry := &Lentry{}
	w := test.NewWindow(lentry)
	w.Resize(fyne.NewSize(200, 400))
	lentry.Refresh()
}

func createLentry(items int) *Lentry {
	var data []string
	for i := 0; i < items; i++ {
		data = append(data, fmt.Sprintf("Test Item %d", i))
	}

	lentry := NewLentry(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(layout.NewHBoxLayout(), NewIcon(theme.DocumentIcon()), NewLabel("Template Object"))
		},
		func(id LentryItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
		},
	)
	lentry.Resize(fyne.NewSize(200, 1000))
	return lentry
}

func changeLentryData(lentry *Lentry) {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	lentry.Length = func() int {
		return len(data)
	}
	lentry.UpdateItem = func(id LentryItemID, item fyne.CanvasObject) {
		item.(*fyne.Container).Objects[1].(*Label).SetText(data[id])
	}
}

func setupLentry(t *testing.T) (*Lentry, fyne.Window) {
	test.NewApp()
	lentry := createLentry(1000)
	w := test.NewWindow(lentry)
	w.Resize(fyne.NewSize(200, 400))
	test.AssertRendersToMarkup(t, "lentry/initial.xml", w.Canvas())
	return lentry, w
}
