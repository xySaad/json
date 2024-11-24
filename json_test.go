package json_test

import (
	"fmt"
	"os"

	"github.com/xySaad/json"
)

func ExampleDecode() {

	rawJson, err := os.ReadFile("./example.json")
	if err != nil {
		panic(err)
	}

	// Compile JSON and retrieve a Getter
	jsonData, err := json.Decode(string(rawJson))
	if err != nil {
		panic(err)
	}

	// Retrieve the 5th artist
	// The expected data is an Object, define a variable to hold the value.
	var ArtistN5 map[string]any
	// use the Get method to store the data in the holder
	err = jsonData.Get(&ArtistN5, "[4]")
	if err != nil {
		panic(err)
	}

	fmt.Println(ArtistN5)

	// You can access any key of the artist object easily by indexing it.
	// Example accessing members:
	fmt.Println(ArtistN5["members"])

	// You can also access a value from a nested JSON directly.
	// Example accessing the id of the 9th artist:
	var id int
	err = jsonData.Get(&id, "[8].id")
	if err != nil {
		panic(err)
	}

	fmt.Println(id)
	// Output:
	// map[concertDates:https://groupietrackers.herokuapp.com/api/dates/5 creationDate:2013 firstAlbum:07-09-2017 id:5 image:https://groupietrackers.herokuapp.com/api/images/xxxtentacion.jpeg locations:https://groupietrackers.herokuapp.com/api/locations/5 members:[Jahseh Dwayne Ricardo Onfroy] name:XXXTentacion relations:https://groupietrackers.herokuapp.com/api/relation/5]
	// [Jahseh Dwayne Ricardo Onfroy]
	// 9
}
