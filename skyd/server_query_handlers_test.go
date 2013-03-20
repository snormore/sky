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

// Ensure that we can query the server for multiple selections with multiple dimensions.
func TestServerMultiDimensionalQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "gender", "object", "string")
		setupTestProperty("foo", "state", "object", "string")
		setupTestProperty("foo", "price", "action", "float")
		setupTestData(t, "foo", [][]string{
			[]string{"0", "2012-01-01T00:00:00Z", `{"data":{"gender":"m", "state":"NY", "price":100}}`},
			[]string{"0", "2012-01-01T00:00:01Z", `{"data":{"price":200}}`},
			[]string{"0", "2012-01-01T00:00:02Z", `{"data":{"state":"CA","price":10}}`},

			[]string{"1", "2012-01-01T00:00:00Z", `{"data":{"gender":"m", "state":"CA", "price":20}}`},
			[]string{"1", "2012-01-01T00:00:01Z", `{"data":{}}`},

			[]string{"2", "2012-01-01T00:00:00Z", `{"data":{"gender":"f", "state":"NY", "price":30}}`},
		})

		// Run query.
		query := `{
			"steps":[
				{"type":"selection","alias":"count","dimensions":["gender","state"],"expression":"count()","steps":[]},
				{"type":"selection","alias":"sum","dimensions":["gender","state"],"expression":"sum(price)","steps":[]},
				{"type":"selection","alias":"minPrice","dimensions":["gender","state"],"expression":"min(price)","steps":[]},
				{"type":"selection","alias":"maxPrice","dimensions":["gender","state"],"expression":"max(price)","steps":[]}
			]
		}`
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8585/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"gender":{"f":{"state":{"NY":{"count":1,"maxPrice":30,"minPrice":30,"sum":30}}},"m":{"state":{"CA":{"count":3,"maxPrice":20,"minPrice":0,"sum":30},"NY":{"count":2,"maxPrice":200,"minPrice":100,"sum":300}}}}}`+"\n", "POST /tables/:name/query failed.")
	})
}

// Ensure that we can perform a non-sessionized funnel analysis.
func TestServerFunnelAnalysisQuery(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "action", "object", "string")
		setupTestData(t, "foo", [][]string{
			// A0[0..0]..A1[1..2] occurs twice for this object.
			[]string{"0", "2012-01-01T00:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"0", "2012-01-01T00:00:01Z", `{"data":{"action":"A1"}}`},
			[]string{"0", "2012-01-01T00:00:02Z", `{"data":{"action":"A2"}}`},
			[]string{"0", "2012-01-01T12:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"0", "2012-01-01T13:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"0", "2012-01-01T14:00:00Z", `{"data":{"action":"A1"}}`},

			// A0[0..0]..A1[1..2] occurs once for this object. (Second time matches A1[1..3]).
			[]string{"1", "2012-01-01T00:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"1", "2012-01-01T00:00:01Z", `{"data":{"action":"A0"}}`},
			[]string{"1", "2012-01-01T00:00:02Z", `{"data":{"action":"A1"}}`},
			[]string{"1", "2012-01-02T00:00:00Z", `{"data":{"action":"A0"}}`},
			[]string{"1", "2012-01-02T00:00:01Z", `{"data":{"action":"A0"}}`},
			[]string{"1", "2012-01-02T00:00:02Z", `{"data":{"action":"A0"}}`},
			[]string{"1", "2012-01-02T00:00:03Z", `{"data":{"action":"A1"}}`},
		})

		// Run query.
		query := `{
			"steps":[
				{"type":"condition","expression":"action == 'A0'","steps":[
					{"type":"condition","expression":"action == 'A1'","within":[1,2],"steps":[
						{"type":"selection","alias":"count","dimensions":["action"],"expression":"count()"}
					]}
				]}
			]
		}`
		//_codegen(t, "foo", query)
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8585/tables/foo/query", "application/json", query)
		assertResponse(t, resp, 200, `{"action":{"A1":{"count":3}}}`+"\n", "POST /tables/:name/query failed.")
	})
}
