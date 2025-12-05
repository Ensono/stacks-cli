package util

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

func CheckConnectivity(target string) error {

	var err error

	// define the error that will be displayed if either of the checks fail
	msg := fmt.Sprintf("Cannot connect to '%s', is the machine offline?", target)

	// check that the address can be resolved
	_, err = net.LookupIP(target)
	if err != nil {
		return errors.New(msg)
	}

	// check that the address can be contacted
	resp, err := http.Get(fmt.Sprintf("https://%s", target))
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return errors.New(msg)
	}

	return err
}
