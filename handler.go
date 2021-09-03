package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

type dataStruct struct {
	period string
}

type sample struct {
	data interface{}
}

func (s *sample) func1() {
	obj := map[string]interface{}{}
	// obj["foo1"] = "bar1"
	s.data = obj
}

func (s *sample) addField(key, value string) {
	v, _ := s.data.(map[string]interface{})
	v[key] = value
	s.data = v
}
func main() {

	s := &sample{}
	s.func1()
	//Son ödeme tarihi ve ödeme tarihi alınarak dönem parametresi alıp faiz hesabı yapılacak.
	//bu parametreler dışarıdan alınacak ve double bir değer dönecek.
	//amaç şu bir faturanın son ödeme tarihi var birde ödeme tarihi var.Mesela Faturanın son ödeme tarihi 30-04-2021
	//ödeme tarihi 16-07-2021 aradaki fark 77 gün.Burada 30 + 30 +17 gün devir olacak.
	var expiredDate = "30-07-2020"
	//son ödeme tarihine 1 ay ekle ve  devret

	var paymentDate = "15-09-2020"
	var donem = "06-2018"
	var interestRate = getData(donem)
	// i, _ := strconv.ParseFloat(interestRate, 64)
	log.Println(interestRate)
	// layout := "2006-01-02"
	layout := "02-01-2006"

	t, err := time.Parse(layout, expiredDate)

	var tt, mm = addOneMounth(t)
	//gelen son ödeme tarihinin başlangıç ve bitiş tarihini buldum
	log.Println(tt)
	log.Println(mm)
	t2, err := time.Parse(layout, paymentDate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	//nekadarlık bir gecikme oldu onu buldum
	var expDay = calculateDate(t, t2)
	log.Println(expDay)
	var deliverDate = calculateDate(tt, mm)
	log.Println(deliverDate)
	var resultCalculate = calculateIntrest(40.500, 15.500, 3000, interestRate, deliverDate)

	log.Println("---------ilk ay hesaplanan değer")
	log.Println(resultCalculate)
	log.Println("---------ilk ay hesaplanan değer")
	float64ToString := fmt.Sprint(resultCalculate)
	s.addField(tt.String(), float64ToString)
	var control = int(expDay - deliverDate)
	// var resultMountByMount = calculateOtherMount(resultCalculate, int(deliverDate), int(expDay))
	// log.Println(resultMountByMount)
	for control > -1 {
		// var firstTt = tt
		var r, tt, d = controlAndCalculateOtherMount(control, tt, resultCalculate)
		float64ToString := fmt.Sprint(d)
		s.addField(tt.String(), float64ToString)
		log.Println("------ikinci değer----")
		log.Println(d)
		log.Println("------ikinci değer----")
		control = r
		if control > -1 {

			// var secondTt = tt
			control, tt, resultCalculate = controlAndCalculateOtherMount(r, tt, d)
			float64ToString := fmt.Sprint(resultCalculate)
			s.addField(tt.String(), float64ToString)
			// control = control

		}

	}
	var sum float64 = 0
	for _, val := range s.data.(map[string]interface{}) {

		// log.Println(key)
		// log.Println(val)
		// log.Println(donem)
		var number = foo(val)
		sum += number
	}
	log.Println(sum)
	log.Println(s)

	// getData("asd")
}

func foo(veri interface{}) float64 {
	f64, _ := strconv.ParseFloat(veri.(string), 64)
	return f64 + 1
}

func calculateOtherMount(firstMountCalculate float64, dayOfMount int, delayMount int) float64 {

	var result = (firstMountCalculate / float64(dayOfMount)) * float64(delayMount)

	return result

}
func controlAndCalculateOtherMount(delayTimeMount int, firstMont time.Time, resultCalculate float64) (int, time.Time, float64) {

	var tt, mm = addOneMounth(firstMont)
	var deliverDate = calculateDate(tt, mm)

	var resultMountByMount = calculateOtherMount(resultCalculate, int(deliverDate), delayTimeMount)

	var control = int(float64(delayTimeMount) - deliverDate)

	return control, tt, resultMountByMount

}

func addOneMounth(date time.Time) (time.Time, time.Time) {

	var after = date.AddDate(0, 1, 0)
	// var app = Bod(after)
	year, month, _ := after.Date()
	firstDayOfThisMonth := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())

	endOfThisMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, date.Location())

	return firstDayOfThisMonth, endOfThisMonth

}
func Bod(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func getData(date string) string {

	var result string = ""
	resp, _ := soup.Get("https://www.tcmb.gov.tr/wps/wcm/connect/TR/TCMB+TR/Main+Menu/Temel+Faaliyetler/Para+Politikasi/Reeskont+ve+Avans+Faiz+Oranlari")
	doc := soup.HTMLParse(resp)
	rows := doc.Find("table").FindAll("tr")[1:]

	for _, row := range rows {
		columns := row.FindAll("td")

		tarih := columns[0].Text()
		// iskonto := columns[1].Text()
		faiz := columns[2].Text()
		res2 := strings.ReplaceAll(faiz, "\u00a0", "")
		res3 := strings.Replace(tarih, ".", "-", -1)
		layout := "02-01-2006"

		t, err := time.Parse(layout, res3)
		if err != nil {
			fmt.Println(err)
		}

		// year, month, _ := t.Date()
		// log.Println(t)
		// log.Println(year)
		// log.Println(int(month))

		formatted := fmt.Sprintf("%02d-%02d",
			t.Month(), t.Year())
		// log.Println(formatted)
		layout2 := "01-2006"
		t2, _ := time.Parse(layout2, date)
		t3, _ := time.Parse(layout2, formatted)
		if t2.Unix() > t3.Unix() {

			log.Println(faiz)

			result = res2

		}

		// fmt.Println(tarih + "\t" + iskonto + "\t" + faiz)
		log.Println(result)
	}
	return result
}
func calculateDate(expiredDate time.Time, paymentDate time.Time) float64 {

	var result = -1.0
	days := paymentDate.Sub(expiredDate).Hours() / 24
	if days >= 0 {

		result = (days + 1)

	}
	return result

}
func calculateIntrest(aylik_ucret float64, indirimler float64, devlete_iletilecek_vergiler float64, interestedRate string, date float64) float64 {

	interestRateRep := strings.Replace(interestedRate, ",", ".", -1)

	f, _ := strconv.ParseFloat(interestRateRep, 64)
	var dayOfInterest = ((f + 5) / 360) / 100
	var a = (((aylik_ucret - indirimler) * 29) * dayOfInterest) * 1000
	var b = devlete_iletilecek_vergiler * 18 * dayOfInterest

	var result = (a + b) * 1.28

	return result

}
