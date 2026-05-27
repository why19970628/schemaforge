package generator

import (
	"strings"
	"testing"
)

func TestJSONFromXMLAttributesAndRepeatedChildren(t *testing.T) {
	out, err := JSONFromXML(`<users org="core"><user>Ada</user><user>Linus</user></users>`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"@org": "core"`, `"user": [`, `"Ada"`, `"Linus"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
}
