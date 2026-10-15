/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"net/http"
	"testing"

	gitlab "gitlab.com/gitlab-org/api/client-go/v3"
)

func TestAreStringSlicesEqual(t *testing.T) {
	cases := map[string]struct {
		a    []string
		b    []string
		want bool
	}{
		"BothNil": {
			a: nil, b: nil, want: true,
		},
		"BothEmpty": {
			a: []string{}, b: []string{}, want: true,
		},
		"NilAndEmpty": {
			a: nil, b: []string{}, want: true,
		},
		"Equal": {
			a: []string{"a", "b"}, b: []string{"a", "b"}, want: true,
		},
		"EqualDifferentOrder": {
			a: []string{"b", "a"}, b: []string{"a", "b"}, want: true,
		},
		"DifferentLength": {
			a: []string{"a", "b"}, b: []string{"a"}, want: false,
		},
		"DifferentElements": {
			a: []string{"a", "b"}, b: []string{"a", "c"}, want: false,
		},
		"DuplicatesInA": {
			a: []string{"a", "a"}, b: []string{"a", "b"}, want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := AreStringSlicesEqual(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("AreStringSlicesEqual(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestIsResponseNotFound(t *testing.T) {
	response := func(status int, contentType string) *gitlab.Response {
		header := http.Header{}
		if contentType != "" {
			header.Set("Content-Type", contentType)
		}
		return &gitlab.Response{Response: &http.Response{StatusCode: status, Header: header}}
	}

	cases := map[string]struct {
		res  *gitlab.Response
		want bool
	}{
		"NilResponse": {
			res: nil, want: false,
		},
		"NilHTTPResponse": {
			res: &gitlab.Response{}, want: false,
		},
		"GitLabNotFound": {
			res: response(http.StatusNotFound, "application/json"), want: true,
		},
		"GitLabNotFoundWithCharset": {
			res: response(http.StatusNotFound, "application/json; charset=utf-8"), want: true,
		},
		"ProxyHTMLNotFound": {
			res: response(http.StatusNotFound, "text/html"), want: false,
		},
		"ProxyPlainTextNotFound": {
			res: response(http.StatusNotFound, "text/plain; charset=utf-8"), want: false,
		},
		"NotFoundWithoutContentType": {
			res: response(http.StatusNotFound, ""), want: false,
		},
		"GitLabServerError": {
			res: response(http.StatusInternalServerError, "application/json"), want: false,
		},
		"BadGateway": {
			res: response(http.StatusBadGateway, "text/html"), want: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := IsResponseNotFound(tc.res); got != tc.want {
				t.Errorf("IsResponseNotFound() = %v, want %v", got, tc.want)
			}
		})
	}
}
