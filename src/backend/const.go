package backend

type student struct {
	Rank       int
	Name       string
	Percentage string
}

type Result struct {
	Absent      string
	Comptt      string
	Other       string
	Passed      string
	Repeat      string
	Total       string
	PI          float64
	Topper      student
	LowestMarks student
}

var basicSubjectCodes = map[string]string{
	// Class 10
	"002": "HINDI",
	"086": "SCIENCE",
	"087": "SOCIAL SC.",
	"122": "SANSKRIT",
	"184": "ENGLISH",
	"241": "MATHEMATICS - BASIC",
	"401": "RETAILING",
	"402": "IT",

	// Common
	"041": "MATHEMATICS",

	// Class 12
	"027": "HISTORY",
	"028": "POLITICAL SCIENCE",
	"029": "GEOGRAPHY",
	"030": "ECONOMICS",
	"042": "PHYSICS",
	"043": "CHEMISTRY",
	"044": "BIOLOGY",
	"048": "PHYSICAL EDUCATION",
	"049": "PAINTING",
	"054": "BUSINESS STUDIES",
	"055": "ACCOUNTANCY",
	"083": "COMPUTER SCIENCE",
	"301": "ENGLISH CORE",
	"302": "HINDI CORE",
	"806": "TOURISM",
	"809": "FOOD PRODUCTION",
}

var outputTemplate = []string{
	"ROLL", "GENDER", "NAME",
	"SUBJECT CODE 1", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
	"SUBJECT CODE 2", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
	"SUBJECT CODE 3", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
	"SUBJECT CODE 4", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
	"SUBJECT CODE 5", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
	"SUBJECT CODE 6", "SUBJECT NAME", "MARKS OBTAINED", "GRADE",
	"TOTAL", "PERCENTAGE", "RESULT",
}
