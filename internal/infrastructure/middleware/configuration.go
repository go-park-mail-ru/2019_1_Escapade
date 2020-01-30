package middleware

type CorsConfiguration struct {
	Origins     []string
	Headers     []string
	Methods     []string
	Credentials string
}

//easyjson:json
type CorsConfigurationJSON struct {
	Origins     []string `json:"origins"`
	Headers     []string `json:"headers"`
	Methods     []string `json:"methods"`
	Credentials string   `json:"credentials"`
}

func (c CorsConfigurationJSON) Get() CorsConfiguration {
	return CorsConfiguration(c)
}
