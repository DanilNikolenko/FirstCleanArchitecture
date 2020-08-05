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
	mu     *sync.Mutex
	active []*models.Application
	cancel []*models.Application
}

func NewApplicationsRepository() *ApplicationsStorage {
	MaxCount := viper.GetInt("maxApps")
	// 1.First create ApplicationsStorage and then write into
	rez := &ApplicationsStorage{
		mu:     new(sync.Mutex),
		active: make([]*models.Application, 0),
		cancel: make([]*models.Application, 0),
	}

	// generate 50 active app
	for i := 0; i < MaxCount; i++ {
		tempApp := models.Application{
			Name:  GetApplicationRandomNAme(),
			Count: 0,
		}

		rez.mu.Lock()
		rez.active = append(rez.active, &tempApp)
		rez.mu.Unlock()
	}

	//create: go func(){ every 200msec delete 1 activeApp and add newApp}
	c := time.Tick(200 * time.Millisecond)
	go func() {
		for range c {
			rez.refreshAvailableAppPool()
		}
	}()

	return rez
}

func (r *ApplicationsStorage) refreshAvailableAppPool() {
	// get random
	randINT := random(0, viper.GetInt("maxApps")-1)

	// add to cancel if it have >0 count shows
	r.mu.Lock()

	if r.active[randINT].Count != 0 {
		temp := &models.Application{
			Name:  r.active[randINT].Name,
			Count: r.active[randINT].Count,
		}
		r.cancel = append(r.cancel, temp)
	}

	// add new app instead old
	tempApp := models.Application{
		Name:  GetApplicationRandomNAme(),
		Count: 0,
	}

	r.active[randINT] = &tempApp

	r.mu.Unlock()
}

func (r *ApplicationsStorage) GetShowedAndCancelApplications(ctx context.Context) (rezActive []models.Application, rezCancel []models.Application, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// add to slice apps with count > 0
	for i, j := range r.active {
		if j.Count > 0 {
			rezActive = append(rezActive, *r.active[i])
		}
	}

	// add all canceled apps
	for i := range r.cancel {
		rezCancel = append(rezCancel, *r.cancel[i])
	}

	return rezActive, rezCancel, nil
}

func (r *ApplicationsStorage) GetRandomAliveApplication(ctx context.Context) (models.Application, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	randINT := random(0, viper.GetInt("maxApps")-1)

	r.active[randINT].Count++
	rez := *r.active[randINT]

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
