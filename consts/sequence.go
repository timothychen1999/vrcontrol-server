package consts

import (
	"encoding/json"
	"log"
	"os"
)

func LoadAssignedSequence(room string) map[string]int {
	//Check if room file exists
	//If not, set it to empty map
	//If yes, load the map from file

	path := "sequence/" + room + ".json"
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return make(map[string]int)
	}
	// Load the map from file
	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file: ", err)
		return make(map[string]int)
	}
	defer file.Close()
	var sequenceMap map[string]int
	err = json.NewDecoder(file).Decode(&sequenceMap)
	if err != nil {
		log.Println("Error decoding file: ", err)
		return make(map[string]int)
	}
	return sequenceMap

}
