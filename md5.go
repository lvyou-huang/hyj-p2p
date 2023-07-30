package mymd5

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

func Md5check2(fileName string, md5bytes []byte) (error, bool) {
	file, err := os.Open(fileName)
	if err != nil {
		return err, false
	}
	defer file.Close()

	return md5Check_2(fileName, md5bytes)
}

func md5Check_2(fileName string, md5sum []byte) (error, bool) {
	fileMd5, err := CalcMd5_2(fileName)
	if err != nil {
		return err, false
	}

	if !bytes.Equal(fileMd5, md5sum) {
		return errors.New("the hash check failed, fileMD5 not equal to md5sum"), false
	}
	return nil, true
}

func CalcMd5_2(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	const bufferSize = 65536

	hash := md5.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(f); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		hash.Write(buf[:n])
	}

	return hash.Sum(nil), nil
}
func HashSHA256File(filePath string) (string, error) {
	var hashValue string
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("openfile error!")
		return hashValue, err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return hashValue, err
	}
	hashInBytes := hash.Sum(nil)
	hashValue = hex.EncodeToString(hashInBytes)
	return hashValue, nil
}
