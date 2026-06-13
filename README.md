# Augment
Small app as an exercise for Augment

## Requirements
1. Create a New Augment Fund
 - Each Augment Fund must have at least:
    - A name (e.g., “Augment Fund II”).
    - A total number of units (e.g., 1,000).
 - Think about how you’ll store and identify each Augment Fund (e.g., with an
integer ID, a UUID, etc.).
2. Retrieve an Augment Fund’s Current Cap Table
 - The cap table shows who owns how many units, and when they acquired those
units.
 - Each line in the cap table must include at least:
    - The owner’s name.
    - The number of units they currently own.
    - The date they acquired (or last updated) that ownership.
 - You can decide how to structure and return this information (e.g., an array of
objects, a JSON response, a CLI printout, etc.).
3. Create a New Transfer
 - A transfer is an exchange of ownership between two people for a specified
number of units within an Augment Fund.
 - You must update the Augment Fund’s cap table accordingly (the “from” person
loses some units, the “to” person gains those units).
 - You should handle relevant validations and constraints (e.g., a person can’t
transfer more units than they own).
4. Show the History of All Transfers for an Augment Fund
 - Provide a way to list all past transfers that happened in a given Augment Fund,
in some sensible order (chronological or reverse-chronological).
 - Each transfer record should, at minimum, identify which Augment Fund it pertains
to, who transferred what to whom, and when it occurred.

## Conventions
- Seller means the one transferring units. Buyer is the receiver of the units. In a real transfer I assume a price would be attached, so I went with these terms rather than "giver" and "receiver" or "from" and "to".
- A seed investor is created in migrations that mimics a fund giving units to someone.
   - This was mostly a mechanism to get investors in the database with units to show the required logic for this exercise.
   - In practice I do not know how exactly units are distributed in the Augment ecosystem.

## Assumptions and choices
- Integer IDs
   - I prefer the incremental nature of IDs to UUIDs. For a simple apps like this, incremental IDs made finding objects and testing just a bit easier
