package context

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/logrusorgru/aurora"

	"github.com/1x-eng/tomatick/pkg/llm"
	"github.com/1x-eng/tomatick/pkg/ui"
)

type ContextManager struct {
	contextDir         string
	au                 aurora.Aurora
	presenter          *ui.ContextPresenter
	currentContextFile string
	llmClient          *llm.PerplexityAI
}

func NewContextManager(contextDir string, au aurora.Aurora, theme *ui.Theme, llmClient *llm.PerplexityAI) *ContextManager {
	return &ContextManager{
		contextDir: contextDir,
		au:         au,
		presenter:  ui.NewContextPresenter(theme),
		llmClient:  llmClient,
	}
}

func (cm *ContextManager) GetSessionContext(llmClient *llm.PerplexityAI) (string, error) {
	// Display the context menu
	fmt.Print(cm.presenter.PresentContextMenu())

	// Add a separator
	fmt.Println()

	var useContextFile bool
	prompt := &survey.Confirm{
		Message: cm.au.BrightBlue("Would you like to load an existing context?").String(),
	}
	survey.AskOne(prompt, &useContextFile)

	var context string
	var err error

	if useContextFile {
		context, err = cm.getContextFromFile()
	} else {
		context, err = cm.getContextFromInput()
	}

	if err != nil {
		return "", err
	}

	return context, nil
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
	cm.currentContextFile = selected

	content, err := os.ReadFile(filepath.Join(cm.contextDir, selected))
	if err != nil {
		return "", fmt.Errorf("failed to read context file: %w", err)
	}

	// Ask if user wants to add additional context
	var addDelta bool
	deltaPrompt := &survey.Confirm{
		Message: cm.au.BrightBlue("Would you like to add *any* additional context for this session?").String(),
	}
	survey.AskOne(deltaPrompt, &addDelta)

	if !addDelta {
		return string(content), nil
	}

	// Get delta context
	fmt.Print(cm.presenter.PresentDeltaContextInput())

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "done" {
			break
		}
		lines = append(lines, line)
	}

	deltaContext := strings.Join(lines, "\n")

	if deltaContext == "" {
		fmt.Println(cm.au.BrightYellow("No additional context provided. Proceeding with original context."))
		return string(content), nil
	}

	// Ask if delta should be appended to saved context
	var saveDelta bool
	savePrompt := &survey.Confirm{
		Message: "Would you like to append this additional context to the saved context file? Otherwise, the delta you've provided will be ephemeral and only available for this session.",
	}
	survey.AskOne(savePrompt, &saveDelta)

	if saveDelta {
		// Append to existing file
		updatedContent := string(content) + "\n\n=== Additional Context ===\n" + deltaContext
		if err := os.WriteFile(filepath.Join(cm.contextDir, selected), []byte(updatedContent), 0644); err != nil {
			fmt.Println(cm.au.Red("Failed to update context file:"), err)
			// Continue with session even if save fails
		}
		return updatedContent, nil
	}

	enrichedContext := string(content) + "\n\n=== Session Context ===\n" + deltaContext

	// Refine enriched context with copilot - if user wants to
	refinedContext, err := cm.RefineContext(enrichedContext, cm.llmClient)
	if err != nil {
		fmt.Println(cm.au.Red("\nError during context refinement. Proceeding with original context."))
		refinedContext = enrichedContext
	}

	var saveRefinedContext bool
	saveRefinedContextPrompt := &survey.Confirm{
		Message: "Would you like to save this refined context for future sessions?",
	}
	survey.AskOne(saveRefinedContextPrompt, &saveRefinedContext)

	if saveRefinedContext {
		if err := cm.saveContext(refinedContext); err != nil {
			fmt.Println(cm.au.Red("Failed to save context:"), err)
		}
	}

	return refinedContext, nil
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

	refinedContext, err := cm.RefineContext(context, cm.llmClient)
	if err != nil {
		fmt.Println(cm.au.Red("\nError during context refinement. Proceeding with original context."))
		refinedContext = context
	}

	var saveContext bool
	prompt := &survey.Confirm{
		Message: "Would you like to save this context for future sessions?",
	}
	survey.AskOne(prompt, &saveContext)

	if saveContext {
		if err := cm.saveContext(refinedContext); err != nil {
			fmt.Println(cm.au.Red("Failed to save context:"), err)
		}
	}

	return refinedContext, nil
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

func (cm *ContextManager) RefineContext(context string, llmClient *llm.PerplexityAI) (string, error) {
	fmt.Println(cm.presenter.PresentRefinementOption())

	var useRefinement bool
	prompt := &survey.Confirm{
		Message: cm.au.BrightBlue("Would you like your copilot to help fine-tune OR refine this context?").String(),
	}
	survey.AskOne(prompt, &useRefinement)

	if !useRefinement {
		return context, nil
	}

	// Show thinking spinner while initializing refinement
	spinner := ui.NewSpinner(cm.presenter.GetTheme().Styles.Spinner.
		Foreground(lipgloss.Color("#818CF8")).
		Bold(true))
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s Initializing context refinement...", spinner.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	refiner := llm.NewContextRefiner(llmClient, context)
	chat, err := refiner.StartRefinement()
	done <- true
	fmt.Print("\r\033[K") // Clear the spinner line

	if err != nil {
		fmt.Println(cm.au.Red("\nError starting context refinement. Proceeding with original context along with any additional context you provided."))
		return context, nil
	}

	refinedContext, err := cm.handleRefinementChat(chat, context)
	if err != nil {
		fmt.Println(cm.au.Red("\nError during context refinement. Proceeding with original context along with any additional context you provided."))
		return context, nil
	}

	return refinedContext, nil
}

