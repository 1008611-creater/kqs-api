package handler

import (
	"fmt"
	"io"
)

const maxOAuthResponseBodyBytes int64 = 1 << 20

func readLimitedOAuthResponse(body io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(body, maxOAuthResponseBodyBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxOAuthResponseBodyBytes {
		return nil, fmt.Errorf("oauth response body exceeds %d bytes", maxOAuthResponseBodyBytes)
	}
	return data, nil
}
