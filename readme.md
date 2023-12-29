# Tomatick Memento

Another pomodoro timer, but this has some bells and whistles :tada:

- A cli tool, coz I don't do frontend (thanks, but no thanks JS)
- Interactive, asks you to list tasks - the idea is to make you think before you ink. Once the pom cycle is complete, it asks you to reflect and record wins/distractions. 
- The tasks, your reflection of wins and distractions all get posted to `mem.ai` (which is what Im moving over to from obsidian, and I think I like it so far; so, Im going all in), formatted as markdown.
- That's it. Are you productive yet?

## Features

- Interactive CLI interface with colors and more jazz :tada:
- Customizable Pomodoro and break durations.
- Reflection logging for each Pomodoro cycle.
- Integration with `mem.ai` for persistent tracking.

## Getting Started

### Prerequisites

- Go (version 1.15 or higher)
- An account with `mem.ai` and an API token

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/1x-eng/tomatick.git
   ```

2. `cd tomatick`

3. `go mod tidy`

4. Setup your .env
```
POMODORO_DURATION=25m
SHORT_BREAK_DURATION=5m
LONG_BREAK_DURATION=15m
CYCLES_BEFORE_LONGBREAK=4
MEM_AI_API_TOKEN=your_mem_ai_api_token
```

### Usage

- Run it, `go run main.go` from repo root.
- Follow interactive prompts, and enjoy :wave:


