package client

import (
	"github.com/apex/log"
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/internal/parallel"
	"github.com/cidertool/cider/pkg/config"
	"github.com/cidertool/cider/pkg/context"
)

func (c *ascClient) UpdateBetaAppLocalizations(ctx *context.Context, appID string, config config.TestflightLocalizations) error {
	var g = parallel.New(ctx.MaxProcesses)
	locListResp, _, err := c.client.TestFlight.ListBetaAppLocalizationsForApp(ctx, appID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for i := range locListResp.Data {
		loc := locListResp.Data[i]
		locale := *loc.Attributes.Locale
		log.WithField("locale", locale).Debug("found beta app locale")
		locConfig, ok := config[locale]
		if !ok {
			log.WithField("locale", locale).Debug("not in configuration. skipping...")
			continue
		}
		found[locale] = true

		g.Go(func() error {
			_, _, err = c.client.TestFlight.UpdateBetaAppLocalization(ctx, loc.ID, &asc.BetaAppLocalizationUpdateRequestAttributes{
				Description:       &locConfig.Description,
				FeedbackEmail:     &locConfig.FeedbackEmail,
				MarketingURL:      &locConfig.MarketingURL,
				PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
				TVOSPrivacyPolicy: &locConfig.TVOSPrivacyPolicy,
			})
			return err
		})
	}

	for locale := range config {
		locale := locale
		if found[locale] {
			continue
		}
		locConfig := config[locale]

		g.Go(func() error {
			_, _, err = c.client.TestFlight.CreateBetaAppLocalization(ctx.Context, asc.BetaAppLocalizationCreateRequestAttributes{
				Description:       &locConfig.Description,
				FeedbackEmail:     &locConfig.FeedbackEmail,
				Locale:            locale,
				MarketingURL:      &locConfig.MarketingURL,
				PrivacyPolicyURL:  &locConfig.PrivacyPolicyURL,
				TVOSPrivacyPolicy: &locConfig.TVOSPrivacyPolicy,
			}, appID)
			return err
		})
	}

	return g.Wait()
}

func (c *ascClient) UpdateBetaBuildDetails(ctx *context.Context, buildID string, config config.Testflight) error {
	_, _, err := c.client.TestFlight.UpdateBuildBetaDetail(ctx, buildID, &config.EnableAutoNotify)
	return err
}

func (c *ascClient) UpdateBetaBuildLocalizations(ctx *context.Context, buildID string, config config.TestflightLocalizations) error {
	var g = parallel.New(ctx.MaxProcesses)
	locListResp, _, err := c.client.TestFlight.ListBetaBuildLocalizationsForBuild(ctx, buildID, nil)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for i := range locListResp.Data {
		loc := locListResp.Data[i]
		locale := *loc.Attributes.Locale
		log.WithField("locale", locale).Debug("found beta build locale")
		locConfig, ok := config[locale]
		if !ok {
			log.WithField("locale", locale).Debug("not in configuration. skipping...")
			continue
		}
		found[locale] = true

		g.Go(func() error {
			_, _, err := c.client.TestFlight.UpdateBetaBuildLocalization(ctx, loc.ID, &locConfig.WhatsNew)
			return err
		})
	}

	for locale := range config {
		locale := locale
		if found[locale] {
			continue
		}
		locConfig := config[locale]

		g.Go(func() error {
			_, _, err := c.client.TestFlight.CreateBetaBuildLocalization(ctx.Context, locale, &locConfig.WhatsNew, buildID)
			return err
		})
	}

	return g.Wait()
}

func (c *ascClient) UpdateBetaLicenseAgreement(ctx *context.Context, appID string, config config.Testflight) error {
	resp, _, err := c.client.TestFlight.GetBetaLicenseAgreementForApp(ctx, appID, nil)
	if err != nil {
		return err
	}

	_, _, err = c.client.TestFlight.UpdateBetaLicenseAgreement(ctx, resp.Data.ID, &config.LicenseAgreement)
	return err
}

func (c *ascClient) AssignBetaGroups(ctx *context.Context, appID string, buildID string, groups []config.BetaGroup) error {
	var g = parallel.New(ctx.MaxProcesses)

	if len(groups) == 0 {
		log.Debug("no groups provided as input to add")
		return nil
	}

	existingGroups, err := c.listBetaGroups(ctx, appID, groups)
	if err != nil {
		return err
	}

	var groupNamesToConfig = make(map[string]*config.BetaGroup, len(groups))
	for i := range groups {
		groupNamesToConfig[groups[i].Name] = &groups[i]
	}

	found := make(map[string]bool)
	for i := range existingGroups {
		group := existingGroups[i]
		if group.Attributes == nil || group.Attributes.Name == nil {
			continue
		}
		name := *group.Attributes.Name
		found[name] = true
		if configGroup, ok := groupNamesToConfig[name]; ok && configGroup != nil {
			g.Go(func() error {
				if _, _, err := c.client.TestFlight.UpdateBetaGroup(ctx, group.ID, &asc.BetaGroupUpdateRequestAttributes{
					FeedbackEnabled:        &configGroup.FeedbackEnabled,
					Name:                   &configGroup.Name,
					PublicLinkEnabled:      &configGroup.EnablePublicLink,
					PublicLinkLimit:        &configGroup.PublicLinkLimit,
					PublicLinkLimitEnabled: &configGroup.EnablePublicLinkLimit,
				}); err != nil {
					return err
				}
				return c.AssignBetaTesters(ctx, appID, buildID, group.ID, configGroup.Testers)
			})
		}
		g.Go(func() error {
			_, err := c.client.TestFlight.AddBuildsToBetaGroup(ctx, group.ID, []string{buildID})
			return err
		})
	}

	for i := range groups {
		group := groups[i]
		if group.Name == "" {
			log.Warn("beta group name missing")
			continue
		} else if found[group.Name] {
			continue
		}

		g.Go(func() error {
			groupResp, _, err := c.client.TestFlight.CreateBetaGroup(ctx, asc.BetaGroupCreateRequestAttributes{
				Name:                   group.Name,
				FeedbackEnabled:        &group.FeedbackEnabled,
				PublicLinkEnabled:      &group.EnablePublicLink,
				PublicLinkLimit:        &group.PublicLinkLimit,
				PublicLinkLimitEnabled: &group.EnablePublicLinkLimit,
			}, appID, []string{}, []string{buildID})
			if err != nil {
				return err
			}
			return c.AssignBetaTesters(ctx, appID, buildID, groupResp.Data.ID, group.Testers)
		})
	}

	return g.Wait()
}

func (c *ascClient) listBetaGroups(ctx *context.Context, appID string, config []config.BetaGroup) ([]asc.BetaGroup, error) {
	var groupNames = make([]string, len(config))
	for i := range config {
		groupNames[i] = config[i].Name
	}

	groupsResp, _, err := c.client.TestFlight.ListBetaGroups(ctx, &asc.ListBetaGroupsQuery{
		FilterName: groupNames,
		FilterApp:  []string{appID},
	})
	if err != nil {
		return nil, err
	}
	return groupsResp.Data, nil
}

func (c *ascClient) AssignBetaTesters(ctx *context.Context, appID string, buildID string, betaGroupID string, testers []config.BetaTester) error {
	var g = parallel.New(ctx.MaxProcesses)

	if len(testers) == 0 {
		return nil
	}

	existingTesters, err := c.listBetaTesters(ctx, appID, testers)
	if err != nil {
		return err
	}

	found := make(map[string]bool)
	for i := range existingTesters {
		tester := existingTesters[i]
		if tester.Attributes == nil || tester.Attributes.Email == nil {
			continue
		}
		found[string(*tester.Attributes.Email)] = true
		g.Go(func() error {
			_, err := c.client.TestFlight.AssignSingleBetaTesterToBuilds(ctx, tester.ID, []string{buildID})
			return err
		})
	}

	for i := range testers {
		tester := testers[i]
		if tester.Email == "" {
			log.Warn("beta tester email missing")
			continue
		} else if found[tester.Email] {
			continue
		}
		g.Go(func() error {
			var betaGroupIDs []string
			if betaGroupID != "" {
				betaGroupIDs = []string{betaGroupID}
			} else {
				betaGroupIDs = make([]string, 0)
			}
			_, _, err := c.client.TestFlight.CreateBetaTester(ctx, asc.BetaTesterCreateRequestAttributes{
				Email:     asc.Email(tester.Email),
				FirstName: &tester.FirstName,
				LastName:  &tester.LastName,
			}, betaGroupIDs, []string{buildID})
			return err
		})
	}

	return g.Wait()
}

func (c *ascClient) listBetaTesters(ctx *context.Context, appID string, config []config.BetaTester) ([]asc.BetaTester, error) {
	emailFilters := make([]string, 0)
	firstNameFilters := make([]string, 0)
	lastNameFilters := make([]string, 0)
	for _, tester := range config {
		if tester.Email != "" {
			emailFilters = append(emailFilters, tester.Email)
		}
		if tester.FirstName != "" {
			firstNameFilters = append(firstNameFilters, tester.FirstName)
		}
		if tester.LastName != "" {
			lastNameFilters = append(lastNameFilters, tester.LastName)
		}
	}

	testersResp, _, err := c.client.TestFlight.ListBetaTesters(ctx, &asc.ListBetaTestersQuery{
		FilterEmail:     emailFilters,
		FilterFirstName: firstNameFilters,
		FilterLastName:  lastNameFilters,
		FilterApps:      []string{appID},
	})
	if err != nil {
		return nil, err
	}
	return testersResp.Data, nil
}

func (c *ascClient) UpdateBetaReviewDetails(ctx *context.Context, appID string, config config.ReviewDetails) error {
	detailsResp, _, err := c.client.TestFlight.GetBetaAppReviewDetailsForApp(ctx, appID, nil)
	if err != nil {
		return err
	}
	attributes := asc.BetaAppReviewDetailUpdateRequestAttributes{}
	if config.Contact != nil {
		attributes.ContactEmail = &config.Contact.Email
		attributes.ContactFirstName = &config.Contact.FirstName
		attributes.ContactLastName = &config.Contact.LastName
		attributes.ContactPhone = &config.Contact.Phone
	}
	if config.DemoAccount != nil {
		attributes.DemoAccountName = &config.DemoAccount.Name
		attributes.DemoAccountPassword = &config.DemoAccount.Password
		attributes.DemoAccountRequired = &config.DemoAccount.Required
	}
	attributes.Notes = &config.Notes
	if len(config.Attachments) > 0 {
		log.Warn("attachments are not supported for beta review details and will be ignored")
	}
	_, _, err = c.client.TestFlight.UpdateBetaAppReviewDetail(ctx, detailsResp.Data.ID, &attributes)
	return err
}

func (c *ascClient) SubmitBetaApp(ctx *context.Context, buildID string) error {
	_, _, err := c.client.TestFlight.CreateBetaAppReviewSubmission(ctx, buildID)
	return err
}
