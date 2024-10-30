package pomodoro

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/1x-eng/tomatick/pkg/ltm"

	"github.com/1x-eng/tomatick/pkg/llm"
	"github.com/1x-eng/tomatick/pkg/markdown"

	"github.com/1x-eng/tomatick/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chzyer/readline"
	"github.com/logrusorgru/aurora"

	"github.com/1x-eng/tomatick/pkg/context"
	"github.com/1x-eng/tomatick/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var commandInstructions = []struct {
	cmd  string
	desc string
}{
	{"Type a task", "Add a new task to your list"},
	{"done", "Finish adding tasks and start the timer"},
	{"list", "Display your current task list"},
	{"edit N text", "Edit task number N with new text"},
	{"remove N", "Remove task number N from the list"},
	{"suggest", "Get AI-powered task suggestions"},
	{"flush", "Clear any existing in-memory AI suggestions"},
	{"help", "Show this help message"},
	{"quit", "End the session and save progress"},
}

type TomatickMemento struct {
	cfg                      *config.Config
	memClient                *ltm.MemAI
	llmClient                *llm.PerplexityAI
	memID                    string
	cycleCount               int
	cyclesSinceLastLongBreak int
	auroraInstance           aurora.Aurora
	sessionContext           string
	theme                    *ui.Theme
	currentSuggestions       []string
	currentTasks             []string
	lastAnalysis             string
}

func NewTomatickMemento(cfg *config.Config) *TomatickMemento {
	return &TomatickMemento{
		cfg:                      cfg,
		memClient:                ltm.NewMemAI(cfg),
		llmClient:                llm.NewPerplexityAI(cfg),
		cycleCount:               0,
		cyclesSinceLastLongBreak: 0,
		auroraInstance:           aurora.NewAurora(true),
		theme:                    ui.NewTheme(),
		currentSuggestions:       make([]string, 0),
	}
}

func (p *TomatickMemento) StartCycle() {
	if p.cycleCount == 0 {
		displayWelcomeMessage(p.auroraInstance)

		contextManager := context.NewContextManager(p.cfg.ContextDir, p.auroraInstance, p.theme)

		sessionContext, err := contextManager.GetSessionContext()
		if err != nil {
			fmt.Println(p.auroraInstance.Red("Error getting context:"), err)
		} else {
			p.sessionContext = sessionContext

			// Confirm context collection
			fmt.Println(p.theme.Styles.Subtitle.Render("\n✓ Context collected successfully"))
			fmt.Println(p.theme.Styles.InfoText.Render("Copilot initialized with your session context"))
			fmt.Println()
		}
	}

	for {
		p.runTomatickMementoCycle()

		if p.cyclesSinceLastLongBreak >= (p.cfg.CyclesBeforeLongBreak - 1) {
			p.takeLongBreak()
			p.cyclesSinceLastLongBreak = 0
		} else {
			if p.cyclesSinceLastLongBreak < p.cfg.CyclesBeforeLongBreak {
				p.takeShortBreak()
			}
			p.cyclesSinceLastLongBreak++
		}

		p.cycleCount++

		if !p.askToContinue() {
			fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightGreen(("\nTomatick workday completed. Goodbye!"))))
			p.printTotalHoursWorked()
			break
		}
	}
}

func (p *TomatickMemento) askToContinue() bool {
	continuePrompt := &survey.Confirm{
		Message: p.auroraInstance.BrightBlue("Would you like to start another Tomatick cycle?").String(),
	}
	var answer bool
	survey.AskOne(continuePrompt, &answer)
	return answer
}

func (p *TomatickMemento) createAndSetMemID() {
	memTitle := fmt.Sprintf("# Tomatick Workday | %s\n#workday #tomatick\n", time.Now().Format("02-01-2006"))
	memID, err := p.memClient.CreateMem(memTitle)

	if err != nil {
		fmt.Println(p.auroraInstance.Bold(p.auroraInstance.Red("Error creating MemAI entry: ")), err)
		return
	}

	p.memID = memID
}

