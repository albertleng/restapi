package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var path = "ids.txt"

func CreateFile() {
	var _, err = os.Stat(path)
	
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer file.Close()
	}
	fmt.Println("Done creating file ", path)
}

func WriteFile(charIds []int) {
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	
	_, err = file.WriteString(getString(charIds))
	if err != nil {
		log.Fatal(err)
		return
	}
	
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
		return
	}
	
	fmt.Println("Done writing to file")
}

func ReadFile() (charIds []int, e error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	strs := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")

	ids := make([]int, len(strs) - 1)
	for i := range ids {
		ids[i], _ = strconv.Atoi(strs[i])
	}

	return ids, err
}

func getString(charIds []int) string {
	var str string
	for _, i := range charIds {
		str += fmt.Sprintf("%d\n", i)
	}
	return str
}