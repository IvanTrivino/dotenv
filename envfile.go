package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
 * Creates a new EnvFile instance.
 */
func NewEnvFile(filenames []string) *EnvFile {
	return &EnvFile{filenames: filenames}
}

type EnvFile struct {
	filenames  []string
	envVarsMap map[string]string
}

/*
 * Main method for configuring the env vars. This method
 * takes the filename values provided during the struct
 * instantiation and pass them accross the whole process.
 */
func (ef *EnvFile) Load() {
	ef.checkFileNamesSlice()

	for _, fname := range ef.filenames {
		_ = ef.openFile(fname)
	}
}

/*
 * Check the length of the filenames provided, if length
 * is equal to zero appends the default filename ".env"
 */
func (ef *EnvFile) checkFileNamesSlice() {
	if len(ef.filenames) == 0 {
		ef.filenames = append(ef.filenames, ".env")
	}
}

/*
 * Open files according with the filenames provided and
 * parse/load the env vars.
 */
func (ef *EnvFile) openFile(fname string) (err error) {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer file.Close()

	ef.parseVars(file)

	return
}

/*
 * Scann the file line by line to extract the env vars.
 */
func (ef *EnvFile) parseVars(file *os.File) {
	scanner := bufio.NewScanner(file)
	var rows []string

	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading the file:", err)
	}

	ef.envVarsMap = ef.getKeyValuePairs(rows)

	err := ef.setEnvVars()
	if err != nil {
		fmt.Println("Error setting env vars", err)
	}
}

/*
 * Get the key-value pairs for each env var. Each key-value pair
 * should be separanted by "=" (e.g. KEY=value), otherwise env
 * var will not be loaded.
 */
func (ef *EnvFile) getKeyValuePairs(rows []string) map[string]string {
	//TODO: parse env lists
	keyValuePairs := make(map[string]string)

	for _, row := range rows {
		isComment := ef.checkIsComment(row)

		if !isComment {
			rowSplit := strings.Split(row, "=")

			if len(rowSplit) == 2 {
				keyValuePairs[rowSplit[0]] = rowSplit[1]
			}
		}
	}

	return keyValuePairs
}

/*
 * Check if the current line is a comment. Every comment must
 * begin with "#", otherwise will be taken as a common line.
 */
func (ef *EnvFile) checkIsComment(row string) (isComment bool) {
	isComment = false

	if (len(row) > 0) && (row[0] == '#') {
		isComment = true
	}

	return
}

/*
 * Set env vars parsed.
 */
func (ef *EnvFile) setEnvVars() (err error) {
	for k := range ef.envVarsMap {
		err = os.Setenv(k, ef.envVarsMap[k])
	}

	return
}
