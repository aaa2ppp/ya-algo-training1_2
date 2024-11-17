package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		debug   bool
	}{
		{
			"1",
			args{strings.NewReader(`9 1
1 0
1 1
1 1
0 0
1
0
1
0
1
`)},
			`1
`,
			true,
		},
		{
			"2",
			args{strings.NewReader(`5 0
1 1
0 0
1
1
0
`)},
			`IMPOSSIBLE
`,
			true,
		},
		{
			"001",
			args{strings.NewReader(`17 1
0 1
1 0
1 0
1 1
0 1
1 1
1 1
1 1
0
1
0
1
0
0
1
0
1`)},
			`2`,
			true,
		},
		{
			"004",
			args{strings.NewReader(`25 0
0 0
0 0
0 1
0 1
0 0
0 1
0 1
0 1
0 1
0 1
0 1
0 1
1
1
1
0
1
0
1
0
1
0
1
0
1
`)},
			`7`,
			true,
		},
		// {
		// 	"3",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// {
		// 	"4",
		// 	args{strings.NewReader(``)},
		// 	``,
		// 	true,
		// },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func(v bool) { debugEnable = v }(debugEnable)
			debugEnable = tt.debug
			out := &bytes.Buffer{}
			run(tt.args.in, out)
			if gotOut := out.String(); trimLines(gotOut) != trimLines(tt.wantOut) {
				t.Errorf("run() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func trimLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r\n")
	}
	for n := len(lines); n > 0 && lines[n-1] == ""; n-- {
		lines = lines[:n-1]
	}
	return strings.Join(lines, "\n")
}