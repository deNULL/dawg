package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/deNULL/dawg"
)

var optTab bool
var optGuide bool
var optRanked bool
var optBuild bool
var optSort bool
var optUtfc bool

var optLexicon string
var optDictionary string

var fileLexicon *os.File
var fileDictionary *os.File
var scanner *bufio.Scanner

var builder *dawg.DawgBuilder

func processLine(line string) bool {
	if len(line) == 0 {
		return false
	}

	//fmt.Printf("inserting %s: %v\n", line, Encode(line))

	var key string = line
	var val int32 = 0
	if optTab {
		parts := strings.Split(line, "\t")
		key = parts[0]
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal(err)
		}
		if value < 0 {
			fmt.Printf("warning: negative value is replaced by 0: %d\n", value)
			value = 0
		} else if value > dawg.MaxValue {
			fmt.Printf("warning: too large value is replaced by %d\n", dawg.MaxValue)
			value = dawg.MaxValue
		}
		val = int32(value)
	}

	var bytes []uint8
	if optUtfc {
		bytes = dawg.UtfcEncode(key)
	} else {
		bytes = []uint8(key)
	}
	if !builder.InsertKeyValue(bytes, len(bytes), val) {
		//fmt.Printf("key %s: %v", line, bytes)
		log.Fatalf("error: failed to insert key: %s\n", key)
	}

	return true
}

func handleBuildDict() {
	builder = dawg.NewDawgBuilder()

	var keyCount int = 0
	if optSort {
		var buffer []string = []string{}
		for scanner.Scan() {
			buffer = append(buffer, scanner.Text())
		}
		sort.SliceStable(buffer, func(i int, j int) bool {
			if !optUtfc {
				return buffer[i] < buffer[j]
			}
			bytes1 := dawg.UtfcEncode(buffer[i])
			bytes2 := dawg.UtfcEncode(buffer[j])
			length := len(bytes1)
			if len(bytes2) < length {
				length = len(bytes2)
			}
			for i := 0; i < length; i++ {
				if bytes1[i] < bytes2[i] {
					return true
				} else if bytes1[i] > bytes2[i] {
					return false
				}
			}
			if len(bytes1) < len(bytes2) {
				return true
			}
			return false
		})

		for i := 0; i < len(buffer); i++ {
			if processLine(buffer[i]) {
				keyCount++
				if keyCount%10000 == 0 {
					fmt.Printf("no. keys: %d\n", keyCount)
				}
			}
		}
	} else {
		for scanner.Scan() {
			var line string = scanner.Text()
			if processLine(line) {
				keyCount++
				if keyCount%10000 == 0 {
					fmt.Printf("no. keys: %d\n", keyCount)
				}
			}
		}
	}

	d := dawg.NewDawg()
	builder.Finish(d)

	fmt.Printf("no. keys: %d\n", keyCount)
	fmt.Printf("no. states: %d\n", d.NumOfStates())
	fmt.Printf("no. transitions: %d\n", d.NumOfTransitions())
	fmt.Printf("no. merged states: %d\n", d.NumOfMergedStates())
	fmt.Printf("no. merging states: %d\n", d.NumOfMergingStates())
	fmt.Printf("no. merged transitions: %d\n", d.NumOfMergedTransitions())

	var numOfUnusedUnits uint32
	dict := d.BuildWithUnused(&numOfUnusedUnits)
	if dict == nil {
		log.Fatalf("error: failed to build dictionary\n")
	}

	var unusedRatio float64 = 100.0 * float64(numOfUnusedUnits) / float64(dict.Size())

	fmt.Printf("no. elements: %d\n", dict.Size())
	fmt.Printf("no. unused elements: %d (%.2f%%)\n", numOfUnusedUnits, unusedRatio)
	fmt.Printf("dictionary size: %d\n", dict.TotalSize())

	/*for i := 0; i < dict.size; i++ {
		fmt.Printf("%.2d: %.8x leaf? %v ext? %v hleaf? %v label? %d (%c) offset = %d => %d, value = %d\n", i, dict.units[i], dictIsLeaf(dict.units[i]), dictHasExtBit(dict.units[i]), dictHasLeaf(dict.units[i]), dictLabel(dict.units[i]), dictLabel(dict.units[i]), dictOffset(dict.units[i]), dictOffset(dict.units[i])^baseType(i), dictValue(dict.units[i]))
	}*/

	if !dict.Write(fileDictionary) {
		log.Fatalf("error: failed to write Dictionary")
	}

	// Builds a guide
	if optRanked {
		guide := dawg.BuildRankedGuide(d, dict)
		if guide == nil {
			log.Fatalf("error: failed to build Guide\n")
		}

		fmt.Printf("no. units: %d\n", guide.Size())
		fmt.Printf("guide size: %d\n", guide.TotalSize())

		if !guide.Write(fileDictionary) {
			log.Fatalf("error: failed to write RankedGuide\n")
		}
	} else if optGuide {
		guide := dawg.BuildGuide(d, dict)
		if guide == nil {
			log.Fatalf("error: failed to build Guide\n")
		}

		fmt.Printf("no. units: %d\n", guide.Size())
		fmt.Printf("guide size: %d\n", guide.TotalSize())

		if !guide.Write(fileDictionary) {
			log.Fatalf("error: failed to write Guide\n")
		}
	}
}

