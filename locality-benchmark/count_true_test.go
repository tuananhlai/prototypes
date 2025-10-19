package locality_test

import (
	"testing"

	"github.com/tuananhlai/prototypes/locality-benchmark"
)

var (
	data     = &locality.Data{}
	listNode = &locality.ListNode{}
)

func init() {
	numRows := len(data)
	numCols := len(data[0])
	x := listNode

	for i := range numRows {
		for j := range numCols {
			val := i == j
			data[i][j] = val
			x.Value = val
			if i < numRows-1 || j < numCols-1 {
				x.Next = &locality.ListNode{}
				x = x.Next
			}
		}
	}
}

func BenchmarkCountTrueElementsRowByRow(b *testing.B) {
	for b.Loop() {
		locality.CountTrueElementsRowByRow(data)
	}
}

func BenchmarkCountTrueElementsColumnByColumn(b *testing.B) {
	for b.Loop() {
		locality.CountTrueElementsColumnByColumn(data)
	}
}

func BenchmarkCountTrueElementsLinkedList(b *testing.B) {
	for b.Loop() {
		locality.CountTrueElementsLinkedList(listNode)
	}
}
