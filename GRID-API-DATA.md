4. Read the documentation
The Documentation section covers all GRID APIs. As an Open Access user, the relevant parts are:

Static Data: everything about the Central Data API
In-Game Data → Series State: everything about the Series State API
API Reference: the full schema for both
Sections on Series Events, File Download, Stats Feed, Widgets and Video describe products that are not part of Open Access. Remember to swap api.grid.gg for api-op.grid.gg in any examples you copy.



Accessing Riot Data through GRID APIs
Riot Open Access gives you official League of Legends and VALORANT esports data through two GraphQL APIs: Central Data and Series State. It is free for non-commercial use, such as research, prototyping and evaluation.

What you get
Riot Open Access
Games
League of Legends, VALORANT
APIs
Central Data API, Series State API, Stats API
Data
Schedules, tournaments, teams, players, and pre- and post-match series data
Latency
5-minute delay on live series
Rate limit
120 requests per series per hour
Need live feeds, higher rate limits or any of the above? Reach out through Upgrade to Commercial Access on the Welcome page.

Endpoints
Open Access uses its own host, api-op.grid.gg. Keys issued for Open Access do not work on the standard api.grid.gg endpoints shown elsewhere in the documentation.

API
Use it for
Endpoint
Central Data
Finding titles, tournaments, series, teams and players
https://api-op.grid.gg/central-data/graphql
Series State
The state of a series: live snapshot (5-minute delay) or final result
https://api-op.grid.gg/live-data-feed/series-state/graphql
Authentication
To read more about how to authenticate your requests, head over here

Your first queries
Everything starts with a series ID from Central Data, which you then pass to Series State.

1. Find the title IDs for League of Legends (ID: 3) and VALORANT (ID: 6) (Central Data).

Try it out

query Titles {
  titles { id name }
}
2. List tournaments for a title (Central Data). Replace the ID with one from step 1.

3. List series in a time window (Central Data) or use the Dashboard to copy series ID from a match.

4. Get the series state (Series State). A finished series returns its final state; a live series returns a snapshot 5 minutes behind.

Try it out

query SeriesState {
  seriesState(id: "YOUR_SERIES_ID") {
    started
    finished
    teams { id name won }
    games {
      teams { id players { id name kills deaths } }
    }
  }
}
To see every available field, open the schema panel in the GQL Playground, or run an introspection query in Postman, Insomnia or Bruno.

Usage terms
When you registered, you agreed that:

Your use is non-commercial. If your project becomes commercial, it must move to the commercial tier before any commercial use of the data.
The data is not used for betting, gambling or wagering, and is not supplied to anyone for those purposes.
You accept the Open Access Terms of Service.
Troubleshooting
Error
Likely cause
Fix
Requester unauthorized to make query (UNAUTHENTICATED)
Wrong host, or missing or wrong key
Use the api-op.grid.gg endpoints and check the x-api-key header
Same error on a valid key
Querying an API or title outside Open Access
Stick to Central Data and Series State for LoL and VALORANT
Rate limit error
Over 120 requests per series in an hour
Cache responses and poll a live series less often


Overview
The GRID APIs are a mixture of REST, GraphQL and WebSocket based APIs.

Within this reference area you'll find individual documentation for each API we offer to you along with which protocol to use for this API.

Feed Overview
Central Data Feed
Query API: a GraphQL API exposing information about titles, tournaments, teams and series (scheduling).
Live Data Feed
Series State API: a GraphQL API exposing the current state of a live series (updated as actions happen during an ongoing game) and the final state of previously played series.
Series Events API: a Websocket API that enables clients to listen for events that are pushed to them as they occur during a live series.
File Download API: a REST API that allows the download of the complete gameplay data of a previously played series.
Access
Access to the GRID APIs is different for every consumer, which can vary on both per API and also for the title, tournament and individual series being accessed.

Bulk Access
Bulk requests are not supported, so instead a single request for every object is required.



Authentication
GRID's APIs use API Keys to authenticate requests, you can manage API Keys by contacting our support team.

These API Keys carry a lot of privilege so try ensure you keep these secret and secure. Don't share these with others, or upload to public spaces such as GitHub.

GRID's APIs utilise two different protocols HTTP and WebSocket, due to limitations our WebSocket API's

REST/GraphQL
As both of these standards utilise HTTP underneath we can use the x-api-key header with the value being that of our API Key provided to you.


var headers = {
	'x-api-key': 'XXX'
}
WebSocket
To access our WebSocket based connections you will need to use a QueryParam appended to the end of the wss url you are using to connect to the API.

wss://api.grid.gg...?key=XXX

All of our requests must be used using SSL therefor please ensure your using HTTPS or WSS endpoints when querying for data.

GraphQL
Authentication
Users of all the GRID APIs must identify themselves before consuming the APIs. The authentication process requires passing a key to the API.

When consuming a GraphQL API, the API key can be passed by setting the x-api-key header in the HTTP request. The way in which this is achieved will vary depending on the specific tool being utilized.

More information about the GraphQL authorization mechanism can be found in the official documentation located at this page.

Pagination
GRID APIs can be used to retrieve all the series for a given time period. While querying for short time periods may yield a small list of series in the response, querying for larger periods, such as a year, would result in larger lists of series. In this latter case, pagination is used to split the results into multiple pages.

GRID APIs maximum page length is 50 results. Thus, a single query will return no more than 50 series. The caller can adjust the page length, and it is currently set to 10 results by default.

Let us consider the following query, which can be found in the GRID Playground. The query requests all the series that started between 03/11/2021 and 04/11/2021, and it will yield the first 10 series from that period.


query GetAllSeriesInNext24Hours {
  allSeries(
    filter:{
      startTimeScheduled:{
        gte: "2021-11-03T00:00:00+01:00"
        lte: "2021-11-04T00:00:00+01:00"
      }
    }
    orderBy: StartTimeScheduled
  ) {
    totalCount,
    pageInfo{
      hasPreviousPage
      hasNextPage
      startCursor
      endCursor
    }
    edges{
      cursor
      node{
        id
        startTimeScheduled
      }
    }
  }
}
The resulting JSON response can be similar to the following:


{
  "data": {
    "allSeries": {
      "totalCount": 32,
      "pageInfo": {
        "hasPreviousPage": false,
        "hasNextPage": true,
        "startCursor": "JAMWBQMjHhoSJBQfEhMCGxITCwtaCwtERUdORUcLC1oLC0VHRUZaRkZaR0QjR05NR0dNR0ct",
        "endCursor": "JAMWBQMjHhoSJBQfEhMCGxITCwtaCwtERUZHRUULC1oLC0VHRUZaRkZaR0QjRkVNR0dNR0ct"
      },
      "edges": [...]
    }
  }
}

To understand this response, it is necessary to examine the following three fields values:

totalCount refers to the number of results matching the given filter.
hasNextPage indicates whether further querying is required to obtain all the results.
endCursor works like a bookmark. It indicates where further queries should start in order to get the next set of results.
Based on the above information, the following query can be executed in order to obtain the next set of results:


query GetAllSeriesInNext24Hours {
  allSeries(
    after: "JAMWBQMjHhoSJBQfEhMCGxITCwtaCwtERUZHRUULC1oLC0VHRUZaRkZaR0QjRkVNR0dNR0ct"
    filter:{
      startTimeScheduled:{
        gte: "2021-11-03T00:00:00+01:00"
        lte: "2021-11-04T00:00:00+01:00"
      }
    }
    orderBy: StartTimeScheduled
  ) {
    totalCount,
    pageInfo{
      hasPreviousPage
      hasNextPage
      startCursor
      endCursor
    }
    edges{
      cursor
      node{
        id
        startTimeScheduled
      }
    }
  }
}
Notably, the after field is set with the value of the endCursor field from the result of the previous query. This ensures that the first item of the returned result is the one following the last of the first query. In order to fetch the entire result set, it is sufficient to keep executing queries, passing the value of the endCursor of the previous result, until the hasNextPage field is set to false.

Instead, in order to get all the results in one query - but always up to 50 - it would be sufficient to specify the first field in the query:


query GetAllSeriesInNext24Hours {
  allSeries(
    first: 50
    filter:{
      startTimeScheduled:{
        gte: "2021-11-03T09:05:39+01:00"
        lte: "2021-11-04T09:05:39+01:00"
      }
    }
    orderBy: StartTimeScheduled
  ) {
    totalCount,
    pageInfo{
      hasPreviousPage
      hasNextPage
      startCursor
      endCursor
    }
    edges{
      cursor
      node{
        id
        startTimeScheduled
      }
    }
  }
}
The first field indicates the number of items we request. This can be any number, from 1 to 50. The first field has a counterpart: last, which allows querying the results the other way around; instead of returning the first X results that match the given filter, the result will include the last X results.

More information on the GraphQL pagination can be found in the official documentation here.

Rate limit
The rate limit restricts the number of queries a user can execute in a time period.

By default, GRID allows 40 requests per minute for Closed Platform and 20 requests per minute for Open Platform. Any more attempts will be rejected.

Errors
Error type
Reason
What to do
UNAUTHENTICATED
The API key provided can’t be found.
Check the right API key is passed. If the API key is correct, users may contact the GRID Support team for assistance.
UNAUTHORIZED
The API key provided has no access to the product/content.
Contact the GRID Support team to get the proper permissions.
Error type
Error detail
Reason
What to do
BAD_REQUEST
INVALID_ARGUMENT
Can either provide first/after OR last/before
When paginating, make sure not to have first and last in the query. If you are using after or before, ensure the correct one is used.
BAD_REQUEST
INVALID_ARGUMENT
from cannot be larger or equal to to
Make sure that the end date for the time window is larger than the start date.
BAD_REQUEST
INVALID_ARGUMENT
Parameter cannot be used as DateTime
Make sure the provided date is in the ISO-8601 format.
BAD_REQUEST
INVALID_ARGUMENT
Invalid cursor value provided
Make sure the provided cursor value is the right one.
BAD_REQUEST
INVALID_ARGUMENT
Invalid {titleId/tournamentId/seriesId/teamId} value provided
The provided entity ID provided does not exist on Central Data.
Example responses
UNAUTHENTICATED

{
  "errors": [
    {
      "message": "Requester unauthorized to make query",
      "extensions": {
        "code": "DOWNSTREAM_SERVICE_ERROR",
        "serviceName": "gql-cd-internal-query",
        "errorType": "UNAUTHENTICATED",
        "exception": {
          "message": "Requester unauthorized to make query",
          "locations": [],
          "path": [
            "queryValue"
          ]
        }
      }
    }
  ],
  "data": {
    "title": null
  }
}

UNAUTHORIZED

{
  "errors": [
    {
      "message": "Requester forbidden to make query",
      "extensions": {
        "code": "DOWNSTREAM_SERVICE_ERROR",
        "serviceName": "gql-cd-internal-query",
        "errorType": "PERMISSION_DENIED",
        "exception": {
          "message": "Requester forbidden to make query",
          "locations": [],
          "path": [
            "queryValue"
          ]
        }
      }
    }
  ],
  "data": {
    "title": null
  }
}

BAD_REQUEST

{
  "errors": [
    {
      "message": "Parameter cannot be used as DateTime",
      "extensions": {
        "code": "DOWNSTREAM_SERVICE_ERROR",
        "serviceName": "gql-cd-internal-query",
        "errorType": "BAD_REQUEST",
        "errorDetail": "INVALID_ARGUMENT",
        "exception": {
          "message": "Parameter cannot be used as DateTime",
          "locations": [],
          "path": [
            "allSeries"
          ]
        }
      }
    }
  ],
  "data": null
}


WebSocket
WebSocket connections allow for two-way communication. Once a connection is opened, either side can push a message to the other side.

WebSocket connections are used to push data to the client. This enables clients to open a connection to listen for data even before it is available. This article presents details about the behavior of WebSocket connections.

Authentication
To confirm the client's identity, information must be passed on to each of the APIs.

When using WebSockets, the client's personal API key is used for authentication, and this information should be passed as a query parameter 'key' when opening a connection. For example, if the client's personal API key is 'secret-personal-api-key', then the URL to open a connection would be:

wss://api.grid.gg/live-data-feed/series/1?key=secret-personal-api-key

Clients can retrieve their personal API key here.

Limits
WebSockets usage is limited by GRID to ensure the system's proper operation, and these limits are applied per personal API key.

Rate limit
The rate limit restricts the client in terms of how many connection attempts can be made within a certain period.

By default, X connection attempts per minute are allowed. Any additional connection attempts will be rejected.

Connection limit
The connection limit restricts the number of simultaneous connections the client can have open at any given time.

By default, Y connections are allowed to be open at the same time. Any additional connection attempts will be rejected.

Duration limit
The duration limit restricts the duration for which a WebSocket connection can remain open. This limit is reset when a Data or Control frame is sent from either side.

By default, a WebSocket connection can remain open for up to 24 hours after the last activity. After this time, the WebSocket connection is automatically closed.

If clients wish to keep inactive WebSocket connections open for a longer period, they can periodically send WebSocket Ping Control frames.



Query and Mutations API
Table of Contents
Queries
allSeries
dataProviders
gameIdByExternalId
organization
organizations
player
playerIdByExternalId
playerRole
playerRoles
players
series
seriesByExternalSeries
seriesFormats
seriesIdByExternalId
team
teamIdByExternalId
teams
title
titles
tournament
tournamentIdByExternalId
tournaments
contentCatalogEntities
contentCatalogEntity
contentCatalogEntityIdByExternalId
contentCatalogVersion
contentCatalogVersionIdByExternalId
contentCatalogVersions
ingestionRequest
ingestionRequests
Mutations
addPlayerExternalLink
addSeriesExternalLink
addTeamExternalLink
addTournamentExternalLink
createPlayer
createPlayerRole
createSeries
createTeam
createTournament
deletePlayer
deletePlayerRole
deleteSeries
deleteTeam
deleteTournament
updatePlayer
updatePlayerRole
updateSeries
updateTeam
updateTournament
createContentCatalogCharacter
createContentCatalogItem
createContentCatalogMap
createContentCatalogStructure
createContentCatalogVersion
deleteContentCatalogEntity
deleteContentCatalogVersion
updateContentCatalogCharacter
updateContentCatalogItem
updateContentCatalogMap
updateContentCatalogStructure
updateContentCatalogVersion
createIngestionRequest
Objects
AddPlayerExternalLinkPayload
AddSeriesExternalLinkPayload
AddTeamExternalLinkPayload
AddTournamentExternalLinkPayload
Asset
Bounds
ContentCatalogCharacter
ContentCatalogEntityConnection
ContentCatalogEntityEdge
ContentCatalogItem
ContentCatalogMap
ContentCatalogVersion
ContentCatalogVersionConnection
ContentCatalogVersionEdge
Coordinates
CreateContentCatalogCharacterPayload
CreateContentCatalogItemPayload
CreateContentCatalogMapPayload
CreateContentCatalogStructurePayload
CreateContentCatalogVersionPayload
CreateIngestionRequestPayload
CreatePlayerPayload
CreatePlayerRolePayload
CreateSeriesPayload
CreateTeamPayload
CreateTournamentPayload
DataProvider
DeleteContentCatalogEntityPayload
DeleteContentCatalogVersionPayload
DeletePlayerPayload
DeletePlayerRolePayload
DeleteSeriesPayload
DeleteTeamPayload
DeleteTournamentPayload
ExternalEntity
ExternalLink
IngestionRequest
IngestionRequestConnection
IngestionRequestEdge
IngestionRequestStatusInformation
LolContentCatalogStructure
LolContentCatalogStructurePosition
MachineDataSource
Money
Nationality
NationalityDescriptor
Organization
OrganizationConnection
OrganizationEdge
OrganizationRelation
PageInfo
Player
PlayerConnection
PlayerDescriptor
PlayerEdge
PlayerRole
PlayerRoleDescriptor
ProductServiceLevel
Reason
Series
SeriesConnection
SeriesDescriptor
SeriesEdge
SeriesFormat
SeriesFormatDescriptor
Team
TeamConnection
TeamDescriptor
TeamEdge
TeamParticipant
TeamRelation
Title
Tournament
TournamentConnection
TournamentDescriptor
TournamentEdge
UpdateContentCatalogCharacterPayload
UpdateContentCatalogItemPayload
UpdateContentCatalogMapPayload
UpdateContentCatalogStructurePayload
UpdateContentCatalogVersionPayload
UpdatePlayerPayload
UpdatePlayerRolePayload
UpdateSeriesPayload
UpdateTeamPayload
UpdateTournamentPayload
VideoStream
Inputs
AddPlayerExternalLinkInput
AddSeriesExternalLinkInput
AddTeamExternalLinkInput
AddTournamentExternalLinkInput
AssetInput
BooleanFilter
BoundsInput
ContentCatalogEntityFilter
ContentCatalogEntityTitleFilter
ContentCatalogEntityTypeFilter
ContentCatalogVersionFilter
ContentCatalogVersionOrder
ContentCatalogVersionTitleFilter
CoordinatesInput
CreateContentCatalogCharacterInput
CreateContentCatalogItemInput
CreateContentCatalogMapInput
CreateContentCatalogStructureInput
CreateContentCatalogVersionInput
CreateIngestionRequestInput
CreateLolContentCatalogStructureInput
CreatePlayerInput
CreatePlayerPlayerRoleInput
CreatePlayerRoleInput
CreatePlayerRoleTitleInput
CreatePlayerTeamInput
CreatePlayerTitleInput
CreateSeriesFormatInput
CreateSeriesInput
CreateSeriesParticipatingPlayersInput
CreateSeriesTeamInput
CreateSeriesTitleInput
CreateSeriesTournamentInput
CreateTeamInput
CreateTeamTitleInput
CreateTournamentInput
CreateTournamentParentInput
CreateTournamentTeamInput
CreateTournamentTitleInput
DateInput
DateTimeFilter
DeletePlayerInput
DeletePlayerRoleInput
DeleteSeriesInput
DeleteTeamInput
DeleteTournamentInput
ExternalEntityInput
ExternalLinkInput
GameFilter
IdFilter
IngestionRequestFilter
IngestionRequestStatusFilter
IntFilter
LolContentCatalogStructurePositionInput
MapFilter
MoneyInput
NationalityDescriptorInput
NationalityFilter
NationalityInput
NullableDateFilter
NullableIdFilter
OrganizationFilter
PlayerDescriptorInput
PlayerFilter
PlayerPlayerRoleFilter
PlayerRoleDescriptorInput
PlayerRoleFilter
PlayerRoleTitleFilter
ProductServiceLevelFilter
SeriesDescriptorInput
SeriesFilter
SeriesFormatDescriptorInput
SeriesLiveFilter
SeriesPlayerFilter
SeriesTournamentFilter
SeriesTournamentTitleFilter
StringFilter
TeamDescriptorInput
TeamFilter
TitleDescriptorInput
TitleFilter
TournamentDescriptorInput
TournamentFilter
TournamentTitleFilter
UpdateContentCatalogCharacterInput
UpdateContentCatalogItemInput
UpdateContentCatalogMapInput
UpdateContentCatalogStructureInput
UpdateContentCatalogVersionInput
UpdateLolContentCatalogStructureInput
UpdatePlayerInput
UpdatePlayerPlayerRoleInput
UpdatePlayerRoleInput
UpdatePlayerRoleTitleInput
UpdatePlayerTeamInput
UpdatePlayerTitleInput
UpdateSeriesFormatInput
UpdateSeriesInput
UpdateSeriesParticipatingPlayersInput
UpdateSeriesTeamInput
UpdateSeriesTitleInput
UpdateSeriesTournamentInput
UpdateTeamInput
UpdateTeamTitleInput
UpdateTournamentInput
UpdateTournamentParentInput
UpdateTournamentTeamInput
UpdateTournamentTitleInput
Interfaces
ContentCatalogEntity
ContentCatalogStructure
OrganizationInterface
TeamInterface
Enums
AssetType
ContentCatalogEntityType
ContentCatalogVersionOrderField
ErrorDetail
ErrorType
IngestionRequestOrderBy
IngestionRequestStatus
OrderDirection
PlayerType
SeriesOrderBy
SeriesType
SeriesWorkflowStatus
ServiceLevel
TeamWorkflowStatus
TournamentType
TournamentVenueType
Scalars
Boolean
Date
DateTime
Decimal
Float
HexColor
ID
Int
Locale
String
Url
Queries
allSeries
Get multiple series by query specifications. Returns SeriesConnection!.

