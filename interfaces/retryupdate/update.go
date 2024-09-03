//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
InnerLoop:
	for {
		var oldValue *string
		var resp *kvapi.GetResponse
		for {
			var err error
			resp, err = c.Get(&kvapi.GetRequest{Key: key})
			var apiErr *kvapi.APIError
			if errors.Is(err, &kvapi.AuthError{}) {
				return err
			} else if errors.As(err, &apiErr) {
				var authErr *kvapi.AuthError
				if errors.As(apiErr, &authErr) {
					return err
				} else if errors.Is(apiErr, kvapi.ErrKeyNotFound) {
					oldValue = nil
					break
				}
				continue
			} else if err != nil {
				continue
			} else {
				oldValue = &resp.Value
				break
			}
		}

		newValue, err := updateFn(oldValue)
		if err != nil {
			return err
		}

		newVersion := uuid.Must(uuid.NewV4())

		for {
			var apiErr *kvapi.APIError
			if oldValue == nil {
				_, err := c.Set(&kvapi.SetRequest{Key: key, Value: newValue, OldVersion: uuid.UUID{0}, NewVersion: newVersion})
				if errors.Is(err, &kvapi.ConflictError{}) {
					continue InnerLoop
				} else if errors.As(err, &apiErr) {
					var authErr *kvapi.AuthError
					var conflictErr *kvapi.ConflictError
					if errors.As(apiErr, &authErr) {
						return err
					} else if errors.As(apiErr, &conflictErr) {
						if conflictErr.ExpectedVersion == newVersion {
							return nil
						}
						continue InnerLoop
					} else if errors.Is(apiErr, kvapi.ErrKeyNotFound) {
						oldValue = nil
						var err error
						newValue, err = updateFn(oldValue)
						if err != nil {
							return err
						}
					}
					continue
				} else if err != nil {
					continue
				}
				return nil
			} else {
				_, err := c.Set(&kvapi.SetRequest{Key: key, Value: newValue, OldVersion: resp.Version, NewVersion: newVersion})
				if errors.Is(err, &kvapi.ConflictError{}) {
					continue InnerLoop
				} else if errors.As(err, &apiErr) {
					var authErr *kvapi.AuthError
					var conflictErr *kvapi.ConflictError
					if errors.As(apiErr, &authErr) {
						return err
					} else if errors.As(apiErr, &conflictErr) {
						if conflictErr.ExpectedVersion == newVersion {
							return nil
						}
						continue InnerLoop
					} else if errors.Is(apiErr, kvapi.ErrKeyNotFound) {
						oldValue = nil
						var err error
						newValue, err = updateFn(oldValue)
						if err != nil {
							return err
						}
					}
					continue
				} else if err != nil {
					continue
				}
				return nil
			}
		}
	}
}
