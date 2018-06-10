package main

import (
	// core
	cregex "github.com/mingrammer/commonregex"

	// helpers
	pp "github.com/sniperkit/colly/plugins/app/debug/pp"
)

func examples_cregex() {

	dateList := cregex.Date(RAW_TEXT)
	pp.Println("regex.Date=", dateList)
	// ['Jan 9th 2012']

	timeList := cregex.Time(RAW_TEXT)
	pp.Println("regex.Time=", timeList)
	// ['5:00PM', '4:00']

	linkList := cregex.Links(RAW_TEXT)
	pp.Println("regex.Links=", linkList)
	// ['www.linkedin.com', 'harold.smith@gmail.com']

	phoneList := cregex.PhonesWithExts(RAW_TEXT)
	pp.Println("regex.PhonesWithExts=", phoneList)
	// ['(519)-236-2723x341']

	emailList := cregex.Emails(RAW_TEXT)
	pp.Println("regex.Emails=", emailList)
	// ['harold.smith@gmail.com']
}
