package ltm

// LongTermMemory defines the interface for long-term memory storage
type LongTermMemory interface {
	// CreateMem creates a new memory entry and returns its ID
	CreateMem(content string) (string, error)
	
	// AppendToMem appends content to an existing memory entry
	AppendToMem(memID, content string) (string, error)
}

// NewLongTermMemory creates a new LongTermMemory implementation based on configuration
func NewLongTermMemory(cfg interface{ GetMemAIToken() string }) LongTermMemory {
	if token := cfg.GetMemAIToken(); token != "" {
		return NewMemAI(cfg)
	}
	return NewNoOpMemory()
}

// NoOpMemory is a no-op implementation of LongTermMemory
type NoOpMemory struct{}

// NewNoOpMemory creates a new NoOpMemory instance
func NewNoOpMemory() *NoOpMemory {
	return &NoOpMemory{}
}

// CreateMem implements LongTermMemory interface
func (n *NoOpMemory) CreateMem(content string) (string, error) {
	return "", nil
}

// AppendToMem implements LongTermMemory interface
func (n *NoOpMemory) AppendToMem(memID, content string) (string, error) {
	return "", nil
} 