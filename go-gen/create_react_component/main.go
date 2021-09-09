package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"unicode"
)

func splitPascal(s string) []string {
	words := []string{}
	word := strings.Builder{}
	var previousR rune

	for i, r := range s {
		if i != 0 {
			if unicode.IsUpper(r) {
				words = append(words, word.String())
				word = strings.Builder{}
			} else if unicode.IsDigit(r) {
				if !unicode.IsDigit(previousR) {
					words = append(words, word.String())
					word = strings.Builder{}
				}
			} else if unicode.IsDigit(previousR) {
				words = append(words, word.String())
				word = strings.Builder{}
			}
		}

		word.WriteRune(r)
		previousR = r
	}

	if word.Len() != 0 {
		words = append(words, word.String())
	}

	return words
}

func joinSnake(words []string) string {
	lowerWords := []string{}
	for _, word := range words {
		lowerWords = append(lowerWords, strings.ToLower(word))
	}

	return strings.Join(lowerWords, "_")
}

func uuid() (string, error) {
	var err error

	bytes := make([]byte, 16)

	_, err = rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func pathExists(pathString string) bool {
	_, err := os.Stat(pathString)
	if err == nil {
		return true
	}
	return false
}

func repeat(s string, n int) []string {
	ss := []string{}
	for i := 0; i < n; i++ {
		ss = append(ss, s)
	}
	return ss
}

func findUtilDirectory() string {
	defaultDots := path.Join("..", "..")
	for i := 0; i < 50; i++ {
		if pathExists(path.Join(path.Join(repeat("..", i)...), ".git")) {
			if i-2 < 1 {
				return defaultDots
			} else {
				return path.Join(repeat("..", i-2)...)
			}
		}
	}
	return defaultDots
}

func findUtilPath() string {
	return path.Join(findUtilDirectory(), "util")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./go-gen-create-react-component name")
		return
	}

	componentNamePascal := os.Args[1]
	componentNameSnake := joinSnake(splitPascal(componentNamePascal))

	err := os.Mkdir(componentNameSnake, 0755)
	if err != nil {
		panic(err)
	}

	u, err := uuid()
	if err != nil {
		panic(err)
	}

	s := "export const cssScope = 'scope-" + u + "';\n"
	ioutil.WriteFile(path.Join(componentNameSnake, "css_scope.ts"), []byte(s), 0644)

	s = "\n"
	ioutil.WriteFile(path.Join(componentNameSnake, "styles.css"), []byte(s), 0644)

	s = `import * as React from 'react';
import * as mobxReactLite from 'mobx-react-lite';
import { h, scopedClasses } from '` + findUtilPath() + `';
import { cssScope } from './css_scope';

type Props = {};

const c = scopedClasses(cssScope);

export const ` + componentNamePascal + ` = mobxReactLite.observer(function ` + componentNamePascal + `(props: Props) {
    const propsRef = React.useRef(props);
    propsRef.current = props;

    return h('div', { className: c('root') },
    );
});
`
	componentFilePath := path.Join(componentNameSnake, componentNameSnake+".ts")
	ioutil.WriteFile(componentFilePath, []byte(s), 0644)

	cmd := exec.Command("/usr/local/bin/code", componentFilePath)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
