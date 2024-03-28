package encription

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	rand2 "crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	custom "github.com/Zaikoa/rapid/src/handling"
)

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand2.Reader, pub, msg, nil)
	if err != nil {
		return nil, custom.NewError(err.Error())
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand2.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, custom.NewError(err.Error())
	}
	return plaintext, nil
}

// Compresses directory or folder into .tar.xz
func Compress(path string) ([]byte, error) {
	cmd := exec.Command("tar", "-cJ", path)
	var compressedData bytes.Buffer
	cmd.Stdout = &compressedData

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return compressedData.Bytes(), nil
}

// Decompresses directory or folder and returns it to original state
func Decompress(name string, path string, file string) error {
	temp := filepath.Join(os.TempDir(), file)
	cmd := exec.Command("tar", "--force-local", "-xf", temp)
	current_dir, _ := os.Getwd()

	dir := filepath.Join(current_dir, path)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Creates the private key text file for the user
func CreatePrivateEncryptFile(privateKey *rsa.PrivateKey) error {
	if err := os.WriteFile("supersecretekey.txt", PrivateKeyToBytes(privateKey), 0644); err != nil {
		return err
	}

	return nil
}

/*
Encrypts location given using public key as a string
*/
func RSAEncryptItem(key string, publickey []byte, nonce []byte) ([]byte, []byte, error) {
	publicKey := BytesToPublicKey(publickey)

	encryptedaes, err := EncryptWithPublicKey([]byte(key), publicKey)
	if err != nil {
		return nil, nil, err
	}
	encryptedNonce, err := EncryptWithPublicKey(nonce, publicKey)
	if err != nil {
		return nil, nil, err
	}
	return encryptedaes, encryptedNonce, nil
}

/*
Decrypts file at location given using private key path
*/
func RSADecryptItem(keypath string, aes []byte, nonce []byte) ([]byte, []byte, error) {
	private_key_bytes, _ := os.ReadFile(keypath)
	privateKey := BytesToPrivateKey(private_key_bytes)

	decryptedAes, err := DecryptWithPrivateKey(aes, privateKey)
	if err != nil {
		return nil, nil, err
	}
	decryptedNonce, err := DecryptWithPrivateKey(nonce, privateKey)
	if err != nil {
		return nil, nil, err
	}
	return decryptedAes, decryptedNonce, nil
}

/*
Ecnrypts file at location given using private key path
*/
func AESEncryptionItem(name string, bytes []byte, keyString string) ([]byte, error) {
	key, _ := hex.DecodeString(keyString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, bytes, nil)
	location := filepath.Join(os.TempDir(), name)
	os.WriteFile(location, ciphertext, 0644)
	return nonce, nil
}

/*
Ecnrypts file at location given using private key path
*/
func AESDecryptItem(name string, keyString []byte, nonce []byte) error {
	key, _ := hex.DecodeString(string(keyString))
	location := filepath.Join(os.TempDir(), name)

	bytes, err := os.ReadFile(location)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	file, _ := aesgcm.Open(nil, nonce, bytes, nil)
	err = os.WriteFile(location, file, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
