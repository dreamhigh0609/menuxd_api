package picture

import (
	"strings"
	"testing"
)

func TestGetExtension(t *testing.T) {
	got := GetExtension("sd4dm4dn4d4.jpg")
	expected := "jpg"

	if got != expected {
		t.Errorf("Got %s, expected %s", got, expected)
	}
}

func TestPictureSlug(t *testing.T) {
	clientID := "2"
	name := "picture01"

	got := GetSlug(clientID, name)
	expected := name + "-" + clientID

	if !strings.HasPrefix(got, expected) {
		t.Errorf("Got %s, expected %s", got, expected)
	}
}
