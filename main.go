//time complexity is O(n(n log n)) ie o(n^2 nlogn)
//space complexity is O(l)

package main

import (
//  	"container/heap"
	"fmt"
//	"math"
)

type Property struct {
  Id string // have to be generated, but as of now, assumed that the user inputs this
  Latitude float32
  Longitude float32
  Price int32
  Numbeds int32
  Numbaths int32
}

type Search struct {
  Id string // have to be generated, but as of now, assumed that the user inputs this
  Latitude float32
  Longitude float32
  MinPrice int32
  MinBeds int32
  MinBaths int32
  MaxPrice int32
  MaxBeds int32
  MaxBaths int32
}

// More Can be added if needed
type extendedProperty struct {
  Property
  version int
  Percentage int
  index int
  next *extendedProperty
}

type PriorityQueue []*extendedProperty

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest percentage so we use greater than here.
	return pq[i].Percentage > pq[j].Percentage
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*extendedProperty)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type Plist struct {
	Id       string
	head       *extendedProperty
}

func createPlist(Id string) *Plist {
	return &Plist{
		Id: Id,
	}
}

func (p *Plist) listProperty() error {
	currentNode := p.head
	if currentNode == nil {
		fmt.Println("Property list is empty.")
		return nil
	}
	fmt.Printf("%+v\n", *currentNode)
	for currentNode.next != nil {
		currentNode = currentNode.next
		fmt.Printf("%+v\n", *currentNode)
	}

	return nil
}

func calroomrange(max int32, min int32, array []int32) {
	if(max == 0) {
		max = int32(float32(min)*.25) + min
		min = int32(float32(min)*.75)
	} else if (min == 0) {
		min = int32(float32(max)*.25) + max
		max = int32(float32(max)*.75)
	} 
	array[0] = int32(float32(min)*0.5)
	array[1] = min
	array[2] = max
	array[3] = int32(float32(max)*0.5) + max
}

func calbudgetrange(max int32, min int32, array []int32) {
	if(max == 0) {
		max = int32(float32(min)*.25) + min
		min = int32(float32(min)*.75)
	} else if (min == 0) {
		min = int32(float32(max)*.25) + max
		max = int32(float32(max)*.75)
	} 
	array[0] = int32(float32(min)*0.5)
	array[1] = int32(float32(min)*0.8)
	array[2] = min
	array[3] = max
	array[4] = int32(float32(max)*0.3) + max
	array[5] = int32(float32(max)*0.6) + max

}

func (r *Search) createarrs(arrayofarray [][]int32) {
	for i := 0; i < 3; i++ {
		arrayofarray[i] = make([]int32, 6)
	}
	calroomrange(r.MaxBeds, r.MinBeds, arrayofarray[0])
	calroomrange(r.MaxBaths, r.MinBaths, arrayofarray[1])
	calbudgetrange(r.MaxPrice, r.MinPrice, arrayofarray[2])
}

func contricompute(i int, j int, k int) (contri int) {
	contri = (i+1)*10 + (j+1)*10 + (k+1)*10
	return 
}

/*
func (s *Search) locationcontri(lat, long) lcontri {
	distance := math.sqrt((delta()**2) + (delta()**2))
	if distance < 2 
		return 30
	else if distance < 5
		return 20
	else if distance < 10
		return 10
	else
		return -1
}

func (s *Search) callroutine(i,j,k int, arrays [][]int32, heaplist *PriorityQueue, p *Plist) {
	//create a temp list
	temp list = sqlquery(ranges of i, j and k using array in the Plist)
	//assume that the returned records are put in a list.
	
	//loop the list and start insert it into the heaplist based on the location contribution.
	while (temp list is not empty) {
		//read from templist
		//calculate contri without location
		restcontri := contricompute(i,j,k)
		//calculatelocation contri
		lcontri := s.locationcontri(cur_lat,cur_long)
		if (lcontri != -1) {
			totalcontri = lcontri + restcontri
			//upload the totalcontri as the priority
			take lock on heapsortedlist
			write to heapsortedlist based on the total contributed percentage.
			release lock
		}
		delete emtries in the list
	}
}
*/

