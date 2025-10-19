package locality

const (
	numRows = 1024
	numCols = 1024
)

type Data [numRows][numCols]bool

func CountTrueElementsRowByRow(data *Data) int {
	count := 0

	for i := range numRows {
		for j := range numCols {
			if data[i][j] {
				count++
			}
		}
	}

	return count
}

func CountTrueElementsColumnByColumn(data *Data) int {
	count := 0

	for j := range numCols {
		for i := range numRows {
			if data[i][j] {
				count++
			}
		}
	}

	return count
}

type ListNode struct {
	Value bool
	Next  *ListNode
}

func CountTrueElementsLinkedList(head *ListNode) int {
	count := 0

	for head != nil {
		if head.Value {
			count++
		}
		head = head.Next
	}

	return count
}
