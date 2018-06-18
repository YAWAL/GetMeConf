package errortype

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotFound - return gRPC not found error
func NotFound(msg string) error {
	return status.Errorf(codes.NotFound, fmt.Sprintf("Requested item(s) not found: %s", msg))
}

// Internal - return gRPC internal error
func Internal(msg string) error {
	return status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %s", msg))
}

// GrpcError - translates a shared-service error into an appropriate gRPC response.
func GrpcError(err error, msg string) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return NotFound(msg)
	default:
		return Internal(msg)
	}
}
