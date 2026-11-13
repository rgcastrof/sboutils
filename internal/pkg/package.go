package pkg

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/fs"
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
	fileChecks []fileCheck
}

func (pkg *Package) Find(pkgName string) error {
	if pkgName == "" {
		return errors.New("package name is empty")
	}

	err := filepath.WalkDir(config.RepoPath, func(s string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() && dir.Name() == pkgName {
			pkg.Name = pkgName
			pkg.Needed = true
			pkg.Path = s
		}

		return nil
	})
	if err != nil {
		return err
	}

	if !pkg.Needed {
		return fmt.Errorf("package %s not found", pkgName)
	}
	return nil
}

func (pkg *Package) getInfo() string {
	return filepath.Join(pkg.Path, pkg.Name + ".info")
}

func (pkg *Package) getInfoElements(variable string) []string {
	var elements []string
	var line string

	shouldContinue := false
	infoPath := pkg.getInfo()

	file, err := os.Open(infoPath)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
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

func (pkg *Package) mapChecks() {
	if pkg.fileChecks == nil {
		urls := pkg.getInfoElements("DOWNLOAD")
		sums := pkg.getInfoElements("MD5SUM")

		if urls == nil {
			urls = pkg.getInfoElements("DOWNLOAD_x86_64")
		}

		if sums == nil {
			sums = pkg.getInfoElements("MD5SUM_x86_64")
		}

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

func (pkg *Package) DownloadURLs() {
	pkg.mapChecks()
	dispatcher := NewDispatcher(pkg.fileChecks, 12)
	err := dispatcher.Run(context.Background(), func(fileMap fileCheck) error {
		distName := filepath.Base(fileMap.URL)
		return download(fileMap.URL, filepath.Join(config.DistDir, distName))
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		fmt.Println("All files downloaded successful")
	}
}

func (pkg *Package) CheckMD5() {
	pkg.mapChecks()
	fmt.Println("Calculating MD5SUM...")
	dispatcher := NewDispatcher(pkg.fileChecks, 12)
	err := dispatcher.Run(context.Background(), func(fileMap fileCheck) error {
		distPath, err := pkg.getDistPath(filepath.Base(fileMap.URL))
		if err != nil {
			return err
		}

		fmt.Printf("Verifying file: %s\n", distPath)
		return validate(distPath, fileMap.MD5)
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	} else {
		fmt.Println("All sums match")
	}
}
