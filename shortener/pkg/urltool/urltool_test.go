package urltool

import "testing"

func TestGetbasePath(t *testing.T) {
	type args struct {
		lurl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "基本示例", args: args{lurl: "https://www.liwenzhou.com/posts/Go/golang-menu/"}, want: "golang-menu", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetbasePath(tt.args.lurl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetbasePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetbasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