func (p *TomatickMemento) asyncAppendToMem(cycleSummary string) {
	_, err := p.memClient.AppendToMem(p.memID, cycleSummary)

	if err != nil {
		fmt.Println(p.auroraInstance.Bold(p.auroraInstance.Red("Error appending to MemAI: ")), err)
	}

}

func (p *TomatickMemento) runTomatickMementoCycle() {
	if p.memID == "" {
		p.createAndSetMemID()
	}

	tasks := p.captureTasks()
	p.startTimer(p.cfg.TomatickMementoDuration, p.auroraInstance.Italic(p.auroraInstance.BrightRed("Tick Tock Tick Tock...")).String())
	p.playSound()

	completedTasks := p.markTasksComplete(tasks)
	reflections := p.captureReflections()

	// Initialize the spinner
	spinner := ui.NewSpinner(p.theme.Styles.Spinner.
		Foreground(lipgloss.Color("#C4B5FD")).
		Bold(true))
	done := make(chan bool)

	// Start spinner in a goroutine
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s Analyzing reflections...", spinner.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Perform AI analysis
	assistant := llm.NewAssistant(p.llmClient, p.sessionContext)
	analysis, err := assistant.AnalyzeProgress(strings.Split(completedTasks, "\n"), reflections)

	// Stop the spinner
	done <- true
	fmt.Print("\r") // Clear spinner line

	if err != nil {
		fmt.Println(p.auroraInstance.Red("Error getting AI analysis:"), err)
	} else {
		presenter := ui.NewAnalysisPresenter(p.theme)
		formattedAnalysis := presenter.Present(analysis)
		fmt.Println(formattedAnalysis)
		p.lastAnalysis = analysis
	}

	cycleSummary := markdown.FormatCycleSummary(completedTasks, reflections)
	if analysis != "" {
		cycleSummary += "\n### Copilot's Analysis\n" + analysis + "\n*\n"
	}

	go p.asyncAppendToMem(cycleSummary)
}

func (p *TomatickMemento) captureTasks() []string {
	header := p.theme.Styles.Title.Render("=== Task Entry Mode ===")
	var sb strings.Builder
	sb.WriteString("\n───────────────────────────────────────────────\n")
	sb.WriteString("                Instructions\n")
	sb.WriteString("───────────────────────────────────────────────\n\n")

	for _, cmd := range commandInstructions {
		sb.WriteString(fmt.Sprintf("• %s: %s\n", cmd.cmd, cmd.desc))
	}

	sb.WriteString("\n───────────────────────────────────────────────")
	instructions := p.theme.Styles.SystemInstruction.Render(sb.String())
	fmt.Println(p.theme.Styles.Subtitle.Render(header + "\n" + instructions))

	assistant := llm.NewAssistant(p.llmClient, p.sessionContext)
	var tasks []string
	rl, _ := readline.New(p.auroraInstance.BrightGreen("➤ ").String())
	defer rl.Close()

	for {
		p.displayTasks(tasks)
		input, _ := rl.Readline()
		input = strings.TrimSpace(input)

		switch strings.ToLower(input) {
		case "done":
			if len(tasks) == 0 {
				fmt.Println(p.auroraInstance.Red("❗ Please add at least one task before finishing."))
				continue
			}
			return tasks
		case "suggest":
			spinner := ui.NewSpinner(p.theme.Styles.Spinner.
				Foreground(lipgloss.Color("#C4B5FD")).
				Bold(true))
			done := make(chan bool)

			// Start spinner in goroutine
			go func() {
				for {
					select {
					case <-done:
						return
					default:
						fmt.Printf("\r%s Getting suggestions...", spinner.Next())
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()

			suggestions, err := assistant.GetTaskSuggestions(tasks, p.lastAnalysis)
			done <- true
			fmt.Print("\r") // Clear spinner line

			if err != nil {
				fmt.Println(p.auroraInstance.Red("❗ Error getting suggestions:"), err)
				continue
			}
			p.currentSuggestions = suggestions // Store suggestions
			p.displaySuggestions(suggestions)
		case "flush":
			p.FlushSuggestions()
		case "quit":
			fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightGreen("Session ended. Goodbye!")))
			os.Exit(0)
		case "help":
			p.displayHelp()
		case "list":
			continue
		case "":
			fmt.Println(p.auroraInstance.Red("❗ Task cannot be empty. Please try again."))
		default:
			if strings.HasPrefix(input, "edit ") {
				p.editTask(&tasks, input)
			} else if strings.HasPrefix(input, "remove ") {
				p.removeTask(&tasks, input)
			} else if strings.HasPrefix(input, "use ") {
				p.useSuggestion(&tasks, input)
			} else {
				tasks = append(tasks, input)
				fmt.Println(p.auroraInstance.Green("✓ Task added successfully."))
			}
		}
	}
}

func (p *TomatickMemento) displaySuggestions(suggestions []string) {
	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightBlue("\n=== Copilot's Suggestions ===")))
	for i, suggestion := range suggestions {
		fmt.Printf("%s %s\n",
			p.theme.Styles.TaskNumber.Render(fmt.Sprintf("%d.", i+1)),
			p.theme.Styles.AIMessage.Render(suggestion))
	}
	fmt.Println(p.auroraInstance.Italic("\nTo use a suggestion, type 'use N' where N is the suggestion number."))
}

