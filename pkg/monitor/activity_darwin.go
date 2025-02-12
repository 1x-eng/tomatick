package monitor

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/1x-eng/tomatick/config"
)

// IMPORTANT: The following block is NOT a regular comment.
// It contains actual Objective-C code that is compiled and linked by CGO (Go's C interop system).
// This code is essential for interfacing with macOS APIs to monitor user activity.
//
// The block includes:
// 1. CGO directives for compilation and linking
// 2. Required macOS framework imports
// 3. Objective-C implementations for:
//    - Getting frontmost application name
//    - Checking system idle time
//    - Detecting screen lock/screensaver state
//
// Without this block, the activity monitoring functionality would not work
// as we wouldn't have access to the necessary macOS system APIs.

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices
#import <Cocoa/Cocoa.h>
#import <ApplicationServices/ApplicationServices.h>
#import <stdlib.h>
#import <string.h>

// getFrontmostAppName returns the name of the currently active application
// Uses NSWorkspace API to safely access this information
// Returns NULL if any step fails to ensure proper error handling
static char* getFrontmostAppName(void) {
    @autoreleasepool {
        NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
        NSRunningApplication *app = [workspace frontmostApplication];
        if (app == nil) {
            return NULL;
        }
        NSString *name = [app localizedName];
        if (name == nil) {
            return NULL;
        }
        const char *utf8String = [name UTF8String];
        if (utf8String == NULL) {
            return NULL;
        }
        // Create a copy that we can free later in Go code
        char *copy = strdup(utf8String);
        return copy;
    }
}

// getIdleTime returns the number of seconds since the last user input
// Uses CGEventSource API to detect keyboard/mouse activity
static double getIdleTime(void) {
    CFTimeInterval idleTime = CGEventSourceSecondsSinceLastEventType(
        kCGEventSourceStateHIDSystemState,
        kCGAnyInputEventType
    );
    return idleTime;
}

// checkScreenSaver checks if the screensaver is active
// Uses CGSession API to safely check screensaver state
static int checkScreenSaver(void) {
    @autoreleasepool {
        CFDictionaryRef sessionDict = CGSessionCopyCurrentDictionary();
        if (sessionDict == NULL) {
            return 0;
        }

        CFBooleanRef screenSaverValue = CFDictionaryGetValue(sessionDict,
            CFSTR("CGSSessionScreenIsLocked"));

        int isLocked = screenSaverValue != NULL &&
            CFBooleanGetValue(screenSaverValue);

        CFRelease(sessionDict);
        return isLocked;
    }
}

// checkScreenLock checks if the screen is locked
// Uses CGSession API to safely check screen lock state
static int checkScreenLock(void) {
    @autoreleasepool {
        CFDictionaryRef sessionDict = CGSessionCopyCurrentDictionary();
        if (sessionDict == NULL) {
            return 0;
        }

        CFBooleanRef screenLockValue = CFDictionaryGetValue(sessionDict,
            CFSTR("CGSSessionScreenIsLocked"));

        int isLocked = screenLockValue != NULL &&
            CFBooleanGetValue(screenLockValue);

        CFRelease(sessionDict);
        return isLocked;
    }
}
*/
import "C"

var (
	// idleThreshold is the time in seconds after which we consider the user idle
	idleThreshold = 10.0 // Increased to 10 seconds to avoid false positives from casual mouse movement

	// activityThreshold is how long continuous activity needs to be present to count as a violation
	activityThreshold = 30.0 // 30 seconds of continuous activity before counting as a violation

	// lastAppName caches the last fetched app name to reduce system calls
	lastAppName     string
	lastAppNameTime time.Time
	appNameCacheTTL = 500 * time.Millisecond

	// Track continuous activity
	lastActivityStart time.Time

	// workApps holds the configured work-related applications
	workApps map[string]bool
)

// getForegroundApp returns the name of the frontmost application
func getForegroundApp() (string, error) {
	// Check cache first
	if time.Since(lastAppNameTime) < appNameCacheTTL {
		return lastAppName, nil
	}

	cAppName := C.getFrontmostAppName()
	if cAppName == nil {
		return "", fmt.Errorf("failed to get frontmost app name")
	}
	// Ensure we free the memory allocated by strdup in the C code
	defer C.free(unsafe.Pointer(cAppName))

	appName := C.GoString(cAppName)

	// Update cache
	lastAppName = appName
	lastAppNameTime = time.Now()

	return appName, nil
}

// isWorkApp checks if the given application is considered a work-related app
func isWorkApp(appName string) bool {
	// First check exact match
	if workApps[appName] {
		return true
	}

	// Then check if any work app name is contained in the given app name
	appNameLower := strings.ToLower(appName)
	for app := range workApps {
		if strings.Contains(appNameLower, strings.ToLower(app)) {
			return true
		}
	}

	return false
}

// hasRecentActivity checks if there has been any keyboard or mouse activity recently
func hasRecentActivity() bool {
	// Check if screen is locked or screensaver is active
	if C.checkScreenLock() == 1 || C.checkScreenSaver() == 1 {
		return false
	}

	idleTime := float64(C.getIdleTime())

	// If user has been idle longer than threshold, reset activity tracking
	if idleTime >= idleThreshold {
		lastActivityStart = time.Time{} // Reset activity start time
		return false
	}

	// If this is the start of new activity, record the time
	if lastActivityStart.IsZero() {
		lastActivityStart = time.Now().Add(-time.Duration(idleTime * float64(time.Second)))
	}

	// Only count as violation if activity has been continuous for longer than activityThreshold
	return time.Since(lastActivityStart).Seconds() >= activityThreshold
}

// InitializeMonitoring sets up any necessary system-level monitoring
func InitializeMonitoring(cfg *config.Config) error {
	// Only proceed if break monitoring is enabled
	if !cfg.Features.BreakMonitoring {
		return fmt.Errorf("break monitoring is not supported on this operating system")
	}

	// Initialize work apps map
	workApps = make(map[string]bool)
	for _, app := range cfg.WorkApps {
		workApps[app] = true
	}

	// Verify we can access the necessary macOS APIs
	if _, err := getForegroundApp(); err != nil {
		return fmt.Errorf("failed to initialize app monitoring: %w", err)
	}

	if C.getIdleTime() < 0 {
		return fmt.Errorf("failed to initialize idle time monitoring")
	}

	return nil
}
