package sqlAssister

import (
	"log"
	"reflect"
)

func (ac *Assister) ScanStruct(arg any) (any, error) {
	vals := reflect.ValueOf(arg)
	vals.NumField()
	numOfFields := vals.NumField()

	for i := 0; i < numOfFields; i++ {
		var z interface{}
		row, err := ac.SingleRowScannerWithArgs("")
		if err != nil {
			log.Println(err)
		}

		err = row.Scan(&z)
		if err != nil {
			return nil, err
		}
		vals.Field(i).Set(reflect.ValueOf(z))
	}

	return vals, nil
}
