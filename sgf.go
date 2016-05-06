package weiqi

import (
	"errors"
	"github.com/dgf1988/mahonia"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"fmt"
)
/*
	tm timelimit
	
*/

/*
CREATE TABLE `sgf` (
`id` INT(11) NOT NULL AUTO_INCREMENT,
`stime` DATE NOT NULL DEFAULT '0000-00-00',
`splace` CHAR(50) NOT NULL DEFAULT '',
`sevent` CHAR(100) NOT NULL DEFAULT '',
`sblack` CHAR(50) NOT NULL DEFAULT 'b',
`swhite` CHAR(50) NOT NULL DEFAULT 'w',
`srule` CHAR(50) NOT NULL DEFAULT '',
`sresult` CHAR(50) NOT NULL DEFAULT '',
`ssteps` MEDIUMTEXT NOT NULL,
`supdate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
*/

const (
	c_SgfCharset = "utf-8"
)

type Sgf struct {
	Id     int64
	Time   time.Time
	Place  string
	Event  string
	Black  string
	White  string
	Rule   string
	Result string
	Sgf    string
	Update time.Time
}

func (this Sgf) ToSgf() string {
	if this.Sgf == "" {
		return ""
	}
	items_sgf := make([]string, 0)
	items_sgf = append(items_sgf, "(;")
	items_sgf = append(items_sgf, fmt.Sprintf("CA[%s]EV[%s]", c_SgfCharset, this.Event))
	if this.Time.IsZero() {
		items_sgf = append(items_sgf, "DT[0000-00-00]")
	} else {
		items_sgf = append(items_sgf, fmt.Sprintf("DT[%s]", this.Time.Format("2006-01-02")))
	}
	items_sgf = append(items_sgf, fmt.Sprintf("PC[%s]", this.Place))
	items_sgf = append(items_sgf, fmt.Sprintf("PB[%s]", this.Black))
	items_sgf = append(items_sgf, fmt.Sprintf("PW[%s]", this.White))
	items_sgf = append(items_sgf, fmt.Sprintf("RL[%s]", this.Rule))
	items_sgf = append(items_sgf, fmt.Sprintf("RE[%s]", this.Result))
	items_sgf = append(items_sgf, "\n")
	items_sgf = append(items_sgf, this.Sgf)
	items_sgf = append(items_sgf, ")")
	return strings.Join(items_sgf, "")
}

type sgfOrderByTimeDesc []Sgf

func (arr sgfOrderByTimeDesc) Len() int           { return len(arr) }
func (arr sgfOrderByTimeDesc) Swap(i, j int)      { arr[i], arr[j] = arr[j], arr[i] }
func (arr sgfOrderByTimeDesc) Less(i, j int) bool { return arr[i].Time.After(arr[j].Time) }

func listSgfOrderByTimeDesc(take, skip int) ([]Sgf, error) {
	var sgfs = make([]Sgf, 0)
	if rows, err := Db.Sgf.Query("order by sgf.time desc limit ?, ?", skip, take); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var sgf Sgf
			if err = rows.Struct(&sgf); err != nil {
				return nil, err
			} else {
				sgfs = append(sgfs, sgf)
			}
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return sgfs, nil
}

func listSgfByNameOrderByTimeDesc(name string) ([]Sgf, error) {
	var sgfs = make([]Sgf, 0)
	if rows, err := Db.Sgf.FindAny(nil, nil, nil, nil, name, name); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var sgf Sgf
			if err = rows.Struct(&sgf); err != nil {
				return nil, err
			} else {
				sgfs = append(sgfs, sgf)
			}
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	sort.Sort(sgfOrderByTimeDesc(sgfs))
	return sgfs, nil
}

func listSgfByNamesOrderByTimeDesc(names ...string) ([]Sgf, error) {
	var sgfs = make([]Sgf, 0)
	for _, name := range names {
		if name == "" {
			continue
		}
		if rows, err := Db.Sgf.Query("where sgf.black = ? or sgf.white = ? order by sgf.time desc", name, name); err != nil {
			return nil, err
		} else {
			defer rows.Close()
			for rows.Next() {
				var sgf Sgf
				if err = rows.Struct(&sgf); err != nil {
					return nil, err
				} else {
					sgfs = append(sgfs, sgf)
				}
			}
			if err = rows.Err(); err != nil {
				return nil, err
			}
		}
	}
	sort.Sort(sgfOrderByTimeDesc(sgfs))
	return sgfs, nil
}

