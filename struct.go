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
	Interpretation string `json:"interpretation"`
}

type DrugReferenceDesc struct {
	Reactions     string           `json:"reactions"`
	RelateDisease string           `json:"relateDisease"`
	References    []DrugReferences `json:"references"`
}

type DrugReferences struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type DrugGenomicsDesc struct {
	Mutation []DrugMutation `json:"mutation"`
}

type DrugMutation struct {
	Locus []DrugLocus `json:"locus"`
	Gene  string      `json:"gene"`
	Rank  int         `json:"rank"`
	Desc  string      `json:"desc"`
}

type DrugLocus struct {
	SnpRs       string `json:"snpRs"`
	Advice      string `json:"advice"`
	Rs          string `json:"rs"`
	Metabolizer string `json:"metabolizer"`
	GeneType    string `json:"genetype"`
}

type DrugMedicineDesc struct {
	Name  DrugName `json:"name"`
	Brief string   `json:"brief"`
}

type DrugName struct {
	En string `json:"en"`
	Cn string `json:"cn"`
}
