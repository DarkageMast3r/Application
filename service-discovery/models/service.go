package models

type Service struct {
	Name       string
	Hosts      []string
	LastServed int
}

var services map[string]Service

func Service_Init() {
	services = make(map[string]Service)
}

func Service_Create(name string, uri string) {
	service, exists := services[name]
	if !exists {
		service.Name = name
		service.Hosts = make([]string, 0)
		service.LastServed = 0
	}
	for _, current := range service.Hosts {
		if current == uri {
			return
		}
	}
	service.Hosts = append(service.Hosts, uri)
	service.Save()
}

func Service_Get_Names() []string {
	names := make([]string, len(services))
	i := 0
	for key := range services {
		names[i] = key
		i++
	}
	return names
}

func Service_Get_By_Name(name string) (Service, bool) {
	service, exists := services[name]
	return service, exists
}

func (service Service) Save() {
	services[service.Name] = service
}
