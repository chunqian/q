// Copyright 2016 Ryan Boehning. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package q

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"path/filepath"
)

// nolint: gochecknoglobals
var (
	// std is the singleton logger.
	std logger
)

var level string

func init() {
	level = "debug"

	dir, _ := os.Executable()
	path := filepath.Dir(dir)
	config, err := toml.LoadFile(path + "/q.toml")
	if err == nil {
		logL := config.Get("LOG_LEVEL")
		if logL == nil {
			fmt.Printf("q.Q Log Level is %s\n", level)
			return
		}
		level = logL.(string)
		fmt.Printf("q.Q Log Level is %s\n", level)
		return
	}

	path, _ = os.Getwd()
	config, err = toml.LoadFile(path + "/q.toml")
	if err == nil {
		logL := config.Get("LOG_LEVEL")
		if logL == nil {
			fmt.Printf("q.Q Log Level is %s\n", level)
			return
		}
		level = logL.(string)
		fmt.Printf("q.Q Log Level is %s\n", level)
		return
	}
}

// Q pretty-prints the given arguments to the $TMPDIR/q log file.
func Q(v ...interface{}) {

	std.mu.Lock()
	defer std.mu.Unlock()

	// Flush the buffered writes to disk.
	defer func() {
		if err := std.flush(); err != nil {
			fmt.Println(err)
		}
	}()

	if level == "prod" {

		args := formatArgsProd(v...)
		funcName, file, line, err := getCallerInfo()
		if err != nil {
			std.outputProd(args...) // no name=value printing
			return
		}

		// Print a header line if this q.Q() call is in a different file or
		// function than the previous q.Q() call, or if the 2s timer expired.
		// A header line looks like this: [14:00:36 main.go main.main:122].
		header := std.header(funcName, file, line)
		if header != "" {
			// fmt.Fprint(&std.buf, "\n", header, "\n")
			fmt.Print("\n", header, "\n")
		}

		// q.Q(foo, bar, baz) -> []string{"foo", "bar", "baz"}
		names, err := argNames(file, line)
		if err != nil {
			std.outputProd(args...) // no name=value printing
			return
		}

		// Convert the arguments to name=value strings.
		args = prependArgNameProd(names, args)
		// args = prependArgName(names, args)
		std.outputProd(args...)
		return
	}

	args := formatArgs(v...)
	funcName, file, line, err := getCallerInfo()
	if err != nil {
		std.output(args...) // no name=value printing
		return
	}

	// Print a header line if this q.Q() call is in a different file or
	// function than the previous q.Q() call, or if the 2s timer expired.
	// A header line looks like this: [14:00:36 main.go main.main:122].
	header := std.header(funcName, file, line)
	if header != "" {
		// fmt.Fprint(&std.buf, "\n", header, "\n")
		fmt.Print("\n", header, "\n")
	}

	// q.Q(foo, bar, baz) -> []string{"foo", "bar", "baz"}
	names, err := argNames(file, line)
	if err != nil {
		std.output(args...) // no name=value printing
		return
	}

	// Convert the arguments to name=value strings.
	args = prependArgName(names, args)
	std.output(args...)
}
