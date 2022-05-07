package main

import (
	bst "GoSchool_Assignment2/binarysearchtree"
	dll "GoSchool_Assignment2/doublylinkedlist"
	util "GoSchool_Assignment2/utility"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

func printPatientsSearchResult(patients []interface{}, searchString string) {
	fmt.Printf("\nSearch Results for [%s]\n", searchString)
	fmt.Println("----------------------------------------------------")
	printList(patients)
}

func printPatientsList(patients []interface{}) {
	fmt.Println("\nDisplaying all patients...")
	fmt.Println("----------------------------------------------------")
	printList(patients)
}

func printList(patients []interface{}) {
	if len(patients) > 0 {
		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 0, 8, 5, '\t', 0)
		fmt.Fprintf(writer, "Selection\tName\tMobile Number\n")
		for idx, v := range patients {
			patient := v.(*Patient)
			fmt.Fprintf(writer, "[%v]\t%s\t%d\n", idx+1, patient.GetName(), patient.GetMobileNum())
		}
		writer.Flush()
		fmt.Println("----------------------------------------------------")
		fmt.Printf("(%d) Results\n", len(patients))
	} else {
		fmt.Println("0 Result")
	}
}

// Create new array to sort patient
func sortPatientsByMobileNumber(pl **dll.DoublyLinkedlist, channel chan error) {
	if (**pl).GetSize() == 0 {
		channel <- errors.New("\nerror: patient list is empty")
	} else {
		patientSlice := make([]*Patient, (**pl).GetSize())
		for k, v := range (**pl).GetList() {
			patient := v.(*Patient)
			patientSlice[k] = patient
		}
		sortedArr := insertionSort(patientSlice, len(patientSlice))
		// Clear linkedlist by setting head to nil, go will proceed with garbage collection
		(**pl).Clear()
		// Add Sorted patients back into linkedlist
		for _, v := range sortedArr {
			(**pl).Add(v)
		}
		channel <- nil
	}
}

func insertionSort(arr []*Patient, n int) []*Patient {
	for i := 1; i < n; i++ {
		data := arr[i]
		last := i
		for (last > 0) && (arr[last-1].GetMobileNum() > data.GetMobileNum()) {
			arr[last] = arr[last-1]
			last--
		}
		arr[last] = data
	}
	return arr
}

func addAppointment(appointmentTree **bst.Binarysearchtree, appointmentDate string, appointmentSession int, dentist *Dentist, patient *Patient, channel chan error) {
	isExist := (**appointmentTree).Contains(appointmentDate, appointmentSession, dentist, patient)
	if isExist {
		channel <- errors.New("error: this appointment already exist, please select another session")
	} else {
		(**appointmentTree).Add(appointmentDate, appointmentSession, dentist, patient)
		channel <- nil
	}
}

func printPatientSchedule(sessionList **dll.DoublyLinkedlist, patient *Patient, pSchedule []*bst.BinaryNode) {
	fmt.Println("\nListing appointments for:", patient.GetName())
	fmt.Println("----------------------------------------------------------------------------------------------------------")

	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 5, '\t', 0)
	fmt.Fprintf(writer, "S/N\tDentist\tDate\tDay\tSession\tStart Time\tEnd Time\n")
	for idx, data := range pSchedule {
		dentist := data.GetDentist().(*Dentist)
		date, _ := time.Parse("2006-01-02", data.GetDate())
		r, _ := (**sessionList).Get(data.GetSession())
		session := r.(Session)
		fmt.Fprintf(writer, "%v\t%s\t%v\t%s\t%d\t%s\t%s\n", idx+1, dentist.GetName(), date.Format("2006-01-02"), date.Weekday(), data.GetSession(), session.GetStartTime(), session.GetEndTime())
	}
	writer.Flush()
}

func searchPatient(sessionList, patientList **dll.DoublyLinkedlist, appointmentTree **bst.Binarysearchtree) {
	var patient *Patient
	var patients []interface{}
	fmt.Println("\nSearch Patient")
	fmt.Println("--------------")
	fmt.Println("\nPlease enter patient name. Press enter to view all patients:")
	searchString := util.ReadInput()
	if len(searchString) == 0 {
		patients = (**patientList).GetList()
		printPatientsList(patients)
	} else {
		patients = (**patientList).SearchByName(searchString)
		printPatientsSearchResult(patients, searchString)
	}
patientInput:
	if len(patients) > 0 {
		fmt.Println("\nPlease select a patient:")
		patientSelection, err := util.ReadInputAsInt()
		if err == nil {
			if patientSelection > len(patients) || patientSelection <= 0 {
				fmt.Println(invalidSelection)
				goto patientInput
			}
			patient = patients[patientSelection-1].(*Patient)
			fmt.Printf("Selected [%s]\n", patient.GetName())
		} else {
			fmt.Println(err)
			fmt.Println(invalidSelection)
			goto patientInput
		}
		fmt.Println("\nBrowse appointments")
		fmt.Println("-------------------")
		fmt.Println("[1] View all appointments")
		fmt.Println("[2] View upcoming appointments")
	browseInput:
		fmt.Println("\nPlease enter which type to display:")
		browseSelection, err := util.ReadInputAsInt()
		if err == nil {
			if browseSelection > 2 || browseSelection <= 0 {
				fmt.Println(invalidSelection)
				goto browseInput
			}
		} else {
			fmt.Println(err)
			goto browseInput
		}
		switch browseSelection {
		case 1:
			pSchedule := (**appointmentTree).GetAllSchedule(patient)
			if len(pSchedule) > 0 {
				printPatientSchedule(sessionList, patient, pSchedule)
			} else {
				fmt.Printf("\nNo appointments made for [%s]\n", patient.GetName())
			}
		case 2:
			pSchedule := (**appointmentTree).GetUpComingSchedule(patient)
			if len(pSchedule) > 0 {
				printPatientSchedule(sessionList, patient, pSchedule)
			} else {
				fmt.Printf("\nNo upcoming appointments made for [%s]\n", patient.GetName())
			}
		}
	}
}
