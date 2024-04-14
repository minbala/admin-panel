package resources

import "github.com/cockroachdb/errors"

var ErrUnmarshalData = errors.New("invalid json data")
var ErrNoRecordFound = errors.New("has not been found")
var ErrClient = errors.New("invalid request data")
var ErrNotAcceptable = errors.New("request is not acceptable")
var ErrInternalServer = errors.New("we are fixing the issue. please, be patient")
var ErrUnAuthorized = errors.New("unauthorized access")
var ErrForbidden = errors.New("forbidden")
var ErrBadGateway = errors.New("bad gateway")
var ErrDuplicateValue = errors.New("duplicate value")
