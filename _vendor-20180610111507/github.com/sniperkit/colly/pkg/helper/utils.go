package helper

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func BetweenStr(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func BeforeStr(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func AfterStr(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

// GetStringInBetween Returns empty string if no start string found
func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str, end)
	return str[s:e]
}

func EnsurePathExists(path string) error {
	i := strings.LastIndexByte(path, '/')
	if i >= 0 {
		dir := path[:i]
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func IsRemoteURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func ResolvePath(basePath string, relativePath string) string {
	if IsRemoteURL(basePath) || IsRemoteURL(relativePath) || filepath.IsAbs(relativePath) {
		return relativePath
	}
	return filepath.Join(filepath.Dir(basePath), relativePath)
}

func RoundU(val float64) int {
	if val > 0 {
		return int(val + 1.0)
	}
	return int(val)
}

// return the source filename after the last slash
func ChopPath(original string) string {
	i := strings.LastIndex(original, "/")
	if i == -1 {
		return original
	} else {
		return original[i+1:]
	}
}

func ParseTimeStamp(utime string) (*time.Time, error) {
	i, err := strconv.ParseInt(utime, 10, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(i, 0)
	return &t, nil
}

func ToBytes(input string) []byte {
	return []byte(input)
}

func MapToString(input map[string]interface{}) string {
	return ToString(input)
}

func ToString(obj interface{}) string {
	return fmt.Sprintf("%v", obj)
}

func ConvertInterfaceArray(input []map[string]interface{}) []interface{} {
	results := make([]interface{}, len(input))
	for _, result := range input {
		// resultSlice := result.(map[string]interface{})
		// pp.Println("resultSlice=", resultSlice)
		results = append(results, result)
	}
	return results
}

func ConvertMapToInterface(input map[string]interface{}) []interface{} {
	results := make([]interface{}, len(input))
	for _, result := range input {
		resultSlice := result.(interface{})
		results = append(results, resultSlice)
	}
	return results
}

func Shuffle(slc []interface{}) []interface{} {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
	return slc
}

func ShuffleInts(slc []int) []int {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
	return slc
}

func ShuffleStrings(slc []string) []string {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
	return slc
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func UpdateRequestDelay(beat int, unit string) (delay time.Duration) {
	input := time.Duration(beat)
	switch strings.ToLower(unit) {
	case "microsecond":
		delay = input * time.Microsecond
	case "millisecond":
		delay = input * time.Millisecond
	case "minute":
		delay = input * time.Minute
	case "hour":
		delay = input * time.Hour
	case "second":
		fallthrough
	default:
		delay = input * time.Second
	}
	return
}
