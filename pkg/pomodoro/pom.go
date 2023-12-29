package pomodoro

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/1x-eng/tomatick/pkg/ltm"
	"github.com/1x-eng/tomatick/pkg/markdown"

	"github.com/1x-eng/tomatick/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

const (
	TomatickMementosBeforeLongBreak = 4
)

type TomatickMemento struct {
	cfg                      *config.Config
	memClient                *ltm.MemAI
	memID                    string
	cycleCount               int
	cyclesSinceLastLongBreak int
}

func NewTomatickMemento(cfg *config.Config) *TomatickMemento {
	return &TomatickMemento{
		cfg:                      cfg,
		memClient:                ltm.NewMemAI(cfg),
		cycleCount:               0,
		cyclesSinceLastLongBreak: 0,
	}
}

func (p *TomatickMemento) StartCycle() {

	if p.cycleCount == 0 {
		displayWelcomeMessage()
	}

	for {
		if p.cyclesSinceLastLongBreak >= TomatickMementosBeforeLongBreak {
			p.takeLongBreak()
			p.cyclesSinceLastLongBreak = 0
		} else {
			p.runTomatickMementoCycle()
			p.cyclesSinceLastLongBreak++
		}

		if !p.askToContinue() {
			fmt.Println("Tomatick workday completed. Goodbye!")
			break
		}
		p.cycleCount++
	}
}

func (p *TomatickMemento) askToContinue() bool {
	continuePrompt := &survey.Confirm{
		Message: "Would you like to start another Tomatick cycle?",
	}
	var answer bool
	survey.AskOne(continuePrompt, &answer)
	return answer
}

func (p *TomatickMemento) runTomatickMementoCycle() {
	if p.memID == "" {
		p.memID = p.memClient.CreateMem(fmt.Sprintf("# Tomatick Workday | %s\n", time.Now().Format("02-01-2006")))
	}

	tasks := p.captureTasks()
	p.startTimer(p.cfg.TomatickMementoDuration, "Tick Tock Tick Tock...")

	completedTasks := p.markTasksComplete(tasks)
	reflections := p.captureReflections()
	cycleSummary := markdown.FormatCycleSummary(completedTasks, reflections)
	p.memClient.AppendToMem(p.memID, cycleSummary)

	p.takeShortBreak()
}

func (p *TomatickMemento) captureTasks() []string {
	fmt.Println("Remember the 80:20 rule and enter tasks that you plan to work on (one per line), type 'done' to finish:")
	scanner := bufio.NewScanner(os.Stdin)
	var tasks []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ToLower(line) == "done" {
			break
		}
		tasks = append(tasks, line)
	}
	return tasks
}

func (p *TomatickMemento) markTasksComplete(tasks []string) string {
	au := aurora.NewAurora(true)
	fmt.Println(au.Bold(au.Magenta("How'd you go? Mark tasks that you completed:")))

	completed := make([]bool, len(tasks))
	for i, task := range tasks {
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Did you complete '%s'?", task),
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
	fmt.Println("Reflect and record your wins & distractions:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func (p *TomatickMemento) startTimer(duration time.Duration, message string) {
	au := aurora.NewAurora(true)
	fmt.Println(au.Bold(au.Cyan(message)))

	p.progress(duration, au)

	playSound()
}

func (p *TomatickMemento) progress(duration time.Duration, au aurora.Aurora) {
	pBar := mpb.New(mpb.WithWidth(60))
	totalSeconds := int(duration.Seconds())
	bar := pBar.AddBar(int64(totalSeconds),
		mpb.PrependDecorators(
			decor.Name("Time left: "),
			decor.Elapsed(decor.ET_STYLE_GO, decor.WC{W: 5}),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	for i := 0; i < totalSeconds; i++ {
		bar.Increment()
		time.Sleep(time.Second)
	}

	pBar.Wait()
}

func (p *TomatickMemento) takeShortBreak() {
	fmt.Println("Starting short break...")
	p.startTimer(p.cfg.ShortBreakDuration, "On short break...")
}

func (p *TomatickMemento) takeLongBreak() {
	fmt.Println("4 TomatickMementos complete. Starting long break.")
	p.startTimer(p.cfg.LongBreakDuration, "On long break...")
}

func playSound() {
	// todo: Add sound playing logic here
}

func displayWelcomeMessage() {
	au := aurora.NewAurora(true)
	asciiArt := `

	████████╗ ██████╗ ███╗   ███╗ █████╗ ████████╗██╗ ██████╗██╗  ██╗
	╚══██╔══╝██╔═══██╗████╗ ████║██╔══██╗╚══██╔══╝██║██╔════╝██║ ██╔╝
	   ██║   ██║   ██║██╔████╔██║███████║   ██║   ██║██║     █████╔╝ 
	   ██║   ██║   ██║██║╚██╔╝██║██╔══██║   ██║   ██║██║     ██╔═██╗ 
	   ██║   ╚██████╔╝██║ ╚═╝ ██║██║  ██║   ██║   ██║╚██████╗██║  ██╗
	   ╚═╝    ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝╚═╝  ╚═╝
	`
	fmt.Println(au.Bold(au.BrightGreen(asciiArt)))
	fmt.Println("\n")
}
