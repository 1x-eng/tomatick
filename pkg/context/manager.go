package context

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/logrusorgru/aurora"

	"github.com/1x-eng/tomatick/pkg/ui"
)

type ContextManager struct {
	contextDir string
	au         aurora.Aurora
	presenter  *ui.ContextPresenter
}

func NewContextManager(contextDir string, au aurora.Aurora, theme *ui.Theme) *ContextManager {
	return &ContextManager{
		contextDir: contextDir,
		au:         au,
		presenter:  ui.NewContextPresenter(theme),
	}
}

func (cm *ContextManager) GetSessionContext() (string, error) {
	// Display the context menu
	fmt.Print(cm.presenter.PresentContextMenu())

	// Add a separator
	fmt.Println()

	var useContextFile bool
	prompt := &survey.Confirm{
		Message: cm.au.BrightBlue("Would you like to load an existing context?").String(),
	}
	survey.AskOne(prompt, &useContextFile)

	if useContextFile {
		return cm.getContextFromFile()
	}

	return cm.getContextFromInput()
}

func (cm *ContextManager) getContextFromFile() (string, error) {
	files, err := os.ReadDir(cm.contextDir)
	if err != nil {
		return "", fmt.Errorf("failed to read context directory: %w", err)
	}

	var options []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			options = append(options, file.Name())
		}
	}

	// Display available contexts
	fmt.Print(cm.presenter.PresentContextList(options))

	if len(options) == 0 {
		return cm.getContextFromInput()
	}

	var selected string
	prompt := &survey.Select{
		Message: "Choose a context:",
		Options: options,
	}
	survey.AskOne(prompt, &selected)

	content, err := os.ReadFile(filepath.Join(cm.contextDir, selected))
	if err != nil {
		return "", fmt.Errorf("failed to read context file: %w", err)
	}

	return string(content), nil
}

func (cm *ContextManager) getContextFromInput() (string, error) {
	fmt.Print(cm.presenter.PresentContextInput())

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "done" {
			break
		}
		lines = append(lines, line)
	}

	context := strings.Join(lines, "\n")

	var saveContext bool
	prompt := &survey.Confirm{
		Message: "Would you like to save this context for future sessions?",
	}
	survey.AskOne(prompt, &saveContext)

	if saveContext {
		if err := cm.saveContext(context); err != nil {
			fmt.Println(cm.au.Red("Failed to save context:"), err)
		}
	}

	return context, nil
}

func (cm *ContextManager) saveContext(context string) error {
	var filename string
	prompt := &survey.Input{
		Message: "Enter a name for your context file (will be saved as .txt):",
	}
	survey.AskOne(prompt, &filename)

	if !strings.HasSuffix(filename, ".txt") {
		filename += ".txt"
	}

	filepath := filepath.Join(cm.contextDir, filename)
	return os.WriteFile(filepath, []byte(context), 0644)
}
