package gencore

import "testing"

func TestReceiverName(t *testing.T) {
	type args struct {
		typeName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"One char type name",
			args{typeName: "a"},
			"a",
		},
		{"More than one char type name",
			args{typeName: "ab"},
			"a",
		},
		{"Upper case prefix type name",
			args{typeName: "Ab"},
			"a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReceiverName(tt.args.typeName); got != tt.want {
				t.Errorf("ReceiverName() = %v, want %v", got, tt.want)
			}
		})
	}
}
