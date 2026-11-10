package pkg

import (
	"os"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

func calculate(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func validate(file, checksum string) error {
	md5Hash, err := calculate(file)
	if err != nil {
		return fmt.Errorf("error calculating sum %s: %w", file, err)
	}
	if strings.Compare(md5Hash, checksum) != 0 {
		return fmt.Errorf("CHECKSUM MISMATCH: %s", file)
	}
	return nil
}
