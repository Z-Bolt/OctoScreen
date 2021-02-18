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
	utils.Logger.Debug("entering SplashPanel.initialize()")

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

	utils.Logger.Debug("leaving SplashPanel.initialize()")
}

func (this *SplashPanel) createActionBar() gtk.IWidget {
	utils.Logger.Debug("entering SplashPanel.createActionBar()")

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

	utils.Logger.Debug("leaving SplashPanel.createActionBar()")

	return actionBar
}

func (this *SplashPanel) putOnHold() {
	utils.Logger.Debug("entering SplashPanel.putOnHold()")

	this.RetryButton.Show()
	ctx, err := this.RetryButton.GetStyleContext()
	if err != nil {
		utils.LogError("SplashPanel.putOnHold()", "RetryButton.GetStyleContext()", err)
	} else {
		ctx.RemoveClass("hidden")
	}
	this.Label.SetText("Cannot connect to the printer.  Tap \"Retry\" to try again.")

	utils.Logger.Debug("leaving SplashPanel.putOnHold()")
}

func (this *SplashPanel) releaseFromHold() {
	utils.Logger.Debug("entering SplashPanel.releaseFromHold()")

	this.RetryButton.Hide()
	ctx, _ := this.RetryButton.GetStyleContext()
	ctx.AddClass("hidden")

	this.Label.SetText("Loading...")
	this.UI.connectionAttempts = 0

	utils.Logger.Debug("leaving SplashPanel.releaseFromHold()")
}

func (this *SplashPanel) showNetwork() {
	utils.Logger.Debug("entering SplashPanel.showNetwork()")

	this.UI.GoToPanel(NetworkPanel(this.UI, this))

	utils.Logger.Debug("leaving SplashPanel.showNetwork()")
}

func (this *SplashPanel) showSystem() {
	utils.Logger.Debug("entering SplashPanel.showSystem()")

	this.UI.GoToPanel(SystemPanel(this.UI, this))

	utils.Logger.Debug("leaving SplashPanel.showSystem()")
}