func (p *TomatickMemento) useSuggestion(tasks *[]string, input string) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		fmt.Println(p.auroraInstance.Red("❗ Invalid use command. Use 'use N'"))
		return
	}

	index, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println(p.auroraInstance.Red("❗ Invalid suggestion number."))
		return
	}

	// Convert to 0-based index
	index--

	if index < 0 || index >= len(p.currentSuggestions) {
		fmt.Println(p.auroraInstance.Red("❗ Invalid suggestion number. Please choose a number between 1 and"), len(p.currentSuggestions))
		return
	}

	// Add the selected suggestion to tasks
	*tasks = append(*tasks, p.currentSuggestions[index])
	fmt.Printf("%s %s\n",
		p.auroraInstance.Green("✓ Added suggestion to tasks:"),
		p.theme.Styles.TaskItem.Render(p.currentSuggestions[index]))
}

func (p *TomatickMemento) editTask(tasks *[]string, input string) {
	parts := strings.SplitN(input, " ", 3)
	if len(parts) != 3 {
		fmt.Println(p.auroraInstance.Red("❗ Invalid edit command. Use 'edit N new_task_description'"))
		return
	}
	index, err := strconv.Atoi(parts[1])
	if err != nil || index < 1 || index > len(*tasks) {
		fmt.Println(p.auroraInstance.Red("❗ Invalid task number. Please try again."))
		return
	}
	(*tasks)[index-1] = parts[2]
	fmt.Println(p.auroraInstance.Green("✓ Task updated successfully."))
}

func (p *TomatickMemento) removeTask(tasks *[]string, input string) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		fmt.Println(p.auroraInstance.Red("❗ Invalid remove command. Use 'remove N'"))
		return
	}
	index, err := strconv.Atoi(parts[1])
	if err != nil || index < 1 || index > len(*tasks) {
		fmt.Println(p.auroraInstance.Red("❗ Invalid task number. Please try again."))
		return
	}
	*tasks = append((*tasks)[:index-1], (*tasks)[index:]...)
	fmt.Println(p.auroraInstance.Green("✓ Task removed successfully."))
}

