package utils

import (
	"errors"
	"fmt"
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
type jsonMap map[string]interface{}

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

func (j Json) Decode(raw string) (jsonMap, error) {
	rawR := []rune(raw)
	state := stateMap{}.init()
	stack := []rune{}
	index := 0
	depth := 0
	var result jsonMap = make(jsonMap)
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
			fmt.Println(prevProp, value)
			appendValue(prevProp, value[:len(value)-1], &result)
			value = ""
		} else if char == '[' {
			inArray = true
		} else if char == ']' {
			inArray = false
		}

		if index > 0 && char == ',' && char == rawR[index-1] {
			return nil, errors.New("expected: " + string(state[(stack)[len(stack)-1]].end) + " found: " + string(char))
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
	fmt.Println(value)
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
		fmt.Println("this is a seperator:", char)
	} else if char == t.end {
		if len(*stack) > 0 {
			if state[(*stack)[len(*stack)-1]].start != t.start {
				if state[(*stack)[len(*stack)-1]].end == ',' && isLastChar {
					fmt.Println("skip last ,")
				} else {
					return errors.New("expected: " + string(state[(*stack)[len(*stack)-1]].end) + " found: " + string(char))
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

func createProperty(propName string, jMap *jsonMap) {
	if len(propName) > 0 {
		(*jMap)[propName[:len(propName)-1]] = nil
	}
}

func appendValue(propName string, value any, jMap *jsonMap) {
	if len(propName) > 0 {
		(*jMap)[propName[:len(propName)-1]] = value
	}
}
