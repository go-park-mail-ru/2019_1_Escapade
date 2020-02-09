package traefik

const (
	ErrNoConfiguration = "Configuration not given"

	LabelEnable   = "traefik.enable=true"
	LabelServices = "traefik.http.services."
	LabelRouters  = "traefik.http.routers."
	LabelNetwork  = "traefik.docker.network="
)
