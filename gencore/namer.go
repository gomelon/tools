package gencore

import (
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
	"path/filepath"
	"strconv"
	"strings"
)

type rawNamer struct {
	pkg     string
	tracker namer.ImportTracker
	namer.Names
}

func NewRawNamer(pkg string, tracker namer.ImportTracker) *rawNamer {
	return &rawNamer{pkg: pkg, tracker: tracker}
}

// Name makes a name the way you'd write it to literally refer to type t,
// making ordinary assumptions about how you've imported t's package (or using
// r.tracker to specifically track the package imports).
func (r *rawNamer) Name(t *types.Type) string {
	if r.Names == nil {
		r.Names = namer.Names{}
	}
	if name, ok := r.Names[t]; ok {
		return name
	}
	if t.Name.Package != "" {
		var name string
		if r.tracker != nil {
			r.tracker.AddType(t)
			if t.Name.Package == r.pkg {
				name = t.Name.Name
			} else {
				name = r.tracker.LocalNameOf(t.Name.Package) + "." + t.Name.Name
			}
		} else {
			if t.Name.Package == r.pkg {
				name = t.Name.Name
			} else {
				name = filepath.Base(t.Name.Package) + "." + t.Name.Name
			}
		}
		r.Names[t] = name
		return name
	}
	var name string
	switch t.Kind {
	case types.Builtin:
		name = t.Name.Name
	case types.Map:
		name = "map[" + r.Name(t.Key) + "]" + r.Name(t.Elem)
	case types.Slice:
		name = "[]" + r.Name(t.Elem)
	case types.Array:
		l := strconv.Itoa(int(t.Len))
		name = "[" + l + "]" + r.Name(t.Elem)
	case types.Pointer:
		name = "*" + r.Name(t.Elem)
	case types.Struct:
		elems := []string{}
		for _, m := range t.Members {
			elems = append(elems, m.Name+" "+r.Name(m.Type))
		}
		name = "struct{" + strings.Join(elems, "; ") + "}"
	case types.Chan:
		// TODO: include directionality
		name = "chan " + r.Name(t.Elem)
	case types.Interface:
		// TODO: add to name test
		if t.Name.Name == "error" {
			name = "error"
		} else {
			var elems []string
			for _, m := range t.Methods {
				// TODO: include function signature
				elems = append(elems, m.Name.Name)
			}
			name = "interface{" + strings.Join(elems, "; ") + "}"
		}
	case types.Func:
		// TODO: add to name test
		var params []string
		for _, pt := range t.Signature.Parameters {
			params = append(params, r.Name(pt))
		}
		results := []string{}
		for _, rt := range t.Signature.Results {
			results = append(results, r.Name(rt))
		}
		name = "func(" + strings.Join(params, ",") + ")"
		if len(results) == 1 {
			name += " " + results[0]
		} else if len(results) > 1 {
			name += " (" + strings.Join(results, ",") + ")"
		}
	default:
		name = "unnameable_" + string(t.Kind)
	}
	r.Names[t] = name
	return name
}
