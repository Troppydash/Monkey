package object

import "testing"

// Test Hashing Strings
func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}

	diff1 := &String{Value: "My name is perry"}
	diff2 := &String{Value: "My name is perry"}

	var n = 10
	for i := 0; i < n; i++ {
		if hello1.HashKey() != hello2.HashKey() {
			t.Errorf("strings with same content have different keys. %v != %v",
				hello1.HashKey(), hello2.HashKey())
		}

		if diff1.HashKey() != diff2.HashKey() {
			t.Errorf("strings with same content have different keys. %v != %v",
				diff1.HashKey(), diff2.HashKey())
		}

		if hello1.HashKey() == diff1.HashKey() {
			t.Errorf("strings with different content have same hash keys. %v == %v",
				hello1.HashKey(), diff1.HashKey())
		}
	}
}
