# Golang nutrition, exercise, and weight tracking API for bodybuilding - WIP

## Needs work

### Currently tracks

#### Exercise
- Movements performed
- Weight movement performed at
- Sets of movement
- Reps per set of movement
- Date of movement

#### Diet & Nutrition
- Meals Eaten
- Dates of meals
- Macronutrients

**TODO:**

- Adaptation such that it is API focused only, instead of providing routes that map to Gin usable templates - Done

- Addition of diet tracking routes in addition to exercise/movement routes - Done

- Addition of weight tracking routes in `main.go`

- Finish implementation of Unit Tests in `main_test.go`

**Uses**

- `github.com/gin-contrib/cors`
- `github.com/gin-gonic/gin`
- `github.com/joho/godotenv`
- `gorm.io/driver/mysql`
- `gorm.io/gorm`

- Default testing framework (Needs Work)