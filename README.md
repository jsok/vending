vending
=======

Vending machine in golang

Configuration
-------------

You should create a JSON file of the form:

```json
{
    "denominations": [1, 5, 20, 50],
    "slots": {
        "A1": {"item": {"name": "Coke Cola", "price": 180}, "inventory": 20},
        "etc": {}
    }
}
```

HTTP API
--------

###  POST `/api/purchase`

Purchase an item.

Example request:

```
curl http://localhost:5000/api/purchase \
 -d "coins[]=100" \
 -d "coins[]=50" \
 -d "coins[]=50" \
 -d "choice=A1"
```

Example response:

```json
{
    "id": "e13b7f80-6d58-11e4-9803-0800200c9a66",
    "item": "Coke Cola",
    "change": {"20": 1}
}
```

### GET `/api/items`

List all the items in the vending machine.

Example request:

```
curl http://localhost:5000/api/items
```

Example response:

```json
{
    "slots": {
        "A1": {"item": "Coke Cola", "price": 180, "available": true},
        "A2": {"item": "Water", "price": 120, "available": false}
    }
}
```

### POST `/api/items/<choice>`

Stock a new item.

Requires auth.

### PUT `/api/items/<choice>`

Refills an item.

Requires auth.

###  DELETE `/api/items/<choice>`

Sets an item as *OUT OF ORDER*.

Requires auth.