var testsgf = `
(
TE[700网晚报杯中国业余赛第13轮]
RD[2016-04-29]
PC[无锡华美达酒店]
TM[85]
LT[10]
LC[1]
KO[7.5]
RE[黑中盘胜]
PB[唐崇哲]
BR[7段]
PW[李岱春]
WR[8段]
GK[1]
TC[]

;B[qd];W[dp];B[cd];W[qp];B[oc];W[ed];B[ec];W[fc];B[dc];W[gd];B[fq];W[cn];B[cf];W[kc]
;B[oq];W[jq];B[pm];W[po];B[lq];W[pk];B[qq];W[rq];B[pq];W[ro];B[dr];W[cq];B[qi];W[rk]
;B[oo];W[on];B[no];W[nn];B[mn];W[mm];B[ln];W[gp];B[gq];W[hp];B[fp];W[fo];B[cr];W[br]
;B[dq];W[cp];B[hq];W[ip];B[oi];W[lm];B[kn];W[nl];B[mc];W[ke];B[ok];W[ol];B[pj];W[qk]
;B[nf];W[ci];B[dj];W[ir];B[hr];W[lp];B[mp];W[ep];B[gs];W[di];B[ej];W[bk];B[ei];W[dg]
;B[dl];W[cg];B[bf];W[dd];B[cc];W[cl];B[hi];W[lr];B[mq];W[im];B[gm];W[jl];B[hk];W[me]
;B[ne];W[sh];B[rh];W[hf];B[lb];W[kb];B[ig];W[hg];B[ih];W[if];B[ld];W[kd];B[bq];W[bp]
;B[bs];W[ar];B[kr];W[kq];B[rr];W[sr];B[rs];W[si];B[sg];W[rj];B[ri];W[sj];B[rg];W[md]
;B[nd];W[lc];B[mb];W[fr];B[cs];W[es];B[eq];W[dm];B[el];W[df];B[eb];W[fb];B[fa];W[ga]
;B[ea];W[ka];B[jf];W[je];B[mj];W[mr];B[nr];W[ks];B[mf];W[le];B[fe];W[fd];B[bg];W[fm]
;B[fl];W[gn];B[em];W[en];B[hm];W[hn];B[jk];W[kk];B[jj];W[kg];B[il];W[jm];B[fg];W[gh]
;B[eh];W[pp];B[kp];W[bh];B[ge];W[he];B[lg];W[lh];B[kh];W[ki];B[jh];W[li];B[lk];W[kj]
;B[kl];W[mg];B[lf];W[ll];B[km];W[mk];B[mi];W[lj];B[in];W[mh];B[nh];W[jo];B[jn];W[ko]
;B[lo];W[kf];B[cj];W[bj];B[io];W[jp];B[ag];W[gi];B[gj];W[is];B[hs];W[jg];B[ah];W[ai]
;B[dh];W[ch];B[ck];W[nj];B[ni];W[op];B[np];W[fi];B[fj];W[la];B[ma];W[ng];B[og];W[ji]
;B[ii];W[lk];B[ms];W[jr];B[lp])
`
var testsgf2 = `
(;EV[第64期日本王座战预选]DT[2016-03-10]PC[]TM[3小时]LT[60]LC[5]KO[6.5]RE[白中盘胜]PB[大场惇也]BR[六段]PW[王立诚]WR[九段]GK[1]TC[]
;B[pd];W[dc];B[dp];W[pp];B[ce];W[ed];B[ci];W[nc];B[qf];W[pb];B[nq];W[qn];B[gq];W[qc];B[pj];W[cp];B[co];W[bo];B[do];W[bq];B[bn];W[dq];B[ao];W[bp];B[eq];W[er];B[dr];W[cq];B[fq];W[cr];B[no];W[pr];B[fr];W[ds];B[pm];W[qm];B[pl];W[pn];B[on];W[il];B[hl];W[im];B[ik];W[hm];B[gl];W[gm];B[fm];W[fn];B[fl];W[kl];B[jk];W[ko];B[io];W[mn];B[jn];W[kn];B[jm];W[jl];B[km];W[lm];B[ll];W[cn];B[dn];W[nn];B[oo];W[ml];B[lk];W[om];B[nm];W[ol];B[nl];W[ok];B[nk];W[oj];B[kq];W[jq];B[jo];W[kp];B[lq];W[kk];B[lj];W[kj];B[li];W[jh];B[kh];W[nj];B[mk];W[pi];B[jg];W[ih];B[ig];W[hh];B[ij];W[ki];B[lg];W[gj];B[hg];W[hk];B[fc];W[fd];B[gc];W[gh];B[fi];W[gi];B[gg];W[fh];B[ef];W[cd];B[be];W[bd];B[eb];W[db];B[de];W[bb];B[lc];W[qh];B[cm];W[ne];B[oe];W[od];B[of];W[nf];B[qd];W[rc];B[oh];W[mh];B[lh];W[nh];B[qg];W[rh];B[og];W[mf];B[go];W[ei];B[oi];W[qj];B[ph];W[pk];B[rd];W[jc];B[kd];W[jd];B[je];W[ke];B[le];W[gd];B[lb];W[mb];B[hc];W[ie];B[kf];W[id];B[oc];W[nd];B[pc];W[ob];B[rb];W[qb];B[sc];W[ib];B[hb];W[ma];B[hd];W[he];B[ge];W[ec];B[ha];W[ia];B[fa];W[gb];B[kc];W[fb];B[ga];W[fb];B[da];W[ea];B[bs];W[ap];B[eb];W[ca];B[gb];W[ea];B[ar];W[br];B[da];W[ff];B[gf];W[ea];B[fs];W[es];B[da];W[la];B[jb];W[ea];B[cs];W[as];B[da];W[ic];B[kb];W[ea];B[ni];W[sb];B[da];W[sd];B[if];W[ea];B[ej];W[fj];B[da];W[mp];B[oq];W[mq];B[mr];W[nr];B[lr];W[or];B[mo];W[lp];B[jr];W[ea];B[eh];W[dh])
`

