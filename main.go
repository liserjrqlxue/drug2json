package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/liserjrqlxue/goUtil/jsonUtil"
	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/version"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
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
		"output json file name, default is -input.txt",
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
		*output = *input + ".txt"
	}
	var outputF = osUtil.Create(*output)
	defer simpleUtil.DeferClose(outputF)

	// 读取药物结果
	var excel = simpleUtil.HandleError(excelize.OpenFile(*input)).(*excelize.File)
	var results = simpleUtil.HandleError(excel.GetRows("检测结果")).([][]string)
	var db = simpleUtil.Slice2MapArray(results)
	var info = make(map[string]map[string]*DrugInfo)
	for _, item := range db {
		var sampleID = item["样本编号"]
		var drugName = item["药物名称"]
		var gene = item["检测基因"]
		var drugs, ok1 = info[sampleID]
		if !ok1 {
			drugs = map[string]*DrugInfo{}
		}
		var drugInfo, ok2 = drugs[drugName]
		if !ok2 { // 药物初值
			drugInfo = &DrugInfo{
				MedicineCate: DrugMedicineCate{
					Id:   item["不知道是什么"],
					Name: item["药物分类"],
				},
				Desc: DrugDesc{
					ReportDesc: DrugReportDesc{
						Guidance:       item["用药建议"],
						Interpretation: item["结果说明"],
					},
					ReferenceDesc: DrugReferenceDesc{
						Reactions:     "不知道是什么",
						RelateDisease: "不知道是什么",
						References: []DrugReferences{
							{
								Id:    "0",
								Title: "[0]数据库没有参考文献",
							},
						},
					},
					GenomicsDesc: DrugGenomicsDesc{
						MutationMap: map[string]DrugMutation{},
						Mutation:    []DrugMutation{},
					},
					MedicineDesc: DrugMedicineDesc{
						Name: DrugName{
							Cn: item["药物名称"],
							En: item["英文名"],
						},
						Brief: "药物背景数据库待提供",
					},
				},
			}
		}

		var drugMutation, ok3 = drugInfo.Desc.GenomicsDesc.MutationMap[gene]
		if !ok3 {
			drugMutation = DrugMutation{
				Locus: []DrugLocus{},
				Gene:  gene,
				Rank:  0, // 不知道是什么，默认设置为0
				Desc:  "不知道是什么",
			}
		}
		drugMutation.Locus = append(
			drugMutation.Locus,
			DrugLocus{
				SnpRs:       item["检测位点"],
				Advice:      item["分条用药建议"],
				Rs:          "不知道为什么是空的",
				Metabolizer: "不知道是什么",
				GeneType:    item["检测结果"],
			},
		)
		// ok3
		drugInfo.Desc.GenomicsDesc.MutationMap[gene] = drugMutation

		// MutationMap -> Mutation
		var mutations []DrugMutation
		for _, mutation := range drugInfo.Desc.GenomicsDesc.MutationMap {
			mutations = append(mutations, mutation)
		}
		drugInfo.Desc.GenomicsDesc.Mutation = mutations

		// ok2
		drugs[drugName] = drugInfo
		// ok1
		info[sampleID] = drugs
	}
	for sampleID, drugs := range info {
		for _, drugInfo := range drugs {
			simpleUtil.HandleError(
				fmt.Fprintf(
					outputF,
					"%s\t%s\t%s\n",
					time.Now().Format("20060102150405"),
					sampleID,
					jsonUtil.MarshalString(drugInfo),
				),
			)
		}
	}
}
