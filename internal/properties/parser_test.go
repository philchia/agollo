package properties

import (
	"bytes"
	"testing"
)

func TestDocumentSetGet(t *testing.T) {

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "get set",
			args: args{
				key:   "key",
				value: "value",
			},
		},

		{
			name: "get set",
			args: args{
				key:   "key=",
				value: "value=",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.Set(tt.args.key, tt.args.value)

			if val, ok := p.Get(tt.args.key); !ok || val != tt.args.value {
				t.Error("set value not equal to got value")
			}
		})
	}
}

func TestSave(t *testing.T) {
	type args struct {
		kv map[string]string
	}
	tests := []struct {
		name        string
		args        args
		wantContent string
		wantErr     bool
	}{
		{
			name: "case1",
			args: args{
				kv: map[string]string{
					"timeout": "10",
				},
			},
			wantErr:     false,
			wantContent: "timeout=10\n",
		},

		{
			name: "case2",
			args: args{
				kv: map[string]string{
					"timeout=": "10",
				},
			},
			wantErr:     false,
			wantContent: "timeout\\==10\n",
		},

		{
			name: "case3",
			args: args{
				kv: map[string]string{
					"timeout=": "=10",
				},
			},
			wantErr:     false,
			wantContent: "timeout\\==\\=10\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			doc := New()
			for k, v := range tt.args.kv {
				doc.Set(k, v)
			}
			err := Save(doc, writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantContent {
				t.Errorf("Save() gotContent = %v, want %v", gotWriter, tt.wantContent)
			}
		})
	}
}
