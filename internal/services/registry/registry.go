package registry

import (
	"errors"
	"fmt"
	"net/url"
)

// ParseAndValidateURL parses and validates the registry URL.
func ParseAndValidateURL(registryURL string) (*url.URL, error) {
	host, err := url.Parse(registryURL)
	if err != nil {
		return nil, err
	}
	if host.Scheme == "" || host.Host == "" {
		return nil, errors.Join(fmt.Errorf("Invalid registry URL: %s", registryURL), err)
	}
	return host, nil
}
