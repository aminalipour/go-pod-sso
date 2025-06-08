package pkg

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/big"
)

type RSAKey struct {
	Modulus  string `xml:"Modulus"`
	Exponent string `xml:"Exponent"`
	D        string `xml:"D"`
	P        string `xml:"P"`
	Q        string `xml:"Q"`
	DP       string `xml:"DP"`
	DQ       string `xml:"DQ"`
	InverseQ string `xml:"InverseQ"`
}

func GetSignatureFromString(privateKey string, payLoad string) (string, error) {
	// convert payload to byte
	payloadByte := []byte(payLoad)

	private_key, err := load_private_key_from_xml([]byte(privateKey))
	if err != nil {
		return "", err
	}

	// sing the private key with payload
	signature, err := sign_payload(private_key, payloadByte)
	if err != nil {
		return "", err
	}
	return signature, nil
}

func GetSignatureFromFile(privateKeyPath string, payLoad string) (string, error) {
	// convert payload to byte
	payloadByte := []byte(payLoad)
	// read the xml from file
	xml_private_key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}
	// generate the private key from xml
	private_key, err := load_private_key_from_xml(xml_private_key)
	if err != nil {
		return "", err
	}
	// sing the private key with payload
	signature, err := sign_payload(private_key, payloadByte)
	if err != nil {
		return "", err
	}
	return signature, nil
}

func base64_to_bigint(b64_str string) *big.Int {
	data, _ := base64.StdEncoding.DecodeString(b64_str)
	bi := new(big.Int)
	bi.SetBytes(data)
	return bi
}

func load_private_key_from_xml(xml_key []byte) (*rsa.PrivateKey, error) {
	var rsa_key RSAKey
	if err := xml.Unmarshal(xml_key, &rsa_key); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %v", err)
	}
	expBytes, err := base64.StdEncoding.DecodeString(rsa_key.Exponent)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %v", err)
	}

	e := 0
	for _, b := range expBytes {
		e = e<<8 + int(b)
	}

	private_key := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: base64_to_bigint(rsa_key.Modulus),
			E: e,
		},
		D:      base64_to_bigint(rsa_key.D),
		Primes: []*big.Int{base64_to_bigint(rsa_key.P), base64_to_bigint(rsa_key.Q)},
	}
	private_key.Precompute()
	return private_key, nil
}

func sign_payload(private_key *rsa.PrivateKey, payload []byte) (string, error) {
	hashed := sha256.Sum256(payload)
	signature, err := rsa.SignPKCS1v15(rand.Reader, private_key, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
