package logic

type Range struct {
	StartByte int
	StopByte  int
	Index     int
}

type Splitter interface {
	GetRanges() []Range
}
