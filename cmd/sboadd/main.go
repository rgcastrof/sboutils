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
	"os"
	"github.com/rgcastrof/sboutils/cli"
	"github.com/rgcastrof/sboutils/internal/pkg"
)

func main() {
	c := cli.SboAdd{
		Arg: "",
		Dflag: false,
		Cflag: false,
		Sflag: false,
	}

	err := cli.Run(&c, os.Args[1:], func() error {
		p, err := pkg.Find(c.Arg)
		if err != nil {
			return err
		}

		if c.Dflag {
			p.DownloadURLs()
		}
		if c.Cflag {
			p.CheckMD5()
		}

		return nil
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
