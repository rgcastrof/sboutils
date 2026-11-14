package pkg

import (
	"net/http"
	"time"
	"os"
	"fmt"
	"io"
	"path/filepath"
)

func download(url, destdir string, fileName string) error {
	if err := os.MkdirAll(destdir, 0755); err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Get(url)
	if err == nil {
		fmt.Printf("Downloading file from url: %s\n", url)
	} else {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	file, err := os.Create(filepath.Join(destdir, fileName))
	if err != nil {
		return fmt.Errorf("failed create download file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed save file: %w", err)
	}
	return nil
}
