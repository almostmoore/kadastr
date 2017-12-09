package parser

import (
	"github.com/almostmoore/kadastr/feature"
	"github.com/almostmoore/kadastr/rapi"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"sync"
	"time"
)

type FeatureParser struct {
	session *mgo.Session
	fRepo   feature.FeatureRepository
	rClient *rapi.Client
}

func NewFeatureParser(session *mgo.Session) FeatureParser {
	return FeatureParser{
		session: session,
		fRepo:   feature.NewFeatureRepository(session),
		rClient: rapi.NewClient(),
	}
}

// Run function starts parsing
func (f *FeatureParser) Run(quarter string, streams int64) {
	var maxUnit int64 = 10000

	done := make(chan bool, streams)
	errors := make(chan bool, streams)
	items := make(chan int64, maxUnit)
	defer close(done)
	defer close(errors)
	defer close(items)

	wg := &sync.WaitGroup{}

	var i int64
	for i = 0; i < streams; i++ {
		wg.Add(1)
		go f.parse(quarter, items, errors, done, wg)
	}

	go f.checkError(errors, done, streams)
	go func() {
		for i = 0; i < maxUnit; i++ {
			items <- i
		}
	}()

	wg.Wait()
}

func (f *FeatureParser) checkError(errors chan bool, done chan bool, streams int64) {
	errCount := 0

	for has := range errors {
		if has {
			errCount += 1
		} else {
			errCount = 0
		}

		if errCount == 200 {
			var i int64
			for i = 0; i < streams; i++ {
				done <- true
			}
		}
	}
}

// parse data from rosreestr
func (f *FeatureParser) parse(quarter string, items <-chan int64, errors, done chan bool, wg *sync.WaitGroup) {
	for {
		select {
		case i := <-items:
			result := f.parseItem(quarter, i)
			errors <- !result
		case <-done:
			wg.Done()
			return
		default:
		}
	}
}

// parseItem Parse item for quarter
func (f *FeatureParser) parseItem(quarter string, item int64) bool {
	time.Sleep(5 * time.Second)
	number := quarter + ":" + strconv.FormatInt(item, 10)
	log.Printf("Парсинг участка %s\n", number)

	ft, err := f.rClient.GetFeature(number)
	if err != nil || ft.CadNumber == "" {
		log.Printf("Участок не найден %s (%s)\n", number, err)
		return false
	}

	_, err = f.fRepo.FindByCadNumber(ft.CadNumber)
	if err == nil {
		log.Printf("Участок %s уже присутствует в базе данных. Пропускаем\n", ft.CadNumber)
		return true
	}

	ft.ID = bson.NewObjectId()
	err = f.fRepo.Insert(ft)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Участок сохранен %s\n", number)
	}

	return true
}
