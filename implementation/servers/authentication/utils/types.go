package utils

var Users = []User{}

type User struct {
	Username string
	Password string
}
type FormData struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type FingerprintAuthenticator interface {
	CheckFingerprint(username string) bool
}

func (m *MockFingerprintAuthenticator) CheckFingerprint(username string) bool {
	// Placeholder implementation, always returns true
	return true
}

type MockFingerprintAuthenticator struct{}
