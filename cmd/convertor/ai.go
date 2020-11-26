package convertor

import (
	"encoding/json"
	fmt "fmt"
	"github.com/aws/aws-sdk-go/service/textract"
	"io/ioutil"
	"os"
)

func Run() interface{} {
	jsonFile, err := os.Open("../../assets/results.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []textract.Block
	json.Unmarshal([]byte(byteValue), &result)
	blocksMap := make(map[string]textract.Block)
	var tableBlocks []textract.Block
	for _, s := range result {
		blocksMap[*s.Id] = s
		if *s.BlockType == "TABLE" {
				tableBlocks = append(tableBlocks, s)
		}
	}

	fmt.Println(len(tableBlocks))


	row := make(map[int64]interface{})
	rowIndex := make(map[int]interface{})
	var sentence string
	var main []string
	for i,s := range tableBlocks {
		_,to := rowIndex[i]
		if !to {
			rowIndex[i] = ""
			row = make(map[int64]interface{})
		}
		for _,relationship := range s.Relationships {
			if *relationship.Type == "CHILD" {
				for _, childId := range relationship.Ids {
					cell := blocksMap[*childId]
					if *cell.BlockType == "CELL" {
						//fmt.Println(rowIndex, *cell.RowIndex)


							_, ok := row[*cell.RowIndex]
							if !ok {

								sentence = ""
								main = nil
							}

							sentence += getText(cell, blocksMap)

							row[*cell.RowIndex] = append(main, sentence)
						}


					}

				}

			}

		rowIndex[i] = row

	}

	return rowIndex
}



func getText(cell textract.Block, blocksMap map[string]textract.Block)string{
	text := ""
	if len(cell.Relationships) > 0 {
		for _, relationship := range cell.Relationships{
			if *relationship.Type == "CHILD" {
				for _, childId := range relationship.Ids {
					word := blocksMap[*childId]
					if *word.BlockType == "WORD" {
						text += *word.Text + " "
					}
					if *word.BlockType == "SELECTION_ELEMENT" {
						if *word.SelectionStatus == "SELECTED"{
							text +=  "X "
						}
					}
				}
			}
		}
	}

	return text + ","
}
