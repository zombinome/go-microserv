package microserv

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func isEmpty(val string) bool {
	return len(val) == 0
}

var randSource = rand.NewSource(time.Now().UnixNano())
var randGen = rand.New(randSource)

const maxRandPart int64 = 1 << 20

var base64CharMap = []rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
	'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F',
	'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
	'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '_',
}

func GenerateCorrelationId() string {
	var now = time.Now().UTC()
	var hours = now.Hour()
	if hours > 12 {
		hours = hours - 12
	}

	// Getting total number of seconds since last midday or midnight
	// (12h should be more than enough to end any request)
	var seconds = uint64(hours*3600 + now.Minute()*60 + now.Second())
	var milliseconds = uint64(now.Nanosecond() / 1000000) // 28 bit

	var timeComponent = (seconds*1000 + milliseconds) << 20
	var randPart = uint64(randGen.Int63n(maxRandPart))

	var fullNumber = timeComponent + randPart
	var chars = make([]rune, 8)
	var mask uint64 = 0b00111111
	var shift = 0
	for i := 7; i >= 0; i-- {
		var charMask = mask << shift

		var runeIndex = (fullNumber & charMask) >> shift
		chars[i] = base64CharMap[runeIndex]
		shift += 6
	}
	return string(chars)
}

func EnsurePathExists(pathToFile string) {
	var parentDir = filepath.Dir(pathToFile)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		os.MkdirAll(parentDir, 0700)
	}
}
