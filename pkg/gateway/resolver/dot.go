package resolver

import (
	"crypto/tls"
	"crypto/x509"
	"dns-proxy/pkg/domain/proxy"
	"errors"
	"log"
	"strconv"
	"time"
)

const CA = `-----BEGIN CERTIFICATE-----
MIIEQzCCAyugAwIBAgIQCidf5wTW7ssj1c1bSxpOBDANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBD
QTAeFw0yMDA5MjMwMDAwMDBaFw0zMDA5MjIyMzU5NTlaMFYxCzAJBgNVBAYTAlVT
MRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxMDAuBgNVBAMTJ0RpZ2lDZXJ0IFRMUyBI
eWJyaWQgRUNDIFNIQTM4NCAyMDIwIENBMTB2MBAGByqGSM49AgEGBSuBBAAiA2IA
BMEbxppbmNmkKaDp1AS12+umsmxVwP/tmMZJLwYnUcu/cMEFesOxnYeJuq20ExfJ
qLSDyLiQ0cx0NTY8g3KwtdD3ImnI8YDEe0CPz2iHJlw5ifFNkU3aiYvkA8ND5b8v
c6OCAa4wggGqMB0GA1UdDgQWBBQKvAgpF4ylOW16Ds4zxy6z7fvDejAfBgNVHSME
GDAWgBQD3lA1VtFMu2bwo+IbG8OXsj3RVTAOBgNVHQ8BAf8EBAMCAYYwHQYDVR0l
BBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBIGA1UdEwEB/wQIMAYBAf8CAQAwdgYI
KwYBBQUHAQEEajBoMCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdpY2VydC5j
b20wQAYIKwYBBQUHMAKGNGh0dHA6Ly9jYWNlcnRzLmRpZ2ljZXJ0LmNvbS9EaWdp
Q2VydEdsb2JhbFJvb3RDQS5jcnQwewYDVR0fBHQwcjA3oDWgM4YxaHR0cDovL2Ny
bDMuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNybDA3oDWgM4Yx
aHR0cDovL2NybDQuZGlnaWNlcnQuY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdENBLmNy
bDAwBgNVHSAEKTAnMAcGBWeBDAEBMAgGBmeBDAECATAIBgZngQwBAgIwCAYGZ4EM
AQIDMA0GCSqGSIb3DQEBDAUAA4IBAQDeOpcbhb17jApY4+PwCwYAeq9EYyp/3YFt
ERim+vc4YLGwOWK9uHsu8AjJkltz32WQt960V6zALxyZZ02LXvIBoa33llPN1d9R
JzcGRvJvPDGJLEoWKRGC5+23QhST4Nlg+j8cZMsywzEXJNmvPlVv/w+AbxsBCMqk
BGPI2lNM8hkmxPad31z6n58SXqJdH/bYF462YvgdgbYKOytobPAyTgr3mYI5sUje
CzqJx1+NLyc8nAK8Ib2HxnC+IrrWzfRLvVNve8KaN9EtBH7TuMwNW4SpDCmGr6fY
1h3tDjHhkTb9PA36zoaJzu0cIw265vZt6hCmYWJC+/j+fgZwcPwL
-----END CERTIFICATE-----
------BEGIN CERTIFICATE-----
-MIIFhjCCBQ2gAwIBAgIQBQdvZtEbaSJWzKzVRv/sUzAKBggqhkjOPQQDAzBWMQsw
-CQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMTAwLgYDVQQDEydEaWdp
-Q2VydCBUTFMgSHlicmlkIEVDQyBTSEEzODQgMjAyMCBDQTEwHhcNMjEwMTExMDAw
-MDAwWhcNMjIwMTE4MjM1OTU5WjByMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2Fs
-aWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEZMBcGA1UEChMQQ2xvdWRm
-bGFyZSwgSW5jLjEbMBkGA1UEAxMSY2xvdWRmbGFyZS1kbnMuY29tMFkwEwYHKoZI
-zj0CAQYIKoZIzj0DAQcDQgAEF60f6DWvcNONnJ5k/UceW5cMCtEQqCYyETZmTRKZ
-w+Exu/UhY3PdpcHBoPBtpMRe4cLb2vkNNIAa97ngOvLVdKOCA58wggObMB8GA1Ud
-IwQYMBaAFAq8CCkXjKU5bXoOzjPHLrPt+8N6MB0GA1UdDgQWBBThtvwG+bmLBfTB
-4kibArkLwbU9eTCBpgYDVR0RBIGeMIGbghJjbG91ZGZsYXJlLWRucy5jb22CFCou
-Y2xvdWRmbGFyZS1kbnMuY29tgg9vbmUub25lLm9uZS5vbmWHBAEBAQGHBAEAAAGH
-BKKfJAGHBKKfLgGHECYGRwBHAAAAAAAAAAAAERGHECYGRwBHAAAAAAAAAAAAEAGH
-ECYGRwBHAAAAAAAAAAAAAGSHECYGRwBHAAAAAAAAAAAAZAAwDgYDVR0PAQH/BAQD
-AgeAMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjCBlwYDVR0fBIGPMIGM
-MESgQqBAhj5odHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRUTFNIeWJy
-aWRFQ0NTSEEzODQyMDIwQ0ExLmNybDBEoEKgQIY+aHR0cDovL2NybDQuZGlnaWNl
-cnQuY29tL0RpZ2lDZXJ0VExTSHlicmlkRUNDU0hBMzg0MjAyMENBMS5jcmwwSwYD
-VR0gBEQwQjA2BglghkgBhv1sAQEwKTAnBggrBgEFBQcCARYbaHR0cDovL3d3dy5k
-aWdpY2VydC5jb20vQ1BTMAgGBmeBDAECAjCBgwYIKwYBBQUHAQEEdzB1MCQGCCsG
-AQUFBzABhhhodHRwOi8vb2NzcC5kaWdpY2VydC5jb20wTQYIKwYBBQUHMAKGQWh0
-dHA6Ly9jYWNlcnRzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydFRMU0h5YnJpZEVDQ1NI
-QTM4NDIwMjBDQTEuY3J0MAwGA1UdEwEB/wQCMAAwggEEBgorBgEEAdZ5AgQCBIH1
-BIHyAPAAdgApeb7wnjk5IfBWc59jpXflvld9nGAK+PlNXSZcJV3HhAAAAXby6BKo
-AAAEAwBHMEUCIQDRsvaM+FOVneTUUwY0ggKKCuqKp7wnHvtWHtEUZB+uZwIgJbGG
-3Rsq548BxED2wxZ4q2G/9jo0/EeIEwdl9GC7NEIAdgAiRUUHWVUkVpY/oS/x922G
-4CMmY63AS39dxoNcbuIPAgAAAXby6BMPAAAEAwBHMEUCIQCV3RpnSizsrJ1vi/48
-/qT1PoclZYI3N51mveRdD2gkWQIgdWX+MLuAa8ziuKGIlqjoAiaOvs/4IfqthaAN
-h6HW8TQwCgYIKoZIzj0EAwMDZwAwZAIwJMLPbL32rtHJ1R9KdC48PdHAPtzXG9OU
-cVv+pYYWJoIBItMKbvyYtdLiueUHaXeWAjBFe2+Cpn22YsMxhdW1NV1PTISIrBoA
-PQyEQNywp8ocEycVHjf5RsOu2f35uSOLfyo=
 -----END CERTIFICATE-----`

