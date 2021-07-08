package autoconfig_test

import (
	"eurocontrol.io/digital-platform-product-deployment/pkg/autoconfig"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestExampleValueOrPanic(t *testing.T) {
	// Try different approach of unit testing with the following example ?
	var want = 10
	var got int
	autoconfig.ValueOrPanic(&got, "just.one.value|10")
	//fmt.Println(got)

	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}

	// Output:
	// 10
}

func TestExampleDurationOrPanic(t *testing.T) {
	os.Setenv("JUST_ONE_VALUE", "10")
	tenMillis := autoconfig.DurationOrPanic("just.one.value|1", "not a correct duration, so getting millis")
	fmt.Println(tenMillis)

	fiveSeconds := autoconfig.DurationOrPanic("yet.another.value|5", "seconds")
	fmt.Println(fiveSeconds)

	// Output:
	// 10ms
	// 5s
}

func TestExampleAutoConfigure_embeddedStruct(t *testing.T) {
	type Em struct {
		Embedded string `value:"embedded.property|it"`
	}
	var c struct {
		Em
		Another string `value:"my.property|works"`
	}
	autoconfig.OrPanic(&c)
	autoconfig.OrPanic(&c.Em)
	fmt.Println(c.Embedded, c.Another)

	// Output:
	// it works
}

func TestExampleAutoConfigure_duration(t *testing.T) {
	var c struct {
		Ten        time.Duration `value:"one|10"`
		TenSeconds time.Duration `value:"two|10"`
		TenMinutes time.Duration `value:"three|10"`
	}
	autoconfig.OrPanic(&c)
	fmt.Println(c.Ten)
	fmt.Println(c.TenSeconds)
	fmt.Println(c.TenMinutes)

	// Output:
	// 10ms
	// 10s
	// 10m0s
}

func printPanic() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}
