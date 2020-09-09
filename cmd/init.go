package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/aaronsky/applereleaser/internal/closer"
	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const configDocString = `# This is a template applereleaser.yaml file with some sane defaults, initially-generated by applereleaser init.
# Check this file into your repository so you can version changes to your apps' configurations in App Store Connect.
# For additional configuration options, see: https://aaronsky.github.io/applereleaser/configuration
#

`

// ErrNotImplemented happens when the author has not implemented a feature yet, but it's still accessible by the user.
var ErrNotImplemented = errors.New("this feature has yet to be implemented")

type initCmd struct {
	cmd  *cobra.Command
	opts initOpts
}

type initOpts struct {
	config     string
	export     bool
	skipPrompt bool
}

func newInitCmd() *initCmd {
	var root = &initCmd{}
	var cmd = &cobra.Command{
		Use:           "init",
		Short:         "Generates an .applereleaser.yml file",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initProject(root.opts)
		},
	}

	cmd.Flags().StringVarP(&root.opts.config, "config", "f", ".applereleaser.yml", "Path of configuration file to create")
	cmd.Flags().BoolVar(&root.opts.export, "export", false, "Populate the project file automatically using your App Store Connect team")
	cmd.Flags().BoolVarP(&root.opts.skipPrompt, "skip-prompt", "y", false, "Skips onboarding prompts. This can result in an overwritten configuration file")

	root.cmd = cmd
	return root
}

func initProject(opts initOpts) error {
	file, err := createFileIfNeeded(opts.config, opts.skipPrompt)
	if err != nil {
		return err
	}
	defer closer.Close(file)

	log.Info(color.New(color.Bold).Sprintf("Populating project file at %s", opts.config))

	project, err := newProject(opts.export, opts.skipPrompt)
	if err != nil {
		return err
	}

	if err := writeProject(project, file); err != nil {
		return err
	}

	log.
		WithField("file", file.Name()).
		Info("config created")
	log.Info("Please edit accordingly to fit your needs.")
	log.Info("For additional configuration options, see: https://aaronsky.github.io/applereleaser/configuration")

	return nil
}