Argument
Type
Description
after
String
Cursor, after which results to query.
before
String
Cursor, before which results to query.
filter
SeriesFilter
Filters for series results.
first
Int
Amount of nodes to return starting from the beginning. Defaults to 10.
last
Int
Amount of nodes to return starting from the end.
orderBy
SeriesOrderBy!
Ordering of resulting series set.
orderDirection
OrderDirection!
Whether to order by ascending or descending direction.
dataProviders
Get supported data providers. Returns [DataProvider!]!.

No arguments.

gameIdByExternalId
Get the GRID Game ID by an external ID. Returns ID.

Argument
Type
Description
dataProviderName
String!
The name of the data provider to look by. An example could be "STEAM"
externalGameId
ID!
The external ID to look for.
organization
Get an organization by ID. Returns Organization.

Argument
Type
Description
id
ID!
The ID of the organization to retrieve.
organizations
Get multiple organizations by query specifications. Returns OrganizationConnection!.

Argument
Type
Description
after
String
Cursor, after which results to query.
before
String
Cursor, before which results to query.
filter
OrganizationFilter
Filters for organization results.
first
Int
Amount of nodes to return starting from the beginning. Defaults to 10.
last
Int
Amount of nodes to return starting from the end.
player
Get a player by ID. Returns Player.

Argument
Type
Description
id
ID!
The ID of the player to retrieve.
playerIdByExternalId
Get the GRID Player ID by an external ID. Returns ID.

Argument
Type
Description
dataProviderName
String!
The name of the data provider to look by. An example could be "STEAM"
externalPlayerId
ID!
The external ID to look for.
titleId
ID
The title player belongs to.
playerRole
Get a player role by ID. Returns PlayerRole.

Argument
Type
Description
id
ID!
None
playerRoles
Get multiple player roles by query specifications. Returns [PlayerRole!]!.

Argument
Type
Description
filter
PlayerRoleFilter
None
players
Get multiple players by query specifications. Returns PlayerConnection!.

Argument
Type
Description
after
String
Cursor, after which results to query.
before
String
Cursor, before which results to query.
filter
PlayerFilter
Filters for player results.
first
Int
Amount of nodes to return starting from the beginning. Defaults to 10.
last
Int
Amount of nodes to return starting from the end.
series
Get a series by ID. Returns Series.

Argument
Type
Description
id
ID!
The ID of the series to retrieve.
seriesByExternalSeries
@deprecated(reason: 'Use seriesIdByExternalId instead') Get a series by its external ID. Returns Series.

Argument
Type
Description
dataProviderName
String!
The name of the data provider to look by. An example could be "STEAM"
externalSeriesId
String!
The external ID to look for.
seriesFormats
Get supported series formats. Returns [SeriesFormat!]!.

No arguments.

seriesIdByExternalId
Get the GRID Series ID by an external ID. Returns ID.

Argument
Type
Description
dataProviderName
String!
The name of the data provider to look by. An example could be "STEAM"
externalSeriesId
ID!
The external ID to look for.
team
Get a team by ID. Returns Team.

Argument
Type
Description
id
ID!
The ID of the team to retrieve.
teamIdByExternalId
Get the GRID Team ID by an external ID. Returns ID.

Argument
Type
Description
dataProviderName
String!
The name of the data provider to look by. An example could be "STEAM"
externalTeamId
ID!
The external ID to look for.
titleId
ID!
The title in which the team participates.
teams
Get multiple teams by query specifications. Returns TeamConnection!.

Argument
Type
Description
after
String
Cursor, after which results to query.
before
String
Cursor, before which results to query.
filter
TeamFilter
Filters for team results.
first
Int
Amount of nodes to return starting from the beginning. Defaults to 10.
last
Int
Amount of nodes to return starting from the end.
title
Get a title by ID. Returns Title.

Argument
Type
Description
id
ID!
The ID of the title to retrieve.
titles
Get all titles. Returns [Title!]!.

Argument
Type
Description
filter
TitleFilter
Filters for titles results.
tournament
Get a tournament by ID. Returns Tournament.

Argument
Type
Description
id
ID!
The ID of the tournament to retrieve.
tournamentIdByExternalId
Get the GRID Tournament ID by an external ID. Returns ID.

Argument
Type
Description
dataProviderName
String!
The name of the data provider to look by. An example could be "STEAM"
externalTournamentId
ID!
The external ID to look for.
tournaments
Get multiple tournaments by query specifications. Returns TournamentConnection!.

Argument
Type
Description
after
String
Cursor, after which results to query.
before
String
Cursor, before which results to query.
filter
TournamentFilter
Filters for tournament results.
first
Int
Amount of nodes to return starting from the beginning. Defaults to 10.
last
Int
Amount of nodes to return starting from the end.
orderDirection
OrderDirection!
Whether to order the results by ID in ascending or descending direction.
contentCatalogEntities
Get multiple content catalog entities by query specification. Returns ContentCatalogEntityConnection!.

Argument
Type
Description
after
String
Returns the elements in the list that come after the specified cursor.
before
String
Returns the elements in the list that come before the specified cursor.
contentCatalogVersionId
ID
The content catalog version Id to search entities by.If not specified, the latest version is selected.
filter
ContentCatalogEntityFilter
Filter specifications.
first
Int
Number of nodes to return starting from the beginning. Defaults to 10.
last
Int
Number of nodes to return starting from the end.
contentCatalogEntity
Get a Content Catalog Entity by ID. Returns ContentCatalogEntity.

Argument
Type
Description
acceptLanguage
Locale
The language to return the localized content in.
contentCatalogVersionId
ID
The content catalog version of the item.If not specified, the query returns the item definition of the latest version.
id
ID!
The ID of a Content Catalog Entity.
contentCatalogEntityIdByExternalId
Get the GRID Content Catalog Entity ID by an external ID. Returns ID.

Argument
Type
Description
contentCatalogVersionId
ID
The content catalog version of the item.If not specified, the query returns the id of the latest item definition.
dataProviderName
String!
The name of the data provider to search by. An example could be 'STEAM'.
externalEntityId
ID!
The external ID to search by.
contentCatalogVersion
Get a Content Catalog Version by ID. Returns ContentCatalogVersion.

Argument
Type
Description
id
ID!
The ID of a Content Catalog Version.
contentCatalogVersionIdByExternalId
Get the GRID Content Catalog Version ID by an external ID. Returns ID.

Argument
Type
Description
dataProviderName
String!
The name of the data provider to search by. An example could be "STEAM".
externalEntityId
ID!
The external ID to search by.
contentCatalogVersions
Get multiple content catalog versions by query specification. Returns ContentCatalogVersionConnection!.

Argument
Type
Description
after
String
Returns the elements in the list that come after the specified cursor.
before
String
Returns the elements in the list that come before the specified cursor.
filter
ContentCatalogVersionFilter
Filter specifications.
first
Int
Number of nodes to return starting from the beginning. Defaults to 10.
last
Int
Number of nodes to return starting from the end.
orderBy
[ContentCatalogVersionOrder!]
Ordering options.
ingestionRequest
Get a single Ingestion Request. Returns IngestionRequest.

Argument
Type
Description
id
ID!
None
ingestionRequests
Get a list of filtered Ingestion Requests. Returns IngestionRequestConnection!.

Argument
Type
Description
after
String
The cursor after which to query results.
before
String
The cursor before which to query results.
filter
IngestionRequestFilter
Filters for Ingestion Request results.
first
Int
The number of nodes to return starting from the beginning. Defaults to 10.
last
Int
The number of nodes to return starting from the end.
orderBy
IngestionRequestOrderBy!
The ordering of the resulting Ingestion Request set.
orderDirection
OrderDirection!
The direction for ordering results (ascending or descending).
Mutations
addPlayerExternalLink
Appends an external link to an existing GRID player. Returns AddPlayerExternalLinkPayload!.

Argument
Type
Description
addPlayerExternalLinkInput
AddPlayerExternalLinkInput!
Expects the ID of the player and the external link to be appended.
addSeriesExternalLink
Appends an external link to an existing GRID series. Returns AddSeriesExternalLinkPayload!.

Argument
Type
Description
addSeriesExternalLinkInput
AddSeriesExternalLinkInput!
Expects the ID of the series and the external link to be appended.
addTeamExternalLink
Appends an external link to an existing GRID team. Returns AddTeamExternalLinkPayload!.

Argument
Type
Description
addTeamExternalLinkInput
AddTeamExternalLinkInput!
Expects the ID of the team and the external link to be appended.
addTournamentExternalLink
Appends an external link to an existing GRID tournament. Returns AddTournamentExternalLinkPayload!.

Argument
Type
Description
addTournamentExternalLinkInput
AddTournamentExternalLinkInput!
Expects the ID of the tournament and the external link to be appended.
createPlayer
Creates a new GRID player. Returns CreatePlayerPayload!.

Argument
Type
Description
createPlayerInput
CreatePlayerInput!
Expects the nickname, the team the player belongs to and the title the player participates in and optionally the external links associated with the player.
createPlayerRole
Creates a new player role. Returns CreatePlayerRolePayload!.

Argument
Type
Description
createPlayerRoleInput
CreatePlayerRoleInput!
Expects the name of the role and title that the role is associated with.
createSeries
Creates a new GRID Series. Returns CreateSeriesPayload!.

Argument
Type
Description
createSeriesInput
CreateSeriesInput!
Expects the title for the series, the tournament where it's being played, the teams playing it, the start date and time, and the format.
createTeam
Creates a new GRID team. Returns CreateTeamPayload!.

Argument
Type
Description
createTeamInput
CreateTeamInput!
Expects the name, the title the team participates in and optionally the primary and/or secondary colors, the logo of the team and the external links associated with the team.
createTournament
Creates a new GRID tournament. Returns CreateTournamentPayload!.

Argument
Type
Description
createTournamentInput
CreateTournamentInput!
Expects the name and short name, and optionally the titles for the tournament, the teams playing in it, and the logo of the tournament.
deletePlayer
Deletes an existing GRID player. Returns DeletePlayerPayload!.

Argument
Type
Description
deletePlayerInput
DeletePlayerInput!
Expects the ID of an existing player.
deletePlayerRole
Deletes a specific player role. Returns DeletePlayerRolePayload!.

Argument
Type
Description
deletePlayerRoleInput
DeletePlayerRoleInput!
Expects the ID of an existing player role.
deleteSeries
Deletes an existing GRID Series. Returns DeleteSeriesPayload!.

Argument
Type
Description
deleteSeriesInput
DeleteSeriesInput!
Expects the ID of an existing Series.
deleteTeam
Deletes an existing GRID team. Returns DeleteTeamPayload!.

Argument
Type
Description
deleteTeamInput
DeleteTeamInput!
Expects the ID of an existing team.
deleteTournament
Deletes an existing GRID tournament. Returns DeleteTournamentPayload!.

Argument
Type
Description
deleteTournamentInput
DeleteTournamentInput!
Expects the ID of an existing tournament.
updatePlayer
Updates an existing GRID player. Returns UpdatePlayerPayload!.

Argument
Type
Description
updatePlayerInput
UpdatePlayerInput!
Expects the ID of the player, and optionally the nickname, the team the player belongs to, the title the player participates in and the external links associated with the player.
updatePlayerRole
Updates a specific player role. Returns UpdatePlayerRolePayload!.

Argument
Type
Description
updatePlayerRoleInput
UpdatePlayerRoleInput!
Expects the ID of the role to be updated and its new name.
updateSeries
Updates an existing GRID Series. Returns UpdateSeriesPayload!.

Argument
Type
Description
updateSeriesInput
UpdateSeriesInput!
Expects the ID, and optionally can receive: the title, the tournament, the teams, the start date and time, and/or the format.
updateTeam
Updates an existing GRID team. Returns UpdateTeamPayload!.

Argument
Type
Description
updateTeamInput
UpdateTeamInput!
Expects the ID of the team and optionally the name, the title the team participates in, the primary and the primary and/or secondary colors, the logo of the team and the external links associated with the team.
updateTournament
Updates an existing GRID tournament. Returns UpdateTournamentPayload!.

Argument
Type
Description
updateTournamentInput
UpdateTournamentInput!
Expects the ID of the tournament, and optionally the name, the short name, the titles for the tournament, the teams playing in it, and/or the logo of the tournament.
createContentCatalogCharacter
Creates a new version for a content catalog character. Returns CreateContentCatalogCharacterPayload!.

Argument
Type
Description
createContentCatalogCharacterInput
CreateContentCatalogCharacterInput!
None
createContentCatalogItem
Creates a new version for a content catalog item. Returns CreateContentCatalogItemPayload!.

Argument
Type
Description
createContentCatalogItemInput
CreateContentCatalogItemInput!
None
createContentCatalogMap
Creates a new version for a content catalog map. Returns CreateContentCatalogMapPayload!.

Argument
Type
Description
createContentCatalogMapInput
CreateContentCatalogMapInput!
None
createContentCatalogStructure
Creates a new version for a content catalog structure. Returns CreateContentCatalogStructurePayload!.

Argument
Type
Description
createContentCatalogStructureInput
CreateContentCatalogStructureInput!
None
createContentCatalogVersion
Creates a new content catalog version. Returns CreateContentCatalogVersionPayload!.

Argument
Type
Description
createContentCatalogVersionInput
CreateContentCatalogVersionInput!
None
deleteContentCatalogEntity
Deletes a specific version of a content catalog entity. Returns DeleteContentCatalogEntityPayload!.

Argument
Type
Description
contentCatalogVersionId
ID!
The version of the content catalog entity to be deleted.
id
ID!
The ID of the entity to be deleted.
deleteContentCatalogVersion
Deletes a specific content catalog version. Returns DeleteContentCatalogVersionPayload!.

Argument
Type
Description
id
ID!
The ID of the content catalog version to delete.
updateContentCatalogCharacter
Updates a specific version of a content catalog character. Returns UpdateContentCatalogCharacterPayload!.

Argument
Type
Description
updateContentCatalogCharacterInput
UpdateContentCatalogCharacterInput!
None
updateContentCatalogItem
Updates a specific version of a content catalog item. Returns UpdateContentCatalogItemPayload!.

Argument
Type
Description
updateContentCatalogItemInput
UpdateContentCatalogItemInput!
None
updateContentCatalogMap
Updates a specific version of a content catalog map. Returns UpdateContentCatalogMapPayload!.

Argument
Type
Description
updateContentCatalogMapInput
UpdateContentCatalogMapInput!
None
updateContentCatalogStructure
Updates a specific version of a content catalog structure. Returns UpdateContentCatalogStructurePayload!.

Argument
Type
Description
updateContentCatalogStructureInput
UpdateContentCatalogStructureInput!
None
updateContentCatalogVersion
Updates a specific content catalog version Returns UpdateContentCatalogVersionPayload!.

Argument
Type
Description
updateContentCatalogVersionInput
UpdateContentCatalogVersionInput!
None
createIngestionRequest
Submit an Ingestion Request. Returns CreateIngestionRequestPayload!.

Argument
Type
Description
input
CreateIngestionRequestInput!
None
Objects
AddPlayerExternalLinkPayload
Field
Type
Description
updatedPlayer
Player!
The player information after the external link has been added.
AddSeriesExternalLinkPayload
Field
Type
Description
updatedSeries
Series!
The series information after the external link has been added.
AddTeamExternalLinkPayload
Field
Type
Description
updatedTeam
Team!
The team information after the external link has been added.
AddTournamentExternalLinkPayload
Field
Type
Description
updatedTournament
Tournament!
The tournament information after the external link has been added.
Asset
An asset.

Field
Type
Description
assetType
AssetType!
The type of the asset.
url
Url!
The asset URL.
Bounds
Bounds information (i.e. min and max coordinates for a map).

Field
Type
Description
max
Coordinates!
Maximum coordinates value.
min
Coordinates!
Minimum coordinates value.
ContentCatalogCharacter
A content catalog character.

Field
Type
Description
contentCatalogVersion
ContentCatalogVersion!
The content catalog version this entity appears in.
externalLinks
[ExternalLink!]!
The content catalog entity version external IDs and providers.
id
ID!
The ID of a content catalog entity.
imageUrl
Url!
The URL to the content catalog entity's image
name
String!
The name of a content catalog entity.
private
Boolean
The visibility of this entity to other users.
updatedAt
DateTime!
The date and time of the last update.
ContentCatalogEntityConnection
The result of querying list of content catalog entities, which provides additional information.

Field
Type
Description
edges
[ContentCatalogEntityEdge!]!
All content catalog entities edges for a query on a given page.
pageInfo
PageInfo!
Information about the current page being pulled with a query.
totalCount
Int!
Total count of content catalog entities matching the filtering in the API.
ContentCatalogEntityEdge
Content catalog entity and cursor information.

Field
Type
Description
cursor
String!
Content catalog entity cursor value for the edge.
node
ContentCatalogEntity!
A content catalog entity.
ContentCatalogItem
A content catalog item.

Field
Type
Description
contentCatalogVersion
ContentCatalogVersion!
The content catalog version this entity appears in.
cost
Float!
The cost of an item. When not set, defaults to 0.
externalLinks
[ExternalLink!]!
The content catalog entity version external IDs and providers.
id
ID!
The ID of the content catalog entity.
imageUrl
Url!
The URL to the content catalog entity's image.
name
String!
The name of the content catalog entity.
private
Boolean
The visibility of this entity to other users.
updatedAt
DateTime!
The date and time of the last update.
ContentCatalogMap
Field
Type
Description
bounds
Bounds!
The boundaries of the map.
contentCatalogVersion
ContentCatalogVersion!
The content catalog version this entity appears in.
externalLinks
[ExternalLink!]!
The content catalog entity version external IDs and providers.
id
ID!
The ID of a content catalog entity.
imageUrl
Url!
The URL to the content catalog entity's image
name
String!
The name of a content catalog entity.
private
Boolean
The visibility of this entity to other users.
updatedAt
DateTime!
The date and time of the last update.
ContentCatalogVersion
A content catalog version.

