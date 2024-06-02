package dns

import (
	"fmt"
	"net/http"
	"os"

	"github.com/akatranlp/akatran/internal/viper"
)

func GetRepoFromViperOrFlag(domain string, provider, token string) (DnsRepository, error) {
	keyStart := "dns::" + domain

	if sub := viper.GetString(keyStart + "::provider"); sub != "" {
		provider = sub
	} else {
		fmt.Fprintln(os.Stderr, "domain not found in config falling back to flag params")
	}

	var repo DnsRepository
	switch provider {
	case CloudflareProvider:
		if sub := viper.GetString(keyStart + "::token"); sub != "" {
			token = sub
		}

		if token == "" {
			return nil, fmt.Errorf("token not provided")
		}
		repo = NewCloudflareRepo(domain, token, http.DefaultClient)
	default:
		return nil, fmt.Errorf("provider not supported: %s", provider)
	}

	return repo, nil
}
