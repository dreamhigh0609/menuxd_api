package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparePass(t *testing.T) {
	n := New()
	n.Password = "Dk123456"
	err := n.PreparePass()
	if err != nil {
		t.Errorf("Expected not errors, got %s", err.Error())
	}
	if ok := n.ComparePass("Dk123456"); !ok {
		t.Errorf("Expected true, got %t", ok)
	}
}

func TestComparePassFail(t *testing.T) {
	n := New()
	n.Password = "Dk123489"
	err := n.PreparePass()
	if err != nil {
		assert.Nil(t, err)
	}

	assert.False(t, n.ComparePass("Dk123456"))

}

func TestConfirmPass(t *testing.T) {
	n := New()
	n.Password = "Dk123456"
	n.ConfirmPassword = "Dk123456"

	if ok := n.ConfirmPass(); !ok {
		t.Errorf("Expected true, got %t", ok)
	}
}

func TestFailConfirmPass(t *testing.T) {
	n := New()
	n.Password = "Dk123456"
	n.ConfirmPassword = "Dk654321"

	if ok := n.ConfirmPass(); ok {
		t.Errorf("Expected false, got %t", ok)
	}
}

func TestValidateEmail(t *testing.T) {
	users := []User{
		{Email: "abcd@gmail-yahoo.com"},
		{Email: "orlando@gmail.com"},
		{Email: "client@example.com"},
	}

	assert := assert.New(t)

	for _, u := range users {
		assert.True(u.ValidateEmail())
	}
}

func TestValidateEmailFail(t *testing.T) {
	users := []User{
		{Email: "@gmail-yahoo.com"},
		{Email: "orlando-gmail"},
		{Email: "orlando daniel"},
	}

	assert := assert.New(t)

	for _, u := range users {
		assert.False(u.ValidateEmail())
	}
}

func TestCleanPass(t *testing.T) {
	u := User{
		Password:        "Xd123456",
		ConfirmPassword: "Xd123456",
		HashPassword:    "dweded32rd43d4dwerdwg4g",
		OldPassword:     "123456",
	}

	u.CleanPass()

	assert.Equal(t, "", u.Password)
	assert.Equal(t, "", u.ConfirmPassword)
	assert.Equal(t, "", u.HashPassword)
	assert.Equal(t, "", u.OldPassword)
}
