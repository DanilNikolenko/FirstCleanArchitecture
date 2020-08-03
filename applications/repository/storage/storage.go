package storage

import (
	"ProjectCleanArchitecture/FirstCleanArchitecture/models"
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	MIN = 97
	MAX = 122
)

type ApplicationsStorage struct {
	active []*models.Application
	cancel []*models.Application
}

func NewApplicationsRepository() *ApplicationsStorage {
	mu := new(sync.Mutex)

	// 1.First create ApplicationsStorage and then write into
	rez := &ApplicationsStorage{}

	// generate 50 active app
	for {
		temp := GetNewApplication()

		mu.Lock()
		rez.active = append(rez.active, temp)
		if len(rez.active) == viper.GetInt("maxApps") {
			break
		}
		mu.Unlock()
	}

	//create: go func(){ every 200msec delete 1 activeApp and add newApp}
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)

			// get random
			randINT := random(0, viper.GetInt("maxApps")-1)

			// add to cancel if it have >0 count shows
			mu.Lock()
			if rez.active[randINT].Count != 0 {
				rez.cancel = append(rez.cancel, rez.active[randINT])
			}

			// add new app instead old
			rez.active[randINT] = GetNewApplication()
			mu.Unlock()
		}
	}()

	return rez
}

func (r ApplicationsStorage) GetShowedAndCancelApplications(ctx context.Context) ([]models.Application, error) {
	mu := new(sync.Mutex)

	// create slice to return
	rez := []models.Application{}

	// add to slice apps with count > 0
	mu.Lock()
	for i, j := range r.active {
		if j.Count > 0 {
			rez = append(rez, *r.active[i])
		}
	}

	// add all canceled apps
	for i := range r.cancel {
		rez = append(rez, *r.active[i])
	}
	mu.Unlock()

	return rez, nil
}
func (r ApplicationsStorage) GetRandomAliveApplication(ctx context.Context) (models.Application, error) {
	mu := new(sync.Mutex)

	randINT := random(0, viper.GetInt("maxApps")-1)

	mu.Lock()
	r.active[randINT].Count++
	rez := *r.active[randINT]
	mu.Unlock()

	return rez, nil
}

// =====================

func GetNewApplication() *models.Application {
	AppName := GetApplicationRandomNAme()
	return &models.Application{
		Name:  AppName,
		Count: 0,
	}
}

func GetApplicationRandomNAme() string {
	var name string
	for i := 0; i < 2; i++ {
		TempRand := random(MIN, MAX)
		name += string(byte(TempRand))
	}

	return name
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
