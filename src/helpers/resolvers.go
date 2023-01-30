package helpers

import (
	"errors"
)

func ResolveStatusCode(code int) (int, error) {
	if code >= 1000 && code < 2000 {
		return 200, nil
	}

	if code >= 2000 && code < 5000 {
		if code == 2006 || code == 2052 {
			return 429, nil
		}

		if code == 2050 || code == 2051 {
			return 403, nil
		}

		if code == 2214 || code == 2316 {
			return 409, nil
		}

		return 400, nil
	}

	if code >= 5000 {
		return 500, nil
	}

	return 500, errors.New("failed to resolve status code")
}
