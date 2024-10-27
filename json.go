package gosonify

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Json struct{}

type token struct {
	start rune
	end   rune
	depth int
	sep   rune
}

type (
	stateMap map[rune]*token
	Object   map[string]interface{}
)

func JsonDecoder() *Json {
	return &Json{}
}

func (s stateMap) init() stateMap {
	object := &token{start: '{', end: '}', depth: 0, sep: ','}
	array := &token{start: '[', end: ']', depth: 0, sep: ','}
	value := &token{start: ':', end: ',', depth: -1}

	return stateMap{
		'{': object,
		'}': object,
		'[': array,
		']': array,
		'"': &token{start: '"', end: '"', depth: -1},
		':': value,
		',': value,
	}
}

func (j Json) Decode(raw string) ([]Object, error) {

	arrayIndex := 0

	object, err := parseObject(raw, &arrayIndex)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func decoderHelper(state stateMap, stack *[]rune, depth *int, char rune, t *token, isLastChar bool, index int) error {
	if len(*stack) > 1 && state[(*stack)[len(*stack)-1]].end == ',' && state[(*stack)[len(*stack)-1]].end == '}' {
		*stack = (*stack)[:len(*stack)-1]
	}
	skip := false

	if len(*stack) > 0 && char == state[(*stack)[len(*stack)-1]].end {
		skip = true
	}

	if char == t.start && !skip {
		isUrl := len(*stack) > 1 && state[(*stack)[len(*stack)-1]].end == '"' && state[(*stack)[len(*stack)-2]].start == ':' && char == ':'
		if !isUrl {
			*stack = append(*stack, char)
			if t.depth != -1 {
				*depth++
				t.depth++
			}
		}
	} else if len(*stack) > 0 && state[(*stack)[len(*stack)-1]].sep == char {
	} else if char == t.end {
		if len(*stack) > 0 {
			if state[(*stack)[len(*stack)-1]].start != t.start {
				if state[(*stack)[len(*stack)-1]].end == ',' && isLastChar {
					*stack = (*stack)[:len(*stack)-1]

				} else if state[(*stack)[len(*stack)-1]].end == ',' && state[(*stack)[len(*stack)-2]].end == '}' {
					*stack = (*stack)[:len(*stack)-1]

				} else {
					return errors.New("expected2: " + string(state[(*stack)[len(*stack)-1]].end) + " found: " + string(char) + " index: " + strconv.Itoa(index))
				}
			}
			*stack = (*stack)[:len(*stack)-1]

		} else {
			return errors.New("expected: EOF" + " found: " + string(char))
		}
		if t.depth != -1 {
			*depth--
			t.depth--
		}
	}
	return nil
}

func createProperty(propName string, jMap *[]Object, arrayIndex *int) {
	if len((*jMap)) > *arrayIndex && len(propName) > 0 {
		_, ok := (*jMap)[*arrayIndex][propName[:len(propName)-1]]
		if ok {
			*arrayIndex++
		}
	}

	if len(propName) > 0 {
		if len((*jMap)) <= *arrayIndex {
			(*jMap) = append((*jMap), make(Object))
		}

		(*jMap)[*arrayIndex][propName[:len(propName)-1]] = nil
	}
}

func appendValue(propName string, value string, jMap *[]Object, arrayIndex *int) error {
	var result interface{}
	if value[0] == '[' && (value[len(value)-1] == '}' || (value[len(value)-2] == ']' && value[len(value)-1] == ',')) {
		result = parseArray(value[:len(value)-1])
	} else if value[0] == '"' && (value[len(value)-1] == '}' || (value[len(value)-2] == '"' && value[len(value)-1] == ',')) {
		value = value[1 : len(value)-2]
		result = value
	} else if value[0] == '{' && (value[len(value)-1] == '}' || (value[len(value)-2] == '"' && value[len(value)-1] == ',')) {
		object, err := parseObject(value, arrayIndex)
		if err != nil {
			return err
		}
		result = object
	} else {
		value = value[:len(value)-1]
		if value == "true" || value == "false" || value == "null" {
			result = value
		} else {
			num, err := strconv.Atoi(value)
			if err != nil {
				return errors.New("invalid value type: " + value)
			}
			result = num
		}
	}

	if len(propName) > 0 {
		(*jMap)[*arrayIndex][propName[:len(propName)-1]] = result
	}
	return nil
}

func parseArray(str string) []string {
	array := []string{}
	index := 0
	rStr := []rune(str)
	for i := 1; i < len(rStr)-1; i++ {
		char := rStr[i]
		if char == ',' {
			index++
			continue
		}
		if char == ' ' {
			continue
		}
		if len(array)-1 < index {
			array = append(array, string(char))
		} else {
			array[index] += string(char)
		}
	}

	return array
}

func parseObject(raw string, arrayIndex *int) ([]Object, error) {
	rawR := []rune(raw)
	state := stateMap{}.init()
	stack := []rune{}
	index := 0
	depth := 0
	result := make([]Object, 0)
	property := ""
	var value string
	inProp := false
	inValue := false
	inArray := false
	prevProp := ""

	for index < len(rawR) {
		char := rawR[index]
		if inProp && !inValue {
			property += string(char)
		} else if inValue {
			value += string(char)
		}
		if char == '{' {
		} else if char == '"' {
			if !inProp {
				inProp = true
			} else {
				createProperty(property, &result, arrayIndex)
				if !inValue {
					prevProp = property
					property = ""
				}
				inProp = false
			}
		} else if char == ':' {
			inValue = true
		} else if char == ',' && (stack)[len(stack)-1] != '{' {
			if !inArray || (inArray && (stack)[len(stack)-1] == ':') {
				inValue = false
				if len(value) > 1 {
					err := appendValue(prevProp, strings.TrimSpace(value), &result, arrayIndex)
					if err != nil {
						return nil, err
					}
				} else {
				}
				value = ""
			}
		} else if char == '[' {
			inArray = true
		} else if char == ']' {
			inArray = false
		} else if char == '}' {
			inValue = false
			if len(value) > 1 {
				err := appendValue(prevProp, strings.TrimSpace(value), &result, arrayIndex)
				if err != nil {
					return nil, err
				}
			} else {
			}
			value = ""
		}

		if index > 0 && char == ',' && char == rawR[index-1] {
			return nil, errors.New("expected1: " + string(state[(stack)[len(stack)-1]].end) + " found: " + string(char))
		}
		t, exist := state[char]
		if exist {
			err := decoderHelper(state, &stack, &depth, char, t, index == len(raw)-1, index)
			if err != nil {
				return result, err
			}
		}
		index++
	}

	if depth != 0 {
		return nil, fmt.Errorf("mismatched brackets: depth is %d at the end of parsing, expected: %q", depth, state[stack[len(stack)-1]].end)
	}
	return result, nil
}
