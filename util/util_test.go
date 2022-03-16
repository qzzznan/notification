package util

import (
	"github.com/golang-jwt/jwt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"testing"
)

func TestJWT(t *testing.T) {
	tokenStr := `eyJraWQiOiJXNldjT0tCIiwiYWxnIjoiUlMyNTYifQ.eyJpc3MiOiJodHRwczovL2FwcGxlaWQuYXBwbGUuY29tIiwiYXVkIjoiY29tLnB1c2hkZWVyLnNlbGYuaW9zIiwiZXhwIjoxNjQ3NDgyMTkwLCJpYXQiOjE2NDczOTU3OTAsInN1YiI6IjAwMDQ0OC4xNDlkNGE1MzAzZTU0ZmRjYWVlZDEzNmVkMTg5YzdlNC4wOTU1IiwiY19oYXNoIjoiTjNrS2IzOENVd01JbG9URi1hMEZMZyIsImVtYWlsIjoidzVyZDJ0YzQ4Y0Bwcml2YXRlcmVsYXkuYXBwbGVpZC5jb20iLCJlbWFpbF92ZXJpZmllZCI6InRydWUiLCJpc19wcml2YXRlX2VtYWlsIjoidHJ1ZSIsImF1dGhfdGltZSI6MTY0NzM5NTc5MCwibm9uY2Vfc3VwcG9ydGVkIjp0cnVlfQ.IgwkstvcKQY2SAZHaCm7k-9M0p0wYSY_o8Vr6k2WP1bNR0dUJfUccaEniXFCfti3vhjrl3wGxtFTEl0YTSHMQ4YsWVJJS571E0604GR8wN-9_Yp_z68ud6lwv35ECCDuBrzpC67-Lt10m0ACDYnsJnxLSh6KHc9hjki0KGeTp2XytZOQ1JQM_kQjWKNOj1CY9ZsXgPS-jNNhvAAkv2ntYX24uh2PVhPPnoOnyYtfzxj9MKbgfbCuKamOl6q-jgWWzeQnqShLUg_FIoM0PV-wd3pxo_xFSG8mFk1TEmUalYF8M6VDNseN-CAOpx9Wl2CcwfWLrxVw2GhvEFuuFIdScw`
	token, err := jwt.Parse(tokenStr, nil)

	t.Log(token.Header)
	t.Log(token.Claims)
	t.Log(token.Method)
	t.Log(err)

	resp, err := http.Get("https://appleid.apple.com/auth/keys")
	if err != nil {
		t.Fatal(err)
	}
	keys := make(map[string][]map[string]string)
	err = jsoniter.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(keys)
	for _, key := range keys {
		for _, v := range key {
			if v["kid"] == token.Header["kid"] {
				t.Log("find key")
			}
			t.Log(v)
		}
	}
}
