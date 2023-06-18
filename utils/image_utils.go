package utils

import (
    "github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)


func ImageNewFromSvg(svg string) (*gtk.Image, error) {
	// pixBuff, err := gdk.PixbufNewFromDataOnly([]byte(svg))
	pixBuff, err := PixbufNewFromDataOnly([]byte(svg))
	if err !=  nil {
		return nil, err
	}
	
	image, err := gtk.ImageNewFromPixbuf(pixBuff)
	if err !=  nil {
		return nil, err
	}
	
	return image, nil
}



// TODO: these were lifted from the latest version of GOTK3.  Update GOTK3 and then remove these.
// PixbufNewFromDataOnly is a convenient alternative to PixbufNewFromData() and also a wrapper around gdk_pixbuf_new_from_data().
func PixbufNewFromDataOnly(pixbufData []byte) (*gdk.Pixbuf, error) {
	pixbufLoader, err := gdk.PixbufLoaderNew()
	if err != nil {
		return nil, err
	}
	// return pixbufLoader.WriteAndReturnPixbuf(pixbufData)
	return WriteAndReturnPixbuf(pixbufLoader, pixbufData)
}

// Convenient function like above for Pixbuf. Write data, close loader and return Pixbuf.
func WriteAndReturnPixbuf(v *gdk.PixbufLoader, data []byte) (*gdk.Pixbuf, error) {
	_, err := v.Write(data)
	if err != nil {
		return nil, err
	}

	if err := v.Close(); err != nil {
		return nil, err
	}

	return v.GetPixbuf()
}
