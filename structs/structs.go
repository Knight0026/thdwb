package structs

import (
	"net/url"
	"thdwb/mustard"
	profiler "thdwb/profiler"
)

type WebBrowser struct {
	Document       *HTMLDocument
	ActiveDocument *Document
	Documents      []*Document
	Viewport       *mustard.CanvasWidget
	History        *History
	Window         *mustard.Window
}

type HTMLDocument struct {
	Title       string
	RootElement *NodeDOM
	URL         *url.URL
	RawDocument string
	OffsetY     int
	Styles      []*StyleElement
	Profiler    *profiler.Profiler

	SelectedElement *NodeDOM
	DebugFlag       bool
}

type History struct {
	pages []*url.URL
}

func (history *History) PageCount() int {
	return len(history.pages)
}

func (history *History) Push(URL *url.URL) {
	history.pages = append(history.pages, URL)
}

func (history *History) Last() *url.URL {
	return history.pages[len(history.pages)-1]
}

func (history *History) Pop() {
	if len(history.pages) > 0 {
		history.pages = history.pages[:len(history.pages)-1]
	}
}

type Document struct {
	Title       string
	Path        string
	ContentType string
	RawDocument string
	DOM         *NodeDOM
}

type RenderBox struct {
	Node *NodeDOM

	Top  float64
	Left float64

	Width  float64
	Height float64

	MarginTop    float64
	MarginLeft   float64
	MarginRight  float64
	MarginBottom float64

	PaddingTop    float64
	PaddingLeft   float64
	PaddingRight  float64
	PaddingBottom float64
}

//NodeDOM "DOM Node Struct definition"
type NodeDOM struct {
	Element string `json:"element"`
	Content string `json:"content"`

	Children   []*NodeDOM   `json:"children"`
	Attributes []*Attribute `json:"attributes"`
	Style      *Stylesheet  `json:"style"`
	Parent     *NodeDOM     `json:"-"`
	RenderBox  *RenderBox   `json:"-"`

	NeedsReflow  bool `json:"-"`
	NeedsRepaint bool `json:"-"`

	Document *HTMLDocument `json:"-"`
}

func (node *NodeDOM) Attr(attrName string) string {
	for _, attribute := range node.Attributes {
		if attribute.Name == attrName {
			return attribute.Value
		}
	}

	return ""
}

func (node *NodeDOM) CalcPointIntersection(x, y float64) *NodeDOM {
	var intersectedNode *NodeDOM
	if x > float64(node.RenderBox.Left) &&
		x < float64(node.RenderBox.Left+node.RenderBox.Width) &&
		y > float64(node.RenderBox.Top) &&
		y < float64(node.RenderBox.Top+node.RenderBox.Height) {
		intersectedNode = node
	}

	for i := 0; i < len(node.Children); i++ {
		tempNode := node.Children[i].CalcPointIntersection(x, y)
		if tempNode != nil {
			intersectedNode = tempNode
		}
	}

	return intersectedNode
}

func (node NodeDOM) RequestRepaint() {
	node.NeedsRepaint = true

	for _, childNode := range node.Children {
		childNode.RequestRepaint()
	}
}

func (node NodeDOM) RequestReflow() {
	node.NeedsReflow = true

	for _, childNode := range node.Children {
		childNode.RequestReflow()
	}
}

//Resource "HTTP resource struct definition"
type Resource struct {
	Body        string
	ContentType string
	Code        int
	URL         *url.URL
}

//Attribute "Generic key:value attribute definition"
type Attribute struct {
	Name  string
	Value string
}

//Stylesheet "Stylesheet definition for DOM Nodes"
type Stylesheet struct {
	Color           *ColorRGBA
	BackgroundColor *ColorRGBA

	FontSize   float64
	FontWeight int

	Display  string
	Position string

	Width  float64
	Height float64
	Top    float64
	Left   float64
}

//StyleElement "hmtl <style> element"
type StyleElement struct {
	Selector string
	Style    *Stylesheet
}

//ColorRGBA "RGBA color model"
type ColorRGBA struct {
	R float64
	G float64
	B float64
	A float64
}
