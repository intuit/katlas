# K-Atlas APIs

### Introduction
The following documents describes K-Atlas API. It provides interfaces for data collector and user to CRUD resource including  metadata, entity and history.

### HTTP Header
|Header |Description|
|:--- |:---|
|Content-Type | application/json|

### HTTP Status Codes
|Status Code |Description|
|:-----------|:----------|
|200 - OK |The request has succeeded|
|201 - Created |The request has been fulfilled and resulted in a new resource being created|
|202 - Accepted |The request has been accepted for processing, but the processing has not been completed|
|204 - No Content |The server has fulfilled the request but does not need to return an entity-body, and might want to return updated meta information|
|400 - Bad Request |The request was malformed|
|404 - Not Found |Resource not found|
|500 - Server Error |The request could not be fulfilled due to an internal error in the server|
|503 - Service Unavailable |The request could not be fulfilled due to an error/unavailability of a downstream dependency|

### Metadata Service
CRUD API for metadata. The metadata describing the types of data.

**Create Metadata**:

Name | Description
:---|:---
`Request HTTP Method`| POST
`Request Path` | /v1/metadata
`Request Header Params`| Header above
`Request Body` | JSON format <br/>1. Use single json to create single metadata <br/>2. Array can be used to create multiple metadatas at a time
`Response` | Response code <br/> Success or error message

**Example**:
```
POST /v1/metadata
with body
[
  {
    "name":"application",
    "objtype":"metadata",
    "fields":[
      {
        "fieldname":"creationtime",
        "fieldtype":"string",
        "mandatory":true,
        "cardinality":"one"
      },
      {
        "fieldname":"objtype",
        "fieldtype":"string",
        "mandatory":true,
        "cardinality":"one"
      },
      {
        "fieldname":"name",
        "fieldtype":"string",
        "mandatory":true,
        "cardinality":"one"
      },
      {
        "fieldname":"resourceid",
        "fieldtype":"string",
        "mandatory":false,
        "cardinality":"one"
      },
      {
        "fieldname":"labels",
        "fieldtype":"json",
        "mandatory":false,
        "cardinality":"one"
      },
      {
        "fieldname":"resourceversion",
        "fieldtype":"string",
        "mandatory":true,
        "cardinality":"one"
      }
    ]
  }
]
```

**Get Metadata**:

Name | Description
:---|:---
`Request HTTP Method`| GET
`Request Path` | /v1/metadata/{type}
`Request Header Params`| Header above
`Request Body` | N/A
`Response` | Metadata with JSON format

**Example**:
```
GET /v1/metadata/application
return
{
  "uid":"0x1f019",
  "name":"application",
  "objtype":"metadata",
  "fields":[
    {
      "fieldname":"resourceid",
      "fieldtype":"string",
      "mandatory":false,
      "cardinality":"One"
    },
    {
      "fieldname":"labels",
      "fieldtype":"json",
      "mandatory":false,
      "cardinality":"One"
    },
    {
      "fieldname":"resourceversion",
      "fieldtype":"string",
      "mandatory":true,
      "cardinality":"One"
    },
    {
      "fieldname":"creationtime",
      "fieldType":"string",
      "mandatory":true,
      "cardinality":"One"
    },
    {
      "fieldname":"objtype",
      "fieldType":"string",
      "mandatory":true,
      "cardinality":"One"
    },
    {
      "fieldname":"name",
      "fieldType":"string",
      "mandatory":true,
      "cardinality":"One"
    }
  ]
}
```

**Update Metadata**:

Name | Description
:---|:---
`Request HTTP Method`| POST
`Request Path` | /v1/metadata/{type}
`Request Header Params`| Header above
`Request Body` | Data with JSON format for update
`Response` | Response code <br/> Metadata name and unified ID. Or error message if any

**Example**:
```
POST /v1/metadata/application
with body
{
  "fields":[
    {
      "fieldname":"belongsTo",
      "fieldtype":"relationship",
      "mandatory":true,
      "refdatatype":"namespace"
    }
  ]
}
return
{
  "status": 200,
  "objects": [{
    "objtype": "metadata",
    "uid" : "0x123ba0"
  }]
}
```

**Delete Metadata**:

Name | Description
:---|:---
`Request HTTP Method`| DELETE
`Request Path` | /v1/metadata/{type}
`Request Header Params`| Header above
`Request Body` | N/A
`Response` | Response code <br/> Metadata name and unified ID. Or error message if any

**Example**:
```
DELETE /v1/metadata/application
return
{
  "status":200,
  "objects":[{
    "objtype":"metadata",
    "uid":"0x467ba0"
  }]
}
```

**Upsert Schema for Metadata**:

Name | Description
:---|:---
`Request HTTP Method`| POST
`Request Path` | /v1/metadata/schema
`Request Header Params`| Header above
`Request Body` | JSON format <br/>1. Use single json to upsert single schema <br/>2. Array can be used to upsert multiple schema at a time
`Response` | Response code <br/> Success or error message

