# Tomatick Memento

Yet another pomodoro timer (yes, really) – but this one's your accountability buddy with just enough cognitive wizardry to keep perfectionists from disappearing down the rabbit hole of endless optimizations. Because sometimes "good enough" today beats "perfect" never. :brain: :rocket:

## Core Features

- **Interactive CLI Interface**: A terminal interface that doesn't suck, with colors and all that jazz

- **AI-Powered Task Management**:
  - Spots when you're overdoing it
  - Learns your work patterns
  - Stops you before you burn out
  - Keeps momentum going
  - Gets you into flow (and keeps you there)

- **Context-Aware Sessions**:
  - Tracks what you're actually doing
  - Matches tasks to when you're actually alert
  - Keeps tabs on progress
  - Tells you when you're taking on too much

- **Intelligent Task Suggestions**:
  - AI suggestions that make sense
  - Built for 25-minute chunks
  - Uses your past performance (good and bad)
  - Keeps you moving forward

- **Deep Performance Analysis**:
  - Tracks mental energy in real-time
  - Optimizes your flow state
  - Monitors energy levels
  - Shows actual progress

- **Sustainable Progress**:
  - Prevents you from burning out
  - Balances work and breaks
  - Tells you when to slow down
  - Helps you recharge properly

- **Persistent Memory Integration**:
  - Works with `mem.ai` because you'll forget
  - Clean markdown logs
  - Tracks everything important
  - Finds useful patterns


## Getting Started

### Prerequisites

- Go (version 1.15 or higher)
- An account with `mem.ai` and an API token
- Perplexity API token for AI features

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/1x-eng/tomatick.git
   ```

2. Navigate to project directory:
   ```bash
   cd tomatick
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Configure your environment:
   ```env
   POMODORO_DURATION=25m
   SHORT_BREAK_DURATION=5m
   LONG_BREAK_DURATION=15m
   CYCLES_BEFORE_LONGBREAK=4
   MEM_AI_API_TOKEN=your_mem_ai_api_token
   PERPLEXITY_API_TOKEN=your_perplexity_api_token
   ```

### Usage

1. Start the application:
   ```bash
   go run main.go
   ```

2. Follow the interactive prompts:
   - Provide session context for AI optimization
   - Add tasks or get AI-powered suggestions
   - Complete focused work sessions
   - Reflect on progress and receive AI analysis

3. Review your progress:
   - Session summaries in `mem.ai`
   - AI-powered performance analysis
   - Strategic recommendations for next sessions

## How It Works

Tomatick Memento uses advanced cognitive optimization algorithms to help you maintain peak performance while preventing burnout. The system:

- Analyzes your work patterns and adapts to your natural rhythms
- Sequences tasks based on cognitive demand and energy states
- Monitors mental load to prevent overextension
- Reduces decision overhead through intelligent suggestions
- Maintains productive momentum without the burnout
- Ensures consistent, sustainable progress

Every suggestion and analysis is precision-engineered to help you achieve maximum impact while maintaining optimal energy levels.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.


