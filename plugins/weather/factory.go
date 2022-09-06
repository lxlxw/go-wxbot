package weather

import (
	"errors"
)

type Info interface {
	GetInfo(url string, locationID string) string
}

func Factory(urlType string) (info Info, err error) {
	switch urlType {
	case "空气指数":
		info = &AirInfo{}
	case "生活指数":
		info = &LifeInfo{}
	case "日出日落":
		info = &SunInfo{}
	case "气象预警":
		info = &WarningInfo{}
	case "实时天气":
		info = &ActualInfo{}
	default:
		errors.New("factory error")
	}
	return info, nil
}
