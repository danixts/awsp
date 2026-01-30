package gatewaytui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danixts/awsp/internal/apigateway"
	"github.com/danixts/awsp/internal/store"
)

type focusPanel int

const (
	focusAPIs focusPanel = iota
	focusResources
)

const (
	panelMinWidth  = 24
	panelMinHeight = 10
	defaultSince   = "1h"

	modalNone = iota
	modalTimeSelect
	modalLogView
	modalThemeSelect
)

func panelBorderStyle(activeBorderColor, inactiveBorderColor string) (active, inactive lipgloss.Style) {
	active = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(activeBorderColor)).
		Padding(0, 0)
	inactive = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(inactiveBorderColor)).
		Padding(0, 0)
	return active, inactive
}

type model struct {
	width  int
	height int

	apiList      list.Model
	resourceList list.Model
	logViewport  viewport.Model

	focus focusPanel

	apis      []apigateway.RestAPI
	endpoints []apigateway.Endpoint
	selectedAPI *apigateway.RestAPI
	profile   string
	region    string

	cachedAPIs      []apigateway.RestAPI
	cachedResources map[string][]apigateway.ResourceWithMethods
	cachedLogGroups map[string]string

	loading  string
	err      string
	quitting bool
	spinner  spinner.Model

	modalState      int
	modalTimeList   list.Model
	modalThemeList  list.Model
	modalLogView    viewport.Model
	modalContent    string
	pendingLogGroup string
	pendingSince    string
	loadingForModal bool
	logStreamCh     <-chan string
	themeIndex      int
}

type logStreamStartedMsg struct {
	ch <-chan string
}

type logLineMsg struct {
	line string
}

type logStreamDoneMsg struct{}

type apisLoadedMsg struct {
	apis []apigateway.RestAPI
	err  error
}

type resourcesLoadedMsg struct {
	apiID    string
	resources []apigateway.ResourceWithMethods
	err      error
}

type logsLoadedMsg struct {
	content string
	err     error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.WindowSize(), loadAPIsCmd, func() tea.Msg { return m.spinner.Tick() })
}

func loadAPIsCmd() tea.Msg {
	apis, err := apigateway.GetRestAPIs()
	return apisLoadedMsg{apis: apis, err: err}
}

func loadResourcesCmd(apiID string) tea.Cmd {
	return func() tea.Msg {
		resources, err := apigateway.GetResources(apiID)
		return resourcesLoadedMsg{apiID: apiID, resources: resources, err: err}
	}
}

func loadLogsStreamCmd(logGroup, since, region string) tea.Cmd {
	return func() tea.Msg {
		ready := make(chan (<-chan string), 1)
		go func() {
			ready <- apigateway.TailLogsStream(logGroup, since, region)
		}()
		return logStreamStartedMsg{ch: <-ready}
	}
}

func waitForLogLineCmd(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-ch
		if ok {
			return logLineMsg{line: line}
		}
		return logStreamDoneMsg{}
	}
}

func newModel(profile, region string) model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	apiList := list.New([]list.Item{}, delegate, panelMinWidth, panelMinHeight)
	apiList.Title = " APIs "
	apiList.SetShowStatusBar(false)
	apiList.SetFilteringEnabled(true)
	apiList.DisableQuitKeybindings()

	resourceDelegate := endpointDelegate{DefaultDelegate: list.NewDefaultDelegate()}
	resourceDelegate.ShowDescription = false
	resourceDelegate.SetSpacing(0)
	resourceList := list.New([]list.Item{}, resourceDelegate, panelMinWidth, panelMinHeight)
	resourceList.Title = " Endpoints "
	resourceList.SetShowStatusBar(false)
	resourceList.SetFilteringEnabled(true)
	resourceList.DisableQuitKeybindings()

	vp := viewport.New(panelMinWidth, panelMinHeight)
	vp.Style = lipgloss.NewStyle()

	timeDelegate := list.NewDefaultDelegate()
	timeDelegate.ShowDescription = false
	timeDelegate.SetSpacing(0)
	timeItems := timeRangeItemsFrom(apigateway.DefaultTimeRanges)
	modalTimeList := list.New(timeItems, timeDelegate, 32, 14)
	modalTimeList.Title = " Time range "
	modalTimeList.SetShowStatusBar(false)
	modalTimeList.DisableQuitKeybindings()

	modalVp := viewport.New(60, 18)
	modalVp.Style = lipgloss.NewStyle()

	themeDelegate := list.NewDefaultDelegate()
	themeDelegate.ShowDescription = false
	themeDelegate.SetSpacing(0)
	themeIndex := 0
	if name, _ := store.LoadGatewayTheme(); name != "" {
		themeIndex = themeIndexByName(name)
	}
	s := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))))
	modalThemeList := list.New(themeItemsFrom(themes), themeDelegate, 28, 12)
	modalThemeList.Title = " Theme "
	modalThemeList.SetShowStatusBar(false)
	modalThemeList.DisableQuitKeybindings()
	modalThemeList.Select(themeIndex)

	return model{
		apiList:         apiList,
		resourceList:    resourceList,
		logViewport:     vp,
		modalTimeList:   modalTimeList,
		modalThemeList:  modalThemeList,
		modalLogView:    modalVp,
		spinner:         s,
		focus:           focusAPIs,
		profile:         profile,
		region:          region,
		themeIndex:      themeIndex,
		loading:         "Loading API Gateways...",
		cachedResources: make(map[string][]apigateway.ResourceWithMethods),
		cachedLogGroups: make(map[string]string),
	}
}
