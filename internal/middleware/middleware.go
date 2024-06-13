package middleware

import (
	"context"
	crand "crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"goth/internal/store"
	"log"
    "math/rand/v2"
	"net/http"
	// "strings"
)

type key string

var NonceKey key = "nonces"

type Nonces struct {
	Htmx            string
	ResponseTargets string
	Tw              string
	HtmxCSSHash     string
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := crand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func GenerateRandomStringOfNumbers(length int) string {
    bytes := make([]byte, length)
    for i := range length {
        // Generate a random num using math/rand
        newInt := rand.IntN(10)
        intAsByte := byte(newInt) + byte('0')
        bytes[i] = intAsByte
    }
	return string(bytes)
}

func CSPMiddleware(next http.Handler) http.Handler {
	// To use the same nonces in all responses, move the Nonces
	// struct creation to here, outside the handler.

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new Nonces struct for every request when here.
		// move to outside the handler to use the same nonces in all responses
		nonceSet := Nonces{
			Htmx:            generateRandomString(16),
			ResponseTargets: generateRandomString(16),
			Tw:              generateRandomString(16),
			HtmxCSSHash:     "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg=",
		}

		// set nonces in context
		ctx := context.WithValue(r.Context(), NonceKey, nonceSet)
		// insert the nonces into the content security policy header

		// cspHeader := fmt.Sprintf("default-src 'self'; script-src 'nonce-%s' 'nonce-%s' ; style-src 'nonce-%s' '%s';",
		// 	nonceSet.Htmx,
		// 	nonceSet.ResponseTargets,
		// 	nonceSet.Tw,
		// 	nonceSet.HtmxCSSHash)

		// w.Header().Set("Content-Security-Policy", cspHeader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TextHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// get the Nonce from the context, it is a struct called Nonces,
// so we can get the nonce we need by the key, i.e. HtmxNonce
func GetNonces(ctx context.Context) Nonces {
	nonceSet := ctx.Value(NonceKey)
	if nonceSet == nil {
		log.Fatal("error getting nonce set - is nil")
	}

	nonces, ok := nonceSet.(Nonces)

	if !ok {
		log.Fatal("error getting nonce set - not ok")
	}

	return nonces
}

func GetHtmxNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)

	return nonceSet.Htmx
}

func GetResponseTargetsNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.ResponseTargets
}

func GetTwNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.Tw
}

type AuthMiddleware struct {
	sessionStore      store.SessionStore
	sessionCookieName string
}

func NewAuthMiddleware(sessionStore store.SessionStore, sessionCookieName string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore:      sessionStore,
		sessionCookieName: sessionCookieName,
	}
}

type SessionContextKey string

var SessionKey SessionContextKey = "session"

func (m *AuthMiddleware) AddSessionToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// sessionCookie, err := r.Cookie(m.sessionCookieName)

		sessionCookie, err := r.Cookie("session")

		if err != nil {
			fmt.Println("error getting session cookie", err)
			next.ServeHTTP(w, r)
			return
		}
        decodedValue, err := b64.StdEncoding.DecodeString(sessionCookie.Value)

        if err != nil {
            next.ServeHTTP(w, r)
            return
        }

        sessionID := string(decodedValue)

        ctx := context.WithValue(r.Context(), SessionKey, sessionID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// func (m *AuthMiddleware) AddUserToContext(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
// 		sessionCookie, err := r.Cookie(m.sessionCookieName)
//
// 		if err != nil {
// 			fmt.Println("error getting session cookie", err)
// 			next.ServeHTTP(w, r)
// 			return
// 		}
//
// 		decodedValue, err := b64.StdEncoding.DecodeString(sessionCookie.Value)
//
// 		if err != nil {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
//
// 		splitValue := strings.Split(string(decodedValue), ":")
//
// 		if len(splitValue) != 2 {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
//
// 		sessionID := splitValue[0]
// 		userID := splitValue[1]
//
//
// 		user, err := m.sessionStore.GetUserFromSession(sessionID, userID)
//
// 		if err != nil {
// 			next.ServeHTTP(w, r)
// 			return
// 		}
//
// 		ctx := context.WithValue(r.Context(), UserKey, user)
//
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// func GetUser(ctx context.Context) *store.User {
// 	user := ctx.Value(UserKey)
// 	if user == nil {
// 		return nil
// 	}
//
// 	return user.(*store.User)
// }