Field
Type
Description
externalLinks
[ExternalLink!]!
The content catalog version external IDs and providers.
id
ID!
The unique identifier of the content catalog version.
name
String!
The unique name of the content catalog version. E.g. 7.34.
private
Boolean
The visibility of this version to other users.
publishedOn
DateTime!
The date and time when the content catalog version becomes effective.
title
Title!
The title which the content catalog version belongs to.
ContentCatalogVersionConnection
The result of querying list of content catalog versions, which provides additional information.

Field
Type
Description
edges
[ContentCatalogVersionEdge!]!
All content catalog version edges for a query on a given page.
pageInfo
PageInfo!
Information about the current page being pulled with a query.
totalCount
Int!
Total count of content catalog versions matching the filtering in the API.
ContentCatalogVersionEdge
Content catalog version and cursor information.

Field
Type
Description
cursor
String!
Content catalog version cursor value for the edge.
node
ContentCatalogVersion!
A content catalog version.
Coordinates
Spatial coordinates.

Field
Type
Description
x
Float!
X coordinate value.
y
Float!
Y coordinate value.
CreateContentCatalogCharacterPayload
The response for the createContentCatalogCharacter mutation.

Field
Type
Description
createdContentCatalogCharacter
ContentCatalogCharacter!
The created content catalog character information.
CreateContentCatalogItemPayload
The response for the createContentCatalogItem mutation.

Field
Type
Description
createdContentCatalogItem
ContentCatalogItem!
The created content catalog version information.
CreateContentCatalogMapPayload
The response for the createContentCatalogMap mutation.

Field
Type
Description
createdContentCatalogMap
ContentCatalogMap!
The created content catalog map information.
CreateContentCatalogStructurePayload
The response for the createContentCatalogStructure mutation.

Field
Type
Description
createdContentCatalogStructure
ContentCatalogStructure!
The created content catalog structure information.
CreateContentCatalogVersionPayload
The response for the createContentCatalogVersion mutation.

Field
Type
Description
createdContentCatalogVersion
ContentCatalogVersion!
The created content catalog version information.
CreateIngestionRequestPayload
----- PAYLOAD ----
Field
Type
Description
createdIngestionRequestId
ID!
None
CreatePlayerPayload
The response for the createPlayer mutation.

Field
Type
Description
createdPlayer
Player!
The created player information.
CreatePlayerRolePayload
The response for the createPlayerRole mutation.

Field
Type
Description
createdPlayerRole
PlayerRole!
The created player role information.
CreateSeriesPayload
The response for the createSeries mutation.

Field
Type
Description
createdSeries
Series!
The created series information.
CreateTeamPayload
The response for the createTeam mutation.

Field
Type
Description
createdTeam
Team!
The created team information.
CreateTournamentPayload
The response for the createTournament mutation.

Field
Type
Description
createdTournament
Tournament!
The created tournament information.
DataProvider
A data provider which provides external entity IDs.

Field
Type
Description
description
String
Description of the data provider.
name
String!
The name of the external data provider.
DeleteContentCatalogEntityPayload
The response for the deleteContentCatalogEntity mutation.

Field
Type
Description
contentCatalogVersionId
ID!
The version of the deleted content catalog entity.
id
ID!
The ID of the deleted content catalog entity.
DeleteContentCatalogVersionPayload
The response for the deleteContentCatalogVersion mutation.

Field
Type
Description
id
ID!
The ID of the deleted content catalog version.
DeletePlayerPayload
The ID of the removed player.

Field
Type
Description
id
ID!
The player ID.
DeletePlayerRolePayload
The response for the deletePlayerRole mutation.

Field
Type
Description
deletedPlayerRole
ID!
The ID of the deleted player role.
DeleteSeriesPayload
The ID of the removed series.

Field
Type
Description
id
ID!
The series ID.
DeleteTeamPayload
The ID of the removed team.

Field
Type
Description
id
ID!
The team ID.
DeleteTournamentPayload
The ID of the removed tournament.

Field
Type
Description
id
ID!
The tournament ID.
ExternalEntity
Information provided by the external data provider.

Field
Type
Description
id
ID!
The ID used by the external data provider.
ExternalLink
Linking information for an external data provider.

Field
Type
Description
dataProvider
DataProvider!
The external data provider.
externalEntity
ExternalEntity!
The external information.
IngestionRequest
A request to ingest the described entities.

Field
Type
Description
createdAt
DateTime!
The creation date and time of the ingestion request.
id
ID!
An identifier of the ingestion request.
playerRoles
[PlayerRoleDescriptor!]!
The player roles described in the ingestion request.
players
[PlayerDescriptor!]!
The players described in the ingestion request.
series
[SeriesDescriptor!]!
The series described in the ingestion request.
statusInformation
IngestionRequestStatusInformation!
Information regarding the current status of the ingestion request.
teams
[TeamDescriptor!]!
The teams described in the ingestion request.
tournaments
[TournamentDescriptor!]!
The tournaments described in the ingestion request.
updatedAt
DateTime!
The last update date and time of the ingestion request.
IngestionRequestConnection
The result of querying a list of Ingestion Requests. It provides pagination information.

Field
Type
Description
edges
[IngestionRequestEdge!]!
The Ingestion Request edges for the current page of the query.
pageInfo
PageInfo!
Information about the current page being retrieved by the query.
totalCount
Int!
Total count of Ingestion Requests in the API.
IngestionRequestEdge
Ingestion Request and cursor information.

Field
Type
Description
cursor
String!
Ingestion Request cursor value for the edge.
node
IngestionRequest!
Ingestion Request.
IngestionRequestStatusInformation
Field
Type
Description
reason
Reason
The reason the request status was changed.
status
IngestionRequestStatus!
Current status of the ingestion request.
LolContentCatalogStructure
A content catalog structure in LoL.

Field
Type
Description
alias
String!
The natural ID of the structure that is unique for LoL.
contentCatalogVersion
ContentCatalogVersion!
The content catalog version this entity appears in.
externalLinks
[ExternalLink!]!
The content catalog entity version external IDs and providers.
gridType
String!
The GRID type of the structure.
id
ID!
The ID of a content catalog entity.
imageUrl
Url!
The URL to the content catalog entity's image
lane
String
The lane the structure is in.
maxHealth
String!
The maximum health of the structure.
name
String!
The name of a content catalog entity.
nexusTurretName
String
The name of the Nexus turret if the structure is one, according to its turretTier.
position
LolContentCatalogStructurePosition!
The position of the structure expressed with two coordinates.
private
Boolean
The visibility of this entity to other users.
respawnTime
String
The time after which the structure respawns.
riotType
String!
The type of the structure as provided by Riot Games.
teamId
Int!
The ID of the team the structure belongs to.
turretTier
String
The tier of the structure if it is a turret, according to its riotType.
updatedAt
DateTime!
The date and time of the last update.
LolContentCatalogStructurePosition
The position of a LoL structure.

Field
Type
Description
x
Float!
None
z
Float!
None
MachineDataSource
Field
Type
Description
id
ID!
None
Money
A monetary value.

Field
Type
Description
amount
Decimal!
The amount in USD.
Nationality
A player's nationality.

Field
Type
Description
code
String!
The ISO 3166-1 Alpha-3 code for the country (e.g., 'USA' for the United States).
name
String!
The full name of the country (e.g., 'United States').
NationalityDescriptor
A description of a player's nationality.

Field
Type
Description
code
String
The ISO 3166-1 Alpha-3 code for the country (e.g., 'USA' for the United States).
name
String
The full name of the country (e.g., 'United States').
Organization
An organization.

Field
Type
Description
id
ID!
The ID of the organization.
name
String!
The name of the organization.
private
Boolean!
The visibility status of the organization.
teams
[TeamRelation!]
The teams that belong to the organization.
updatedAt
DateTime!
The date and time of the last update.
OrganizationConnection
The result of querying a list of organizations. It provides pagination information.

Field
Type
Description
edges
[OrganizationEdge!]!
All players edges for query on a given page.
pageInfo
PageInfo!
Information about current page being pulled with query.
totalCount
Int!
Total count of players in the API.
OrganizationEdge
Organization and cursor information.

Field
Type
Description
cursor
String!
Organizations cursor value for the edge.
node
Organization!
Organizations.
OrganizationRelation
An organization relation. Removes connection to Teams from the Organization.

Field
Type
Description
id
ID!
The ID of the organization.
name
String!
The name of the organization.
PageInfo
Information about current paginated data set position.

Field
Type
Description
endCursor
String
The cursor of the last edge of the page.
hasNextPage
Boolean!
If there is a page after the one being returned.
hasPreviousPage
Boolean!
If there is a page before the one being returned.
startCursor
String
The cursor of the first edge of the page.
Player
A player.

Field
Type
Description
id
ID!
The ID of a player
age
Int
The player's age.
externalLinks
[ExternalLink!]!
The external player IDs and providers.
fullName
String
The player's full name.
imageUrl
Url!
The player image URL.
nationality
[Nationality!]!
The player's nationality.
nickname
String!
The name of the player.
private
Boolean!
The visibility status of the player.
roles
[PlayerRole!]!
The player's roles (e.g. MID, SUPPORT).
team
Team
The team the player is currently playing.
title
Title!
The title the player is currently playing.
type
PlayerType!
The player's type.
updatedAt
DateTime!
The date and time of the last update.
PlayerConnection
The result of querying a list of players. It provides pagination information.

Field
Type
Description
edges
[PlayerEdge!]!
All players edges for query on a given page.
pageInfo
PageInfo!
Information about current page being pulled with query.
totalCount
Int!
Total count of players in the API.
PlayerDescriptor
A description of a player.

Field
Type
Description
dataSource
String
The data source that sourceId belongs to. Must be set together with sourceId.
dateOfBirth
Date
The described player's date of birth.
fullName
String
The described player's full name.
id
ID
ID of this descriptor.
nationality
[NationalityDescriptor!]
The description of the nationality of the described player.
nickname
String!
The name of the described player.
roles
[PlayerRoleDescriptor!]
The description of player roles (e.g., MID, SUPPORT) of the described player.
sourceId
ID
The player identifier on the source platform.
team
TeamDescriptor!
The description of a team where the described player currently plays.
PlayerEdge
Player and cursor information.

Field
Type
Description
cursor
String!
Players cursor value for the edge.
node
Player!
Players.
PlayerRole
A player role

Field
Type
Description
id
ID!
The ID of a player role
name
String!
The name of the player role.
private
Boolean!
The visibility of this entity to other users.
title
Title!
The title of the player role.
PlayerRoleDescriptor
A description of player's role (e.g., MID, SUPPORT).

Field
Type
Description
name
String!
The name of the described player role.
ProductServiceLevel
The availability of a product for a given series.

Field
Type
Description
productName
String!
The name of the product.
serviceLevel
ServiceLevel!
The availability level for the product.
Reason
A reason for which the status of the ingestion request was changed.

Field
Type
Description
message
String!
A textual description of the reason the request status was changed.
Series
A series.

Field
Type
Description
id
ID!
The ID of a series.
externalLinks
[ExternalLink!]!
The external series IDs and providers.
format
SeriesFormat!
The format of the series.
machineDataSource
MachineDataSource
The machine data source ID for the series.
players
[Player!]!
The players participating in the series.
private
Boolean!
The visibility status of the series.
productServiceLevels
[ProductServiceLevel!]!
The service level for the series.
startTimeScheduled
DateTime!
The scheduled start time of the series.
streams
[VideoStream!]!
The list of stream URLs available for the series.
teams
[TeamParticipant!]!
The teams participating in the series.
title
Title!
The title being played in the series.
tournament
Tournament!
The tournament in which the series is being played.
type
SeriesType!
The type of the series.
updatedAt
DateTime!
The date and time of the last update.
workflowStatus
SeriesWorkflowStatus!
The workflow status of the series.
SeriesConnection
The result of querying a list of series. It provides pagination information.

Field
Type
Description
edges
[SeriesEdge!]!
All series edges for query on a given page.
pageInfo
PageInfo!
Information about current page being pulled with query.
totalCount
Int!
Total count of series in the API.
SeriesDescriptor
A description of a series.

Field
Type
Description
dataSource
String
The data source that sourceId belongs to. Must be set together with sourceId.
format
SeriesFormatDescriptor!
The format of the described series.
id
ID
ID of this descriptor.
players
[PlayerDescriptor!]
The players participating in the described series.
sourceId
ID
The series identifier on the source platform.
startTimeScheduled
DateTime!
The scheduled start time of the described series.
teams
[TeamDescriptor!]!
The teams participating in the described series.
tournament
TournamentDescriptor!
The description of a tournament in which the described series is being played.
SeriesEdge
Series and cursor information.

Field
Type
Description
cursor
String!
Series cursor value for the edge.
node
Series!
Series.
SeriesFormat
The format of the series, such as: BO1, BO3, SA2.

Field
Type
Description
id
ID
The ID of the format.
name
String!
The name of a format.
nameShortened
String!
The shortened name of a format (max. 5 characters).
SeriesFormatDescriptor
The format of the series (e.g., Bo1, Bo3, SA2).

Field
Type
Description
nameShortened
String!
The shortened name of the described format.
Team
A team.

Field
Type
Description
id
ID!
The ID of a team.
assets
[Asset!]!
All logo variants for this team.
colorPrimary
HexColor!
The primary color of the team.
colorSecondary
HexColor!
The secondary color of the team.
externalLinks
[ExternalLink!]!
The external team IDs and providers.
logoUrl
Url!
The team logo URL.
name
String!
The name of the team.
nameShortened
String
The shortened name of a team.
organization
OrganizationRelation
The organization the team plays for.
private
Boolean!
The visibility status of the team.
rating
Float
The team rating.
title
Title!
Title the team is related to. 'titles' field will be removed soon and every team will have a single title.
titles
[Title!]!
Use title field instead. The titles the team is related to.
updatedAt
DateTime!
The date and time of the last update.
workflowStatus
TeamWorkflowStatus!
The workflow status of the team.
TeamConnection
The result of querying a list of teams. It provides pagination information.

Field
Type
Description
edges
[TeamEdge!]!
All team edges for query on a given page.
pageInfo
PageInfo!
Information about current page being pulled with query.
totalCount
Int!
Total count of teams in the API.
TeamDescriptor
A description of a team.

Field
Type
Description
dataSource
String
The data source that sourceId belongs to. Must be set together with sourceId.
id
ID
ID of this descriptor.
name
String!
The name of the described team.
nameShortened
String
The shortened name of the described team.
sourceId
ID
The team identifier on the source platform.
TeamEdge
Team and cursor information.

Field
Type
Description
cursor
String!
Team cursor value cursor value for the edge.
node
Team!
Team.
TeamParticipant
Team participating in the series.

Field
Type
Description
baseInfo
Team!
The team base information.
scoreAdvantage
Int!
The team score advantage for the team in the series.
TeamRelation
A team relation. Removes connection to Organization from the Team.

Field
Type
Description
assets
[Asset!]!
All logo variants for this team.
colorPrimary
HexColor!
The primary color of the team.
colorSecondary
HexColor!
The secondary color of the team.
id
ID!
The ID of the team.
logoUrl
Url!
The team logo URL.
name
String!
The name of the team.
nameShortened
String
The shortened name of a team.
rating
Float
The team rating.
titles
[Title!]!
The titles the team is related to.
Title
A title.

Field
Type
Description
id
ID!
The ID of the title.
logoUrl
Url!
The title logo URL.
name
String!
The name of the title.
nameShortened
String!
The shortened name of the title (max. 20 characters).
private
Boolean
The visibility status of the title.
Tournament
A tournament.

Field
Type
Description
id
ID!
The ID of a tournament.
children
[Tournament!]!
The children tournament of the tournament.
endDate
Date
The end date of the tournament.
externalLinks
[ExternalLink!]!
The external tournament IDs and providers.
logoUrl
Url!
The tournament logo URL.
name
String!
The name of the tournament.
nameShortened
String!
The shortened name of the tournament (max. 30 characters).
parent
Tournament
The parent tournament of the tournament.
private
Boolean!
The visibility status of the tournament.
prizePool
Money
The tournament prize pool in USD.
startDate
Date
The start date of the tournament.
teams
[Team!]!
The teams linked to the tournament.
titles
[Title!]!
The titles linked to the tournament.
tournamentType
TournamentType!
The type of the tournament within its tournament tree.
updatedAt
DateTime!
The date and time of the last update.
venueType
TournamentVenueType!
The type of venue where the tournament is held.
TournamentConnection
The result of querying a list of tournaments. It provides pagination information.

Field
Type
Description
edges
[TournamentEdge!]!
All tournament edges for query on a given page.
pageInfo
PageInfo!
Information about current page being pulled with query.
totalCount
Int!
Total count of tournaments in the API.
TournamentDescriptor
A description of a tournament.

Field
Type
Description
dataSource
String
The data source that sourceId belongs to. Must be set together with sourceId.
endDate
Date
The end date of the described tournament.
id
ID
ID of this descriptor.
name
String!
The name of the described tournament.
nameShortened
String
The shortened name of the described tournament (max. 30 characters).
parent
TournamentDescriptor
The description of the parent tournament of the described tournament.
prizePool
Money
The prize pool of the described tournament in USD.
sourceId
ID
The tournament identifier on the source platform.
startDate
Date
The start date of the described tournament.
venueType
TournamentVenueType
The type of venue where the described tournament is held.
TournamentEdge
Tournament and cursor information.

Field
Type
Description
cursor
String!
Tournament cursor value for the edge.
node
Tournament!
Tournament.
UpdateContentCatalogCharacterPayload
The response for the updateContentCatalogCharacter mutation.

Field
Type
Description
updatedContentCatalogCharacter
ContentCatalogCharacter!
The updated content catalog character information.
UpdateContentCatalogItemPayload
The response for the updateContentCatalogItem mutation.

Field
Type
Description
updatedContentCatalogItem
ContentCatalogItem!
The updated content catalog version information.
UpdateContentCatalogMapPayload
The response for the updateContentCatalogMap mutation.

Field
Type
Description
updatedContentCatalogMap
ContentCatalogMap!
The updated content catalog map information.
UpdateContentCatalogStructurePayload
The response for the updateContentCatalogStructure mutation.

Field
Type
Description
updatedContentCatalogStructure
ContentCatalogStructure!
The updated content catalog structure information.
UpdateContentCatalogVersionPayload
The response for the updateContentCatalogVersion mutation.

Field
Type
Description
updatedContentCatalogVersion
ContentCatalogVersion!
The updated content catalog version information.
UpdatePlayerPayload
The response for the updatePlayer mutation.

Field
Type
Description
updatedPlayer
Player!
The updated player information.
UpdatePlayerRolePayload
The response for the updatePlayerRole mutation.

