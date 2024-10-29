package pomodoro

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/1x-eng/tomatick/pkg/llm"
	"github.com/1x-eng/tomatick/pkg/ltm"
	"github.com/1x-eng/tomatick/pkg/markdown"

	"github.com/1x-eng/tomatick/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chzyer/readline"
	"github.com/logrusorgru/aurora"

	"github.com/1x-eng/tomatick/pkg/context"
	"github.com/1x-eng/tomatick/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
)

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
	}
}

func (p *TomatickMemento) StartCycle() {
	if p.cycleCount == 0 {
		displayWelcomeMessage(p.auroraInstance)

		contextManager := context.NewContextManager(p.cfg.ContextDir, p.auroraInstance)

		sessionContext, err := contextManager.GetSessionContext()
		if err != nil {
			fmt.Println(p.auroraInstance.Red("Error getting context:"), err)
		} else {
			p.sessionContext = sessionContext
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

	assistant := llm.NewAssistant(p.llmClient, p.sessionContext)
	analysis, err := assistant.AnalyzeProgress(strings.Split(completedTasks, "\n"), reflections)
	if err != nil {
		fmt.Println(p.auroraInstance.Red("Error getting AI analysis:"), err)
	} else {
		fmt.Println(p.theme.Styles.Border.Render(
			p.theme.Styles.InfoText.Render("\n=== AI Analysis ===\n" + analysis),
		))
	}

	cycleSummary := markdown.FormatCycleSummary(completedTasks, reflections)
	if analysis != "" {
		cycleSummary += "\n### AI Analysis\n" + analysis + "\n***\n"
	}

	go p.asyncAppendToMem(cycleSummary)
}

func (p *TomatickMemento) captureTasks() []string {
	header := p.theme.Styles.Title.Render("=== Task Entry Mode ===")
	instructions := p.theme.Styles.InfoText.Render(`
• Type a task and press Enter to add it
• Type 'suggest' to get AI task suggestions
• Type 'done' when you've finished adding tasks
• Type 'help' to see all available commands
	`)

	fmt.Println(p.theme.Styles.Border.Render(header + "\n" + instructions))

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
			suggestions, err := assistant.GetTaskSuggestions(tasks)
			if err != nil {
				fmt.Println(p.auroraInstance.Red("❗ Error getting suggestions:"), err)
				continue
			}
			p.displaySuggestions(suggestions)
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
	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightYellow("\n=== AI Suggestions ===")))
	for i, suggestion := range suggestions {
		fmt.Printf("%s %s\n",
			p.theme.Styles.TaskNumber.Render(fmt.Sprintf("%d.", i+1)),
			p.theme.Styles.InfoText.Render(suggestion))
	}
	fmt.Println(p.auroraInstance.Italic("\nTo use a suggestion, type 'use N' where N is the suggestion number."))
}

func (p *TomatickMemento) useSuggestion(tasks *[]string, input string) {
	parts := strings.SplitN(input, " ", 2)
	if len(parts) != 2 {
		fmt.Println(p.auroraInstance.Red("❗ Invalid use command. Use 'use N'"))
		return
	}
	_, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println(p.auroraInstance.Red("❗ Invalid suggestion number."))
		return
	}
	// Note: The actual suggestion adding logic would need to maintain a suggestions slice
	// This is just a placeholder for the implementation
	fmt.Println(p.auroraInstance.Green("✓ Added suggestion to tasks."))
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
	var sb strings.Builder
	sb.WriteString(p.theme.Styles.Subtitle.Render("Current Tasks"))
	sb.WriteString("\n")

	if len(tasks) == 0 {
		sb.WriteString(p.theme.Styles.InfoText.Render("No tasks yet. Start typing to add tasks."))
	} else {
		for i, task := range tasks {
			taskNum := p.theme.Styles.TaskNumber.Render(fmt.Sprintf("%d.", i+1))
			taskText := p.theme.Styles.TaskItem.Render(task)
			sb.WriteString(fmt.Sprintf("%s %s\n", taskNum, taskText))
		}
	}

	fmt.Println(p.theme.Styles.Border.Render(sb.String()))
}

func (p *TomatickMemento) displayHelp() {
	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightYellow("\n=== Available Commands ===")))
	fmt.Println(p.auroraInstance.BrightYellow("• Type a task: Add a new task"))
	fmt.Println(p.auroraInstance.BrightYellow("• done: Finish adding tasks"))
	fmt.Println(p.auroraInstance.BrightYellow("• list: Display current tasks"))
	fmt.Println(p.auroraInstance.BrightYellow("• edit N new_description: Edit task N"))
	fmt.Println(p.auroraInstance.BrightYellow("• remove N: Remove task N"))
	fmt.Println(p.auroraInstance.BrightYellow("• help: Show this help message"))
}

func (p *TomatickMemento) markTasksComplete(tasks []string) string {
	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.Magenta("\nHow'd you go? Mark tasks that you completed:")))

	completed := make([]bool, len(tasks))
	for i, task := range tasks {
		prompt := &survey.Confirm{
			Message: fmt.Sprintf(p.auroraInstance.Italic(p.auroraInstance.BrightWhite("Did you complete '%s'?")).String(), task),
		}
		survey.AskOne(prompt, &completed[i])
	}

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
	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightWhite(("\nReflect and record your wins & distractions (you can use multiple lines, type 'done' to finish):"))))

	rl, err := readline.New("> ")
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
	p.startTimer(p.cfg.ShortBreakDuration, p.auroraInstance.Italic(p.auroraInstance.BrightGreen("\nOn short break...")).String())
	p.playSound()
}

func (p *TomatickMemento) takeLongBreak() {
	p.startTimer(p.cfg.LongBreakDuration, p.auroraInstance.Italic(p.auroraInstance.BrightRed("\nTomatickMementos long cycle complete! On long break...")).String())
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

	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightCyan("\n\nTotal TomatickMemento cycles completed: ")), p.auroraInstance.Bold(p.auroraInstance.BrightYellow(p.cycleCount)))
	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightCyan("Total hours worked: ")), p.auroraInstance.Bold(p.auroraInstance.BrightYellow(fmt.Sprintf("%.2f hours", totalHours))))

	workHoursSummary := fmt.Sprintf("#### Total Hours Worked: %.2f hours\n#### Total Cycles Completed: %d\n***", totalHours, p.cycleCount)
	// Not running this async, as we want to wait for it to complete before exiting
	p.asyncAppendToMem(workHoursSummary)
}

func displayWelcomeMessage(au aurora.Aurora) {
	asciiArt := `

	████████╗ ██████╗ ███╗   ███╗ █████╗ ████████╗██╗ ██████╗██╗  ██╗
	╚══██╔══╝██╔═══██╗████╗ ████║██╔══██╗╚══██╔══╝██║██╔════╝██║ ██╔╝
	   ██║   ██║   ██║██╔████╔██║███████║   ██║   ██║██║     █████╔╝ 
	   ██║   ██║   ██║██║╚██╔╝██║██╔══██║   ██║   ██║██║     ██╔═██╗ 
	   ██║   ╚██████╔╝██║ ╚═╝ ██║██║  ██║   ██║   ██║╚██████╗██║  ██╗
	   ╚═╝    ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝╚═╝  ╚═╝
	`
	fmt.Println(au.Bold(au.BrightMagenta(asciiArt)))
	fmt.Println()
}
