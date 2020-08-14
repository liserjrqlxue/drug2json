package main

import (
	"flag"
	"log"
	"os"

	"github.com/liserjrqlxue/goUtil/jsonUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/version"
)

var (
	input = flag.String(
		"input",
		"",
		"input to be convert",
	)
	output = flag.String(
		"output",
		"",
		"output json file name, default is -input.json",
	)
)

func main() {
	version.LogVersion()
	flag.Parse()
	if *input == "" {
		flag.Usage()
		log.Print("-input is required!")
		os.Exit(1)
	}
	if *output == "" {
		*output = *input + ".json"
	}
	var drugInfo = DrugInfo{
		MedicineCate: DrugMedicineCate{},
		Desc: DrugDesc{
			ReportDesc: DrugReportDesc{},
			ReferenceDesc: DrugReferenceDesc{
				References: []DrugReferences{
					{},
				},
			},
			GenomicsDesc: DrugGenomicsDesc{
				Mutation: []DrugMutation{
					{Locus: []DrugLocus{
						{},
					}},
				},
			},
			MedicineDesc: DrugMedicineDesc{
				Name: DrugName{},
			},
		},
	}
	simpleUtil.CheckErr(jsonUtil.Json2File(*output, drugInfo))
	var drugInfoJson = string(simpleUtil.HandleError(jsonUtil.JsonIndent(drugInfo, "", "\t")).([]byte))
	log.Print(drugInfoJson)
	log.Print(jsonUtil.MarshalString(drugInfo))
}
