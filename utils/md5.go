package utils

import (
	"fmt"
	md5simd "github.com/minio/md5-simd"
	"os"
)

func GetFileMD5(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	server := md5simd.NewServer()
	defer server.Close()

	// Create hashing object (conforming to hash.Hash)
	md5Hash := server.NewHash()
	defer md5Hash.Close()

	// Write one (or more) blocks
	md5Hash.Write(content)

	// Return digest
	return fmt.Sprintf("%x", md5Hash.Sum([]byte{})), nil
}

func GetMd5FromBytes(data []byte) (string, error) {
	server := md5simd.NewServer()
	defer server.Close()

	// Create hashing object (conforming to hash.Hash)
	md5Hash := server.NewHash()
	defer md5Hash.Close()

	// Write one (or more) blocks
	_, err := md5Hash.Write(data)
	if err != nil {
		return "", err
	}
	// Return digest
	return fmt.Sprintf("%x", md5Hash.Sum([]byte{})), nil

}
