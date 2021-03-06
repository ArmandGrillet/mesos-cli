// Code generated by goderive DO NOT EDIT.

package httpclient

import (
	http "net/http"
)

// deriveDeepCopy recursively copies the contents of src into dst.
func deriveDeepCopy(dst, src http.Header) {
	for src_key, src_value := range src {
		if src_value == nil {
			dst[src_key] = nil
		}
		if src_value == nil {
			dst[src_key] = nil
		} else {
			if dst[src_key] != nil {
				if len(src_value) > len(dst[src_key]) {
					if cap(dst[src_key]) >= len(src_value) {
						dst[src_key] = (dst[src_key])[:len(src_value)]
					} else {
						dst[src_key] = make([]string, len(src_value))
					}
				} else if len(src_value) < len(dst[src_key]) {
					dst[src_key] = (dst[src_key])[:len(src_value)]
				}
			} else {
				dst[src_key] = make([]string, len(src_value))
			}
			copy(dst[src_key], src_value)
		}
	}
}
