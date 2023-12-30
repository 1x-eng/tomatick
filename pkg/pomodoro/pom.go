package pomodoro

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

type TomatickMemento struct {
	cfg                      *config.Config
	memClient                *ltm.MemAI
	memID                    string
	cycleCount               int
	cyclesSinceLastLongBreak int
	auroraInstance           aurora.Aurora
}

func NewTomatickMemento(cfg *config.Config) *TomatickMemento {
	return &TomatickMemento{
		cfg:                      cfg,
		memClient:                ltm.NewMemAI(cfg),
		cycleCount:               0,
		cyclesSinceLastLongBreak: 0,
		auroraInstance:           aurora.NewAurora(true),
	}
}

func (p *TomatickMemento) StartCycle() {

	if p.cycleCount == 0 {
		displayWelcomeMessage(p.auroraInstance)
	}

	for {
		if p.cyclesSinceLastLongBreak >= p.cfg.CyclesBeforeLongBreak {
			p.takeLongBreak()
			p.cyclesSinceLastLongBreak = 0
		} else {
			p.runTomatickMementoCycle()
			p.cyclesSinceLastLongBreak++
		}

		if !p.askToContinue() {
			fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightGreen(("\nTomatick workday completed. Goodbye!"))))
			break
		}
		p.cycleCount++
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
	cycleSummary := markdown.FormatCycleSummary(completedTasks, reflections)

	go p.asyncAppendToMem(cycleSummary)

	p.takeShortBreak()
}

func (p *TomatickMemento) captureTasks() []string {
	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightYellow(("\nRemember the 80:20 rule and enter tasks that you plan to work on (one per line), type 'done' to finish:"))))

	scanner := bufio.NewScanner(os.Stdin)
	var tasks []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ToLower(line) == "done" {
			fmt.Println()
			break
		}
		tasks = append(tasks, line)
	}
	return tasks
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

	scanner := bufio.NewScanner(os.Stdin)
	var reflections []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ToLower(line) == "done" {
			fmt.Println()
			break
		}
		reflections = append(reflections, line)
	}
	return strings.Join(reflections, "\n")
}

func (p *TomatickMemento) startTimer(duration time.Duration, message string) {
	fmt.Println(p.auroraInstance.Bold(p.auroraInstance.BrightBlue(message)))

	p.progress(duration)
}

func (p *TomatickMemento) progress(duration time.Duration) {
	pBar := mpb.New(mpb.WithWidth(60))
	totalSeconds := int(duration.Seconds())
	bar := pBar.AddBar(int64(totalSeconds),
		mpb.PrependDecorators(
			decor.Name(p.auroraInstance.Bold(p.auroraInstance.BrightCyan("Time elapsed: ")).String()),
			decor.Elapsed(decor.ET_STYLE_GO, decor.WC{W: 5}),
		),
		mpb.AppendDecorators(decor.OnComplete(
			decor.Spinner(nil, decor.WC{W: 5}), p.auroraInstance.Bold(p.auroraInstance.BrightGreen("Done!")).String(),
		)),
	)

	for i := 0; i < totalSeconds; i++ {
		bar.Increment()
		time.Sleep(time.Second)
	}

	pBar.Wait()
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