func (cm *ContextManager) handleRefinementChat(chat *llm.RefinementChat, originalContext string) (string, error) {
	fmt.Print(cm.presenter.PresentRefinementStart())

	scanner := bufio.NewScanner(os.Stdin)
	refinementComplete := false
	var lastResponse string

	for !refinementComplete {
		// Show thinking indicator with new style
		fmt.Printf("\n%s %s\n",
			cm.au.BrightCyan("💭"),
			cm.presenter.GetTheme().Styles.ThinkingText.Render("Analyzing your context..."))

		// Show thinking spinner
		spinner := ui.NewSpinner(cm.presenter.GetTheme().Styles.Spinner.
			Foreground(lipgloss.Color("#818CF8")).
			Bold(true))
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				default:
					fmt.Printf("\r%s Analyzing context...", spinner.Next())
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()

		// If this is the first iteration, send empty string to get initial question
		// Otherwise, send the last user input
		response, err := chat.Chat(lastResponse)
		done <- true
		fmt.Print("\r\033[K")

		if err != nil {
			return originalContext, fmt.Errorf("chat error: %w", err)
		}

		// Check if refinement is complete
		if strings.Contains(strings.ToLower(response), "context refinement complete") {
			refinementComplete = true
			// Clean up the response to remove redundant header
			cleanResponse := response
			if idx := strings.Index(strings.ToLower(response), "here's what i've learned:"); idx != -1 {
				cleanResponse = response[idx+len("here's what i've learned:"):]
			}
			fmt.Printf("\n%s %s\n",
				cm.au.BrightCyan("🤖"),
				cm.presenter.GetTheme().Styles.AIMessage.Render(cleanResponse))
			break
		}

		// Format AI's response with new style
		fmt.Printf("\n%s %s\n",
			cm.au.BrightCyan("🤖"),
			cm.presenter.GetTheme().Styles.AIMessage.Render(response))

		// Format user input prompt with new style
		fmt.Printf("\n%s %s\n",
			cm.au.BrightBlue("ℹ️"),
			cm.presenter.GetTheme().Styles.SystemInstruction.Render("Your response:"))
		fmt.Printf("%s ", cm.au.BrightBlue("👤").String())

		var lines []string
		for scanner.Scan() {
			line := scanner.Text()
			if line == "done" || line == "exit" {
				break
			}
			lines = append(lines, line)
		}

		input := strings.TrimSpace(strings.Join(lines, "\n"))
		if input == "exit" {
			return originalContext, nil
		}

		if input == "" {
			continue
		}

		lastResponse = input
	}

	// Get and present the refined context
	spinner := ui.NewSpinner(cm.presenter.GetTheme().Styles.Spinner.
		Foreground(lipgloss.Color("#818CF8")).
		Bold(true))
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s Finalizing refined context...", spinner.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	refinedContext, err := chat.GetRefinedContext()
	done <- true
	fmt.Print("\r\033[K")

	if err != nil {
		return originalContext, fmt.Errorf("chat error: %w", err)
	}

	// Present the refined context for approval
	fmt.Print(cm.presenter.PresentRefinedContext(refinedContext))

	var approved bool
	for !approved {
		approvalPrompt := &survey.Confirm{
			Message: "Do you approve this refined context?",
		}
		survey.AskOne(approvalPrompt, &approved)

		if !approved {
			fmt.Print(cm.au.BrightBlue("What changes would you like to make? Type your changes and enter 'done' when finished:\n").String())

			var lines []string
			for scanner.Scan() {
				line := scanner.Text()
				if line == "done" {
					break
				}
				lines = append(lines, line)
			}

			userFeedback := strings.Join(lines, "\n")

			// Show thinking spinner
			spinner := ui.NewSpinner(cm.presenter.GetTheme().Styles.Spinner.
				Foreground(lipgloss.Color("#818CF8")).
				Bold(true))
			done := make(chan bool)

			go func() {
				for {
					select {
					case <-done:
						return
					default:
						fmt.Printf("\r%s Incorporating feedback...", spinner.Next())
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()

			refinedContext, err = chat.RequestContextModification(originalContext, userFeedback)
			done <- true
			fmt.Print("\r\033[K")

			if err != nil {
				return originalContext, fmt.Errorf("modification error: %w", err)
			}

			fmt.Print(cm.presenter.PresentRefinedContext(refinedContext))
		}
	}

	// Ask about persisting the refined context
	var persistRefinedContext bool
	persistPrompt := &survey.Confirm{
		Message: "Would you like to save this refined context permanently? By chosing 'yes', understand that, the original context of the chosen file will be overwritten. Choose 'no' to use it only for this session.",
	}
	survey.AskOne(persistPrompt, &persistRefinedContext)

	if persistRefinedContext {
		var filename string
		if cm.currentContextFile != "" {
			// If we loaded from a file, use that filename
			filename = cm.currentContextFile
		} else {
			filename = "copilot_refined_context_" + time.Now().Format("2006-01-02_15-04-05") + ".txt"
		}

		if err := os.WriteFile(filepath.Join(cm.contextDir, filename), []byte(refinedContext), 0644); err != nil {
			return originalContext, fmt.Errorf("failed to save refined context: %w", err)
		}

		fmt.Printf("\n%s Context saved to: %s\n",
			cm.presenter.GetTheme().Emoji.Success,
			cm.presenter.GetTheme().Styles.InfoText.Render(filename))
	}

	return refinedContext, nil
}
