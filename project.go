package weiqi

import (
    "database/sql"
    "github.com/dgf1988/weiqi/db"
)

/*
CREATE TABLE `project` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`name` CHAR(50) NOT NULL,
	`textid` INT(11) NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=2
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

type Project struct {
    Id int64
    Name string
    Text string
    Items []Item
}

func (p *Project) AppendItem(item Item) {
    p.Items = append(p.Items, item)
}

func (p *Project) AddItem(key, value string) {
    p.Items = append(p.Items, Item{Key:key, Value:value})
}

func (p *Project) Save() (id int64, err error) {
    if p.Id > 0 {
        //先更新名称
        if _, err = Db.Project.Update(p.Id).Values(nil, p.Name); err != nil {
            return
        }

        //再更新文本
        var textid int64
        if err = Db.Project.Get(p.Id).Scan(nil, nil, &textid); err == nil {
            if _, err = Db.Text.Update(textid).Values(nil, p.Text); err != nil {
                return
            }
        } else if err == sql.ErrNoRows {
            if textid, err = Db.Text.Add(nil, p.Text); err != nil {
                return
            }
            if _, err = Db.Project.Update(p.Id).Values(nil, nil, textid); err != nil {
                return
            }
        } else {
            return
        }
    } else {
        if err = Db.Project.Get(nil, p.Name).Scan(&p.Id); err == nil {
            //名称已经存在
            return p.Save()
        } else if err != sql.ErrNoRows {
            return
        }

        //添加文本
        var textid int64
        if textid, err = Db.Text.Add(nil, p.Text); err != nil {
            return
        }
        //再添加项目
        if id, err = Db.Project.Add(nil, p.Name, textid); err != nil {
            return
        }
    }

    var itemid int64
    for _, item := range p.Items {
        itemid, err = item.Save()
        if err != nil {
            return
        }
        if itemid > 0 {
            _, err = Db.ProjectItem.Add(nil, p.Id, itemid)
            if err != nil {
                return
            }
        }
    }
    return
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
                if err = Db.Item.Get(itemid).Struct(&item); err == nil {
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

