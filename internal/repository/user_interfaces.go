package repository
import "github.com/Shubham23jha/go-gin-postgres-clean/internal/models"

type UserRepository interface {
	Create(user *models.User) error

	
	FindAll() ([]models.User, error)
}
