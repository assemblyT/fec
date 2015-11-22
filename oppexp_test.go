package fec

import (
	"fmt"
	"testing"
	"time"
)

func TestFromString(t *testing.T) {
	var parsed Oppexp
	err := FromStringSlice(&parsed, []string{
		"C00004606", "N", "2015",
		"M4", "15951124870", "21B",
		"F3X", "SB", "DIRECT MAIL SYSTEMS, INC.",
		"CLEARWATER", "FL", "33762",
		"03/13/2015", "600", "",
		"POSTAGE", "003", "Solicitation and Fundraising Expenses ",
		"", "", "ORG", "4041320151241802165",
		"1002259", "SB21B.20726", ""})
	if err != nil {
		fmt.Println(err)
		t.Error("failed to parse")
	}

	expected := Oppexp{
		"C00004606", 'N', 2015,
		"M4", "15951124870",
		"21B", "F3X", "SB",
		"DIRECT MAIL SYSTEMS, INC.", "CLEARWATER", "FL",
		"33762",
		time.Date(2015, time.March, 13, 0, 0, 0, 0, DefaultLocation()),
		600,
		"", "POSTAGE", "003", "Solicitation and Fundraising Expenses ",
		0, "", "ORG", 4041320151241802165, 1002259,
		"SB21B.20726", ""}

	if !parsed.Equals(expected) {
		fmt.Println(parsed)
		fmt.Println(expected)
		t.Error("Parsed oppexp did not equal the expected one")
	}
}

func TestIsEmpty(t *testing.T) {
	var nilOppexp Oppexp

	if !nilOppexp.IsEmpty() {
		t.Error("Uninitialized Oppexp should be empty")
	}
}
