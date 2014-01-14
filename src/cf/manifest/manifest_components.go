package manifest

import (
	"cf"
	"errors"
	"fmt"
	"generic"
	"strconv"
)

type manifestComponents struct {
	Applications   cf.AppSet
	GlobalServices []string
	GlobalEnvVars  generic.Map
}

func newManifestComponents(data generic.Map) (m manifestComponents, errs ManifestErrors) {
	m.Applications = cf.NewEmptyAppSet()
	m.GlobalEnvVars = generic.NewMap()
	m.GlobalServices = []string{}

	if data.Has("applications") {
		m.Applications = cf.NewAppSet(data.Get("applications"))
		for _, app := range m.Applications {
			appErrs := validateAppParams(app)
			if !appErrs.Empty() {
				errs = append(errs, appErrs...)
			}

			for _, fieldName := range []string{"instances"} {
				if app.Has(fieldName) && app.Get(fieldName) != nil {
					value, err := strconv.Atoi(app.Get(fieldName).(string))
					if err != nil {
						errs = append(errs, errors.New(fmt.Sprintf("Expected %s to be a number.", fieldName)))
					} else {
						app.Set(fieldName, value)
					}
				}
			}

			if app.Has("services") {
				appServices, err := servicesComponent(app.Get("services"))
				if err != nil {
					errs = append(errs, errors.New("Expected local services to be an array of service instance names."))
				} else {
					app.Set("services", appServices)
				}
			} else {
				app.Set("services", []string{})
			}

			if app.Has("env") {
				env, ok := app.Get("env").(map[string]interface{})
				if !ok {
					errs = append(errs, errors.New("Expected local env vars to be a set of key => value."))
				} else {
					app.Set("env", generic.NewMap(env))
				}
			} else {
				app.Set("env", generic.NewMap())
			}
		}
	}

	if data.Has("env") {
		env, ok := data.Get("env").(map[string]interface{})
		if ok {
			m.GlobalEnvVars = generic.NewMap(env)
		} else {
			errs = append(errs, errors.New("Expected global env vars to be a set of key => value."))
		}
	}

	if data.Has("services") {
		globalServices, err := servicesComponent(data.Get("services"))
		if err != nil {
			errs = append(errs, errors.New("Expected global services to be an array of service instance names."))
		} else {
			m.GlobalServices = globalServices
		}
	}

	return
}

func validateAppParams(appParams cf.AppParams) (errs ManifestErrors) {
	keysToCheck := []string{"name", "command", "space_guid", "buildpack", "disk_quota", "instances", "memory"}
	for _, key := range keysToCheck {
		if appParams.Has(key) && appParams.Get(key) == nil {
			errs = append(errs, errors.New(key+" should not be null"))
		}
	}

	return
}

func servicesComponent(input interface{}) (result []string, err error) {
	switch input := input.(type) {
	case []interface{}:
		for _, value := range input {
			stringValue, ok := value.(string)
			if !ok {
				err = errors.New("validation error")
				return
			}
			result = append(result, stringValue)
		}
	default:
		err = errors.New("validation error")
		return
	}
	return
}

func mergeSets(set1, set2 []string) (result []string) {
	for _, aString := range set1 {
		result = append(result, aString)
	}

	for _, aString := range set2 {
		if !stringInSlice(aString, result) {
			result = append(result, aString)
		}
	}
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
