package llm

import (
	"fmt"
)

type RefinementChat struct {
	perplexity *PerplexityAI
	history    []Message
}

func NewRefinementChat(p *PerplexityAI, initialMessages []Message) *RefinementChat {
	return &RefinementChat{
		perplexity: p,
		history:    initialMessages,
	}
}

func (rc *RefinementChat) Chat(userInput string) (string, error) {
	if userInput != "" {
		rc.history = append(rc.history, Message{
			Role:    "user",
			Content: userInput,
		})
	}

	if len(rc.history) <= 2 {
		initialResponse, err := rc.perplexity.GetResponse(rc.history)
		if err != nil {
			return "", err
		}
		rc.history = append(rc.history, Message{
			Role:    "assistant",
			Content: initialResponse,
		})
		return initialResponse, nil
	}

	response, err := rc.perplexity.GetResponse(rc.history)
	if err != nil {
		return "", err
	}

	rc.history = append(rc.history, Message{
		Role:    "assistant",
		Content: response,
	})

	return response, nil
}

func (rc *RefinementChat) GetRefinedContext() (string, error) {
	finalPrompt := `Based on our discussion, provide the complete fine-tuned, exhaustive and refined context. Format your response exactly as:
"Context refinement complete. Here's what I've learned:
• [point 1]
• [point 2]
• [point 3]
..."

Include ANY and ALL key details while maintaining temporal references where applicable. DO NOT OVERSIMPLIFY OR OVERSUMMARIZE. AIM TO BE AS DETAILED AS POSSIBLE.`

	response, err := rc.Chat(finalPrompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (rc *RefinementChat) RequestContextModification(originalContext, userFeedback string) (string, error) {
	modificationPrompt := fmt.Sprintf(`Previous context:
%s

User feedback for modifications:
%s

Based on this feedback, provide a refined and fine-tuned version of the context. 

YOU MUST FOLLOW THESE RULES:
1. Carefully analyze both the original context and user feedback & incorporate all user requested changes
2. Preserve all critical details from the original context
3. Maintain chronological accuracy and temporal relationships
4. Retain all specific requirements, constraints, and conditions
5. Resolve any contradictions between original context and new modifications
6. Add clarifications where ambiguity exists
7. Validate that no important information is lost during refinement
8. Maintain logical flow and coherence
9. Format the response exactly as:
"Context refinement complete. Here's what I've learned:
• [point 1]
• [point 2]
..."`, originalContext, userFeedback)

	return rc.Chat(modificationPrompt)
}
