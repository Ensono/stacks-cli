package util

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

func CheckConnectivity(target string) error {

	var err error

	// check that the address can be resolved
	_, err = net.LookupIP(target)
	if err != nil {
		return fmt.Errorf("cannot connect to '%s', is the machine offline?: %w", target, err)
	}

	// check that the address can be contacted
	resp, err := http.Get(fmt.Sprintf("https://%s", target))
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New("cannot connect to '" + target + "', is the machine offline?")
	}

	return err
}
