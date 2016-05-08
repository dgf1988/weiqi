package weiqi

import (
    "database/sql"
    "github.com/dgf1988/weiqi/db"
    "fmt"
)

/*
CREATE TABLE `item` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`key` CHAR(40) NOT NULL,
	`vlaue` VARCHAR(255) NOT NULL DEFAULT '',
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
CREATE TABLE `project` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`name` CHAR(50) NOT NULL,
	`text` INT(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
CREATE TABLE `projectitem` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`projectid` INT(11) NOT NULL,
	`itemid` INT(11) NOT NULL,
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

*/

type Item struct {
    Key string
    Value string
}

type Project struct {
    Id int64
    Name string
    Text string
    Items []Item
}

func (p *Project) AddItem(key, value string) {
    p.Items = append(p.Items, Item{key, value})
}

func addProject(project Project) (int64, error) {
    var err error
    if project.Id <= 0 {
        var textid int64
        textid, err = Db.Text.Add(nil, project.Text)
        if err != nil {
            return -1, err
        }

        var id int64
        id, err = Db.Project.Add(nil, project.Name, textid)
        if err != nil {
            return -1, err
        }

        var itemid int64
        for _, item := range project.Items {
            itemid, err = Db.Item.Add(nil, item.Key, item.Value)
            if err != nil {
                return -1, err
            }
            _, err = Db.ProjectItem.Add(nil, id, itemid)
            if err != nil {
                return -1, err
            }
        }
        return id, nil
    }
    return -1, fmt.Errorf("add project: id error %d", project.Id)
}

func getProject(id int64) (*Project, error ) {
    var err error

    var project Project
    var textid int64
    if err = Db.Project.Get(id).Scan(&project.Id, &project.Name, &textid); err != nil {
        return nil, err
    }

    var text Text
    if err = Db.Text.Get(textid).Struct(&text); err == sql.ErrNoRows {

    } else if err == nil {
        project.Text = text.Text
    } else {
        return nil, err
    }

    var rows db.Rows
    if rows, err = Db.ProjectItem.FindAll(nil, project.Id); err != nil {
        return nil, err
    } else {
        defer rows.Close()
        for rows.Next() {
            var itemid int64
            if err = rows.Scan(nil, nil, &itemid); err != nil {
                return nil, err
            } else {
                var item Item
                if err = Db.Item.Get(itemid).Scan(nil, &item.Key, &item.Value); err == nil {
                    project.Items = append(project.Items, item)
                } else if err == sql.ErrNoRows {
                    continue
                } else {
                    return nil, err
                }
            }
        }
        if err = rows.Err(); err != nil {
            return nil, err
        }
    }
    return &project, nil
}

