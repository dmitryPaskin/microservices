package responder

import "google.golang.org/grpc/status"

func extractGrpcErrorMessage(err error) string {
	if st, ok := status.FromError(err); ok {
		return st.Message()
	}

	return err.Error()
}
