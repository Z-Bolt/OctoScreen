package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/logger"
)

// MustWindow returns a new gtk.Window, if error panics.
func MustWindow(windowType gtk.WindowType) *gtk.Window {
	win, err := gtk.WindowNew(windowType)
	if err != nil {
		logger.LogError("PANIC!!! - MustWindow()", "gtk.WindowNew()", err)
		panic(err)
	}

	win.SetResizable(false)

	return win
}

// MustGrid returns a new gtk.Grid, if error panics.
func MustGrid() *gtk.Grid {
	grid, err := gtk.GridNew()
	if err != nil {
		logger.LogError("PANIC!!! - MustGrid()", "gtk.GridNew()", err)
		panic(err)
	}

	return grid
}

// MustBox returns a new gtk.Box, with the given configuration, if err panics.
func MustBox(orientation gtk.Orientation, spacing int) *gtk.Box {
	box, err := gtk.BoxNew(orientation, spacing)
	if err != nil {
		logger.LogError("PANIC!!! - MustBox()", "gtk.BoxNew()", err)
		panic(err)
	}

	return box
}

// MustProgressBar returns a new gtk.ProgressBar, if err panics.
func MustProgressBar() *gtk.ProgressBar {
	progressBar, err := gtk.ProgressBarNew()
	if err != nil {
		logger.LogError("PANIC!!! - MustProgressBar()", "gtk.ProgressBarNew()", err)
		panic(err)
	}

	return progressBar
}

// MustLabel returns a new gtk.Label, if err panics.
func MustLabel(format string, args ...interface{}) *gtk.Label {
	label, err := gtk.LabelNew("")
	if err != nil {
		logger.LogError("PANIC!!! - MustLabel()", "gtk.LabelNew()", err)
		panic(err)
	}

	label.SetMarkup(fmt.Sprintf(format, args...))

	return label
}

// MustLabelWithCssClass returns a stylized new gtk.Label, if err panics.
func MustLabelWithCssClass(format string, className string, args ...interface{}) *gtk.Label {
	label, err := gtk.LabelNew("")
	if err != nil {
		logger.LogError("PANIC!!! - MustLabelWithCssClass()", "gtk.LabelNew()", err)
		panic(err)
	}

	ctx, _ := label.GetStyleContext()
	ctx.AddClass(className)

	label.SetMarkup(fmt.Sprintf(format, args...))

	return label
}

// MustLabelWithCssClass returns a stylized new gtk.Label, if err panics.
func MustLabelWithCssClasses(format string, classNames []string, args ...interface{}) *gtk.Label {
	label, err := gtk.LabelNew("")
	if err != nil {
		logger.LogError("PANIC!!! - MustLabelWithCssClasses()", "gtk.LabelNew()", err)
		panic(err)
	}

	label.SetMarkup(fmt.Sprintf(format, args...))

	ctx, _ := label.GetStyleContext()
	for i := 0; i < len(classNames); i++ {
		ctx.AddClass(classNames[i])
	}

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
func MustButtonImageStyle(buttonlabel, imageFileName string, style string, clicked func()) *gtk.Button {
	button := MustButtonImage(buttonlabel, imageFileName, clicked)
	ctx, _ := button.GetStyleContext()
	ctx.AddClass(style)

	return button
}

func MustButtonImage(buttonlabel, imageFileName string, clicked func()) *gtk.Button {
	image := MustImageFromFile(imageFileName)
	button, err := gtk.ButtonNewWithLabel(buttonlabel)
	if err != nil {
		logger.LogError("PANIC!!! - MustButtonImage()", "gtk.ButtonNewWithLabel()", err)
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
		logger.LogError("PANIC!!! - MustToggleButton()", "gtk.ToggleButtonNewWithLabel()", err)
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
		logger.LogError("PANIC!!! - MustButton()", "gtk.ButtonNew()", err)
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
		logger.LogError("PANIC!!! - MustButtonText()", "gtk.ButtonNewWithLabel()", err)
		panic(err)
	}

	if clicked != nil {
		button.Connect("clicked", clicked)
	}

	return button
}

func MustImageFromFileWithSize(imageFileName string, width, height int) *gtk.Image {
	if imageFileName == "" {
		logger.Error("MustImageFromFileWithSize() - imageFileName is empty")
		//debug.PrintStack()			need to import "runtime/debug"
	}

	imageFilePath := imagePath(imageFileName)
	if !FileExists(imageFilePath) {
		logger.Error("MustImageFromFileWithSize() - imageFilePath is '" + imageFilePath + "', but doesn't exist")
		//debug.PrintStack()			need to import "runtime/debug"
	}

	pixbuf, err := gdk.PixbufNewFromFileAtScale(imageFilePath, width, height, true)
	if err != nil {
		logger.LogError("gtk.MustImageFromFileWithSize()", "PixbufNewFromFileAtScale()", err)
	}

	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		logger.LogError("PANIC!!! - MustImageFromFileWithSize()", "gtk.ImageNewFromPixbuf()", err)
		panic(err)
	}

	return image
}

