package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/pa-pe/luca-control/config"
	"github.com/pa-pe/luca-control/src/utils"
	"github.com/pa-pe/luca-control/src/web/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Structure for storing login attempt data
type loginAttempt struct {
	Count       int
	LastAttempt time.Time
}

// Global variable for storing information about login attempts
var loginAttempts sync.Map

func AuthRequired(db *gorm.DB, isFirstRun *bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Redirect to the initial setup page if it's the first run
		if *isFirstRun {
			c.Redirect(http.StatusSeeOther, "/initial-setup")
			c.Abort()
			return
		}

		isAuthorized, _, _ := GetAuth(c, db)
		if !isAuthorized {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetAuth(c *gin.Context, db *gorm.DB) (bool, *models.WebUser, *models.WebSession) {
	cookie, err := c.Cookie("session")
	if err != nil {
		return false, nil, nil
	}

	isValid, err, session := checkAndGetSession(db, cookie)
	if err != nil || !isValid {
		log.Printf("[Web Auth] Fail cookie, ip=%s", c.ClientIP())
		return false, nil, nil
	}

	c.Set("currentAuthSession", session)

	var currentAuthUser models.WebUser
	if err := db.Where("id = ?", session.WebUserID).First(&currentAuthUser).Error; err != nil {
		return false, nil, nil
	}
	c.Set("currentAuthUser", currentAuthUser)

	return true, &currentAuthUser, session
}

func ShowLoginPage(c *gin.Context, db *gorm.DB) {
	isAuthorized, _, _ := GetAuth(c, db)
	if isAuthorized {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.HTML(http.StatusOK, "login.html", nil)
}

func HandleLogin(c *gin.Context, db *gorm.DB) {
	ip := c.ClientIP()

	if !isLoginAllowed(ip) {
		c.HTML(http.StatusTooManyRequests, "login.html", gin.H{"error": "Too many login attempts. Please try again later."})
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.WebUser
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		attempts := registerFailedLogin(ip)
		log.Printf("[Web Auth] Fail username=%s, ip=%s, attempts=%d", username, ip, attempts)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckStrHash(password, user.Password) {
		attempts := registerFailedLogin(ip)
		log.Printf("[Web Auth] Fail password for username=%s, ip=%s, attempts=%d", username, ip, attempts)
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	// Successful authorization, clearing data on attempts
	loginAttempts.Delete(ip)

	sessionID, sessionKey, err := createSession(db, user.ID)
	if err != nil {
		log.Printf("[Web Auth] Fail with internal server error while createSession, username=%s, ip=%s", username, ip)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Unable to create session"})
		return
	}

	log.Printf("[Web Auth] Success username=%s, ip=%s", username, ip)

	user.Password = "hidden"
	c.Set("currentAuthUser", user)

	// Setting a cookie with a session ID and key
	cookieValue := fmt.Sprintf("%d:%s", sessionID, sessionKey)
	c.SetCookie("session", cookieValue, config.WebAuthSessionDurationInHour*60*60, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func HandleLogout(c *gin.Context, db *gorm.DB) {
	isAuthorized, currentAuthUser, currenSession := GetAuth(c, db)
	if isAuthorized {
		//db.Where("id < ?", currenSession.ID).Delete(currenSession)
		db.Delete(currenSession)
		c.SetCookie("session", "", -1, "/", "", false, true) // Удаление куки
		log.Printf("[Web Auth] Logout username=%s, ip=%s", currentAuthUser.Username, c.ClientIP())
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func generateSessionKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), nil
}

func createSession(db *gorm.DB, userID int) (int, string, error) {

	sessionKey, err := generateSessionKey()
	if err != nil {
		return 0, "", err
	}

	hashedKey, err := utils.HashStr(sessionKey)

	if err != nil {
		return 0, "", err
	}

	expiresAt := time.Now().Add(config.WebAuthSessionDurationInHour * time.Hour)

	session := models.WebSession{
		WebUserID:  userID,
		SessionKey: hashedKey,
		CreatedAt:  time.Now(),
		ExpiresAt:  expiresAt,
	}

	if err := db.Create(&session).Error; err != nil {
		return 0, "", err
	}

	return session.ID, sessionKey, nil
}

func checkAndGetSession(db *gorm.DB, cookie string) (bool, error, *models.WebSession) {
	// Splitting cookies into sessionID and sessionKey
	parts := strings.SplitN(cookie, ":", 2)
	if len(parts) != 2 {
		return false, nil, nil
	}

	sessionID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return false, nil, nil
	}

	sessionKey := parts[1]

	var session models.WebSession
	if err := db.Where("id = ?", sessionID).Take(&session).Error; err != nil {
		return false, err, nil
	}

	if time.Now().After(session.ExpiresAt) {
		return false, nil, nil
	}

	if utils.CheckStrHash(sessionKey, session.SessionKey) != true {
		return false, nil, nil
	}

	go cleanUpExpiredSessions(db)

	return true, nil, &session
}

func cleanUpExpiredSessions(db *gorm.DB) {
	db.Where("expires_at < ?", time.Now()).Delete(&models.WebSession{})
}

// Checking if login attempt is allowed
func isLoginAllowed(ip string) bool {
	// Getting information about attempts by IP address
	if attempt, ok := loginAttempts.Load(ip); ok {
		attemptData := attempt.(loginAttempt)

		// Reset the expired attempts
		if time.Since(attemptData.LastAttempt) > config.WebAuthAttemptResetDuration {
			loginAttempts.Delete(ip)
			return true
		}

		if attemptData.Count >= config.WebAuthMaxAttempts {
			return false
		}
	}

	return true
}

func registerFailedLogin(ip string) int {
	attemptData := loginAttempt{Count: 1, LastAttempt: time.Now()}

	// update existingAttempt
	if attempt, ok := loginAttempts.Load(ip); ok {
		existingAttempt := attempt.(loginAttempt)
		attemptData.Count = existingAttempt.Count + 1
		attemptData.LastAttempt = time.Now()
	}

	loginAttempts.Store(ip, attemptData)

	return attemptData.Count
}

func GetCurrentAuthUser(c *gin.Context) models.WebUser {
	cCurrentAuthUser, exists := c.Get("currentAuthUser")
	if !exists {
		log.Fatal("GetCurrentAuthUser internal error")
	}
	currentAuthUser, ok := cCurrentAuthUser.(models.WebUser)
	if !ok {
		log.Fatal("GetCurrentAuthUser internal error with models.WebUser")
	}
	return currentAuthUser
}
