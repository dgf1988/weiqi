package weiqi

import (
    "database/sql"
    "github.com/dgf1988/db"
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
        if p.Id, err = Db.Project.Add(nil, p.Name, textid); err != nil {
            return
        }
        id = p.Id
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

func DelProject(id int64) (err error) {
    var textid int64
    err = Db.Project.Get(id).Scan(nil, nil, &textid)
    if err == nil {
        //删除文本
        _, err = Db.Text.Del(textid)
        if err != nil {
            return
        }
        var rows *db.Rows
        if rows, err = Db.ProjectItem.FindMany(nil, id); err != nil {
            return
        } else {
            defer rows.Close()
            var itemid int64
            for rows.Next() {
                err = rows.Scan(nil, nil, &itemid)
                if err != nil {
                    return
                }
                //删除
                _, err = Db.Item.Del(itemid)
                if err != nil {
                    return
                }
            }
            err = rows.Err()
            if err != nil {
                return
            }
        }
        //删除
        _, err = Db.ProjectItem.Del(nil, id)
        if err != nil {
            return
        }
        _, err = Db.Project.Del(id)
        return
    }
    return
}

func GetProject(id int64) (*Project, error ) {
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

    var rows *db.Rows
    if rows, err = Db.ProjectItem.FindMany(nil, project.Id); err != nil {
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

func ListProject(take, skip int) (listproject []Project, err error) {
    listproject = make([]Project, 0)
    var rows *db.Rows
    rows, err = Db.Project.List(take, skip)
    if err != nil {
        return
    }
    defer rows.Close()
    for rows.Next() {
        var p Project

        var textid int64
        err = rows.Scan(&p.Id, &p.Name, &textid)
        if err != nil {
            return
        }
        err = Db.Text.Get(textid).Scan(nil, &p.Text)
        if err != nil && err != sql.ErrNoRows {
            return
        }

        var projectitemrows *db.Rows
        projectitemrows, err = Db.ProjectItem.FindMany(nil, p.Id)
        if err != nil {
            return
        }
        defer projectitemrows.Close()
        for projectitemrows.Next() {
            var itemid int64
            var projectitemid int64
            err = projectitemrows.Scan(&projectitemid, nil, &itemid)
            if err != nil {
                return
            }
            var item Item
            err = Db.Item.Get(itemid).Scan(&item.Id, &item.Key, &item.Value)
            if err == nil {
                p.AppendItem(item)
            } else if err == sql.ErrNoRows {
                _, err = Db.ProjectItem.Del(projectitemid)
                if err != nil {
                    return
                }
            } else {
                return
            }
        }
        err = projectitemrows.Err()
        if err != nil {
            return
        }
        listproject = append(listproject, p)
    }
    err = rows.Err()
    return
}

