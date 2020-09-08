package ui


import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/Z-Bolt/OctoScreen/utils"
)

type SplashPanel struct {
	CommonPanel
	Label			*gtk.Label
	RetryButton		*gtk.Button
}

func NewSplashPanel(ui *UI) *SplashPanel {
	instane := &SplashPanel {
		CommonPanel: NewCommonPanel(ui, nil),
	}
	instane.initialize()

	return instane
}

func (this *SplashPanel) initialize() {
	logo := utils.MustImageFromFile("logos/logo.png")
	this.Label = utils.MustLabel("...")
	this.Label.SetHExpand(true)
	this.Label.SetLineWrap(true)
	this.Label.SetMaxWidthChars(30)
	this.Label.SetText("Initializing printer...")

	main := utils.MustBox(gtk.ORIENTATION_VERTICAL, 15)
	main.SetVAlign(gtk.ALIGN_END)

	// main.SetVExpand(true)
	// main.SetHExpand(true)
	main.SetVExpand(false)
	main.SetHExpand(false)

	main.Add(logo)
	main.Add(this.Label)

	box := utils.MustBox(gtk.ORIENTATION_VERTICAL, 0)
	box.Add(main)
	box.Add(this.createActionBar())

	this.Grid().Add(box)
}

func (this *SplashPanel) createActionBar() gtk.IWidget {
	actionBar := utils.MustBox(gtk.ORIENTATION_HORIZONTAL, 5)
	actionBar.SetHAlign(gtk.ALIGN_END)

	this.RetryButton = utils.MustButtonImageStyle("Retry", "refresh.svg", "color2", this.releaseFromHold)
	this.RetryButton.SetProperty("width-request", this.Scaled(100))
	this.RetryButton.SetProperty("visible", true)
	actionBar.Add(this.RetryButton)
	ctx, _ := this.RetryButton.GetStyleContext()
	ctx.AddClass("hidden")

	systemButton := utils.MustButtonImageStyle("System", "info.svg", "color3", this.showSystem)
	systemButton.SetProperty("width-request", this.Scaled(100))
	actionBar.Add(systemButton)

	networkButton := utils.MustButtonImageStyle("Network", "network.svg", "color4", this.showNetwork)
	networkButton.SetProperty("width-request", this.Scaled(100))
	actionBar.Add(networkButton)

	return actionBar
}

func (this *SplashPanel) putOnHold() {
	this.RetryButton.Show()
	ctx, _ := this.RetryButton.GetStyleContext()
	ctx.RemoveClass("hidden")
	this.Label.SetText("Cannot connect to the printer.  Tap \"Retry\" to try again.")
}

func (this *SplashPanel) releaseFromHold() {
	this.RetryButton.Hide()
	ctx, _ := this.RetryButton.GetStyleContext()
	ctx.AddClass("hidden")

	this.Label.SetText("Loading...")
	this.UI.connectionAttempts = 0
}

func (this *SplashPanel) showNetwork() {
	this.UI.Add(NetworkPanel(this.UI, this))
}

func (this *SplashPanel) showSystem() {
	this.UI.Add(SystemPanel(this.UI, this))
}
