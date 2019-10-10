package server

import (
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

//TODO
// realize docker, grpc, tcp, http, script check - https://www.consul.io/docs/agent/checks.html

// ConsulClient register service and start healthchecking
func ConsulClient(serviceName, host, portString string, portInt int, tags []string,
	consulPort string, ttl time.Duration, check func() (bool, error),
	finish chan interface{}) (*consulapi.Client, string, error) {
	var (
		config = &consulapi.Config{
			Address:   host + consulPort,
			Scheme:    "http",
			Transport: cleanhttp.DefaultPooledTransport(),
		}
		consul, err = consulapi.NewClient(config)
		serviceID   = serviceName + portString
	)
	if err != nil {
		return consul, serviceID, err
	}
	host = FixHost(host)
	if strings.Contains(host, "http://") {
		host = strings.Replace(host, "http://", "", 1)
	}
	if strings.Contains(host, "https://") {
		host = strings.Replace(host, "https://", "", 1)
	}

	fmt.Println("tll:", ttl.String())

	agent := consul.Agent()
	consul.Agent().ServiceDeregister(serviceID)
	err = agent.ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Port:    portInt,
		Address: host,
		Tags:    tags,
		// https://www.consul.io/docs/agent/checks.html
		Check: &consulapi.AgentServiceCheck{
			TTL:                            ttl.String(),
			DeregisterCriticalServiceAfter: time.Minute.String(),
		},
	})
	if err != nil {
		utils.Debug(false, "cant add service to consul", err)
		return consul, serviceID, err
	}

	go updateTTL(agent, serviceID, ttl, check, finish)
	return consul, serviceID, nil
}

func updateTTL(agent *consulapi.Agent, serviceID string, ttl time.Duration,
	check func() (bool, error), finish chan interface{}) {
	if ttl.Seconds() > 5 {
		ttl = ttl - 5*time.Second
	}
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			update(agent, serviceID, check)
		case <-finish:
			close(finish)
			return
		}
	}
}

func update(agent *consulapi.Agent, serviceID string, check func() (bool, error)) {

	var (
		isWarning, err = check()
		status         = consulapi.HealthPassing
		message        = "Alive and reachable"
	)
	if err != nil {
		message = err.Error()
		if isWarning {
			status = consulapi.HealthWarning
			utils.Debug(false, "healthcheck function warning:", message)
		} else {
			status = consulapi.HealthCritical
			utils.Debug(false, "healthcheck function error:", message)
		}
	}
	err = agent.UpdateTTL("service:"+serviceID, message, status)
	if err != nil {
		utils.Debug(false, "agent of", serviceID, " UpdateTTL error:", err.Error())
	}
}
