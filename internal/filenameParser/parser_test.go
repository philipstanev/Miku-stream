package nameparser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	input1 := "[Sav1or] 5-toubun no Hanayome (The Quintessential Quintuplets) - S01E01 - The Quintessential Quintuplets [BD][1080p][AV1][OPUS][Dual Audio].mkv"
	expected1 := Episode{"5-toubun no Hanayome (The Quintessential Quintuplets)", "The Quintessential Quintuplets", "S01E01"}

	got1, err := ParseToStruct(input1)
	if err != nil {
		t.Errorf("Error in parsing: %s", err.Error())
	}

	fmt.Println(input1)

	if !reflect.DeepEqual(expected1, got1) {
		t.Errorf("Input: %s \n Expected: %+v \n Got: %+v", input1, expected1, got1)
	}

	input2 := "Death Note - 01x13 - Confession.mkv"
	expected2 := Episode{"Death Note", "Confession", "01x13"}

	got2, err := ParseToStruct(input2)
	if err != nil {
		t.Errorf("Error in parsing: %s", err.Error())
	}

	fmt.Println(input2)

	if !reflect.DeepEqual(expected2, got2) {
		t.Errorf("Input: %s \n Expected: %+v \n Got: %+v", input2, expected2, got2)
	}

}
