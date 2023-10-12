package mysql

import (
	"errors"

	"gorm.io/gorm"
)

func WrapInsertSingle(table string, value interface{}) error {
	DB := GetDB()
	result := DB.Table(table).Create(value)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapInsertBatch(table string, value interface{}, batchSize int) error {
	DB := GetDB()
	result := DB.Table(table).CreateInBatches(value, batchSize)
	if result == nil || result.Error != nil || result.RowsAffected != int64(batchSize) {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapFindFirst(table string, v interface{}, condition map[string]interface{}) (error, bool) {
	DB := GetDB()
	result := DB.Table(table).Where(condition).First(v)
	if result == nil || (result.Error != nil && result.Error != gorm.ErrRecordNotFound) {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg), false
	}
	return nil, result.Error != gorm.ErrRecordNotFound
}

//query use map now.actually struct/string also ok?
func WrapUpdateByCondition(table string, query interface{}, updateMap map[string]interface{}) error {
	DB := GetDB()
	result := DB.Table(table).Where(query).Updates(updateMap)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}

	return nil
}

func WrapFindAllByQueryCondition(table string, query string, condition []interface{}, v interface{}) error {
	DB := GetDB()
	result := DB.Table(table).Where(query, condition).Find(v)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapFindAllByConditionCheckOrderLimit(table string, condition map[string]interface{}, v interface{}, order interface{}, num int, count *int64) error {
	DB := GetDB()
	result := DB.Table(table).Where(condition).Count(count)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}

	if *count == 0 {
		return nil
	}

	r := result.Order(order).Limit(num).Find(v)
	if r == nil || r.Error != nil {
		errmsg := ""
		if r == nil {
			errmsg = "res is nil"
		} else {
			errmsg = r.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapFindAllByCondition(table string, condition map[string]interface{}, v interface{}) error {
	DB := GetDB()
	result := DB.Table(table).Where(condition).Find(v)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapFindAllByConditionOffsetLimit(table string, condition map[string]interface{}, v interface{}, order interface{}, offset, num int) error {
	DB := GetDB()
	count := int64(0)
	result := DB.Table(table).Where(condition).Count(&count)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}

	if count == 0 {
		return nil
	}

	r := result.Order(order).Offset(offset).Limit(num).Find(v)
	if r == nil || r.Error != nil {
		errmsg := ""
		if r == nil {
			errmsg = "res is nil"
		} else {
			errmsg = r.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil

}

func WrapCountByCondition(table string, condition map[string]interface{}, v *int64) error {
	DB := GetDB()
	result := DB.Table(table).Where(condition).Count(v)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapFindAll(table string, v interface{}) error {
	DB := GetDB()
	result := DB.Table(table).Find(v)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapDeleteByPrikey(table string, v interface{}, id uint64) error {
	DB := GetDB()
	result := DB.Table(table).Delete(v, id)
	if result == nil || result.Error != nil || int(result.RowsAffected) != 1 {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapBeginTx() *gorm.DB {
	DB := GetDB()
	return DB.Begin()
}

func WrapCountByConditionFilter(table string, condition map[string]interface{}, query string, filter []interface{}, v *int64) error {
	DB := GetDB()
	result := DB.Table(table).Where(condition).Where(query, filter).Count(v)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapFindAllByConditionFilterOffsetLimit(table string, condition map[string]interface{}, query string, filter []interface{}, v interface{}, order interface{}, offset, num int) error {
	DB := GetDB()
	r := DB.Table(table).Where(condition).Where(query, filter).Order(order).Offset(offset).Limit(num).Find(v)
	if r == nil || r.Error != nil {
		errmsg := ""
		if r == nil {
			errmsg = "res is nil"
		} else {
			errmsg = r.Error.Error()
		}

		return errors.New(errmsg)
	}
	return nil
}

func WrapFindAllByQueryConditionFilter(table string, query string, condition []interface{}, filter map[string]interface{}, v interface{}) error {
	DB := GetDB()
	result := DB.Table(table).Where(filter).Where(query, condition).Find(v)
	if result == nil || result.Error != nil {
		errmsg := ""
		if result == nil {
			errmsg = "res is nil"
		} else {
			errmsg = result.Error.Error()
		}
		return errors.New(errmsg)
	}
	return nil
}
