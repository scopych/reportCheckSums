package main

// Reads csv file with field report and checks sums
import (
	//"bufio"
	"encoding/csv"
	"fmt"
	"io"
	//"log"
	"os"
	"regexp"
	"strconv"
)

type Servant struct {
	name          string
	placements    int
	videos        int
	hours         int
	returnVisits  int
	numBibleStuds int
	comments      string
	servsAs       string
}

type Overal struct {
	category      string
	quantity      int
	placements    int
	videos        int
	hours         int
	returnVisits  int
	numBibleStuds int
}

func main() {
	groups := make([][]Servant, 10)
	grp := make([]Servant, 25)

	var vozvOver Overal
	var podsobOver Overal
	var pionersOver Overal
	var regular Overal
	var neregular Overal

	structsOver := [5]*Overal{
		&vozvOver,
		&podsobOver,
		&pionersOver,
		&regular,
		&neregular,
	}

	// Open the file
	csvfile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))
	groupN := 0
	checkStart := false
	itogStart := false
	// Iterate through the records
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			//      log.Fatal(err)
			panic(err)
		}
		// Check if record where 1-t colon is number, 2-d colon is first name and last name
		fio, _ := regexp.MatchString(`[А-Я](\p{Cyrillic})+ +[А-Я](\p{Cyrillic})+`, record[2])
		n, _ := regexp.MatchString(`[[:digit:]]+`, record[1])
		if n && fio {
			grp = append(grp, filer(record))
		}
		vsego, _ := regexp.MatchString(`Всего`, record[2])
		if vsego {
			checkStart = true
			fmt.Println(record[2], "Группа №", groupN+1)
			continue
		}
		if checkStart {
			if record[1] == "" && record[2] == "" {
				groupN++
				groups = append(groups, grp)
				grp = grp[:0]
				checkStart = false
			} else {
				nums, servs := check(grp, record[2], structsOver)
				difprn(record, nums, servs)
			}
		}
		itog, _ := regexp.MatchString(`Общие цифры`, record[2])
		if itog {
			itogStart = true
			fmt.Println(record[2])
			continue
		}
		if itogStart {
			if record[1] == "" && record[2] == "" {
				println()
			} else {
				nums, servs := calculate(record[2], structsOver)
				difprn(record, nums, servs)
			}
		}
	}
}

func filer(info []string) (servant Servant) {
	name := info[2]
	placements, _ := strconv.Atoi(info[3])
	videos, _ := strconv.Atoi(info[4])
	hours, _ := strconv.Atoi(info[5])
	returnVisits, _ := strconv.Atoi(info[6])
	bibleStuds, _ := strconv.Atoi(info[7])
	comments := info[8]
	servsAs := info[9]

	servant = Servant{name, placements, videos, hours, returnVisits, bibleStuds, comments, servsAs}
	return servant
}

func check(grp []Servant, grpVozv string, structsOver [5]*Overal) (nums [5]int, servs int) {

	var nazn string
	var filStruc int

	switch grpVozv {
	case "Возвещатели":
		nazn = "Возвещатель"
		filStruc = 0
	case "Подсобные пионеры":
		nazn = "Подсобный пионер"
		filStruc = 1
	case "Пионеры":
		nazn = "Пионер"
		filStruc = 2
	}
	var placems, vids, hrs, returns, bibstuds int
	for i := 0; i < len(grp); i++ {
		if grp[i].servsAs == nazn {
			servs++
			placems += grp[i].placements
			vids += grp[i].videos
			hrs += grp[i].hours
			returns += grp[i].returnVisits
			bibstuds += grp[i].numBibleStuds
		}
	}
	nums = [5]int{placems, vids, hrs, returns, bibstuds}
	structsOver[filStruc].category = grpVozv
	structsOver[filStruc].quantity += servs
	structsOver[filStruc].placements += placems
	structsOver[filStruc].videos += vids
	structsOver[filStruc].hours += hrs
	structsOver[filStruc].returnVisits += returns
	structsOver[filStruc].numBibleStuds += bibstuds
	return nums, servs
}

func difprn(record []string, nums [5]int, servs int) {
	checked := make([]string, 0, 7)
	checked = append(checked, strconv.Itoa(servs))
	checked = append(checked, "")
	for _, val := range nums {
		checked = append(checked, strconv.Itoa(val))
	}
	rec := record[:8]
	for i, v := range rec {
		if i == 0 {
			continue
		}
		if i == 2 {
			fmt.Printf("%-17s%s", v, checked[i-1])
		} else if v != checked[i-1] {
			fmt.Printf("%10s/%-3s ", ("error=>" + v), checked[i-1])
		} else {
			fmt.Printf("%10s/%-3s ", v, checked[i-1])
		}
	}
	println()
}

func calculate(categorVozv string, structsOver [5]*Overal) (nums [5]int, quant int) {
	var filStruc int
	switch categorVozv {
	case "Возвещатели":
		filStruc = 0
	case "Подсобные пионеры":
		filStruc = 1
	case "Пионеры":
		filStruc = 2
	case "Все регулярные":
		filStruc = 3
	case "Нерегулярные":
		filStruc = 4
	}

	if filStruc == 0 || filStruc == 1 || filStruc == 2 {
		quant = structsOver[filStruc].quantity
		nums[0] = structsOver[filStruc].placements
		nums[1] = structsOver[filStruc].videos
		nums[2] = structsOver[filStruc].hours
		nums[3] = structsOver[filStruc].returnVisits
		nums[4] = structsOver[filStruc].numBibleStuds
	} else {
		quant = structsOver[0].quantity + structsOver[1].quantity + structsOver[2].quantity
		nums[0] = structsOver[0].placements + structsOver[1].placements + structsOver[2].placements
		nums[1] = structsOver[0].videos + structsOver[1].videos + structsOver[2].videos
		nums[2] = structsOver[0].hours + structsOver[1].hours + structsOver[2].hours
		nums[3] = structsOver[0].returnVisits + structsOver[1].returnVisits + structsOver[2].returnVisits
		nums[4] = structsOver[0].numBibleStuds + structsOver[1].numBibleStuds + structsOver[2].numBibleStuds

		*structsOver[filStruc] = Overal{
			categorVozv,
			quant,
			nums[0],
			nums[1],
			nums[2],
			nums[3],
			nums[4],
		}
	}
	return nums, quant
}
