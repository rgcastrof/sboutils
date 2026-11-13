package pkg

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/rgcastrof/sboutils/internal/config"
)

type fileCheck struct {
	URL string
	MD5 string
}

type Package struct {
	Name string
	Path string
	Needed bool

type Semaphore struct {
	Channel chan struct{}
}

func (s *Semaphore) Acquire() {
	s.Channel <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.Channel
	fileChecks []fileCheck
}

func (pkg *Package) getInfo() (string, error) {
	if len(pkg.Name) == 0 {
		return "", errors.New("Package name is empty")
	}
	return filepath.Join(pkg.Path, pkg.Name + ".info"), nil
}

func (pkg *Package) getInfoElements(variable string) []string {
	var elements []string
	var line string

	shouldContinue := false
	infoPath, err := pkg.getInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
		os.Exit(1)
	}

	file, err := os.Open(infoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		line = scanner.Text()
		if strings.HasPrefix(line, variable) {
			if line[len(line)-1] == '\\' {
				shouldContinue = true
				elements = append(elements, strings.TrimSpace(line[len(variable + "=\""):len(line)-1]))
				continue
			} else {
				element := line[len(variable + "=\""):len(line)-1]
				if len(element) == 0 || strings.Compare(element, "UNSUPPORTED") == 0 {
					return nil
				}
				return append(elements, element)
			}
		}

		if shouldContinue {
			elements = append(elements, strings.TrimSpace(line[:len(line)-1]))
			if !strings.HasSuffix(line, "\\") {
				break
			}
		}

	}
	return elements
}

func (pkg *Package) getDistPath(distName string) (string, error) {
	dir, err := os.ReadDir(config.DistDir)
	if err != nil {
		return "", fmt.Errorf("Cannot read directory: %s", err.Error())
	}

	for _, entry := range dir {
		if strings.Compare(entry.Name(), distName) == 0 {
			return filepath.Join(config.DistDir, entry.Name()), nil
		}
	}
	return "", fmt.Errorf("File %s not found", distName)
}

func (pkg *Package) DownloadURLs() {
	var wg sync.WaitGroup
	sem := &Semaphore{Channel: make(chan struct{}, 12)}
func (pkg *Package) mapChecks() {
	if pkg.fileChecks == nil {
		urls := pkg.getInfoElements("DOWNLOAD")
		sums := pkg.getInfoElements("MD5SUM")

	pkg.urls = pkg.getInfoElements("DOWNLOAD")
	pkg.sums = pkg.getInfoElements("MD5SUM")
		if urls == nil {
			urls = pkg.getInfoElements("DOWNLOAD_x86_64")
		}

	if pkg.urls == nil {
		pkg.urls = pkg.getInfoElements("DOWNLOAD_x86_64")
	}
		if sums == nil {
			sums = pkg.getInfoElements("MD5SUM_x86_64")
		}

	if pkg.sums == nil {
		pkg.sums = pkg.getInfoElements("MD5SUM_x86_64")
		if (len(urls) == len(sums)) {
			for i, u := range urls {
				pkg.fileChecks = append(pkg.fileChecks, fileCheck{
					URL: u,
					MD5: sums[i],
				})
			}
		}
	}
}

	wg.Add(len(pkg.urls))
	for _, u := range pkg.urls {
		distName := pkg.getDistName(u)
		go func() {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()
			fmt.Printf("Downloading file from url: %s\n", u)
			err := downloader.Download(u, filepath.Join(config.DistDir, distName))
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				return
			}
		}()
	}
	wg.Wait()
}

func (pkg *Package) CheckMD5() {
	var wg sync.WaitGroup
	sem := &Semaphore{Channel: make(chan struct{}, 12)}

	fmt.Println("Calculating MD5SUM...")
	if len(pkg.urls) == len(pkg.sums) {
		wg.Add(len(pkg.sums))
		for i := range pkg.urls {
			distName := pkg.getDistName(pkg.urls[i])
			distPath := pkg.getDistPath(distName)
			go func() {
				defer wg.Done()

				sem.Acquire()
				defer sem.Release()
				err := md5sum.Validate(distPath, pkg.sums[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
					return
					
				}
			}()
		}
		wg.Wait()
	}
}
