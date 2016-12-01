package global

import (
	"encoding/json"
	"github.com/sydnash/lotou/log"
	"strconv"
)

type PropertySet struct {
	Property map[int]interface{}
}

func (s *PropertySet) setPropertyByType(ptype int, v interface{}) {
	_, ok := propertyToBase[ptype]
	if !ok {
		log.Error("property:%d is not exist", ptype)
		return
	}
	s.Property[ptype] = v
}
func (s *PropertySet) SetPropertyInt32(ptype int, v int32) {
	s.setPropertyByType(ptype, v)
}
func (s *PropertySet) SetPropertyInt64(ptype int, v int64) {
	s.setPropertyByType(ptype, v)
}
func (s *PropertySet) SetPropertyString(ptype int, v string) {
	s.setPropertyByType(ptype, v)
}

func (s *PropertySet) setProertyByKey(key string, v interface{}) {
	ptype, ok := keyToProperty[key]
	if !ok {
		log.Error("propertye key:%s is not exist.", key)
		return
	}
	s.setPropertyByType(ptype, v)
}
func (s *PropertySet) SetPropertyByKeyInt32(key string, v int32) {
	s.setProertyByKey(key, v)
}
func (s *PropertySet) SetPropertyByKeyInt64(key string, v int64) {
	s.setProertyByKey(key, v)
}
func (s *PropertySet) SetPropertyByKeyString(key string, v string) {
	s.setProertyByKey(key, v)
}

func (s *PropertySet) checkPropertyByType(ptype int) interface{} {
	v, ok := s.Property[ptype]
	if !ok {
		log.Error("property type:%d is not exist", ptype)
	}
	return v
}

func (s *PropertySet) GetPropertyInt32(ptype int) int32 {
	v := s.checkPropertyByType(ptype)
	if v == nil {
		return 0
	}
	return v.(int32)
}
func (s *PropertySet) GetPropertyInt64(ptype int) int64 {
	v := s.checkPropertyByType(ptype)
	if v == nil {
		return 0
	}
	return v.(int64)
}
func (s *PropertySet) GetPropertyString(ptype int) string {
	v := s.checkPropertyByType(ptype)
	if v == nil {
		return ""
	}
	return v.(string)
}

func (s *PropertySet) HasFlag(ptype int, flagType int) bool {
	v := s.GetPropertyInt32(ptype)
	return (v&int32(flagType) != 0)
}
func (s *PropertySet) ClearFlag(ptype int, flagType int) {
	v := s.GetPropertyInt32(ptype)
	v &= ^int32(flagType)
	s.SetPropertyInt32(ptype, v)
}
func (s *PropertySet) SetFlag(ptype int, flagType int) {
	v := s.GetPropertyInt32(ptype)
	v |= int32(flagType)
	s.SetPropertyInt32(ptype, v)
}

func NewPropertySet() *PropertySet {
	ret := PropertySet{}
	ret.Property = make(map[int]interface{})
	for k, v := range propertyToBase {
		ret.Property[k] = v.Value
	}
	return &ret
}

func (s *PropertySet) LoadJson(str string) {
	b := []byte(str)
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	if m != nil {
		for k, v := range m {
			ptype, ok := keyToProperty[k]
			if ok {
				base, ok := propertyToBase[ptype]
				if ok {
					switch base.ValueType {
					case KValueTypeInt32:
						str, ok := v.(string)
						if ok {
							t, err := strconv.ParseInt(str, 10, 32)
							if err != nil {
								log.Error("parse int faield:%s", err)
								continue
							}
							s.SetPropertyInt32(ptype, int32(t))
						}
					case KValueTypeInt64:
						str, ok := v.(string)
						if ok {
							t, err := strconv.ParseInt(str, 10, 64)
							if err != nil {
								log.Error("parse int faield:%s", err)
								continue
							}
							s.SetPropertyInt64(ptype, int64(t))
						}
					case KValueTypeString:
						switch v.(type) {
						case string:
							s.SetPropertyString(ptype, v.(string))
						default:
							nb, err := json.Marshal(v)
							if err == nil {
								s.SetPropertyString(ptype, string(nb))
							}
						}
					}
				}
			}
		}
	}
}
