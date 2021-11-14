package main

import (
	"testing"
)

func TestCompleted(t *testing.T) {
	itemList := ReadItemStructData("TODOList.todo")

	for _, element := range itemList {
		if element.Done {
			t.Errorf("Completed item in TODO List!")
		}
	}
}

func TestUncompleted(t *testing.T) { // t here is test handler...

	itemList := ReadItemStructData("CompletedItemList.todo")

	for _, element := range itemList {
		if !element.Done {
			t.Errorf("Uncompleted item in Completed Item List!")
		}
	}

}
