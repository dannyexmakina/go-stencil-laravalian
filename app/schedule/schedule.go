package schedule

import (
	"log"

	"github.com/robfig/cron/v3"
)

// Register は app/Console/Kernel.php の schedule() 相当
func Register(c *cron.Cron) error {
	_, err := c.AddFunc("* * * * *", func() {
		log.Println("every minute task")
	})
	return err
}
