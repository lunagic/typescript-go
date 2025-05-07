package typescript

import "reflect"

type ServiceConfigFunc func(s *Service)

func WithCustomNamespace(namespace string) ServiceConfigFunc {
	return func(s *Service) {
		s.outputNamespace = namespace
	}
}

func WithTypes(registry map[string]reflect.Type) ServiceConfigFunc {
	return func(s *Service) {
		s.outputTypes = registry
	}
}

func WithRoutes(routes map[string]Route) ServiceConfigFunc {
	return func(s *Service) {
		s.outputRoutes = routes
	}
}

func WithData(data map[string]any) ServiceConfigFunc {
	return func(s *Service) {
		s.outputData = data
	}
}
