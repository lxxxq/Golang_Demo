package repository

import (
	"demo/entity"
	mysql "demo/mysql_common/mysql"
	"log"
)

func SelectByID(id uint64) entity.Student {
	var student entity.Student
	db := mysql.RDBs["mysql"]

	err := db.Db.Table("student_info").Where("id = ?", id).Find(&student).Error

	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return student
}

func SelectAll(page int, pageSize int, name string) ([]entity.Student, int64, error) {
	var students []entity.Student
	var count int64

	d := mysql.RDBs["mysql"]
	d.Db.Table("student_info").Count(&count)

	if name != "" {
		d.Db = d.Db.Table("student_info").Where("name LIKE ?", "%"+name+"%")
	}

	d.Db.Table("student_info").Limit(pageSize).Offset((page) * pageSize).Find(&students)

	return students, count, nil
}

func DeleteByID(id uint64) error {
	db := mysql.RDBs["mysql"]

	err := db.Db.Table("student_info").Where("id = ?", id).Delete(&entity.Student{}).Error

	if err != nil {
		log.Printf("error: %v\n", err)
	}
	return err
}

func UpdateByID(students []entity.Student) error {
	db := mysql.RDBs["mysql"]
	tx := db.Db.Begin()

	for _, student := range students {
		err := db.Db.Table("student_info").Where("id = ?", student.ID).
			Updates(entity.Student{Name: student.Name, Age: student.Age}).Error

		if err != nil {
			tx.Rollback()
			log.Printf("error: %v\n", err)

			return err
		}
	}
	tx.Commit()
	return nil
}

func InsetData(students []entity.Student) error {
	db := mysql.RDBs["mysql"]
	tx := db.Db.Begin()

	err := db.Db.Table("student_info").Create(students).Error

	if err != nil {
		tx.Rollback()
		log.Printf("error: %v\n", err)
		return err
	}
	tx.Commit()

	return nil
}
