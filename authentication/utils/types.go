package utils

var Users = []User{}

type User struct {
	ID         int
	Username   string
	Password   string
	TrustLevel int
	Groups     []string
}

type FingerprintAuthenticator interface {
	CheckFingerprint(username string) bool
}

func (m *MockFingerprintAuthenticator) CheckFingerprint(username string) bool {
	// Placeholder implementation, always returns true
	return true
}

type MockFingerprintAuthenticator struct{}
