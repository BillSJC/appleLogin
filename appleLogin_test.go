package appleLogin

import "testing"

//a random p8 cert to go test
const certStr = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgusZ/Y029MmQ4mXWn
fnzXUMI/DgtJIJdvG3cZtOsL3pmgCgYIKoZIzj0DAQehRANCAASQloEXsIF31S59
n5/2YdbDaijlx2eIyIfkv7tre3GxgG8NILwvNCrg6L9Tm9JkVjsLucwXcQ+ezINf
YJBJn/t2
-----END PRIVATE KEY-----`

func TestAppleConfig_LoadP8CertByByte(t *testing.T) {

	a := InitAppleConfig("", "", "", "")
	err := a.LoadP8CertByByte([]byte(certStr))
	if err != nil {
		t.Fatal(err)
	}
}
