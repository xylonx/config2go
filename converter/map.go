package converter

import (
	"reflect"
	"strings"

	"github.com/xylonx/config2go/util"
)

// FIXME: using it to let converter ignore field name
// using more elegant way to solve it.
const magicPrefix = "__$slice$__"

// MapParser - parse map[string]interface{} into Node Tree
type MapParser struct {
	Data map[string]interface{}
}

var _ Parser = &MapParser{}

func NewMapParser(data map[string]interface{}) Parser {
	return &MapParser{Data: data}
}

func (m *MapParser) ParseToNodeTree() (*Node, error) {
	fakeRoot := &Node{}
	if err := m.parse(fakeRoot); err != nil {
		return nil, err
	}
	return fakeRoot, nil
}

func (m *MapParser) parse(root *Node) (err error) {
	for k := range m.Data {
		vt := reflect.TypeOf(m.Data[k]).Kind()

		tag := k
		kk := k
		if strings.HasPrefix(k, magicPrefix) {
			tag = strings.Replace(k, magicPrefix, "", 1)
			k = ""
		}

		switch vt { // nolint:exhaustive
		// basic type
		case reflect.Bool:
			root.Child = append(root.Child, Node{TagKey: tag, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Bool})
		case reflect.String:
			root.Child = append(root.Child, Node{TagKey: tag, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.String})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			root.Child = append(root.Child, Node{TagKey: tag, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Int64})
		case reflect.Float32, reflect.Float64:
			root.Child = append(root.Child, Node{TagKey: tag, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Float64})

		case reflect.Map:
			v := m.Data[kk] //nolint:forcetypeassert
			var vv map[string]interface{}
			vv, ok := v.(map[string]interface{})
			if !ok {
				vvv, ok := v.(map[interface{}]interface{})
				if !ok {
					return ErrorUnsupportedMap
				}
				vv = mapInterface2MapString(vvv)
			}

			mn := Node{TagKey: tag, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Map}
			if err = (&MapParser{Data: vv}).parse(&mn); err != nil {
				return err
			}
			root.Child = append(root.Child, mn)

		case reflect.Slice:
			v := reflect.ValueOf(m.Data[kk])
			if v.Len() < 1 {
				continue
			}
			if err = validateSliceAlignment(root); err != nil {
				return err
			}

			sn := Node{TagKey: tag, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Slice}
			vv := v.Index(0).Interface()
			if err = (&MapParser{Data: map[string]interface{}{magicPrefix + k: vv}}).parse(&sn); err != nil {
				return err
			}
			root.Child = append(root.Child, sn)
		}
	}
	// switch m.Data

	return nil
}

func mapInterface2MapString(d map[interface{}]interface{}) map[string]interface{} {
	r := make(map[string]interface{})
	for k := range d {
		r[k.(string)] = d[k]
	}
	return r
}

// TODO: check slice type align in
func validateSliceAlignment(*Node) error {
	return nil
}
