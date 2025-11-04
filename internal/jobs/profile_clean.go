package main

import (
	"context"
	"log"
	"strings"

	"github.com/BeatEcoprove/identityService/internal/repositories"
	"github.com/BeatEcoprove/identityService/pkg/adapters"
)

type ProfileCleanUpJob struct {
	ctx         context.Context
	pubSub      adapters.RedisConsumer
	profileRepo repositories.IProfileRepository
}

func NewProfileCleanUpJob(
	pubSub adapters.RedisConsumer,
	profileRepo repositories.IProfileRepository,
) *ProfileCleanUpJob {
	return &ProfileCleanUpJob{
		ctx:         context.Background(),
		pubSub:      pubSub,
		profileRepo: profileRepo,
	}
}

func (pc *ProfileCleanUpJob) StartConsuming() {
	const redisChannel = "__keyevent@0__:expired"

	if err := pc.pubSub.EnableOpt(pc.ctx, "notify-keyspace-events", "Ex"); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	channel := pc.pubSub.Subscribe(pc.ctx, redisChannel)
	defer channel.Close()

	log.Printf("Starting Redis Expire listener")

	_, err := channel.Receive(pc.ctx)

	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	ch := channel.Channel()

	for msg := range ch {
		pc.handleExpireKey(msg.Payload)
	}
}

func (pc *ProfileCleanUpJob) handleExpireKey(key string) error {
	log.Printf("Key expired: %s", key)

	if !strings.HasPrefix(key, "prefix") {
		log.Println("No valid value")
		return nil
	}

	profileID := strings.TrimPrefix(key, "profile:pending:")

	profile, err := pc.profileRepo.Get(profileID)

	if err != nil {
		log.Println("Profile not found")
		return err
	}

	if err := pc.profileRepo.Delete(profile); err != nil {
		log.Println("Failed to remove profile")
		return err
	}

	return nil
}
