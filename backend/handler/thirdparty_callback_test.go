package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	"github.com/teamhanko/hanko/backend/thirdparty"
	"github.com/teamhanko/hanko/backend/utils"
)

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignUp_Google() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	gock.New(thirdparty.GoogleOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GoogleUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GoogleUser{
			ID:            "google_abcde",
			Email:         "test-google-signup@example.com",
			EmailVerified: true,
		})

	cfg := s.setUpConfig([]string{"google"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "google", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-google-signup@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("google", "google_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signup_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_Google() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.GoogleOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GoogleUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GoogleUser{
			ID:            "google_abcde",
			Email:         "test-with-google-identity@example.com",
			EmailVerified: true,
		})

	cfg := s.setUpConfig([]string{"google"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "google", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-google-identity@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("google", "google_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), "", "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignUp_GitHub() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	gock.New(thirdparty.GithubOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GithubUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GithubUser{
			ID:   1234,
			Name: "John Doe",
		})

	gock.New(thirdparty.GitHubEmailsEndpoint).
		Get("/").
		Reply(200).
		JSON([]*thirdparty.GithubUserEmail{
			{
				Email:    "test-github-signup@example.com",
				Primary:  true,
				Verified: true,
			},
		})

	cfg := s.setUpConfig([]string{"github"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "github", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-github-signup@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("github", "1234")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signup_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_GitHub() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.GithubOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GithubUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GithubUser{
			ID:   1234,
			Name: "John Doe",
		})

	gock.New(thirdparty.GitHubEmailsEndpoint).
		Get("/").
		Reply(200).
		JSON([]*thirdparty.GithubUserEmail{
			{
				Email:    "test-with-github-identity@example.com",
				Primary:  true,
				Verified: true,
			},
		})

	cfg := s.setUpConfig([]string{"github"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "github", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-github-identity@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("github", "1234")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignUp_Apple() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	fakeIdToken := s.setUpAppleIdToken("apple_abcde", "fakeClientID", "test-apple-signup@example.com", true, false)
	gock.New(thirdparty.AppleTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken", "id_token": fakeIdToken})

	fakeJwkSet := s.setUpFakeJwkSet()
	gock.New(thirdparty.AppleKeysEndpoint).
		Get("/").
		Reply(200).
		JSON(fakeJwkSet)

	cfg := s.setUpConfig([]string{"apple"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "apple", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-apple-signup@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("apple", "apple_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signup_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_Apple() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	fakeIdToken := s.setUpAppleIdToken("apple_abcde", "fakeClientID", "test-with-apple-identity@example.com", true, false)
	gock.New(thirdparty.AppleTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken", "id_token": fakeIdToken})

	fakeJwkSet := s.setUpFakeJwkSet()
	gock.New(thirdparty.AppleKeysEndpoint).
		Get("/").
		Reply(200).
		JSON(fakeJwkSet)

	cfg := s.setUpConfig([]string{"apple"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "apple", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-apple-identity@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("apple", "apple_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_Apple_WithBooleanEmailVerifiedClaim() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	fakeIdToken := s.setUpAppleIdToken("apple_abcde", "fakeClientID", "test-with-apple-identity@example.com", true, true)
	gock.New(thirdparty.AppleTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken", "id_token": fakeIdToken})

	fakeJwkSet := s.setUpFakeJwkSet()
	gock.New(thirdparty.AppleKeysEndpoint).
		Get("/").
		Reply(200).
		JSON(fakeJwkSet)

	cfg := s.setUpConfig([]string{"apple"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "apple", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-apple-identity@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("apple", "apple_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignUp_Discord() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	gock.New(thirdparty.DiscordOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.DiscordUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.DiscordUser{
			ID:       "discord_abcde",
			Email:    "test-discord-signup@example.com",
			Verified: true,
		})

	cfg := s.setUpConfig([]string{"discord"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "discord", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-discord-signup@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("discord", "discord_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signup_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_Discord() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.DiscordOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.DiscordUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.DiscordUser{
			ID:       "discord_abcde",
			Email:    "test-with-discord-identity@example.com",
			Verified: true,
		})

	cfg := s.setUpConfig([]string{"discord"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "discord", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-discord-identity@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("discord", "discord_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), "", "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignUp_Microsoft() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	fakeIdToken := s.setUpMicrosoftIdToken("microsoft_abcde", "fakeClientID", "test-microsoft-signup@example.com", true)
	gock.New(thirdparty.MicrosoftOAuthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken", "id_token": fakeIdToken})

	fakeJwkSet := s.setUpFakeJwkSet()
	gock.New(thirdparty.MicrosoftKeysEndpoint).
		Get("/").
		Reply(200).
		JSON(fakeJwkSet)

	cfg := s.setUpConfig([]string{"microsoft"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "microsoft", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-microsoft-signup@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("microsoft", "microsoft_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signup_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_Microsoft() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	fakeIdToken := s.setUpMicrosoftIdToken("microsoft_abcde", "fakeClientID", "test-with-microsoft-identity@example.com", true)
	gock.New(thirdparty.MicrosoftOAuthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken", "id_token": fakeIdToken})

	fakeJwkSet := s.setUpFakeJwkSet()
	gock.New(thirdparty.MicrosoftKeysEndpoint).
		Get("/").
		Reply(200).
		JSON(fakeJwkSet)

	cfg := s.setUpConfig([]string{"microsoft"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "microsoft", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-microsoft-identity@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("microsoft", "microsoft_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), "", "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignUp_Facebook() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	gock.New(thirdparty.FacebookOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.FacebookUserInfoEndpoint).
		Get("/me").
		Reply(200).
		JSON(&thirdparty.FacebookUser{
			ID:    "facebook_abcde",
			Email: "test-facebook-signup@example.com",
		})

	cfg := s.setUpConfig([]string{"facebook"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "facebook", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-facebook-signup@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("facebook", "facebook_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signup_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_Facebook() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.FacebookOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.FacebookUserInfoEndpoint).
		Get("/me").
		Reply(200).
		JSON(&thirdparty.FacebookUser{
			ID:    "facebook_abcde",
			Email: "test-with-facebook-identity@example.com",
		})

	cfg := s.setUpConfig([]string{"facebook"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "facebook", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-facebook-identity@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("facebook", "facebook_abcde")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), "", "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignUp_WithUnclaimedEmail() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.GoogleOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GoogleUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GoogleUser{
			ID:            "google_unclaimed_email",
			Email:         "unclaimed-email@example.com",
			EmailVerified: true,
		})

	cfg := s.setUpConfig([]string{"google"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "google", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("unclaimed-email@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.True(email.IsPrimary())

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("google", "google_unclaimed_email")
		s.NotNil(identity)

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signup_succeeded"}, user.ID.String(), email.Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_ProviderEMailChangedToExistingEmail() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.GoogleOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GoogleUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GoogleUser{
			ID:            "google_abcde",
			Email:         "test-with-google-identity-changed@example.com",
			EmailVerified: true,
		})

	cfg := s.setUpConfig([]string{"google"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "google", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-google-identity-changed@example.com")
		s.NoError(err)
		s.NotNil(email)

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("google", "google_abcde")
		s.NotNil(identity)
		s.Equal("test-with-google-identity-changed@example.com", identity.Data["email"])

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), user.Emails.GetPrimary().Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_ProviderEMailChangedToUnclaimedEmail() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.GoogleOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GoogleUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GoogleUser{
			ID:            "google_abcde",
			Email:         "unclaimed-email@example.com",
			EmailVerified: true,
		})

	cfg := s.setUpConfig([]string{"google"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "google", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("unclaimed-email@example.com")
		s.NoError(err)
		s.NotNil(email)

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)

		identity := email.Identities.GetIdentity("google", "google_abcde")
		s.NotNil(identity)
		s.Equal("unclaimed-email@example.com", identity.Data["email"])

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), user.Emails.GetPrimary().Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_SignIn_ProviderEMailChangedToNonExistentEmail() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.GoogleOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GoogleUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GoogleUser{
			ID:            "google_abcde",
			Email:         "non-existent-email@example.com",
			EmailVerified: true,
		})

	cfg := s.setUpConfig([]string{"google"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "google", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("non-existent-email@example.com")
		s.NoError(err)
		s.NotNil(email)

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)
		s.Len(user.Emails, 3)

		identity := email.Identities.GetIdentity("google", "google_abcde")
		s.NotNil(identity)
		s.Equal("non-existent-email@example.com", identity.Data["email"])

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_signin_succeeded"}, user.ID.String(), user.Emails.GetPrimary().Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_Link_ExistingAccountNoIdentities() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.GoogleOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GoogleUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GoogleUser{
			ID:            "google_1234",
			Email:         "test-no-identity@example.com",
			EmailVerified: true,
		})

	cfg := s.setUpConfig([]string{"google"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "google", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-no-identity@example.com")
		s.NoError(err)
		s.NotNil(email)

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)
		s.Len(user.Emails, 1)

		identity := email.Identities.GetIdentity("google", "google_1234")
		s.NotNil(identity)
		s.Equal("test-no-identity@example.com", identity.Data["email"])

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_linking_succeeded"}, user.ID.String(), user.Emails.GetPrimary().Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}

func (s *thirdPartySuite) TestThirdPartyHandler_Callback_Link_GoogleToAccountWithGithubIdentity() {
	defer gock.Off()
	if testing.Short() {
		s.T().Skip("skipping test in short mode.")
	}

	err := s.LoadFixtures("../test/fixtures/thirdparty")
	s.NoError(err)

	gock.New(thirdparty.GoogleOauthTokenEndpoint).
		Post("/").
		Reply(200).
		JSON(map[string]string{"access_token": "fakeAccessToken"})

	gock.New(thirdparty.GoogleUserInfoEndpoint).
		Get("/").
		Reply(200).
		JSON(&thirdparty.GoogleUser{
			ID:            "google_1234",
			Email:         "test-with-github-identity@example.com",
			EmailVerified: true,
		})

	cfg := s.setUpConfig([]string{"google", "github"}, []string{"https://example.com"})

	state, err := thirdparty.GenerateState(cfg, "google", "https://example.com")
	s.NoError(err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/thirdparty/callback?code=abcde&state=%s", state), nil)
	req.AddCookie(&http.Cookie{
		Name:  utils.HankoThirdpartyStateCookie,
		Value: string(state),
	})

	c, rec := s.setUpContext(req)
	handler := s.setUpHandler(cfg)

	if s.NoError(handler.Callback(c)) {
		s.Equal(http.StatusTemporaryRedirect, rec.Code)

		s.assertLocationHeaderHasToken(rec)
		s.assertStateCookieRemoved(rec)

		email, err := s.Storage.GetEmailPersister().FindByAddress("test-with-github-identity@example.com")
		s.NoError(err)
		s.NotNil(email)
		s.Len(email.Identities, 2)

		user, err := s.Storage.GetUserPersister().Get(*email.UserID)
		s.NoError(err)
		s.NotNil(user)
		s.Len(user.Emails, 1)

		identity := email.Identities.GetIdentity("google", "google_1234")
		s.NotNil(identity)
		s.Equal("test-with-github-identity@example.com", identity.Data["email"])

		logs, lerr := s.Storage.GetAuditLogPersister().List(0, 0, nil, nil, []string{"thirdparty_linking_succeeded"}, user.ID.String(), user.Emails.GetPrimary().Address, "", "")
		s.NoError(lerr)
		s.Len(logs, 1)
	}
}
