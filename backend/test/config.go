package test

import "github.com/teamhanko/hanko/backend/config"

var DefaultConfig = config.Config{
	Webauthn: config.WebauthnSettings{
		RelyingParty: config.RelyingParty{
			Id:          "localhost",
			DisplayName: "Test Relying Party",
			Icon:        "",
			Origins:     []string{"http://localhost:8080", "http://localhost:8888"},
		},
		Timeout:          60000,
		UserVerification: "preferred",
	},
	Secrets: config.Secrets{
		Keys: []string{"abcdefghijklmnop"},
	},
	Email: config.Email{
		Enabled:              true,
		UseForAuthentication: true,
	},
	EmailDelivery: config.EmailDelivery{
		Enabled: true,
		SMTP: config.SMTP{
			Host: "localhost",
			Port: "2500",
		},
		FromAddress: "test@hanko.io",
		FromName:    "Hanko Test",
	},
	Passcode: config.Passcode{
		TTL: 300,
	},
	Session: config.Session{
		Lifespan: "1h",
		Cookie: config.Cookie{
			SameSite: "none",
		},
		Limit: 5,
	},
	Service: config.Service{
		Name: "Test",
	},
	Account: config.Account{
		AllowSignup:   true,
		AllowDeletion: false,
	},
	Passkey: config.Passkey{
		Enabled:          true,
		UserVerification: "preferred",
	},
}