type dot struct {
	ip          string
	port        int
	rootCert    string
	readTimeOut uint
}

func NewDNSOverTlsResolver(ip string, port int, rto uint) proxy.Resolver {
	return &dot{ip, port, CA, rto}
}

func (dot *dot) GetTLSConnection() (*tls.Conn, error) {
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(dot.rootCert)) {
		log.Println("Fail to parse rootCert")
		return nil, errors.New("Fail to parse rootCert")
	}
	dnsTLSConn, err := tls.Dial("tcp", dot.ip+":"+strconv.Itoa(dot.port), &tls.Config{
		RootCAs: roots,
	})
	if err != nil {
		log.Println("Error connecting to CloudFlare")
		return nil, err
	}
	_ = dnsTLSConn.SetReadDeadline(time.Now().Add(time.Duration(dot.readTimeOut) * time.Millisecond))
	return dnsTLSConn, nil
}

func (dot *dot) Solve(um proxy.UnsolvedMsg) (proxy.SolvedMsg, error) {
	// Levanto una conexión con CF
	conn, err := dot.GetTLSConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	_, e := conn.Write(um)
	if e != nil {
		log.Printf("%v", e)
	}
	var reply [2045]byte
	n, err := conn.Read(reply[:])
	if err != nil {
		log.Printf("Could read response from CloudFlare: %v \n", err)
		return nil, err
	} else {
		log.Println("Succesfuly fullfiled the request")
	}

	return reply[:n], nil
}