func (s *Search) loopandquery(arrays [][]int32, heaplist *PriorityQueue, p *Plist) {
	for i := 2; i >= 0; i-- {
		for j := 1; j>= 0; j-- {
			for k := 1; k>= 0; k-- {
				//s.callroutine(i,j,k, arrays, PriorityQueue, p)
			}
		}
	}
}

//assume we give the value as zero for the ones which user doesnt enter.
func (p *Plist) searchProperty() error {
	var Id string
	var Lat float32
	var Long float32
	var MaxPrice int32
	var MaxBeds int32
	var MaxBaths int32
	var MinPrice int32
	var MinBeds int32
	var MinBaths int32
	arrays := make([][]int32, 3)
	//heapsortedlist := make(PriorityQueue, 100)
	fmt.Scanln(&Id)
	fmt.Scanln(&Lat)
	fmt.Scanln(&Long)
	fmt.Scanln(&MinPrice)
	fmt.Scanln(&MaxPrice)
	fmt.Scanln(&MinBeds)
	fmt.Scanln(&MaxBeds)
	fmt.Scanln(&MinBaths)
	fmt.Scanln(&MaxBaths)
	request := &Search{
		Id:   Id,
		Latitude: Lat,
		Longitude:  Long,
		MinPrice: MinPrice,
		MinBeds: MinBeds,
		MinBaths: MinBaths,
		MaxPrice: MaxPrice,
		MaxBeds: MaxBeds,
		MaxBaths: MaxBaths,
	}
	//heap.Init(&heapsortedlist)
	request.createarrs(arrays)
	fmt.Println(arrays)
	//request.loopandquery(arrays, &heapsortedlist, p)
	return nil
}

func (l *Plist) addProperty() error {
	var Id string
	var Lat float32
	var Long float32
	var Price int32
	var Beds int32
	var Baths int32
	fmt.Scanln(&Id)
	fmt.Scanln(&Lat)
	fmt.Scanln(&Long)
	fmt.Scanln(&Price)
	fmt.Scanln(&Beds)
	fmt.Scanln(&Baths)
	p := &Property{
		Id:   Id,
		Latitude: Lat,
		Longitude:  Long,
		Price: Price,
		Numbeds: Beds,
		Numbaths: Baths,
	}
	ep := &extendedProperty{
		Property: *p,
	}
	if l.head == nil {
		l.head = ep
	} else {
		currentNode := l.head
		for currentNode.next != nil {
			currentNode = currentNode.next
		}
		currentNode.next = ep
	}
	return nil
}

func main() {
	fmt.Println("Creating list")
	Plist := createPlist("firstlist")

	fmt.Println("Hello, Pressing as mentioned")
	for {
		fmt.Printf("1 for adding property\n")
		fmt.Printf("2 for adding requirement\n")
		fmt.Printf("3 for searching Property by requirement\n")
		fmt.Printf("4 for searching requirements by property\n")
		fmt.Printf("5 for listing property\n")
		fmt.Printf("6 to exit\n")
		var input int
		fmt.Scanln(&input)
		fmt.Println("Input: ", input)
		switch {
		case input == 1:
			fmt.Println("Adding Property")
			Plist.addProperty()
		case input == 2:
			fmt.Println("Adding requirement")
			Plist.addProperty()
		case input == 3:
			fmt.Println("Searching Property")
			Plist.searchProperty()
		case input == 4:
			fmt.Println("Searching requirement")
			Plist.listProperty()
		case input == 5:
			fmt.Println("Searching requirement")
			Plist.listProperty()
		case input == 6:
			fmt.Println("Exiting")
			return
		default:
			fmt.Println("Invalid")
		}
	}
}
