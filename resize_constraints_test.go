package icat

import (
	"image"
	"testing"
)

func Test_resizeFactor(t *testing.T) {
	type args struct {
		imageBounds image.Rectangle
		height      int
		width       int
	}
	tests := []struct {
		name        string
		args        args
		want_size   int
		want_option rune
	}{
		{
			name: "no resize",
			args: args{
				imageBounds: image.Rect(0, 0, 100, 100),
				height:      200,
				width:       200,
			},
			want_size:   0,
			want_option: '0',
		},
		{
			name: "resize height",
			args: args{
				imageBounds: image.Rect(0, 0, 100, 100),
				height:      50,
				width:       200,
			},
			want_size:   50,
			want_option: 'y',
		},
		{
			name: "resize width",
			args: args{
				imageBounds: image.Rect(0, 0, 100, 100),
				height:      200,
				width:       50,
			},
			want_size:   50,
			want_option: 'x',
		},
		{
			name: "resize both x and y",
			args: args{
				imageBounds: image.Rect(0, 0, 100, 100),
				height:      50,
				width:       50,
			},
			want_size:   50,
			want_option: 'x',
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got_size, got_option := resizeConstraints(tt.args.imageBounds, tt.args.height, tt.args.width)
			if got_size != int(tt.want_size) || got_option != tt.want_option {
				t.Errorf("resizeConstraints() = %v, %s, want %v, %s", got_size, string(got_option), tt.want_size, string(tt.want_option))
			}
		})
	}
}
