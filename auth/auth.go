package auth

import (
	"net/http"
	"net/url"
)

// Authenticator handles request authentication for the Kraken API.
type Authenticator struct {
	signer *Signer
	nonce  NonceGenerator
	otp    func() string // Optional 2FA OTP generator
}

// NewAuthenticator creates a new Authenticator with the given credentials.
func NewAuthenticator(apiKey, apiSecret string, opts ...AuthOption) (*Authenticator, error) {
	signer, err := NewSigner(apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	a := &Authenticator{
		signer: signer,
		nonce:  NewTimestampNonceGenerator(),
	}

	// Apply options
	for _, opt := range opts {
		opt(a)
	}

	return a, nil
}

// AuthOption configures the Authenticator.
type AuthOption func(*Authenticator)

// WithOTP sets the OTP generator for 2FA authentication.
// The generator function will be called for each authenticated request.
func WithOTP(generator func() string) AuthOption {
	return func(a *Authenticator) {
		a.otp = generator
	}
}

// WithNonceGenerator sets a custom nonce generator.
func WithNonceGenerator(g NonceGenerator) AuthOption {
	return func(a *Authenticator) {
		a.nonce = g
	}
}

// APIKey returns the public API key.
func (a *Authenticator) APIKey() string {
	return a.signer.APIKey()
}

// GenerateNonce generates a new nonce value.
func (a *Authenticator) GenerateNonce() string {
	return a.nonce.Next()
}

// GenerateOTP returns a one-time password (OTP) if configured.
// If no OTP generator is set, it returns an empty string.
func (a *Authenticator) GenerateOTP() string {
	if a.otp == nil {
		return ""
	}
	return a.otp()
}

// Authenticate adds authentication headers and parameters to a request.
//
// This method:
//  1. Generates a nonce if not already present in params
//  2. Adds the OTP if configured
//  3. Generates the API-Sign signature
//  4. Sets the API-Key and API-Sign headers
//
// Parameters:
//   - req: The HTTP request to authenticate
//   - path: The URI path for signature calculation
//   - params: The request parameters (will be modified to include nonce and possibly OTP)
func (a *Authenticator) Authenticate(req *http.Request, path string, params url.Values) error {
	// Generate nonce if not present
	nonce := params.Get("nonce")
	if nonce == "" {
		nonce = a.nonce.Next()
		params.Set("nonce", nonce)
	}

	// Add OTP if configured
	if a.otp != nil {
		otp := a.otp()
		if otp != "" {
			params.Set("otp", otp)
		}
	}

	// Generate signature
	signature, err := a.signer.Sign(path, nonce, params)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("API-Key", a.signer.APIKey())
	req.Header.Set("API-Sign", signature)

	return nil
}

// AuthenticateJSON adds authentication headers for a JSON request.
//
// Parameters:
//   - req: The HTTP request to authenticate
//   - path: The URI path for signature calculation
//   - nonce: The nonce value (must be included in the JSON body)
//   - jsonBody: The JSON request body as a string
func (a *Authenticator) AuthenticateJSON(req *http.Request, path string, nonce string, jsonBody string) error {
	// Generate signature
	signature, err := a.signer.SignJSON(path, nonce, jsonBody)
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("API-Key", a.signer.APIKey())
	req.Header.Set("API-Sign", signature)

	return nil
}

// HasOTP returns true if an OTP generator is configured.
func (a *Authenticator) HasOTP() bool {
	return a.otp != nil
}
