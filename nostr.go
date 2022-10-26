package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fiatjaf/go-lnurl"
)

func GetNostrPubkey(name, domain string) (string, error) {
	return "e0d05a5b8c7789eb83f87672f4eb0dca78f99292ab038e5c66f84d97d77b95ae", nil
}

func handleNostr(w http.ResponseWriter, r *http.Request) {
	// username := mux.Vars(r)["user"]
	username := r.URL.Query().Get("name")

	// we can remove the following because there can be only one domain
	// for your nip05 id as of now
	domains := getDomains(s.Domain)
	domain := ""

	if len(domains) == 1 {
		domain = domains[0]
	} else {
		hostname := r.URL.Host
		if hostname == "" {
			hostname = r.Host
		}

		for _, one := range getDomains(s.Domain) {
			if strings.Contains(hostname, one) {
				domain = one
				break
			}
		}
		if domain == "" {
			json.NewEncoder(w).Encode(lnurl.ErrorResponse("incorrect domain"))
			return
		}
	}

	pubKey, err := GetNostrPubkey(username, domain)
	if err != nil {
		log.Error().Err(err).Str("name", username).Str("domain", domain).Msg("failed to get name")
		json.NewEncoder(w).Encode(lnurl.ErrorResponse(fmt.Sprintf(
			"failed to get name %s@%s", username, domain)))
		return
	}

	log.Info().Str("username", username).Str("domain", domain).Msg("got nostr nip-05 request")

	json.NewEncoder(w).Encode(fmt.Sprintf(
		`{"names":{"%s": "%s"}}`, username, pubKey,
	))
}