Field
Type
Description
updatedPlayerRole
PlayerRole!
The updated player role information.
UpdateSeriesPayload
The response for the updateSeries mutation.

Field
Type
Description
updatedSeries
Series!
The updated series information.
UpdateTeamPayload
The response for the updateTeam mutation.

Field
Type
Description
updatedTeam
Team!
The updated team information.
UpdateTournamentPayload
The response for the updateTournament mutation.

Field
Type
Description
updatedTournament
Tournament!
The updated tournament information.
VideoStream
Live video stream details.

Field
Type
Description
url
String!
URL link to live video stream.
Inputs
AddPlayerExternalLinkInput
Field
Type
Description
externalLink
ExternalLinkInput!
The external link to be added.
playerId
ID!
The ID of the player.
AddSeriesExternalLinkInput
Field
Type
Description
externalLink
ExternalLinkInput!
The external link to be added.
seriesId
ID!
The ID of the series.
AddTeamExternalLinkInput
Field
Type
Description
externalLink
ExternalLinkInput!
The external link to be added.
teamId
ID!
The ID of the team.
AddTournamentExternalLinkInput
Field
Type
Description
externalLink
ExternalLinkInput!
The external link to be added.
tournamentId
ID!
The ID of the tournament.
AssetInput
Field
Type
Description
assetData
String!
The logo of the team. Has to be a Base64 encoded image, supported types: image/jpg, image/jpeg, image/png, image/gif, image/svg+xml, image/webp. Maximum allowed image size is 3.5 MB
assetType
AssetType!
The type of the asset.
BooleanFilter
Boolean filter for true/false fields.

Field
Type
Description
equals
Boolean
Value to look for, can be nullable.
BoundsInput
The fields used to mutate map bounds information.

Field
Type
Description
max
CoordinatesInput!
Maximum coordinates value.
min
CoordinatesInput!
Minimum coordinates value.
ContentCatalogEntityFilter
Content catalog entity selection filter.

Field
Type
Description
entityType
ContentCatalogEntityTypeFilter
Filter by content catalog entity type.
id
IdFilter
Filter by content catalog entity id.
name
StringFilter
Filter by content catalog entity name.
private
BooleanFilter
Filter by entity visibility.
title
ContentCatalogEntityTitleFilter
Filter by content catalog entity title.
updatedAt
DateTimeFilter
Filter by last update.
ContentCatalogEntityTitleFilter
Filter by content catalog entity title.

Field
Type
Description
id
IdFilter
Filter by content catalog entity title id.
ContentCatalogEntityTypeFilter
Filter by content catalog entity type.

Field
Type
Description
in
[ContentCatalogEntityType!]!
Filter by content catalog entity type.
ContentCatalogVersionFilter
Filter by a content catalog version.

Field
Type
Description
name
StringFilter
Filter by content catalog version name.
private
BooleanFilter
Filter by version visibility.
publishedOn
DateTimeFilter
Filter by content catalog version published date.
title
ContentCatalogVersionTitleFilter
Filter by content catalog title.
ContentCatalogVersionOrder
Ordering options for content catalog versions results.

Field
Type
Description
direction
OrderDirection!
The ordering direction.
field
ContentCatalogVersionOrderField!
The field to order the content catalog entities versions by.
ContentCatalogVersionTitleFilter
Filter by a content catalog version title.

Field
Type
Description
id
IdFilter
Filter by title ID.
CoordinatesInput
The fields used to mutate spatial coordinates.

Field
Type
Description
x
Float!
X coordinate value.
y
Float!
Y coordinate value.
CreateContentCatalogCharacterInput
The fields used to create a new content catalog character or a version of an existing character.

Field
Type
Description
contentCatalogVersionId
ID!
The ID of the content catalog version the character should be linked to.
externalLinks
[ExternalLinkInput!]
The external character IDs.
id
ID
The UUID of an existing character to create a new version of it. If not passed, a new character will be created.
imageData
String
The image encoded as base64 string. If not provided, it defaults to a placeholder image.
name
String
The name of the character. Optional if creating a version of an existing character.
private
Boolean
The visibility of this entity to other users.
CreateContentCatalogItemInput
The fields used to create a new content catalog item or a version of an existing item.

Field
Type
Description
contentCatalogVersionId
ID!
The version of the content catalog the item should be linked to.
cost
Float
The cost of the item. When not provided, it defaults to 0.
externalLinks
[ExternalLinkInput!]
The external item IDs.
id
ID
The UUID of an existing item to create a new version of it.If not passed, a new item will be created.
imageData
String
The image encoded as base64 string.If not provided, it defaults to a placeholder image.
name
String
The name of the item.Optional if creating a version of an existing item.
private
Boolean
The visibility of this entity to other users.
CreateContentCatalogMapInput
The fields used to create a new content catalog map or a version of an existing map.

Field
Type
Description
bounds
BoundsInput
The map boundaries. If not specified, it defaults the coordinates to zero.
contentCatalogVersionId
ID!
The ID of the content catalog version the map should be linked to.
externalLinks
[ExternalLinkInput!]
The external map IDs.
id
ID
The UUID of an existing map to create a new version of it. If not passed, a new map will be created.
imageData
String
The image encoded as base64 string. If not provided, it defaults to a placeholder image.
name
String
The name of the map. Optional if creating a version of an existing map.
private
Boolean
The visibility of this entity to other users.
CreateContentCatalogStructureInput
The data used to create a new content catalog structure or a version of an existing structure for a specific title.

Field
Type
Description
lolContentCatalogStructure
CreateLolContentCatalogStructureInput
None
CreateContentCatalogVersionInput
The fields used to create a content catalog version.

Field
Type
Description
externalLinks
[ExternalLinkInput!]!
The external content catalog version ids.
name
String!
The UNIQUE name of the content catalog version.
private
Boolean
The visibility of this entity to other users.
publishedOn
DateTime!
The date and time when the catalog version becomes effective.
titleId
ID!
The ID of the title for the content catalog version.
CreateIngestionRequestInput
Input for creating a new Ingestion Request.

Field
Type
Description
players
[PlayerDescriptorInput!]!
The described players.
series
[SeriesDescriptorInput!]!
The described series.
title
TitleDescriptorInput!
The description of the title the ingestion request applies to.
tournaments
[TournamentDescriptorInput!]!
The described tournaments.
CreateLolContentCatalogStructureInput
The fields used to create a new content catalog structure or a version of an existing structure in LoL.

Field
Type
Description
alias
String!
The natural ID of the structure that is unique for LoL.
contentCatalogVersionId
ID!
The ID of the content catalog version the structure should be linked to.
externalLinks
[ExternalLinkInput!]
The external structure IDs.
gridType
String!
The GRID type of the structure.
id
ID
The UUID of an existing structure to create a new version of it. If not passed, a new structure will be created.
imageData
String
The image encoded as a base64 string. If not provided, it defaults to a placeholder image.
lane
String
The lane the structure is in.
maxHealth
String!
The maximum health of the structure.
name
String
The name of the structure. Optional if creating a version of an existing structure.
nexusTurretName
String
The name of the Nexus turret if the structure is one, according to its turretTier.
position
LolContentCatalogStructurePositionInput!
The position of the structure expressed with two coordinates.
private
Boolean
The visibility of this entity to other users.
respawnTime
String
The time after which the structure respawns.
riotType
String!
The type of the structure as provided by Riot Games.
teamId
Int!
The ID of the team the structure belongs to.
turretTier
String
The tier of the structure if it is a turret, according to its riotType.
CreatePlayerInput
The fields used to create a player.

Field
Type
Description
dateOfBirth
Date
The player's date of birth.
externalLinks
[ExternalLinkInput!]!
The external player IDs.
fullName
String
The player's full name.
imageData
String
The image of the player. Has to be a Base64 encoded image, supported types: image/jpg, image/jpeg, image/png, image/gif, image/svg+xml, image/webp. Maximum allowed image size is 3.5 MB
nationality
[NationalityInput!]!
The player's nationality.
nickname
String!
The nickname of the player. Must be unique and limited to 100 characters.
private
Boolean
The visibility of this entity to other users.
role
[CreatePlayerPlayerRoleInput!]!
The player's roles.
team
CreatePlayerTeamInput!
The team the player belongs to.
title
CreatePlayerTitleInput!
The title the player participates in.
CreatePlayerPlayerRoleInput
The player role field used to create a player.

Field
Type
Description
id
ID!
The ID of the player role for the player.
CreatePlayerRoleInput
The fields used to create a new player role.

Field
Type
Description
name
String!
The player role name.
private
Boolean
The visibility of this entity to other users.
title
CreatePlayerRoleTitleInput!
The player role title.
CreatePlayerRoleTitleInput
Field
Type
Description
titleId
ID!
None
CreatePlayerTeamInput
The team field used to create a player.

Field
Type
Description
teamId
ID!
The ID of the team for the player.
CreatePlayerTitleInput
The title field used to create a player.

Field
Type
Description
titleId
ID!
The ID of the title for the player.
CreateSeriesFormatInput
The format fields used to create a series.

Field
Type
Description
id
ID!
The series format identifier. Check available formats via SeriesFormat query.
CreateSeriesInput
The fields used to create a series.

Field
Type
Description
externalLinks
[ExternalLinkInput!]!
The external series IDs.
format
CreateSeriesFormatInput!
The series format i.a. best-of-3 (Bo3), score-after-2 (SA2). Referenced by ID.
players
[CreateSeriesParticipatingPlayersInput!]!
The players participating in the series being created.
private
Boolean
The visibility of this entity to other users.
startTimeScheduled
DateTime!
The start time in UTC of this series, formatted as an ISO date string. For example: 2023-04-12T12:12:37-03:00.
teams
[CreateSeriesTeamInput!]!
The teams playing in the series.
title
CreateSeriesTitleInput!
The title for the series.
tournament
CreateSeriesTournamentInput!
The tournament for the series.
CreateSeriesParticipatingPlayersInput
A player participating in the series being created.

Field
Type
Description
id
ID!
The player ID.
CreateSeriesTeamInput
The team fields used to create a series.

Field
Type
Description
scoreAdvantage
Int!
The score advantage for the team in this series.
teamId
ID!
The ID of the team playing in the series.
CreateSeriesTitleInput
The title fields used to create a series.

Field
Type
Description
titleId
ID!
The ID of the title for the series.
CreateSeriesTournamentInput
The tournament fields used to create a series.

Field
Type
Description
tournamentId
ID!
The ID of the tournament for the series.
CreateTeamInput
The fields used to create a team.

Field
Type
Description
assets
[AssetInput!]!
Assets for the team.
colorPrimary
HexColor!
The team's primary color in hexadecimal. Defaults to grey.
colorSecondary
HexColor!
The team's secondary color in hexadecimal. Defaults to white.
externalLinks
[ExternalLinkInput!]!
The external team IDs.
name
String!
The name of the team. Must be unique and limited to 100 characters.
nameShortened
String
The shortened name of the team. It's limited to 20 characters and it's not unique.
private
Boolean
The visibility of this entity to other users.
title
CreateTeamTitleInput!
The title the team participates in.
CreateTeamTitleInput
The title fields used to create a team.

Field
Type
Description
titleId
ID!
The ID of the title for the team.
CreateTournamentInput
The fields used to create a tournament.

Field
Type
Description
endDate
DateInput
The end date of the tournament formatted as ISO 8601 i.e. yyyy-MM-dd
externalLinks
[ExternalLinkInput!]!
The external tournament IDs.
logoData
String
The logo of the tournament. Has to be a Base64 encoded image, supported types: image/jpg, image/jpeg, image/png, image/gif, image/svg+xml, image/webp. Maximum allowed image size is 3.5 MB
name
String!
The name of the tournament. Must be unique and limited to 150 characters.
nameShortened
String!
The short name of the tournament. Must be unique and limited to 30 characters.
parent
CreateTournamentParentInput
The parent tournament of the tournament
private
Boolean
The visibility of this entity to other users.
prizePool
MoneyInput
The prize pool of the tournament in USD.
startDate
DateInput
The start date of the tournament formatted as ISO 8601 i.e. yyyy-MM-dd
teams
[CreateTournamentTeamInput!]!
The participating teams of the tournament.
titles
[CreateTournamentTitleInput!]!
The titles being played in the tournament.
tournamentType
TournamentType!
The type of the tournament. Defaults to UNSPECIFIED.
venueType
TournamentVenueType!
The venue type of the tournament. Defaults to UNKNOWN.
CreateTournamentParentInput
The parent tournament fields used to create a tournament together with a parent.

Field
Type
Description
tournamentId
ID!
The ID of the parent tournament
CreateTournamentTeamInput
The team fields used to update a tournament.

Field
Type
Description
teamId
ID!
The ID of the team for the tournament.
CreateTournamentTitleInput
The title fields used to create a tournament.

Field
Type
Description
titleId
ID!
The ID of the title for the tournament.
DateInput
Field
Type
Description
date
Date
ISO 8601 formatted date i.e. yyyy-MM-dd
DateTimeFilter
Filter by a datetime range.

Field
Type
Description
gte
String
DateTime for the beginning of the time range (optional), formatted as an ISO date string. Time window includes the given value.
lte
String
DateTime for the end of the time range (optional), formatted as an ISO date string. Time window includes the given value.
DeletePlayerInput
The ID of the player to delete.

Field
Type
Description
id
ID!
The player ID.
DeletePlayerRoleInput
The fields used to delete a player role.

Field
Type
Description
id
ID!
The ID of the player role to be deleted.
DeleteSeriesInput
The ID of the series to delete.

Field
Type
Description
id
ID!
The series ID.
DeleteTeamInput
The ID of the team to delete.

Field
Type
Description
id
ID!
The team ID.
DeleteTournamentInput
The ID of the tournament to delete.

Field
Type
Description
id
ID!
The tournament ID.
ExternalEntityInput
The field for the external ID value.

Field
Type
Description
id
ID!
The external ID value.
ExternalLinkInput
The fields the external entity consist of.

Field
Type
Description
dataProviderName
String!
The data provider name. Strict control over possible values.
externalEntityInput
ExternalEntityInput!
The external id value.
GameFilter
Filter for series by game data points.

Field
Type
Description
id
IdFilter
Filter by game IDs.
map
MapFilter
Filter by game map.
IdFilter
ID filter for entity identifiers.

Field
Type
Description
in
[ID!]
Array of IDs to look for.
IngestionRequestFilter
Filter for ingestion request data points.

Field
Type
Description
status
IngestionRequestStatusFilter
Filter by status.
IngestionRequestStatusFilter
Field
Type
Description
in
[IngestionRequestStatus!]
None
IntFilter
Filter by integer value or range.

Field
Type
Description
equals
Int
Filter by value equal.
gte
Int
Filter by value greater than or equal.
lte
Int
Filter by value less than or equal.
LolContentCatalogStructurePositionInput
The fields used to mutate position data of a LoL structure.

Field
Type
Description
x
Float!
None
z
Float!
None
MapFilter
Filter for map data points.

Field
Type
Description
name
StringFilter
Filter by map name.
MoneyInput
A monetary value.

Field
Type
Description
amount
Decimal
The amount of the prize pool.
NationalityDescriptorInput
A description of a player's nationality. At least one field must be set.

Field
Type
Description
code
String
The ISO 3166-1 Alpha-3 code for the country (e.g., 'USA' for the United States).
name
String
The full name of the country (e.g., 'United States').
NationalityFilter
Filter by player nationality.

Field
Type
Description
code
StringFilter
Filter by nationality code.
name
StringFilter
Filter by nationality name.
NationalityInput
Field
Type
Description
nationalityCode
String!
None
NullableDateFilter
Filter by a date range. If both 'gte' and 'lte' are null, it'll fetch the entities for which the field the filter is applied to is null.

Field
Type
Description
gte
Date
Date for the beginning of the time range (optional), formatted as an ISO date string. Time window includes the given value.
lte
Date
Date for the end of the time range (optional), formatted as an ISO date string. Time window includes the given value.
NullableIdFilter
ID filter with nullable value.

Field
Type
Description
id
ID
ID to look for, can be nullable.
OrganizationFilter
Filter for organization by data points.

Field
Type
Description
name
StringFilter
Limits results by string filter of organization's name.
private
BooleanFilter
Filter by organization visibility.
updatedAt
DateTimeFilter
Filter by last update.
PlayerDescriptorInput
A description of a player.

Field
Type
Description
dataSource
String
The data source that sourceId belongs to. Must be set together with sourceId.
dateOfBirth
Date
The described player's date of birth.
fullName
String
The described player's full name.
id
ID
GRID ID of the described player.
nationality
[NationalityDescriptorInput!]
The description of the nationality of the described player.
nickname
String!
The name of the described player.
roles
[PlayerRoleDescriptorInput!]
The description of the roles (e.g., MID, SUPPORT) of the described player.
sourceId
ID
The player identifier on the source platform.
team
TeamDescriptorInput!
The description of a team where the described player currently plays.
PlayerFilter
Filter for player data points.

Field
Type
Description
age
IntFilter
Filter by age.
fullName
StringFilter
Filter by full name.
nationality
NationalityFilter
Filter by nationality.
nickname
StringFilter
Filter by the player's nickname.
private
BooleanFilter
Filter by player visibility.
roles
PlayerPlayerRoleFilter
Filter players by role.
teamIdFilter
NullableIdFilter
Filter by a specific team ID. If this is null, no filter will be applied. If teamId.id is null, then players without team will be fetched.
titleId
ID
Filter by a specific title ID.
types
[PlayerType!]
Filter by player type.
updatedAt
DateTimeFilter
Filter by last update.
PlayerPlayerRoleFilter
Filter by player role.

Field
Type
Description
name
StringFilter
Filter by player role name.
title
PlayerRoleTitleFilter
Filter by player role title.
PlayerRoleDescriptorInput
A description of player's role (e.g., MID, SUPPORT).

Field
Type
Description
name
String!
The name of the described player role.
PlayerRoleFilter
Filter by player role.

Field
Type
Description
name
StringFilter
Filter by player role name.
title
PlayerRoleTitleFilter
Filter by player role title.
PlayerRoleTitleFilter
Filter player roles by title.

Field
Type
Description
id
IdFilter!
Filter player roles by title ID.
ProductServiceLevelFilter
Filters for series by product and/or by service level.

Field
Type
Description
productName
String!
The name of the product to filter by.
serviceLevel
ServiceLevel!
The availability level to filter by.
SeriesDescriptorInput
A description of a series.

Field
Type
Description
dataSource
String
The data source that sourceId belongs to. Must be set together with sourceId.
format
SeriesFormatDescriptorInput!
The format of the described series.
id
ID
GRID ID of the described series.
players
[PlayerDescriptorInput!]
The players participating in the described series.
sourceId
ID
The series identifier on the source platform.
startTimeScheduled
DateTime!
The scheduled start time of the described series.
teams
[TeamDescriptorInput!]!
The teams participating in the described series.
tournament
TournamentDescriptorInput!
The description of a tournament in which the described series is being played.
SeriesFilter
Filter for series data points.

