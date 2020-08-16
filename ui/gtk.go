package ui

import (
	//"log"
	"fmt"
	"path/filepath"

	"github.com/Z-Bolt/OctoScreen/utils"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// MustWindow returns a new gtk.Window, if error panics.
func MustWindow(windowType gtk.WindowType) *gtk.Window {
	win, err := gtk.WindowNew(windowType)
	if err != nil {
		panic(err)
	}

	win.SetResizable(false)

	return win
}

// MustGrid returns a new gtk.Grid, if error panics.
func MustGrid() *gtk.Grid {
	grid, err := gtk.GridNew()
	if err != nil {
		panic(err)
	}

	return grid
}

// MustBox returns a new gtk.Box, with the given configuration, if err panics.
func MustBox(orientation gtk.Orientation, spacing int) *gtk.Box {
	box, err := gtk.BoxNew(orientation, spacing)
	if err != nil {
		panic(err)
	}

	return box
}

// MustProgressBar returns a new gtk.ProgressBar, if err panics.
func MustProgressBar() *gtk.ProgressBar {
	progressBar, err := gtk.ProgressBarNew()
	if err != nil {
		panic(err)
	}

	return progressBar
}

// MustLabel returns a new gtk.Label, if err panics.
func MustLabel(format string, args ...interface{}) *gtk.Label {
	label, err := gtk.LabelNew("")
	if err != nil {
		panic(err)
	}

	label.SetMarkup(fmt.Sprintf(format, args...))

	return label
}

// LabelWithImage represents a gtk.Label with a image to the right.
type LabelWithImage struct {
	Label *gtk.Label
	*gtk.Box
}

// LabelImageSize default width and height of the image for a LabelWithImage
const LabelImageSize = 20

// MustLabelWithImage returns a new LabelWithImage based on a gtk.Box containing
// a gtk.Label with a gtk.Image, the image is scaled at LabelImageSize.
func MustLabelWithImage(imageFileName, format string, args ...interface{}) *LabelWithImage {
	label := MustLabel(format, args...)
	box := MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	box.Add(MustImageFromFileWithSize(imageFileName, LabelImageSize, LabelImageSize))
	box.Add(label)

	return &LabelWithImage{Label: label, Box: box}
}

// MustButtonImageStyle returns a new gtk.Button with the given label, image and clicked callback, if error panics.
func MustButtonImageStyle(label, imageFileName string, style string, clicked func()) *gtk.Button {
	button := MustButtonImage(label, imageFileName, clicked)
	ctx, _ := button.GetStyleContext()
	ctx.AddClass(style)

	return button
}

func MustButtonImage(label, imageFileName string, clicked func()) *gtk.Button {
	image := MustImageFromFile(imageFileName)
	button, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		panic(err)
	}

	button.SetImage(image)
	button.SetAlwaysShowImage(true)
	button.SetImagePosition(gtk.POS_TOP)
	button.SetVExpand(true)
	button.SetHExpand(true)

	if clicked != nil {
		button.Connect("clicked", clicked)
	}

	return button
}

func MustToggleButton(label string, imageFileName string, clicked func()) *gtk.ToggleButton {
	image := MustImageFromFile(imageFileName)
	button, err := gtk.ToggleButtonNewWithLabel(label)
	if err != nil {
		panic(err)
	}

	button.SetImage(image)
	button.SetAlwaysShowImage(true)
	button.SetImagePosition(gtk.POS_TOP)
	button.SetVExpand(true)
	button.SetHExpand(true)

	if clicked != nil {
		button.Connect("clicked", clicked)
	}

	return button
}

func MustButton(image *gtk.Image, clicked func()) *gtk.Button {
	button, err := gtk.ButtonNew()
	if err != nil {
		panic(err)
	}

	button.SetImage(image)
	button.SetImagePosition(gtk.POS_TOP)

	if clicked != nil {
		button.Connect("clicked", clicked)
	}

	return button
}

func MustButtonText(label string, clicked func()) *gtk.Button {
	button, err := gtk.ButtonNewWithLabel(label)
	if err != nil {
		panic(err)
	}

	if clicked != nil {
		button.Connect("clicked", clicked)
	}

	return button
}

func MustImageFromFileWithSize(imageFileName string, width, height int) *gtk.Image {
	if imageFileName == "" {
		Logger.Error("MustImageFromFileWithSize() - imageFileName is empty")
		//debug.PrintStack()			need to import "runtime/debug"
	}

	imageFilePath := imagePath(imageFileName)
	if !utils.FileExists(imageFilePath) {
		Logger.Error("MustImageFromFileWithSize() - imageFilePath is '" + imageFilePath + "', but doesn't exist")
		//debug.PrintStack()			need to import "runtime/debug"
	}

	p, err := gdk.PixbufNewFromFileAtScale(imageFilePath, width, height, true)
	if err != nil {
		Logger.Error(err)
	}

	image, err := gtk.ImageNewFromPixbuf(p)
	if err != nil {
		panic(err)
	}

	return image
}

// MustImageFromFile returns a new gtk.Image based on the given file, if error panics.
func MustImageFromFile(imageFileName string) *gtk.Image {
	if imageFileName == "" {
		Logger.Error("MustImageFromFile() - imageFileName is empty")
		//debug.PrintStack()			need to import "runtime/debug"
	}

	imageFilePath := imagePath(imageFileName)
	if !utils.FileExists(imageFilePath) {
		Logger.Error("MustImageFromFile() - imageFilePath is '" + imageFilePath + "', but doesn't exist")
		//debug.PrintStack()			need to import "runtime/debug"
	}

	image, err := gtk.ImageNewFromFile(imageFilePath)
	if err != nil {
		panic(err)
	}

	return image
}

// MustCSSProviderFromFile returns a new gtk.CssProvider for a given css file, if error panics.
func MustCSSProviderFromFile(css string) *gtk.CssProvider {
	p, err := gtk.CssProviderNew()
	if err != nil {
		panic(err)
	}

	if err := p.LoadFromPath(filepath.Join(StylePath, css)); err != nil {
		panic(err)
	}

	return p
}

func imagePath(imageFileName string) string {
	return filepath.Join(StylePath, ImageFolder, imageFileName)
}

// MustOverlay returns a new gtk.Overlay, if error panics.
func MustOverlay() *gtk.Overlay {
	o, err := gtk.OverlayNew()
	if err != nil {
		panic(err)
	}

	return o
}