- I did not do auto-creation of buyers and sellers
   - Sending units (shares) to someone who does not exist in the system seems like a good way to lose track of those units as they will not have an owner.
   - See [Improvements](#improvements) for more on how I consider handling this.
- MySQL or another DB is how this would be implemented properly for production, so I wanted to make an app in the fashion that I would if I were working on this to get it released.

## How To Use This Application

### Prerequisites
-Docker

### Starting The App
From the project root, simply run `docker compose up --build`. This will create the database, start the app, and run the migrations. Give it time on first startup - it can take up to 5 minutes.

If running on WSL, use localhost or the WSL IP if localhost does not resolve.

### Running Tests
From the project root start the app with `docker compose up --build`, then run `DB_HOST=localhost go test ./...` once the app and DB are ready.

### Endpoints And Their Return Formats
See the sections below for examples of the json objects returned for each endpoint.

#### Create Fund (POST)

Endpoint: `/fund/create`

Data:
- name: string
- units: int

EG: 
```bash
curl -i -X POST http://localhost:8080/fund/create \
-H "Content-Type: application/json" \
-d '{"name":"Test3","units":100}'
```

On Success Returns:
```
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sat, 13 Jun 2026 04:43:01 GMT
Content-Length: 111

{"id":2,"name":"Test3","units":100,"created_on":"2026-06-13T04:43:01Z","last_modified":"2026-06-13T04:43:01Z"}
```

On Failure:
```
HTTP/1.1 409 Conflict
Content-Type: application/json
Date: Sat, 13 Jun 2026 04:47:39 GMT
Content-Length: 48

{"error":"Fund with name Test3 already exists"}
```

#### Create Investor (POST)

Endpoint: `/investor/create`

Data:
- username: string
- email: string
- first_name: string
- last_name: string

EG: 
```bash
curl -i -X POST http://localhost:8080/investor/create \
-H "Content-Type: application/json" \
-d '{"username":"User1","email":"cboudreau@augment.com","first_name":"Cameron","last_name":"Boudreau"}'
```

On Success Returns:
```
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sat, 13 Jun 2026 04:53:50 GMT
Content-Length: 181

{"id":5,"username":"User1","email":"cboudreau@augment.com","first_name":"Cameron","last_name":"Boudreau","created_on":"2026-06-13T04:53:50Z","last_modified":"2026-06-13T04:53:50Z"}
```

On Failure:
```
HTTP/1.1 409 Conflict
Content-Type: application/json
Date: Sat, 13 Jun 2026 04:54:28 GMT
Content-Length: 69

{"error":"Investor with email cboudreau@augment.com already exists"}
```

#### Create Transfer (POST)

Endpoint: `/transfer/create`

Data:
- fund: int (ID of existing fund)
- buyer: int (ID of existing Investor)
- seller: int (ID of existing Investor)
- units: decimal (up to 4 places)

EG: 
```bash
curl -i -X POST http://localhost:8080/transfer/create \
-H "Content-Type: application/json" \
-d '{"fund":1,"buyer":3,"seller":2,"units":10}'
```

On Success Returns:
```
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sat, 13 Jun 2026 04:58:51 GMT
Content-Length: 88

{"id":4,"fund":1,"buyer":3,"seller":2,"units":"10","created_on":"2026-06-13T04:58:51Z"}
```

On Failure:
```
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Sat, 13 Jun 2026 04:59:45 GMT
Content-Length: 93

{"error":"Seller 2 does not have 500 units available in fund 1 to transfer. 390 available."}
```

#### Get Cap Table For Fund (GET)

Endpoint: `/cap/fund`

Query param:
- id: int (ID of existing fund)

EG: 
```bash
curl -i http://localhost:8080/cap/fund?id=1
```

On Success Returns:
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 13 Jun 2026 05:19:00 GMT
Content-Length: 263

[
   {"investor":"testFirst1 testLast1","units":"390","last_changed":"2026-06-13T04:58:51Z"},
   {"investor":"testFirst2 testLast2","units":"60","last_changed":"2026-06-13T04:58:51Z"},
   {"investor":"testFirst3 testLast3","units":"5","last_changed":"2026-06-13T04:36:47Z"}
]
```

On Failure:
```
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Sat, 13 Jun 2026 05:26:54 GMT
Content-Length: 43

{"error":"Fund with ID 5 does not exist."}
```

#### Get Cap Table History For Fund (GET)

Endpoint: `/cap/fund/history`

Query param:
- id: int (ID of existing fund)

EG: 
```bash
curl -i http://localhost:8080/cap/fund/history?id=1
```

On Success Returns:
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 13 Jun 2026 05:31:42 GMT
Content-Length: 461

[
   {"fund":1,"buyer":"testFirst1 testLast1","seller":"fund fund","units":"400","created_on":"2026-06-13T04:36:47Z"},
   {"fund":1,"buyer":"testFirst2 testLast2","seller":"fund fund","units":"50","created_on":"2026-06-13T04:36:47Z"},
   {"fund":1,"buyer":"testFirst3 testLast3","seller":"fund fund","units":"5","created_on":"2026-06-13T04:36:47Z"},
   {"fund":1,"buyer":"testFirst2 testLast2","seller":"testFirst1 testLast1","units":"10","created_on":"2026-06-13T04:58:51Z"}
]
```

On Failure:
```
HTTP/1.1 404 Not Found
Content-Type: application/json
Date: Sat, 13 Jun 2026 05:33:56 GMT
Content-Length: 43

{"error":"Fund with ID 5 does not exist."}
```

## Improvements
- Log Levels
- Email the intended recipient to invite them to Augment when someone tries to transfer units to an Investor that does not exist.
   - Temp hold on units for X time until the person creates account or the go back to the original owner?
   - Not sure on regulation for what's possible here
   - Probably need to confirm more info from the seller
- Need cost info on transfers
   - Pricing/most recent cost of transaction on the fund?
- I would move some of the validations done in the db file (checking if a fund with the same name, a username, or email exists already) and into validations like cap does.
- Add Swagger documentation
- More user friendly error messages and common return pattern with more information