// MustImageFromFile returns a new gtk.Image based on the given file, if error panics.
func MustImageFromFile(imageFileName string) *gtk.Image {
	if imageFileName == "" {
		logger.Error("MustImageFromFile() - imageFileName is empty")
		//debug.PrintStack()			need to import "runtime/debug"
	}

	imageFilePath := imagePath(imageFileName)
	if !FileExists(imageFilePath) {
		logger.Error("MustImageFromFile() - imageFilePath is '" + imageFilePath + "', but doesn't exist")
		//debug.PrintStack()			need to import "runtime/debug"
	}

	image, err := gtk.ImageNewFromFile(imageFilePath)
	if err != nil {
		logger.LogError("PANIC!!! - MustImageFromFile()", "gtk.ImageNewFromFile()", err)
		panic(err)
	}

	return image
}

func ImageFromUrl(imageUrl string) (*gtk.Image, error) {
	if imageUrl == "" {
		logger.Error("ImageFromUrl() - imageUrl is empty")
		return nil, errors.New("imageUrl is empty")
	}

	httpResponse, getErr:= http.Get(imageUrl)
	if getErr != nil {
		return nil, getErr
	}

	defer func() {
		io.Copy(ioutil.Discard, httpResponse.Body)
		httpResponse.Body.Close()
	}()

	buffer := new(bytes.Buffer)
	readLength, readErr := buffer.ReadFrom(httpResponse.Body)
	if readErr != nil {
		return nil, readErr
	} else if readLength < 1 {
		return nil, errors.New("bytes read was zero")
	}

	pixbufLoader, newPixbufLoaderErr := gdk.PixbufLoaderNew()
	if newPixbufLoaderErr != nil {
		return nil, newPixbufLoaderErr
	}
	defer pixbufLoader.Close()

	writeLength, writeErr := pixbufLoader.Write(buffer.Bytes())
	if writeErr != nil {
		return nil, writeErr
	} else if writeLength < 1 {
		return nil, errors.New("bytes written was zero")
	}

	pixbuf, _ := pixbufLoader.GetPixbuf()
	image, imageNewFromPixbufErr := gtk.ImageNewFromPixbuf(pixbuf)

	return image, imageNewFromPixbufErr
}


// MustCSSProviderFromFile returns a new gtk.CssProvider for a given css file, if error panics.
func MustCssProviderFromFile(cssFileName string) *gtk.CssProvider {
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		logger.LogError("PANIC!!! - MustCssProviderFromFile()", "gtk.CssProviderNew()", err)
		panic(err)
	}

	cssFilePath := cssFilePath(cssFileName)
	if err := cssProvider.LoadFromPath(cssFilePath); err != nil {
		logger.LogError("PANIC!!! - MustCssProviderFromFile()", "cssProvider.LoadFromPath()", err)
		panic(err)
	}

	return cssProvider
}

func cssFilePath(cssFileName string) string {
	octoScreenConfigInstance := GetOctoScreenConfigInstance()
	return filepath.Join(octoScreenConfigInstance.CssStyleFilePath, cssFileName)
}

func imagePath(imageFileName string) string {
	octoScreenConfigInstance := GetOctoScreenConfigInstance()
	return filepath.Join(octoScreenConfigInstance.CssStyleFilePath, ImageFolder, imageFileName)
}

// MustOverlay returns a new gtk.Overlay, if error panics.
func MustOverlay() *gtk.Overlay {
	overlay, err := gtk.OverlayNew()
	if err != nil {
		logger.LogError("PANIC!!! - MustOverlay()", "gtk.OverlayNew()", err)
		panic(err)
	}

	return overlay
}
