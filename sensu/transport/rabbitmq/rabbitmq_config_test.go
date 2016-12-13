package rabbitmq

import (
	"errors"
	"fmt"
	"net/url"
	"testing"
)

func validateStringParameter(
	actual string,
	expected string,
	parameterName string,
	t *testing.T) {

	if actual != expected {
		t.Errorf(
			"Expected %s to be \"%s\" but got \"%s\" instead!",
			parameterName,
			expected,
			actual,
		)
	}
}

func validateError(actual error, expected error, t *testing.T) {
	if actual == nil {
		t.Errorf("Expected error to be set but got nil instead!")
	}

	if actual.Error() != expected.Error() {
		t.Errorf(
			"Expected error to be \"%s\" but got \"%s\" instead!",
			expected,
			actual,
		)
	}
}

func TestNewTransportConfigHappyFlow(t *testing.T) {
	expectedHost := "localhost"
	expectedPort := "5673"
	expectedVhost := "/sensu"
	expectedUser := "test_user"
	expectedPassword := "test_password"

	config, err := NewTransportConfig(fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		expectedUser,
		expectedPassword,
		expectedHost,
		expectedPort,
		url.QueryEscape(expectedVhost),
	))

	if err != nil {
		t.Errorf(
			"Expected error to be nil but got \"%s\" instead!",
			err,
		)
	}

	if config == nil {
		t.Errorf("Expected config to not be nil")
	}

	validateStringParameter(config.Host, expectedHost, "host", t)
	validateStringParameter(string(config.Port), expectedPort, "port", t)
	validateStringParameter(config.Vhost, expectedVhost, "vhost", t)
	validateStringParameter(config.User, expectedUser, "user", t)
	validateStringParameter(config.Password, expectedPassword, "password", t)
}

var transportConfigErrorTestScenarios = []struct {
	uri           string
	expectedError error
}{
	{
		"",
		errors.New("Failed to determine the port for host: "),
	},
	{
		"://",
		errors.New("parse ://: missing protocol scheme"),
	},
	{
		"amqp://example.com/",
		errors.New("Failed to determine the port for host: example.com"),
	},
	{
		"amqp://example.com::",
		errors.New("too many colons in address example.com::"),
	},
	{
		"amqp://example.com:5672",
		errNoUserInURI,
	},
}

func TestNewTransportConfigErrorScenarios(t *testing.T) {
	for _, test := range transportConfigErrorTestScenarios {
		config, err := NewTransportConfig(test.uri)

		validateError(err, test.expectedError, t)

		if config != nil {
			t.Errorf("Expected config to be nil")
		}
	}
}
