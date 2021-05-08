package dotenv

import (
	"os"
)

/*
 * Loads the envorinment variables given the file
 * names provided. Default file name is .env
 */
func Env(filenames ...string) {
	env := NewEnvFile(filenames)
	env.Load()
}

/*
 * Returns the value of a environment variable
 */
func Get(name string) (envVar string) {
	envVar = os.Getenv(name)
	return
}
