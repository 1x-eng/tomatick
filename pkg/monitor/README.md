# Tomatick Activity Monitor

This package provides macOS-native activity monitoring for Tomatick, enabling intelligent break tracking and adaptive notifications.

## Architecture

The monitor package consists of several key components:

```
monitor/
├── activity.go           # Core activity monitoring logic
├── activity_darwin.go    # macOS-specific implementations
├── notifications.go      # Intelligent break notifications
├── integration.go        # Integration with Tomatick core
└── README.md            # This file
```

## Core Features

### 1. Activity Monitoring
- Real-time tracking of:
  - Active applications
  - Keyboard/mouse activity
  - Screen lock/screensaver state
- Smart activity detection:
  - 10-second idle threshold for casual movements
  - 30-second continuous activity threshold for violations
  - Intelligent work app detection
- Memory-efficient event buffering

### 2. Intelligent Notifications
- Context-aware break reminders
- Adaptive messaging based on:
  - Current activity context
  - Break violation patterns
  - Violation frequency
- Burnout risk detection after 3+ violations
- Clean, centered notifications with:
  - Consistent width and padding
  - Clear activity summaries
  - Violation count tracking
  - Automatic cleanup

### 3. Break Violation Rules
- Violations only triggered by:
  1. Active work in work-related apps (30+ seconds)
  2. Sustained activity in any app (30+ seconds)
- Casual interactions ignored:
  - Brief mouse movements
  - Quick app switches
  - Short keyboard usage
- Detailed violation context:
  - App name and activity type
  - Duration of activity
  - Time of violation

## Usage

Basic integration:

```go
monitor, err := monitor.NewTomatickMonitor(cfg, llmClient)
if err != nil {
    log.Fatal(err)
}

// Start break monitoring
monitor.OnBreakStart()

// Monitor will automatically:
// - Track activity
// - Send intelligent notifications
// - Adapt to user patterns
// - Track violation counts

// End monitoring and get summary
summary := monitor.OnBreakEnd()
```

## Notification System

### Design Philosophy
- Data-driven but concise
- Context-aware and personalized
- Adaptive to user patterns
- Focus on sustainable productivity

### Notification Format
```
╭─────────────────────────────────────╮
│            🔔 Break Time            │
│                                     │
│ [Activity Summary]                  │
│ [Personalized Advice]              │
│                                     │
│ Break violation #N today            │
╰─────────────────────────────────────╯
```

### Violation Tracking
- Counts violations per session
- Adapts message severity based on count
- Provides violation metadata
- Generates end-of-break summaries

## Work App Detection

The monitor maintains a list of work-related applications:

- VS Code (`Code`)
- Cursor Editor (`Cursor`)
- iTerm2 (`iTerm2`)
- Insomnia (`Insomnia`)
- pgAdmin (`pgAdmin`, `pgAdmin 4`)
- Chrome (`Chrome`, `Google Chrome`)
- Terminal (`Terminal`)

Customize by modifying the `workApps` map in `activity_darwin.go`.

## Performance Optimizations

1. **App Name Caching**:
   - 500ms cache duration
   - Reduces system API calls
   - Configurable via `appNameCacheTTL`

2. **Notification Throttling**:
   - 60-second minimum interval
   - Prevents notification spam
   - Configurable via `notifyThreshold`

3. **Event Buffering**:
   - Efficient event channel
   - Automatic cleanup
   - Memory-safe implementation

## Privacy & Security

- All monitoring is local
- No external data transmission
- Active only during breaks
- Respects system privacy settings
- No keystroke logging
- Only detects activity presence

## Configuration

Key configuration options:

```go
type Config struct {
    // Activity thresholds
    IdleThreshold time.Duration
    
    // Notification settings
    NotifyThreshold time.Duration
    
    // User preferences
    UserName string
}
```

## Error Handling

The system includes graceful handling of:
- System API failures
- LLM service disruptions
- Invalid configurations
- Resource constraints

## Future Enhancements

Planned features:
1. Machine learning for pattern detection
2. Custom work app definitions
3. Break success analytics
4. Integration with health platforms 