package jwt

import "testing"

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{
			name: "Happy path",
			n:    10,
			want: 20,
		},
		{
			name: "Edge case: n=0",
			n:    0,
			want: 0,
		},
		{
			name: "Edge case: n=1",
			n:    1,
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			str, _ := generateRandomString(tt.n)
			if got := len(str); got != tt.want {
				t.Errorf("GenerateRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}
