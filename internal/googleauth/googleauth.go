package googleauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"go.uber.org/zap"
)

// OAuthService handles Google OAuth2 authentication.
// It holds the oauth2.Config and a logger.
type OAuthService struct {
	Config *oauth2.Config
	Logger *zap.Logger
	// stateStore is a simple in-memory store for OAuth state tokens.
	// For production, consider a more robust distributed store (e.g., Redis) if running multiple instances.
	stateStore map[string]time.Time
}

const (
	stateTokenExpiry   = 10 * time.Minute // OAuth state token expires in 10 minutes
	stateCookieName    = "oauthstate"
	redirectURLDev     = "http://localhost:8080/api/auth/google/callback" // Default, adjust if your port differs
	// Add production redirect URL when ready: const redirectURLProd = "https://ana.world/api/auth/google/callback"
)

// NewOAuthService creates a new OAuthService.
// credPath should be the path to your credentials.json file.
// redirectURL should be the callback URL registered in Google Cloud Console.
// scopes are the Google API scopes you are requesting.
func NewOAuthService(credPath string, redirectURL string, logger *zap.Logger, scopes []string) (*OAuthService, error) {
	absCredPath, err := filepath.Abs(credPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for credentials: %w", err)
	}

	b, err := ioutil.ReadFile(absCredPath)
	if err != nil {
		logger.Error("Unable to read client secret file", zap.String("path", absCredPath), zap.Error(err))
		return nil, fmt.Errorf("unable to read client secret file at %s: %w", absCredPath, err)
	}

	oauthConfig, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		logger.Error("Unable to parse client secret file to config", zap.Error(err))
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	oauthConfig.RedirectURL = redirectURL

	return &OAuthService{
		Config:     oauthConfig,
		Logger:     logger.Named("OAuthService"),
		stateStore: make(map[string]time.Time),
	}, nil
}

// generateStateToken creates a random string for CSRF protection and stores it.
func (s *OAuthService) generateStateToken(c *gin.Context) (string, error) {
	// Generate a random state token for CSRF protection
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		s.Logger.Error("Failed to generate state token bytes", zap.Error(err))
		return "", fmt.Errorf("failed to generate state token: %w", err)
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)

	// Store the state token with an expiry in the service's stateStore
	s.stateStore[state] = time.Now().Add(stateTokenExpiry)
	s.Logger.Info("Generated and stored state token", zap.String("state", state))

	// Clean up expired state tokens (simple cleanup, can be improved for efficiency)
	go s.cleanupExpiredStates()

	// Also set as a cookie for potential stateless verification, though primary check is server-side store
	c.SetCookie(stateCookieName, state, int(stateTokenExpiry.Seconds()), "/api/auth/google", "localhost", false, true)

	return state, nil
}

// validateStateToken checks if the provided state is valid and removes it from the store.
func (s *OAuthService) validateStateToken(state string) bool {
	expiry, ok := s.stateStore[state]
	if !ok {
		s.Logger.Warn("State token not found in store", zap.String("state", state))
		return false
	}
	delete(s.stateStore, state) // Consume the state token

	if time.Now().After(expiry) {
		s.Logger.Warn("State token expired", zap.String("state", state))
		return false
	}
	s.Logger.Info("State token validated successfully", zap.String("state", state))
	return true
}

// cleanupExpiredStates iterates through the state store and removes expired tokens.
func (s *OAuthService) cleanupExpiredStates() {
	s.Logger.Debug("Running state token cleanup")
	for st, expiry := range s.stateStore {
		if time.Now().After(expiry) {
			delete(s.stateStore, st)
			s.Logger.Debug("Deleted expired state token", zap.String("state", st))
		}
	}
}

// HandleLogin redirects the user to Google's OAuth 2.0 consent page.
func (s *OAuthService) HandleLogin(c *gin.Context) {
	state, err := s.generateStateToken(c)
	if err != nil {
		s.Logger.Error("Failed to generate state for login", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate authentication. Please try again."})
		return
	}

	// ApprovalForce ensures the user is prompted for consent every time during development.
	// Remove oauth2.ApprovalForce for a smoother UX in production after initial testing,
	// especially if AccessTypeOffline is used and you have a refresh token.
	authURL := s.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	s.Logger.Info("Redirecting to Google for authentication", zap.String("url", authURL))
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// HandleCallback handles the OAuth 2.0 callback from Google.
// It exchanges the authorization code for tokens.
func (s *OAuthService) HandleCallback(c *gin.Context) {
	stateFromQuery := c.Query("state")
	code := c.Query("code")

	// Validate state token from query against server-side store
	if !s.validateStateToken(stateFromQuery) {
		// As a fallback, check cookie if state wasn't in query (though it should be)
		stateFromCookie, cookieErr := c.Cookie(stateCookieName)
		if cookieErr != nil || !s.validateStateToken(stateFromCookie) {
			s.Logger.Error("Invalid or missing state token during callback", zap.String("queryState", stateFromQuery), zap.String("cookieState", stateFromCookie))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state token. Authentication failed."})
			return
		}
	}
	// Clear the cookie after use
	c.SetCookie(stateCookieName, "", -1, "/api/auth/google", "localhost", false, true)

	if code == "" {
		errorDesc := c.Query("error_description")
		if errorDesc == "" {
			errorDesc = "Authorization code not found in callback from Google."
		}
		s.Logger.Error("Failed to get authorization code from Google", zap.String("error_description", errorDesc))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get authorization code: " + errorDesc})
		return
	}

	s.Logger.Info("Received authorization code", zap.String("code", code))

	// Exchange code for token
	token, err := s.Config.Exchange(context.Background(), code)
	if err != nil {
		s.Logger.Error("Failed to exchange authorization code for token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token. " + err.Error()})
		return
	}

	s.Logger.Info("Successfully exchanged code for token",
		zap.String("accessToken", token.AccessToken),
		zap.Bool("hasRefreshToken", token.RefreshToken != ""),
		zap.Time("expiry", token.Expiry),
	)

	// TODO: Securely store the token (token.AccessToken and token.RefreshToken if present)
	// For now, just return a success message with the access token (for demonstration)
	// In a real app, you would associate this token with a user session/ID.
	c.JSON(http.StatusOK, gin.H{
		"message":        "Authentication successful! Tokens obtained.",
		"access_token":   token.AccessToken, // DO NOT expose this in production like this directly to client for general use
		"refresh_token":  token.RefreshToken, // Especially DO NOT expose this. Store server-side only.
		"token_type":     token.TokenType,
		"expiry":         token.Expiry,
		"scope_returned": token.Extra("scope"),
	})
}
