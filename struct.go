package main

type DrugInfo struct {
	Available    int              `json:"available"`
	ProductCode  string           `json:"productCode"`
	IsPositive   string           `json:"isPositive"`
	Gender       string           `json:"gender"`
	PhoneNum     string           `json:"phoneNum"`
	Birthdate    string           `json:"birthdate"`
	ProductName  string           `json:"productName"`
	SampleType   string           `json:"sampleType"`
	SampleNum    string           `json:"sampleNum"`
	MedicineCate DrugMedicineCate `json:"medicineCate"`
	Desc         DrugDesc         `json:"desc"`
}

type DrugMedicineCate struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type DrugDesc struct {
	ReportDesc    DrugReportDesc    `json:"reportDesc"`
	ReferenceDesc DrugReferenceDesc `json:"referenceDesc"`
	GenomicsDesc  DrugGenomicsDesc  `json:"genomicsDesc"`
	MedicineDesc  DrugMedicineDesc  `json:"medicineDesc"`
}

type DrugReportDesc struct {
	Guidance       string `json:"guidance"`
	Interpretation string `json:"interpretation"` // 解读
}

type DrugReferenceDesc struct {
	Reactions     string           `json:"reactions"`
	RelateDisease string           `json:"relateDisease"`
	References    []DrugReferences `json:"references"`
}

type DrugReferences struct {
	Id    string `json:"id"`
	Title string `json:"title"` // 参考文献
}

type DrugGenomicsDesc struct {
	Mutation    []DrugMutation          `json:"mutation"`
	MutationMap map[string]DrugMutation `json:"-"`
}

type DrugMutation struct {
	Locus []DrugLocus `json:"locus"`
	Gene  string      `json:"gene"` // 基因
	Rank  int         `json:"rank"`
	Desc  string      `json:"desc"`
}

type DrugLocus struct {
	SnpRs       string `json:"snpRs"`  // 检测位点
	Advice      string `json:"advice"` // 建议
	Rs          string `json:"rs"`
	Metabolizer string `json:"metabolizer"`
	GeneType    string `json:"genetype"` // 基因型
}

type DrugMedicineDesc struct {
	Name  DrugName `json:"name"`
	Brief string   `json:"brief"` // 背景
}

type DrugName struct {
	En string `json:"en"`
	Cn string `json:"cn"` // 药物名称
}
