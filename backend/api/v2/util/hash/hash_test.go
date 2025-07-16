package hash

import (
	"testing"
)

func TestSha256(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"password", "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"},
		{"123456", "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}

	for _, test := range tests {
		actual := Sha256(test.input)
		if actual != test.expected {
			t.Errorf("Sha256(%q) = %q, expected %q", test.input, actual, test.expected)
		}
	}
}

func TestCryptSha512(t *testing.T) {
	tests := []struct {
		password string
		salt     string
		want string
	}{
		{
			"password",
			"somesalt",
			"$6$rounds=5000$somesalt$A7P/0Yfu8RprY88D5T1n.xKT749BOn/IXBvmR1gXZzU7imsoTfZhCQ1916CB7WNX9eOOeSmBmmMrl5fQn9LAP1",
		},
		{
			"123456",
			"anothersalt",
			"$6$rounds=5000$anothersalt$ZffCt8y5Hl8gLYS79/rhyT76C12kNhuOvkFR8Ll0RXcjQz2Nzuh3VUdT//e21HYfH6fP9btOp2aG22O3S7q1z/"},
	}

	for _, test := range tests {
		got, err := CryptSha512(test.password, test.salt)
		if err != nil {
			t.Errorf("CryptSha512(%q, %q) returned an error: %v", test.password, test.salt, err)
		}
		if got != test.want {
			t.Errorf("CryptSha512(%q, %q) = %q, want %q", test.password, test.salt, got, test.want)
		}
	}
}