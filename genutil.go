// Package genutil collects in place the utility functions used in various golang scripts
package genutil

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type bslice []byte

//================================================================================

// GzFile is used to write to regular or gz file, removing existing compression variant first
type GzFile struct {
	fo   *os.File
	ww   *bufio.Writer
	wwgz *gzip.Writer
}

func (us GzFile) Write(pp []byte) (nn int, err error) {
	switch {
	case us.wwgz != nil:
		nn, err = us.wwgz.Write(pp)
	case us.ww != nil:
		nn, err = us.ww.Write(pp)
	}
	return
}

// WriteString writes to the (un)compressed stream
func (us GzFile) WriteString(ss string) (nn int, err error) {
	switch {
	case us.wwgz != nil:
		nn, err = us.wwgz.Write([]byte(ss))
	case us.ww != nil:
		nn, err = us.ww.WriteString(ss)
	}
	return
}

// Close flushes and closes
func (us GzFile) Close() {
	switch {
	case us.wwgz != nil:
		us.wwgz.Flush()
		us.wwgz.Close()
	}
	if us.ww != nil {
		us.ww.Flush()
		us.fo.Close()
	}
}

// OpenGzFile Opens a file for buffered writing, optionally using gzip compression
func OpenGzFile(_fname string) GzFile {
	self := new(GzFile)
	var err error

	switch {
	case strings.HasPrefix(_fname, "/dev/"):
	default:
		ofname, ofcode := WritableFilename(_fname)
		if false {
			fmt.Println("Removed existing file: %s, ofcode=%d\n", ofname, ofcode)
		}
	}

	self.fo, err = os.Create(_fname)
	if err != nil {
		panic(err)
	}
	self.ww = bufio.NewWriter(self.fo)
	switch {
	case strings.HasSuffix(_fname, ".gz"):
		self.wwgz = gzip.NewWriter(self.ww)
	}
	return (*self)
}

//================================================================================

var (
	comma  byte = ','
	slash  byte = '|'
	sepmap      = map[string]string{
		"comma": ",", ",": ",",
		"pipe": "|", "|": "|",
		"bang": "!", "!": "!",
		"plus": "+", "+": "+",
		"space": " ", " ": " ", "blank": " ",
		"tab": "	", "	": "	",
		"colon": ":", ":": ":",
		"newline": "\n", `
`: "\n",
	}
)

// ================================================================================

// Hostname retrieves hostname
func Hostname() string {
	return strings.TrimSpace(BashExecOrDie(false, fmt.Sprintf("hostname -s"), "."))
}

// Millions is shorthand
func Millions(num float64) string {
	return fmt.Sprintf("%.2fMM", num/1000000.0)
}

// Kilos is shorthand
func Kilos(num float64) string {
	return fmt.Sprintf("%.2fK", num/1000.0)
}

// KB2GB is shorthand
func KB2GB(_kb string) string {
	return fmt.Sprintf("%.02f", StrToFloat(_kb)/1000000.0)
}

// Float64ToHuman is shorthand
func Float64ToHuman(num float64) string {
	if num > 1000000.0 {
		return Millions(num)
	}

	return Kilos(num)
}

// Thousands drops fractional part and inserts commas to separate thousands
func Thousands(_num float64) string {
	isneg := _num < 0
	str := fmt.Sprintf("%.0f", math.Abs(_num))
	nstr := len(str)
	ostr := ""
	for ii := 0; ii < nstr; ii += 3 {
		ostr1 := ""
		if nstr-ii-3 >= 0 {
			ostr1 = str[nstr-ii-3 : nstr-ii]
		} else {
			ostr1 = str[:nstr-ii]
		}
		if ii > 0 {
			ostr = "," + ostr
		}
		ostr = ostr1 + ostr
	}
	if isneg {
		return "-" + ostr
	}
	return ostr
}