Field
Type
Description
live
SeriesLiveFilter
Filter by live data points.
livePlayerIds
IdFilter
Filter by participating player IDs.
players
SeriesPlayerFilter
Filter by participating players data points.
private
BooleanFilter
Filter by series visibility.
productServiceLevels
ProductServiceLevelFilter
Filter by product service levels.
startTimeScheduled
DateTimeFilter
Filter by a time range.
teamId
ID
Use teamIds filter instead. Filter by team ID.
teamIds
IdFilter
Filter by team IDs.
titleId
ID
Use titleIds filter instead. Filter by title ID.
titleIds
IdFilter
Filter by title IDs.
tournament
SeriesTournamentFilter
Filter by tournament data points.
tournamentId
ID
Use the tournament filter instead. Filter by tournament ID.
tournamentIds
IdFilter
Use the tournament filter instead. Filter by tournament IDs.
type
SeriesType
Use types filter instead. Filter by series type.
types
[SeriesType!]
Filter by series types.
updatedAt
DateTimeFilter
Filter by last update.
workflowStatuses
[SeriesWorkflowStatus!]!
Filter by series workflow statuses.
SeriesFormatDescriptorInput
The format of the series (e.g., Bo1, Bo3, SA2).

Field
Type
Description
nameShortened
String!
The shortened name of the described format.
SeriesLiveFilter
Filter for series by live data points.

Field
Type
Description
games
GameFilter
Filter by game.
SeriesPlayerFilter
Filter series by participating players data points.

Field
Type
Description
age
IntFilter
Filter series where at least one participant's age meets the specified criteria.
SeriesTournamentFilter
Filter for series by tournament data points.

Field
Type
Description
endDate
NullableDateFilter
Filter by tournament end date.
id
IdFilter
Filter by tournament IDs.
includeChildren
BooleanFilter
Include children tournaments in the filter.
name
StringFilter
Filter by tournament name.
nameShortened
StringFilter
Filter by tournament shortened name.
startDate
NullableDateFilter
Filter by tournament start date.
title
SeriesTournamentTitleFilter
Filter by tournament titles data points.
SeriesTournamentTitleFilter
Filter for series by the tournament title data points.

Field
Type
Description
id
IdFilter
Array of IDs to look for.
StringFilter
String filter for character-based fields.

Field
Type
Description
contains
String
Case-insensitive string filter to match field containing it.
equals
String
Case-insensitive string filter to match field exactly it.
TeamDescriptorInput
A description of a team.

Field
Type
Description
dataSource
String
The data source that sourceId belongs to. Must be set together with sourceId.
id
ID
GRID ID of the described team.
name
String!
The name of the described team.
nameShortened
String
The shortened name of the described team.
sourceId
ID
The team identifier on the source platform.
TeamFilter
Filter for team data points.

Field
Type
Description
name
StringFilter
Filter by the team's name.
nameShortened
StringFilter
Filter by the team's shortened name.
organizationId
ID
Filter by a specific organization.
private
BooleanFilter
Filter by team visibility.
titleId
ID
Filter by a specific title.
updatedAt
DateTimeFilter
Filter by last update.
workflowStatuses
[TeamWorkflowStatus!]!
Filter by team workflow statuses.
TitleDescriptorInput
A description of a title.

Field
Type
Description
name
String!
The name of the described title.
TitleFilter
Filter for titles.

Field
Type
Description
private
BooleanFilter
Filter by title visibility.
TournamentDescriptorInput
A description of a tournament.

Field
Type
Description
dataSource
String
The data source that sourceId belongs to. Must be set together with sourceId.
endDate
Date
The end date of the described tournament.
id
ID
GRID ID of the described tournament.
name
String!
The name of the described tournament.
nameShortened
String
The shortened name of the described tournament (max. 30 characters).
parent
TournamentDescriptorInput
The description of the parent tournament of the described tournament.
prizePool
MoneyInput
The prize pool of the described tournament in USD.
sourceId
ID
The tournament identifier on the source platform.
startDate
Date
The start date of the described tournament.
venueType
TournamentVenueType
The type of venue where the described tournament is held.
TournamentFilter
Filter for tournament data points.

Field
Type
Description
endDate
NullableDateFilter
Filter by tournament end date.
hasChildren
BooleanFilter
Filter tournaments that have children. If false, only tournaments at the bottom-most level of the hierarchy are returned.
hasParent
BooleanFilter
Filter tournaments that have parents. If false, only tournaments at the top-most level of the hierarchy are returned.
name
StringFilter
Filter by tournament name.
nameShortened
StringFilter
Filter by tournament short name.
private
BooleanFilter
Filter by tournaments visibility.
productServiceLevels
ProductServiceLevelFilter
Filter tournaments for which at least one series has the specified product service level.
startDate
NullableDateFilter
Filter by tournament start date.
title
TournamentTitleFilter
Filter by tournament titles.
titleId
ID
Use the titles filter instead. Limits results to specific title.
tournamentType
[TournamentType!]
Filters tournaments by tournament type.
updatedAt
DateTimeFilter
Filter by last update.
venueType
[TournamentVenueType!]
Filters tournaments by venue type.
TournamentTitleFilter
Filter for tournaments by title data points.

Field
Type
Description
id
IdFilter
Filter by tournament IDs.
UpdateContentCatalogCharacterInput
The fields used to update a new content catalog character or a version of an existing character.

Field
Type
Description
contentCatalogVersionId
ID!
The content catalog version.
externalLinks
[ExternalLinkInput!]
The external character IDs. Overwrites all the previous values.
id
ID!
The UUID of an existing character.
imageData
String
The image encoded as base64 string. If not provided, it defaults to a placeholder image.
name
String
The new name of the character.
private
Boolean
The visibility of this entity to other users.
UpdateContentCatalogItemInput
The fields used to create a new content catalog item or a version of an existing item.

Field
Type
Description
contentCatalogVersionId
ID!
The content catalog version.
cost
Float
The new cost of the item.
externalLinks
[ExternalLinkInput!]
The external item IDs. Overrides all previous values.
id
ID!
The UUID of an existing item.
imageData
String
The new base64 encoded image of the item.
name
String
The new name of the item.
private
Boolean
The visibility of this entity to other users.
UpdateContentCatalogMapInput
The fields used to update a new content catalog map or a version of an existing map.

Field
Type
Description
bounds
BoundsInput
The new map boundaries.
contentCatalogVersionId
ID!
The content catalog version.
externalLinks
[ExternalLinkInput!]
The external map IDs. Overwrites all the previous values.
id
ID!
The UUID of an existing map.
imageData
String
The image encoded as base64 string. If not provided, it defaults to a placeholder image.
name
String
The new name of the map.
private
Boolean
The visibility of this entity to other users.
UpdateContentCatalogStructureInput
The data used to update a content catalog structure for a specific title.

Field
Type
Description
lolContentCatalogStructure
UpdateLolContentCatalogStructureInput
None
UpdateContentCatalogVersionInput
The fields used to update a content catalog version.

Field
Type
Description
contentCatalogVersionId
ID!
The unique identifier of the content catalog version.
private
Boolean
The visibility of this entity to other users.
UpdateLolContentCatalogStructureInput
The fields used to update a content catalog structure in LoL.

Field
Type
Description
alias
String!
The natural ID of the structure that is unique for LoL.
contentCatalogVersionId
ID!
The content catalog version.
externalLinks
[ExternalLinkInput!]
The external structure IDs. Overwrites all the previous values.
gridType
String!
The GRID type of the structure.
id
ID!
The UUID of an existing structure.
imageData
String
The image encoded as a base64 string. If not provided, it defaults to a placeholder image.
lane
String
The lane the structure is in.
maxHealth
String!
The maximum health of the structure.
name
String
The new name of the structure.
nexusTurretName
String
The name of the Nexus turret if the structure is one, according to its turretTier.
position
LolContentCatalogStructurePositionInput!
The position of the structure expressed with two coordinates.
private
Boolean
The visibility of this entity to other users.
respawnTime
String
The time after which the structure respawns.
riotType
String!
The type of the structure as provided by Riot Games.
teamId
Int!
The ID of the team the structure belongs to.
turretTier
String
The tier of the structure if it is a turret, according to its riotType.
UpdatePlayerInput
The fields used to update a player.

Field
Type
Description
dateOfBirth
Date
The player's date of birth.
externalLinks
[ExternalLinkInput!]
The external player IDs. Overrides all previous values.
fullName
String
The player's full name.
id
ID!
The ID of the player.
imageData
String
The image of the player. Has to be a Base64 encoded image, supported types: image/jpg, image/jpeg, image/png, image/gif, image/svg+xml, image/webp. Maximum allowed image size is 3.5 MB
nationality
[NationalityInput!]!
The player's nationality.
nickname
String
The nickname of the player. Must be unique and limited to 100 characters.
private
Boolean
The visibility of this entity to other users.
role
[UpdatePlayerPlayerRoleInput!]
The player's roles.
team
UpdatePlayerTeamInput
The team the player belongs to. Overrides previous value.
title
UpdatePlayerTitleInput
The title the player participates in. Overrides previous value.
UpdatePlayerPlayerRoleInput
The player role field used to update a player.

Field
Type
Description
id
ID!
The ID of the player role for the player.
UpdatePlayerRoleInput
The fields used to update a player role.

Field
Type
Description
id
ID!
The ID of the player role to be updated.
name
String
The new role name.
private
Boolean
The visibility of this entity to other users.
UpdatePlayerRoleTitleInput
Field
Type
Description
titleId
ID!
None
UpdatePlayerTeamInput
The team field used to update a player.

Field
Type
Description
teamId
ID
The ID of the team for the player.
UpdatePlayerTitleInput
The title field used to update a player.

Field
Type
Description
titleId
ID!
The ID of the title for the player.
UpdateSeriesFormatInput
The format fields used to update a series.

Field
Type
Description
id
ID!
The series format identifier. Check available formats via SeriesFormat query.
UpdateSeriesInput
The fields used to update a series.

Field
Type
Description
externalLinks
[ExternalLinkInput!]
The external series IDs. Overrides all previous values.
format
UpdateSeriesFormatInput
The new series format i.a. best-of-3 (Bo3), score-after-2 (SA2). Referenced by ID.
id
ID!
The ID of the series.
players
[UpdateSeriesParticipatingPlayersInput!]
The players participating in the series being updated.
private
Boolean
The visibility of this entity to other users.
startTimeScheduled
DateTime
The new start time in UTC of this series, formatted as an ISO date string. For example: 2023-04-12T12:12:37-03:00.
teams
[UpdateSeriesTeamInput!]
The new teams for the series. Overrides the previous value and is only possible if at least one series team is not yet defined.
title
UpdateSeriesTitleInput
The new title for the series.
tournament
UpdateSeriesTournamentInput
The new tournament for the series.
UpdateSeriesParticipatingPlayersInput
A player participating in the series being updated.

Field
Type
Description
id
ID!
The player ID.
UpdateSeriesTeamInput
The team fields used to update a series.

Field
Type
Description
scoreAdvantage
Int!
The score advantage for the team in this series.
teamId
ID!
The ID of the team playing in the series.
UpdateSeriesTitleInput
The title fields used to update a series.

Field
Type
Description
titleId
ID!
The ID of the title for the series.
UpdateSeriesTournamentInput
The tournament fields used to update a series.

Field
Type
Description
tournamentId
ID!
The ID of the tournament for the series.
UpdateTeamInput
The fields used to update a team.

Field
Type
Description
assets
[AssetInput!]
Replaces all assets logo variants. Omit to leave unchanged. Pass empty list to delete all.
colorPrimary
HexColor
The new team's primary color in hexadecimal. Defaults to grey.
colorSecondary
HexColor
The new team's secondary color in hexadecimal. Defaults to white.
externalLinks
[ExternalLinkInput!]
The external team IDs. Overrides all previous values.
id
ID!
The ID of the team.
name
String
The new name of the team. Must be unique and limited to 100 characters.
nameShortened
String
The new shortened name of the team. It's limited to 20 characters and it's not unique.
private
Boolean
The visibility of this entity to other users.
title
UpdateTeamTitleInput
The title the team participates in. Overrides previous value.
UpdateTeamTitleInput
The title fields used to update a team.

Field
Type
Description
titleId
ID!
The ID of the title for the team.
UpdateTournamentInput
The fields used to update a tournament.

Field
Type
Description
endDate
DateInput
The end date of the tournament formatted as ISO 8601 i.e. yyyy-MM-dd
externalLinks
[ExternalLinkInput!]
The external tournament IDs. Overrides all previous values.
id
ID!
The ID of the tournament.
logoData
String
The new logo of the tournament. Has to be a Base64 encoded image, supported types: image/jpg, image/jpeg, image/png, image/gif, image/svg+xml, image/webp. Maximum allowed image size is 3.5 MB
name
String
The new name of the tournament. Must be unique and limited to 150 characters.
nameShortened
String
The new short name of the tournament. Must be unique and limited to 30 characters.
parent
UpdateTournamentParentInput
The parent tournament of the tournament
private
Boolean
The visibility of this entity to other users.
prizePool
MoneyInput
The prize pool of the tournament in USD.
startDate
DateInput
The start date of the tournament formatted as ISO 8601 i.e. yyyy-MM-dd
teams
[UpdateTournamentTeamInput!]
The new teams that will be playing in the tournament. Overrides previous value.
titles
[UpdateTournamentTitleInput!]
The new titles being played in the tournament. Overrides previous value.
tournamentType
TournamentType
The type of the tournament.
venueType
TournamentVenueType
The venue type of the tournament.
UpdateTournamentParentInput
The parent tournament fields used to create a tournament together with a parent.

Field
Type
Description
tournamentId
ID
None
UpdateTournamentTeamInput
The team fields used to update a tournament.

Field
Type
Description
teamId
ID!
The ID of the team for the tournament.
UpdateTournamentTitleInput
The title fields used to update a tournament.

Field
Type
Description
titleId
ID!
The ID of the title for the tournament.
Interfaces
ContentCatalogEntity
A content catalog entity.

Field
Type
Description
contentCatalogVersion
ContentCatalogVersion!
The content catalog version this entity appears in.
id
ID!
The ID of the content catalog entity.
imageUrl
Url!
The URL to the content catalog entity's image
name
String!
The name of the content catalog entity.
private
Boolean
The visibility of this entity to other users.
updatedAt
DateTime!
The date and time of the last update.
ContentCatalogStructure
A content catalog structure.

Field
Type
Description
contentCatalogVersion
ContentCatalogVersion!
The content catalog version this entity appears in.
externalLinks
[ExternalLink!]!
The content catalog entity version external IDs and providers.
id
ID!
The ID of a content catalog entity.
imageUrl
Url!
The URL to the content catalog entity's image
name
String!
The name of a content catalog entity.
private
Boolean
The visibility of this entity to other users.
updatedAt
DateTime!
The date and time of the last update.
OrganizationInterface
An interface for organization types.

Field
Type
Description
id
ID!
The ID of the organization.
name
String!
The name of the organization.
TeamInterface
An interface for team types.

Field
Type
Description
colorPrimary
HexColor!
The primary color of the team.
colorSecondary
HexColor!
The secondary color of the team.
id
ID!
The ID of the team.
logoUrl
Url!
The team logo URL.
name
String!
The name of the team.
nameShortened
String
The shortened name of a team.
rating
Float
The team rating.
titles
[Title!]!
The titles the team is related to.
Enums
AssetType
Field
Description
ALTERNATIVE_LOGO
-
BACKGROUND_LOGO
-
DARK_LOGO
-
LIGHT_LOGO
-
MAIN_LOGO
-
ContentCatalogEntityType
The type of a content catalog entity.

Field
Description
CHARACTER
Content catalog character.
ITEM
Content catalog item.
MAP
Content catalog map.
STRUCTURE
Content catalog structure.
ContentCatalogVersionOrderField
The field which the result will be ordered by.

Field
Description
PUBLISHED_ON
Order by publishedOn.
TITLE
Order by title name.
ErrorDetail
Field
Description
DEADLINE_EXCEEDED
The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long enough for the deadline to expire. HTTP Mapping: 504 Gateway Timeout Error Type: UNAVAILABLE
ENHANCE_YOUR_CALM
The server detected that the client is exhibiting a behavior that might be generating excessive load. HTTP Mapping: 420 Enhance Your Calm Error Type: UNAVAILABLE
FIELD_NOT_FOUND
The requested field is not found in the schema. This differs from NOT_FOUND in that NOT_FOUND should be used when a query is valid, but is unable to return a result (if, for example, a specific video id doesn't exist). FIELD_NOT_FOUND is intended to be returned by the server to signify that the requested field is not known to exist. This may be returned in lieu of failing the entire query. See also PERMISSION_DENIED for cases where the requested field is invalid only for the given user or class of users. HTTP Mapping: 404 Not Found Error Type: BAD_REQUEST
INVALID_ARGUMENT
The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name). HTTP Mapping: 400 Bad Request Error Type: BAD_REQUEST
INVALID_CURSOR
The provided cursor is not valid. The most common usage for this error is when a client is paginating through a list that uses stateful cursors. In that case, the provided cursor may be expired. HTTP Mapping: 404 Not Found Error Type: NOT_FOUND
MISSING_RESOURCE
Unable to perform operation because a required resource is missing. Example: Client is attempting to refresh a list, but the specified list is expired. This requires an action by the client to get a new list. If the user is simply trying GET a resource that is not found, use the NOT_FOUND error type. FAILED_PRECONDITION.MISSING_RESOURCE is to be used particularly when the user is performing an operation that requires a particular resource to exist. HTTP Mapping: 400 Bad Request or 500 Internal Server Error Error Type: FAILED_PRECONDITION
SERVICE_ERROR
Service Error. There is a problem with an upstream service. This may be returned if a gateway receives an unknown error from a service or if a service is unreachable. If a request times out which waiting on a response from a service, DEADLINE_EXCEEDED may be returned instead. If a service returns a more specific error Type, the specific error Type may be returned instead. HTTP Mapping: 502 Bad Gateway Error Type: UNAVAILABLE
TCP_FAILURE
Request failed due to network errors. HTTP Mapping: 503 Unavailable Error Type: UNAVAILABLE
THROTTLED_CONCURRENCY
Request throttled based on server concurrency limits. HTTP Mapping: 503 Unavailable Error Type: UNAVAILABLE
THROTTLED_CPU
Request throttled based on server CPU limits HTTP Mapping: 503 Unavailable. Error Type: UNAVAILABLE
TOO_MANY_REQUESTS
The server detected that the client is exhibiting a behavior that might be generating excessive load. HTTP Mapping: 429 Too Many Requests Error Type: UNAVAILABLE
UNIMPLEMENTED
The operation is not implemented or is not currently supported/enabled. HTTP Mapping: 501 Not Implemented Error Type: BAD_REQUEST
UNKNOWN
Unknown error. This error should only be returned when no other error detail applies. If a client sees an unknown errorDetail, it will be interpreted as UNKNOWN. HTTP Mapping: 500 Internal Server Error
ErrorType
Field
Description
BAD_REQUEST
Bad Request. There is a problem with the request. Retrying the same request is not likely to succeed. An example would be a query or argument that cannot be deserialized. HTTP Mapping: 400 Bad Request
FAILED_PRECONDITION
The operation was rejected because the system is not in a state required for the operation's execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementers can use the following guidelines to decide between FAILED_PRECONDITION and UNAVAILABLE: - Use UNAVAILABLE if the client can retry just the failing call. - Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory. HTTP Mapping: 400 Bad Request or 500 Internal Server Error
INTERNAL
Internal error. An unexpected internal error was encountered. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors. HTTP Mapping: 500 Internal Server Error
NOT_FOUND
The requested entity was not found. This could apply to a resource that has never existed (e.g. bad resource id), or a resource that no longer exists (e.g. cache expired.) Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, NOT_FOUND may be used. If a request is denied for some users within a class of users, such as user-based access control, PERMISSION_DENIED must be used. HTTP Mapping: 404 Not Found
PERMISSION_DENIED
The caller does not have permission to execute the specified operation. PERMISSION_DENIED must not be used for rejections caused by exhausting some resource or quota. PERMISSION_DENIED must not be used if the caller cannot be identified (use UNAUTHENTICATED instead for those errors). This error Type does not imply the request is valid or the requested entity exists or satisfies other pre-conditions. HTTP Mapping: 403 Forbidden
UNAUTHENTICATED
The request does not have valid authentication credentials. This is intended to be returned only for routes that require authentication. HTTP Mapping: 401 Unauthorized
UNAVAILABLE
Currently Unavailable. The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. HTTP Mapping: 503 Unavailable
UNKNOWN
Unknown error. For example, this error may be returned when an error code received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error. If a client sees an unknown errorType, it will be interpreted as UNKNOWN. Unknown errors MUST NOT trigger any special behavior. These MAY be treated by an implementation as being equivalent to INTERNAL. When possible, a more specific error should be provided. HTTP Mapping: 520 Unknown Error
IngestionRequestOrderBy
Field to order results by.

