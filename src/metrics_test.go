package main

import (
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func TestGetMetric(t *testing.T) {
	testCases := []struct {
		header string
		data   string
		result map[string]interface{}
	}{
		{"", "", make(map[string]interface{})},
		{" ", " ", make(map[string]interface{})},
		{"snafu", "fubar", make(map[string]interface{})},
		{"1 2 3 4 5 6", "", make(map[string]interface{})},
		{"1 2 3 4 5 6", " ts ", map[string]interface{}{"source": "4", "destination": "5", "ts": "true"}},
		{"1 2 3 4 5 6", " advmss:9999 ", map[string]interface{}{"source": "4", "destination": "5", "advmss": int64(9999)}},
		{"1 2 3 4 5 6", " minrtt:99.99 ", map[string]interface{}{"source": "4", "destination": "5", "minrtt": float64(99.99)}},
		{"1 2 3 4 5 6", " delivery_rate 2.7Mbps", map[string]interface{}{"source": "4", "destination": "5", "delivery_rate": float64(2700000)}},
		{"1 2 3 4 5 6", " rtt:32.354/16.077", map[string]interface{}{"source": "4", "destination": "5", "rtt_average": float64(32.354), "rtt_std_dev": float64(16.077)}},
		{"1 2 3 4 5 6", " wscale:8,7", map[string]interface{}{"source": "4", "destination": "5", "snd_wscale": int64(8), "rcv_wscale": int64(7)}},
	}

	for _, tc := range testCases {
		r := getMetric(tc.header, tc.data)
		if !reflect.DeepEqual(r, tc.result) {
			t.Errorf("Invalid result: %s", pretty.Diff(r, tc.result))
		}
	}
}

func TestGetFilter(t *testing.T) {
	testCases := []struct {
		src    string
		dst    string
		result string
	}{
		{"", "", ""},
		{"1.2.3.4", "", "( src 1.2.3.4 )"},
		{"1.2.3.4 5.6.7.8", "", "( src 1.2.3.4 or src 5.6.7.8 )"},
		{"   1.2.3.4    ", "5.6.7.8", "( src 1.2.3.4 or dst 5.6.7.8 )"},
	}

	for _, tc := range testCases {
		r := getFilter(tc.src, tc.dst)
		if r != tc.result {
			t.Errorf("Invalid result: %s", pretty.Diff(r, tc.result))
		}
	}
}

func TestGetCommandArgs(t *testing.T) {
	testCases := []struct {
		resolve bool
		result  string
	}{
		{true, "-iotr"},
		{false, "-iot"},
	}
	Args.SSArgs = "-iot"
	for _, tc := range testCases {
		Args.Resolve = tc.resolve
		r := getCommandArgs()
		if r != tc.result {
			t.Errorf("Invalid result: %s", pretty.Diff(r, tc.result))
		}
	}

}