// SepSplit2 is shorthand splitter
func SepSplit2(str, _sep string) (part0, part1 string) {
	parts := strings.SplitN(str, _sep, 2)
	switch len(parts) - 1 {
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// SepSplit4 is shorthand splitter
func SepSplit4(str, _sep string) (part0, part1, part2, part3 string) {
	parts := strings.SplitN(str, _sep, 4)
	switch len(parts) - 1 {
	case 3:
		part3 = parts[3]
		fallthrough
	case 2:
		part2 = parts[2]
		fallthrough
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// EqualsSplit2 is shorthand splitter
func EqualsSplit2(str string) (part0, part1 string) {
	parts := strings.SplitN(str, "=", 2)
	switch len(parts) - 1 {
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// EqualsSplit2Trimmed is shorthand splitter
func EqualsSplit2Trimmed(str string) (part0, part1 string) {
	parts := strings.SplitN(str, "=", 2)
	switch len(parts) - 1 {
	case 1:
		part1 = strings.TrimSpace(parts[1])
		fallthrough
	case 0:
		part0 = strings.TrimSpace(parts[0])
	}
	return
}

// EqualsSplit6 is shorthand splitter
func EqualsSplit6(str string) (part0, part1, part2, part3, part4, part5 string) {
	parts := strings.SplitN(str, "=", 5)
	switch len(parts) - 1 {
	case 5:
		part5 = parts[5]
		fallthrough
	case 4:
		part4 = parts[4]
		fallthrough
	case 3:
		part3 = parts[3]
		fallthrough
	case 2:
		part2 = parts[2]
		fallthrough
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// ColonSplit2 is shorthand splitter
func ColonSplit2(str string) (part1, part2 string) {
	parts := strings.SplitN(str, ":", 2)
	if len(parts) == 2 {
		part1 = parts[0]
		part2 = parts[1]
	} else {
		part1 = parts[0]
		part2 = ""
	}

	return
}

// ColonSplit3 is shorthand splitter
func ColonSplit3(str string) (part1, part2, part3 string) {
	parts := strings.SplitN(str, ":", 3)
	if len(parts) == 3 {
		part1 = parts[0]
		part2 = parts[1]
		part3 = parts[2]
	} else {
		part1 = parts[0]
		part2 = ""
		part3 = ""
	}

	return
}

// ColonSplit4 is shorthand splitter
func ColonSplit4(str string) (part0, part1, part2, part3 string) {
	parts := strings.SplitN(str, ":", 4)
	switch len(parts) - 1 {
	case 3:
		part3 = parts[3]
		fallthrough
	case 2:
		part2 = parts[2]
		fallthrough
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// ColonSplit5 is shorthand splitter
func ColonSplit5(str string) (part0, part1, part2, part3, part4 string) {
	parts := strings.SplitN(str, ":", 4+1)
	switch len(parts) - 1 {
	case 4:
		part4 = parts[4]
		fallthrough
	case 3:
		part3 = parts[3]
		fallthrough
	case 2:
		part2 = parts[2]
		fallthrough
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// ColonSplit5Len is shorthand splitter, which also returns the number of parts found
func ColonSplit5Len(str string) (nn int, part0, part1, part2, part3, part4 string) {
	parts := strings.SplitN(str, ":", 4+1)
	nn = len(parts)
	switch nn - 1 {
	case 4:
		part4 = parts[4]
		fallthrough
	case 3:
		part3 = parts[3]
		fallthrough
	case 2:
		part2 = parts[2]
		fallthrough
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// HatSplit2 is shorthand splitter
func HatSplit2(str string) (part0, part1 string) {
	parts := strings.SplitN(str, "^", 2)
	switch len(parts) - 1 {
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// DashSplit2 is shorthand splitter
func DashSplit2(str string) (part0, part1 string) {
	parts := strings.SplitN(str, "-", 2)
	switch len(parts) - 1 {
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// CommaSplit2 is shorthand splitter
func CommaSplit2(str string) (part0, part1 string) {
	parts := strings.SplitN(str, ",", 2)
	switch len(parts) - 1 {
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// CommaSplit7 is shorthand splitter
func CommaSplit7(str string) (part0, part1, part2, part3, part4, part5, part6 string) {
	parts := strings.SplitN(str, ",", 7)
	switch len(parts) - 1 {
	case 6:
		part6 = parts[6]
		fallthrough
	case 5:
		part5 = parts[5]
		fallthrough
	case 4:
		part4 = parts[4]
		fallthrough
	case 3:
		part3 = parts[3]
		fallthrough
	case 2:
		part2 = parts[2]
		fallthrough
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// SpaceSplit2 is shorthand splitter
func SpaceSplit2(str string) (part0, part1 string) {
	parts := strings.SplitN(str, " ", 2)
	switch len(parts) - 1 {
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// SlashSplit2 is shorthand splitter
func SlashSplit2(str string) (part0, part1 string) {
	parts := strings.SplitN(str, "/", 2)
	switch len(parts) - 1 {
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// SlashSplit3 is shorthand splitter
func SlashSplit3(str string) (part0, part1, part2 string) {
	parts := strings.SplitN(str, "/", 3)
	switch len(parts) - 1 {
	case 2:
		part2 = parts[2]
		fallthrough
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// AnySplit splits input string into parts, using any char in splitstr
func AnySplit(str, splitstr string) []string {
	aset := NewBoolMap()
	aset[str] = true
	for ii := 0; ii < len(splitstr); ii++ {
		bset := NewBoolMap()
		for kk := range aset {
			UpdateBoolMapFromCsv(&bset, kk, splitstr[ii:ii+1])
		}
		aset = bset
	}
	return KeysBoolMap(&aset)
}

// AnySplit2 splits input string in upto 2 parts, using any char in splitstr
func AnySplit2(str, splitstr string) (splitter, part0, part1 string) {
	maxii, maxlen := -1, -1
	for ii := 0; ii < len(splitstr); ii++ {
		splitter := splitstr[ii : ii+1]
		parts := strings.SplitN(str, splitter, 2)
		if len(parts) > maxlen {
			maxii = ii
			maxlen = len(parts)
		}
	}
	splitter = splitstr[maxii : maxii+1]
	parts := strings.SplitN(str, splitter, 2)
	switch len(parts) - 1 {
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// AnySplit3 splits input string in upto 3 parts, using any char in splitstr
func AnySplit3(str, splitstr string) (splitter0, splitter2, part0, part1, part2 string) {
	splitter0, splitter2, part0, part1, part2 = "", "", "", "", ""
	idx0 := strings.IndexAny(str, splitstr)
	idx1 := -1
	if (idx0 >= 0) && (idx0 < len(str)) {
		idx1 = strings.IndexAny(str[idx0+1:], splitstr)
	}
	switch {
	case idx0 < 0:
		part0 = str
	case idx1 < 0:
		splitter0, part0, part1 = str[idx0:idx0+1], str[:idx0], str[idx0+1:]
	default:
		idx1 += idx0 + 1
		splitter0, splitter2 = str[idx0:idx0+1], str[idx1:idx1+1]
		part0, part1, part2 = str[:idx0], str[idx0+1:idx1], str[idx1+1:]
	}
	fmt.Printf("genutil.AnySplit3: str(%s) splitstr(%s) splitter0(%s) splitter2(%s) part0(%s) part1(%s) part2(%s)\n", splitstr, str, splitter0, splitter2, part0, part1, part2)
	return
}

// SlashSplit5 is shorthand splitter
func SlashSplit5(str string) (part0, part1, part2, part3, part4 string) {
	parts := strings.SplitN(str, "/", 5)
	switch len(parts) - 1 {
	case 4:
		part4 = parts[4]
		fallthrough
	case 3:
		part3 = parts[3]
		fallthrough
	case 2:
		part2 = parts[2]
		fallthrough
	case 1:
		part1 = parts[1]
		fallthrough
	case 0:
		part0 = parts[0]
	}
	return
}

// StrDropComponent drops the indicated component
func StrDropComponent(_str, _sep string, _drop int, _doPanic bool) (string, string) {
	arr := strings.Split(_str, _sep)
	nn := len(arr)
	if _drop < nn {
		switch _drop {
		case 0:
			return StrTernary(nn > 1, strings.Join(arr[1:], _sep), ""), arr[0]
		case (nn - 1):
			return strings.Join(arr[:(nn-1)], _sep), arr[nn-1]
		}
		return strings.Join(arr[:(_drop-1)], _sep) + _sep + strings.Join(arr[(_drop+1):], _sep), arr[_drop]
	}
	if _doPanic {
		log.Panicf("genutil.StrDropComponent: Input string (%s) does not not have element at index (%d)\n", _str, _drop)
	}
	return _str, ""
}

// StrReplaceComponent replaces the indicated component
func StrReplaceComponent(_str, _sep string, _reploc int, _rep string, _doPanic bool) (string, string) {
	arr := strings.Split(_str, _sep)
	nn := len(arr)
	if _reploc < nn {
		switch _reploc {
		case 0:
			return StrTernary(nn > 1, _rep+_sep+strings.Join(arr[1:], _sep), _rep), arr[0]
		case (nn - 1):
			return strings.Join(arr[:(nn-1)], _sep) + _sep + _rep, arr[nn-1]
		}
		return strings.Join(arr[:(_reploc-1)], _sep) + _sep + _rep + _sep + strings.Join(arr[(_reploc+1):], _sep), arr[_reploc]
	}
	if _doPanic {
		log.Panicf("genutil.StrDropComponent: Input string (%s) does not not have element at index (%d)\n", _str, _reploc)
	}
	return _str, ""
}

// GetSplitTrimmedPart returns the numbered part (or _badstr if there was an error)
func GetSplitTrimmedPart(_str, _sep, _badstr string, _partno int) string {
	parts := strings.Split(strings.TrimSpace(_str), _sep)
	nn := len(parts)
	if _partno >= 0 {
		if _partno > nn-1 {
			return _badstr
		}
		return parts[_partno]
	}
	_partno += nn // -1 becomes (nn-1) which is the last part
	if _partno < 0 {
		return _badstr
	}
	return parts[_partno]
}

// GetSplitTrimmedPartFloat64 returns the numbered part (or _bad if there was an error)
func GetSplitTrimmedPartFloat64(_str, _sep string, _bad float64, _partno int) float64 {
	parts := strings.Split(strings.TrimSpace(_str), _sep)
	if len(parts) < _partno+1 {
		return _bad
	}
	return StrToFloat(parts[_partno])
}

// GetSplitTrimmedPartInt64 returns the numbered part (or _bad if there was an error)
func GetSplitTrimmedPartInt64(_str, _sep string, _bad int64, _partno int) int64 {
	parts := strings.Split(strings.TrimSpace(_str), _sep)
	if len(parts) < _partno+1 {
		return _bad
	}
	return ToInt(parts[_partno], _bad)
}

//================================================================================

// AbsInt64 is shorthand
func AbsInt64(_ival int64) int64 {
	if _ival < 0 {
		return -_ival
	}
	return _ival
}

// AbsInt is shorthand
func AbsInt(_ival int) int {
	if _ival < 0 {
		return -_ival
	}
	return _ival
}

// MinInt64 is shorthand
func MinInt64(_ival1, _ival2 int64) int64 {
	if _ival1 < _ival2 {
		return _ival1
	}
	return _ival2
}

// MinInt is shorthand
func MinInt(_ival1, _ival2 int) int {
	if _ival1 < _ival2 {
		return _ival1
	}
	return _ival2
}

// MaxInt64 is shorthand
func MaxInt64(_ival1, _ival2 int64) int64 {
	if _ival1 > _ival2 {
		return _ival1
	}
	return _ival2
}

// MaxInt is shorthand
func MaxInt(_ival1, _ival2 int) int {
	if _ival1 > _ival2 {
		return _ival1
	}
	return _ival2
}

// MinFloat is shorthand
func MinFloat(_ival1, _ival2 float64) float64 {
	if _ival1 < _ival2 {
		return _ival1
	}
	return _ival2
}

// MaxFloat is shorthand
func MaxFloat(_ival1, _ival2 float64) float64 {
	if _ival1 > _ival2 {
		return _ival1
	}
	return _ival2
}

// StrSin checks if _str is in one of the items in list specified as "(foo|bar|baz|...)"
func StrSin(_str, _sin string) bool {
	// fmt.Println("_str=", _str, "_sin=", _sin)
	if _str == "" {
		return true
	}
	if _sin == "" {
		return false
	}
	from := 0
	if len(_sin) > 1 {
		if _sin[0] == '(' {
			from++
		}
	}
	to := len(_sin)
	if len(_sin) > 1 {
		if _sin[to-1] == ')' {
			to--
		}
	}
	sinparts := strings.Split(_sin[from:to], "|")
	for _, sin := range sinparts {
		if _str == sin {
			return true
		}
	}
	return false
}

// StrIsInt checks whether the string is an int
func StrIsInt(_str string) bool {
	for _, ch := range _str {
		if !(('0' <= ch) && (ch <= '9')) {
			return false
		}
	}
	return true
}

// StrAbs returns the abs value of a string, as string
func StrAbs(_num string) string {
	return fmt.Sprintf("%f", math.Abs(StrToFloat(_num)))
}

// ToBool converts string to bool, with default
func ToBool(_str string, _def bool) bool {
	bval, err := strconv.ParseBool(_str)
	if err == nil {
		return bval
	}
	return _def
}

// ToInt converts string to int64, with default
func ToInt(_str string, _def int64) int64 {
	if strings.Contains(_str, "e+") || strings.Contains(_str, "E+") {
		nm := ToFloat(([]byte)(_str))
		return (int64(nm))
	}
	num, err := strconv.ParseInt(_str, 10, 64)
	if err == nil {
		return num
	}
	return _def
}

// Toint converts string to int, with default
func Toint(_str string, _def int) int {
	if strings.Contains(_str, "e+") || strings.Contains(_str, "E+") {
		nm := ToFloat(([]byte)(_str))
		return (int(nm))
	}
	num, err := strconv.ParseInt(_str, 10, 64)
	if err == nil {
		return int(num)
	}
	return _def
}

// Toint0 converts string to int, with default 0
func Toint0(_str string) int {
	if strings.Contains(_str, "e+") || strings.Contains(_str, "E+") {
		nm := ToFloat(([]byte)(_str))
		return (int(nm))
	}
	num, err := strconv.ParseInt(_str, 10, 64)
	if err == nil {
		return int(num)
	}
	return 0
}

// ByteToint0 converts bytestring to int, with default 0
func ByteToint0(_str []byte) int {
	if bytes.Contains(_str, []byte{'e', '+'}) || bytes.Contains(_str, []byte{'E', '+'}) {
		nm := ToFloat(_str)
		return (int(nm))
	}
	num, err := strconv.ParseInt(string(_str), 10, 64)
	if err == nil {
		return int(num)
	}
	return 0
}

// StrToFloat converts string to float
func StrToFloat(_bsl string) float64 {
	if len(_bsl) <= 0 {
		return 0.0
	}
	f, _ := strconv.ParseFloat(_bsl, 64)
	return f
}

// StrToFloatAbs converts string to absolute float
func StrToFloatAbs(_bsl string) float64 {
	if len(_bsl) <= 0 {
		return 0.0
	}
	f, _ := strconv.ParseFloat(_bsl, 64)
	return math.Abs(f)
}

// TwoStrToFloat converts 2 strings to 2 floats
func TwoStrToFloat(_bsl1, _bsl2 string) (float64, float64) {
	f1, f2 := 0.0, 0.0
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseFloat(_bsl1, 64)
	}
	if len(_bsl2) > 0 {
		f2, _ = strconv.ParseFloat(_bsl2, 64)
	}
	return f1, f2
}

// ToFloat converts string to float, without default
func ToFloat(_bsl []byte) float64 {
	if len(_bsl) <= 0 {
		return 0.0
	}
	f, _ := strconv.ParseFloat(string(_bsl), 64)
	return f
}

// StrMultFloat returns the result as a float, of multiplying a string and a float
func StrMultFloat(_bsl1 string, _num float64) string {
	f1 := 0.0
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseFloat(_bsl1, 64)
	}
	return fmt.Sprintf("%f", f1*_num)
}

// StrAddFloat is shorthand
func StrAddFloat(_bsl1 string, _num float64) string {
	f1 := 0.0
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseFloat(_bsl1, 64)
	}
	return fmt.Sprintf("%f", f1+_num)
}

// StrAddInt is shorthand
func StrAddInt(_bsl1 string, _num int64) string {
	f1 := int64(0)
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseInt(_bsl1, 10, 64)
	}
	return fmt.Sprintf("%d", f1+_num)
}

// StrAddint is shorthand
func StrAddint(_bsl1 string, _num int) string {
	f1 := int64(0)
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseInt(_bsl1, 10, 64)
	}
	return fmt.Sprintf("%d", f1+int64(_num))
}

// StrDivFloat is shorthand
func StrDivFloat(_bsl1 string, _num float64) string {
	f1 := 0.0
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseFloat(_bsl1, 64)
	}
	return fmt.Sprintf("%f", f1/_num)
}

// StrAbsDivFloat is shorthand
func StrAbsDivFloat(_bsl1 string, _num float64) string {
	f1 := 0.0
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseFloat(_bsl1, 64)
	}
	return fmt.Sprintf("%f", math.Abs(f1/_num))
}

// StrInvert is shorthand
func StrInvert(_bsl1 string) string {
	f1 := 0.0
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseFloat(_bsl1, 64)
	}
	return fmt.Sprintf("%f", math.Abs(1.0/f1))
}

// StrNegate is shorthand
func StrNegate(_bsl1 string) string {
	f1 := 0.0
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseFloat(_bsl1, 64)
	}
	return fmt.Sprintf("%f", -f1)
}

// StrSignAsFloat is shorthand
func StrSignAsFloat(_bsl1 string, _snum float64) string {
	f1 := 0.0
	if len(_bsl1) > 0 {
		f1, _ = strconv.ParseFloat(_bsl1, 64)
	}
	switch _snum < 0.0 {
	case true:
		f1 = -math.Abs(f1)
	case false:
		f1 = math.Abs(f1)
	}
	return fmt.Sprintf("%f", f1)
}

// StrIntsAdd is shorthand
func StrIntsAdd(_bsl1, _bsl2 string) string {
	if len(_bsl1) <= 0 {
		return _bsl2
	}
	if len(_bsl2) <= 0 {
		return _bsl1
	}
	n1, _ := strconv.ParseInt(_bsl1, 10, 64)
	n2, _ := strconv.ParseInt(_bsl2, 10, 64)
	return fmt.Sprintf("%d", n1+n2)
}

// StrFloatsAdd is shorthand
func StrFloatsAdd(_bsl1, _bsl2 string) string {
	if len(_bsl1) <= 0 {
		return _bsl2
	}
	if len(_bsl2) <= 0 {
		return _bsl1
	}
	f1, _ := strconv.ParseFloat(_bsl1, 64)
	f2, _ := strconv.ParseFloat(_bsl2, 64)
	return fmt.Sprintf("%f", f1+f2)
}

// StrFloatsDiff is shorthand
func StrFloatsDiff(_bsl1, _bsl2 string) string {
	if len(_bsl1) <= 0 {
		return StrNegate(_bsl2)
	}
	if len(_bsl2) <= 0 {
		return _bsl1
	}
	f1, _ := strconv.ParseFloat(_bsl1, 64)
	f2, _ := strconv.ParseFloat(_bsl2, 64)
	return fmt.Sprintf("%f", f1-f2)
}

// StrFloatsAbsDiff is shorthand
func StrFloatsAbsDiff(_bsl1, _bsl2 string) string {
	if len(_bsl1) <= 0 {
		return StrAbs(_bsl2)
	}
	if len(_bsl2) <= 0 {
		return StrAbs(_bsl1)
	}
	f1, _ := strconv.ParseFloat(_bsl1, 64)
	f2, _ := strconv.ParseFloat(_bsl2, 64)
	return fmt.Sprintf("%f", math.Abs(f1-f2))
}

// StrFloatsMult is shorthand
func StrFloatsMult(_bsl1, _bsl2 string) string {
	if len(_bsl1) <= 0 {
		return _bsl2
	}
	if len(_bsl2) <= 0 {
		return _bsl1
	}
	f1, _ := strconv.ParseFloat(_bsl1, 64)
	f2, _ := strconv.ParseFloat(_bsl2, 64)
	return fmt.Sprintf("%f", f1*f2)
}

// StrFloatsMult3Zero returns 0 if any items are missing
func StrFloatsMult3Zero(_bsl1, _bsl2, _bsl3 string) string {
	if len(_bsl1) <= 0 {
		return "0.0"
	}
	if len(_bsl2) <= 0 {
		return "0.0"
	}
	if len(_bsl3) <= 0 {
		return "0.0"
	}
	f1, _ := strconv.ParseFloat(_bsl1, 64)
	f2, _ := strconv.ParseFloat(_bsl2, 64)
	f3, _ := strconv.ParseFloat(_bsl3, 64)
	return fmt.Sprintf("%f", f1*f2*f3)
}

// StrFloatsDiv is shorthand
func StrFloatsDiv(_bsl1, _bsl2, _def string) string {
	if len(_bsl1) <= 0 {
		return _def
	}
	if len(_bsl2) <= 0 {
		return _def
	}
	f2, _ := strconv.ParseFloat(_bsl2, 64)
	if math.Abs(f2) == 0.0 {
		return _def
	}
	f1, _ := strconv.ParseFloat(_bsl1, 64)
	return fmt.Sprintf("%f", f1/f2)
}

// StrFloatsAplusBminusC is shorthand
func StrFloatsAplusBminusC(_bsl1, _bsl2, _bsl3 string) string {
	var a, b, c float64
	if len(_bsl1) > 0 {
		a, _ = strconv.ParseFloat(_bsl1, 64)
	}
	if len(_bsl2) > 0 {
		b, _ = strconv.ParseFloat(_bsl2, 64)
	}
	if len(_bsl3) > 0 {
		c, _ = strconv.ParseFloat(_bsl3, 64)
	}
	return fmt.Sprintf("%f", a+b-c)
}

// StrFloatsAplusminusBminusC is shorthand
func StrFloatsAplusminusBminusC(_bsl1 string, _plusminus int64, _bsl2, _bsl3 string) string {
	var a, b, c float64
	if len(_bsl1) > 0 {
		a, _ = strconv.ParseFloat(_bsl1, 64)
	}
	if len(_bsl2) > 0 {
		b, _ = strconv.ParseFloat(_bsl2, 64)
	}
	if len(_bsl3) > 0 {
		c, _ = strconv.ParseFloat(_bsl3, 64)
	}
	return fmt.Sprintf("%f", a+float64(_plusminus)*(b-c))
}

// SliceFloatsAdd adds slice elements of the slice
func SliceFloatsAdd(_arr []float64) float64 {
	sum := 0.0
	for _, elt := range _arr {
		sum += elt
	}
	return sum
}

// IsDigit is shorthand
func IsDigit(_bb byte) bool {
	return ('0' <= _bb) && (_bb <= '9')
}

// IsYYYYMMDD is shorthand
func IsYYYYMMDD(_str string) bool {
	if len(_str) != 8 {
		return false
	}
	cc, _, mm, dd := _str[:2], _str[2:4], _str[4:6], _str[6:]
	switch cc {
	default:
		return false
	case "19", "20":
	}
	if !(IsDigit(_str[2]) && IsDigit(_str[3]) && IsDigit(_str[6]) && IsDigit(_str[7])) {
		return false
	}
	switch mm {
	default:
		return false
	case "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12":
	}
	if Toint0(dd) > 31 {
		return false
	}
	return true
}

// StryyyymmddLTEQ returns true if firstdate <= seconddate
// use AddCalDate if you want to compare offsetted dates
func StryyyymmddLTEQ(_dt1, _dt2 string) bool {
	len1, len2 := len(_dt1), len(_dt2)
	switch {
	case len1 < 8 && len2 == 8:
		return true
	case len1 == 8 && len2 < 8:
		return false
	case len1 == 8 && len2 == 8:
		for ii := 0; ii < 8; ii++ {
			switch {
			case _dt1[ii] > _dt2[ii]:
				return false
			case _dt1[ii] < _dt2[ii]:
				return true
			}
		}
		return true
	}
	pc, file, line, ok := runtime.Caller(1)
	fmt.Println("genutil.StryyyymmddLTEQ : bad dates dt1=", _dt1, " dt2=", _dt2, " callerFile=", file, " callerLine=", line, " pc=", pc, " ok=", ok)
	return false
}

// StryyyymmddLT returns true if firstdate < seconddate
// use AddCalDate if you want to compare offsetted dates
func StryyyymmddLT(_dt1, _dt2 string) bool {
	len1, len2 := len(_dt1), len(_dt2)
	switch {
	case len1 < 8 && len2 == 8:
		return true
	case len1 == 8 && len2 < 8:
		return false
	case len1 == 8 && len2 == 8:
		if _dt1 == _dt2 {
			return false
		}
		for ii := 0; ii < 8; ii++ {
			switch {
			case _dt1[ii] > _dt2[ii]:
				return false
			case _dt1[ii] < _dt2[ii]:
				return true
			}
		}
		return true
	}
	pc, file, line, ok := runtime.Caller(1)
	fmt.Println("genutil.StryyyymmddLT : bad dates dt1=", _dt1, " dt2=", _dt2, " callerFile=", file, " callerLine=", line, " pc=", pc, " ok=", ok)
	return false
}

// StryyyymmddInRange checks if STARTDATE <= yyyymmdd <= ENDDATE
func StryyyymmddInRange(_yyyymmdd, _startdate, _enddate string) bool {
	if StryyyymmddLTEQ(_startdate, _yyyymmdd) && StryyyymmddLTEQ(_yyyymmdd, _enddate) {
		return true
	}
	return false
}

// StryyyymmddInRangeOpenOpen checks if STARTDATE < yyyymmdd < ENDDATE
func StryyyymmddInRangeOpenOpen(_yyyymmdd, _startdate, _enddate string) bool {
	if StryyyymmddLT(_startdate, _yyyymmdd) && StryyyymmddLT(_yyyymmdd, _enddate) {
		return true
	}
	return false
}

// StryyyymmddInRangeClosedOpen checks if STARTDATE <= yyyymmdd < ENDDATE
func StryyyymmddInRangeClosedOpen(_yyyymmdd, _startdate, _enddate string) bool {
	if StryyyymmddLTEQ(_startdate, _yyyymmdd) && StryyyymmddLT(_yyyymmdd, _enddate) {
		return true
	}
	return false
}

// StryyyymmddInRangeOpenClosed checks if STARTDATE < yyyymmdd <= ENDDATE
func StryyyymmddInRangeOpenClosed(_yyyymmdd, _startdate, _enddate string) bool {
	if StryyyymmddLT(_startdate, _yyyymmdd) && StryyyymmddLTEQ(_yyyymmdd, _enddate) {
		return true
	}
	return false
}

// YYYY_MM_DD2yyyymmdd converts bytestring date to int64
func YYYY_MM_DD2yyyymmdd(_bsl []byte) int64 {
	if len(_bsl) < 10 {
		return 19010101
	}
	yyyy := ToInt(string(_bsl[0:4]), 1901)
	mm := ToInt(string(_bsl[5:7]), 0)
	dd := ToInt(string(_bsl[8:10]), 0)
	return yyyy*10000 + mm*100 + dd
}

//
// MMDDYYYY2yyyymmdd converts bytestring date of 02/03/2014 format to int64
func MMDDYYYY2yyyymmdd(_bsl []byte) int64 {
	if len(_bsl) < 10 {
		return 19010101
	}
	mm := ToInt(string(_bsl[0:2]), 0)
	dd := ToInt(string(_bsl[3:5]), 0)
	yyyy := ToInt(string(_bsl[6:10]), 1901)

	return yyyy*10000 + mm*100 + dd

}

// StrYYYY_MM_DD2yyyymmdd removes underscores or spaces from date
func StrYYYY_MM_DD2yyyymmdd(_dt string) string {
	if len(_dt) < 10 {
		return "19010101"
	}
	return _dt[0:4] + _dt[5:7] + _dt[8:10]
}

// YYYY_MM_DD_HH_MM_SS2yyyymmdd_hhmmss converts "2020-01-09 16:45:07" format dates to (YYYYMMDD, HHMMSS) string pair
func YYYY_MM_DD_HH_MM_SS2yyyymmdd_hhmmss(_bsl []byte) (int64, int64) {
	if len(_bsl) < 19 {
		return 19010101, -1
	}
	yyyy := ToInt(string(_bsl[0:4]), 1901)
	mm := ToInt(string(_bsl[5:7]), 0)
	dd := ToInt(string(_bsl[8:10]), 0)
	hh := ToInt(string(_bsl[11:13]), 1901)
	MM := ToInt(string(_bsl[14:16]), 0)
	ss := ToInt(string(_bsl[17:19]), 0)
	return yyyy*10000 + mm*100 + dd, hh*10000 + MM*100 + ss
}

// YYYY_MM_DD_HH_MM_SS_mmm_zz2yyyymmdd_hhmmss_mmm_zz converts "2020-01-09 16:45:07.mmm-zz" format dates to (YYYYMMDD, HHMMSS, mmm, zz) string pair
// Here zz is timezone from pgsql in hours from GMT
func YYYY_MM_DD_HH_MM_SS_mmm_zz2yyyymmdd_hhmmss_mmm_zz(_bsl []byte) (int64, int64, int64, int64) {
	if len(_bsl) < 26 {
		return 19010101, -1, -1, -1
	}
	yyyy := ToInt(string(_bsl[0:4]), 1901)
	mm := ToInt(string(_bsl[5:7]), 0)
	dd := ToInt(string(_bsl[8:10]), 0)
	hh := ToInt(string(_bsl[11:13]), 1901)
	MM := ToInt(string(_bsl[14:16]), 0)
	ss := ToInt(string(_bsl[17:19]), 0)
	mmm := ToInt(string(_bsl[21:23]), 0)
	zz := ToInt(string(_bsl[23:25]), 0) // could be signed
	return yyyy*10000 + mm*100 + dd, hh*10000 + MM*100 + ss, mmm, zz
}

// Hhmmss2Seconds converts to (possibly fractional) seconds
func Hhmmss2Seconds(_hhmmss string) float64 {
	hh := ToInt(_hhmmss[0:2], 0)
	mm := ToInt(_hhmmss[2:4], 0)
	ss := ToInt(_hhmmss[4:6], 0)
	hhmmss := float64(hh*3600.0 + mm*60.0 + ss)
	if (len(_hhmmss) > 6) && (_hhmmss[6] == '.') {
		hhmmss += StrToFloat("0" + _hhmmss[6:])
	}
	return hhmmss
}

// Hhmmss2Timetz converts specified HHMMSS time to today in the specified timezone, return in time.Time
// It returns false if tz is invalid
func Hhmmss2Timetz(_localTime, _timezone string) (time.Time, bool) {
	location, err := time.LoadLocation(_timezone)
	if err != nil {
		return time.Now(), false
	}
	nowTZ := time.Now().In(location)
	yyyy, mo, dd := nowTZ.Year(), nowTZ.Month(), nowTZ.Day()

	hh, _ := strconv.ParseInt(_localTime[0:2], 10, 64)
	mm, _ := strconv.ParseInt(_localTime[2:4], 10, 64)
	ss, _ := strconv.ParseInt(_localTime[4:6], 10, 64)
	tLoc := time.Date(int(yyyy), mo, int(dd), int(hh), int(mm), int(ss), 0, location)
	return tLoc, true
}

// Timetz2Timetz convert input time to the specified timezone
func Timetz2Timetz(_time time.Time, _timezone string) time.Time {
	location, err := time.LoadLocation(_timezone)
	if err != nil {
		return time.Now()
	}
	locTime := _time.In(location)
	return locTime
}

// Hhmmsstz2Timetz converts HHMMSS time in one timezone (fromTZ) to another (toTZ)
func Hhmmsstz2Timetz(_time, _fromTZ, _toTZ string) (string, bool) {
	lcltime, okl := Hhmmss2Timetz(_time, _fromTZ)
	if !okl {
		return "", false
	}
	regionTime := Timetz2Timetz(lcltime, _toTZ)
	return regionTime.Format("150405"), true
}

// DD_MMM_YY2yyyymmdd converts date format
func DD_MMM_YY2yyyymmdd(_bsl []byte) int64 {
	if len(_bsl) < 9 {
		return 20990101
	}
	dd := int((_bsl[0]-'0')*10 + (_bsl[1] - '0')) // fmt.Println("dd is ", dd)
	yy := int((_bsl[7]-'0')*10 + (_bsl[8] - '0')) // fmt.Println("yy is ", yy)
	mm := 0
	// JAN FEB MAR APR MAY JUN JUL AUG SEP OCT NOV DEC
	switch _bsl[3] {
	case 'J': // JAN, JUN, JUL
		switch _bsl[4] {
		case 'A':
			mm = 1 // JAN
		case 'U':
			switch _bsl[5] {
			case 'N':
				mm = 6 // JUN
			case 'L':
				mm = 7 // JUL
			}
		}
	case 'F':
		mm = 2 // FEB
	case 'M': // MAR, MAY
		switch _bsl[5] {
		case 'R':
			mm = 3 // MAR
		case 'Y':
			mm = 5 // MAY
		}
	case 'A': // APR, AUG
		switch _bsl[4] {
		case 'P':
			mm = 4 // APR
		case 'U':
			mm = 8 // AUG
		}
	case 'S':
		mm = 9 // SEP
	case 'O':
		mm = 10 // OCT
	case 'N':
		mm = 11 // NOV
	case 'D':
		mm = 12 // DEC
	}
	cc := 1900
	if yy < 30 {
		cc += 100
	}
	yyyymmdd := dd + mm*100 + 10000*(cc+yy)
	// fmt.Println("bsl= ", _bsl, " yy=", yy, " mm=", mm, " dd=", dd, " yyyymmdd=", yyyymmdd)
	return int64(yyyymmdd)
}

// SplitYYYYMMDD splits date into parts
func SplitYYYYMMDD(_yyyymmdd int64) (yyyy, mm, dd int64) {
	yyyy = int64(float64(_yyyymmdd) / 10000.0)
	mm = int64(float64(_yyyymmdd-10000*yyyy) / 100.0)
	dd = _yyyymmdd - 10000*yyyy - 100*mm
	return
}

// Yyyymmdd2SimpleJulian_Since_1900 returns simple julian of input date
func Yyyymmdd2SimpleJulian_Since_1900(_yyyymmdd int64) int64 {
	yyyy, mm, dd := SplitYYYYMMDD(_yyyymmdd)
	//return (yyyy - 1900)*365 + (mm - 1) * 31 + dd

	days_for_leap := int64((yyyy - 1900) / 4.0)

	days_to_mon := int64(0)
	ii := int64(1)
	days_in_mon := map[int64]int64{1: 31, 2: 28, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31}
	for ; ii < mm; ii++ {
		days_to_mon += days_in_mon[ii]
	}

	return (yyyy-1900)*365 + days_for_leap + days_to_mon + dd
}

// YYYYMMDD2M_D_YYYY is shorthand
func YYYYMMDD2M_D_YYYY(_dt string) string {
	return fmt.Sprintf("%d/%d/%s", ToInt(_dt[4:6], 0), ToInt(_dt[6:], 0), _dt[:4])
}

// MMM2MM is shorthand
func MMM2MM(_MMM string) string {
	switch _MMM {
	case "JAN", "Jan":
		return "01"
	case "FEB", "Feb":
		return "02"
	case "MAR", "Mar":
		return "03"
	case "APR", "Apr":
		return "04"
	case "MAY", "May":
		return "05"
	case "JUN", "Jun":
		return "06"
	case "JUL", "Jul":
		return "07"
	case "AUG", "Aug":
		return "08"
	case "SEP", "Sep":
		return "09"
	case "OCT", "Oct":
		return "10"
	case "NOV", "Nov":
		return "11"
	case "DEC", "Dec":
		return "12"
	}
	panic("genutil.MMM2MM: unknown 3-letter uppercase MONTH=" + _MMM)
	return ""
}

// MMSlashDDSlashYYYY2YYYYMMDD is shorthand
// convert 11/27/2013 (MM/DD/YYYY) --> 20131127 (YYYYMMDD)
// convert 2/9/2013   (M/D/YYYY)   --> 20130209 (YYYYMMDD)
func MMSlashDDSlashYYYY2YYYYMMDD(_dt string) string {
	lendt := len(_dt)
	if lendt == 0 {
		return ""
	}
	parts := strings.Split(_dt, "/")
	if len(parts) != 3 {
		return ""
	}
	if len(parts[0]) == 1 {
		parts[0] = "0" + parts[0]
	}
	if len(parts[1]) == 1 {
		parts[1] = "0" + parts[1]
	}
	return parts[2] + parts[0] + parts[1]
}

// DDSlashMMSlashYYYY2YYYYMMDD is shorthand
// convert 27/11/2013 (DD/MM/YYYY) --> 20131127 (YYYYMMDD)
// convert 9/2/2013   (D/M/YYYY)   --> 20130209 (YYYYMMDD)
func DDSlashMMSlashYYYY2YYYYMMDD(_dt string) string {
	lendt := len(_dt)
	if lendt == 0 {
		return ""
	}
	parts := strings.Split(_dt, "/")
	if len(parts) != 3 {
		return ""
	}
	if len(parts[0]) == 1 {
		parts[0] = "0" + parts[0]
	}
	if len(parts[1]) == 1 {
		parts[1] = "0" + parts[1]
	}
	return parts[2] + parts[1] + parts[0]
}

// DDDashMMDashYY2YYYYMMDD is shorthand
// convert 9-FEB-13 (DD-MMM-YY)   --> 20130209 (YYYYMMDD)
func DDDashMMDashYY2YYYYMMDD(_dt string) string {
	lendt := len(_dt)
	if lendt == 0 {
		return ""
	}
	parts := strings.Split(_dt, "-")
	if len(parts) != 3 {
		return ""
	}
	if len(parts[0]) == 1 {
		parts[0] = "0" + parts[0]
	}
	if len(parts[1]) != 3 {
		return ""
	}
	parts[1] = MMM2MM(parts[1])
	if len(parts[2]) == 2 {
		year, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return ""
		}
		if year < 30 {
			parts[2] = "20" + parts[2]
		} else {
			parts[2] = "19" + parts[2]
		}
	}
	return parts[2] + parts[1] + parts[0]
}

// Date2YYYYMMDD converts a date (by guessing) from one of several formats to YYYYMMDD
func Date2YYYYMMDD(_today, _dt string) string {
	lendt := len(_dt)
	if lendt == 0 {
		return ""
	}
	yynow, _ /*mmnow*/, _ /*ddnow*/ := _today[2:4], _today[4:6], _today[6:8]

	switch lendt {
	case 11:
		parts := strings.Split(_dt, "-")
		if (len(parts) == 3) && (len(parts[0]) == 3) && (len(parts[1]) == 2) && (len(parts[2]) == 4) {
			return parts[2] + MMM2MM(parts[0]) + parts[1]
		}
		if (len(parts) == 3) && (len(parts[0]) == 2) && (len(parts[1]) == 3) && (len(parts[2]) == 4) {
			return parts[2] + MMM2MM(parts[1]) + parts[0]
		}
		parts = strings.Split(_dt, "/")
		if (len(parts) == 3) && (len(parts[0]) == 3) && (len(parts[1]) == 2) && (len(parts[2]) == 4) {
			return parts[2] + MMM2MM(parts[0]) + parts[1]
		}
	case 10:
		parts := strings.Split(_dt, "-")
		if (len(parts) == 3) && (len(parts[0]) == 4) && (len(parts[1]) == 2) && (len(parts[2]) == 2) {
			return _dt[0:4] + _dt[5:7] + _dt[8:10]
		}
		parts = strings.Split(_dt, "/")
		if (len(parts) == 3) && (len(parts[0]) == 4) && (len(parts[1]) == 2) && (len(parts[2]) == 2) {
			return _dt[0:4] + _dt[5:7] + _dt[8:10]
		}
		if (len(parts) == 3) && (len(parts[0]) == 2) && (len(parts[1]) == 2) && (len(parts[2]) == 4) {
			if false {
				mmddyyyy := parts[2] + parts[0] + parts[1]
				ddmmyyyy := parts[2] + parts[1] + parts[0]
				diff0 := ToInt(mmddyyyy, 0) - ToInt(_today, 0)
				diff1 := ToInt(mmddyyyy, 0) - ToInt(_today, 0)
				return StrTernary(AbsInt64(diff0) < AbsInt64(diff1), mmddyyyy, ddmmyyyy)
			}
		}
	case 8:
		if strings.IndexAny(_dt, "/-") < 0 {
			return _dt
		} // YYYYMMDD already
		parts := strings.Split(_dt, "/")
		if (len(parts) == 3) && (len(parts[0]) == 2) && (len(parts[1]) == 2) && (len(parts[2]) == 2) {
			if _dt[6:8] == yynow {
				return "20" + _dt[6:8] + _dt[0:2] + _dt[2:4]
			}
		}
		parts = strings.Split(_dt, "-")
		if (len(parts) == 3) && (len(parts[0]) == 2) && (len(parts[1]) == 2) && (len(parts[2]) == 2) {
			if _dt[6:8] == yynow {
				return "20" + _dt[6:8] + _dt[0:2] + _dt[2:4]
			}
		}
	case 6:
		if strings.IndexAny(_dt, "/-") < 0 {
			return ""
		} //
		parts := strings.Split(_dt, "/")
		if (len(parts) == 3) && (len(parts[0]) == 1) && (len(parts[1]) == 1) && (len(parts[2]) == 2) {
			if _dt[4:6] == yynow {
				return fmt.Sprintf("20%s0%s0%s", _dt[4:6], _dt[0:1], _dt[1:2])
			}
		}
		parts = strings.Split(_dt, "-")
		if (len(parts) == 3) && (len(parts[0]) == 1) && (len(parts[1]) == 1) && (len(parts[2]) == 2) {
			if _dt[4:6] == yynow {
				return fmt.Sprintf("20%s0%s0%s", _dt[4:6], _dt[0:1], _dt[1:2])
			}
		}
	case 7:
		if strings.IndexAny(_dt, "/-") < 0 {
			return ""
		} //
		parts := strings.Split(_dt, "/")
		if (len(parts) == 3) && (len(parts[0]) == 1) && (len(parts[1]) == 2) && (len(parts[2]) == 2) {
			if parts[2] == yynow {
				return fmt.Sprintf("20%s0%s%s", parts[2], parts[0], parts[1])
			}
		}
		if (len(parts) == 3) && (len(parts[0]) == 2) && (len(parts[1]) == 1) && (len(parts[2]) == 2) {
			if parts[2] == yynow {
				return fmt.Sprintf("20%s%s0%s", parts[2], parts[0], parts[1])
			}
		}
		parts = strings.Split(_dt, "-")
		if (len(parts) == 3) && (len(parts[0]) == 1) && (len(parts[1]) == 2) && (len(parts[2]) == 2) {
			if parts[2] == yynow {
				return fmt.Sprintf("20%s0%s%s", parts[2], parts[0], parts[1])
			}
		}
		if (len(parts) == 3) && (len(parts[0]) == 2) && (len(parts[1]) == 1) && (len(parts[2]) == 2) {
			if parts[2] == yynow {
				return fmt.Sprintf("20%s%s0%s", parts[2], parts[0], parts[1])
			}
		}
	case 0:
		return ""
	}
	panic("genutil.Date2YYYYMMDD: could not parse date(" + _dt + ") : " + CallerInfo2())
	return ""
}

func FilenameExpandUser(_fname string) string {
	switch {
	case len(_fname) < 2:
		return _fname
	case _fname[:2] != "~/":
		return _fname
	}
	usr, _ := user.Current()
	dir := usr.HomeDir
	return strings.Replace(_fname, "~/", dir+"/", 1)
}

// PathExists returns whether the given file or directory exists or not
func PathExists(_path string) (bool, error) {
	_, err := os.Stat(_path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// PathIsDir checks if path is dir
func PathIsDir(_path string) bool {
	stat, err := os.Stat(_path)
	if os.IsNotExist(err) {
		return false
	}
	if err == nil {
		return stat.IsDir()
	}
	return false
}

// PathOK is shorthand
func PathOK(_path string) bool {
	_, err := os.Stat(_path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// AnyPathOK checks for a readable file in one of the supported compression formats
func AnyPathOK(_path string) bool {
	_, _, ofcode := ReadableFilename(_path)
	if ofcode == 0 {
		return false
	}
	return true
}

// FileExecutable returns whether the given file or directory exists or not, and is executable
func FileExecutable(_path string) bool {
	info, err := os.Stat(_path)
	if err != nil {
		return false
	}
	mode := info.Mode()
	isx := (mode & 0111) > 0
	return isx
}

// GetParentDir returns parent dir of supplied dir, setting bool to false if no parent
func GetParentDir(_dir string) (parent string, ok bool) {
	if parent == "/" {
		return "", false
	}
	parent = filepath.Dir(_dir)
	if parent == "." {
		return "", false
	}
	return parent, true
}

// ResolveDir returns Directory after resolving multiple dots, for example ..../foo will return absolute path in great-grand-parent dir
// Optionally it will check if the directory exists
// Return empty string on any error
func ResolveDir(_dir string, doCheck bool) string {
	if len(_dir) < 0 {
		return _dir
	}
	dir := filepath.Clean(_dir)
	if len(_dir) < 0 {
		return _dir
	}
	lessDots := strings.TrimLeft(dir, ".")
	if len(lessDots)+2 >= len(dir) {
		dd, _ := filepath.Abs(dir) // convert to abs
		if doCheck {
			if dd[:1] != "/" {
				return ""
			} // panic("Path should be absolute")
			if !PathIsDir(dd) {
				return ""
			} // log.Panicf("Path=%s does not exist\n", opt.Path)
		}
		return dd
	}
	// user provided too many leading dots, will need to handle them differently
	dots := dir[:len(_dir)-len(lessDots)]
	// fmt.Println("dots=", dots)
	dd, err := os.Getwd()
	if err != nil {
		return ""
	}
	ok := true
	for ii := 1; ii < len(dots); ii++ {
		dd, ok = GetParentDir(dd)
		if !ok {
			return ""
		}
	}
	// fmt.Println("dd=", dd)
	dd += lessDots
	if doCheck {
		if dd[:1] != "/" {
			return ""
		} // panic("Path should be absolute")
		if !PathIsDir(dd) {
			return ""
		} // log.Panicf("Path=%s does not exist\n", opt.Path)
	}
	return dd
}

// AllDirs returns array of all parental dirs.   You should provide it a good dir
func AllDirs(_dir string) (paths, dirs []string) {
	paths = []string{}
	dirs = []string{}
	path, dir := _dir, _dir
	paths = append(paths, path)
	// fmt.Println(path)
	for len(path) > 0 {
		if filepath.Dir(path) == path {
			break
		}
		path = filepath.Dir(path)
		paths = append(paths, path)
		dirs = append(dirs, StrTernary(len(path) > 1, dir[(1+len(path)):], dir[1:]))
		dir = path
	}
	return
}

// ReadableFilename returns information for subsequent reading of the specified file
// If not found, it looks for compression variants of the file
func ReadableFilename(_fname string) (ofname string, ofcmd *exec.Cmd, ofcode int) {
	ofname = "/dev/null"
	// ofcmd = nil
	ofcode = 0

	// ================================================================================
	// First extract the file exactly as the user specified it
	// ================================================================================
	fok := PathOK(_fname)
	switch {
	case strings.HasSuffix(_fname, ".xz") && fok:
		ofname = _fname
		ofcmd = exec.Command("/usr/bin/xzcat", _fname)
		ofcode = 1
		return
	case strings.HasSuffix(_fname, ".gz") && fok:
		ofname = _fname
		ofcmd = exec.Command("/bin/zcat", _fname)
		ofcode = 2
		return
	case strings.HasSuffix(_fname, ".bz2") && fok:
		ofname = _fname
		ofcmd = exec.Command("/usr/bin/bzcat", _fname)
		ofcode = 3
		return
	case strings.HasSuffix(_fname, ".zip") && fok:
		ofname = _fname
		ofcmd = exec.Command("/usr/bin/unzip -p", _fname)
		ofcode = 4
		return
	case strings.HasSuffix(_fname, ".bash"):
		if FileExecutable(_fname) {
			ofcmd = exec.Command(_fname)
			ofcode = 5
		} else {
			panic("bash file exists without execute permissions:" + _fname)
		}
		return
	case fok:
		ofname = _fname
		ofcmd = exec.Command("/bin/cat", _fname)
		ofcode = 6
		return
	}

	// ================================================================================
	// Next extract the file in a different format, but prefer .xz
	// ================================================================================
	var tmpf string
	switch {
	case strings.HasSuffix(_fname, ".xz"):
		tmpf = _fname[:len(_fname)-3]
	case strings.HasSuffix(_fname, ".gz"):
		tmpf = _fname[:len(_fname)-3]
	case strings.HasSuffix(_fname, ".bz2"):
		tmpf = _fname[:len(_fname)-4]
	case strings.HasSuffix(_fname, ".zip"):
		tmpf = _fname[:len(_fname)-4]
	default:
		tmpf = _fname
	}

	if PathOK(tmpf + ".xz") {
		ofname = tmpf + ".xz"
		ofcmd = exec.Command("/usr/bin/xzcat", ofname)
		ofcode = 7
		return
	}
	if PathOK(tmpf + ".gz") {
		ofname = tmpf + ".gz"
		ofcmd = exec.Command("/bin/zcat", ofname)
		ofcode = 8
		return
	}
	if PathOK(tmpf + ".bz2") {
		ofname = tmpf + ".bz2"
		ofcmd = exec.Command("/usr/bin/bzcat", ofname)
		ofcode = 9
		return
	}
	if PathOK(tmpf + ".zip") {
		ofname = tmpf + ".zip"
		ofcmd = exec.Command("/usr/bin/unzip -p", ofname)
		ofcode = 10
		return
	}
	if PathOK(tmpf) {
		ofname = tmpf
		ofcmd = exec.Command("/bin/cat", ofname)
		ofcode = 11
		return
	}
	return
}

// PathRemoveOrPanic panics if it fails to remove a directory
func PathRemoveOrPanic(_fname string) bool {
	err := os.Remove(_fname)
	if err != nil {
		panic(err)
	}
	return true
}

// WritableFilename returns information for subsequent writing of the specified file
// Any compression variants of the file are removed.
func WritableFilename(_fname string) (ofname string, ofcode int) {
	ofname = "/dev/null"
	ofcode = 0

	// ================================================================================
	// First remove any file exactly as the user specified it
	// ================================================================================
	fok := PathOK(_fname)
	switch {
	case strings.HasSuffix(_fname, ".xz") && fok:
		ofname, _, ofcode = _fname, PathRemoveOrPanic(_fname), 1
		return
	case strings.HasSuffix(_fname, ".gz") && fok:
		ofname, _, ofcode = _fname, PathRemoveOrPanic(_fname), 2
		return
	case strings.HasSuffix(_fname, ".bz2") && fok:
		ofname, _, ofcode = _fname, PathRemoveOrPanic(_fname), 3
		return
	case strings.HasSuffix(_fname, ".zip") && fok:
		ofname, _, ofcode = _fname, PathRemoveOrPanic(_fname), 4
		return
	case fok:
		ofname, _, ofcode = _fname, PathRemoveOrPanic(_fname), 6
		return
	}

	// ================================================================================
	// Next remove any variant of the file
	// ================================================================================
	tmpf := ""
	switch {
	case strings.HasSuffix(_fname, ".xz"):
		tmpf = _fname[:len(_fname)-3]
	case strings.HasSuffix(_fname, ".gz"):
		tmpf = _fname[:len(_fname)-3]
	case strings.HasSuffix(_fname, ".bz2"):
		tmpf = _fname[:len(_fname)-4]
	case strings.HasSuffix(_fname, ".zip"):
		tmpf = _fname[:len(_fname)-4]
	default:
		tmpf = _fname
	}

	switch {
	case PathOK(tmpf + ".xz"):
		ofname, _, ofcode = tmpf+".xz", PathRemoveOrPanic(tmpf+".xz"), 7
		return
	case PathOK(tmpf + ".gz"):
		ofname, _, ofcode = tmpf+".gz", PathRemoveOrPanic(tmpf+".gz"), 8
		return
	case PathOK(tmpf + ".bz2"):
		ofname, _, ofcode = tmpf+".bz2", PathRemoveOrPanic(tmpf+".bz2"), 9
		return
	case PathOK(tmpf + ".zip"):
		ofname, _, ofcode = tmpf+".zip", PathRemoveOrPanic(tmpf+".zip"), 10
		return
	case PathOK(tmpf):
		ofname, _, ofcode = tmpf, PathRemoveOrPanic(tmpf), 11
		return
	}
	return
}

// CompressType returns a numeric code based on the compression type indicated in the filename
func CompressType(_fname string) int {
	switch {
	case strings.HasSuffix(_fname, ".xz"):
		return 1
	case strings.HasSuffix(_fname, ".gz"):
		return 2
	case strings.HasSuffix(_fname, ".bz2"):
		return 3
	case strings.HasSuffix(_fname, ".zip") || strings.HasSuffix(_fname, ".ZIP"):
		return 4
	case strings.HasSuffix(_fname, ".bash"):
		return 5
	case strings.HasSuffix(_fname, ".zip"):
		return 4
	}
	return 0
}

// CompressionBasename returns uncompressed filename of the input filename
func CompressionBasename(_fname string) string {
	nn := len(_fname)
	switch {
	case strings.HasSuffix(_fname, ".xz"):
		return CompressionBasename(_fname[:(nn - 3)])
	case strings.HasSuffix(_fname, ".gz"):
		return CompressionBasename(_fname[:(nn - 3)])
	case strings.HasSuffix(_fname, ".bz2"):
		return CompressionBasename(_fname[:(nn - 4)])
	case strings.HasSuffix(_fname, ".zip"):
		return CompressionBasename(_fname[:(nn - 4)])
	case strings.HasSuffix(_fname, ".ZIP"):
		return CompressionBasename(_fname[:(nn - 4)])
	}
	return _fname
}

// RemoveCompressionVariants removes all compression variants of the specified filename, optionally preserving the base filename
func RemoveCompressionVariants(_fname string, _keepbase bool) {
	fbase := CompressionBasename(_fname)
	for _, ext := range []string{"", ".xz", ".gz", ".bz2", ".zip", ".ZIP"} {
		if _keepbase && (ext == "") {
			continue
		}
		ff := fbase + ext
		fok := PathOK(ff)
		if !fok {
			continue
		}
		PathRemoveOrPanic(ff)
	}
}

// ReadableFilenameCommand returns the command portion of the output of ReadableFilename()
func ReadableFilenameCommand(_fname string) string {
	_, ofcmd, _ := ReadableFilename(_fname)
	if ofcmd == nil {
		return ""
	}
	ostr, space := "", ""
	for _, arg := range ofcmd.Args {
		ostr += space + arg
		space = " "
	}
	return ostr
}

// ReadableFilenameTimestamp returns the timestamp of the output of ReadableFilename()
func ReadableFilenameTimestamp(_fname string) string {
	fname, ofcmd, _ := ReadableFilename(_fname)
	if ofcmd == nil {
		return ""
	}
	stat, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return ""
	}
	return stat.ModTime().Format("Mon 20060102 15:04:05 MST") // modification time
}

// OpenAny returns buffered reader for the content of the specified file, or available compression variant
func OpenAny(_fname string) *bufio.Reader {
	ofname, ofcmd, ofcode := ReadableFilename(_fname)
	switch ofcode {
	case 1, 7, 4, 10, 5:
		fi, err := ofcmd.StdoutPipe()
		ofcmd.Start()
		if err != nil {
			log.Panicf("genutil.OpenAny: err(%s) fname(%s) ofcmd(%s) ofcode(%d)", err.Error(), _fname, ofcmd, ofcode)
		}
		// defer fi.Close()
		r := bufio.NewReaderSize(fi, 20*4096)
		return r
	case 2, 8:
		fi, err := os.Open(ofname)
		if err != nil {
			log.Panicf("genutil.OpenAny: err(%s) fname(%s) ofname(%s) ofcode(%d)", err.Error(), _fname, ofname, ofcode)
		}
		// defer fi.Close()
		gzr, err := gzip.NewReader(fi)
		r := bufio.NewReaderSize(gzr, 20*4096)
		return r
	case 3, 9:
		fi, err := os.Open(ofname)
		if err != nil {
			log.Panicf("genutil.OpenAny: err(%s) fname(%s) ofname(%s) ofcode(%d)", err.Error(), _fname, ofname, ofcode)
		}
		// defer fi.Close()
		bzr := bzip2.NewReader(fi)
		r := bufio.NewReaderSize(bzr, 20*4096)
		return r
	case 6, 11:
		fi, err := os.Open(ofname)
		if err != nil {
			log.Panicf("genutil.OpenAny: err(%s) fname(%s) ofname(%s) ofcode(%d)", err.Error(), _fname, ofname, ofcode)
		}
		// defer fi.Close()
		r := bufio.NewReaderSize(fi, 20*4096)
		return r
	default:
	}
	return nil
}

// OpenAnyIO returns unbuffered reader for the content of the specified file, or available compression variant
func OpenAnyIO(_fname string) *io.Reader {
	ofname, ofcmd, ofcode := ReadableFilename(_fname)
	switch ofcode {
	case 1, 7, 4, 10, 5:
		fi, err := ofcmd.StdoutPipe()
		ofcmd.Start()
		if err != nil {
			log.Panicf("genutil.OpenAny: err(%s) fname(%s) ofcmd(%s) ofcode(%d)", err.Error(), _fname, ofcmd, ofcode)
		}
		// defer fi.Close()
		r := io.Reader(fi)
		return &r
	case 2, 8:
		fi, err := os.Open(ofname)
		if err != nil {
			log.Panicf("genutil.OpenAnyIO: err(%s) fname(%s) ofname(%s) ofcode(%d)", err.Error(), _fname, ofname, ofcode)
		}
		// defer fi.Close()
		gzr, err := gzip.NewReader(fi)
		r := io.Reader(gzr)
		return &r
	case 3, 9:
		fi, err := os.Open(ofname)
		if err != nil {
			log.Panicf("genutil.OpenAnyIO: err(%s) fname(%s) ofname(%s) ofcode(%d)", err.Error(), _fname, ofname, ofcode)
		}
		// defer fi.Close()
		bzr := bzip2.NewReader(fi)
		r := io.Reader(bzr)
		return &r
	case 6, 11:
		fi, err := os.Open(ofname)
		if err != nil {
			log.Panicf("genutil.OpenAnyIO: err(%s) fname(%s) ofname(%s) ofcode(%d)", err.Error(), _fname, ofname, ofcode)
		}
		// defer fi.Close()
		r := io.Reader(fi)
		return &r
	default:
	}
	return nil
}

// OpenAnyErr returns buffered reader for the content of the specified file, or available compression variant
// It is more error conscious than OpenAny()
func OpenAnyErr(_fname string) (*bufio.Reader, error) {
	ofname, ofcmd, ofcode := ReadableFilename(_fname)
	if ofcmd == nil {
		return nil, errors.New("os.exec.Command returned nil pointer")
	}
	switch ofcode {
	case 1, 7, 4, 10, 5:
		fi, err := ofcmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = ofcmd.Start()
		if err != nil {
			return nil, err
		}
		// defer fi.Close()
		r := bufio.NewReaderSize(fi, 20*4096)
		return r, nil
	case 2, 8:
		fi, err := os.Open(ofname)
		if err != nil {
			return nil, err
		}
		// defer fi.Close()
		gzr, err := gzip.NewReader(fi)
		if err != nil {
			fi.Close()
			return nil, err
		}
		r := bufio.NewReaderSize(gzr, 20*4096)
		return r, nil
	case 3, 9:
		fi, err := os.Open(ofname)
		if err != nil {
			return nil, err
		}
		// defer fi.Close()
		bzr := bzip2.NewReader(fi)
		r := bufio.NewReaderSize(bzr, 20*4096)
		return r, nil
	case 6, 11:
		fi, err := os.Open(ofname)
		if err != nil {
			return nil, err
		}
		// defer fi.Close()
		r := bufio.NewReaderSize(fi, 20*4096)
		return r, nil
	default:
	}
	return nil, fmt.Errorf("OpenAnyErr : unknown ofcode = %d", ofcode)
}

// WriteStringToFile is shorthand
func WriteStringToFile(_str, _fname string) {
	fo, err := os.Create(_fname)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	io.WriteString(fo, _str)
}

// WriteStringToGzipFile is shorthand
func WriteStringToGzipFile(_str, _fname string) {
	fo, err := os.Create(_fname)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	ww0 := bufio.NewWriter(fo)
	defer ww0.Flush()
	ww := gzip.NewWriter(ww0)
	defer ww.Close()
	io.WriteString(ww, _str)
}

// NewBoolMap returns a map of string to true
func NewBoolMap() map[string]bool {
	aset := make(map[string]bool)
	return aset
}

// NewBoolMapFromCsv returns a map where each element of the supplied string is set true
func NewBoolMapFromCsv(_csv, _sep string) map[string]bool {
	aset := make(map[string]bool)
	parts := strings.Split(_csv, _sep)
	for _, part := range parts {
		str := strings.TrimSpace(part)
		if len(str) > 0 {
			aset[str] = true
		}
	}
	return aset
}

// UpdateBoolMapFromCsv updates the map setting elements of the string to true
func UpdateBoolMapFromCsv(_aset *map[string]bool, _csv, _sep string) {
	parts := strings.Split(_csv, _sep)
	for _, part := range parts {
		str := strings.TrimSpace(part)
		if len(str) > 0 {
			(*_aset)[str] = true
		}
	}
}

// UpdateBoolMap updates the map, setting elements of the slice to true
func UpdateBoolMap(_aset *map[string]bool, _keys []string) {
	for _, key := range _keys {
		str := strings.TrimSpace(key)
		if len(str) > 0 {
			(*_aset)[str] = true
		}
	}
}

// KeysBoolMap is shorthand
func KeysBoolMap(_aset *map[string]bool) []string {
	keys := []string{}
	for kk := range *_aset {
		keys = append(keys, kk)
	}
	return keys
}

// NewInt64BoolMap is shorthand
func NewInt64BoolMap() map[int64]bool {
	aset := make(map[int64]bool)
	return aset
}

// FileList returns files in dir
func FileList(_dname string) []string {
	flist, _ := ioutil.ReadDir(_dname)
	rlist := []string{}
	for _, finfo := range flist {
		rlist = append(rlist, finfo.Name())
	}
	return rlist
}

// SliceContainsStr checks if given string is in the slice
func SliceContainsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// SliceIndexStr calculates index of string in slice
func SliceIndexStr(s []string, e string) int {
	for ii, a := range s {
		if a == e {
			return ii
		}
	}
	return -1
}

// ListContainsStr checks if given string is in list
func ListContainsStr(_str string, _list ...string) bool {
	for _, ss := range _list {
		if ss == _str {
			return true
		}
	}
	return false
}

// ListContainsByte checks if given byte is in list
func ListContainsByte(_bb byte, _list ...byte) bool {
	for _, bb := range _list {
		if bb == _bb {
			return true
		}
	}
	return false
}

// MakeDirOrDie panics if unable to create the dir (or if it exists)
func MakeDirOrDie(_dirBase, _dirName string) string {
	if len(_dirBase) <= 0 {
		panic("genutil.MakeDirOrDie: empty dirBase")
	}
	if strings.HasSuffix(_dirBase, "/") {
		panic("genutil.MakeDirOrDie: dirBase should not end in /")
	}
	if len(_dirName) <= 0 {
		panic("genutil.MakeDirOrDie: empty dirName")
	}
	if !PathIsDir(_dirBase) {
		panic("genutil.MakeDirOrDie: dirBase is not a dir: " + _dirBase)
	}
	newpath := _dirBase + "/" + _dirName
	if PathOK(newpath) {
		panic("genutil.MakeDirOrDie: path already exists: " + newpath)
	}
	var perm os.FileMode = 0775
	if err := os.Mkdir(newpath, perm); err != nil {
		panic("genutil.MakeDirOrDie: error creating dir with 0775 perm : " + newpath)
	}
	return newpath
}

// EnsureDirOrDie dies if the dir did not exist and could not be created
func EnsureDirOrDie(_dirBase, _dirName string) string {
	if len(_dirBase) <= 0 {
		panic("genutil.EnsureDirOrDie: empty dirBase")
	}
	if strings.HasSuffix(_dirBase, "/") {
		panic("genutil.EnsureDirOrDie: dirBase should not end in /")
	}
	if len(_dirName) <= 0 {
		panic("genutil.EnsureDirOrDie: empty dirName")
	}
	if !PathIsDir(_dirBase) {
		panic("genutil.EnsureDirOrDie: dirBase is not a dir: " + _dirBase)
	}
	newpath := _dirBase + "/" + _dirName
	if PathIsDir(newpath) {
		return newpath
	}
	var perm os.FileMode = 0775
	if err := os.Mkdir(newpath, perm); err != nil {
		panic("genutil.EnsureDirOrDie: error creating dir with 0775 perm : " + newpath)
	}
	return newpath
}

// FileInfoSysStr formats file info into readable form
// call it with stat.Sys()
func FileInfoSysStr(_sys interface{}, _sep string) string {
	str := ""
	unixStat, ok := _sys.(*syscall.Stat_t)
	if ok {
		str += fmt.Sprintf("inumber=%d%s", unixStat.Ino, _sep)
		str += fmt.Sprintf("uid=%d%s", unixStat.Uid, _sep)
		str += fmt.Sprintf("gid=%d%s", unixStat.Gid, _sep)
		// str	+= fmt.Sprintf("Mtim=%d.%d%s", unixStat.Mtim.Sec,unixStat.Mtim.NSec, _sep)
		str += fmt.Sprintf("Nlink=%d", unixStat.Nlink) // Number of hard links
	}
	return str
}

// FileInfo formats file info into readable form
func FileInfo(_fname, _sep string, _fullinfo bool) string {
	stat, err := os.Stat(_fname)
	if os.IsNotExist(err) {
		return fmt.Sprintf("fname=%s%sstatus=notexists", _fname, _sep)
	}
	str := fmt.Sprintf("fname=%s%ssize=%d%smode=%s%smodtime=%s",
		_fname, _sep, stat.Size(), _sep, stat.Mode().String(), _sep, stat.ModTime().Format("Mon 20060102 15:04:05 MST"))
	if _fullinfo {
		str += fmt.Sprintf("%sname=%s%sisdir=%t%s%s", _sep, stat.Name(), _sep, stat.IsDir(), _sep, FileInfoSysStr(stat.Sys(), _sep))
	}
	return str
}

// FileSize returns -1 if file not found
func FileSize(_fname string) int {
	stat, err := os.Stat(_fname)
	if os.IsNotExist(err) {
		return -1
	}
	return int(stat.Size())
}

// CheckFileIsReadableAndNonzeroOrDie is shorthand
func CheckFileIsReadableAndNonzeroOrDie(_fname string) {
	stat, err := os.Stat(_fname)
	if os.IsNotExist(err) {
		panic("genutil.CheckFileIsReadableAndNonzeroOrDie: bad file: " + _fname)
	}
	perm := stat.Mode()
	if (perm & 0x0004) == 0x0000 {
		panic("genutil.CheckFileIsReadableAndNonzeroOrDie: bad filemode(" + string(perm) + ") for file:" + _fname)
	}
}

// BashExecOrDie executes the string cmd with /bin/bash and panics on any kind of failure
func BashExecOrDie(_verbose bool, _cmd, _dir string) string {
	if _verbose {
		fmt.Println("BashExecOrDie:info cmd is: (" + _cmd + ")")
	}
	if len(_cmd) < 0 {
		panic("genutil.BashExecOrDie: empty command")
	}
	cmd := exec.Command("/bin/bash", "-c", _cmd)
	cmd.Dir = _dir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic("genutil.BashExecOrDie: failed to get stdout pipe from command")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic("genutil.BashExecOrDie: failed to get stderr pipe from command")
	}
	err = cmd.Start()
	if err != nil {
		panic("genutil.BashExecOrDie: could not run the command")
	}
	buf, err := ioutil.ReadAll(stdout)
	if err != nil {
		panic("genutil.BashExecOrDie: could not retrieve output from command")
	}
	bufe, err := ioutil.ReadAll(stderr)
	if err != nil {
		panic("genutil.BashExecOrDie: could not retrieve error from command")
	}
	cmd.Wait()
	if (len(buf) > 0) && (buf[len(buf)-1] == '\n') {
		buf = buf[:len(buf)-1]
	}

	if len(bufe) <= 0 {
		return string(buf)
	}
	return string(buf) + "\n" + string(bufe)
}

// ExecCommandOrDie executes the given command and panics on any kind of failure
func ExecCommandOrDie(_verbose bool, _cmd string) {
	if _verbose {
		fmt.Println("ExecCommandOrDie:info cmd is: (" + _cmd + ")")
	}
	if len(_cmd) < 0 {
		panic("genutil.ExecCommandOrDie: empty command")
	}
	parts := strings.Split(_cmd, " ")
	if len(parts) < 1 {
		panic("genutil.ExecCommandOrDie: bad command (" + _cmd + ")")
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	if _verbose {
		fmt.Println("ExecCommandOrDie:info running cmd")
	}
	_, err := cmd.StdoutPipe()
	cmd.Run()
	if err != nil {
		panic("genutil.ExecCommandOrDie: command (" + _cmd + ") failed with error:" + fmt.Sprint("%v", err))
	}
	if _verbose {
		fmt.Println("ExecCommandOrDie:info done cmd")
	}
}

//
// IsZipFilename checks for any kind of .zip or .ZIP or .ZiP file
func IsZipFilename(_fname string) bool {
	fmt.Println("IsZipFilename=", _fname)
	flen := len(_fname)
	if flen < 5 {
		return false
	}
	suff := strings.ToLower(_fname[flen-4:])
	if suff == ".zip" {
		return true
	}
	return false
}

// ZipFirstFileInfo returns name, date, time, size
func ZipFirstFileInfo(_zipfile string, _verbose bool) (string, string, string, int) {
	cmd := fmt.Sprintf("/usr/bin/unzip -l %s", _zipfile)
	out := BashExecOrDie(false, cmd, ".")
	seenArchive, seenHeader, parts := false, false, []string{}
	for _, ln := range strings.Split(out, "\n") {
		switch {
		case (!seenArchive) && strings.Contains(ln, "Archive"):
			seenArchive = true
			continue
		case (!seenHeader) && strings.Contains(ln, "Length ") && strings.Contains(ln, " Name"):
			seenHeader = true
			continue
		case strings.HasPrefix(strings.Trim(ln, " "), "-"):
			continue
		}
		parts = strings.Fields(ln)
		if _verbose {
			fmt.Printf("ln=%s\n", ln)
		}
		break
	}
	if _verbose {
		fmt.Printf("parts=%s\n", strings.Join(parts, ","))
	}
	return parts[3], fmt.Sprintf("%d", MMDDYYYY2yyyymmdd([]byte(parts[1]))), parts[2], Toint0(parts[0])
}

// GetYyyymmddFromFilenameYymmddFromEndWithSuffixLen grabs the YYMMDD from filenames of form foo_YYMMDD.zip, but extend with the decade
func GetYyyymmddFromFilenameYymmddFromEndWithSuffixLen(_fname string, _suffLen int, _def string) string {
	flen := len(_fname)
	if flen < _suffLen+6 {
		return _def
	}
	yymmdd := _fname[flen-_suffLen-6 : flen-_suffLen]
	switch yymmdd[0:1] {
	case "0", "1", "2", "3":
		return "20" + yymmdd
	}
	return "19" + yymmdd
}

// GetYyyymmddFromFilenameYyyymmFromEndWithSuffixLen grab the YYMMDD from filenames of form foo_YYYYMM.zip, but extend with 1st day of month
func GetYyyymmddFromFilenameYyyymmFromEndWithSuffixLen(_fname string, _suffLen int, _def string) string {
	flen := len(_fname)
	if flen < _suffLen+6 {
		return _def
	}
	yyyymm := _fname[flen-_suffLen-6 : flen-_suffLen]
	return yyyymm + "01"
}

// SplitOrNull on empty input returns null slice, unlike plain strings.Split which will return 1 element slice
func SplitOrNull(_str, _sep string) []string {
	if _str == "" {
		return []string{}
	}
	parts := strings.Split(_str, _sep)
	return parts
}

// GetNumLines counts number of lines in any compression variant of file
func GetNumLines(_fname string) int64 {
	r := OpenAny(_fname)
	var err error
	var num int64
	for {
		_, _, err = r.ReadLine()
		if err != nil {
			break
		}
		num++
	}

	return num
}

// Today is shorthand
func Today() string { return fmt.Sprintf("%d", Time2YYYYMMDD(time.Now())) }

// Now is shorthand
func Now() string {
	now := time.Now()
	return fmt.Sprintf("%02d%02d%02d", now.Hour(), now.Minute(), now.Second())
}

// AddCalDate adds number of dates to specified date
func AddCalDate(_date string, _offset int) string {
	if len(_date) < 8 {
		return ""
	}
	yyyy, mm, dd := ToInt(_date[:4], 0), ToInt(_date[4:6], 0), ToInt(_date[6:], 0)
	dt := time.Date(int(yyyy), time.Month(mm), int(dd), 0, 0, 0, 0, time.Now().Location())
	newDate := dt.AddDate(0, 0, _offset)
	return fmt.Sprintf("%d", Time2YYYYMMDD(newDate))
}

// CalDatelist creates list of dates from the range, possibly including/excluding the begin/end dates
func CalDatelist(_begdate, _enddate string, _includeBeg, _includeEnd bool) []string {
	if len(_begdate) < 8 {
		panic(fmt.Sprintf("CalDatelist: bad begdate(%s)", _begdate))
	}
	if len(_enddate) < 8 {
		panic(fmt.Sprintf("CalDatelist: bad enddate(%s)", _enddate))
	}
	if !StryyyymmddLTEQ(_begdate, _enddate) {
		return []string{}
	}
	dts := []string{}
	yyyy0, mm0, dd0 := ToInt(_begdate[:4], 0), ToInt(_begdate[4:6], 0), ToInt(_begdate[6:], 0)
	dt0 := time.Date(int(yyyy0), time.Month(mm0), int(dd0), 0, 0, 0, 0, time.Now().Location())
	if _includeBeg {
		dts = append(dts, _begdate)
	}
	for {
		dt0 = dt0.AddDate(0, 0, 1)
		dt := fmt.Sprintf("%d", Time2YYYYMMDD(dt0))
		if !StryyyymmddLT(dt, _enddate) {
			break
		}
		dts = append(dts, dt)
	}
	if _includeEnd && StryyyymmddLT(_begdate, _enddate) {
		dts = append(dts, _enddate)
	}
	return dts
}

// TodayTZ returns today in specified timezone
func TodayTZ(_timezone string) string {
	location, err := time.LoadLocation(_timezone)
	if err != nil {
		panic(err)
	}
	todaytz := time.Now().In(location)
	return fmt.Sprintf("%d", Time2YYYYMMDD(todaytz))
}

// NowTZ returns today in specified timezone
func NowTZ(_timezone string) string {
	location, err := time.LoadLocation(_timezone)
	if err != nil {
		panic(err)
	}
	return time.Now().In(location).Format("150405")
}

// GetLastSunday returns the most recent sunday
func GetLastSunday(_timezone string) string {
	location, err := time.LoadLocation(_timezone)
	if err != nil {
		panic(err)
	}
	t1 := time.Now().In(location)
	if t1.Weekday() == time.Sunday {
		return t1.Format("20060102")
	}
	return t1.AddDate(0, 0, int(time.Sunday)-int(t1.Weekday())).Format("20060102")
}

// GetLogicalDate returns today. Or tomorrow if it is now past the specified time.
func GetLogicalDate(_timezone string, _time string) string {
	location, err := time.LoadLocation(_timezone)
	if err != nil {
		panic(err)
	}
	t1 := time.Now().In(location)
	t2, ok := Hhmmss2Timetz(_time, _timezone)
	if !ok {
		panic("invalid time")
	}
	if t1.After(t2) {
		return t1.AddDate(0, 0, 1).Format("20060102")
	}
	return t1.Format("20060102")
}

// DateParts is a utility to convert some date mnemonics
// Do not call this directly for NBD/PBD etc, although today might be sort of safe
func DateParts(_date string) (int, int) {
	print := false
	if print {
		fmt.Println("genutil.Dateparts: testing date:", _date)
	}

	if len(_date) <= 1 {
		return 0, 0
	}

	dt, offset := "", 0
	switch {
	case strings.Index(_date, "-") >= 0:
		parts := strings.Split(_date, "-")
		dt, offset = parts[0], -Toint0(parts[1])
	case strings.Index(_date, "+") >= 0:
		parts := strings.Split(_date, "+")
		dt, offset = parts[0], Toint0(parts[1])
	default:
		dt, offset = _date, 0
	}

	if print {
		fmt.Println("genutil.Dateparts: testing date:", _date, "dt=", dt, "offset=", offset)
	}

	switch strings.ToUpper(dt) {
	case "TD", "TODAY":
		return Time2YYYYMMDD(time.Now()), offset
	case "NBD":
		return Time2YYYYMMDD(time.Now()), offset + 1
	case "PBD":
		return Time2YYYYMMDD(time.Now()), offset - 1
	case "NNBD":
		return Time2YYYYMMDD(time.Now()), offset + 2
	case "PPBD":
		return Time2YYYYMMDD(time.Now()), offset - 2
	default:
		if print {
			fmt.Println("genutil.Dateparts: default dt=", strings.ToUpper(dt))
		}
	}

	return Toint0(dt), offset
}

// Time2YYYYMMDD converts time.Time to  date string YYYYMMDD
func Time2YYYYMMDD(_tt time.Time) int {
	yyyy, mo, dd := _tt.Date()
	dt := yyyy*10000 + int(mo)*100 + dd
	return dt
}

// FillDate instantiates the date pattern from specified date
func FillDate(pat string, ctime time.Time) string {

	YYYY := ctime.Format("2006")
	YY := ctime.Format("06")
	MM := ctime.Format("01")
	DD := ctime.Format("02")

	pat = strings.Replace(pat, "$YYYY", YYYY, -1)
	pat = strings.Replace(pat, "$YY", YY, -1)
	pat = strings.Replace(pat, "$MM", MM, -1)
	pat = strings.Replace(pat, "$DD", DD, -1)

	return pat
}

// SearchForFileWithPattern is shorthand
func SearchForFileWithPattern(pat string) (bool, string) {
	matched, err := filepath.Glob(pat)
	if matched != nil && err == nil {
		return true, matched[0]
	}
	return false, ""
}

// GetLatestDatedDir is shorthand
func GetLatestDatedDir(parentdir string) string {
	out := BashExecOrDie(false, fmt.Sprintf("ls -1t %s | grep [12][0-9][0-9][0-9] | head -1", parentdir), "/tmp/")
	out = strings.Trim(out, "\r\n\t ")
	return out
}

// GetLatestFileWithPattern is shorthand
func GetLatestFileWithPattern(pattern string) string {
	out := BashExecOrDie(false, fmt.Sprintf("ls -1t %s | head -1", pattern), "/tmp/")
	out = strings.Trim(out, "\r\n\t ")
	return out

}

// GetSecondLatestFileWithPattern is shorthand
func GetSecondLatestFileWithPattern(pattern string) string {
	out := BashExecOrDie(false, fmt.Sprintf("ls -1t %s | head -2 | tail -1", pattern), "/tmp/")
	out = strings.Trim(out, "\r\n\t ")
	return out
}

// PreviousYYYYMMDD is shorthand
func PreviousYYYYMMDD(_dt string, _num int) string {
	yyyy, mm, dd := ToInt(_dt[:4], 0), ToInt(_dt[4:6], 0), ToInt(_dt[6:], 0)
	dt := _dt
	for ii := 0; ii < _num; ii++ {
		dd--
		if dd == 0 {
			dd = 31
			mm--
		}
		if mm == 0 {
			mm = 12
			yyyy--
		}
		dt = fmt.Sprintf("%04d%02d%02d", yyyy, mm, dd)
	}
	return dt
}

// PreviousYYYYMM is shorthand
func PreviousYYYYMM(_dt string) string {
	yyyy, mm := ToInt(_dt[:4], 0), ToInt(_dt[4:6], 0)
	if mm == 1 {
		yyyy, mm = yyyy-1, 12
	} else {
		mm--
	}
	return fmt.Sprintf("%04d%02d", yyyy, mm)
}

// FileAsofPrevious replaces YYYYMMDD with older dates until it finds a readable file (any compression variant)
// Today is not considered
func FileAsofPrevious(_path, _dt string, _num int) string {
	yyyy, mm, dd := ToInt(_dt[:4], 0), ToInt(_dt[4:6], 0), ToInt(_dt[6:], 0)
	for ii := 0; ii < _num; ii++ {
		dd--
		if dd < 0 {
			dd = 31
			mm--
		}
		if mm <= 0 {
			mm = 12
			yyyy--
		}
		dt := fmt.Sprintf("%04d%02d%02d", yyyy, mm, dd)
		trypath := _path
		trypath = strings.Replace(trypath, "YYYYMMDD", dt, -1)
		ofname, _, ofcode := ReadableFilename(trypath)
		if ofcode != 0 {
			return ofname
		}
	}
	return ""
}

// FileAsofCurrent replaces YYYYMMDD with older dates until it finds a readable file (any compression variant)
// Today is considered
func FileAsofCurrent(_path, _dt string, _num int) string {
	if false {
		fmt.Println("genutil.FileAsofCurrent: _path=", _path, "dt=", _dt)
	}
	yyyy, mm, dd := ToInt(_dt[:4], 0), ToInt(_dt[4:6], 0), ToInt(_dt[6:], 0)
	for ii := 0; ii < _num; ii++ {
		dt := fmt.Sprintf("%04d%02d%02d", yyyy, mm, dd)
		trypath := _path
		trypath = strings.Replace(trypath, "YYYYMMDD", dt, -1)
		ofname, _, ofcode := ReadableFilename(trypath)
		if ofcode != 0 {
			return ofname
		}
		dd--
		if dd < 0 {
			dd = 31
			mm--
		}
		if mm <= 0 {
			mm = 12
			yyyy--
		}
	}
	return ""
}

// CallerInfo2 is used to embellish error messages with the caller name
func CallerInfo2() string {
	pc, file, line, ok := runtime.Caller(2)
	return fmt.Sprintf(" callerFile=%s callerLine=%d pc=%d ok=%t", file, line, pc, ok)
}

// FlipIfFalseStr helps compensate for golang not having ternary op a
func FlipIfFalseStr(_flipIfFalse bool, _val1, _val2 string) (string, string) {
	if _flipIfFalse {
		return _val1, _val2
	}
	return _val2, _val1
}

// FlipIfFalseInt helps compensate for golang not having ternary op a
func FlipIfFalseInt(_flipIfFalse bool, _val1, _val2 int) (int, int) {
	if _flipIfFalse {
		return _val1, _val2
	}
	return _val2, _val1
}

// FlipIfFalseInt64 helps compensate for golang not having ternary op a
func FlipIfFalseInt64(_flipIfFalse bool, _val1, _val2 int64) (int64, int64) {
	if _flipIfFalse {
		return _val1, _val2
	}
	return _val2, _val1
}

// FlipIfFalseFloat helps compensate for golang not having ternary op a
func FlipIfFalseFloat(_flipIfFalse bool, _val1, _val2 float64) (float64, float64) {
	if _flipIfFalse {
		return _val1, _val2
	}
	return _val2, _val1
}

// ================================================================================

// UpdateMaxId updates the max id
func UpdateMaxId(_idmap *map[string]int64, _kk, _newid string) int64 {
	newid := ToInt(_newid, 0)
	lastid, ok := (*_idmap)[_kk]
	if ok && (lastid >= newid) {
		return lastid
	}
	(*_idmap)[_kk] = newid
	return newid
}

// IncrementMaxId increments the max id
func IncrementMaxId(_idmap *map[string]int64, _kk string) int64 {
	lastid, ok := (*_idmap)[_kk]
	if !ok {
		lastid = -1
	}
	lastid++
	(*_idmap)[_kk] = lastid
	return lastid
}

//================================================================================

// StrAorB is shorthand
func StrAorB(_a, _b string) string {
	if len(_a) > 0 {
		return _a
	}
	return _b
}

// StrTernary is shorthand for the missing golang string ternary operatory
func StrTernary(_aIfTrue bool, _a, _b string) string {
	if _aIfTrue {
		return _a
	}
	return _b
}

// FloatTernary is shorthand
func FloatTernary(_aIfTrue bool, _a, _b float64) float64 {
	if _aIfTrue {
		return _a
	}
	return _b
}

// IntTernary is shorthand
func IntTernary(_aIfTrue bool, _a, _b int) int {
	if _aIfTrue {
		return _a
	}
	return _b
}

// Int64Ternary is shorthand
func Int64Ternary(_aIfTrue bool, _a, _b int64) int64 {
	if _aIfTrue {
		return _a
	}
	return _b
}

// EmptyIfZero returns empty string or the currency amount if nonzero
func EmptyIfZero(_num, _ccy string) string {
	num := StrToFloat(_num)
	if math.Abs(num) <= 0 {
		return ""
	}
	return fmt.Sprintf("%s%f", _ccy, num)
}

// StrNonzeroAorB returns empty string or one of the currency amounts if nonzero
func StrNonzeroAorB(_a, _accy, _b, _bccy string) string {
	lena, lenb := len(_a), len(_b)
	if lena+lenb == 0 {
		return ""
	}
	if lena <= 0 {
		return EmptyIfZero(_b, _bccy)
	}
	if lenb <= 0 {
		return EmptyIfZero(_a, _accy)
	}
	a, b := StrToFloat(_a), StrToFloat(_b)
	if math.Abs(a)+math.Abs(b) == 0.0 {
		return ""
	}
	if math.Abs(a) == 0 {
		return fmt.Sprintf("%s%f", _bccy, b)
	}
	return fmt.Sprintf("%s%f", _accy, a)
}

// CalcPriceIfZero returns price calculation if the input price was bad
func CalcPriceIfZero(_px, _badpx string, _val, _qty float64) string {
	if math.Abs(StrToFloat(_px)) > 0.0001 {
		return _px
	}
	if math.Abs(_qty) < 0.0001 {
		return _badpx
	}
	return fmt.Sprintf("%f", _val/_qty)
}

// CleanString replaces comma with semi
func CleanString(_str string) string {
	_str = strings.Replace(_str, ",", ";", -1)
	return _str
}

// CleanStringMaximally removes various chars from name
func CleanStringMaximally(_str string) string {
	_str = strings.Replace(_str, ",", "", -1)
	_str = strings.Replace(_str, "-", "", -1)
	_str = strings.Replace(_str, "_", "", -1)
	_str = strings.Replace(_str, "/", "", -1)
	_str = strings.Replace(_str, ":", "", -1)
	_str = strings.Replace(_str, ",", "", -1)
	_str = strings.Replace(_str, "$", "", -1)
	_str = strings.Replace(_str, "%", "", -1)
	_str = strings.Replace(_str, ".", "", -1)
	_str = strings.Replace(_str, "@", "", -1)
	_str = strings.Replace(_str, " ", "", -1)
	_str = strings.Replace(_str, "	", "", -1)
	return _str
}

//================================================================================

// CsvCountTuple counts the number of parts under specified separator
func CsvCountTuple(_csvlist, _sep string) int {
	parts := strings.Split(_csvlist, _sep)
	count := 0
	for _, ss := range parts {
		ss1 := strings.TrimSpace(ss)
		if len(ss1) <= 0 {
			continue
		}
		count++
	}
	return count
}

// CsvLastTuple returns the last item in the tuple, or blank
func CsvLastTuple(_csvlist, _sep string) string {
	if len(_csvlist) < 1 {
		return ""
	}
	if _csvlist == _sep {
		return ""
	}
	parts := strings.Split(_csvlist, _sep)
	lenii := len(parts)
	for ii := 0; ii < lenii; ii++ {
		ss1 := strings.TrimSpace(parts[lenii-1-ii])
		if len(ss1) <= 0 {
			continue
		}
		return ss1
	}
	return ""
}

// CsvLastTupleElem returns the numbered sub-element of the csv's last element (itself viewed as an elemlist)
func CsvLastTupleElem(_csvlist, _sep, _elemsep, _badval string, _partno int) string {
	if len(_csvlist) < 1 {
		return _badval
	}
	if _csvlist == _sep {
		return _badval
	}
	parts := strings.Split(_csvlist, _sep)
	lenii := len(parts)
	for ii := 0; ii < lenii; ii++ {
		ss1 := strings.TrimSpace(parts[lenii-1-ii])
		if len(ss1) <= 0 {
			continue
		}
		elems := strings.Split(ss1, _elemsep)
		if len(elems) < _partno+1 {
			return _badval
		}
		return strings.TrimSpace(elems[_partno])
	}
	return _badval
}

// GetKV obtains the value from csvlist of kvps or the default
func GetKV(_list, _kk, _def string) string {
	parts := strings.Split(_list, ";")
	for _, kvp := range parts {
		kvparts := strings.Split(kvp, "=")
		if len(kvparts) < 2 {
			continue
		}
		if kvparts[0] == _kk {
			return kvparts[1]
		}
	}
	return _def
}

// ModifyKV updates the value in a csvlist of kvps
func ModifyKV(_list, _kk, _val string) string {
	parts := strings.Split(_list, ";")
	kvmap := map[string]string{}
	for _, kvp := range parts {
		kvparts := strings.Split(kvp, "=")
		if len(kvparts) < 2 {
			continue
		}
		kvmap[kvparts[0]] = kvparts[1]
	}
	kvmap[_kk] = _val
	return GenKVFromMap(kvmap)
}

// GetMapFromKV returns the csvlist of kvps as a map
func GetMapFromKV(_list string) map[string]string {
	parts := strings.Split(_list, ";")
	kvmap := map[string]string{}
	for _, kvp := range parts {
		kvparts := strings.Split(kvp, "=")
		if len(kvparts) < 2 {
			continue
		}
		kvmap[kvparts[0]] = kvparts[1]
	}
	return kvmap
}

// GenKVFromMap returns the map as a csvlist of kvps
func GenKVFromMap(_kvmap map[string]string) string {
	parts := []string{}
	for kk, val := range _kvmap {
		parts = append(parts, fmt.Sprintf("%s=%s", kk, val))
	}
	return strings.Join(parts, ";")
}

// GetKVFloat obtains the value from csvlist of kvps or the default
func GetKVFloat(_list, _kk string, _def float64) float64 {
	parts := strings.Split(_list, ";")
	for _, kvp := range parts {
		kvparts := strings.Split(kvp, "=")
		if len(kvparts) < 2 {
			continue
		}
		if kvparts[0] == _kk {
			return StrToFloat(kvparts[1])
		}
	}
	return _def
}

// GetNocasekeyKV (case-insensitively) obtains the value from csvlist of kvps or the default.
func GetNocasekeyKV(_list, _kk, _def string) string {
	parts := strings.Split(_list, ";")
	for _, kvp := range parts {
		kvparts := strings.Split(kvp, "=")
		if len(kvparts) < 2 {
			continue
		}
		if strings.ToLower(kvparts[0]) == strings.ToLower(_kk) {
			return kvparts[1]
		}
	}
	return _def
}

// OverrideWithKVMap does map lookup with a default
func OverrideWithKVMap(_mp map[string]string, _key, _alt string) string {
	if val, ok := _mp[_key]; ok {
		return val
	}
	return _alt
}

// IndexNl finds newline in byte buffer of len _buflen starting at _rowBegin
func IndexNl(_buffer []byte, _buflen, _rowBegin int) int {
	idx := _rowBegin
	for ; idx < _buflen; idx++ {
		if _buffer[idx] != '\n' {
			continue
		}
		idx++
		return idx
	}
	return _buflen // return end of buf
}

// JoinSlice joins slice elements using named separator
func JoinSlice(_strarr []string, _sep string) string {
	switch strings.ToLower(_sep) {
	case "pipe", "|":
		return strings.Join(_strarr, "|")
	case "plus", "+":
		return strings.Join(_strarr, "+")
	case "comma", ",":
		return strings.Join(_strarr, ",")
	case "space", " ":
		return strings.Join(_strarr, " ")
	case "tab", "	":
		return strings.Join(_strarr, "	")
	}
	return strings.Join(_strarr, _sep)
}

// JoinSliceWithReverse joins slice elements using named separator, and optionally in reverse
func JoinSliceWithReverse(_strarr []string, _sep string, _reverse bool) string {
	switch strings.ToLower(_sep) {
	case "pipe", "|":
		_sep = "|"
	case "plus", "+":
		_sep = "+"
	case "comma", ",":
		_sep = ","
	case "space", " ":
		_sep = " "
	case "tab", "	":
		_sep = "	"
	}

	if !_reverse {
		return strings.Join(_strarr, _sep)
	}
	str := ""
	nn := len(_strarr)
	for ii := 0; ii < nn; ii++ {
		if ii > 0 {
			str += _sep
		}
		str += _strarr[nn-ii-1]
	}
	return str
}

// JoinSliceLimitingColumns joins slice elements using named separator, and breaking into new "rows" when max cols is reached
func JoinSliceLimitingColumns(_strarr []string, _sep, _rowsep string, _maxcol int) string {
	_sep = sepmap[_sep]
	_rowsep = sepmap[_rowsep]
	inlen := len(_strarr)
	nrow := int(inlen / _maxcol)
	ostr := ""
	nn := 0
	for ii := 0; ii < nrow; ii++ {
		if ii > 0 {
			ostr += _rowsep
		}
		ostr += strings.Join(_strarr[nn:(nn+_maxcol)], _sep)
		nn += _maxcol
	}
	if nn < inlen {
		if nn > 0 {
			ostr += _rowsep
		}
		ostr += strings.Join(_strarr[nn:], _sep)
	}
	return ostr
}

// CopyStrSlice copies a string slice, optionally with prefix and suffix
func CopyStrSlice(_strarr []string, _prefix, _suffix string) []string {
	newarr := append([]string(nil), _strarr...)
	for ii := range newarr {
		newarr[ii] = _prefix + newarr[ii] + _suffix
	}
	return newarr
}

// SepReplace replaces one named separator with another
func SepReplace(_str, _insep, _outsep string) string {
	return strings.Replace(_str, sepmap[_insep], sepmap[_outsep], -1)
}

// SepMap obtains the separator
func SepMap(_sep string, _anycase bool) string {
	switch _anycase {
	case true:
		return sepmap[strings.ToLower(_sep)]
	}
	return sepmap[_sep]
}

// Str2Bool is shorthand
func Str2Bool(_str string) bool {
	switch strings.ToLower(_str) {
	case "yes", "true":
		return true
	}
	return false
}

// IsPositiveInteger returns true if the string is int > 0
func IsPositiveInteger(_str string) bool {
	nn := len(_str)
	if nn < 1 {
		return false
	}
	for ii := 0; ii < nn; ii++ {
		switch _str[ii] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return false
		}
	}
	return true
}

// SplitFilename returns (filename, extension, ok) from the input filename, for an allowed list of extensions, e.g. ["txt", "DAT", ...]
func SplitFilename(_fname string, _extlist []string) (string, string, bool) {
	if len(_fname) < 1 {
		return "", "", false
	}
	parts := strings.Split(_fname, ".")
	nparts := len(parts)
	if nparts < 1 {
		return "", "", false
	}
	lastpart := parts[nparts-1]
	for _, ext := range _extlist {
		if strings.ToLower(ext) == strings.ToLower(lastpart) {
			return strings.Join(parts[:nparts-1], "."), lastpart, true
		}
	}
	return "", "", false
}

// SplitFilename2 splits based on (consumed) string or position (negative counts from right)
func SplitFilename2(_str string, _ii interface{}) (string, string) {
	// fmt.Println("SplitFilename2 str=", _str)
	ns := len(_str)
	switch _ii.(type) {
	case int:
		ii := _ii.(int)
		if ii < 0 {
			ii = ns + ii
		}
		if ii >= ns {
			return _str, ""
		}
		return _str[:ii], _str[ii:]
	case string:
		ss := _ii.(string)
		ix := strings.Index(_str, ss)
		if ix < 0 {
			return _str, ""
		}
		return _str[:ix], _str[(ix + len(ss)):]
	}
	return "", ""
}

// SplitFilename3 splits into 3 parts, using delim which is either string or positional index (which may be negative for counting from end)
func SplitFilename3(_str string, _ii0, _ii1 interface{}) (string, string, string) {
	// fmt.Println("SplitFilename3 str=", _str)
	ns := len(_str)
	aa, bb, cc, rest := "", "", "", ""
	ii := [2]int{-999999, -999999}
	switch _ii0.(type) {
	case int:
		ii[0] = _ii0.(int)
		if ii[0] < 0 {
			ii[0] = ns + ii[0]
		} // fmt.Println("SplitFilename3 first int")

	}
	switch _ii1.(type) {
	case int:
		ii[1] = _ii1.(int)
		if ii[1] < 0 {
			ii[1] = ns + ii[1]
		} // fmt.Println("SplitFilename3 second int")
	}
	switch {
	case (ii[0] > -999999) && (ii[1] > -999999): // fmt.Println("SplitFilename3 case both ii0=", ii[0], "ii1=", ii[1])
		if ii[0] < ns {
			aa, rest = _str[:ii[0]], _str[ii[0]:]
			if (ii[1] >= ii[1]) && (ii[1] < ns) {
				bb, cc = _str[ii[0]:ii[1]], _str[ii[1]:]
			} else if ii[1] < ii[0] {
				cc = rest
			} else if ii[1] > ns {
				bb = rest
			}
		} else {
			aa = _str
		}
	case (ii[0] > -999999): // fmt.Println("SplitFilename3 case first")
		if ii[0] < ns {
			aa = _str[:ii[0]]
			bb, cc = SplitFilename2(_str[ii[0]:], _ii1)
		} else {
			aa = _str
		}
	case ii[1] > -999999: // fmt.Println("SplitFilename3 case second strlen=", ns, "ii1=", ii[1])
		if ii[1] < ns {
			cc = _str[ii[1]:]
			aa, bb = SplitFilename2(_str[:ii[1]], _ii0)
		} else {
			aa, bb = SplitFilename2(_str, _ii0)
		}
	default: // fmt.Println("SplitFilename3 case default")
		aa, rest = SplitFilename2(_str, _ii0)
		bb, cc = SplitFilename2(rest, _ii1)
	}
	// fmt.Println("SplitFilename3 str=", _str, "aa=", aa, "bb=", bb, "cc=", cc)
	return aa, bb, cc
}

// SplitFilename4 splits into 4 parts, using delim which is either string or positional index (which may be negative for counting from end)
func SplitFilename4(_str string, _ii0, _ii1, _ii2 interface{}) (string, string, string, string) {
	// // fmt.Println("SplitFilename4 str=", _str)
	aa, bb, cc, dd, rest := "", "", "", "", ""
	ns := len(_str)
	ii := [3]int{-999999, -999999, -999999}
	switch _ii0.(type) {
	case int:
		ii[0] = _ii0.(int)
		if ii[0] < 0 {
			ii[0] = ns + ii[0]
		} else if ii[0] > ns {
			ii[0] = ns
		} // // fmt.Println("SplitFilename4 first int")
	}
	switch _ii1.(type) {
	case int:
		ii[1] = _ii1.(int)
		if ii[1] < 0 {
			ii[1] = ns + ii[1]
		} else if ii[1] > ns {
			ii[1] = ns
		} // // fmt.Println("SplitFilename4 second int")
	}
	switch _ii2.(type) {
	case int:
		ii[2] = _ii2.(int)
		if ii[2] < 0 {
			ii[2] = ns + ii[2]
		} else if ii[2] > ns {
			ii[2] = ns
		} // // fmt.Println("SplitFilename4 second int")
	}
	switch {
	case (ii[0] > -999999) && (ii[1] > -999999) && (ii[2] > -999999): // fmt.Println("SplitFilename4 case all three ii0=", ii[0], "ii1=", ii[1], "ii2=", ii[2])
		if ii[0] >= ns {
			aa = _str
			break
		}
		aa, rest = _str[:ii[0]], _str[ii[0]:]
		bb, cc, dd = SplitFilename3(rest, ii[1]-len(aa), ii[2]-len(aa))
	case (ii[0] > -999999) && (ii[1] > -999999): // fmt.Println("SplitFilename4 case both ii0=", ii[0], "ii1=", ii[1])
		if ii[0] >= ns {
			aa = _str
			break
		}
		aa, rest = _str[:ii[0]], _str[ii[0]:]
		bb, cc, dd = SplitFilename3(rest, ii[1]-len(aa), _ii2)
	case (ii[0] > -999999) && (ii[2] > -999999): // fmt.Println("SplitFilename4 case both ii0=", ii[0], "ii2=", ii[2])
		if ii[0] >= ns {
			aa = _str
			break
		}
		aa, rest = _str[:ii[0]], _str[ii[0]:]
		bb, cc, dd = SplitFilename3(rest, _ii1, ii[2]-len(aa))
	case (ii[1] > -999999) && (ii[2] > -999999): // fmt.Println("SplitFilename4 case both ii1=", ii[1], "ii2=", ii[2])
		rest, cc, dd = SplitFilename3(_str, ii[1], ii[2])
		aa, bb = SplitFilename2(rest, _ii0)
	case (ii[0] > -999999): // fmt.Println("SplitFilename4 case first")
		if ii[0] >= ns {
			aa = _str
			break
		}
		aa = _str[:ii[0]]
		bb, cc, dd = SplitFilename3(_str[ii[0]:], _ii1, _ii2)
	case (ii[2] > -999999): // fmt.Println("SplitFilename4 case third")
		if ii[2] < ns {
			dd = _str[ii[2]:]
			aa, bb, cc = SplitFilename3(_str[:ii[2]], _ii0, _ii1)
		} else {
			aa, bb, cc = SplitFilename3(_str, _ii0, _ii1)
		}
	case (ii[1] > -999999): // fmt.Println("SplitFilename4 case second")
		if ii[1] < ns {
			aa, bb = SplitFilename2(_str[:ii[1]], _ii0)
			cc, dd = SplitFilename2(_str[ii[1]:], _ii2)
		} else {
			aa, bb = SplitFilename2(_str, _ii0)
		}
	default: // fmt.Println("SplitFilename4 case default")
		aa, rest = SplitFilename2(_str, _ii0)
		bb, cc, dd = SplitFilename3(rest, _ii1, _ii2)
	}
	// fmt.Println("SplitFilename4 str=", _str, "aa=", aa, "bb=", bb, "cc=", cc, "dd=", dd)
	return aa, bb, cc, dd
}

// SplitToIntSlice convert "1,2,3" to slice of ints
func SplitToIntSlice(_str, _sep string) []int {
	osl := []int{}
	if len(_str) > 0 {
		parts := strings.Split(_str, _sep)
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if len(part) < 1 {
				continue
			}
			osl = append(osl, Toint0(part))
		}
	}
	return osl
}

// SplitToStrSlice converts "1,2,3" to slice of strings, ignoring blanks
func SplitToStrSlice(_str, _sep string) []string {
	osl := []string{}
	if len(_str) > 0 {
		parts := strings.Split(_str, _sep)
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if len(part) < 1 {
				continue
			}
			osl = append(osl, part)
		}
	}
	return osl
}

// IntSliceContains is shorthand
func IntSliceContains(_sl []int, _num int) bool {
	for _, ss1 := range _sl {
		if _num == ss1 {
			return true
		}
	}
	return false
}

// StrSliceContains is shorthand
func StrSliceContains(_sl []string, _num string) bool {
	for _, ss1 := range _sl {
		if _num == ss1 {
			return true
		}
	}
	return false
}

// IsCancellingStrings checks if strings are offsetting numbers, and informs which is the negative one
func IsCancellingStrings(_lstr, _rstr string, _leftneg *bool) bool {
	*_leftneg = false
	llen, rlen := len(_lstr), len(_rstr)
	switch llen - rlen {
	case 1:
		if (_lstr[0] == '-') && (_lstr[1:] == _rstr) {
			*_leftneg = true
			return true
		}
	case -1:
		if (_rstr[0] == '-') && (_rstr[1:] == _lstr) {
			*_leftneg = false
			return true
		}
	}
	return false
}

// EqualAndZeroStrings tells if strings are equal, and both zero
func EqualAndZeroStrings(_lstr, _rstr string) bool {
	if _lstr != _rstr {
		return false
	}
	if StrToFloat(_lstr) != 0.0 {
		return false
	}
	return true
}

// EqualFloats tells if floats are equal
func EqualFloats(_f1, _f2 float64) bool {
	return (_f1 - _f2) < 0.0000001
}

// StrCapped returns truncated string if exceeds cap
func StrCapped(_str string, _cap int) string {
	if len(_str) <= _cap {
		return _str
	}
	return _str[:_cap]
}

// Flatline2 combines 2 strings, potentially eliminating newline on the first
func Flatline2(_str0, _str1 string) string {
	// fmt.Printf("Flatline2 called str0(%s) str1(%s)\n", _str0, _str1)
	switch {
	case _str0 == "":
		return _str1
	case strings.HasSuffix(_str0, "\n"):
		return _str0[:len(_str0)-1] + _str1
	}
	return _str0 + _str1
}

// SmartFlatline2 combines 2 strings, potentially eliminating newline on the first.
// In addition, it guesses whether to insert line-ending period.
func SmartFlatline2(_str0, _str1 string) string {
	// fmt.Printf("SmartFlatline2 called str0(%s) str1(%s)\n", _str0, _str1)
	_str0 = strings.TrimSpace(_str0)
	_str1 = strings.TrimSpace(_str1)

	switch {
	case _str0 == "":
		return _str1
	case _str1 == "":
		return _str0
	}

	len0 := len(_str0)
	lc := _str0[len0-1:]
	isEnding0 := lc == "."

	fc := _str1[:1]
	isUpper1 := fc == strings.ToUpper(fc)

	switch {
	case isEnding0 && isUpper1:
		return _str0 + " " + _str1
	case (!isEnding0) && isUpper1:
		return _str0 + ". " + _str1
	case (!isEnding0) && !isUpper1:
		return _str0 + " " + _str1
	}
	return _str0 + " " + _str1
}

// ChompStr eliminates suffix string if present
func ChompStr(_str, _chompstr string) string {
	if strings.HasSuffix(_str, _chompstr) {
		return _str[:len(_str)-len(_chompstr)]
	}
	return _str
}

// ChompChar eliminates suffix char if present
func ChompChar(_str string, _chompchar uint8) string {
	nn := len(_str)
	if (nn > 0) && (_str[nn-1] == _chompchar) {
		return _str[:nn-1]
	}
	return _str
}

// ChompParens eliminates dual delimiting strings if present
func ChompParens(_str string, _trimSpace bool) string {
	if _trimSpace {
		_str = strings.TrimSpace(_str)
	}
	num := len(_str)
	if num < 1 {
		return _str
	}
	ch0, chN, doit := _str[0], _str[num-1], false
	switch {
	case (ch0 == '(') && (chN == ')'):
		doit = true
	case (ch0 == '[') && (chN == ']'):
		doit = true
	case (ch0 == '{') && (chN == '}'):
		doit = true
	}
	if doit {
		_str = _str[1:(num - 1)]
	}
	if _trimSpace {
		_str = strings.TrimSpace(_str)
	}
	return _str
}

// ChompQuotes eliminates dual delimiting quotes if present
func ChompQuotes(_str string, _trimSpace bool) string {
	if _trimSpace {
		_str = strings.TrimSpace(_str)
	}
	num := len(_str)
	if num < 1 {
		return _str
	}
	ch0, chN, doit := _str[0], _str[num-1], false
	switch {
	case (ch0 == '\'') && (chN == '\''):
		doit = true
	case (ch0 == '"') && (chN == '"'):
		doit = true
	}
	if doit {
		_str = _str[1:(num - 1)]
	}
	if _trimSpace {
		_str = strings.TrimSpace(_str)
	}
	return _str
}

// StryyyymmddLTTernary is shorthand
func StryyyymmddLTTernary(_dt1, _dt2, _trueStr, _falseStr string) string {
	if StryyyymmddLT(_dt1, _dt2) {
		return _trueStr
	}
	return _falseStr
}

// StryyyymmddLTEQTernary is shorthand
func StryyyymmddLTEQTernary(_dt1, _dt2, _trueStr, _falseStr string) string {
	if StryyyymmddLTEQ(_dt1, _dt2) {
		return _trueStr
	}
	return _falseStr
}

// SetupLogger returns a logger
func SetupLogger(_logfilepath, _logcontentprefix string) (lglocal *log.Logger, err error) {
	lglocal, err = nil, nil
	if err = os.MkdirAll(path.Dir(_logfilepath), 0755); err != nil {
		return
	}
	fp, err := os.Create(_logfilepath)
	if err != nil {
		return
	}
	lglocal = log.New(fp, _logcontentprefix, log.LstdFlags)
	return
}

// GetFileLineCount counts non-comment lines of a file
func GetFileLineCount(_fname, _comments string) (int64, error) {
	comments := strings.Split(_comments, ",")
	bio, err := OpenAnyErr(_fname)
	if err != nil {
		return 0, err
	}
	count := int64(0)
	var line []byte
	for {
		line, err = bio.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		if err == bufio.ErrBufferFull {
			continue
		}
		line = line[0 : len(line)-1]
		//fmt.Printf("Line:*%s*\n", string(line))
		if IsCommentLine(line, comments) { /*fmt.Println("comment");*/ continue
		}
		count++
	}
	return count, nil
}

// IsCommentLine checks if a line is one of the list of comment types
func IsCommentLine(_line []byte, _commenttags []string) bool {
	for _, commenttag := range _commenttags {
		switch commenttag {
		case "Whitespace":
			if (len(_line) > 0) && (len(bytes.Trim(_line, " ")) == 0) {
				return true
			}
		case "WhitespaceHash":
			if tmp := bytes.TrimLeft(_line, " "); len(tmp) > 0 {
				if tmp[0] == '#' {
					return true
				}
			}
		}
	}
	return false
}

// Resplit replaces separator with another, then splits it on yet another separator
func Resplit(_str, _osep, _njoin, _nsep string) []string {
	return strings.Split(strings.Join(strings.Split(_str, _osep), _njoin), _nsep)
}

// SpaceSplitter informs if the input char is a split char
func SpaceSplitter(inp rune) bool {
	switch inp {
	case 32:
		return true
	case 10:
		return true
	default:
		return false
	}
	return false
}

// CleanAndSplitOnSpaces trims each line obtained by splitting input paragraph at newline.
// Each line is split on space, and joined back by specified separator
func CleanAndSplitOnSpaces(_str, _sep string) []string {
	outlines := []string{}
	lines := strings.Split(_str, "\n")
	for _, str := range lines {
		str = strings.TrimSpace(str)
		if len(str) < 1 {
			continue
		}
		// outitems	:= []string{}
		newitems := strings.FieldsFunc(str, SpaceSplitter)
		items := []string{}
		for _, item := range newitems {
			item = strings.TrimSpace(item)
			items = append(items, item)
		}
		line := strings.Join(items, _sep)
		// fmt.Printf("line is %s\n", line)
		outlines = append(outlines, line)
	}
	return outlines
}

// CleanAndSplitOnSeparator trims each line obtained by splitting input paragraph at newline.
// Each line is split on separator insep, and joined back by specified separator
func CleanAndSplitOnSeparator(_str, _insep, _sep string) []string {
	outlines := []string{}
	lines := strings.Split(_str, "\n")
	for _, str := range lines {
		str = strings.TrimSpace(str)
		if len(str) < 1 {
			continue
		}
		// outitems	:= []string{}
		newitems := strings.Split(str, _insep)
		items := []string{}
		for _, item := range newitems {
			item = strings.TrimSpace(item)
			items = append(items, item)
		}
		line := strings.Join(items, _sep)
		// fmt.Printf("line is %s\n", line)
		outlines = append(outlines, line)
	}
	return outlines
}

// UniqueKeys returns array of unique strings within the input string array
func UniqueKeys(_strs ...[]string) (strs []string) {
	mp := map[string]bool{}
	for _, ss := range _strs {
		for _, str := range ss {
			mp[str] = true
		}
	}
	for kk := range mp {
		strs = append(strs, kk)
	}
	return
}

// SortedUniqueKeys returns sorted array of unique strings within the input string array
func SortedUniqueKeys(_strs ...[]string) (strs []string) {
	mp := map[string]bool{}
	for _, ss := range _strs {
		for _, str := range ss {
			mp[str] = true
		}
	}
	for kk := range mp {
		strs = append(strs, kk)
	}
	sort.Strings(strs)
	return
}

// ShrinkSep removes Semi at EOL and shrinks multiple Semi to single Semi prior to EOL
func ShrinkSep(_str string, _ch byte) string {
	sep := string(_ch)
	twosep := sep + sep
	done := false
	for !done {
		switch {
		case strings.Contains(_str, twosep):
			_str = strings.Replace(_str, twosep, sep, -1)
		case strings.Contains(_str, sep) && (_str[len(_str)-1] == _ch):
			_str = _str[:len(_str)-1]
		case strings.Contains(_str, sep) && (_str[0] == _ch):
			_str = _str[1:]
		default:
			done = true
		}
	}
	return _str
}

// StrReplaceWithMap replaces substrings in the input string based on the map passed
func StrReplaceWithMap(_instr string, _mp map[string]string) string {
	outstr := _instr
	for key, val := range _mp {
		outstr = strings.Replace(outstr, key, val, -1)
	}
	return outstr
}

// Next 4 functions are for printing colour text.
// Usage example:  fmt.Println(GreenBold("Success:") + "Limit check passed")

// Green sets a color
func Green(in string) (out string) {
	return "\033[32m" + in + "\033[0m"
}

// GreenBold sets a color
func GreenBold(in string) (out string) {
	return "\033[1;32m" + in + "\033[0m"
}

// Red sets a color
func Red(in string) (out string) {
	return "\033[31m" + in + "\033[0m"
}

// RedBold sets a color
func RedBold(in string) (out string) {
	return "\033[1;31m" + in + "\033[0m"
}
