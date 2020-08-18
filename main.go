package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/liserjrqlxue/goUtil/jsonUtil"
	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/version"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// os
var (
	ex, _  = os.Executable()
	exPath = filepath.Dir(ex)
)

// flag
var (
	input = flag.String(
		"input",
		"",
		"input to be convert",
	)
	prefix = flag.String(
		"prefix",
		"",
		"output json file -prefix.sampleID.txt",
	)
	bgFile = flag.String(
		"background",
		filepath.Join(exPath, "background-child-V1.1-xcj20181206.xlsx"),
		"background database",
	)
	bgSheetName = flag.String(
		"backgroundSheet",
		"Sheet2",
		"background database sheet name",
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

	var backgroundDb = simpleUtil.Slice2MapMapArray(
		simpleUtil.HandleError(
			simpleUtil.HandleError(
				excelize.OpenFile(*bgFile),
			).(*excelize.File).
				GetRows(*bgSheetName),
		).([][]string),
		"药物中文名称",
	)

	var excel = simpleUtil.HandleError(excelize.OpenFile(*input)).(*excelize.File)
	// 读取样品信息
	var sampleInfo = simpleUtil.HandleError(excel.GetRows("样本信息")).([][]string)
	var sampleInfoMap = simpleUtil.Slice2MapMapArray(sampleInfo, "样品编号")
	// 读取药物结果
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
				Available:   1, // 默认设置为 1
				ProductCode: sampleInfoMap[sampleID]["产品编号"],
				IsPositive:  strconv.FormatBool(item["用药建议"] == "常规用药"),
				Gender:      sampleInfoMap[sampleID]["性别"],
				PhoneNum:    sampleInfoMap[sampleID]["电话"],
				Birthdate:   sampleInfoMap[sampleID]["出生日期"],
				ProductName: sampleInfoMap[sampleID]["产品名称"],
				SampleType:  sampleInfoMap[sampleID]["样品类型"],
				SampleNum:   sampleID,
				MedicineCate: DrugMedicineCate{
					Id:   "", // 默认设置为空
					Name: item["药物分类"],
				},
				Desc: DrugDesc{
					ReportDesc: DrugReportDesc{
						Guidance:       item["用药建议"],
						Interpretation: item["结果说明"],
					},
					ReferenceDesc: DrugReferenceDesc{
						Reactions:     "", // 默认设置为空
						RelateDisease: "", // 默认设置为空
						References:    str2DrugReferencesArray(backgroundDb[drugName]["中文版报告用到的文献"]),
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
						Brief: backgroundDb[drugName]["药物知识-中文"],
					},
				},
			}
		}

		var drugMutation, ok3 = drugInfo.Desc.GenomicsDesc.MutationMap[gene]
		if !ok3 {
			drugMutation = DrugMutation{
				Locus: []DrugLocus{},
				Gene:  gene,
				Rank:  1,  // 默认设置为 1
				Desc:  "", // 默认设置为空
			}
		}
		drugMutation.Locus = append(
			drugMutation.Locus,
			DrugLocus{
				SnpRs:       item["检测位点"],
				Advice:      item["分条用药建议"],
				Rs:          "", // 默认设置为空
				Metabolizer: "", // 默认设置为空
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
		var outputF = osUtil.Create(*prefix + "." + sampleID + ".txt")
		fmt.Printf("%s\n", sampleID)
		fmt.Println("----------------------------------------------------------------------------------------")
		for drugName, drugInfo := range drugs {
			simpleUtil.HandleError(
				fmt.Fprintf(
					outputF,
					"%s\t%s\t%s\n",
					time.Now().Format("20060102150405"),
					sampleID,
					jsonUtil.MarshalString(drugInfo),
				),
			)
			if drugInfo.Desc.ReportDesc.Guidance == "常规用药" {
				continue
			}
			fmt.Printf(
				"%s\t%s\t%s\t%s\t%s\t%s\n",
				drugInfo.MedicineCate.Name,
				drugName,
				drugInfo.Desc.GenomicsDesc.Mutation[0].Gene,
				drugInfo.Desc.GenomicsDesc.Mutation[0].Locus[0].SnpRs,
				drugInfo.Desc.GenomicsDesc.Mutation[0].Locus[0].GeneType,
				drugInfo.Desc.ReportDesc.Guidance,
			)
			for i, mutation := range drugInfo.Desc.GenomicsDesc.Mutation {
				for j, locus := range mutation.Locus {
					if i == 0 {
						if j == 0 {
							continue
						} else {
							fmt.Printf(
								"%-20s\t%s\t%s\t%s\t%s\t%s\n",
								"",
								"",
								"",
								locus.SnpRs,
								locus.GeneType,
								"",
							)
						}
					} else {
						if j == 0 {
							fmt.Printf(
								"%-20s\t%s\t%s\t%s\t%s\t%s\n",
								"",
								"",
								mutation.Gene,
								locus.SnpRs,
								locus.GeneType,
								"",
							)
						} else {
							fmt.Printf(
								"%-20s\t%s\t%s\t%s\t%s\t%s\n",
								"",
								"",
								"",
								locus.SnpRs,
								locus.GeneType,
								"",
							)
						}
					}
				}
			}
		}
		fmt.Println("----------------------------------------------------------------------------------------")
		simpleUtil.DeferClose(outputF)
	}
}

func str2DrugReferencesArray(str string) (references []DrugReferences) {
	var i = 1
	for _, ref := range strings.Split(str, "\n") {
		if ref != "" {
			var reference = DrugReferences{
				Id:    strconv.Itoa(i),
				Title: ref,
			}
			i++
			references = append(references, reference)
		}
	}
	return
}
