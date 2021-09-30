package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func pathExists(pathString string) bool {
	_, err := os.Stat(pathString)
	if err == nil {
		return true
	}
	return false
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getDirList(root string) []string {
	result := []string{}

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		check(err)

		if info.IsDir() {
			result = append(result, path)
		}

		return nil
	})

	return result
}

func findSassScope(s string) string {
	lines := strings.Split(s, "\n")

	for _, line := range lines {
		prefix := `$scope: "`
		suffix := `";`
		if strings.HasPrefix(line, prefix) && strings.HasSuffix(line, suffix) {
			return line[len(prefix) : len(line)-len(suffix)]
		}
	}

	return ""
}

func findTypescriptFiles(path string) []string {
	result := []string{}

	files, err := ioutil.ReadDir(path)
	check(err)

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".ts") {
			result = append(result, file.Name())
		}
	}

	return result
}

func findLastImportLine(lines []string) int {
	maxI := -1

	for i, line := range lines {
		if strings.HasPrefix(line, "import ") {
			maxI = i
		}
	}

	return maxI
}

func insertAfter(xs []string, index int, s string) []string {
	result := []string{}

	for i, x := range xs {
		result = append(result, x)
		if i == index {
			result = append(result, s)
		}
	}

	return result
}

func reverse(xs []string) []string {
	result := []string{}

	for i := len(xs) - 1; i >= 0; i-- {
		result = append(result, xs[i])
	}

	return result
}

func sassRemoveScopeLines(lines []string) []string {
	result := []string{}

	for _, line := range lines {
		if !strings.HasPrefix(line, "$scope: ") && line != ".#{$scope} {" {
			result = append(result, line)
		}
	}

	return result
}

func sassRemoveLastCurly(lines []string) []string {
	result := []string{}

	found := false
	for _, line := range reverse(lines) {
		if line == "}" {
			if found {
				result = append(result, line)
			} else {
				found = true
			}
		} else {
			result = append(result, line)
		}
	}

	return reverse(result)
}

func sassUnindent(lines []string) []string {
	result := []string{}

	for _, line := range lines {
		prefix := "    "
		if strings.HasPrefix(line, prefix) {
			result = append(result, line[len(prefix):])
		} else {
			result = append(result, line)
		}
	}

	return result
}

func sassRemoveEmptyStart(lines []string) []string {
	result := []string{}

	done := false
	for _, line := range lines {
		if line == "" {
			if done {
				result = append(result, line)
			}
		} else {
			done = true
			result = append(result, line)
		}
	}

	return result
}

func sassRemoveEmptyEnd(lines []string) []string {
	result := reverse(sassRemoveEmptyStart(reverse(lines)))
	result = append(result, "")
	return result
}

func migrateSassFile(path string) string {
	dat, err := ioutil.ReadFile(filepath.Join(path, "styles.scss"))
	check(err)

	fileContents := string(dat)

	scope := findSassScope(fileContents)

	if scope == "" {
		panic("Could not find sass scope")
	}

	lines := strings.Split(fileContents, "\n")
	lines = sassRemoveScopeLines(lines)
	lines = sassRemoveLastCurly(lines)
	lines = sassUnindent(lines)
	lines = sassRemoveEmptyStart(lines)
	lines = sassRemoveEmptyEnd(lines)
	fileContents = strings.Join(lines, "\n")
	fileContents = strings.ReplaceAll(fileContents, "&", ".SCOPE")

	err = ioutil.WriteFile(filepath.Join(path, "styles.css"), []byte(fileContents), 0644)
	check(err)

	err = os.Remove(filepath.Join(path, "styles.scss"))
	check(err)

	return scope
}

func typescriptContainsScope(lines []string, scope string) bool {
	for _, line := range lines {
		if strings.Contains(line, "'"+scope+"'") {
			return true
		}
	}

	return false
}

func typescriptReplaceScope(lines []string, scope string) []string {
	result := []string{}

	for _, line := range lines {
		line = strings.ReplaceAll(line, "'"+scope+"'", "cssScope")
		result = append(result, line)
	}

	return result
}

func typescriptAddCssScopeImport(lines []string) []string {
	lastImportLine := findLastImportLine(lines)
	if lastImportLine >= 0 {
		return insertAfter(lines, lastImportLine, "import { cssScope } from './css_scope';")
	}
	return lines
}

func typescriptHasCImport(lines []string) bool {
	for _, line := range lines {
		if line == "import { c } from './css_scope';" {
			return true
		}
	}

	return false
}

func typescriptRemoveCImport(lines []string) []string {
	result := []string{}

	for _, line := range lines {
		if line != "import { c } from './css_scope';" {
			result = append(result, line)
		}
	}

	return result
}

func typescriptAddScopedClassesImport(lines []string) []string {
	result := []string{}

	done := false
	for _, line := range lines {
		prefix := "import { "
		suffix := "/util';"
		if !done && strings.HasPrefix(line, prefix) && strings.HasSuffix(line, suffix) {
			line = prefix + "scopedClasses, " + line[len(prefix):]
			done = true
		}
		result = append(result, line)
	}

	return result
}

func typescriptAddCBeforeFunctions(lines []string) []string {
	result := []string{}

	done := false
	for _, line := range lines {
		if !done && strings.Contains(line, "function ") {
			result = append(result, "const c = scopedClasses(cssScope);", "")
			done = true
		}
		result = append(result, line)
	}

	return result
}

func migrateTypescriptFiles(path string, scope string) {
	for _, fileName := range findTypescriptFiles(path) {
		dat, err := ioutil.ReadFile(filepath.Join(path, fileName))
		check(err)

		fileContents := string(dat)
		lines := strings.Split(fileContents, "\n")

		if pathExists(filepath.Join(path, "css_scope.ts")) {
			if !typescriptHasCImport(lines) {
				continue
			}
			lines = typescriptRemoveCImport(lines)
			lines = typescriptAddCssScopeImport(lines)
			lines = typescriptAddScopedClassesImport(lines)
			lines = typescriptAddCBeforeFunctions(lines)
		} else {
			if !typescriptContainsScope(lines, scope) {
				continue
			}
			lines = typescriptReplaceScope(lines, scope)
			lines = typescriptAddCssScopeImport(lines)
		}

		fileContents = strings.Join(lines, "\n")

		err = ioutil.WriteFile(filepath.Join(path, fileName), []byte(fileContents), 0644)
		check(err)
	}
}

func main() {
	for _, path := range getDirList(".") {
		if pathExists(filepath.Join(path, "styles.scss")) {
			scope := migrateSassFile(path)
			migrateTypescriptFiles(path, scope)

			err := ioutil.WriteFile(filepath.Join(path, "css_scope.ts"), []byte("export const cssScope = '"+scope+"';\n\n"), 0644)
			check(err)
		}
	}
}
