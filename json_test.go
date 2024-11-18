package json_test

import (
	"fmt"
	"io"
	"net/http"

	"github.com/xySaad/json"
)

func ExampleDecode() {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	rawJson, err := io.ReadAll(resp.Body)
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
	var ArtistN5 json.Object
	// use the Get method to store the data in the holder
	err = jsonData.Get(&ArtistN5, "[4]")
	if err != nil {
		panic(err)
	}

	fmt.Println(ArtistN5)

	// You access any key of the artist object easily by indexing it.
	// Example accessing members:
	fmt.Println(ArtistN5["members"])

	// You also access a value from a nested JSON directly.
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
