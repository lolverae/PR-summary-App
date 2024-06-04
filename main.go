package main

import (
	svc "prsummaryapp/services"
)

func main() {
	prInfo := svc.GenerateReport()
	svc.SendReport(prInfo)
}
