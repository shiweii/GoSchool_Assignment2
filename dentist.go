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

func printDentistList(dentists []interface{}) {
	fmt.Println("\nDentist")
	fmt.Println("-------")
	for idx, v := range dentists {
		dentist := v.(*Dentist)
		fmt.Printf("[%v] %v\n", idx+1, dentist.GetName())
	}
}

func getDentistFromSelection(dentists []interface{}) (*Dentist, bool) {
	var dentist *Dentist
	fmt.Println("\nPlease select a dentist:")
	dentistSelection, err := util.ReadInputAsInt()
	if err == nil {
		if dentistSelection > len(dentists) || dentistSelection <= 0 {
			fmt.Println(invalidSelection)
			return nil, false
		}
		dentist = dentists[dentistSelection-1].(*Dentist)
		fmt.Printf("Selected dentist [%s]\n", dentist.GetName())
	} else {
		fmt.Println(err)
		return nil, false
	}
	return dentist, true
}

func printDentistSchedule(sessionList **dll.DoublyLinkedlist, inputDate string, dentist *Dentist, schedule *[]*bst.BinaryNode) {
	fmt.Println("\nListing appointments for Dentist:", dentist.GetName())
	fmt.Println("-----------------------------------------------------------------------------------------------------")
	writer := new(tabwriter.Writer)
	writer.Init(os.Stdout, 0, 8, 5, '\t', 0)
	fmt.Fprintln(writer, "Session\tDate\tDay\tStart Time\tEnd Time\tStatus")
	for _, v := range (**sessionList).GetList() {
		session := v.(Session)
		status := "Available"
		date, _ := time.Parse("2006-01-02", inputDate)
		for _, data := range *schedule {
			if data.GetSession() == session.GetSession() {
				status = "Booked"
			}
		}
		fmt.Fprintf(writer, "%d\t%+v\t%s\t%s\t%s\t%s\n", session.GetSession(), inputDate, date.Weekday(), session.GetStartTime(), session.GetEndTime(), status)
	}
	writer.Flush()
}

func printDentistAdminSchedule(sessionList **dll.DoublyLinkedlist, dentist *Dentist, dSchedule []*bst.BinaryNode) {

	if len(dSchedule) > 0 {
		fmt.Println("\nListing appointments for:", dentist.GetName())
		fmt.Println("--------------------------------------------------------------------------------------------------------------")

		writer := new(tabwriter.Writer)
		writer.Init(os.Stdout, 0, 8, 5, '\t', 0)
		fmt.Fprintf(writer, "S/N\tPatient\tMobile Number\tDate\tDay\tSession\tStart Time\tEnd Time\n")
		for idx, data := range dSchedule {
			patient := data.GetPatient().(*Patient)
			date, _ := time.Parse("2006-01-02", data.GetDate())
			r, _ := (**sessionList).Get(data.GetSession())
			session := r.(Session)
			fmt.Fprintf(writer, "%v\t%s\t%d\t%v\t%s\t%d\t%s\t%s\n", idx+1, patient.GetName(), patient.GetMobileNum(), date.Format("2006-01-02"), date.Weekday(), data.GetSession(), session.GetStartTime(), session.GetEndTime())
		}
		writer.Flush()
	} else {
		fmt.Println(noAppointment)
	}
}

func viewDentistSchedule(sessionList, dentistList **dll.DoublyLinkedlist, appointmentTree **bst.Binarysearchtree) error {

	var (
		ok      bool
		dentist *Dentist
	)

	if (**dentistList).GetSize() == 0 {
		return errors.New(noDentist)
	}

	fmt.Println("\nSearch Doctor Schedule")
	fmt.Println("---------------------")

	dentists := (**dentistList).GetList()
	printDentistList(dentists)

dentistInput:
	dentist, ok = getDentistFromSelection(dentists)
	if !ok {
		goto dentistInput
	}

apptInput:
	fmt.Println("\nPlease enter appointment date (yyyy-mm-dd):")
	inputDate := util.ReadInput()
	appointmentDate, err := time.Parse("2006-01-02", inputDate)
	if err != nil {
		fmt.Println(invalidDate)
		goto apptInput
	} else {
		currentTime := time.Now()
		if !currentTime.Before(appointmentDate) {
			fmt.Println(olderDate)
			goto apptInput
		}
	}

	schedule := (**appointmentTree).GetSchedule(appointmentDate.Format("2006-01-02"), dentist)
	printDentistSchedule(sessionList, appointmentDate.Format("2006-01-02"), dentist, &schedule)

	return nil
}

func browseDoctorAppointment(sessionList, dentistList **dll.DoublyLinkedlist, appointmentTree **bst.Binarysearchtree) error {

	if (**dentistList).GetSize() == 0 {
		return errors.New(noDentist)
	}

	var (
		ok      bool
		dentist *Dentist
	)

	dentists := (**dentistList).GetList()
	printDentistList(dentists)

dentistInput:
	dentist, ok = getDentistFromSelection(dentists)
	if !ok {
		goto dentistInput
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
		dSchedule := (**appointmentTree).GetAllSchedule(dentist)
		printDentistAdminSchedule(sessionList, dentist, dSchedule)
	case 2:
		dSchedule := (**appointmentTree).GetUpComingSchedule(dentist)
		printDentistAdminSchedule(sessionList, dentist, dSchedule)
	}
	return nil
}
