// Package client provides a full-featured App Store Connect API client
package client

import (
	"crypto/md5" // #nosec
	"fmt"
	"io"

	"github.com/aaronsky/applereleaser/pkg/config"
	"github.com/aaronsky/applereleaser/pkg/context"
	"github.com/aaronsky/asc-go/asc"
)

// Client is an abstraction of an App Store Connect API client's functionality.
type Client interface {
	// GetAppForBundleID returns the App resource matching the given bundle ID
	GetAppForBundleID(ctx *context.Context, id string) (*asc.App, error)
	GetAppInfo(ctx *context.Context, app *asc.App) (*asc.AppInfo, error)
	// GetRelevantBuild returns the latest Build resource for the given app. Returns an error if
	// the latest build is still processing.
	GetRelevantBuild(ctx *context.Context, app *asc.App) (*asc.Build, error)
	// ReleaseForAppIsInitial returns true if the App resource has never released before,
	// i.e. has one or less associated App Store Version relationships.
	ReleaseForAppIsInitial(ctx *context.Context, app *asc.App) (bool, error)

	// Testflight

	// UpdateBetaAppLocalizations updates an App's beta app localizations, and creates any new ones that do not exist.
	// It will not delete or update any locales that are associated with the app but are not configured in applereleaser.
	UpdateBetaAppLocalizations(ctx *context.Context, app *asc.App, config config.TestflightLocalizations) error
	// UpdateBetaBuildDetails updates an App's beta build details, or creates new ones if they do not yet exist.
	UpdateBetaBuildDetails(ctx *context.Context, build *asc.Build, config config.TestflightForApp) error
	// UpdateBetaBuildLocalizations updates an App's beta build localizations, and creates any new ones that do not exist.
	// It will not delete or update any locales that are associated with the app but are not configured in applereleaser.
	UpdateBetaBuildLocalizations(ctx *context.Context, build *asc.Build, config config.TestflightLocalizations) error
	// UpdateBetaLicenseAgreement updates an App's beta license agreement, or creates a new one if one does not yet exist.
	UpdateBetaLicenseAgreement(ctx *context.Context, app *asc.App, config config.TestflightForApp) error
	AssignBetaGroups(ctx *context.Context, build *asc.Build, groups []string) error
	AssignBetaTesters(ctx *context.Context, build *asc.Build, testers []config.BetaTester) error
	// UpdateBetaReviewDetails updates an App's beta review details, or creates new ones if they do not yet exist.
	UpdateBetaReviewDetails(ctx *context.Context, app *asc.App, config config.ReviewDetails) error
	// SubmitBetaApp submits the given beta build for review
	SubmitBetaApp(ctx *context.Context, build *asc.Build) error

	// App Store

	UpdateApp(ctx *context.Context, app *asc.App, appInfo *asc.AppInfo, config config.App) error
	UpdateAppLocalizations(ctx *context.Context, app *asc.App, appInfo *asc.AppInfo, config config.AppLocalizations) error
	CreateVersionIfNeeded(ctx *context.Context, app *asc.App, build *asc.Build, config config.Version) (*asc.AppStoreVersion, error)
	UpdateVersionLocalizations(ctx *context.Context, version *asc.AppStoreVersion, config config.VersionLocalizations) error
	UpdateIDFADeclaration(ctx *context.Context, version *asc.AppStoreVersion, config config.IDFADeclaration) error
	UploadRoutingCoverage(ctx *context.Context, version *asc.AppStoreVersion, config config.File) error
	// UpdateReviewDetails updates an App's review details, or creates new ones if they do not yet exist.
	UpdateReviewDetails(ctx *context.Context, version *asc.AppStoreVersion, config config.ReviewDetails) error
	// SubmitApp submits the given app store version for review
	SubmitApp(ctx *context.Context, version *asc.AppStoreVersion) error
}

// New returns a new Client.
func New(ctx *context.Context) Client {
	client := asc.NewClient(ctx.Credentials.Client())
	return &ascClient{client: client}
}

type ascClient struct {
	client *asc.Client
}

func (c *ascClient) GetAppForBundleID(ctx *context.Context, id string) (*asc.App, error) {
	resp, _, err := c.client.Apps.ListApps(ctx, &asc.ListAppsQuery{
		FilterBundleID: []string{id},
	})
	if err != nil {
		return nil, fmt.Errorf("app not found matching %s: %w", id, err)
	} else if len(resp.Data) == 0 {
		return nil, fmt.Errorf("app not found matching %s", id)
	}
	return &resp.Data[0], nil
}

func (c *ascClient) GetAppInfo(ctx *context.Context, app *asc.App) (*asc.AppInfo, error) {
	resp, _, err := c.client.Apps.ListAppInfosForApp(ctx, app.ID, nil)
	if err != nil {
		return nil, err
	}
	for _, info := range resp.Data {
		if info.Attributes == nil {
			continue
		} else if info.Attributes.AppStoreState == nil {
			continue
		}
		state := *info.Attributes.AppStoreState
		if state == asc.AppStoreVersionStatePrepareForSubmission {
			return &info, nil
		}
	}

	return nil, fmt.Errorf("app info not found matching %s", app.ID)
}

func (c *ascClient) GetRelevantBuild(ctx *context.Context, app *asc.App) (*asc.Build, error) {
	if ctx.Version == "" {
		return nil, fmt.Errorf("no version provided to lookup build with")
	}
	resp, _, err := c.client.Builds.ListBuilds(ctx, &asc.ListBuildsQuery{
		FilterApp:                      []string{app.ID},
		FilterPreReleaseVersionVersion: []string{ctx.Version},
	})
	if err != nil {
		return nil, fmt.Errorf("build not found matching app %s and version %s: %w", *app.Attributes.BundleID, ctx.Version, err)
	} else if len(resp.Data) == 0 {
		return nil, fmt.Errorf("build not found matching app %s and version %s", *app.Attributes.BundleID, ctx.Version)
	}
	build := resp.Data[0]
	if build.Attributes == nil {
		return nil, fmt.Errorf("build %s has no attributes", build.ID)
	}
	if build.Attributes.ProcessingState == nil {
		return nil, fmt.Errorf("build %s has no processing state", build.ID)
	}
	if *build.Attributes.ProcessingState != "VALID" {
		return nil, fmt.Errorf("latest build %s has a processing state of %s. it would be dangerous to proceed", build.ID, *build.Attributes.ProcessingState)
	}
	return &build, nil
}

func (c *ascClient) ReleaseForAppIsInitial(ctx *context.Context, app *asc.App) (bool, error) {
	resp, _, err := c.client.Apps.ListAppStoreVersionsForApp(ctx, app.ID, nil)
	if err != nil {
		return false, err
	}
	return len(resp.Data) <= 1, nil
}

func md5Checksum(f io.Reader) (string, error) {
	/* #nosec */
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
