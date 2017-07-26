// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type dummyTokenSource struct {
	oauth2.TokenSource
}

func TestTokenSource(t *testing.T) {
	ctx := context.Background()

	// Pass in a TokenSource, get it back.
	ts := &dummyTokenSource{}
	ds := &DialSettings{TokenSource: ts}
	got, err := Creds(ctx, ds)
	if err != nil {
		t.Fatal(err)
	}
	want := &google.DefaultCredentials{TokenSource: ts}
	if !cmp.Equal(got, want) {
		t.Error("did not get the same TokenSource back")
	}

	// Load a valid JSON file. No way to really test the contents; we just
	// verify that there is no error.
	ds = &DialSettings{ServiceAccountJSONFilename: "service-account.json"}
	if _, err := Creds(ctx, ds); err != nil {
		t.Errorf("got %v, wanted no error", err)
	}

	// If both a file and TokenSource are passed, the file takes precedence
	// (existing behavior).
	// TODO(jba): make this an error?
	ds = &DialSettings{
		TokenSource:                ts,
		ServiceAccountJSONFilename: "service-account.json",
	}
	got, err = Creds(ctx, ds)
	if err != nil {
		t.Fatal(err)
	}
	if cmp.Equal(got, want) {
		t.Error("got the same TokenSource back, wanted one from the JSON file")
	}

	// TODO(jba): find a way to test the call to google.DefaultTokenSource.
}