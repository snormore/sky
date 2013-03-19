package skyd

import (
	"testing"
)

// Ensure that we can query the server for a count of events.
func TestServerSimpleCountQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "fruit", "action", "string")
		setupTestData(t, "foo", [][]string{
			[]string{"0", "2012-01-01T00:00:00Z", `{"data":{"fruit":"apple"}}`},
			[]string{"1", "2012-01-01T00:00:00Z", `{"data":{"fruit":"grape"}}`},
			[]string{"1", "2012-01-01T00:00:01Z", `{}`},
			[]string{"2", "2012-01-01T00:00:00Z", `{"data":{"fruit":"orange"}}`},
			[]string{"3", "2012-01-01T00:00:00Z", `{"data":{"fruit":"apple"}}`},
		})

		// Run query.
		query := `{
			"steps":[
				{"type":"selection","alias":"count","dimensions":[],"expression":"count()","steps":[]}
			]
		}`
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8585/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"count":5}`+"\n", "POST /tables/:name/query failed.")
	})
}

// Ensure that we can query the server for a count of events with a single dimension.
func TestServerOneDimensionCountQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "fruit", "action", "string")
		setupTestData(t, "foo", [][]string{
			[]string{"0", "2012-01-01T00:00:00Z", `{"data":{"fruit":"apple"}}`},
			[]string{"1", "2012-01-01T00:00:00Z", `{"data":{"fruit":"grape"}}`},
			[]string{"1", "2012-01-01T00:00:01Z", `{}`},
			[]string{"2", "2012-01-01T00:00:00Z", `{"data":{"fruit":"orange"}}`},
			[]string{"3", "2012-01-01T00:00:00Z", `{"data":{"fruit":"apple"}}`},
		})

		// Run query.
		query := `{
			"steps":[
				{"type":"selection","alias":"count","dimensions":["fruit"],"expression":"count()","steps":[]}
			]
		}`
		//_codegen(t, "foo", query)
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8585/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"fruit":{"":{"count":1},"apple":{"count":2},"grape":{"count":1},"orange":{"count":1}}}`+"\n", "POST /tables/:name/query failed.")
	})
}
