package json

import (
	"errors"
	"fmt"
	"strings"
)

type Json struct{}

type token struct {
	start rune
	end   rune
	depth int
	sep   rune
}

/*
34 "
58 ,
*/

type stateMap map[rune]*token
type Object map[string]interface{}

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

func (j Json) Decode(raw string) (Object, error) {
	rawR := []rune(raw)
	state := stateMap{}.init()
	stack := []rune{}
	index := 0
	depth := 0
	var result Object = make(Object)
	property := ""
	var value string
	inProp := false
	inValue := false
	inArray := false
	prevProp := ""

	for index < len(rawR) {
		char := rawR[index]
		if inArray {
			// store array
		}
		if inProp && !inValue {
			property += string(char)
		} else if inValue {
			value += string(char)
		}

		if char == '"' {
			if !inProp {
				inProp = true
			} else {
				createProperty(property, &result)
				if !inValue {
					prevProp = property
					property = ""
				}
				inProp = false
			}
		} else if char == ':' {
			inValue = true
		} else if char == ',' && !inArray && (stack)[len(stack)-depth] != '{' {
			inValue = false
			if len(value) > 2 {
				err := appendValue(prevProp, strings.TrimSpace(value), &result)
				if err != nil {
					return nil, err
				}
			}
			value = ""
		} else if char == '[' {
			inArray = true
		} else if char == ']' {
			inArray = false
		} else if char == '}' {
			inValue = false
			if len(value) > 2 {
				err := appendValue(prevProp, strings.TrimSpace(value), &result)
				if err != nil {
					return nil, err
				}
			}
			value = ""
		}

		if index > 0 && char == ',' && char == rawR[index-1] {
			return nil, errors.New("expected1: " + string(state[(stack)[len(stack)-1]].end) + " found: " + string(char))
		}
		t, exist := state[char]
		if exist {
			err := decoderHelper(state, &stack, &depth, char, t, index == len(raw)-1)
			if err != nil {
				return nil, err
			}
		}
		index++
	}

	if depth != 0 {
		return nil, fmt.Errorf("mismatched brackets: depth is %d at the end of parsing, expected: %q", depth, state[stack[len(stack)-1]].end)
	}
	// fmt.Println(value)
	return result, nil
}

func decoderHelper(state stateMap, stack *[]rune, depth *int, char rune, t *token, isLastChar bool) error {
	/* debuging

	fmt.Println(*stack)
	fmt.Print("character: ", string(char), " token depth: ", t.depth, " stack length: ", len(*stack), " ")

	if len(*stack) > 0 {
		fmt.Println("token in stack:", string(state[(*stack)[len(*stack)-1]].start), string(state[(*stack)[len(*stack)-1]].end))
	} else {
		fmt.Println("")
	}

	*/

	if len(*stack) > 1 && state[(*stack)[len(*stack)-1]].end == ',' && state[(*stack)[len(*stack)-1]].end == '}' {
		*stack = (*stack)[:len(*stack)-1]

	}
	skip := false

	if len(*stack) > 0 && char == state[(*stack)[len(*stack)-1]].end {
		skip = true
	}

	if char == t.start && !skip {
		*stack = append(*stack, char)
		if t.depth != -1 {
			*depth++
			t.depth++
		}
	} else if len(*stack) > 0 && state[(*stack)[len(*stack)-1]].sep == char {
		// fmt.Println("this is a seperator:", char)
	} else if char == t.end {
		if len(*stack) > 0 {
			if state[(*stack)[len(*stack)-1]].start != t.start {
				if state[(*stack)[len(*stack)-1]].end == ',' && isLastChar {
					// fmt.Println("skip last ,")
				} else if state[(*stack)[len(*stack)-1]].end == ',' && state[(*stack)[len(*stack)-2]].end == '}' {
					// fmt.Println("prev object closed yet")
				} else {
					return errors.New("expected2: " + string(state[(*stack)[len(*stack)-1]].end) + " found: " + string(char))
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
	// for debuggin
	// fmt.Printf("Token Start: %s, End: %s, TokenDepth: %d, Depth: %d\n", string(t.start), string(t.end), t.depth, *depth)
	return nil
}

func createProperty(propName string, jMap *Object) {
	if len(propName) > 0 {
		(*jMap)[propName[:len(propName)-1]] = nil
	}
}

func appendValue(propName string, value string, jMap *Object) error {
	var result interface{}
	if value[0] == '[' && (value[len(value)-1] == '}' || (value[len(value)-2] == ']' && value[len(value)-1] == ',')) {
		result = parseArray(value[:len(value)-1])
	} else if value[0] == '"' && (value[len(value)-1] == '}' || (value[len(value)-2] == '"' && value[len(value)-1] == ',')) {
		value = value[1 : len(value)-2]
		result = value
	} else if value[0] == '{' && (value[len(value)-1] == '}' || (value[len(value)-2] == '"' && value[len(value)-1] == ',')) {
		object, err := parseObject(value)
		if err != nil {
			fmt.Println(value)
			return err
		}
		result = object
	} else {
		fmt.Println("invalid value type,", value)
	}

	if len(propName) > 0 {
		(*jMap)[propName[:len(propName)-1]] = result
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

func parseObject(raw string) (Object, error) {
	result := make(Object)
	rawR := []rune(raw)
	state := stateMap{}.init()
	stack := []rune{}
	index := 0
	depth := 0
	property := ""
	var value string
	inProp := false
	inValue := false
	inArray := false
	prevProp := ""

	for index < len(rawR) {
		char := rawR[index]
		if inArray {
			// store array
		}
		if inProp && !inValue {
			property += string(char)
		} else if inValue {
			value += string(char)
		}

		if char == '"' {
			if !inProp {
				inProp = true
			} else {
				createProperty(property, &result)
				if !inValue {
					prevProp = property
					property = ""
				}
				inProp = false
			}
		} else if char == ':' {
			inValue = true
		} else if char == ',' && !inArray {
			inValue = false
			if len(value) > 2 {
				err := appendValue(prevProp, strings.TrimSpace(value), &result)
				if err != nil {
					return nil, err
				}
			}
			value = ""
		} else if char == '[' {
			inArray = true
		} else if char == ']' {
			inArray = false
		} else if char == '}' {
			inValue = false
			if len(value) > 2 {
				err := appendValue(prevProp, strings.TrimSpace(value), &result)
				if err != nil {
					return nil, err
				}
			}
			value = ""
		}

		if index > 0 && char == ',' && char == rawR[index-1] {
			return nil, errors.New("expected1: " + string(state[(stack)[len(stack)-1]].end) + " found: " + string(char))
		}
		t, exist := state[char]
		if exist {
			err := decoderHelper(state, &stack, &depth, char, t, index == len(raw)-1)
			if err != nil {
				return nil, err
			}
		}
		index++
	}

	if depth != 0 {
		return nil, fmt.Errorf("mismatched brackets: depth is %d at the end of parsing, expected: %q", depth, state[stack[len(stack)-1]].end)
	}

	return result, nil
}
