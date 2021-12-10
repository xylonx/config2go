package converter

import (
	"reflect"

	"github.com/xylonx/config2go/util"
)

type MapConverter struct {
	Data map[string]interface{}
}

var _ Converter = &MapConverter{}

func NewMapConverter(data map[string]interface{}) Converter {
	return &MapConverter{Data: data}
}

func (m *MapConverter) ConvertToNodeTree() (*Node, error) {
	fakeRoot := &Node{}
	m.convertConfig2Tree(fakeRoot)
	return fakeRoot, nil
}

func (m *MapConverter) convertConfig2Tree(root *Node) {
	config := m.Data
	for k := range config {
		vt := reflect.TypeOf(config[k]).Kind()
		switch vt { // nolint:exhaustive
		// basic type
		case reflect.Bool:
			root.Child = append(root.Child, Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Bool})
		case reflect.String:
			root.Child = append(root.Child, Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.String})
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			root.Child = append(root.Child, Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Int64})
		case reflect.Float32, reflect.Float64:
			root.Child = append(root.Child, Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Float64})

		// Map
		case reflect.Map:
			v := config[k].(map[string]interface{}) //nolint:forcetypeassert
			sn := Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Map}
			(&MapConverter{Data: v}).convertConfig2Tree(&sn)
			root.Child = append(root.Child, sn)

		// Slice
		case reflect.Slice:
			v := reflect.ValueOf(config[k])
			if v.Len() < 1 {
				continue
			}
			sn := Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Slice}
			switch v.Index(0).Interface().(type) {
			case map[interface{}]interface{}:
				tv := v.Index(0).Interface().(map[interface{}]interface{}) // nolint:forcetypeassert
				vv := make(map[string]interface{}, len(tv))
				for tvk := range tv {
					vv[tvk.(string)] = tv[tvk]
				}
				mn := Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Map}
				(&MapConverter{Data: vv}).convertConfig2Tree(&mn)
				sn.Child = append(sn.Child, mn)
			case bool:
				sn.Child = append(sn.Child, Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Bool})
			case string:
				sn.Child = append(sn.Child, Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.String})
			case int:
				sn.Child = append(sn.Child, Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Int64})
			case float32:
				sn.Child = append(sn.Child, Node{TagKey: k, FieldName: util.ConvertString2UpperCamel(k), Type: reflect.Float64})
			}
			root.Child = append(root.Child, sn)
		}
	}
}
