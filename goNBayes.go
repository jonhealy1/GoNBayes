package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

/*	words that won't be used for calculating probabilities
*/
var stopWords = map[string]bool{"a": true, "and": true, "for": true, "the": true,
	"is": true, "of": true, "in": true, "to": true}

/*	read a file and add it to a slice separated by new lines,
	for the data each element in a slice is a sentence representing
	a fortune cookie saying	
*/
func readFile(s string) []string {
	file, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	str := string(file)
	nstr := strings.Split(str, "\n")
	return nstr
}

func main() {

	trainSentences := readFile("traindata.txt")
	trainLabels := readFile("trainlabels.txt")
	testSentences := readFile("testdata.txt")
	testLabels := readFile("testlabels.txt")

	fmt.Println("train length: ", len(trainLabels))
	fmt.Println("train length: ", len(trainSentences))
	fmt.Println("test length: ", len(testLabels))
	fmt.Println("test length: ", len(testSentences))

	zero := "0"
	one := "1"
	zeroTrain := sepSent(trainSentences, trainLabels, zero)
	oneTrain := sepSent(trainSentences, trainLabels, one)

	pOne := float64(len(oneTrain)) / (float64(len(oneTrain)) + float64(len(zeroTrain)))

	var wordsZero []string
	var wordsOne []string

	wordsZero = wordList(zeroTrain)
	wordsOne = wordList(oneTrain)

	var wZero map[string]int
	wZero = make(map[string]int)

	var wOne map[string]int
	wOne = make(map[string]int)

	var wAll map[string]int
	wAll = make(map[string]int)

	for p := 0; p < len(wordsZero); p++ {
		wZero[wordsZero[p]]++
		wAll[wordsZero[p]]++
	}
	for p := 0; p < len(wordsOne); p++ {
		wOne[wordsOne[p]]++
		wAll[wordsOne[p]]++
	}

	totNumZeros := len(wZero)
	fmt.Println("#zeros", totNumZeros)
	totNumOnes := len(wOne)
	fmt.Println("#ones", totNumOnes)
	totNumWords := len(wAll)
	fmt.Println("#all", totNumWords)

	passer := getResults(trainSentences, trainLabels, wZero,
		wordsZero, wordsOne, wOne,
		pOne, totNumWords)

	passer2 := getResults(testSentences, testLabels, wZero,
		wordsZero, wordsOne, wOne,
		pOne, totNumWords)

	fmt.Println("Training Success Rate: ", passer, "%")
	fmt.Println("Testing Success Rate: ", passer2, "%")
}

func sepSent(trainSentences []string, trainLabels []string,
	value string) []string {
	var Train []string
	for i := 0; i < len(trainLabels); i++ {
		if trainLabels[i] == value {
			Train = append(Train, trainSentences[i])
		}
	}
	return Train
}

/*	Split a list of items on an empty space or a new line
*/
func Split(r rune) bool {
	return r == ' ' || r == '\n'
}

/*	Iterate through all of the data sentences calculate the probability,
	decide whether it sentence should be 0 or 1 then compare
	to the corresponding Labels to see if the predictions are correct.
	Return the percentages.
*/
func getResults(testSentences []string, testLabels []string, wZero map[string]int,
	wordsZero []string, wordsOne []string, wOne map[string]int,
	pOne float64, totNumWords int) float64 {
	var test string
	var truth int
	var ntruth int
	alpha := 0.01

	
	for k := 0; k < len(testSentences)-1; k++ {
		test = testSentences[k]
		test2 := strings.FieldsFunc(test, Split)
		pSentZero := math.Log2(1-pOne) + math.Log2(alpha)
		pSentOne := math.Log2(pOne) + math.Log2(alpha)
		for p := 0; p < len(test2); p++ {
			if !stopWords[test2[p]] {
				/*	add logs for every word in a sentence (test2[p])
					number of occurences of word in a class (0 or 1) + 1
					divided by number of words in class (0 or 1) + 
					number of unique words in both classes (0 and 1 */
				pSentZero += math.Log2((float64(wZero[test2[p]] + 1)) / (float64(len(wordsZero) + totNumWords)))
				//fmt.Println("Zero:", test2[p], "#: ", wZero[test2[p]], " : ", pSentZero)
				pSentOne += math.Log2((float64(wOne[test2[p]] + 1)) / (float64(len(wordsOne) + totNumWords)))
				//fmt.Println("One:", test2[p], "#: ", wOne[test2[p]], " : ", pSentOne)
			}
		}
		pSentOne = math.Pow(2, pSentOne)
		pSentZero = math.Pow(2, pSentZero)
		if pSentOne > pSentZero {
			fmt.Print("#", k, " ", test2, "Belongs to One - ")
			if testLabels[k] == "1" {
				fmt.Println("true")
				truth++
			} else {
				fmt.Println("false")
				ntruth++
			}
			//totalForOne++
		} else {
			fmt.Print("#", k, " ", test2, "Belongs to Zero - ")
			if testLabels[k] == "0" {
				fmt.Println("true")
				truth++
			} else {
				fmt.Println("false")
				ntruth++
			}
			//totalForZero++
		}
	}
	return (float64(truth) / (float64(ntruth) + float64(truth))) * 100
}

/*	Separate the slice of sentences into words belonging to class 1 or 0
	and append to a new sentence. Also parse according to the stop words
	which are common words that are not used to calculate the predictions. 
*/
func wordList(p []string) []string {
	var words []string
	var wordsOne []string
	for j := 0; j < len(p); j++ {
		words = strings.FieldsFunc(p[j], Split)
		for k := 0; k < len(words); k++ {
			if !stopWords[words[k]] {
				wordsOne = append(wordsOne, words[k])
			}
		}
	}
	return wordsOne
}
