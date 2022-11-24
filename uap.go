package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const applicantsFilePath string = "applicants.txt"

var n int

// student struct consisting of firstName, lastName and grade
type applicant struct {
	firstName, lastName, firstChoice, secondChoice, thirdChoice            string
	physicsGrade, chemistryGrade, mathsGrade, compSciGrade, admissionScore float64
	enrolled                                                               bool
}

// course struct consisting of name and slice of enrolledStudents
type course struct {
	name             string
	enrolledStudents []applicant
}

// ApplicantList creating a type for the applicants slice
type ApplicantList []applicant

func main() {

	n = StudentsPerCourse()

	biotech := course{"Biotech", make([]applicant, 0)}
	chemistry := course{"Chemistry", make([]applicant, 0)}
	engineering := course{"Engineering", make([]applicant, 0)}
	mathematics := course{"Mathematics", make([]applicant, 0)}
	physics := course{"Physics", make([]applicant, 0)}

	// get applicants from file into a slice of struct:applicant
	applicants := ApplicantList(ApplicantListFromFile())

	// first choice
	applicants.AssignChoiceStudentsFor(&biotech, 1)
	applicants.AssignChoiceStudentsFor(&chemistry, 1)
	applicants.AssignChoiceStudentsFor(&engineering, 1)
	applicants.AssignChoiceStudentsFor(&mathematics, 1)
	applicants.AssignChoiceStudentsFor(&physics, 1)

	// second choice
	applicants.AssignChoiceStudentsFor(&biotech, 2)
	applicants.AssignChoiceStudentsFor(&chemistry, 2)
	applicants.AssignChoiceStudentsFor(&engineering, 2)
	applicants.AssignChoiceStudentsFor(&mathematics, 2)
	applicants.AssignChoiceStudentsFor(&physics, 2)

	// third choice
	applicants.AssignChoiceStudentsFor(&biotech, 3)
	applicants.AssignChoiceStudentsFor(&chemistry, 3)
	applicants.AssignChoiceStudentsFor(&engineering, 3)
	applicants.AssignChoiceStudentsFor(&mathematics, 3)
	applicants.AssignChoiceStudentsFor(&physics, 3)

	// slice of all course lists
	courseLists := []course{biotech, chemistry, engineering, mathematics, physics}

	//Sort and Print all Course Lists
	for _, c := range courseLists {
		c.SortAndPrintCourseList()
	}

}

func (s applicant) Choice(c int) string {
	switch c {
	case 1:
		return s.firstChoice
	case 2:
		return s.secondChoice
	case 3:
		return s.thirdChoice
	}
	return ""
}

func (a *ApplicantList) AssignChoiceStudentsFor(c *course, choice int) {
	b := 0
	a.sortBy(c.name)
	apps := *a
	for _, student := range apps {
		if student.Choice(choice) == c.name && SpaceOnCourse(*c) {
			AddToCourseList(student, c)
			student.enrolled = true
		}
		if student.enrolled == false {
			apps[b] = student
			b++
		}
	}

	*a = apps[:b]
}

func getSortingGrade(a applicant, course string) float64 {
	var courseGrade float64

	switch course {
	case "Physics":
		courseGrade = MeanAverage(a.physicsGrade, a.mathsGrade)
	case "Chemistry":
		courseGrade = a.chemistryGrade
	case "Mathematics":
		courseGrade = a.mathsGrade
	case "Engineering":
		courseGrade = MeanAverage(a.compSciGrade, a.mathsGrade)
	case "Biotech":
		courseGrade = MeanAverage(a.chemistryGrade, a.physicsGrade)
	}

	// using admission score if higher than courseGrade
	if courseGrade >= a.admissionScore {
		return courseGrade
	}

	return a.admissionScore
}

func (a ApplicantList) sortBy(course string) {
	sort.Slice(a, func(i, j int) bool {
		student1 := a[i]
		student2 := a[j]
		firstComparison := getSortingGrade(student1, course)
		secondComparison := getSortingGrade(student2, course)
		if firstComparison == secondComparison {
			return a[i].firstName < a[j].firstName
		}
		return firstComparison > secondComparison
	})
}

func SpaceOnCourse(universityCourse course) bool {
	return len(universityCourse.enrolledStudents) < n
}

func (c course) SortAndPrintCourseList() {
	// sorting student ist by course grade
	ApplicantList(c.enrolledStudents).sortBy(c.name)
	//fileName := fmt.Sprintf("%s%s.txt", baseFilePath, strings.ToLower(c.name))
	fileName := fmt.Sprintf("%s.txt", strings.ToLower(c.name))

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, line := range c.enrolledStudents {
		fmt.Fprintf(file, "%s %s %.2f\n", line.firstName, line.lastName, getSortingGrade(line, c.name)) // writes each line of the 'data' slice of strings
		if err != nil {
			log.Fatal(err)
		}
	}
}

func MeanAverage(a float64, b float64) float64 {
	return (a + b) / 2
}

func AddToCourseList(student applicant, universityCourse *course) {
	universityCourse.enrolledStudents = append(universityCourse.enrolledStudents, student)
}

func StudentsPerCourse() int {
	var n int
	_, err := fmt.Scan(&n)
	if err != nil {
		return 0
	}
	return n
}

func strToFloat(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal("Error converting")
	}
	return i
}

func ApplicantListFromFile() []applicant {
	// get file from path
	file, err := os.Open(applicantsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// read students line by line using NewScanner
	scanner := bufio.NewScanner(file)

	// create empty slice of applicants
	applicants := make([]applicant, 0)

	// iterate over each line of file and append to applicant to slice
	for scanner.Scan() {
		// split string input into space separated chunks
		var str = strings.Split(scanner.Text(), " ")

		// create instance of applicant
		applicant := applicant{
			str[0],
			str[1],
			str[7],
			str[8],
			str[9],
			strToFloat(str[2]),
			strToFloat(str[3]),
			strToFloat(str[4]),
			strToFloat(str[5]),
			strToFloat(str[6]),
			false,
		}
		applicants = append(applicants, applicant)
	}

	return applicants
}

