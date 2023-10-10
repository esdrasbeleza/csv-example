package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Person struct {
	ID        int
	Name      string
	Surname   string
	AddressId int
	Children  []Child
}

type Child struct {
	ID       int
	ParentID int
	Name     string
	Surname  string
}

func main() {
	peopleFile, err := os.Open("csv1_people.csv")
	if err != nil {
		panic(err)
	}
	defer peopleFile.Close()

	addressesFile, err := os.Open("csv2_addresses.csv")
	if err != nil {
		panic(err)
	}
	defer addressesFile.Close()

	childrenFile, err := os.Open("csv3_children.csv")
	if err != nil {
		panic(err)
	}
	defer childrenFile.Close()

	people := readPeopleCSV(peopleFile)
	addresses := readAddressesCSV(addressesFile)
	children := readChildrenCSV(childrenFile)

	// Create a map to store children by parent ID
	childrenMap := make(map[int][]Child)
	for _, child := range children {
		childrenMap[child.ParentID] = append(childrenMap[child.ParentID], child)
	}

	// Print requested information
	for _, person := range people {
		fmt.Printf("1. Person Full Name: %s %s\n", person.Name, person.Surname)

		address := addresses[person.AddressId]
		fmt.Printf("2. Person Address: %s\n", address)

		fmt.Printf("3. Full Names of Their Children:\n")
		if childNames, ok := getChildNames(childrenMap[person.ID]); ok {
			for _, name := range childNames {
				fmt.Printf("   - %s\n", name)
			}
		} else {
			fmt.Println("   - None")
		}

		fmt.Println()
	}
}

func readPeopleCSV(file *os.File) []Person {
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var people []Person
	for i, line := range lines {
		if i == 0 { // skip first line
			continue
		}
		id, _ := strconv.Atoi(line[0])        // TODO: handle error
		addressId, _ := strconv.Atoi(line[3]) // TODO: handle error
		people = append(people, Person{
			ID:        id,
			Name:      line[1],
			Surname:   line[2],
			AddressId: addressId,
		})
	}

	return people
}

func readAddressesCSV(file *os.File) map[int]string {
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	addresses := make(map[int]string)
	for i, line := range lines {
		if i == 0 { // skip first line
			continue
		}
		id, _ := strconv.Atoi(line[0])
		addresses[id] = line[1] + " " + line[2] + ", " + line[3] + ", " + line[4]
	}

	return addresses
}

func readChildrenCSV(file *os.File) []Child {
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var children []Child
	for i, line := range lines {
		if i == 0 { // skip first line
			continue
		}
		id, _ := strconv.Atoi(line[0])
		parentID, _ := strconv.Atoi(line[1])
		children = append(children, Child{
			ID:       id,
			ParentID: parentID,
			Name:     line[2],
			Surname:  line[3],
		})
	}

	return children
}

func getChildNames(children []Child) ([]string, bool) {
	if len(children) == 0 {
		return nil, false
	}

	var names []string
	for _, child := range children {
		names = append(names, child.Name+" "+child.Surname)
	}
	return names, true
}
