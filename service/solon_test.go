package service

import (
	"encoding/json"
	uuid2 "github.com/google/uuid"
	"github.com/odysseia-greek/plato/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSolonClient(t *testing.T) {
	scheme := "http"
	baseUrl := "somelocalhost.com"
	token := "s.49uwenfke9fue"
	tokenResponse := models.TokenResponse{Token: token}
	postResponse := models.SolonResponse{Created: true}
	config := ClientConfig{
		Solon: OdysseiaApi{
			Scheme: scheme,
			Url:    baseUrl,
		},
	}

	requestBody := models.SolonCreationRequest{
		Role:     "testrole",
		Access:   []string{"test"},
		PodName:  "somepodname",
		Username: "testuser",
	}

	t.Run("Get", func(t *testing.T) {
		codes := []int{
			200,
		}

		r, err := tokenResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)

		uuid := uuid2.New().String()

		resp, err := testClient.Solon().OneTimeToken(uuid)
		assert.Nil(t, err)
		defer resp.Body.Close()

		var sut models.TokenResponse
		err = json.NewDecoder(resp.Body).Decode(&sut)
		assert.Nil(t, err)
		assert.Nil(t, err)
		assert.Equal(t, token, sut.Token)
	})

	t.Run("GetWithError", func(t *testing.T) {
		codes := []int{
			500,
		}

		r, err := tokenResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)
		uuid := uuid2.New().String()
		sut, err := testClient.Solon().OneTimeToken(uuid)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("Post", func(t *testing.T) {
		codes := []int{
			201,
		}

		r, err := postResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)
		uuid := uuid2.New().String()
		resp, err := testClient.Solon().Register(requestBody, uuid)
		assert.Nil(t, err)
		defer resp.Body.Close()

		var sut models.SolonResponse
		err = json.NewDecoder(resp.Body).Decode(&sut)
		assert.Nil(t, err)

		assert.Nil(t, err)
		assert.True(t, sut.Created)
	})

	t.Run("PostWithError", func(t *testing.T) {
		codes := []int{
			500,
		}

		r, err := postResponse.Marshal()
		assert.Nil(t, err)

		responses := []string{
			string(r),
		}

		uuid := uuid2.New().String()
		testClient, err := NewFakeClient(config, codes, responses)
		assert.Nil(t, err)
		sut, err := testClient.Solon().Register(requestBody, uuid)
		assert.NotNil(t, err)
		assert.Nil(t, sut)
		assert.Contains(t, err.Error(), "500")
	})
}
