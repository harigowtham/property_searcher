# property_searcher
Project with Golang: To upload properties and search among them

## To explain about the way I would approach this problem is, there are 2 categories:

1. The property uploader.
2. Property searcher 

## The property uploader:
The property is the one who wants to uploads a property, he will be asked for the various attributes like:

1. The location in terms of latitude and longitude coordinates.
2. The expected price of the property(in future we can have a max and min for this)
3. The number of bedrooms.
4. The number of bathrooms.
5. In the future, we can add support for additional info like the area of the property, parking information, backyard, etc.

Once he uploads the above values, it will upload them to the database. After uploading, we can then call a function to cross-check if there are any requirements that match the above property uploaded. If there are matches we can let the searcher know if there is a match for his search. More on how to do this is explained later.

## The property searcher:
The property searcher when wants to search for a property, he will be asked for a few sets of questions like below:

```
	1. The location around which he wants to search for the property.
		1. If the property is within 2 miles, we set the match as complete 30% for the property. If 5 then 20% and if 10 then 10%.
	2. The budget, get the min and max to be queried. 
		1. If one of it is unavailable, we take 75% of the value as min and value+25% as the max. 
		2. If the budget of the search is within the min and max then we take it as a full 30% match.
		3. If it is less, then 50% of the min will result in a 10% match and 80% of the min will result in a 20% match. (values can be changed)
		4. If it is more than the budget, then we will take budget+30% as the 20% match and budget+ 60% as the 10% match. 
		5. This way we calculate the match if the budget goes both the ways.

	3.  The bedroom and bathroom follow the same approach as they are of the same weights.
		1. If they fall under the min and max mentioned then we give it complete 20%.
		2. If the min or max is not mentioned, we take the value’s 75% and have it as min and the value+ 25% as the max 
		3. If the number of rooms mentioned is less, then we check if it's 50% more. If yes, we take it as a ten percent match.
		4. If it’s more then we check if its less than max+50%, then allocate it 10% match.
```