func (p *TomatickMemento) displayTasks(tasks []string) {
	fmt.Println()
	border := p.theme.Styles.Subtitle.Render("═══════════════════════════════════════")

	header := fmt.Sprintf("%s Current Tasks %s", p.theme.Emoji.TaskPending, p.theme.Emoji.TaskPending)
	fmt.Println(p.theme.Styles.Title.Render(header))
	fmt.Println(border)

	var sb strings.Builder
	if len(tasks) == 0 {
		sb.WriteString(p.theme.Styles.InfoText.Render("No tasks yet. Start typing to add tasks!"))
	} else {
		for i, task := range tasks {
			taskNum := p.theme.Styles.TaskNumber.Render(fmt.Sprintf("%d.", i+1))
			taskText := p.theme.Styles.TaskItem.Render(task)
			sb.WriteString(fmt.Sprintf("%s %s %s\n",
				p.theme.Emoji.TaskPending,
				taskNum,
				taskText))
		}
	}

	fmt.Println(p.theme.Styles.Subtitle.Render(sb.String()))
	fmt.Println(border)
	fmt.Println()
}

func (p *TomatickMemento) displayHelp() {
	header := fmt.Sprintf("\n%s  Available Commands %s", p.theme.Emoji.Help, p.theme.Emoji.Help)
	fmt.Println(p.theme.Styles.Title.Render(header))

	border := p.theme.Styles.Subtitle.Render(strings.Repeat("─", 50))
	fmt.Println(border)

	for _, cmd := range commandInstructions {
		fmt.Printf("%s %s: %s\n",
			p.theme.Emoji.TaskPending,
			p.theme.Styles.TaskNumber.Render(cmd.cmd),
			p.theme.Styles.InfoText.Render(cmd.desc))
	}

	fmt.Println(border)
}

func (p *TomatickMemento) markTasksComplete(tasks []string) string {
	fmt.Println(fmt.Sprintf("\n%s Progress Check %s", p.theme.Emoji.Analysis, p.theme.Emoji.Analysis))
	border := p.theme.Styles.Subtitle.Render(strings.Repeat("─", 50))
	fmt.Println(border)

	completed := make([]bool, len(tasks))
	for i, task := range tasks {
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("%s %s",
				p.theme.Emoji.TaskPending,
				p.theme.Styles.TaskItem.Render(task)),
		}
		survey.AskOne(prompt, &completed[i])

		// Immediate visual feedback
		status := p.theme.Emoji.TaskComplete
		if !completed[i] {
			status = p.theme.Emoji.TaskPending
		}
		fmt.Printf("%s %s\n",
			status,
			p.theme.Styles.InfoText.Render(task))
	}

	fmt.Println(border)

	var completedTasks []string
	for i, task := range tasks {
		if completed[i] {
			completedTasks = append(completedTasks, "- [x] "+task)
		} else {
			completedTasks = append(completedTasks, "- [ ] "+task)
		}
	}
	return strings.Join(completedTasks, "\n")
}

func (p *TomatickMemento) captureReflections() string {
	header := fmt.Sprintf("\n%s Reflection Time %s",
		p.theme.Emoji.Reflection,
		p.theme.Emoji.Reflection)

	fmt.Println(p.theme.Styles.Title.Render(header))
	fmt.Println(p.theme.Styles.InfoText.Render(
		"Share your thoughts on progress, challenges, and insights (type 'done' to finish):"))

	rl, err := readline.New(p.theme.Emoji.TaskPending + " ")
	if err != nil {
		fmt.Println("Error initializing readline:", err)
		return ""
	}
	defer rl.Close()

	var reflections []string
	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}
		if strings.ToLower(strings.TrimSpace(line)) == "done" {
			fmt.Println()
			break
		}
		reflections = append(reflections, line)
	}
	return strings.Join(reflections, "\n")
}

func (p *TomatickMemento) startTimer(duration time.Duration, message string) {
	model := ui.NewProgressModel(duration, message, p.theme)
	program := tea.NewProgram(model)
	if err := program.Start(); err != nil {
		fmt.Println("Error running timer:", err)
		return
	}
}

func (p *TomatickMemento) takeShortBreak() {
	message := fmt.Sprintf("\n%s Time for a refreshing break! %s\n%s Remember to stretch and rest your eyes %s",
		p.theme.Emoji.Break,
		p.theme.Emoji.Success,
		p.theme.Emoji.Timer,
		p.theme.Emoji.Break)

	p.startTimer(
		p.cfg.ShortBreakDuration,
		p.theme.Styles.InfoText.Render(message),
	)
	p.playSound()
}

