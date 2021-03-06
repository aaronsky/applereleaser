/**
Copyright (C) 2020 Aaron Sky.

This file is part of Cider, a tool for automating submission
of apps to Apple's App Stores.

Cider is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Cider is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Cider.  If not, see <http://www.gnu.org/licenses/>.
*/

package client

import (
	"net/http"
	"testing"

	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/pkg/config"
	"github.com/stretchr/testify/assert"
)

const (
	testID = "TEST"
)

// Test UpdateBetaAppLocalizations

func TestUpdateBetaAppLocalizations_Happy(t *testing.T) {
	t.Parallel()

	localizations := config.TestflightLocalizations{
		"en-US": {
			Description:       "TEST",
			FeedbackEmail:     "TEST",
			MarketingURL:      "TEST",
			PrivacyPolicyURL:  "TEST",
			TVOSPrivacyPolicy: "TEST",
			WhatsNew:          "Going away",
		},
		"ja": {
			Description:       "TEST",
			FeedbackEmail:     "TEST",
			MarketingURL:      "TEST",
			PrivacyPolicyURL:  "TEST",
			TVOSPrivacyPolicy: "TEST",
			WhatsNew:          "Going away",
		},
	}

	ctx, client := newTestContext(
		response{
			Response: asc.BetaAppLocalizationsResponse{
				Data: []asc.BetaAppLocalization{
					{
						Attributes: &asc.BetaAppLocalizationAttributes{
							Locale: asc.String("en-US"),
						},
						ID: "TEST",
					},
					{
						Attributes: &asc.BetaAppLocalizationAttributes{
							Locale: asc.String("en-GB"),
						},
						ID: "TEST",
					},
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaAppLocalizations(ctx.Context, testID, localizations)
	assert.NoError(t, err)
}

func TestUpdateBetaAppLocalizations_ErrList(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaAppLocalizations(ctx.Context, testID, config.TestflightLocalizations{})
	assert.Error(t, err)
}

func TestUpdateBetaAppLocalizations_ErrUpdate(t *testing.T) {
	t.Parallel()

	localizations := config.TestflightLocalizations{
		"en-US": {
			Description: "TEST",
			WhatsNew:    "Going away",
		},
	}

	ctx, client := newTestContext(
		response{
			Response: asc.BetaAppLocalizationsResponse{
				Data: []asc.BetaAppLocalization{
					{
						Attributes: &asc.BetaAppLocalizationAttributes{
							Locale: asc.String("en-US"),
						},
						ID: "TEST",
					},
				},
			},
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaAppLocalizations(ctx.Context, testID, localizations)
	assert.Error(t, err)
}

func TestUpdateBetaAppLocalizations_ErrCreate(t *testing.T) {
	t.Parallel()

	localizations := config.TestflightLocalizations{
		"en-US": {
			Description: "TEST",
			WhatsNew:    "Going away",
		},
		"ja": {
			Description: "TEST",
			WhatsNew:    "Going away",
		},
	}

	ctx, client := newTestContext(
		response{
			Response: asc.BetaAppLocalizationsResponse{
				Data: []asc.BetaAppLocalization{
					{
						Attributes: &asc.BetaAppLocalizationAttributes{
							Locale: asc.String("en-US"),
						},
						ID: "TEST",
					},
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaAppLocalizations(ctx.Context, testID, localizations)
	assert.Error(t, err)
}

// Test UpdateBetaBuildDetails

func TestUpdateBetaBuildDetails_Happy(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaBuildDetails(ctx.Context, testID, config.Testflight{})
	assert.NoError(t, err)
}

func TestUpdateBetaBuildDetails_Err(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaBuildDetails(ctx.Context, testID, config.Testflight{
		EnableAutoNotify: true,
	})
	assert.Error(t, err)
}

// Test UpdateBetaBuildLocalizations

func TestUpdateBetaBuildLocalizations_Happy(t *testing.T) {
	t.Parallel()

	localizations := config.TestflightLocalizations{
		"en-US": {
			Description: "TEST",
			WhatsNew:    "Going away",
		},
		"ja": {
			Description: "TEST",
			WhatsNew:    "Going away",
		},
	}

	ctx, client := newTestContext(
		response{
			Response: asc.BetaAppLocalizationsResponse{
				Data: []asc.BetaAppLocalization{
					{
						Attributes: &asc.BetaAppLocalizationAttributes{
							Locale: asc.String("en-US"),
						},
						ID: "TEST",
					},
					{
						Attributes: &asc.BetaAppLocalizationAttributes{
							Locale: asc.String("en-GB"),
						},
						ID: "TEST",
					},
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaBuildLocalizations(ctx.Context, testID, localizations)
	assert.NoError(t, err)
}

func TestUpdateBetaBuildLocalizations_ErrList(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaBuildLocalizations(ctx.Context, testID, config.TestflightLocalizations{})
	assert.Error(t, err)
}

func TestUpdateBetaBuildLocalizations_ErrUpdate(t *testing.T) {
	t.Parallel()

	localizations := config.TestflightLocalizations{
		"en-US": {
			Description: "TEST",
			WhatsNew:    "Going away",
		},
	}

	ctx, client := newTestContext(
		response{
			Response: asc.BetaAppLocalizationsResponse{
				Data: []asc.BetaAppLocalization{
					{
						Attributes: &asc.BetaAppLocalizationAttributes{
							Locale: asc.String("en-US"),
						},
						ID: "TEST",
					},
				},
			},
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaBuildLocalizations(ctx.Context, testID, localizations)
	assert.Error(t, err)
}

func TestUpdateBetaBuildLocalizations_ErrCreate(t *testing.T) {
	t.Parallel()

	localizations := config.TestflightLocalizations{
		"en-US": {
			Description: "TEST",
			WhatsNew:    "Going away",
		},
		"ja": {
			Description: "TEST",
			WhatsNew:    "Going away",
		},
	}

	ctx, client := newTestContext(
		response{
			Response: asc.BetaAppLocalizationsResponse{
				Data: []asc.BetaAppLocalization{
					{
						Attributes: &asc.BetaAppLocalizationAttributes{
							Locale: asc.String("en-US"),
						},
						ID: "TEST",
					},
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaBuildLocalizations(ctx.Context, testID, localizations)
	assert.Error(t, err)
}

// Test UpdateBetaLicenseAgreement

func TestUpdateBetaLicenseAgreement_Happy(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":{"id":"TEST"}}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaLicenseAgreement(ctx.Context, testID, config.Testflight{LicenseAgreement: "TEST"})
	assert.NoError(t, err)
}

func TestUpdateBetaLicenseAgreement_NoLicense(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext()
	defer ctx.Close()

	err := client.UpdateBetaLicenseAgreement(ctx.Context, testID, config.Testflight{})
	assert.NoError(t, err)
}

func TestUpdateBetaLicenseAgreement_ErrGet(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaLicenseAgreement(ctx.Context, testID, config.Testflight{LicenseAgreement: "TEST"})
	assert.Error(t, err)
}

func TestUpdateBetaLicenseAgreement_ErrUpdate(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":{"id":"TEST"}}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaLicenseAgreement(ctx.Context, testID, config.Testflight{LicenseAgreement: "TEST"})
	assert.Error(t, err)
}

// Test AssignBetaGroups

func TestAssignBetaGroups_Happy(t *testing.T) {
	t.Parallel()

	testEmail := asc.Email("email2@test.com")
	ctx, client := newTestContext(
		response{
			Response: asc.BetaGroupsResponse{
				Data: []asc.BetaGroup{
					{
						ID: testID,
						Attributes: &asc.BetaGroupAttributes{
							Name: asc.String(testID + "1"),
						},
					},
					{
						ID: testID,
					},
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
		response{
			Response: asc.BetaTestersResponse{
				Data: []asc.BetaTester{
					{
						ID: testID,
						Attributes: &asc.BetaTesterAttributes{
							Email: &testEmail,
						},
					},
				},
			},
		},
		response{
			Response: asc.BetaTestersResponse{
				Data: []asc.BetaTester{
					{
						ID: testID,
						Attributes: &asc.BetaTesterAttributes{
							Email: &testEmail,
						},
						Relationships: &asc.BetaTesterRelationships{
							BetaGroups: &asc.PagedRelationship{
								Data: []asc.RelationshipData{
									{
										ID: testID + "1",
									},
								},
							},
						},
					},
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
	)

	defer ctx.Close()

	err := client.AssignBetaGroups(ctx.Context, testID, testID, []config.BetaGroup{
		{
			Name: testID + "1",
			Testers: []config.BetaTester{
				{
					Email: "email@test.com",
				},
				{
					Email: "email2@test.com",
				},
			},
		},
		{Name: testID + "2"},
	})
	assert.NoError(t, err)
}

func TestAssignBetaGroups_WarnNoGroupsInput(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext()
	defer ctx.Close()

	err := client.AssignBetaGroups(ctx.Context, testID, testID, []config.BetaGroup{})
	assert.NoError(t, err)
}

func TestAssignBetaGroups_ErrList(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.AssignBetaGroups(ctx.Context, testID, testID, []config.BetaGroup{{}})
	assert.Error(t, err)
}

func TestAssignBetaGroups_ErrUpdate(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			Response: asc.BetaGroupsResponse{
				Data: []asc.BetaGroup{
					{
						ID: testID,
						Attributes: &asc.BetaGroupAttributes{
							Name: asc.String(testID + "1"),
						},
					},
				},
			},
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.AssignBetaGroups(ctx.Context, testID, testID, []config.BetaGroup{
		{Name: testID + "1"},
	})
	assert.Error(t, err)
}

func TestAssignBetaGroups_ErrAssign(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			Response: asc.BetaGroupsResponse{
				Data: []asc.BetaGroup{
					{
						ID:         testID,
						Attributes: &asc.BetaGroupAttributes{Name: asc.String(testID)},
					},
				},
			},
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.AssignBetaGroups(ctx.Context, testID, testID, []config.BetaGroup{
		{Name: testID},
	})
	assert.Error(t, err)
}

func TestAssignBetaGroups_ErrCreate(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":[]}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.AssignBetaGroups(ctx.Context, testID, testID, []config.BetaGroup{
		{Name: testID},
	})
	assert.Error(t, err)
}

// Test AssignBetaTesters

func TestAssignBetaTesters_Happy(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			Response: asc.BetaTestersResponse{
				Data: []asc.BetaTester{
					{ID: testID},
					{ID: testID},
				},
			},
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.AssignBetaTesters(ctx.Context, testID, testID, []config.BetaTester{
		{
			Email:     "test@email.com",
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			Email:     "test+1@email.com",
			FirstName: "Jane",
			LastName:  "Doe",
		},
		{
			Email:     "test+2@email.com",
			FirstName: "Joel",
			LastName:  "Doe",
		},
	})
	assert.NoError(t, err)
}

func TestAssignBetaTesters_WarnNoTestersInput(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext()
	defer ctx.Close()

	err := client.AssignBetaTesters(ctx.Context, testID, testID, []config.BetaTester{})
	assert.NoError(t, err)
}

func TestAssignBetaTesters_ErrList(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.AssignBetaTesters(ctx.Context, testID, testID, []config.BetaTester{{}})
	assert.Error(t, err)
}

func TestAssignBetaTesters_WarnNoTestersMatching(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":[]}`,
		},
	)
	defer ctx.Close()

	err := client.AssignBetaTesters(ctx.Context, testID, testID, []config.BetaTester{{}})
	assert.NoError(t, err)
}

func TestAssignBetaTesters_ErrAssign(t *testing.T) {
	t.Parallel()

	testEmail := asc.Email("test@email.com")
	ctx, client := newTestContext(
		response{
			Response: asc.BetaTestersResponse{
				Data: []asc.BetaTester{
					{
						ID: testID,
						Attributes: &asc.BetaTesterAttributes{
							Email: &testEmail,
						},
					},
				},
			},
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)

	defer ctx.Close()

	err := client.AssignBetaTesters(ctx.Context, testID, testID, []config.BetaTester{{}})
	assert.Error(t, err)
}

// Test UpdateBetaReviewDetails

func TestUpdateBetaReviewDetails_Happy(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":{"id":"TEST"}}`,
		},
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaReviewDetails(ctx.Context, testID, config.ReviewDetails{
		Contact:     &config.ContactPerson{},
		DemoAccount: &config.DemoAccount{},
		Attachments: []config.File{
			{Path: "friend"},
		},
	})
	assert.NoError(t, err)
}

func TestUpdateBetaReviewDetails_ErrGet(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaReviewDetails(ctx.Context, testID, config.ReviewDetails{})
	assert.Error(t, err)
}

func TestUpdateBetaReviewDetails_ErrUpdate(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{"data":{"id":"TEST"}}`,
		},
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.UpdateBetaReviewDetails(ctx.Context, testID, config.ReviewDetails{})
	assert.Error(t, err)
}

// Test SubmitBetaApp

func TestSubmitBetaApp_Happy(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.SubmitBetaApp(ctx.Context, testID)
	assert.NoError(t, err)
}

func TestSubmitBetaApp_Err(t *testing.T) {
	t.Parallel()

	ctx, client := newTestContext(
		response{
			StatusCode:  http.StatusNotFound,
			RawResponse: `{}`,
		},
	)
	defer ctx.Close()

	err := client.SubmitBetaApp(ctx.Context, testID)
	assert.Error(t, err)
}
