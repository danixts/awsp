package cloudwatch

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danixts/awsp/internal/adapter/awscli"
	"github.com/danixts/awsp/internal/domain"
	"github.com/danixts/awsp/internal/store"
)

var cw = &awscli.CloudWatchAdapter{Executor: &awscli.CLIExecutor{}}

type focusPanel int

const (
	focusLogGroups focusPanel = iota
	focusLogs
)

const (
	panelMinWidth   = 24
	panelMinHeight  = 10
	maxLogLinesTail = 500
)

const (
	modalNone = iota
	modalTimeSelect
	modalLogView
)

type model struct {
	width  int
	height int

	logGroupList list.Model
	logViewport  viewport.Model

	focus focusPanel

	logGroups       []domain.LogGroup
	selectedGroup   *domain.LogGroup
	profile         string
	region          string

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
	loadingForModal  bool
	logStreamCh      <-chan string
	cancelLogStream  context.CancelFunc
	themeIndex       int
	currentLogGroup  string
}

type logGroupsLoadedMsg struct {
	groups []domain.LogGroup
	err    error
}

type logStreamStartedMsg struct {
	ch <-chan string
}

type logLineMsg struct {
	line string
}

type logStreamDoneMsg struct{}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.WindowSize(), loadLogGroupsCmd(""), func() tea.Msg { return m.spinner.Tick() })
}

func loadLogGroupsCmd(prefix string) tea.Cmd {
	return func() tea.Msg {
		groups, err := cw.ListLogGroups(prefix)
		return logGroupsLoadedMsg{groups: groups, err: err}
	}
}

func loadLogsStreamCmd(ctx context.Context, logGroup, since, region string) tea.Cmd {
	return func() tea.Msg {
		ch, err := cw.TailLogs(ctx, logGroup, since, region)
		if err != nil {
			errCh := make(chan string, 1)
			errCh <- "Error: " + err.Error()
			close(errCh)
			return logStreamStartedMsg{ch: errCh}
		}
		return logStreamStartedMsg{ch: ch}
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

type logGroupItem struct {
	domain.LogGroup
}

func (g logGroupItem) Title() string {
	return g.Name
}

func (g logGroupItem) Description() string {
	if g.StoredBytes > 0 {
		return fmt.Sprintf("%.1f MB", float64(g.StoredBytes)/1024/1024)
	}
	return ""
}

func (g logGroupItem) FilterValue() string { return g.Name }

type timeRangeItem struct {
	Label string
	Since string
}

func (t timeRangeItem) Title() string       { return t.Label }
func (t timeRangeItem) Description() string { return "" }
func (t timeRangeItem) FilterValue() string { return t.Label }

type logGroupDelegate struct {
	list.DefaultDelegate
}

func (d logGroupDelegate) Height() int   { return 1 }
func (d logGroupDelegate) Spacing() int { return 0 }
func (d logGroupDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.DefaultDelegate.Update(msg, m)
}

type timeRangeDelegate struct {
	list.DefaultDelegate
}

func (d timeRangeDelegate) Height() int   { return d.DefaultDelegate.Height() }
func (d timeRangeDelegate) Spacing() int { return d.DefaultDelegate.Spacing() }
func (d timeRangeDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.DefaultDelegate.Update(msg, m)
}

func logGroupItemsFrom(groups []domain.LogGroup) []list.Item {
	out := make([]list.Item, 0, len(groups))
	for _, g := range groups {
		out = append(out, logGroupItem{g})
	}
	return out
}

func timeRangeItemsFrom(ranges []domain.TimeRange) []list.Item {
	out := make([]list.Item, 0, len(ranges))
	for _, r := range ranges {
		out = append(out, timeRangeItem{Label: r.Label, Since: r.Since})
	}
	return out
}

func newModel(profile, region string) model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.SetSpacing(0)

	themeIndex := 0
	if name, _ := store.LoadGatewayTheme(); name != "" {
		themeIndex = themeIndexByName(name)
	}
	t := &themes[themeIndex]

	logGroupList := list.New([]list.Item{}, delegate, panelMinWidth, panelMinHeight)
	logGroupList.Title = " Log Groups "
	logGroupList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
	logGroupList.SetShowStatusBar(true)
	logGroupList.SetFilteringEnabled(true)
	logGroupList.DisableQuitKeybindings()

	vp := viewport.New(panelMinWidth, panelMinHeight)
	vp.Style = lipgloss.NewStyle()

	timeDelegate := list.NewDefaultDelegate()
	timeDelegate.ShowDescription = false
	timeDelegate.SetSpacing(0)
	timeItems := timeRangeItemsFrom(domain.DefaultTimeRanges)
	modalTimeList := list.New(timeItems, timeDelegate, 32, 14)
	modalTimeList.Title = " Time range "
	modalTimeList.SetShowStatusBar(false)
	modalTimeList.DisableQuitKeybindings()

	modalVp := viewport.New(60, 18)
	modalVp.Style = lipgloss.NewStyle()

	s := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading))))

	return model{
		logGroupList: logGroupList,
		logViewport:  vp,
		modalTimeList: modalTimeList,
		modalLogView:  modalVp,
		spinner:       s,
		focus:         focusLogGroups,
		profile:       profile,
		region:        region,
		themeIndex:    themeIndex,
		loading:       "Loading log groups...",
	}
}
