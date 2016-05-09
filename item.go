package weiqi

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
*/

type Item struct {
    Id int64
    Key string
    Value string
}

func (i *Item) Save() (insertid int64, err error) {
    if i.Id > 0 {
        _, err = Db.Item.Update(i.Id).Values(nil, i.Key, i.Value)
        return
    } else {
        return Db.Item.Add(nil, i.Key, i.Value)
    }
}

func getItem(id int64) (item *Item, err error) {
    item = new(Item)
    err = Db.Item.Get(id).Struct(item)
    return
}