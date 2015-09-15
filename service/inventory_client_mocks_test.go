package pezauth_test

const (
  TwoItemsSample = `{
  "status": "success",
  "data": [
    {
      "id": "55f077966696f40147000001",
      "sku": "2C.small",
      "tier": 0,
      "offeringType": "",
      "size": "",
      "attributes": {},
      "status": "",
      "lease_id": ""
    },
    {
      "id": "55f0779a6696f40147000002",
      "sku": "2C.small",
      "tier": 0,
      "offeringType": "",
      "size": "",
      "attributes": {},
      "status": "",
      "lease_id": ""
    }
  ]
}`

  LeasedItemsSample = `{
  "status": "success",
  "data": [
    {
      "id": "55f1dff033d183001d000001",
      "sku": "2C.small",
      "tier": 2,
      "offering_type": "C",
      "size": "small",
      "attributes": {},
      "status": "available",
      "lease_id": ""
    },
    {
      "id": "55f230977bed08001d000001",
      "sku": "2C.small",
      "tier": 2,
      "offering_type": "C",
      "size": "small",
      "attributes": {},
      "status": "leased",
      "lease_id": "55f22fe5b0cc8b001d000001"
    }
  ]
}`

 IndividualLease = `{
  "status": "success",
  "data": {
    "id": "55f22fe5b0cc8b001d000001",
    "inventory_item_id": "",
    "user": "dnem",
    "duration": "28 days",
    "start_date": "2015-01-02 15:04:05.000000000 +0000 UTC",
    "end_date": "2015-01-02 15:04:05.000000000 +0000 UTC",
    "status": "active",
    "attributes": {}
  }
}`

  NoInventoryDataSample = `{
    "status" : "fail",
    "message" : "Things went horribly wrong"
  }`

)
