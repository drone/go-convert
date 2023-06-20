package browser_tools

import (
	"fmt"
	"strings"

	circle "github.com/drone/go-convert/convert/circle/yaml"
	harness "github.com/drone/spec/dist/go"
)

func Convert(command string, step *circle.Custom) *harness.Step {
	switch command {
	case "":
		return nil // not supported
	case "install-browser-tools":
		return convertInstallBrowserTools(step)
	case "install-chrome":
		return convertInstallChrome(step)
	case "install-chromedriver":
		return convertInstallChromeDriver(step)
	case "install-firefox":
		return convertInstallFirefox(step)
	case "install-geckodriver":
		return convertInstallGeckoDriver(step)
	default:
		return nil // not supported
	}
}

func convertInstallBrowserTools(step *circle.Custom) *harness.Step {
	// default directories
	firefoxInstallDir := "/usr/local/bin"
	geckodriverInstallDir := "/usr/local/bin"
	chromedriverInstallDir := "/usr/local/bin"

	// default versions
	firefoxVersion := "latest"
	geckodriverVersion := "latest"
	chromeVersion := "latest"

	// default installation flags
	installFirefox := true
	installGeckodriver := true
	installChrome := true
	installChromedriver := true

	// check parameters and override defaults if provided
	if s, ok := step.Params["firefox-install-dir"].(string); ok && s != "" {
		firefoxInstallDir = s
	}
	if s, ok := step.Params["geckodriver-install-dir"].(string); ok && s != "" {
		geckodriverInstallDir = s
	}
	if s, ok := step.Params["chromedriver-install-dir"].(string); ok && s != "" {
		chromedriverInstallDir = s
	}
	if s, ok := step.Params["firefox-version"].(string); ok && s != "" {
		firefoxVersion = s
	}
	if s, ok := step.Params["geckodriver-version"].(string); ok && s != "" {
		geckodriverVersion = s
	}
	if s, ok := step.Params["chrome-version"].(string); ok && s != "" {
		chromeVersion = s
	}
	if b, ok := step.Params["install-firefox"].(bool); ok {
		installFirefox = b
	}
	if b, ok := step.Params["install-geckodriver"].(bool); ok {
		installGeckodriver = b
	}
	if b, ok := step.Params["install-chrome"].(bool); ok {
		installChrome = b
	}
	if b, ok := step.Params["install-chromedriver"].(bool); ok {
		installChromedriver = b
	}

	var runCommands []string
	if installFirefox {
		runCommands = append(runCommands, "curl https://raw.githubusercontent.com/CircleCI-Public/browser-tools-orb/main/src/scripts/install-firefox.sh | bash")
	}
	if installGeckodriver {
		runCommands = append(runCommands, "curl https://raw.githubusercontent.com/CircleCI-Public/browser-tools-orb/main/src/scripts/install-geckodriver.sh | bash")
	}
	if installChrome {
		runCommands = append(runCommands, "curl https://raw.githubusercontent.com/CircleCI-Public/browser-tools-orb/main/src/scripts/install-chrome.sh | bash")
	}
	if installChromedriver {
		runCommands = append(runCommands, "curl https://raw.githubusercontent.com/CircleCI-Public/browser-tools-orb/main/src/scripts/install-chromedriver.sh | bash")
	}
	runCommand := strings.Join(runCommands, " && ")

	return &harness.Step{
		Name: "install_browser_tools",
		Type: "script",
		Spec: &harness.StepExec{
			Run: runCommand,
			Envs: map[string]string{
				"ORB_PARAM_FIREFOX_INSTALL_DIR": firefoxInstallDir,
				"ORB_PARAM_FIREFOX_VERSION":     firefoxVersion,
				"ORB_PARAM_CHROME_VERSION":      chromeVersion,
				"ORB_PARAM_DRIVER_INSTALL_DIR":  chromedriverInstallDir,
				"ORB_PARAM_GECKO_INSTALL_DIR":   geckodriverInstallDir,
				"ORB_PARAM_GECKO_VERSION":       geckodriverVersion,
				"HOME":                          "/root", //required for firefox to install
			},
		},
	}
}

func convertInstallChrome(step *circle.Custom) *harness.Step {
	channel := "stable" // default value
	if c, ok := step.Params["channel"].(string); ok && c != "" {
		channel = c
	}

	chromeVersion := "latest" // default value
	if cv, ok := step.Params["chrome-version"].(string); ok && cv != "" {
		chromeVersion = cv
	}

	replaceExisting := "0" // default value
	if re, ok := step.Params["replace-existing"].(bool); ok && re {
		replaceExisting = "1"
	}

	runCommand := "curl https://raw.githubusercontent.com/CircleCI-Public/browser-tools-orb/main/src/scripts/install-chrome.sh | bash"

	return &harness.Step{
		Name: "install_chrome",
		Type: "script",
		Spec: &harness.StepExec{
			Run: runCommand,
			Envs: map[string]string{
				"ORB_PARAM_CHANNEL":          channel,
				"ORB_PARAM_CHROME_VERSION":   chromeVersion,
				"ORB_PARAM_REPLACE_EXISTING": replaceExisting,
			},
		},
	}
}

func convertInstallChromeDriver(step *circle.Custom) *harness.Step {
	installDir := "/usr/local/bin" // default value
	if id, ok := step.Params["install-dir"].(string); ok && id != "" {
		installDir = id
	}

	runCommand := "curl https://raw.githubusercontent.com/CircleCI-Public/browser-tools-orb/main/src/scripts/install-chromedriver.sh | bash"

	return &harness.Step{
		Name: "install_chromedriver",
		Type: "script",
		Spec: &harness.StepExec{
			Run: runCommand,
			Envs: map[string]string{
				"ORB_PARAM_DRIVER_INSTALL_DIR": installDir,
			},
		},
	}
}

func convertInstallFirefox(step *circle.Custom) *harness.Step {
	installDir, _ := step.Params["install-dir"].(string)
	version, _ := step.Params["version"].(string)

	// Set the default values if they are not provided
	if installDir == "" {
		installDir = "/usr/local/bin"
	}
	if version == "" {
		version = "latest"
	}

	runCommand := fmt.Sprintf(`curl -s 'https://raw.githubusercontent.com/CircleCI-Public/browser-tools-orb/main/src/scripts/install-firefox.sh' | bash`)

	return &harness.Step{
		Name: "firefox_setup",
		Type: "script",
		Spec: &harness.StepExec{
			Run: runCommand,
			Envs: map[string]string{
				"ORB_PARAM_FIREFOX_INSTALL_DIR": installDir,
				"ORB_PARAM_FIREFOX_VERSION":     version,
				"HOME":                          "/root", //required for firefox to install
			},
		},
	}
}

func convertInstallGeckoDriver(step *circle.Custom) *harness.Step {
	installDir := "/usr/local/bin" // default value
	if id, ok := step.Params["install-dir"].(string); ok && id != "" {
		installDir = id
	}

	version := "latest" // default value
	if v, ok := step.Params["version"].(string); ok && v != "" {
		version = v
	}

	runCommand := "curl https://raw.githubusercontent.com/CircleCI-Public/browser-tools-orb/main/src/scripts/install-geckodriver.sh | bash"

	return &harness.Step{
		Name: "install_geckodriver",
		Type: "script",
		Spec: &harness.StepExec{
			Run: runCommand,
			Envs: map[string]string{
				"ORB_PARAM_GECKO_INSTALL_DIR": installDir,
				"ORB_PARAM_GECKO_VERSION":     version,
			},
		},
	}
}
