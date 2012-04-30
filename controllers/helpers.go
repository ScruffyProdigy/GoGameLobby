package controllers

import (
	"strings"
)

const unsafe = "/\\\"`^%+?#&{}|<>"

func removeChars(original,removed string) (result string) {
	i := 0
	for i < len(original) {			
		if strings.ContainsRune(removed,rune(original[i])) {
			original = original[:i] + original[i+1:]
		} else {
			i++
		}
	}
	return original
}
//"$&+,/:;=?@#%"
func urlify(a string) string {
	b := strings.ToLower(a)				// lower case it
	c := strings.Replace(b," ","-",-1)	// replace the spaces with more url friendly dashes
	d := removeChars(c,unsafe)	// remove all unsafe characters
	return d
}