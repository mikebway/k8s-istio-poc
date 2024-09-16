// Most of the contents of this file was lifted from the Golang http package, cookie.go file.
//
// See https://github.com/golang/go/blob/master/src/net/http/cookie.go
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"golang.org/x/net/http/httpguts"
	"net/textproto"
	"strings"
)

// RequestCookie represents all the cookie values that we get in the request, just the name and the value.
// It's not worth passing around mostly empty http.Cookie structures so we use this.
type RequestCookie struct {
	Name  string
	Value string
}

// readCookies parses all "Cookie" values from the supplied header string value.
//
// This function is heavily derived lifted from the Golang http package, cookie.go file.
// See https://github.com/golang/go/blob/master/src/net/http/cookie.go
func readCookies(rawCookies string) map[string]*RequestCookie {

	// If there are no cookies, then there are no cookies!
	rawCookies = textproto.TrimString(rawCookies)
	if len(rawCookies) == 0 {
		return nil
	}

	// Gather the cookie structures in a slice of cookies
	cookies := make(map[string]*RequestCookie, strings.Count(rawCookies, ";"))
	var part string
	for len(rawCookies) > 0 { // continue since we have rest
		part, rawCookies, _ = strings.Cut(rawCookies, ";")
		part = textproto.TrimString(part)
		if part == "" {
			continue
		}
		name, val, _ := strings.Cut(part, "=")
		name = textproto.TrimString(name)
		if !isCookieNameValid(name) {
			continue
		}
		val, ok := parseCookieValue(val, true)
		if !ok {
			continue
		}
		cookies[name] = &RequestCookie{Name: name, Value: val}
	}
	return cookies
}

// parseCookieValue was copied, unchanged from https://github.com/golang/go/blob/master/src/net/http/cookie.go
//
// The original was uncommented.
func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	// Strip the quotes, if present.
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	for i := 0; i < len(raw); i++ {
		if !validCookieValueByte(raw[i]) {
			return "", false
		}
	}
	return raw, true
}

// validCookieValueByte was copied, unchanged from https://github.com/golang/go/blob/master/src/net/http/cookie.go
//
// The original was uncommented.
func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}

// isCookieNameValid was copied, unchanged from https://github.com/golang/go/blob/master/src/net/http/cookie.go
//
// The original was uncommented.
func isCookieNameValid(raw string) bool {
	if raw == "" {
		return false
	}
	return strings.IndexFunc(raw, isNotToken) < 0
}

// isNotToken was copied, unchanged from https://github.com/golang/go/blob/master/src/net/http/cookie.go
//
// The original was uncommented.
func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}
