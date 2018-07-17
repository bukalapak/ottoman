package adjust

import "net/http"

type Requester func(r *http.Request) error

func Request(r *http.Request, options ...Requester) error {
	for _, opt := range options {
		if err := opt(r); err != nil {
			return err
		}
	}

	return nil
}