func createFileIfNeeded(path string, skipPrompt bool) (*os.File, error) {
	f, err := os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, 0600)
	if err == nil {
		return f, nil
	}

	if !os.IsExist(err) {
		return nil, err
	}

	if skipPrompt {
		log.Warn("file exists, overwriting")
	} else {
		prompt := promptui.Prompt{
			Label:     "Overwrite file?",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filepath.Clean(path), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

func newProject(export, skipPrompt bool) (*config.Project, error) {
	var project *config.Project
	var err error

	switch {
	case export:
		project, err = newProjectFromAPI()
	case !skipPrompt:
		project, err = newProjectFromPrompts()
	default:
		project = newProjectFromDefaults()
	}

	return project, err
}

func newProjectFromAPI() (*config.Project, error) {
	return nil, ErrNotImplemented
}

func newProjectFromPrompts() (*config.Project, error) {
	values := projectInitValues{
		Apps: map[string]projectInitAppValues{},
	}
	projectNamePrompt := promptui.Prompt{Label: "Project Name"}
	projectName, err := projectNamePrompt.Run()
	if err != nil {
		return nil, err
	}
	values.Name = projectName

	var continueAppsSetup = true
	for continueAppsSetup {
		name, app, err := promptAppValues()
		if err != nil {
			return nil, err
		}
		values.Apps[name] = *app

		continuePrompt := promptui.Prompt{
			Label:     "Add more apps?",
			IsConfirm: true,
		}
		_, err = continuePrompt.Run()
		continueAppsSetup = err == nil
	}

	proj := newProjectFromValues(values)
	return &proj, nil
}

func promptAppValues() (name string, app *projectInitAppValues, err error) {
	var prompt promptui.Prompt
	var selec promptui.Select

	app = &projectInitAppValues{
		AvailableInNewTerritories:    false,
		Territories:                  []string{"USA"},
		EnableAutoNotify:             true,
		ReviewDetailsAccountRequired: false,
		PhasedReleaseEnabled:         true,
	}

	log.Info("Let's set up an app in your project!")

	prompt = promptui.Prompt{Label: "App Name"}
	name, err = prompt.Run()
	if err != nil {
		return name, app, err
	}
	app.NameInPrimaryLocale = name

	prompt = promptui.Prompt{Label: "Bundle ID"}
	bundleID, err := prompt.Run()
	if err != nil {
		return name, app, err
	}
	app.BundleID = bundleID

	selec = promptui.Select{
		Label: "Platform",
		Items: []config.Platform{
			config.PlatformiOS,
			config.PlatformTvOS,
			config.PlatformMacOS,
		},
	}
	_, platform, err := selec.Run()
	if err != nil {
		return name, app, err
	}
	app.Platform = platform

	prompt = promptui.Prompt{
		Label:   "Primary Locale",
		Default: "en-US",
	}
	primaryLocale, err := prompt.Run()
	if err != nil {
		return name, app, err
	}
	app.PrimaryLocale = primaryLocale

	prompt = promptui.Prompt{
		Label: "Price Tier",
	}
	tier, err := prompt.Run()
	if err != nil {
		return name, app, err
	}
	app.PricingTier = tier

	return name, app, nil
}

func newProjectFromDefaults() *config.Project {
	proj := newProjectFromValues(projectInitValues{
		Name: "My Project",
		Apps: map[string]projectInitAppValues{
			"My App": {
				BundleID:                     "com.app.bundleid",
				PrimaryLocale:                "en-US",
				AvailableInNewTerritories:    false,
				PricingTier:                  "0",
				Territories:                  []string{"USA"},
				NameInPrimaryLocale:          "My App",
				EnableAutoNotify:             true,
				ReviewDetailsAccountRequired: false,
				Platform:                     "iOS",
				PhasedReleaseEnabled:         true,
			},
		},
	})
	return &proj
}

type projectInitValues struct {
	Name string
	Apps map[string]projectInitAppValues
}

type projectInitAppValues struct {
	BundleID                     string
	PrimaryLocale                string
	NameInPrimaryLocale          string
	Platform                     string
	PricingTier                  string
	PhasedReleaseEnabled         bool
	EnableAutoNotify             bool
	ReviewDetailsAccountRequired bool
	AvailableInNewTerritories    bool
	Territories                  []string
}

func newProjectFromValues(values projectInitValues) config.Project {
	var project = config.Project{
		Name: values.Name,
		Apps: map[string]config.App{},
	}

	for name, app := range values.Apps {
		availableInNewTerritories := app.AvailableInNewTerritories
		project.Apps[name] = config.App{
			BundleID:      app.BundleID,
			PrimaryLocale: app.PrimaryLocale,
			Availability: &config.Availability{
				AvailableInNewTerritories: &availableInNewTerritories,
				Pricing: []config.PriceSchedule{
					{
						Tier: app.PricingTier,
					},
				},
				Territories: app.Territories,
			},
			Localizations: config.AppLocalizations{
				app.PrimaryLocale: config.AppLocalization{
					Name: app.NameInPrimaryLocale,
				},
			},
			Testflight: config.TestflightForApp{
				EnableAutoNotify: app.EnableAutoNotify,
				Localizations: config.TestflightLocalizations{
					app.PrimaryLocale: config.TestflightLocalization{},
				},
				ReviewDetails: &config.ReviewDetails{
					Contact: &config.ContactPerson{},
					DemoAccount: &config.DemoAccount{
						Required: app.ReviewDetailsAccountRequired,
					},
				},
			},
			Versions: config.Version{
				Platform: config.Platform(app.Platform),
				Localizations: config.VersionLocalizations{
					app.PrimaryLocale: config.VersionLocalization{},
				},
				PhasedReleaseEnabled: app.PhasedReleaseEnabled,
				ReleaseType:          config.ReleaseTypeAfterApproval,
				ReviewDetails: &config.ReviewDetails{
					Contact: &config.ContactPerson{},
					DemoAccount: &config.DemoAccount{
						Required: app.ReviewDetailsAccountRequired,
					},
				},
			},
		}
	}

	return project
}

func writeProject(project *config.Project, f io.StringWriter) error {
	contents, err := project.String()
	if err != nil {
		return err
	}

	if _, err := f.WriteString(configDocString + contents); err != nil {
		return err
	}

	return nil
}

// 	return nil
// }

// func promptApp() (string, config.App, error) {

// }
