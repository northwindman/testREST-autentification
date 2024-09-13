package tokens

import (
	"errors"
	"github.com/northwindman/testREST-autentification/internal/http-server/handlers/user/refresh"
	myjwt "github.com/northwindman/testREST-autentification/internal/lib/tokens/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// Mock for refresh.New
type RefreshMock struct {
	mock.Mock
}

func (m *RefreshMock) New(length int) (string, error) {
	args := m.Called(length)
	return args.String(0), args.Error(1)
}

// Mock для myjwt.New
type JwtMock struct {
	mock.Mock
}

func (m *JwtMock) New(ip string, email string, secret string) (string, error) {
	args := m.Called(ip, email, secret)
	return args.String(0), args.Error(1)
}

func TestGenTokens_Success(t *testing.T) {
	refreshMock := new(RefreshMock)
	jwtMock := new(JwtMock)

	expectedRefreshToken := "refresh-token"
	expectedAccessToken := "access-token"

	refreshMock.On("New", 32).Return(expectedRefreshToken, nil)
	jwtMock.On("New", "127.0.0.1", "test@example.com", "secret").Return(expectedAccessToken, nil)

	originalNewRefresh := refresh.New
	originalNewJwt := myjwt.New
	defer func() {
		refresh.New = originalNewRefresh
		myjwt.New = originalNewJwt
	}()
	refresh.New = refreshMock.New
	myjwt.New = jwtMock.New

	tokens, err := GenTokens("127.0.0.1", "test@example.com", "secret", 32)
	require.NoError(t, err)

	assert.Equal(t, expectedRefreshToken, tokens.RefreshToken)
	assert.Equal(t, expectedAccessToken, tokens.AccessToken)

	refreshMock.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
}

func TestGenTokens_RefreshTokenError(t *testing.T) {
	refreshMock := new(RefreshMock)
	jwtMock := new(JwtMock)

	expectedError := errors.New("failed to generate refresh token")

	refreshMock.On("New", 32).Return("", expectedError)
	jwtMock.On("New", "127.0.0.1", "test@example.com", "secret").Return("", nil)

	originalNewRefresh := refresh.New
	originalNewJwt := myjwt.New
	defer func() {
		refresh.New = originalNewRefresh
		myjwt.New = originalNewJwt
	}()
	refresh.New = refreshMock.New
	myjwt.New = jwtMock.New

	tokens, err := GenTokens("127.0.0.1", "test@example.com", "secret", 32)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, tokens.RefreshToken)
	assert.Empty(t, tokens.AccessToken)

	refreshMock.AssertExpectations(t)
	jwtMock.AssertExpectations(t)
}