func parseSgf(strSgf string) (*Sgf, error) {
	if strSgf == "" {
		return nil, errors.New("sgf: 棋谱为空")
	}
	strSgf = strings.TrimSpace(strSgf)

	var sgf Sgf
	var re = regexp.MustCompile(`([A-Za-z]{1,2})\[([^\]]+)\]`)

	var titles = regexp.MustCompile(`([A-Za-z]{2}\[[^\]]+\])`).FindAllString(strSgf, -1)
	for _, title := range titles {
		var match = re.FindStringSubmatch(title)
		if len(match) != 3 {
			continue
		}
		switch match[1] {
		case "TE", "EV":
			sgf.Event = match[2]
		case "RD", "DT":
			sgf.Time, _ = time.Parse("2006-01-02", match[2])
		case "PC":
			sgf.Place = match[2]
		case "PB":
			sgf.Black = match[2]
		case "PW":
			sgf.White = match[2]
		case "RE":
			sgf.Result = match[2]
		case "TM":
			sgf.Rule += "用时：" + match[2] + "  "
		case "LT":
			sgf.Rule += "读秒：" + match[2] + "  "
		case "LC":
			sgf.Rule += match[2] + "次"
		}
	}
	sgf.Sgf = strSgf
	return &sgf, nil
}

func httpGetSgf(src string, charset string) (*Sgf, error) {
	var strsgf string
	var err error
	var code int
	if strsgf, code, err = httpGetString(src); err != nil {
		return nil, err
	} else if code != 200 {
		return nil, errors.New("sgf: remote: code: " + strconv.Itoa(code))
	}
	var coding = mahonia.NewDecoder(charset)
	var sgf *Sgf
	if sgf, err = parseSgf(coding.ConvertString(strsgf)); err != nil {
		return nil, err
	}
	return sgf, nil
}
