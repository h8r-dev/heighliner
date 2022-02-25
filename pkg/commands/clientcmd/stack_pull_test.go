package clientcmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test_pullStack(t *testing.T) {
	type args struct {
		c    *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "no arg",
			args:    args{&cobra.Command{}, []string{}},
			wantErr: true,
		},
		{
			name:    "unexisted stack name",
			args:    args{&cobra.Command{}, []string{"unexisted-stack-name"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pullStack(tt.args.c, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("pullStack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
