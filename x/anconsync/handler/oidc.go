package handler

import (
	"io"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"

	"golang.org/x/oauth2"

	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anconprotocol/sdk"

	"github.com/gin-gonic/gin"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
)

type OidcHandler struct {
	*sdk.AnconSyncContext

	Provider *oidc.Provider
	Config   oauth2.Config
	Verifier *oidc.IDTokenVerifier
}

func NewOidcHandler(ctx *sdk.AnconSyncContext, clientId string, clientSecretKey string, redirectURL string) *OidcHandler {
	clientID = clientId
	clientSecret = clientSecretKey
	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		panic(err)
	}
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &OidcHandler{
		AnconSyncContext: ctx,
		Provider:         provider,
		Config:           config,
		Verifier:         verifier,
	}
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func (ctx *OidcHandler) OIDCRequest(c *gin.Context) {

	state, err := randString(16)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("an error occurred while requesting token").Error(),
		})
		return
	}
	nonce, err := randString(16)
	if err != nil {
		c.JSON(400, gin.H{
			"error": fmt.Errorf("an error occurred while requesting token").Error(),
		})
		return
	}

	// fqdn, err := fqdn.FqdnHostname()
	// if err != nil {
	// 	c.JSON(400, gin.H{
	// 		"error": fmt.Errorf("an error occurred while requesting token").Error(),
	// 	})
	// 	return
	// }

	c.SetCookie("state", state, 1000000, "/", "", true, false)
	c.SetCookie("nonce", nonce, 1000000, "/", "", true, false)

	c.Redirect(http.StatusFound, ctx.Config.AuthCodeURL(state, oidc.Nonce(nonce)))

}
func (ctx *OidcHandler) OIDCCallback(c *gin.Context) {

	state, err := c.Cookie("state")
	if err != nil {
		c.Data(400, "text/text", []byte("state not found"))
		return

	}
	if c.Query("state") != state {
		c.Data(400, "text/text", []byte("state did not match"))
		return
	}

	oauth2Token, err := ctx.Config.Exchange(c.Request.Context(), c.Query("code"))
	if err != nil {
		c.Data(500, "text/text", []byte("Failed to exchange token: "+err.Error()))
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.Data(500, "text/text", []byte("No id_token field in oauth2 token."))
		return

	}
	idToken, err := ctx.Verifier.Verify(c.Request.Context(), rawIDToken)
	if err != nil {
		c.Data(500, "text/text", []byte("Failed to verify ID Token: "+err.Error()))
		return
	}

	nonce, err := c.Cookie("nonce")
	if err != nil {
		c.Data(400, "text/text", []byte("nonce not found"))
		return

	}
	if idToken.Nonce != nonce {
		c.Data(400, "text/text", []byte("nonce did not match"))
		return

	}

	oauth2Token.AccessToken = "*REDACTED*"

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		c.Data(500, "text/text", []byte("Failed to verify ID Token: "+err.Error()))
		return
	}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		c.Data(500, "text/text", []byte("Failed to verify ID Token: "+err.Error()))
		return
	}
	c.Data(200, "application/json", data)
}
