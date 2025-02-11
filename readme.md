# Tomatick Memento

Yet another pomodoro timer (yes, really) – but this one's your accountability buddy with just enough cognitive wizardry to keep perfectionists from disappearing down the rabbit hole of endless optimizations. Because sometimes "good enough" today beats "perfect" never. :brain: :rocket:

(& Tomatick ?? What the heck is that? Well, it's a combination of "tomato" and "tick" (as in a timer) 🤷 Dad joke? Maybe.)
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


## See it in Action

<div align="center">
  <a href="https://www.youtube.com/watch?v=UVKKza3Isu0" target="_blank">
    <img 
      src="https://img.youtube.com/vi/UVKKza3Isu0/maxresdefault.jpg" 
      alt="Tomatick Demo" 
      style="width:700px;"/>
  </a>
</div>

▶️ Watch a quick demo of Tomatick in action, to help get you started.


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

Tomatick Memento combines traditional pomodoro timing with data analysis to help optimize your work sessions. The system:

- Tracks work patterns to identify your productive periods
- Organizes tasks based on estimated cognitive load
- Provides alerts when signs of mental fatigue appear
- Offers task suggestions based on historical data
- Helps maintain consistent work intervals

The goal is to help you work effectively while avoiding exhaustion.

## Why Another Pomodoro Timer?

Tomatick addresses some common challenges in knowledge work:

- **Practical Over Perfect**: We focus on getting things done rather than perfect optimization
- **Reduced Context Switching**: Groups similar tasks together to maintain focus
- **Sustainable Pacing**: Helps prevent overwork through regular breaks and session monitoring
- **Simplified Decisions**: Reduces small decisions during work sessions
- **Progress Tracking**: Measures concrete task completion rather than just time spent

## Why mem.ai?

I'm not affiliated with mem.ai, but kinda like their product for good reasons:

### The Second Brain Synergy
- mem.ai isn't just another note-taking app - it's the ability to create a web out of your thoughts and ideas. Thanks to LLMs, you can interact with this knowledge web naturally, reducing cognitive load on your brain. While their product is still in alpha (they call it 'production grade', but, it is very rough around the edges), I haven't found anything else that approaches personal knowledge management quite this way. Not yet, anyway.
- Your tasks don't exist in isolation; they're connected to your notes, research, and thoughts
- Temporal search means you can find task patterns across time ("What was I working on last quarter?")
- Chat with your notes to understand how your productivity intersects with your knowledge base

### Perfect Integration
- Natural language task logging that becomes part of your knowledge graph
- Seamless connection between your productivity data and your broader knowledge base (should you choose to go all in with mem.ai that is)
- Future-proof your productivity data - it's all searchable and connected
- Cross-pollination between your tasks and your notes (because context matters)

### Long-term Value
- Your productivity data becomes part of your personal knowledge base
- Understand how your tasks relate to your broader goals and notes
- Query your past productivity patterns alongside your notes
- Build a genuine understanding of your work patterns over time

## Why Perplexity?

It's just a start.

### Current Benefits
- Quick to implement and get started
- Built-in internet search capabilities
- Reasonable cost structure
- Good balance of features and simplicity

### Future Plans
Goal is to make Tomatick model-agnostic. Soon, you'll be able to:
- Choose from multiple AI providers
- Bring your own API keys
- Mix and match models for different features
- Self-host your own models (coming soon)

My goal is to make Tomatick model-agnostic, giving you the freedom to choose how you want to use AI with the tool. I don't want to lock you into any specific AI provider.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Break Monitoring System

### Smart Break Detection
- Real-time activity monitoring during breaks
- Intelligent violation detection:
  - 30-second continuous activity threshold
  - Casual interactions ignored
  - Work app detection
- Privacy-focused, local-only monitoring

### Intelligent Notifications
- Context-aware break reminders
- Clean, centered notifications with:
  - Activity summaries
  - Personalized advice
  - Violation tracking
  - Automatic cleanup
- Burnout prevention through pattern recognition

Example notification:
```
╭──────────────────────────────────────────────────────────────╮
│                        🔔 Break Time                         │
│                                                             │
│ Active work in VS Code detected (2.5 minutes)               │
│ Consider stepping away from your workstation to refresh.     │
│                                                             │
│ Break violation #2 today                                    │
╰──────────────────────────────────────────────────────────────╯
```

### Break Violation Rules
- Violations only triggered by:
  1. Active work in work-related apps (30+ seconds)
  2. Sustained activity in any app (30+ seconds)
- Casual interactions ignored:
  - Brief mouse movements
  - Quick app switches
  - Short keyboard usage

### Performance Optimizations
- Smart app name caching (500ms)
- Notification throttling (60s)
- Efficient event buffering
- Memory-safe implementation

### Privacy First
- All monitoring is local
- No keystroke logging
- Break-time only monitoring
- Transparent data handling

For detailed documentation of the monitoring system, see [pkg/monitor/README.md](pkg/monitor/README.md).