func handleLoadDict() {
	dict := dawg.ReadDictionary(fileDictionary)
	if dict == nil {
		log.Fatalf("error: failed to read Dictionary\n")
	}

	var completer dawg.SomeCompleter
	if optRanked {
		guide := dawg.ReadRankedGuide(fileDictionary)
		if guide == nil {
			log.Fatalf("error: failed to read RankedGuide\n")
		}
		completer = dawg.NewRankedCompleter(dict, guide)
	} else if optGuide {
		guide := dawg.ReadGuide(fileDictionary)
		if guide == nil {
			log.Fatalf("error: failed to read Guide\n")
		}
		completer = dawg.NewCompleter(dict, guide)
	}

	for scanner.Scan() {
		var key string = scanner.Text()

		fmt.Printf("%s:", key)
		var index uint32 = dict.Root()

		if optRanked || optGuide {
			if dict.FollowString(key, &index) {
				completer.Start(index)
				for completer.Next() {
					fmt.Printf(" %s%s = %d;", key, completer.Key(), completer.Value())
				}
			}
		} else {
			for i := 0; i < len(key); i++ {
				if !dict.Follow(key[i], &index) {
					break
				}

				// Reads a value
				if dict.HasValue(index) {
					fmt.Printf(" %s = %d;", key[:i+1], dict.Value(index))
				}
			}
		}
		fmt.Println()
	}
}

func main() {
	flag.BoolVar(&optBuild, "b", false, "build dictionary")
	flag.BoolVar(&optTab, "t", false, "handle tab as separator")
	flag.BoolVar(&optGuide, "g", false, "build/load dictionary with guide")
	flag.BoolVar(&optRanked, "r", false, "build/load dictionary with ranked guide")
	flag.BoolVar(&optSort, "s", false, "sort lexicon before building dict")
	flag.BoolVar(&optUtfc, "u", false, "use utf-c instead of utf-8 for encoding keys")
	flag.StringVar(&optLexicon, "l", "-", "lexicon file")
	flag.StringVar(&optDictionary, "d", "-", "dictionary file")

	flag.Parse()

	if optLexicon == "-" && flag.NArg() > 0 {
		optLexicon = flag.Arg(0)
	}

	if optDictionary == "-" && flag.NArg() > 1 {
		optDictionary = flag.Arg(1)
	}

	var err error
	fileLexicon, err = os.Open(optLexicon)
	if err != nil {
		log.Fatal(err)
	}
	defer fileLexicon.Close()
	scanner = bufio.NewScanner(fileLexicon)

	if optBuild {
		fileDictionary, err = os.Create(optDictionary)
	} else {
		fileDictionary, err = os.Open(optDictionary)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer fileDictionary.Close()

	if optBuild {
		handleBuildDict()
	} else {
		handleLoadDict()
	}
}