Field
Description
CreatedOn
Orders the IngestionRequest results by their creation time.
ID
Orders the IngestionRequest results by their IDs.
IngestionRequestStatus
The status of an ingestion request.

Field
Description
COMPLETED
The ingestion request processing has completed successfully.
FAILED
The ingestion request processing failed due to a system error.
IN_PROGRESS
The ingestion request is being processed.
PENDING
The ingestion request is yet to be processed.
REJECTED
The ingestion request has been rejected by a reviewer.
REQUIRES_REVIEW
The ingestion request includes one or more descriptors which require review.
OrderDirection
Order in which the results will be provided.

Field
Description
ASC
Orders the results ascending.
DESC
Orders the results descending.
PlayerType
The type of player being returned.

Field
Description
ESPORTS
An esports player.
REGULAR
A non-esports player.
SeriesOrderBy
Field to order results by.

Field
Description
ID
Orders the series result by ID.
StartTimeScheduled
Orders the series result by the scheduled start time.
UpdatedAt
Orders the series result by the updated time.
SeriesType
The type of series being returned.

Field
Description
COMPETITIVE
A competitive non-esports series.
ESPORTS
A competitive series.
LOOPFEED
A GRID test game that gets played repeatedly in a loop.
SCRIM
A practice series.
SeriesWorkflowStatus
The workflow status of the series.

Field
Description
DELETED
The series has been deleted.
PUBLISHED
The series is created.
ServiceLevel
The availability level for a given product.

Field
Description
FULL
The product service is fully available.
LIMITED
Some functionalities of the product service may be limited or unavailable.
NONE
The product service is not available.
TeamWorkflowStatus
The workflow status of the team.

Field
Description
DELETED
The team has been deleted.
PUBLISHED
The team is created.
TournamentType
Defines the type of a tournament within a tournament tree.

Field
Description
LEAGUE
The tournament is a league. A league sits at the top of a tournament tree.
UNSPECIFIED
The tournament type is not specified.
TournamentVenueType
Defines the venue type of a tournament.

Field
Description
HYBRID
The tournament has both online and LAN phases.
LAN
The tournament is held entirely at a physical location (LAN).
ONLINE
The tournament is held entirely online.
UNKNOWN
The venue type is not specified.
Scalars
Boolean
The Boolean scalar type represents true or false.

Date
Date formatted as ISO 8601.

DateTime
DateTime formatted as ISO 8601.

Decimal
Decimal number representation.

Float
The Float scalar type represents signed double-precision fractional values as specified by IEEE 754.

HexColor
Hexadecimal value of an RGB color.

ID
The ID scalar type represents a unique identifier, often used to refetch an object or as key for a cache. The ID type appears in a JSON response as a String; however, it is not intended to be human-readable. When expected as an input type, any string (such as "4") or integer (such as 4) input value will be accepted as an ID.

Int
The Int scalar type represents non-fractional signed whole numeric values. Int can represent values between -(2^31) and 2^31 - 1.

Locale
A BCP 47 language tag, e.g. 'en', 'en-GB', 'sv'.

String
The String scalar type represents textual data, represented as UTF-8 character sequences. The String type is most often used by GraphQL to represent free-form human-readable text.

Url
URL pointing to resource.



Series Events API
Transactions
A transaction is a single package of data sent over the WebSocket. It contains metadata (such as timestamps and sequence numbers) and an array of events. Transactions guarantee that all included events are processed together in the correct order, preserving the integrity of the series timeline.

Field Name
Type
Description
id
String
Unique ID of the transaction
correlationId
String
Internal reference for troubleshooting
occurredAt
String
ISO 8601 timestamp when the transaction occurred
publishedAt
String
ISO 8601 timestamp when the transaction was published to the WebSocket
seriesId
String
ID of the series
sequenceNumber
Int
Order of the transaction in the feed for the series (may skip if filtered)
sessionSequenceNumber
Int
Order of the transaction in the context of the current consumer session
events
Event[]
Array of one or more event objects (see below)
Events
An event is a single, discrete change or action that occurred within the series, such as a player kill, round start, or objective completion. Events include detailed information about the actor, action, target, and the resulting state change.

Field Name
Type
Description
id
String
Unique ID of the event
type
String
Event type in the format actor-action-target (e.g., player-killed-player)
includesFullState
Boolean
Whether the event includes the full state (based on configuration)
actor
Entity
The entity that performed the action
action
String
The action performed (e.g., killed, planted)
target
Entity
The entity affected by the action
seriesStateDelta
Object
The fields that changed at the series level as a result of this event
seriesState
Object
The complete series state after the event (if included)
revertedEventId
String?
(Optional) ID of the event that this event is reverting (only present for reverted events)
Events are structured as {actor}-{action}-{target} (e.g., player-killed-player). The API supports both shared (cross-title) and title-specific events.

Reverted Events
Any event can be reverted if an unexpected issue occurred (for example, due to an admin reset or technical reasons). A reverted event has the same structure as a normal event, but the action is prefixed with an exclamation mark (!), e.g., player-!killed-player.
The revertedEventId field points to the original event being reverted. The delta values in the reverted event are the inverse of the original (e.g., a kill count of -1 instead of +1).

Entity Fields
Field Name
Type
Description
type
String
Type of the entity (e.g., player, team, game)
id
String
Unique series-level entity ID
stateDelta
Object
Fields changed at the entity level
state
Object
Entity state after the event
Entity states are usually nested in one another, as can be seen here:

Series, which consists of
Games, which consist of
Segments (e.g., round), which consist of
Segments (e.g., buy phase) and so on
Supported Event Types
The following tables provide a comprehensive overview of all Series Events API event types, clearly indicating which are shared across all titles and which are title-specific.
For clarity, events are grouped by their actor and action, with title-specific events marked accordingly.

Events Available Across All Titles
Actor
Action
Target
Description
Titles
grid
started
feed
Feed created (first event)
All
grid
ended
feed
Feed ended (last event)
All
grid
invalidated
series
Series state cannot be confirmed
All
grid
validated
series
Series state is correct again
All
tournament
started
series
Series begins as part of a tournament
All
tournament
ended
series
Series in a tournament finished
All
series
started
game
Game started in a series
All
series
ended
game
Game ended in a series
All
team
won
game
Team won a game
All
team
won
series
Team won a series
All
player
killed
player
Player killed an opponent
All
player
teamkilled
player
Player killed a teammate
All
player
selfkilled
player
Player killed themselves
All
player
rejoined
series
Player rejoined after being inactive
All
player
left
series
Player left the series
All
player
revived
player
Player revived a teammate
All
game
set
clock
Game clock set anew
All
game
started
clock
Game clock started
All
game
stopped
clock
Game clock stopped
All
Title-Specific Events
Actor
Action
Target
Description
Titles
freezetime
ended
timeout
Timeout ended during freezetime
CS2, VALORANT
freezetime
started
timeout
Timeout started during freezetime
CS2, VALORANT
game
ended
round
Round ended
CS2, VALORANT, Rainbow Six
game
killed
ATierNPC
A Tier NPC killed a player
LoL
game
killed
player
Game killed a player (environmental death)
CS2, DOTA, VALORANT, Rainbow Six
game
killed
STierNPC
S Tier NPC killed a player
LoL
game
respawned
ATierNPC
A Tier NPC respawned
LoL
game
respawned
player
Player respawned
DOTA, LoL
game
respawned
roshan
Roshan respawned
DOTA
game
respawned
STierNPC
S Tier NPC respawned
LoL
game
set
gameClock
Game clock set anew
DOTA, VALORANT
game
set
npcRespawnClock
NPC respawn clock set
DOTA, LoL
game
started
npcRespawnClock
NPC respawn clock started
DOTA, LoL
game
started
round
Round started
CS2, VALORANT, Rainbow Six
game
stopped
npcRespawnClock
NPC respawn clock stopped
DOTA, LoL
player
!purchased
item
Player undid a purchase
LoL
player
captured
outpost
Player captured an outpost
DOTA
player
completed
beginDefuseBomb
Player started defusing the Spike
VALORANT
player
completed
beginDefuseWithKit
Player started defusing bomb with kit
CS2
player
completed
beginDefuseWithoutKit
Player started defusing bomb without kit
CS2
player
completed
captureUltimateOrb
Player captured Ultimate Orb
VALORANT
player
completed
defuseBomb
Player defused the Spike
VALORANT
player
completed
destroyAncient
Player destroyed the Ancient
DOTA
player
completed
destroyBarracksMelee
Player destroyed Melee Barracks
DOTA
player
completed
destroyBarracksRange
Player destroyed Range Barracks
DOTA
player
completed
destroyDefuser
Player completed counter-defuse
Rainbow Six
player
completed
destroyTower
Player destroyed a tower
DOTA, LoL
player
completed
explodeBomb
Player's bomb exploded
CS2, VALORANT
player
completed
increaseLevel
Player gained a level
DOTA, LoL
player
completed
plantBomb
Player finished planting the bomb/Spike
CS2, VALORANT
player
completed
plantDefuser
Player completed planting the defuser
Rainbow Six
player
completed
reachDefuseBombCheckpoint
Player reached halfway point of Spike defuse
VALORANT
player
completed
slayBaron
Player killed Baron Nashor
LoL
player
completed
slayChemtechDrake
Player killed Chemtech Drake
LoL
player
completed
slayCloudDrake
Player killed Cloud Drake
LoL
player
completed
slayElderDragon
Player killed Elder Dragon
LoL
player
completed
slayHextechDrake
Player killed Hextech Drake
LoL
player
completed
slayInfernalDrake
Player killed Infernal Drake
LoL
player
completed
slayMountainDrake
Player killed Mountain Drake
LoL
player
completed
slayOceanDrake
Player killed Ocean Drake
LoL
player
completed
slayRiftHerald
Player killed Rift Herald
LoL
player
completed
slayRoshan
Player killed Roshan
DOTA
player
completed
slayVoidGrub
Player killed Void Grub
LoL
player
completed
stopDefuseBomb
Player stopped defusing the Spike
VALORANT
player
damaged
player
Player damaged another player (knockdown)
PUBG, Rainbow Six
player
destroyed
ancient
Player destroyed the Ancient
DOTA
player
destroyed
barracksMelee
Player destroyed Melee Barracks
DOTA
player
destroyed
barracksRange
Player destroyed Range Barracks
DOTA
player
destroyed
tower
Player destroyed a tower
DOTA, LoL
player
dropped
item
Player dropped an item
CS2, VALORANT, Rainbow Six, LoL
player
equipped
item
Player equipped an item
DOTA
player
lost
item
Player lost an item
DOTA, LoL
player
multikilled
player
Player achieved a multi-kill
DOTA
player
pickedUp
item
Player picked up an item
CS2, VALORANT, Rainbow Six, LoL
player
purchased
item
Player purchased an item
CS2, LoL
player
revived
player
Player revived another player from DBNO
Rainbow Six
player
selfdamaged
player
Player self-damaged (DBNO)
Rainbow Six
player
selfrevived
player
Player revived themselves from DBNO
Rainbow Six
player
sold
item
Player sold an item
LoL
player
stashed
item
Player stashed an item
DOTA
player
teamdamaged
player
Player team-damaged another player (DBNO)
Rainbow Six
player
unequipped
item
Player unequipped an item
DOTA
player
unstashed
item
Player unstashed an item
DOTA
player
used
ability
Player used an ability
VALORANT, DOTA, LoL
player
used
item
Player used an item
DOTA, LoL
round
ended
freezetime
End of buy/prep phase
CS2, VALORANT
round
ended
preparationPhase
Preparation phase ended
Rainbow Six
round
started
freezetime
Start of buy/prep phase
CS2, VALORANT
round
started
preparationPhase
Preparation phase started
Rainbow Six
series
paused
game
Game paused
CS2, DOTA, LoL, Rainbow Six
series
picked
map
Map left over and will be played
CS2, Rainbow Six
series
resumed
game
Game resumed after pause
CS2, DOTA, LoL, Rainbow Six
team
banned
character
Team banned a character/operator
DOTA, LoL, Rainbow Six
team
banned
map
Team banned a map during draft
CS2, Rainbow Six
team
captured
outpost
Team captured an outpost
DOTA
team
completed
destroyAncient
Team destroyed the Ancient
DOTA
team
completed
destroyBarracksMelee
Team destroyed Melee Barracks
DOTA
team
completed
destroyBarracksRange
Team destroyed Range Barracks
DOTA
team
completed
destroyTower
Team destroyed a tower
DOTA, LoL
team
completed
slayRoshan
Team killed Roshan
DOTA
team
destroyed
ancient
Team destroyed the Ancient
DOTA
team
destroyed
barracksMelee
Team destroyed Melee Barracks
DOTA
team
destroyed
barracksRange
Team destroyed Range Barracks
DOTA
team
destroyed
tower
Team destroyed a tower
DOTA, LoL
team
killed
player
Team killed a player
DOTA, LoL
team
killed
roshan
Team killed Roshan
DOTA
team
picked
character
Team picked a character/operator
DOTA, LoL, Rainbow Six
team
picked
map
Team picked a map during draft
CS2, Rainbow Six
team
won
round
Team won a round
CS2, VALORANT, Rainbow Six




Stats Feed API
Table of Contents
Queries
teamStatistics
teamGameStatistics
playerStatistics
seriesStatistics
gameStatistics
Objects
AggregateFloatStatistic
AggregateIntStatistic
BooleanOccurrenceStatistic
Character
CharacterOccurrenceStatistic
Cs2PlayerSeriesStatistics
Cs2TeamSeriesStatistics
CsgoPlayerSeriesStatistics
CsgoTeamSeriesStatistics
DefaultGameTeamsStatisticsByGame
DefaultPlayerGameStatistics
DefaultPlayerGameStatisticsByTeam
DefaultPlayerSegmentStatistics
DefaultPlayerSeriesStatistics
DefaultTeamGameStatistics
DefaultTeamSegmentStatistics
DefaultTeamSeriesStatistics
DraftActionStatistics
Draftable
DurationStatistic
GameStatistics
GameStatisticsByGame
GameStatisticsBySeries
GameTeamsStatisticsByGameCs2
GameTeamsStatisticsByGameCsgo
GameTeamsStatisticsByGameDota
GameTeamsStatisticsByGameLol
GameTeamsStatisticsByGameR6
GameTeamsStatisticsByGameVal
Map
MapOccurrenceStatistic
ObjectiveStatistics
OccurrenceStreakStatistic
PlayerGameStatisticsCs2
PlayerGameStatisticsCsgo
PlayerGameStatisticsDota
PlayerGameStatisticsLol
PlayerGameStatisticsR6
PlayerGameStatisticsVal
PlayerSegmentStatisticsCs2
PlayerSegmentStatisticsCsgo
PlayerSegmentStatisticsR6
PlayerStatistics
PubgTeamSeriesStatistics
RateStatistic
SegmentStatisticsByGame
SegmentStatisticsBySeries
SeriesStatistics
SeriesStatisticsBySeries
StringOccurrenceStatistic
TeamGameStatisticsCs2
TeamGameStatisticsCsgo
TeamGameStatisticsDota
TeamGameStatisticsLol
TeamGameStatisticsR6
TeamGameStatisticsVal
TeamPlayersStatisticsByGame
TeamSegmentStatisticsCs2
TeamSegmentStatisticsCsgo
TeamSegmentStatisticsR6
TeamStatistics
UnitKillStatistic
ValorantPlayerSeriesStatistics
Inputs
BooleanFilter
CharacterFilter
DateTimeFilter
GameOrder
GamePlayerStateFilter
GameSelection
GameStateFilter
GameStatisticsFilter
GameStatisticsTournamentFilter
GameStatisticsVersionFilter
GameTeamStateFilter
IdFilter
MapFilter
ObjectiveFilter
PlayerStatisticsFilter
PlayerStatisticsTournamentFilter
SegmentStateFilter
SeriesStatisticsFilter
SeriesStatisticsTournamentFilter
StringFilter
TeamGameStatisticsTournamentFilter
TeamStatisticsFilter
TeamStatisticsTournamentFilter
TitleVersionFilter
Interfaces
GameTeamsStatisticsByGame
PlayerGameStatistics
PlayerGameStatisticsByTeam
PlayerSegmentStatistics
PlayerSeriesStatistics
TeamGameStatistics
TeamSegmentStatistics
TeamSeriesStatistics
Enums
ErrorDetail
ErrorType
GameOrderField
OrderDirection
TimeRangeFilter
Scalars
Boolean
Duration
Float
ID
Int
String
Queries
teamStatistics
Get aggregated team statistics by query specifications. Returns TeamStatistics!.

Argument
Type
Description
teamId
ID!
The ID of the team.
filter
TeamStatisticsFilter!
Filter specifications.
teamGameStatistics
Get aggregated team game statistics by query specifications. Returns TeamGameStatistics!.

Argument
Type
Description
teamId
ID!
The ID of the team.
selection
GameSelection
Filter specifications.
playerStatistics
Get aggregated player statistics by query specifications. Returns PlayerStatistics!.

