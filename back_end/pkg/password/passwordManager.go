package passwordpackage

import (
	commonImport "admin-panel/pkg/common"
	"admin-panel/pkg/resources"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type Argon2ID struct {
	format  string
	version int
	time    uint32
	memory  uint32
	keyLen  uint32
	saltLen uint32
	threads uint8
}

var manager *Argon2ID

func SetUp() {
	manager = &Argon2ID{
		format:  "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		version: argon2.Version,
		time:    1,
		memory:  64 * 1024,
		keyLen:  32,
		saltLen: 16,
		threads: 4,
	}
}

func (a Argon2ID) Hash(plain string) (string, error) {
	salt := make([]byte, a.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(plain), salt, a.time, a.memory, a.threads, a.keyLen)
	return fmt.Sprintf(
			a.format,
			a.version,
			a.memory,
			a.time,
			a.threads,
			base64.RawStdEncoding.EncodeToString(salt),
			base64.RawStdEncoding.EncodeToString(hash),
		),
		nil
}

func (a Argon2ID) Verify(plain, hash string) (bool, error) {
	hashParts := strings.Split(hash, "$")
	if len(hashParts) < 6 {
		return false, resources.ErrInternalServer
	}
	_, err := fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &a.memory, &a.time, &a.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])
	if err != nil {
		return false, err
	}

	hashToCompare := argon2.IDKey([]byte(plain), salt, a.time, a.memory, a.threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1, nil
}

func (a Argon2ID) CheckPassword(password string, hashPassword string) error {
	checkResult, err := a.Verify(password, hashPassword)
	if err != nil {
		return err
	}
	if !checkResult {
		return fmt.Errorf("password is not match. %w", resources.ErrClient)
	}
	return nil
}

func ProvidePasswordManager() commonImport.PasswordManager {
	if manager == nil {
		SetUp()
	}
	return manager
}