**Example**:
```
POST /v1/schema
with body
[
  {
    "predicate":"resourceid",
    "type":"string",
    "index":true,
    "upsert":true,
    "tokenizer":[
      "hash",
      "trigram"
    ]
  },
  {
    "predicate":"label",
    "type":"string",
    "index":true,
    "upsert":true,
    "tokenizer":[
      "term",
      "trigram"
    ]
  }
]
```

**Delete Schema**:

Name | Description
:---|:---
`Request HTTP Method`| DELETE
`Request Path` | /v1/schema/{name}
`Request Header Params`| Header above
`Request Body` | N/A
`Response` | Response code <br/> message to show success or fail with schema name

**Example**:
```
DELETE /v1/schema/resourceid
return
{
  "status":200,
  "message": "schema resourceid drop successfully"
}
```

### Entity Service
CRUD API for entity

**Create Entity**:
Create an entity based on given metadata

Name | Description
:---|:---
`Request HTTP Method`| POST
`Request Path` | /v1/entity/{metadata}
`Request Header Params`| Header above
`Request Body` | JSON format data
`Response` | Response code <br/> Entity type and unified ID. Or error message if any

**Example**:
```
POST /v1/entity/pod
with body
{
  "description":"describe this object",
  "resourceid":"unique_id_of_pod",
  "name":"pod01",
  "resourceversion":"6365014",
  "starttime":"2018-09-01T10:01:03Z",
  "status":"Running",
  "ip":"172.20.32.128",
  "namespace":"ns"
}
return
{
  "status":200,
  "objects":[{
    "objtype":"pod",
    "uid":"0x467ba0"
  }]
}
```

**Get Entity**:
Get an entity based on given metadata and uid

Name | Description
:---|:---
`Request HTTP Method`| GET
`Request Path` | /v1/entity/{metadata}/{uid}
`Request Header Params`| Header above
`Request Body` | N/A
`Response` | Response code <br/> Entity or error message if any

**Example**:
```
GET /v1/entity/pod/0x467ba0
return
{
  "status":"200",
  "objects":[{
    "uid":"0x467ba0",
    "objtype":"pod",
    "description":"describe this object",
    "resourceid":"unique_id_of_pod",
    "name":"pod01",
    "resourceversion":"6365014",
    "starttime":"2018-09-01T10:01:03Z",
    "status":"Running",
    "ip":"172.20.32.128",
    "namespace":{
      "uid":"0x56291a"
    }
  }]
}
```
**Update Entity**:
Update an entity based on given metadata

Name | Description
:---|:---
`Request HTTP Method`| POST
`Request Path` | /v1/entity/{metadata}/{uid}
`Request Header Params`| Header above
`Request Body` | JSON format data
`Response` | Response code <br/> Entity type and unified ID. Or error message if any

**Example**:
```
POST /v1/entity/pod/0x467ba0
with body
{
  "resourceversion":"6365015",
  "status":"Failed"
}
return
{
  "status":200,
  "objects":[{
    "objtype":"pod",
    "uid":"0x467ba0"
  }]
}
```

**Delete Entity**:
Delete an entity based on given metadata and resourceid

Name | Description
:---|:---
`Request HTTP Method`| DELETE
`Request Path` | /v1/entity/{metadata}/{resourceid}
`Request Header Params`| Header above
`Request Body` | N/A
`Response` | Response code <br/> Entity type and unified ID. Or error message if any

**Example**:
```
DELETE /v1/entity/pod/resourceid01
return
{
  "status":200,
  "objects":[{
    "objtype":"pod",
    "uid":"0x467ba0"
  }]
}
```

### Query Service
Query to get resources

**Keyword Search**:

Name | Description
:---|:---
`Request HTTP Method`| GET
`Request Path` | /v1/query
`Request Header Params`| Header above
`Request Query Params` | The keyword to be matched as substring
`Request Body` | N/A
`Response` | Response code <br/> Entity type and unified ID. Or error message if any

**Example**:
```
GET /v1/query?keyword=webapp
return
{
  "status":200,
  "objects":[
    {
      "uid":"0x467ba0",
      "objtype":"pod",
      "description":"describe this object",
      "resourceid":"webapp",
      "name":"webapp",
      "resourceversion":"6365014",
      "starttime":"2018-09-01T10:01:03Z",
      "status":"Running",
      "ip":"172.20.32.128"
    },
    {
      "uid":"0x467ba2",
      "objtype":"namespace",
      "resourceid":"webapp-ns",
      "name":"webapp-ns",
      "resourceversion":"6365014"
    }
  ]
}
```

**Key/value Query**:

Name | Description
:---|:---
`Request HTTP Method`| GET
`Request Path` | /v1/query
`Request Header Params`| Header above
`Request Query Params` | The key=value pairs to be matched, using `print` to specify which fields to be returned, seprated by comma
`Request Body` | N/A
`Response` | Response code <br/> Entity type and unified ID with specified fields (return all fields by default). Or error message if any

**Example**:
```
GET /v1/query?name=webapp&objtype=pod&print=name,resourceid
return
{
  "status":200,
  "objects":[{
    "uid":"0x467ba0",
    "objtype":"pod",
    "resourceid":"webapp",
    "name":"webapp"
  }]
}
```

**QSL query**:
refer https://github.com/intuit/katlas/blob/master/docs/qsl-api.md
