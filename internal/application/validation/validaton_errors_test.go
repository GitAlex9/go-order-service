package validation

import "testing"

func TestValidationErrors_Add(t *testing.T) {
	verr := &ValidationErrors{}

	verr.Add("name", "cannot be empty")
	verr.Add("email", "invalid format")

	if got, want := len(verr.Errors), 2; got != want {
		t.Fatalf("len(Errors) got = %d, want %d", got, want)
	}
	if got, want := verr.Errors[0].Field, "name"; got != want {
		t.Errorf("Errors[0].Field got = %q, want %q", got, want)
	}
}

func TestValidationErrors_HasErrors(t *testing.T) {
	tests := []struct {
		name string
		verr *ValidationErrors
		want bool
	}{
		{"sem erros", &ValidationErrors{}, false},
		{"com erros", func() *ValidationErrors {
			v := &ValidationErrors{}
			v.Add("field", "message")
			return v
		}(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.verr.HasErrors()
			if got != tt.want {
				t.Errorf("HasErrors() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	verr := &ValidationErrors{}
	verr.Add("name", "cannot be empty")
	verr.Add("email", "invalid format")

	got := verr.Error()
	want := "name: cannot be empty; email: invalid format"

	if got != want {
		t.Errorf("Error() got = %q, want %q", got, want)
	}
}
