package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `json:"name"`
	APIKey    string    `gorm:"uniqueIndex;size:255" json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type RateLimit struct {
	Key         string `gorm:"primaryKey; size:255"`
	Count       int
	WindowStart time.Time `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func initDG() *gorm.DB {
	dsn := "host=localhost user=psssh password=p1223. dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	db.AutoMigrate(&User{}, &RateLimit{})
	return db
}

// here we generate a random API key for the user. The key is a combination of a prefix "API-" and a random 16-byte string encoded in hexadecimal format.
// This ensures that each API key is unique and difficult to guess, providing a secure way to identify users when they make API requests.
func generateAPIKey() string {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	return "API-" + hex.EncodeToString(bytes)
}

// RateLimiter is a FACTORY FUNCTION — it doesn't run on requests itself.
// Instead, it CONFIGURES and RETURNS a middleware that will run on every request.
//
// Think of it like a stamp maker: you tell it "make me a stamp that allows
// 100 requests per minute", and it hands you back that stamp (the middleware).
// The stamp is then applied to every incoming request.
//
// The parameters here become permanently "baked into" the returned middleware
// via Go closures — the inner functions will remember these values even after
// RateLimiter() itself has finished executing.
func RateLimiter(db *gorm.DB, limit int, window time.Duration, exceedStatus int) echo.MiddlewareFunc {

	// LAYER 1: This function receives `next`, which is the next handler in
	// Echo's middleware chain. Middleware in Echo works like a chain of guards —
	// each guard decides whether to let the request pass to the next guard,
	// or to reject it right there. `next` is essentially "the rest of the chain".
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		// LAYER 2: THIS is the function that runs on every single HTTP request.
		// `c` (echo.Context) is a bundle containing everything about the current
		// request (headers, body, URL params) AND the response writer.
		// Returning an error or calling c.JSON() sends a response back to the client.
		return func(c echo.Context) error {

			// ── STEP 1: IDENTIFY THE CLIENT ──────────────────────────────────
			// We need a unique identifier per client to track their request count.
			// We use the "X-API-Key" HTTP header for this. The "X-" prefix is a
			// convention meaning "custom/non-standard header".
			key := c.Request().Header.Get("X-API-Key")

			// If no key is provided, we can't track the client at all —
			// so we reject the request immediately. Note: we return here, which
			// means `next(c)` is NEVER called, so the actual route handler
			// is completely bypassed. This is the power of middleware.
			if key == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "API key required",
				})
			}

			// ── STEP 2: CAPTURE CURRENT TIME ─────────────────────────────────
			// We capture `now` ONCE and reuse it everywhere below.
			// This is important for consistency — if we called time.Now() in
			// multiple places, tiny differences in time could cause subtle bugs
			// (e.g., a window check and a retryAfter calculation might disagree
			// by a few nanoseconds). UTC avoids timezone-related surprises.
			now := time.Now().UTC()

			// ── STEP 3: LOOK UP THE CLIENT'S RECORD IN THE DATABASE ──────────
			// `r1` will hold the rate limit record for this API key.
			// This struct likely has: Key (string), Count (int), WindowStart (time.Time)
			var r1 RateLimit

			// db.First() searches for the FIRST row matching "key = ?".
			// We only care about `.Error` here — we don't use the return value
			// of First() directly because GORM populates `r1` by reference.
			err := db.First(&r1, "key=?", key).Error

			// IMPORTANT: GORM has two kinds of "errors" here —
			//   1. gorm.ErrRecordNotFound → no row for this key yet (NORMAL, first-time client)
			//   2. Any other error         → actual problem (DB down, bad query, etc.)
			// We only want to abort on REAL errors, not on "not found".
			if err != nil && err != gorm.ErrRecordNotFound {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"error": "Database error",
				})
			}

			// ── STEP 4: THE CORE WINDOW LOGIC ────────────────────────────────
			// Two situations mean the same thing: "start a fresh window for this client".
			//
			// Situation A: err == gorm.ErrRecordNotFound
			//   → This client has NEVER made a request before. No row exists in the DB.
			//
			// Situation B: now.Sub(r1.WindowStart) >= window
			//   → A record EXISTS but it's from an OLD window. For example, if window=60s
			//     and the record's WindowStart was 90 seconds ago, the window has expired.
			//     now.Sub(r1.WindowStart) calculates the elapsed time since the window started.
			//
			// Both cases are treated identically: create/overwrite the record with
			// Count=1 (this request counts as the first) and WindowStart=now.
			if err == gorm.ErrRecordNotFound || now.Sub(r1.WindowStart) >= window {
				r1 = RateLimit{
					Key:         key,
					Count:       1,   // This request is the 1st in the new window
					WindowStart: now, // The new window starts RIGHT NOW
				}
				// db.Save() is an UPSERT:
				//   → If r1 has a primary key (from a previous DB lookup), it UPDATEs.
				//   → If r1 is a fresh struct (no primary key), it INSERTs.
				// This elegantly handles both Situation A and B with one call.
				db.Save(&r1)

				// Window is fresh and count is 1 — request is allowed.
				// `next(c)` passes control to the next handler in the chain
				// (ultimately reaching the actual route handler).
				return next(c)
			}

			// ── STEP 5: CHECK IF LIMIT IS EXCEEDED ───────────────────────────
			// If we reach this point, we know:
			//   ✓ The client has a valid API key
			//   ✓ The client has an ACTIVE (non-expired) window in the DB
			// Now the only question is: have they used up their quota?
			if r1.Count >= limit {
				// Calculate how many seconds remain in the current window.
				// Formula: total window duration − time already elapsed in this window
				//
				// Example: window=60s, WindowStart was 40s ago
				//   → retryAfter = 60 - 40 = 20 seconds
				//
				// This tells the client: "wait 20 seconds, then try again."
				retryAfter := int(window.Seconds() - now.Sub(r1.WindowStart).Seconds())

				// "Retry-After" is a standard HTTP header (defined in RFC 9110).
				// Setting it is good API etiquette — well-behaved clients (and SDKs)
				// will read this and automatically back off for the right amount of time.
				c.Response().Header().Set("Retry-After", fmt.Sprint(retryAfter))

				// We return WITHOUT calling next(c) — the request is BLOCKED.
				// exceedStatus is typically 429 (Too Many Requests), but the caller
				// can configure it to anything (e.g., 503 for a different UX).
				return c.JSON(exceedStatus, echo.Map{
					"error":       "Rate limit exceeded",
					"retry_after": retryAfter,
					"status":      "exceeded",
				})
			}

			// ── STEP 6: ALLOW REQUEST AND INCREMENT COUNTER ──────────────────
			// If we reach here, the client is within their limit for the active window.
			// We increment their count to "use up" one slot, persist it to the DB,
			// then pass the request along.
			//
			// CORRELATION WITH STEP 5: The next time this client makes a request,
			// Step 3 will load this updated count from the DB, and Step 5 will
			// compare it against `limit` again. Once Count reaches `limit`, the
			// next request will be blocked — until Step 4 resets the window.
			r1.Count++
			db.Save(&r1) // Persist the incremented count — this is what "uses up" a slot.
			return next(c)
		}
	}
}

func main() {
	db := initDG()
	e := echo.New()

	// Create user endpoint
	e.POST("/signup", func(c echo.Context) error {
		type SignupRequest struct {
			Name string `json:"name"`
		}
		req := new(SignupRequest)
		if err := c.Bind(req); err != nil || strings.TrimSpace(req.Name) == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "name is required"})
		}
		apiKey := generateAPIKey()
		user := User{Name: req.Name, APIKey: apiKey}
		db.Create(&user)

		return c.JSON(http.StatusOK, echo.Map{"message": "user created successfully", "api_key": apiKey})
	})
	ratelimit := RateLimiter(db, 15, 15*time.Second, http.StatusForbidden)
	e.GET("/data", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "welcome to the data api",
			"time":    time.Now().Format(time.RFC3339),
		})
	}, ratelimit)

	log.Println("server is running on 8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
