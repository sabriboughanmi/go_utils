package utils

import (
	"math/rand"
	"strings"
	"time"
)

//ReplaceKeys formats a string using key value pairs
func ReplaceKeys(s string, keyValues map[string]string) string {
	var newString = s
	for key, value := range keyValues {
		newString = strings.Replace(newString, "{"+key+"}", value, 1)
	}
	return newString
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//Returns a Random AlphaNumeric String with a given length.
//Note! passing the same seed returns the same Random String
func RandomStringWithSeed(length, seed int64) string {
	return RandomstringCharsetSeed(length, seed, charset)
}

//Returns a Random String from specified Charset.
//Note! passing the same seed returns the same Random String
func RandomstringCharsetSeed(length, seed int64, charset string) string {
	var seededRand = rand.New(rand.NewSource(seed))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

//Returns a Random AlphaNumeric String with a given length.
//Note! this Function sleeps for a time.Nanosecond to return a different Value every call, if you care about Performance use  RandomStringWithSeed instead.
func RandomString(length int) string {
	return RandomstringCharset(length, charset)
}

//Returns a Random String from specified Charset.
//Note! this Function sleeps for a time.Nanosecond to return a different Value every call, if you care about Performance use  RandomstringCharsetSeed instead.
func RandomstringCharset(length int, charset string) string {
	time.Sleep(time.Nanosecond)
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}



// ContainsProfanity returns if a string contains BadWords.
func ContainsProfanity(str string) (bool, string) {
	for _, forbiddenWord := range BadWords {
		if test := strings.Index(strings.ToLower(str), forbiddenWord); test > -1 {
			//if test := strings.EqualFold(strings.ToLower(word), forbiddenWord); test == true {
			return true, forbiddenWord
		}
	}
	return false, ""
}

// ContainsSpecialCharacters returns if a string contains SpecialCharacters.
func ContainsSpecialCharacters(str string) (bool, string) {
	for _, forbiddenCharacter := range SpecialCharacters {
		if test := strings.Index(strings.ToLower(str), forbiddenCharacter); test > -1 {
			//if test := strings.EqualFold(strings.ToLower(word), forbiddenWord); test == true {
			return true, forbiddenCharacter
		}
	}
	return false, ""
}
