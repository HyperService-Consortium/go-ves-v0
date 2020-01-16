package bni

import "testing"

func TestMockBNIStorage(t *testing.T) {
	b := &mockBNIStorage{}
	tester := bNIStorageTestSet{b}
	b.insertMockData(tester.MockingData())
	tester.RunTests(t)
}
