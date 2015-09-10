package integrations_test

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

  NoInventoryDataSample = `{
    "status" : "fail",
    "message" : "Things went horribly wrong"
  }`

)
