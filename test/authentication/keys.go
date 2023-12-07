package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
)

// This is an only test key, to a only test service.
// NEVER add a private key to any code or GitHub repository.
const privateKeyPem = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAvpa5w4Vm8aOVCnI46O9f7Ixp3hir1TGgdo6p25ZHR/plk4Nd
tQI04TT2Uo7iCQD1FSJat7hYu0HYwsG5qMh1fZwi+GFf3Yqfxy5kpgUsatvC1wZg
lcccmV+qpL+Nj5bsaV7HrTyRPkru1twSXnOcAcZUesQdDo56otJfTDEvbdGBetbk
apIkcjoWZHy39KPg4aWMlJ7GLpHAEvEVTh/6Impu/lUSZMy/V9D1IdgjKFmu0BF1
nMLdxTAwZVU7YNiOxovQU6Hw/UlZRyVlVubzhSp5HA9dib/n0AaIv97VelgDBGo7
OcWlISOb0kIz0VvZtgfXQItqz0tyvSEJAWWytQIDAQABAoIBAQCJONiTL++It3DN
zqOvAvqbxBVNdZHytAKGmf0uPysfZefQp9rGQsp0A7/+fSW9udS73LpYYQByAtzg
jL7yCTKet9Zt4x400mRr8rlA16S9Y3ELhGnyLSQoQbsJV2nDIrUgwL8fueYRZb5F
MLqBCtgeZv/YTG9nVapypOk8YIV3mQx+M0J9PgIjkEZmu7dN0no6Rl3lRnzkW+0Y
Wu+AkSnhGATITeO17p/ONUChNMWaZ+kVABesDaPQeONmzWRHPreR+IbEhZSQmm4k
4l+qkfojIfBL+lT6sslmjX/c9enOEuINQiAJHRb7R5cM4pqjGJ3yOUnMkC5cHRRG
jYMt0yfBAoGBAN3qMH4Q5UIKTdcfMh5V7z1z6+gN8eBviaWCGk8YToEygCglDl7C
ybqVAXEwd0qO0SqvsXIdv4+9OR67oIFhl6sDeHz+G2cuYSUXruvCnUWLUTfnRoqF
9VHi6cWz5N7nfp6Dy52pyQy+e6PUoB4+EZBehfhxON69FZzcMMaTc5YtAoGBANvc
x5lg2E0isXrmiCUY1447sIFnhWDdPxC+kERwV9sJZ2hMs7fl7KKUXS1ichKUURUx
75jVTRygm839KkFViFvSLdn2rYpkpt4tAfqklmPZVvDxtokjSAIRUPDHpqL9rmQD
bGyQnuu8/unT+rGIo4yiFpFhVxEe2rhgr0a4FiupAoGBAJQePxW18z+cHw6KDOrA
kvmiiQAPZrVV3TryVtsaLzP+4Blremb3fqwhzp+dKNJD9wqV0EuJ3ZV0SE7iDySs
Xg5QN7i95s5832xhnWhRMqX7ck9s9+F3viFU4pIKG6ZIP3RQJbTrYX03GtFkFyd4
aELDRIpqD/pjnKxhL9ErFAhVAoGAApwqWm3F45SH2telwhr7ZBrdS4v5D19RAlfg
yo8y28zOx3Qxpfs6xetQ99r1U7cjB0diesP9eFuHvfhFaiUjy0NBfBbrlHsBaB3M
qjcN+f14hL+51QLwNeYSuekE12Z/jXxk6x0EZfQGaqwzi6v9lQvPjMZFDFT7b7jm
G8bPrJECgYAB4jWa7ZIbXoSrPI++rO4KjIwk7ZOHdc+loQCILEZ2tAsRk13W1K2n
TWGZwV3znT3ivHpyi1E3am4muGIZSur3xU9Ugn41jQHky1Z5SMEYew4cdm3PSZlb
CdPf87BbBtvFw5DM9FiDUcSkdn2pwawWHQgLYW/SByk4zJ2cZluQaA==
-----END RSA PRIVATE KEY-----
`

var privateKey *rsa.PrivateKey
var PublicKey rsa.PublicKey

func init() {
	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		log.Fatal("invalid private key")
	}

	var err error
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	PublicKey = privateKey.PublicKey
}
