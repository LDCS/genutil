package genutil

import (
	"math"
	"sort"
)

// /* http://nerdyworm.com/blog/2013/05/15/sorting-a-slice-of-structs-in-go/ */
//================================================================================
type myelemFloat64SortAscending struct {
	kk_  string
	val_ float64
}
type myelemSliceFloat64SortAscending []myelemFloat64SortAscending

func (Slice myelemSliceFloat64SortAscending) Len() int           { return len(Slice) }
func (Slice myelemSliceFloat64SortAscending) Less(i, j int) bool { return Slice[i].val_ < Slice[j].val_ }
func (Slice myelemSliceFloat64SortAscending) Swap(i, j int)      { Slice[i], Slice[j] = Slice[j], Slice[i] }

// SortedKeysByVal_String2Float64_Ascending sorts by value for that maptype
func SortedKeysByVal_String2Float64_Ascending(_mp *map[string]float64) []string {
	mlen := len(*_mp)
	var mp myelemSliceFloat64SortAscending = make([]myelemFloat64SortAscending, mlen)
	ii := 0
	for kk, vv := range *_mp {
		mp[ii].kk_, mp[ii].val_ = kk, vv
		ii++
	}
	sort.Sort(mp)
	kkarr := make([]string, mlen)
	ii = 0
	for idx := range mp {
		kkarr[ii] = mp[idx].kk_
		ii++
	}
	return kkarr
}

//================================================================================
type myelemFloat64SortAbsAscending struct {
	kk_  string
	val_ float64
}
type myelemSliceFloat64SortAbsAscending []myelemFloat64SortAbsAscending

func (Slice myelemSliceFloat64SortAbsAscending) Len() int { return len(Slice) }
func (Slice myelemSliceFloat64SortAbsAscending) Less(i, j int) bool {
	return math.Abs(Slice[i].val_) < math.Abs(Slice[j].val_)
}
func (Slice myelemSliceFloat64SortAbsAscending) Swap(i, j int) {
	Slice[i], Slice[j] = Slice[j], Slice[i]
}

// SortedKeysByVal_String2Float64_AbsAscending sorts by value for that maptype
func SortedKeysByVal_String2Float64_AbsAscending(_mp *map[string]float64) []string {
	mlen := len(*_mp)
	var mp myelemSliceFloat64SortAbsAscending = make([]myelemFloat64SortAbsAscending, mlen)
	ii := 0
	for kk, vv := range *_mp {
		mp[ii].kk_, mp[ii].val_ = kk, vv
		ii++
	}
	sort.Sort(mp)
	kkarr := make([]string, mlen)
	ii = 0
	for idx := range mp {
		kkarr[ii] = mp[idx].kk_
		ii++
	}
	return kkarr
}

//================================================================================
type myelemFloat64SortDescending struct {
	kk_  string
	val_ float64
}
type myelemSliceFloat64SortDescending []myelemFloat64SortDescending

func (Slice myelemSliceFloat64SortDescending) Len() int { return len(Slice) }
func (Slice myelemSliceFloat64SortDescending) Less(i, j int) bool {
	return Slice[i].val_ > Slice[j].val_
}
func (Slice myelemSliceFloat64SortDescending) Swap(i, j int) { Slice[i], Slice[j] = Slice[j], Slice[i] }

// SortedKeysByVal_String2Float64_Descending sorts by value for that maptype
func SortedKeysByVal_String2Float64_Descending(_mp *map[string]float64) []string {
	mlen := len(*_mp)
	var mp myelemSliceFloat64SortDescending = make([]myelemFloat64SortDescending, mlen)
	ii := 0
	for kk, vv := range *_mp {
		mp[ii].kk_, mp[ii].val_ = kk, vv
		ii++
	}
	sort.Sort(mp)
	kkarr := make([]string, mlen)
	ii = 0
	for idx := range mp {
		kkarr[ii] = mp[idx].kk_
		ii++
	}
	return kkarr
}

//================================================================================
type myelemFloat64SortAbsDescending struct {
	kk_  string
	val_ float64
}
type myelemSliceFloat64SortAbsDescending []myelemFloat64SortAbsDescending

func (Slice myelemSliceFloat64SortAbsDescending) Len() int { return len(Slice) }
func (Slice myelemSliceFloat64SortAbsDescending) Less(i, j int) bool {
	return math.Abs(Slice[i].val_) > math.Abs(Slice[j].val_)
}
func (Slice myelemSliceFloat64SortAbsDescending) Swap(i, j int) {
	Slice[i], Slice[j] = Slice[j], Slice[i]
}

// SortedKeysByVal_String2Float64_AbsDescending sorts by value for that maptype
func SortedKeysByVal_String2Float64_AbsDescending(_mp *map[string]float64) []string {
	mlen := len(*_mp)
	var mp myelemSliceFloat64SortAbsDescending = make([]myelemFloat64SortAbsDescending, mlen)
	ii := 0
	for kk, vv := range *_mp {
		mp[ii].kk_, mp[ii].val_ = kk, vv
		ii++
	}
	sort.Sort(mp)
	kkarr := make([]string, mlen)
	ii = 0
	for idx := range mp {
		kkarr[ii] = mp[idx].kk_
		ii++
	}
	return kkarr
}

//================================================================================
//================================================================================
//================================================================================
//================================================================================
