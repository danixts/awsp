package gatewaytui

import (
	"context"

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
	panelMinWidth   = 24
	panelMinHeight  = 10
	defaultSince    = "1h"
	maxLogLinesTail = 500
)

const (
	modalNone = iota
	modalTimeSelect
	modalLogView
)

func PanelBorderStyle(activeBorderColor, inactiveBorderColor string) (active, inactive lipgloss.Style) {
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

	modalState       int
	modalTimeList    list.Model
	modalLogView     viewport.Model
	modalContent     string
	modalLogLines    []string
	pendingLogGroup  string
	pendingLogKey    string
	pendingSince     string
	loadingForModal  bool
	logStreamCh      <-chan string
	cancelLogStream  context.CancelFunc
	themeIndex       int
	currentMethod    string
	currentPath      string
}

type logGroupLoadedMsg struct {
	logGroup string
	logKey   string
	err      error
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
	preload  bool
}

type logsLoadedMsg struct {
	content string
	err     error
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.WindowSize(), loadAPIsCmd(), func() tea.Msg { return m.spinner.Tick() })
}

func loadAPIsCmd() tea.Cmd {
	return func() tea.Msg {
		apis, err := apigateway.GetRestAPIs()
		return apisLoadedMsg{apis: apis, err: err}
	}
}

func loadResourcesCmd(apiID string, preload bool) tea.Cmd {
	return func() tea.Msg {
		resources, err := apigateway.GetResources(apiID)
		return resourcesLoadedMsg{apiID: apiID, resources: resources, err: err, preload: preload}
	}
}

func loadLogGroupCmd(apiID, resourceID, method, logKey string) tea.Cmd {
	return func() tea.Msg {
		lg, err := apigateway.GetIntegrationLambdaLogGroup(apiID, resourceID, method)
		return logGroupLoadedMsg{logGroup: lg, logKey: logKey, err: err}
	}
}

func loadLogsStreamCmd(ctx context.Context, logGroup, since, region string) tea.Cmd {
	return func() tea.Msg {
		ready := make(chan (<-chan string), 1)
		go func() {
			ready <- apigateway.TailLogsStream(ctx, logGroup, since, region)
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
	apiDel := apiDelegate{DefaultDelegate: list.NewDefaultDelegate()}
	apiDel.ShowDescription = false
	apiDel.SetSpacing(0)

	themeIndex := 0
	if name, _ := store.LoadGatewayTheme(); name != "" {
		themeIndex = themeIndexByName(name)
	}
	t := &Themes[themeIndex]

	apiList := list.New([]list.Item{}, apiDel, panelMinWidth, panelMinHeight)
	apiList.Title = " APIs "
	apiList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
	apiList.SetShowStatusBar(false)
	apiList.SetFilteringEnabled(true)
	apiList.DisableQuitKeybindings()

	resourceDelegate := endpointDelegate{DefaultDelegate: list.NewDefaultDelegate()}
	resourceDelegate.ShowDescription = false
	resourceDelegate.SetSpacing(0)
	resourceList := list.New([]list.Item{}, resourceDelegate, panelMinWidth, panelMinHeight)
	resourceList.Title = " Endpoints "
	resourceList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
	resourceList.SetShowStatusBar(false)
	resourceList.SetFilteringEnabled(true)
	resourceList.DisableQuitKeybindings()

	vp := viewport.New(panelMinWidth, panelMinHeight)
	vp.Style = lipgloss.NewStyle()

	timeDel := timeRangeDelegate{DefaultDelegate: list.NewDefaultDelegate()}
	timeDel.ShowDescription = false
	timeDel.SetSpacing(0)
	timeItems := timeRangeItemsFrom(apigateway.DefaultTimeRanges)
	modalTimeList := list.New(timeItems, timeDel, 32, 14)
	modalTimeList.Title = " Time range "
	modalTimeList.SetShowStatusBar(false)
	modalTimeList.DisableQuitKeybindings()

	modalVp := viewport.New(60, 18)
	modalVp.Style = lipgloss.NewStyle()

	s := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading))))

	return model{
		apiList:         apiList,
		resourceList:    resourceList,
		logViewport:     vp,
		modalTimeList:   modalTimeList,
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
