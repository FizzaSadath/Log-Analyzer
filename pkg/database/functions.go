package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (e Entry) String() string {
	return fmt.Sprintf("%s : %s : %s : %s : %s : %s",
		e.TimeStamp.Format("2006-01-02 15:04:05"),
		e.Level.Level,         // LogLevel.Level (string)
		e.Component.Component, // LogComponent.Component (string)
		e.Host.Host,           // LogHost.Host (string)
		e.RequestId,
		e.Message,
	)
}

func CreateDB(dbUrl string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error
			Colorful:                  true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("couldn't open database %v", err)
	}
	return db, nil

}
func InitDB(db *gorm.DB) error {
	db.AutoMigrate(&LogLevel{}, &LogComponent{}, &LogHost{}, &Entry{})
	db.Create(&[]LogLevel{
		{Level: "INFO"},
		{Level: "WARN"},
		{Level: "ERROR"},
		{Level: "DEBUG"},
	})
	db.Create(&[]LogComponent{
		{Component: "api-server"},
		{Component: "database"},
		{Component: "cache"},
		{Component: "worker"},
		{Component: "auth"},
	})
	db.Create(&[]LogHost{
		{Host: "web01"},
		{Host: "web02"},
		{Host: "cache01"},
		{Host: "worker01"},
		{Host: "db01"},
	})
	return nil
}

func AddDB(db *gorm.DB, e Entry) error {
	ctx := context.Background()
	err := gorm.G[Entry](db).Create(ctx, &e)
	if err != nil {
		return err
	}
	return nil
}
func parseQuery(parts []string) ([]queryComponent, error) {
	var ret []queryComponent

	pattern := `^(?P<key>[^\s=!<>]+)\s*(?P<operator>=|!=|>=|<=|>|<)\s*(?P<value>.+)$`
	r := regexp.MustCompile(pattern)

	for _, part := range parts {
		part = strings.TrimSpace(part)

		matches := r.FindStringSubmatch(part)
		if matches == nil {
			return nil, fmt.Errorf("invalid condition: %s", part)
		}

		// Allow INFO|ERROR
		rawValue := matches[r.SubexpIndex("value")]
		rawValue = strings.ReplaceAll(rawValue, "|", ",")

		vals := strings.Split(rawValue, ",")

		ret = append(ret, queryComponent{
			key:      matches[r.SubexpIndex("key")],
			operator: matches[r.SubexpIndex("operator")],
			value:    vals,
		})
	}

	return ret, nil
}

func QueryDB(db *gorm.DB, query []string) ([]Entry, error) {
	var ret []Entry

	// Parse the query string
	parsed, err := parseQuery(query)
	if err != nil {
		return nil, err
	}

	fmt.Println("Parsed conditions:", parsed)

	q := db

	for _, c := range parsed {

		key := strings.ToLower(c.key)

		switch key {

		case "level":
			// Convert values like INFO to levelID
			var ids []uint
			for _, v := range c.value {
				var lvl LogLevel
				if err := db.First(&lvl, "level = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown level '%s'", v)
				}
				ids = append(ids, lvl.ID)
			}
			c.key = "level_id"
			c.value = toStringSlice(ids)

		case "component":
			var ids []uint
			for _, v := range c.value {
				var comp LogComponent
				if err := db.First(&comp, "component = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown component '%s'", v)
				}
				ids = append(ids, comp.ID)
			}
			c.key = "component_id"
			c.value = toStringSlice(ids)

		case "host":
			var ids []uint
			for _, v := range c.value {
				var h LogHost
				if err := db.First(&h, "host = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown host '%s'", v)
				}
				ids = append(ids, h.ID)
			}
			c.key = "host_id"
			c.value = toStringSlice(ids)
		}

		// Apply WHERE condition
		if len(c.value) == 1 {
			q = q.Where(fmt.Sprintf("%s %s ?", c.key, c.operator), c.value[0])
		} else {
			if c.operator == "!=" {
				q = q.Where(fmt.Sprintf("%s NOT IN ?", c.key), c.value)
			} else {
				q = q.Where(fmt.Sprintf("%s IN ?", c.key), c.value)
			}
		}
	}
	q = q.
		Preload("Level").
		Preload("Component").
		Preload("Host")
	if err := q.Find(&ret).Error; err != nil {
		return nil, err
	}

	return ret, nil
}

// convert slice of foreign keys to string
func toStringSlice(nums []uint) []string {
	s := make([]string, len(nums))
	for i, n := range nums {
		s[i] = fmt.Sprint(n)
	}
	return s
}

// for web
func SplitUserFilter(input string) []string {
	var parts []string
	current := ""
	tokens := strings.Fields(input)

	for _, tok := range tokens {
		// If token contains an operator, then new condition
		if strings.Contains(tok, "=") ||
			strings.Contains(tok, ">=") ||
			strings.Contains(tok, "<=") ||
			strings.Contains(tok, ">") ||
			strings.Contains(tok, "<") {

			// Save previous condition
			if current != "" {
				parts = append(parts, current)
			}
			current = tok
		} else {
			// continuation (timestamps)
			current += " " + tok
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
func GetAllLogs(db *gorm.DB) ([]Entry, error) {
	var result []Entry
	err := db.Preload("Level").
		Preload("Component").
		Preload("Host").
		Find(&result).Error
	return result, err
}
func FilterLogsWeb(db *gorm.DB, levels, components, hosts []string, requestID, timestampCond string) ([]Entry, error) {
	var queries []string

	// Levels
	if len(levels) > 0 && len(levels) < 4 { // skip if all selected
		queries = append(queries, "level="+strings.Join(levels, "|"))
	}

	// Components
	if len(components) > 0 && len(components) < 5 {
		queries = append(queries, "component="+strings.Join(components, "|"))
	}

	// Hosts
	if len(hosts) > 0 && len(hosts) < 5 {
		queries = append(queries, "host="+strings.Join(hosts, "|"))
	}

	// RequestID filter (direct equality)
	if strings.TrimSpace(requestID) != "" {
		queries = append(queries, "request_id="+requestID)
	}

	// Timestamp filter (operator + value)
	if strings.TrimSpace(timestampCond) != "" {
		queries = append(queries, "time_stamp "+timestampCond)
	}

	// If no filters, return all logs
	if len(queries) == 0 {
		return GetAllLogs(db)
	}

	// Call existing QueryDB which takes []string
	return QueryDB(db, queries)
}
