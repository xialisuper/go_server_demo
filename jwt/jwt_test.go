package jwt

// func TestCreateJwtToken(t *testing.T) {
// 	type args struct {
// 		userId          string
// 		expireTimeInSec int64
// 		jwtSecret       string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "TestCreateJwtToken",
// 			args: args{
// 				userId:          "123456",
// 				expireTimeInSec: 10,
// 				jwtSecret:       "secret",
// 			},
// 			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEyMzQ1NiwiZXhwIjoxNjM0OTY3NzIwfQ.9270373117",
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := CreateJwtToken(tt.args.userId, tt.args.expireTimeInSec, tt.args.jwtSecret)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CreateJwtToken() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("CreateJwtToken() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
