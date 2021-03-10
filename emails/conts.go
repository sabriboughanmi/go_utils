package emails


import "errors"

const (
	MaxLineLength      = 76                             // MaxLineLength is the maximum line length per RFC 2045
	defaultContentType = "text/plain; charset=us-ascii" // defaultContentType is the default Content-Type according to RFC 2045, section 5.2
)


// ErrMissingBoundary is returned when there is no boundary given for a multipart entity
var ErrMissingBoundary = errors.New("No boundary found for multipart entity")

// ErrMissingContentType is returned when there is no "Content-Type" header for a MIME entity
var ErrMissingContentType = errors.New("No Content-Type found for MIME entity")

var ErrSenderMustBeSpecified = errors.New("sender Must be Specified")

var ErrReceiversMustBeSpecified = errors.New("at least one receiver must be specified")