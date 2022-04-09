package ui


import (
	"github.com/gotk3/gotk3/gtk"

	"github.com/Z-Bolt/OctoScreen/logger"
	"github.com/Z-Bolt/OctoScreen/utils"
)


type SplashPanel struct {
	CommonPanel
	Label			*gtk.Label
	RetryButton		*gtk.Button
}

func CreateSplashPanel(ui *UI) *SplashPanel {
	instance := &SplashPanel {
		CommonPanel: CreateCommonPanel("SplashPanel", ui),
	}
	instance.initialize()

	return instance
}

func (this *SplashPanel) initialize() {
	logger.TraceEnter("SplashPanel.initialize()")

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

	logger.TraceLeave("SplashPanel.initialize()")
}

func (this *SplashPanel) createActionBar() gtk.IWidget {
	logger.TraceEnter("SplashPanel.createActionBar()")

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

	logger.TraceLeave("SplashPanel.createActionBar()")
	return actionBar
}

func (this *SplashPanel) putOnHold() {
	logger.TraceEnter("SplashPanel.putOnHold()")

	this.RetryButton.Show()
	ctx, err := this.RetryButton.GetStyleContext()
	if err != nil {
		logger.LogError("SplashPanel.putOnHold()", "RetryButton.GetStyleContext()", err)
	} else {
		ctx.RemoveClass("hidden")
	}
	this.Label.SetText("Cannot connect to the printer.  Tap \"Retry\" to try again.")

	logger.TraceLeave("SplashPanel.putOnHold()")
}

func (this *SplashPanel) releaseFromHold() {
	logger.TraceEnter("SplashPanel.releaseFromHold()")

	this.RetryButton.Hide()
	ctx, _ := this.RetryButton.GetStyleContext()
	ctx.AddClass("hidden")

	this.Label.SetText("Loading...")
	this.UI.connectionAttempts = 0

	logger.TraceLeave("SplashPanel.releaseFromHold()")
}

func (this *SplashPanel) showNetwork() {
	logger.TraceEnter("SplashPanel.showNetwork()")

	this.UI.GoToPanel(GetNetworkPanelInstance(this.UI))

	logger.TraceLeave("SplashPanel.showNetwork()")
}

func (this *SplashPanel) showSystem() {
	logger.TraceEnter("SplashPanel.showSystem()")

	this.UI.GoToPanel(GetSystemPanelInstance(this.UI))

	logger.TraceLeave("SplashPanel.showSystem()")
}
