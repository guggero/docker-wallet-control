package main

import (
  "testing"
  "github.com/guggero/docker-wallet-control/util"
  "github.com/stretchr/testify/assert"
  "net/http"
)

func setupTestdata() {
  appConfig = &util.Configuration{
    User: []util.UserConfig{
      // correct hash, password is 'test'
      {Username: "test", Password: "1bc1a361f17092bc7af4b2f82bf9194ea9ee2ca49eb2e53e39f555bc1eeaed74", Salt: "salt"},
      // incorrect hash
      {Username: "test2", Password: "aaaabbbbbccccc", Salt: "salt"},
    },
  }
}

func TestAuthenticateUserCorrect(t *testing.T) {
  // given
  setupTestdata()

  // when
  result := authenticateUser("test", "test")

  // then
  assert.Equal(t, appConfig.User[0].Username, result)
}

func TestAuthenticateUserIncorrect(t *testing.T) {
  // given
  setupTestdata()

  // when
  result := authenticateUser("test2", "test2")

  // then
  assert.Equal(t, "", result)
}

func TestGetAuthenticatedUser(t *testing.T) {
  // given
  setupTestdata()
  request := http.Request{
    Header: http.Header{
      // test:test in base64
      "Authorization": {"Basic dGVzdDp0ZXN0"},
    },
  }

  // when
  result := getAuthenticatedUser(&request);

  // then
  assert.Equal(t, appConfig.User[0].Username, result)
}