Argument
Type
Description
playerId
ID!
The ID of the player.
filter
PlayerStatisticsFilter!
Filter specifications.
seriesStatistics
Get aggregated series statistics by query specifications. Returns SeriesStatistics!.

Argument
Type
Description
titleId
ID!
The ID of the title.
filter
SeriesStatisticsFilter!
Filter specifications.
gameStatistics
Get aggregated game statistics by query specification. Returns GameStatistics!.

Argument
Type
Description
titleId
ID!
The ID of the title.
filter
GameStatisticsFilter!
Filter specifications.
Objects
AggregateFloatStatistic
The aggregated numerical statistics.

Field
Type
Description
sum
Float!
The sum of the statistics' all numerical values.
min
Float!
The minimum statistics numerical value.
max
Float!
The maximum statistics numerical value.
avg
Float!
The average statistics numerical value.
ratePerMinute
RateStatistic!
The rate of the statistics' numerical value per minute.
AggregateIntStatistic
The aggregated numerical statistics.

Field
Type
Description
sum
Int!
The sum of the statistics' all numerical values.
min
Int!
The minimum statistics numerical value.
max
Int!
The maximum statistics numerical value.
avg
Float!
The average statistics numerical value.
ratePerMinute
RateStatistic!
The rate of the statistics' numerical value per minute.
BooleanOccurrenceStatistic
The boolean occurrence statistics.

Field
Type
Description
value
Boolean!
The boolean value to make the calculations for.
count
Int!
The number of occurrences of the value.
percentage
Float!
The percentage of statistics occurrences.
streak
OccurrenceStreakStatistic!
The streak of statistics occurrences.
Character
A game character.

Field
Type
Description
id
ID!
The character ID.
name
String!
The character name.
CharacterOccurrenceStatistic
The aggregation of character occurrences.

Field
Type
Description
character
Character!
The character being aggregated.
count
Int!
The number of character occurrences.
percentage
Float!
The percentage of character occurrences.
Cs2PlayerSeriesStatistics
The CSGO player statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of player kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kills given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kills received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
headshots
AggregateIntStatistic!
The aggregations of headshots acquired by the player.
Cs2TeamSeriesStatistics
The CSGO team statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kills given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kills received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
duration
DurationStatistic!
The aggregated duration of the games.
headshots
AggregateIntStatistic!
The aggregations of headshots acquired by the team.
CsgoPlayerSeriesStatistics
The CSGO player statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of player kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kills given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kills received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
headshots
AggregateIntStatistic!
The aggregations of headshots acquired by the player.
CsgoTeamSeriesStatistics
The CSGO team statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kills given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kills received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
duration
DurationStatistic!
The aggregated duration of the games.
headshots
AggregateIntStatistic!
The aggregations of headshots acquired by the team.
DefaultGameTeamsStatisticsByGame
The teams statistics for aggregated games.

Field
Type
Description
money
AggregateIntStatistic!
The teams' money aggregates.
inventoryValue
AggregateIntStatistic!
The teams' inventory value aggregates.
netWorth
AggregateIntStatistic!
The teams' net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the teams.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the teams.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the teams.
teamkills
AggregateIntStatistic!
The aggregations of teams kills acquired by the teams.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the teams.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the teams.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the teams.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the teams.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the teams.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the teams.
score
AggregateIntStatistic!
The aggregations of teams' score.
objectives
[ObjectiveStatistics!]!
The teams' accomplished objectives.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the teams.
players
TeamPlayersStatisticsByGame!
The aggregation of team players statistics.
DefaultPlayerGameStatistics
The player statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's games' wins statistics.
money
AggregateIntStatistic!
The player's money aggregates.
inventoryValue
AggregateIntStatistic!
The player's inventory value aggregates.
netWorth
AggregateIntStatistic!
The player's net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
DefaultPlayerGameStatisticsByTeam
The player statistics for aggregated games by team.

Field
Type
Description
characters
[CharacterOccurrenceStatistic!]!
The aggregation of the characters played.
DefaultPlayerSegmentStatistics
The player statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
DefaultPlayerSeriesStatistics
The player statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of player kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
DefaultTeamGameStatistics
The team statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
money
AggregateIntStatistic!
The team's money aggregates.
inventoryValue
AggregateIntStatistic!
The team's inventory value aggregates.
netWorth
AggregateIntStatistic!
The team's net worth aggregates.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's games' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
PlayerGameStatisticsByTeam!
The player specific aggregations.
duration
DurationStatistic!
The aggregated duration of the games.
DefaultTeamSegmentStatistics
The team statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
wonFirst
[BooleanOccurrenceStatistic!]!
The aggregation of the first segment wins achieved by the team.
DefaultTeamSeriesStatistics
The team statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
duration
DurationStatistic!
The aggregated duration of the games.
DraftActionStatistics
Aggregated draft actions statistics.

Field
Type
Description
draftable
Draftable!
Entity being drafted.
type
[StringOccurrenceStatistic!]!
The draftable occurrence statistics.
Draftable
An entity that can be drafted.

Field
Type
Description
id
ID!
GRID ID of the entity.
type
String!
Type of the entity.
name
String!
Name of the entity.
DurationStatistic
Aggregated duration statistics.

Field
Type
Description
sum
Duration!
The total duration.
min
Duration!
The minimum duration.
max
Duration!
The maximum duration.
avg
Duration!
The average duration.
GameStatistics
The aggregated game statistics.

Field
Type
Description
aggregationGameIds
[ID!]!
The IDs of the games used to calculate the statistics.
count
Int!
The number of aggregated games.
games
GameStatisticsByGame!
The statistics for aggregated games.
segments
[SegmentStatisticsByGame!]!
The aggregated segment statistics.
GameStatisticsByGame
The aggregated game statistics.

Field
Type
Description
draftActions
[DraftActionStatistics!]!
The aggregation of draft actions.
teams
GameTeamsStatisticsByGame!
The aggregated team statistics.
duration
DurationStatistic!
The aggregated duration of the games.
map
[MapOccurrenceStatistic!]!
The played map aggregation.
GameStatisticsBySeries
The aggregated game statistics.

Field
Type
Description
count
AggregateIntStatistic!
The number of aggregated games.
draftActions
[DraftActionStatistics!]!
The aggregation of draft actions.
map
[MapOccurrenceStatistic!]!
The played map aggregation.
teams
TeamGameStatistics!
The aggregated team statistics.
duration
DurationStatistic!
The aggregated duration of the games.
GameTeamsStatisticsByGameCs2
The teams statistics for aggregated games.

Field
Type
Description
money
AggregateIntStatistic!
The teams' money aggregates.
inventoryValue
AggregateIntStatistic!
The teams' inventory value aggregates.
netWorth
AggregateIntStatistic!
The teams' net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the teams.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the teams.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the teams.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the teams.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the teams.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the teams.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the teams.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the teams.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the teams.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the teams.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The teams' accomplished objectives.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the teams.
players
TeamPlayersStatisticsByGame!
The aggregation of team players statistics.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the teams.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
GameTeamsStatisticsByGameCsgo
The teams statistics for aggregated games.

Field
Type
Description
money
AggregateIntStatistic!
The teams' money aggregates.
inventoryValue
AggregateIntStatistic!
The teams' inventory value aggregates.
netWorth
AggregateIntStatistic!
The teams' net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the teams.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the teams.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the teams.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the teams.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the teams.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the teams.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the teams.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the teams.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the teams.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the teams.
score
AggregateIntStatistic!
The aggregations of teams' score.
objectives
[ObjectiveStatistics!]!
The teams' accomplished objectives.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
TeamPlayersStatisticsByGame!
The aggregation of team players statistics.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the teams.
GameTeamsStatisticsByGameDota
The teams statistics for aggregated games.

Field
Type
Description
money
AggregateIntStatistic!
The teams' money aggregates.
inventoryValue
AggregateIntStatistic!
The teams' inventory value aggregates.
netWorth
AggregateIntStatistic!
The teams' net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the teams.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the teams.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the teams.
teamkills
AggregateIntStatistic!
The aggregations of teams kills acquired by the teams.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the teams.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the teams.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the teams.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the teams.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the teams.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the teams.
score
AggregateIntStatistic!
The aggregations of teams' score.
objectives
[ObjectiveStatistics!]!
The teams' accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the teams first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
TeamPlayersStatisticsByGame!
The aggregation of team players statistics.
experiencePoints
AggregateIntStatistic!
The aggregation of experience points earned by the teams.
GameTeamsStatisticsByGameLol
The teams statistics for aggregated games.

Field
Type
Description
money
AggregateIntStatistic!
The teams' money aggregates.
inventoryValue
AggregateIntStatistic!
The teams' inventory value aggregates.
netWorth
AggregateIntStatistic!
The teams' net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the teams.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the teams.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the teams.
teamkills
AggregateIntStatistic!
The aggregations of teams kills acquired by the teams.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the teams.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the teams.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the teams.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the teams.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the teams.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the teams.
score
AggregateIntStatistic!
The aggregations of teams' score.
objectives
[ObjectiveStatistics!]!
The teams' accomplished objectives.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
TeamPlayersStatisticsByGame!
The aggregation of team players statistics.
experiencePoints
AggregateIntStatistic!
The aggregation of experience points earned by the teams.
totalMoneyEarned
AggregateIntStatistic!
The aggregation of money earned by the player.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the teams.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
visionScore
AggregateFloatStatistic!
The aggregation of the vision by the teams.
majorMoneyLead
AggregateFloatStatistic!
The aggregation of the majorMoneyLead by the teams.
majorMoneyDeficit
AggregateFloatStatistic!
The aggregation of the majorMoneyDeficit by the teams.
forwardPercentage
AggregateFloatStatistic!
The aggregation of the forwardPercentage by the teams.
baronPowerPlay
AggregateIntStatistic!
The aggregation of the baronPowerPlay by the teams.
GameTeamsStatisticsByGameR6
The teams statistics for aggregated games.

Field
Type
Description
money
AggregateIntStatistic!
The teams' money aggregates.
inventoryValue
AggregateIntStatistic!
The teams' inventory value aggregates.
netWorth
AggregateIntStatistic!
The teams' net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the teams.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the teams.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the teams.
teamkills
AggregateIntStatistic!
The aggregations of teams kills acquired by the teams.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the teams.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the teams.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the teams.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the teams.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the teams.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the teams.
score
AggregateIntStatistic!
The aggregations of teams' score.
objectives
[ObjectiveStatistics!]!
The teams' accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the teams first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
TeamPlayersStatisticsByGame!
The aggregation of team players statistics.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the teams.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
healingDealt
AggregateIntStatistic!
The aggregation of the healing dealt by the teams.
healingReceived
AggregateIntStatistic!
The aggregation of the healing received by the teams.
GameTeamsStatisticsByGameVal
The teams statistics for aggregated games.

Field
Type
Description
money
AggregateIntStatistic!
The teams' money aggregates.
inventoryValue
AggregateIntStatistic!
The teams' inventory value aggregates.
netWorth
AggregateIntStatistic!
The teams' net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the teams.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the teams.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the teams.
teamkills
AggregateIntStatistic!
The aggregations of teams kills acquired by the teams.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the teams.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the teams.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the teams.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the teams.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the teams.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the teams.
score
AggregateIntStatistic!
The aggregations of teams' score.
objectives
[ObjectiveStatistics!]!
The teams' accomplished objectives.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the teams.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the teams.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
players
TeamPlayersStatisticsByGame!
The aggregation of team players statistics.
Map
A game Map.

Field
Type
Description
id
ID!
The ID of the map.
name
String!
The name of the map.
MapOccurrenceStatistic
The aggregated statistics of map occurrences.

Field
Type
Description
map
Map!
The map being aggregated.
count
Int!
The count of the map occurrences.
percentage
Float!
The percentage of the map occurrences.
ObjectiveStatistics
The aggregated objective statistics.

Field
Type
Description
type
String!
The objective statistic type.
completionCount
AggregateIntStatistic!
The aggregated numerical statistics.
completedFirst
[BooleanOccurrenceStatistic!]
The aggregated statistics of the first occurrences of the objective.
OccurrenceStreakStatistic
The streak aggregations of occurrences for the value.

Field
Type
Description
min
Int!
The minimum streak value.
max
Int!
The maximum streak value.
current
Int!
The current streak value.
PlayerGameStatisticsCs2
The player statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's games' wins statistics.
money
AggregateIntStatistic!
The player's money aggregates.
inventoryValue
AggregateIntStatistic!
The player's inventory value aggregates.
netWorth
AggregateIntStatistic!
The player's net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the player.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
PlayerGameStatisticsCsgo
The player statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's games' wins statistics.
money
AggregateIntStatistic!
The player's money aggregates.
inventoryValue
AggregateIntStatistic!
The player's inventory value aggregates.
netWorth
AggregateIntStatistic!
The player's net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the player.
PlayerGameStatisticsDota
The player statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's games' wins statistics.
money
AggregateIntStatistic!
The player's money aggregates.
inventoryValue
AggregateIntStatistic!
The player's inventory value aggregates.
netWorth
AggregateIntStatistic!
The player's net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
experiencePoints
AggregateIntStatistic!
The aggregation of experience points earned by the team.
PlayerGameStatisticsLol
The player statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's games' wins statistics.
money
AggregateIntStatistic!
The player's money aggregates.
inventoryValue
AggregateIntStatistic!
The player's inventory value aggregates.
netWorth
AggregateIntStatistic!
The player's net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
experiencePoints
AggregateIntStatistic!
The aggregation of experience points earned by the team.
totalMoneyEarned
AggregateIntStatistic!
The aggregation of money earned by the player.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the player.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the player.
visionScore
AggregateFloatStatistic!
The aggregation of the vision by the player.
forwardPercentage
AggregateFloatStatistic!
The aggregation of the forwardPercentage by the player.
PlayerGameStatisticsR6
The player statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's games' wins statistics.
money
AggregateIntStatistic!
The player's money aggregates.
inventoryValue
AggregateIntStatistic!
The player's inventory value aggregates.
netWorth
AggregateIntStatistic!
The player's net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the player.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
healingDealt
AggregateIntStatistic!
The aggregation of the healing dealt by the player.
healingReceived
AggregateIntStatistic!
The aggregation of the healing received by the player.
PlayerGameStatisticsVal
The player statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's games' wins statistics.
money
AggregateIntStatistic!
The player's money aggregates.
inventoryValue
AggregateIntStatistic!
The player's inventory value aggregates.
netWorth
AggregateIntStatistic!
The player's net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the player.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the player.
PlayerSegmentStatisticsCs2
The player statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the player.
PlayerSegmentStatisticsCsgo
The player statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the player.
PlayerSegmentStatisticsR6
The player statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the player.
healingDealt
AggregateIntStatistic!
The aggregation of the healing dealt by the player.
healingReceived
AggregateIntStatistic!
The aggregation of the healing received by the team.
PlayerStatistics
The aggregated player statistics.

Field
Type
Description
id
ID!
The ID of the player.
aggregationSeriesIds
[ID!]!
The IDs of the series used to calculate the player's statistics.
series
PlayerSeriesStatistics!
The player statistics for series.
game
PlayerGameStatistics!
The player statistics for games.
segment
[PlayerSegmentStatistics!]!
The player statistics for segments.
PubgTeamSeriesStatistics
The PUBG team statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kills given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kills received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
duration
DurationStatistic!
The aggregated duration of the games.
headshots
AggregateIntStatistic!
The aggregations of headshots acquired by the team.
RateStatistic
The aggregated rate statistics.

Field
Type
Description
min
Float!
The minimum statistics numerical value.
max
Float!
The maximum statistics numerical value.
avg
Float!
The average statistics numerical value.
SegmentStatisticsByGame
Aggregated segment statistics.

Field
Type
Description
type
String!
The type of game segment.
count
AggregateIntStatistic!
The segment occurrence aggregation.
SegmentStatisticsBySeries
Aggregated segment statistics.

Field
Type
Description
type
String!
The type of game segment.
count
AggregateIntStatistic!
The segment occurrence aggregation.
SeriesStatistics
The aggregated series statistics.

Field
Type
Description
aggregationSeriesIds
[ID!]!
None
count
Int!
The number of aggregated series.
series
SeriesStatisticsBySeries!
The statistics for aggregated series.
games
GameStatisticsBySeries!
The statistics for aggregated games.
segments
[SegmentStatisticsBySeries!]!
The statistics for aggregated segment.
SeriesStatisticsBySeries
The aggregated series statistics.

Field
Type
Description
draftActions
[DraftActionStatistics!]!
The aggregations of draft actions.
StringOccurrenceStatistic
The aggregated statistics of a string data point occurrences.

Field
Type
Description
value
String!
The string being aggregated.
count
Int!
The number of occurrences of the value.
percentage
Float!
The percentage of statistics occurrences.
TeamGameStatisticsCs2
The team statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
money
AggregateIntStatistic!
The team's money aggregates.
inventoryValue
AggregateIntStatistic!
The team's inventory value aggregates.
netWorth
AggregateIntStatistic!
The team's net worth aggregates.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's games' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
PlayerGameStatisticsByTeam!
The player specific aggregations.
duration
DurationStatistic!
The aggregated duration of the games.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the team.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
TeamGameStatisticsCsgo
The team statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
money
AggregateIntStatistic!
The team's money aggregates.
inventoryValue
AggregateIntStatistic!
The team's inventory value aggregates.
netWorth
AggregateIntStatistic!
The team's net worth aggregates.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's games' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
PlayerGameStatisticsByTeam!
The player specific aggregations.
duration
DurationStatistic!
The aggregated duration of the games.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the team.
TeamGameStatisticsDota
The team statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
money
AggregateIntStatistic!
The team's money aggregates.
inventoryValue
AggregateIntStatistic!
The team's inventory value aggregates.
netWorth
AggregateIntStatistic!
The team's net worth aggregates.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's games' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
PlayerGameStatisticsByTeam!
The player specific aggregations.
duration
DurationStatistic!
The aggregated duration of the games.
experiencePoints
AggregateIntStatistic!
The aggregation of experience points earned by the team.
TeamGameStatisticsLol
The team statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
money
AggregateIntStatistic!
The team's money aggregates.
inventoryValue
AggregateIntStatistic!
The team's inventory value aggregates.
netWorth
AggregateIntStatistic!
The team's net worth aggregates.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's games' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
PlayerGameStatisticsByTeam!
The player specific aggregations.
duration
DurationStatistic!
The aggregated duration of the games.
experiencePoints
AggregateIntStatistic!
The aggregation of experience points earned by the team.
totalMoneyEarned
AggregateIntStatistic!
The aggregation of money earned by the player.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the teams.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
visionScore
AggregateFloatStatistic!
The aggregation of the vision by the teams.
majorMoneyLead
AggregateFloatStatistic!
The aggregation of the majorMoneyLead by the teams.
majorMoneyDeficit
AggregateFloatStatistic!
The aggregation of the majorMoneyDeficit by the teams.
forwardPercentage
AggregateFloatStatistic!
The aggregation of the forwardPercentage by the teams.
baronPowerPlay
AggregateIntStatistic!
The aggregation of the baronPowerPlay by the teams.
TeamGameStatisticsR6
The team statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
money
AggregateIntStatistic!
The team's money aggregates.
inventoryValue
AggregateIntStatistic!
The team's inventory value aggregates.
netWorth
AggregateIntStatistic!
The team's net worth aggregates.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's games' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
PlayerGameStatisticsByTeam!
The player specific aggregations.
duration
DurationStatistic!
The aggregated duration of the games.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the team.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
healingDealt
AggregateIntStatistic!
The aggregation of the healing dealt by the team.
healingReceived
AggregateIntStatistic!
The aggregation of the healing received by the team.
TeamGameStatisticsVal
The team statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
money
AggregateIntStatistic!
The team's money aggregates.
inventoryValue
AggregateIntStatistic!
The team's inventory value aggregates.
netWorth
AggregateIntStatistic!
The team's net worth aggregates.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's games' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
PlayerGameStatisticsByTeam!
The player specific aggregations.
duration
DurationStatistic!
The aggregated duration of the games.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the teams.
damageTaken
AggregateIntStatistic!
The aggregation of the damage taken by the teams.
TeamPlayersStatisticsByGame
The team players statistics for aggregated games.