func (p *TomatickMemento) takeLongBreak() {
	message := fmt.Sprintf("\n%s Excellent work! Time for a longer break %s\n%s Take a walk or do some light exercise %s",
		p.theme.Emoji.Success,
		p.theme.Emoji.Break,
		p.theme.Emoji.Timer,
		p.theme.Emoji.Break)

	p.startTimer(
		p.cfg.LongBreakDuration,
		p.theme.Styles.InfoText.Render(message),
	)
	p.playSound()
}

func (p *TomatickMemento) playSound() {
	soundPath := filepath.Join("assets", "softbeep.mp3")
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// Use 'afplay' on macOS
		cmd = exec.Command("afplay", soundPath)
	case "linux":
		// Use 'aplay' on Linux, works with WAV files. (may need to convert mp3 to wav?)
		cmd = exec.Command("aplay", soundPath)
	case "windows":
		// Use PowerShell command on Windows
		cmd = exec.Command("powershell", "-c", "(New-Object Media.SoundPlayer '"+soundPath+"').PlaySync();")
	}

	if cmd != nil {
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error playing sound:", err)
		}
	}
}

func (p *TomatickMemento) printTotalHoursWorked() {
	totalDuration := p.cfg.TomatickMementoDuration * time.Duration(p.cycleCount)
	totalHours := totalDuration.Hours()

	fmt.Println(p.theme.Styles.Title.Render("\n📊 Session Summary"))
	border := p.theme.Styles.Subtitle.Render(strings.Repeat("═", 50))
	fmt.Println(border)

	stats := []struct {
		label string
		value string
		emoji string
	}{
		{"Cycles Completed", fmt.Sprintf("%d", p.cycleCount), "🔄"},
		{"Hours Worked", fmt.Sprintf("%.2f hours", totalHours), "⏱️"},
	}

	for _, stat := range stats {
		fmt.Printf("%s %s: %s\n",
			stat.emoji,
			p.theme.Styles.TaskNumber.Render(stat.label),
			p.theme.Styles.SuccessText.Render(stat.value))
	}
	fmt.Println(border)

	workHoursSummary := fmt.Sprintf("#### Total Hours Worked: %.2f hours\n#### Total Cycles Completed: %d\n*",
		totalHours, p.cycleCount)
	p.asyncAppendToMem(workHoursSummary)
}

func displayWelcomeMessage(au aurora.Aurora) {
	asciiArt := `
	████████╗ ██████╗ ███╗   ███╗ █████╗ ████████╗██╗ ██████╗██╗  ██╗
	╚══██╔══╝██╔═══██╗████╗ ████║██╔══██╗╚══██╔══╝██║██╔════╝██║ ██╔╝
	   ██║   ██║   ██║██╔████╔██║███████║   ██║   ██║██║     █████╔╝ 
	   ██║   ██║   ██║██║╚██╔╝██║██╔══██║   ██║   ██║██║     ██╔═██╗ 
	   ██║   ╚██████╔╝██║ ╚═╝ ██║██║  ██║   ██║   ██║╚██████╗██║  ██╗
	   ╚═╝    ╚═══╝ ╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝╚═╝  ╚═╝
	`
	welcomeText := `
	🌟 Your Productivity Partner 🌟
	
	🎯 Focus Enhancement  |  🧠 Cognitive Optimization  |  📈 Progress Tracking
	`

	fmt.Println(au.Bold(au.BrightMagenta(asciiArt)))
	fmt.Println(au.Bold(au.BrightCyan(welcomeText)))
	fmt.Println(strings.Repeat("─", 80))
}

func (p *TomatickMemento) FlushSuggestions() {
	p.currentSuggestions = []string{}
	p.lastAnalysis = ""
	fmt.Println(p.auroraInstance.Green("✓ Copilot suggestions and analysis cache flushed successfully."))
}
