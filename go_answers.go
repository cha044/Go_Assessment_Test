//1. Create a post api which will have rate limit on basis of each and unique ip.
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

type RateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func endpointHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, Message{
		Status: "Successful",
		Body:   "This is the endpoint for rate limiting",
	})
}

func main() {
	router := gin.Default()

	limiter := &RateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   rate.Limit(1),
		b:   5,
	}

	router.POST("/api/endpoint", rateLimiter(limiter), endpointHandler)

	err := router.Run(":8080")
	if err != nil {
		log.Println("There was an error listening to port: 8080", err)
	}
}

func (rl *RateLimiter) CreateLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]

	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.ips[ip] = limiter
	}

	return limiter
}

func rateLimiter(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := rl.CreateLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Requests exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}


//2. Demonstrate Command Design Pattern in Golang with Unit Tests
package main

import "fmt"

type Command interface {
	Execute()
	// Undo()
}

type Transaction interface {
	Browse()
	Buy()
}

type Listing struct {
	address  string
	isBought bool
}

func (pl *Listing) Browse() {
	fmt.Printf("browsing listing at %s\n", pl.address)
}

func (pl *Listing) Buy() {
	pl.isBought = true
	fmt.Printf("bought listing at %s\n", pl.address)
}

type BrowsePropertyCommand struct {
	property Transaction
}

func (c *BrowsePropertyCommand) Execute() {
	c.property.Browse()
}

func (c *BrowsePropertyCommand) Undo() {
	fmt.Println("undoing browse action")
}

type BuyPropertyCommand struct {
	property Transaction
}

func (c *BuyPropertyCommand) Execute() {
	c.property.Buy()
}

func (c *BuyPropertyCommand) Undo() {
	fmt.Println("undoing buy action")
}

func main() {
	property1 := &Listing{address: "567 Rosewood St."}
	property2 := &Listing{address: "809 Pine Ave."}

	browsePropCommand1 := &BrowsePropertyCommand{property: property1}
	buyPropCommand1 := &BuyPropertyCommand{property: property1}
	browsePropCommand2 := &BrowsePropertyCommand{property: property2}
	buyPropCommand2 := &BuyPropertyCommand{property: property2}

	browsePropCommand1.Execute()
	buyPropCommand1.Execute()
	browsePropCommand2.Execute()
	buyPropCommand2.Undo()
}


//unit test
func TestBrowse(t *testing.T) {
	property := &Listing{address: "567 Rosewood St."}
	browseCommand := &BrowsePropertyCommand{property: property}

	//checking browsed property
	if property.address != "567 Rosewood St." {
		t.Errorf("property address is not '567 Rosewood St.', got %s", property.address)
	} else {
		t.Logf("property address is '567 Rosewood St.', got %s", property.address)
	}
}

func TestBuy(t *testing.T) {
	property := &Listing{address: "809 Pine Ave."}
	buyCommand := &BuyPropertyCommand{property: property}


	//checking if property has been bought
	if !property.isBought {
		t.Errorf("property not bought, status stays available")
	} else {
		t.Logf("property successfully bought")
	}
}

func TestUndo(t *testing.T) {
	property := &Listing{address: "789 Oak St."}
	buyCommand := &BuyPropertyCommand{property: property}

	//checking if property purchase is undone
	if property.isBought {
		t.Errorf("property still bought")
	} else {
		t.Logf("buy property undone successfully")
	}
}