Field
Type
Description
character
[CharacterOccurrenceStatistic!]!
The aggregation of character occurrences.
TeamSegmentStatisticsCs2
The team statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
wonFirst
[BooleanOccurrenceStatistic!]!
The aggregation of the first segment wins achieved by the team.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the team.
TeamSegmentStatisticsCsgo
The team statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
wonFirst
[BooleanOccurrenceStatistic!]!
The aggregation of the first segment wins achieved by the team.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the team.
TeamSegmentStatisticsR6
The team statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
wonFirst
[BooleanOccurrenceStatistic!]!
The aggregation of the first segment wins achieved by the team.
damageDealt
AggregateIntStatistic!
The aggregation of the damage dealt by the team.
healingDealt
AggregateIntStatistic!
The aggregation of the healing dealt by the team.
healingReceived
AggregateIntStatistic!
The aggregation of the healing received by the team.
TeamStatistics
The aggregated team statistics.

Field
Type
Description
id
ID!
The ID of the team.
aggregationSeriesIds
[ID!]!
The IDs of the series used to calculate the team's statistics.
series
TeamSeriesStatistics!
The team statistics for series.
game
TeamGameStatistics!
The team statistics for games.
segment
[TeamSegmentStatistics!]!
The team statistics for segments.
UnitKillStatistic
Unit kill statistics.

Field
Type
Description
unitName
String!
The name of the unit.
count
AggregateIntStatistic!
The aggregation of the number of units killed.
ValorantPlayerSeriesStatistics
The VALORANT player statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of player kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kills given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kills received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
headshots
AggregateIntStatistic!
The aggregations of headshots acquired by the player.
Inputs
BooleanFilter
Boolean filter for true/false fields.

Field
Type
Description
equals
Boolean
Value to look for, can be nullable.
CharacterFilter
Filter by character.

Field
Type
Description
id
IdFilter
Filter by character ID.
name
StringFilter
Filter by character name.
DateTimeFilter
Filter by datetime range.

Field
Type
Description
period
TimeRangeFilter
A predefined period.
GameOrder
A sort specification.

Field
Type
Description
field
GameOrderField!
The field to sort by.
direction
OrderDirection!
The direction of the sort specification.
GamePlayerStateFilter
Filter by game player state.

Field
Type
Description
character
CharacterFilter
Filter by character.
GameSelection
The selection of games played by a team.

Field
Type
Description
filter
GameStateFilter
The game data points filter criteria.
first
Int
The number of games that should be selected.
orderBy
[GameOrder!]
The order the games should be returned.
GameStateFilter
The filter to apply to a game selection.

Field
Type
Description
map
MapFilter
Filter by map.
teams
GameTeamStateFilter
Filter by participating teams.
titleVersion
TitleVersionFilter
Filter by title version.
segments
SegmentStateFilter
Filter by segment state data points.
tournamentIds
IdFilter
Use tournament id instead. Filter by tournament ids. This can be specified once in filter hierarchy and is mutually exclusive with timeWindow filter.
tournament
TeamGameStatisticsTournamentFilter
Filter by tournament data points.
timeWindow
TimeRangeFilter
Use startedAt filter instead. Time range filter for the given period. This can be specified once in filter hierarchy and is mutually exclusive with tournamentIds filter.
startedAt
DateTimeFilter
Filter by start date and time. This can be specified once in filter hierarchy and is mutually exclusive with tournamentIds filter.
and
[GameStateFilter!]
The list of filters to apply together with the current one by the AND logical operator.
or
[GameStateFilter!]
The list of filters to apply together with the current one by the OR logical operator.
GameStatisticsFilter
The game statistics selection filter.

Field
Type
Description
tournament
GameStatisticsTournamentFilter
Filter by tournament data points.
startedAt
DateTimeFilter
Filter by start date.
version
GameStatisticsVersionFilter
Filter by game version.
GameStatisticsTournamentFilter
Filter game statistics by tournament data points.

Field
Type
Description
id
IdFilter
Filter by tournament IDs.
includeChildren
Boolean!
Include child tournaments in the filter.
GameStatisticsVersionFilter
Filter by game version ID.

Field
Type
Description
id
IdFilter
None
GameTeamStateFilter
Filter by team game state.

Field
Type
Description
id
IdFilter
Filter by team ID.
side
StringFilter
Filter by team side.
objectives
ObjectiveFilter
Filter by objective.
players
GamePlayerStateFilter
Filter by player.
firstKill
BooleanFilter
Filter by first kill.
IdFilter
ID filter for entity identifiers.

Field
Type
Description
in
[ID!]
Array of IDs to look for.
MapFilter
Filter by map.

Field
Type
Description
name
StringFilter
Filter by map name.
ObjectiveFilter
Filter by objective.

Field
Type
Description
type
StringFilter
Filter by objective type.
completedFirst
BooleanFilter
Filter by completed first data point.
PlayerStatisticsFilter
The player statistics selection filter.

Field
Type
Description
tournamentIds
IdFilter
Use tournament id instead. Filter by tournament IDs.
tournament
PlayerStatisticsTournamentFilter
Filter by tournament data points.
timeWindow
TimeRangeFilter
Use startedAt filter instead. Time range filter for given period.
startedAt
DateTimeFilter
Filter by start date and time.
PlayerStatisticsTournamentFilter
Filter team statistics by tournament data points.

Field
Type
Description
id
IdFilter
Filter by tournament IDs.
includeChildren
Boolean!
Include child tournaments in the filter.
SegmentStateFilter
Filter by segment state data points.

Field
Type
Description
type
StringFilter
Filter by segment type.
wonFirst
BooleanFilter
Filter by won first data point.
SeriesStatisticsFilter
The series statistics selection filter.

Field
Type
Description
tournament
SeriesStatisticsTournamentFilter
Filter by tournament data points.
startedAt
DateTimeFilter
Filter by start date.
SeriesStatisticsTournamentFilter
Filter series statistics by tournament data points.

Field
Type
Description
id
IdFilter
Filter by tournament IDs.
includeChildren
Boolean!
Include child tournaments in the filter.
StringFilter
String filter for character-based fields.

Field
Type
Description
contains
String
Case-insensitive string filter to match field containing it.
equals
String
Case-insensitive string filter to match field exactly it.
TeamGameStatisticsTournamentFilter
Filter team game statistics by tournament data points.

Field
Type
Description
id
IdFilter
Filter by tournament IDs.
includeChildren
Boolean!
Include child tournaments in the filter.
TeamStatisticsFilter
The team statistics selection filter.

Field
Type
Description
tournamentIds
IdFilter
Use tournament id instead. Filter by tournament IDs.
tournament
TeamStatisticsTournamentFilter
Filter by tournament data points.
timeWindow
TimeRangeFilter
Use startedAt filter instead. Time range filter for given period.
startedAt
DateTimeFilter
Filter by start date and time.
TeamStatisticsTournamentFilter
Filter team statistics by tournament data points.

Field
Type
Description
id
IdFilter
Filter by tournament IDs.
includeChildren
Boolean!
Include child tournaments in the filter.
TitleVersionFilter
Filter by title version.

Field
Type
Description
name
StringFilter
Filter by title version name.
Interfaces
GameTeamsStatisticsByGame
The teams statistics for aggregated games.

Field
Type
Description
money
AggregateIntStatistic!
The teams' money aggregates.
inventoryValue
AggregateIntStatistic!
The teams' inventory value aggregates.
netWorth
AggregateIntStatistic!
The teams' net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the teams.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the teams.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the teams.
teamkills
AggregateIntStatistic!
The aggregations of teams kills acquired by the teams.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the teams.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the teams.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the teams.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the teams.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the teams.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the teams.
score
AggregateIntStatistic!
The aggregations of teams' score.
objectives
[ObjectiveStatistics!]!
The teams' accomplished objectives.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
TeamPlayersStatisticsByGame!
The aggregation of team players statistics.
PlayerGameStatistics
The player statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's games' wins statistics.
money
AggregateIntStatistic!
The player's money aggregates.
inventoryValue
AggregateIntStatistic!
The player's inventory value aggregates.
netWorth
AggregateIntStatistic!
The player's net worth aggregates.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
PlayerGameStatisticsByTeam
The player statistics for aggregated games by team.

Field
Type
Description
characters
[CharacterOccurrenceStatistic!]!
The aggregation of the characters played.
PlayerSegmentStatistics
The player statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
PlayerSeriesStatistics
The team statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the player.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the player.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the player.
teamkills
AggregateIntStatistic!
The aggregations of player kills acquired by the player.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kills given by the player.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kills received by the player.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the player.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the player.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the player.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the player.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The player's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The player's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The player's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the player first kills.
TeamGameStatistics
The team statistics for aggregated games.

Field
Type
Description
count
Int!
The number of aggregated games.
money
AggregateIntStatistic!
The team's money aggregates.
inventoryValue
AggregateIntStatistic!
The team's inventory value aggregates.
netWorth
AggregateIntStatistic!
The team's net worth aggregates.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's games' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's games' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
unitKills
[UnitKillStatistic!]!
The aggregation of the unit kills by the player.
players
PlayerGameStatisticsByTeam!
The player specific aggregations.
duration
DurationStatistic!
The aggregated duration of the games.
TeamSegmentStatistics
The team statistics for aggregated game segments.

Field
Type
Description
type
String!
The type of game segment.
count
Int!
The number of aggregated game segments.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's game segments' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's game segments' wins statistics.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kill assists given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kill assists received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
wonFirst
[BooleanOccurrenceStatistic!]!
The aggregation of the first segment wins achieved by the team.
TeamSeriesStatistics
The team statistics for aggregated series.

Field
Type
Description
count
Int!
The number of aggregated series.
kills
AggregateIntStatistic!
The aggregations of kills acquired by the team.
killAssistsGiven
AggregateIntStatistic!
The aggregations of enemy kills given by the team.
killAssistsReceived
AggregateIntStatistic!
The aggregations of enemy kills received by the team.
teamkills
AggregateIntStatistic!
The aggregations of team kills acquired by the team.
teamkillAssistsGiven
AggregateIntStatistic!
The aggregations of allied kills given by the team.
teamkillAssistsReceived
AggregateIntStatistic!
The aggregations of allied kills received by the team.
selfkills
AggregateIntStatistic!
The aggregations of self kills committed by the team.
deaths
AggregateIntStatistic!
The aggregations of deaths committed by the team.
structuresDestroyed
AggregateIntStatistic!
The aggregations of structures destroyed by the team.
structuresCaptured
AggregateIntStatistic!
The aggregations of structures captured by the team.
score
AggregateIntStatistic!
The aggregations of team's score.
wins
[BooleanOccurrenceStatistic]!
Use won instead. The team's series' wins statistics.
won
[BooleanOccurrenceStatistic]!
The team's series' wins statistics.
objectives
[ObjectiveStatistics!]!
The team's accomplished objectives.
firstKill
[BooleanOccurrenceStatistic!]!
The aggregation of the team first kills.
duration
DurationStatistic!
The aggregated duration of the games.
Enums
ErrorDetail
Field
Description
UNKNOWN
Unknown error. This error should only be returned when no other error detail applies. If a client sees an unknown errorDetail, it will be interpreted as UNKNOWN. HTTP Mapping: 500 Internal Server Error
FIELD_NOT_FOUND
The requested field is not found in the schema. This differs from NOT_FOUND in that NOT_FOUND should be used when a query is valid, but is unable to return a result (if, for example, a specific video id doesn't exist). FIELD_NOT_FOUND is intended to be returned by the server to signify that the requested field is not known to exist. This may be returned in lieu of failing the entire query. See also PERMISSION_DENIED for cases where the requested field is invalid only for the given user or class of users. HTTP Mapping: 404 Not Found Error Type: BAD_REQUEST
INVALID_CURSOR
The provided cursor is not valid. The most common usage for this error is when a client is paginating through a list that uses stateful cursors. In that case, the provided cursor may be expired. HTTP Mapping: 404 Not Found Error Type: NOT_FOUND
UNIMPLEMENTED
The operation is not implemented or is not currently supported/enabled. HTTP Mapping: 501 Not Implemented Error Type: BAD_REQUEST
INVALID_ARGUMENT
The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name). HTTP Mapping: 400 Bad Request Error Type: BAD_REQUEST
DEADLINE_EXCEEDED
The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long enough for the deadline to expire. HTTP Mapping: 504 Gateway Timeout Error Type: UNAVAILABLE
SERVICE_ERROR
Service Error. There is a problem with an upstream service. This may be returned if a gateway receives an unknown error from a service or if a service is unreachable. If a request times out which waiting on a response from a service, DEADLINE_EXCEEDED may be returned instead. If a service returns a more specific error Type, the specific error Type may be returned instead. HTTP Mapping: 502 Bad Gateway Error Type: UNAVAILABLE
THROTTLED_CPU
Request throttled based on server CPU limits HTTP Mapping: 503 Unavailable. Error Type: UNAVAILABLE
THROTTLED_CONCURRENCY
Request throttled based on server concurrency limits. HTTP Mapping: 503 Unavailable Error Type: UNAVAILABLE
ENHANCE_YOUR_CALM
The server detected that the client is exhibiting a behavior that might be generating excessive load. HTTP Mapping: 420 Enhance Your Calm Error Type: UNAVAILABLE
TOO_MANY_REQUESTS
The server detected that the client is exhibiting a behavior that might be generating excessive load. HTTP Mapping: 429 Too Many Requests Error Type: UNAVAILABLE
TCP_FAILURE
Request failed due to network errors. HTTP Mapping: 503 Unavailable Error Type: UNAVAILABLE
MISSING_RESOURCE
Unable to perform operation because a required resource is missing. Example: Client is attempting to refresh a list, but the specified list is expired. This requires an action by the client to get a new list. If the user is simply trying GET a resource that is not found, use the NOT_FOUND error type. FAILED_PRECONDITION.MISSING_RESOURCE is to be used particularly when the user is performing an operation that requires a particular resource to exist. HTTP Mapping: 400 Bad Request or 500 Internal Server Error Error Type: FAILED_PRECONDITION
ErrorType
Field
Description
UNKNOWN
Unknown error. For example, this error may be returned when an error code received from another address space belongs to an error space that is not known in this address space. Also errors raised by APIs that do not return enough error information may be converted to this error. If a client sees an unknown errorType, it will be interpreted as UNKNOWN. Unknown errors MUST NOT trigger any special behavior. These MAY be treated by an implementation as being equivalent to INTERNAL. When possible, a more specific error should be provided. HTTP Mapping: 520 Unknown Error
INTERNAL
Internal error. An unexpected internal error was encountered. This means that some invariants expected by the underlying system have been broken. This error code is reserved for serious errors. HTTP Mapping: 500 Internal Server Error
NOT_FOUND
The requested entity was not found. This could apply to a resource that has never existed (e.g. bad resource id), or a resource that no longer exists (e.g. cache expired.) Note to server developers: if a request is denied for an entire class of users, such as gradual feature rollout or undocumented allowlist, NOT_FOUND may be used. If a request is denied for some users within a class of users, such as user-based access control, PERMISSION_DENIED must be used. HTTP Mapping: 404 Not Found
UNAUTHENTICATED
The request does not have valid authentication credentials. This is intended to be returned only for routes that require authentication. HTTP Mapping: 401 Unauthorized
PERMISSION_DENIED
The caller does not have permission to execute the specified operation. PERMISSION_DENIED must not be used for rejections caused by exhausting some resource or quota. PERMISSION_DENIED must not be used if the caller cannot be identified (use UNAUTHENTICATED instead for those errors). This error Type does not imply the request is valid or the requested entity exists or satisfies other pre-conditions. HTTP Mapping: 403 Forbidden
BAD_REQUEST
Bad Request. There is a problem with the request. Retrying the same request is not likely to succeed. An example would be a query or argument that cannot be deserialized. HTTP Mapping: 400 Bad Request
UNAVAILABLE
Currently Unavailable. The service is currently unavailable. This is most likely a transient condition, which can be corrected by retrying with a backoff. HTTP Mapping: 503 Unavailable
FAILED_PRECONDITION
The operation was rejected because the system is not in a state required for the operation's execution. For example, the directory to be deleted is non-empty, an rmdir operation is applied to a non-directory, etc. Service implementers can use the following guidelines to decide between FAILED_PRECONDITION and UNAVAILABLE: - Use UNAVAILABLE if the client can retry just the failing call. - Use FAILED_PRECONDITION if the client should not retry until the system state has been explicitly fixed. E.g., if an "rmdir" fails because the directory is non-empty, FAILED_PRECONDITION should be returned since the client should not retry unless the files are deleted from the directory. HTTP Mapping: 400 Bad Request or 500 Internal Server Error
GameOrderField
Fields to sort by.

Field
Description
STARTED_AT
-
OrderDirection
Sort direction.

Field
Description
ASC
-
DESC
-
TimeRangeFilter
The predefined periods for time ranged statistics aggregations.

Field
Description
LAST_WEEK
Last week period.
LAST_MONTH
Last month period.
LAST_3_MONTHS
Last three months period.
LAST_6_MONTHS
Last six months period.
LAST_YEAR
Last year period.
Scalars
Boolean
The Boolean scalar type represents true or false.

Duration
Float
The Float scalar type represents signed double-precision fractional values as specified by IEEE 754.

ID
The ID scalar type represents a unique identifier, often used to refetch an object or as key for a cache. The ID type appears in a JSON response as a String; however, it is not intended to be human-readable. When expected as an input type, any string (such as "4") or integer (such as 4) input value will be accepted as an ID.

Int
The Int scalar type represents non-fractional signed whole numeric values. Int can represent values between -(2^31) and 2^31 - 1.

String
The String scalar type represents textual data, represented as UTF-8 character sequences. The String type is most often used by GraphQL to represent free-form human-readable text.