As of now, if the min and max are neglected, it has to be specified as 0. (Haven't handled the case it to automatically set it as 0) The above are the assumptions made as per the percentage contribution. The other assumption is that each degree change in the latitude or longitude coordinates represents a certain unit such as mile.

Ie: the difference between latitude 7.0000 and 7.0001 is one mile. 

The above assumption is done for the case of calculation.

## Approach:

### Storing all the property:

When a property is uploaded by the uploader, we add the property to the database.

In the program, I have maintained it in a list. As the SQL binding isn’t done.

About the id, it is supposed to be generated uniquely, but to simply things here, it has been got from the user. The object property as per the given requirements. The extended property has the pointer to the next property which is used to create the list. It has provisioning for adding new values such as the area of the property, version of the property, etc.

```
Note: While uploading, I haven’t checked the condition if the property is already uploaded.

This has to be coded and we have to fail the upload of duplicate properties.
```

### Searching for a requirement:

Once the above data is collected, we create a request object and create a heap to store the result. When we have a huge amount of data, going through the list and parsing it will be taking a huge toll on performance. So we have to offload as much as possible to the database to do the filtering. 

As the database needs values (at least with a range) that are not relative to the values in the DB. We can filter, the number of bedrooms, number of baths rooms and the budget using the DB as they are not relative. The location can’t be done with DB as the value are relative to the query  (the location coordinates in the query) 

Now about how the various range of weights can be calculated for the budget, bathrooms, and bedrooms. Here, as we have come up with a way to know how much of each is needed to contribute something. We will query for that range and mark it against the contribution by that.

```
For example, if a search is 10000 to 15000 as budget and 3 to 5 for both baths and beds. Then the budget, bed and bath range of a particular property between this will contribute to 30+20+20 which is 70%.

If the property has the budget and bed between the above range, but with only 2 baths, then the bath contribution is 10% instead of 20%. This is because , the value between 3 to 5 baths is 20%, the value between 1 to 3 and 5 to 7 are 10% as per the above explanation.

So we query for the bathroom values between this range will be having a contribution as 10% out of 20%. 

Similarly, the range of beds and budgets can change. And they will have their own contribution. This way the contribution has to be dynamically computed. 

To do this dynamic contribution, I have created an array for each bed, bath, and budget. 

The array for the budget will have the value [5000 8000 10000 15000 19500 24000]

The array of bed will have the value [1 3 5 7 0 0] 

and the array for baths will have the value [1 3 5 7 0 0]

Where the range between 5000 to 8000 & 19500 to 24000 will contribute 10% for the budget. 8000 to 10000 and 15000 to 19500 will contribute 20% and between 10000 to 15000 will contribute to 30%. Here the array is built in a way the array[0] is 50% of the min value. Array[1] is 80% of min and array[2] is the min, array[3] is the max.

Array[4] and array[5] is the max+30% and max+60% values.

So if we query the SQL with a range of array[0] to array[1] || array[4] to array[5] then we will get all the budget contribution that is 10%

The same can be done for 20% and 30%.

If we repeat the process for the bed and bath we can get the property that contributes 10% or 20% for each.

This way we can get contributions from 30% to 70% of the match from the budget, bed and bath contributions.
```

To do this, I have created an array and I loop through the array to create all the combinations. And these combinations are queried for and the resultant recorded are assumed to be parsed to a list and give for further processing.

Now, with this list where we have each property against its contribution. We calculate the location contribution. For this, we take the latitude and longitude coordinates of the search and the property and find the distance between them using the distance between the two points approach. As I mentioned earlier, I have assumed that each latitude/longitude coordinate change is assumed to be one mile (else we have to create a function to map accordingly). This way the distance between the two coordinates will give the distance between the property and the search location. If the distance is less than 2 miles then we will add 30% to the earlier computed contribution. If its within 2 to 5 miles then 20% and if its less than 10% and more than 5 then we will add 10% respectively. This way for each of the property we should be able to calculate the percentage.

So once we get the property and its corresponding match for the property with the search, we will insert it into a maxheap. And then remove the property from the list created from querying. Insertion is done only if it's more than a 40% match. The rest are discarded. This way all the property with the best match criteria at the top. The size of the heap will be equivalent to the size of the output expected.  

We can do the same in the opposite way for searching if a search request matches a property. 

### The time complexity:

The creation of arrays and rest of the code takes O(1) while the computation of the location contribution takes O(n). And then uploading it to the heapsort takes O(n logn). So the total time complexity is O(n*nlog n). 

### Space complexity:

The space complexity is O(n) for the program, and for the sort its O(1). So the total space complexity is O(n) at any point in time. 

#### Why heap sort over quick or merge sort?

The time complexity for all these sorts is pretty much the same. And also the space complexity of quick sort is O(n), while heapsort is O(1) but the program contributes O(n) so the total comes to O(n). The space complexities of all are the same for this implementation.

In spite of the time and space complexity being the same. The quicksort has added advantage as it is good in sorting a nearly sorted array faster than heap sort.

But the reasons behind going for heap sort are:

1. At any given point if we take the heap sorted list, it's going to give us the best result at that time. Which is not true for a quick sorted list.
2. When we do it multithreaded way, to reduce the execution time, we can write to the same heapsorted list from various threads(taking locks). This will consume lesser time than a quick sort or a quick sort that has to be sorted and merged for every thread.

If the above is not a necessary we can use quicksort.

To decrease the execution time, we can call a Goroutine to do the querying, computing the total contribution (adding the location contri to the already calculated rest of the contri) and then insert it to the common heap list by taking a lock. Once inserted, we remove the lock and then we delete the element from the SQL list. This should help us reduce the execution time. 

In the program, I wasn’t able to implement the search part, as the logic for using SQL and searching in a list are totally different. And the SQL implementation takes more time. 

About the goroutines, it’s the same. Without the list from each query, calling Goroutine didn't make sense. For the heap. I need to debug it a little further to make it work the way I want. 

The commented out lines need some work(creating heap and the code to put it into the heap). 

