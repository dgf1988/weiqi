package weiqi

import (
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	SessionCookieName = "WEIQI"
)

var (
	//锁
	SessionLocker sync.RWMutex
	//会话存储
	Sessions map[string]*Session = make(map[string]*Session)
)

type Session struct {
	Id      string
	User    *User
	Timeout time.Time
}

func newSession(u *User) *Session {
	return &Session{
		sessionId(),
		u,
		time.Now().AddDate(0, 0, 1),
	}
}

func logSession() {
	for _, s := range Sessions {
		logDebug(s.Id, s.Timeout, s.User)
	}
}

func getSession(r *http.Request) *Session {
	s := &Session{}
	c, _ := r.Cookie(SessionCookieName)
	if c == nil {
		return nil
	}
	s.Id = c.Value

	SessionLocker.Lock()
	defer SessionLocker.Unlock()
	s_get, ok := Sessions[s.Id]
	//id 相同， 并且user不为nil
	if !ok || s_get.Id != s.Id || s_get.User == nil {
		return nil
	}
	s.Timeout = s_get.Timeout
	s.User = s_get.User
	return s
}

func clearSession(w http.ResponseWriter, r *http.Request) {
	SessionLocker.Lock()
	defer SessionLocker.Unlock()

	c, _ := r.Cookie(SessionCookieName)
	if c != nil {
		c.MaxAge = -1
		c.Expires = time.Now().AddDate(0, -1, 0)
		http.SetCookie(w, c)
	}

	_, ok := Sessions[c.Value]
	if ok {
		delete(Sessions, c.Value)
	}
}

func clearSessionMany() int {
	SessionLocker.Lock()
	defer SessionLocker.Unlock()

	keys := make([]string, 0)
	for k, s := range Sessions {
		if s.Timeout.Before(time.Now()) {
			keys = append(keys, k)
		}
	}
	for i := range keys {
		delete(Sessions, keys[i])
	}
	return len(keys)
}

func sessionId() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	return getMd5(strconv.FormatInt(nano, 10) + strconv.FormatInt(rndNum, 10))
}

func (s *Session) Add(w http.ResponseWriter) {
	SessionLocker.Lock()
	defer SessionLocker.Unlock()

	http.SetCookie(w, &http.Cookie{Name: SessionCookieName, Value: s.Id, Expires: s.Timeout})
	Sessions[s.Id] = s
}

//快速获取会话中的用户。
//并不是什么时候都需要持有一个会话对象。
func getSessionUser(r *http.Request) *User {
	s := getSession(r)
	if s != nil {
		return s.User
	}
	return nil
}