//4. Demonstrate Database Transaction of at least 10 different tables with proper error handling using
//a single api(use GIN).
func cusPropertyPurchase(c *gin.Context) {
    //necessary transaction data

    //user detail from user tbl
    userID := c.Param("userID")

    //cus details from cus tbl
    cusName := c.PostForm("cus_name")
    cusID := getCustomerID //assuming the function existed

    //preferences from preferences tbl
    noOfBedRm := c.PostForm("bedrm_no")
    noOfBathRm := c.PostForm("bathrm_no")

    //listing details from listing tbl
    listingID := 101

    //agent details from agent tbl
    agentID := 201 


    //simple cus acc creation
	sqlIns = sqlIns + " INSERT  INTO    	cus"
	sqlIns = sqlIns + " ("
	sqlIns = sqlIns + "                 cus_id,"
	sqlIns = sqlIns + "                 cus_name"
	sqlIns = sqlIns + "                 who_added" //userID is the inserted value
	sqlIns = sqlIns + "                 when_added"
	sqlIns = sqlIns + " )"
	sqlIns = sqlIns + " VALUES"
	sqlIns = sqlIns + " ("
	sqlIns = sqlIns + "                 ?,"
	sqlIns = sqlIns + "                 ?,"
	sqlIns = sqlIns + "                 ?,"
	sqlIns = sqlIns + "                 UNIX_TIMESTAMP()"
	sqlIns = sqlIns + " )"

	_, err = stmt.Exec(cusID, cusName)

	_, err := dbmap.Prepare(sqlIns)
		if err != nil {
			fmt.Println("insert1")
			fmt.Println(err)
			c.JSON(404, gin.H{"error": err})
			return 
		}
 

    //set cus property preferences
	sqlIns = sqlIns + " INSERT  INTO    	preferences"
	sqlIns = sqlIns + " ("
	sqlIns = sqlIns + "                 cusID,"
	sqlIns = sqlIns + "                 bedrm_no"
	sqlIns = sqlIns + "                 bathrm_no" 
	sqlIns = sqlIns + "                 bathrm_no"
	sqlIns = sqlIns + "                 who_added " //userID is the inserted value
	sqlIns = sqlIns + "                 when_added "
	sqlIns = sqlIns + " )"
	sqlIns = sqlIns + " VALUES"
	sqlIns = sqlIns + " ("
	sqlIns = sqlIns + "                 ?,"
	sqlIns = sqlIns + "                 ?,"
	sqlIns = sqlIns + "                 ?,"
	sqlIns = sqlIns + "                 ?,"
	sqlIns = sqlIns + "                 ?,"
	sqlIns = sqlIns + "                 UNIX_TIMESTAMP()"
	sqlIns = sqlIns + " )"

	_, err = stmt.Exec(cusID, noOfBedRm, noOfBathRm)

	_, err := dbmap.Prepare(sqlIns)
		if err != nil {
			fmt.Println("insert2")
			fmt.Println(err)
			c.JSON(404, gin.H{"error": err})
			return 
		}
    

    	//customer browse listing based on their preferences 
    	//... (retrieve and show listing based on cus preferences)

    	//customer viewing listing details
    	//... (retrieving listing details )

    	//proceeding to listing purchase
	//... (transaction details)
   
	//updating listing status to sold

	sqlUpd := ""
	sqlUpd = sqlUpd + " UPDATE  listing"
	sqlUpd = sqlUpd + " SET     status = 'sold' "
	sqlUpd = sqlUpd + "         who_updated = ?,"
	sqlUpd = sqlUpd + "         when_updated = UNIX_TIMESTAMP()"
	sqlUpd = sqlUpd + " WHERE   listing_id = ?"
	stmt, err := dbmap.Prepare(sqlUpd)

	_, err = stmt.Exec(status,
		userID,
		listingID)

	if err != nil {
		fmt.Println(err)
		c.JSON(404, gin.H{"error": err})
		return
	}


    c.JSON(http.StatusOK, gin.H{"message": "Customer property purchase completed successfully"})
}



//5. Write SQL Queries for Department of Top Three Salaries
USE goproject;

	SELECT
		e.id AS employeeId,
		e.name AS employee,
		e.salary AS salary,
		e.department_id,
		DENSE_RANK() OVER (PARTITION BY e.department_id ORDER BY e.salary DESC) 
	FROM 
		employee e;
        
	SELECT
		d.name AS department,
		r.employee,
		r.salary
	FROM
		department d
	JOIN
		TopSalaries r ON d.id = r.department_id
	WHERE
		r.SalaryRank <= 3;



//6. Return the minimum cost to reach the top of the floor.
//example given
cost := []int{10, 15, 20}


func minCost(val1, val2 int) int {
	if val1 < val2 {
		return val1
	}

	return val2
}


func minCostToTopFloor(cost []int) int {

	//initializing var for cost length 
	var costLength = len(cost)

	//initializing two previous steps with indices 0 and 1
	firstStep := cost[0]
	secondStep := cost[1]


	//current steps 'for loop' to get minimum cost
	for i := 2; i < costLength; i++ {
		currentCost := cost[i] + min(firstStep, secondStep)
		firstStep, secondStep = secondStep, currentCost
	}

	return minCost(firstStep, secondStep)	
}
