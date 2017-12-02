package sessions

import (
	"strings"
	"fmt"
	"errors"
	"net/http"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, fmt.Errorf("sessions ID can not be empty")
	}
	sessionID, err := NewSessionID(signingKey)
	if err != nil {
		return sessionID, fmt.Errorf("invalid session id: %v", err)
	} else {
		store.Save(sessionID, sessionState)
		w.Header().Add(headerAuthorization, schemeBearer + string(sessionID))
		return sessionID, nil
	}
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	if len(signingKey) == 0 {
		return InvalidSessionID, fmt.Errorf("sessions ID can not be empty")
	}
	if len(r.Header.Get(headerAuthorization)) == 0 {
		authParam := r.URL.Query().Get("auth")
		header := strings.Split(authParam, " ")
		if len(header) == 2 {
			sessionID, err := ValidateID(header[1], signingKey)
			if err != nil || header[0] != strings.TrimSpace(schemeBearer) {
				return sessionID, fmt.Errorf("not a valid sessionID: %v", err)
			} else {
				return sessionID, nil
			}
		} else {
			return InvalidSessionID, ErrNoSessionID
		}
	} else {
		value := r.Header.Get(headerAuthorization)
		header := strings.Split(value, " ")
		if len(header) == 2 {
			sessionID, err := ValidateID(header[1], signingKey)
			if err != nil || header[0] != strings.TrimSpace(schemeBearer) {
				return InvalidSessionID, ErrInvalidScheme
			} else {
				return sessionID, nil
			}
		} else {
			return InvalidSessionID, ErrInvalidID
		}
	}
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	sessionID, err := GetSessionID(r, signingKey) 
	if err != nil {
		return sessionID, fmt.Errorf("not a valid sessionID: %v", err)
	} else {
		err := store.Get(sessionID, sessionState)
		if err != nil {
			return sessionID, ErrStateNotFound
		}
		return sessionID, nil
	}
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	sessionID, err := GetSessionID(r, signingKey) 
	if err != nil {
		return sessionID, fmt.Errorf("not a valid sessionID: %v", err)
	} else {
		store.Delete(sessionID)
		return sessionID, nil
	}
}
