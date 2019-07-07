package browser

import (
	"bytes"
	"github.com/dkoston/cdp-examples/src/logger"
	"os"
	"os/exec"
	"runtime"
)

// LaunchChrome will launch chrome on a mac with remote debugging open on port 9222
func LaunchChrome() {
	readBuffer := bytes.Buffer{}

	cmd := exec.Command(getChromePath(), "--remote-debugging-port="+devToolsPort)
	cmd.Stdout = &readBuffer
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		logger.Get().Fatalf("[ERROR] Failed to launch google chrome with remote debugging: %v\n. %v", err, os.Stderr)
	}
}

// getOperatingSystem returns the name of the OS at runtime
func getOperatingSystem() string {
	return runtime.GOOS
}

// getChromePath gets the likely Chrome executable path based on the operating system
func getChromePath() string {
	operatingSystem := getOperatingSystem()
	switch operatingSystem {
	case "darwin":
		return getPathMac()
	case "windows":
		return getPathWindows()
	default:
		logger.Get().Fatalf("operating system not currently supported: %s", operatingSystem)
	}
	return ""
}

// getPathMac checks the possible paths for Chrome (default installs only) on Mac
func getPathMac() string {
	expectedPath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		logger.Get().Fatalf("Unable to find Chrome executable")
	}
	return expectedPath
}

// getPathWindows checks the possible paths for Chrome (default installs only) on windows
func getPathWindows() string {
	possiblePaths := []string{
		"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe", // Windows 10
		"C:\\Program Files (x86)\\Google\\Application\\chrome.exe", // Windows 7
	}

	for i := 0; i < len(possiblePaths); i++ {
		_, err := os.Stat(possiblePaths[i])
		if os.IsNotExist(err) {
			continue
		}
		return possiblePaths[i]
	}

	logger.Get().Fatalf("Unable to find Chrome executable")
	return ""
}