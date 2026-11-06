/*
 * sboutils
 * Copyright (C) 2026  Rogério Girão
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 * Project Page: https://github.com/rgcastrof/sboutils
 * Rogério Girão rogeriogirao1@proton.me
 */

package main

import (
	"fmt"
	"path/filepath"
	"io/fs"
	"os"
	"strings"
	"github.com/rgcastrof/sboutils/internal/pkg"
	"github.com/rgcastrof/sboutils/internal/config"
)

func main() {
	pkg := pkg.Package{Name:"", Path:"", Needed:false}

	for _, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			pkg.Name = arg
			continue
		}
	}

	filepath.WalkDir(config.RepoPath, func(s string, dir fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			os.Exit(1)
		}

		if dir.IsDir() && dir.Name() == pkg.Name {
			pkg.Needed = true
			pkg.Path = s
		}

		return nil
	})

	if pkg.Needed {
		pkg.DownloadURLs()
		pkg.CheckMD5()
	} else {
		fmt.Fprintf(os.Stderr, "Error: pkg not found.\n")
		os.Exit(1)
	}
}
