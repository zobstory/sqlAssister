package utils

import "errors"

func QueryCheckerWithArgs(query string, args ...any) error {
	switch {

	case len(query) == 0 && len(args) == 0:
		return errors.New("both query / statement & args are not present")

	case len(query) == 0 && len(args) > 0:
		return errors.New("no query / statement present")

	case len(query) > 0 && len(args) == 0:
		return errors.New("no args present")

	default:
		return nil
	}
}

func QueryChecker(query string) error {
	if len(query) == 0 {
		return errors.New("no query / statement present")
	}

	return nil
}
