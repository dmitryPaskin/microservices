package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRateLimitExceeded = status.Error(codes.ResourceExhausted, "rate limit per minute exceeded")
)
