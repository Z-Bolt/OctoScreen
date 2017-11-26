package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

const EnvImagesFolder = "8TFT_IMAGES"

var ImagesFolder string

func init() {
	ImagesFolder = os.Getenv(EnvImagesFolder)
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	notebook, err := gtk.NotebookNew()
	if err != nil {
		panic(err)
	}

	AddPage(notebook, "test", ButtonsPage())
	//AddPage(notebook, "bar")
	win.Add(notebook)

	// Set the default window size.
	win.SetDefaultSize(480, 320)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func AddPage(nb *gtk.Notebook, title string, w gtk.IWidget) {
	l, _ := gtk.LabelNew(title)

	//box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	//box.Add(lab)

	nb.AppendPage(w, l)
}

// http://python-gtk-3-tutorial.readthedocs.io/en/latest/layout.html#grid
func ButtonsPage() gtk.IWidget {
	grid, _ := gtk.GridNew()

	grid.Attach(NewButtonImage("Status", "status.svg"), 1, 0, 1, 1)
	grid.Attach(NewButtonImage("Heat Up", "heat-up.svg"), 2, 0, 1, 1)
	grid.Attach(NewButtonImage("Move", "move.svg"), 3, 0, 1, 1)
	grid.Attach(NewButtonImage("Home", "home.svg"), 4, 0, 1, 1)
	grid.Attach(NewButtonImage("Extruct", "extruct.svg"), 1, 1, 1, 1)
	grid.Attach(NewButtonImage("HeatBed", "bed.svg"), 2, 1, 1, 1)
	grid.Attach(NewButtonImage("Fan", "fan.svg"), 3, 1, 1, 1)
	grid.Attach(NewButtonImage("Settings", "settings.svg"), 4, 1, 1, 1)
	return grid
}

func NewButtonImage(label, img string) gtk.IWidget {
	i, err := gtk.ImageNewFromFile(filepath.Join(ImagesFolder, img))
	if err != nil {
		panic(err)
	}

	b, _ := gtk.ButtonNewWithLabel(label)
	b.SetImage(i)
	b.SetAlwaysShowImage(true)
	b.SetImagePosition(gtk.POS_TOP)
	b.SetVExpand(true)
	b.SetHExpand(true)
	return b

}
