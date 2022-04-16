package gencore

//func TestTagParser(t *testing.T) {
//	type args struct {
//		namespace string
//		name      string
//		comments  []string
//	}
//	tests := []struct {
//		name         string
//		args         args
//		key          string
//		wantedGet    string
//		wantedLookup bool
//		wantErr      bool
//	}{
//		{
//			name: "positive,newline",
//			args: args{
//				namespace: "+gomelon",
//				name:      "Query",
//				comments:  []string{"FindByName", "+gomelon.Query `", "sql:\"select 1 from dual\"`"},
//			},
//			key:          "sql",
//			wantedGet:    "select 1 from dual",
//			wantedLookup: true,
//			wantErr:      false,
//		},
//		{
//			name: "positive,inline",
//			args: args{
//				namespace: "+gomelon",
//				name:      "Query",
//				comments:  []string{"FindByName", "+gomelon.Query `sql:\"select 1 from dual\"`"},
//			},
//			key:          "sql",
//			wantedGet:    "select 1 from dual",
//			wantedLookup: true,
//			wantErr:      false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			parser := NewTagParser()
//
//			get := parser.Parse(tt.args.namespace, tt.args.name, tt.args.comments)
//			if get != tt.wantedGet {
//				t.Errorf("ParseTagAnnotation() got.Get(%v) = %v, want %v", tt.key, get, tt.wantedGet)
//			}
//			_, lookup := parser.Lookup(tt.key)
//			if lookup != tt.wantedLookup {
//				t.Errorf("ParseTagAnnotation() got.Lookup(%v) = _,%v, want %v", tt.key, lookup, tt.wantedLookup)
//			}
//		})
//	}
//}
//
//type TestAnnotation struct {
//	DefaultValue int `default:"1"`
//	NoneValue    int
//	IntValue     int
//	StringValue  string
//	BoolValue    bool
//	Float32Value float32
//	Float64Value float64
//}
//
//func TestTagParser_Parse(t1 *testing.T) {
//	type fields struct {
//		namespace string
//		name      string
//		has       bool
//		tag       reflect.StructTag
//	}
//	type args struct {
//		annotation interface{}
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "positive",
//			fields: fields{
//				namespace: "+gomelon",
//				name:      "TestAnnotation",
//				has:       true,
//				tag:       `NoneValue:"2" IntValue:"1" StringValue:"a" BoolValue:"true" Float32Value:"3.14" Float64Value:"3.14"`,
//			},
//			args: args{
//				annotation: &TestAnnotation{
//					DefaultValue: 1,
//					NoneValue:    2,
//					IntValue:     1,
//					StringValue:  "a",
//					BoolValue:    true,
//					Float32Value: 3.14,
//					Float64Value: 3.14,
//				},
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t1 *testing.T) {
//			t := &TagAnnotationParser{
//				namespace: tt.fields.namespace,
//				name:      tt.fields.name,
//				has:       tt.fields.has,
//				tag:       tt.fields.tag,
//			}
//			ta := &TestAnnotation{}
//			err := t.Parse(ta)
//
//			if (err != nil) != tt.wantErr {
//				t1.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
//			}
//
//			if err == nil && !reflect.DeepEqual(ta, tt.args.annotation) {
//				t1.Errorf("Parse() result = %#v, wanted %#v", ta, tt.args.annotation)
//			}
//		})
//	}
//}
