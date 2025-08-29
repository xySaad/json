another trash code. made it for the purpose of learning more about golang reflection, dynamic types, and json structure

# Json Parser

A Go package for parsing and manipulating JSON data. This package provides functionality to decode JSON strings into Go data structures and allows for easy retrieval of values using a path-like syntax.

### The package doesn't support structs yet, only maps

## Features

- Decode JSON strings into Go data structures (objects and arrays).

- Retrieve values from nested JSON objects using a simple path syntax.

- Handle errors gracefully with informative messages.

## Instalation

To use this package, you need to have Go installed on your machine. You can get the package by running:

```bash
go get github.com/xySaad/json-handler
```

## Usage

### Decoding JSON

To decode a JSON string, use the Decode function. It will return a getter interface that you can use to retrieve values. [see example](/json_test.go)

## Contributions

Contributions are welcome! If you find any bugs or have suggestions for improvements, feel free to open an issue or submit a pull request.
