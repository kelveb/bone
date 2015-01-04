/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/squiidz        ***
*********************************/

package bone

import "net/http"

// Mux have routes and a notFound handler
// Route: all the registred route
// notFound: 404 handler, default http.NotFound if not provided
type Mux struct {
	Routes   map[string][]*Route
	Static   map[string]*Route
	notFound http.HandlerFunc
}

var (
	method = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS"}
	vars   = make(map[*http.Request]*Route)
)

// New create a pointer to a Mux instance
func New() *Mux {
	return &Mux{
		Routes: make(map[string][]*Route),
		Static: make(map[string]*Route),
	}
}

// Serve http request
func (m *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	reqLen := len(reqPath)

	// Check if the request path doesn't end with /
	if !m.valid(reqPath) {
		http.Redirect(rw, req, reqPath[:reqLen-1], http.StatusMovedPermanently)
		return
	}
	// Loop over all the registred route.
	for _, r := range m.Routes[req.Method] {
		// If the route is equal to the request path.
		if reqPath == r.Path {
			r.Handler.ServeHTTP(rw, req)
			return
		} else if r.Pattern.Exist {
			if v, ok := r.Match(req.URL.Path); ok {
				r.insert(req, v)
				r.Handler.ServeHTTP(rw, req)
				return
			}
			continue
		}
		continue
	}
	// If no valid Route found, check for static file
	for _, s := range m.Static {
		if reqLen >= s.Size && reqPath[:s.Size] == s.Path {
			s.Handler.ServeHTTP(rw, req)
			return
		}
		continue
	}
	m.BadRequest(rw, req)
}

// HandleFunc is use to pass a func(http.ResponseWriter, *Http.Request) instead of http.Handler
func (m *Mux) HandleFunc(path string, handler http.HandlerFunc) {
	m.Handle(path, handler)
}

// Handle add a new route to the Mux without a HTTP method
func (m *Mux) Handle(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	if m.isStatic(path) {
		m.Static[path] = r.Get()
		return
	}
	for _, mt := range method {
		m.Routes[mt] = append(m.Routes[mt], r)
		byLength(m.Routes[mt]).Sort()
	}
}

// Get add a new route to the Mux with the Get method
func (m *Mux) Get(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["GET"] = append(m.Routes["GET"], r.Get())
	byLength(m.Routes["GET"]).Sort()
}

// Post add a new route to the Mux with the Post method
func (m *Mux) Post(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["POST"] = append(m.Routes["POST"], r.Post())
	byLength(m.Routes["POST"]).Sort()
}

// Put add a new route to the Mux with the Put method
func (m *Mux) Put(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["PUT"] = append(m.Routes["PUT"], r.Put())
	byLength(m.Routes["PUT"]).Sort()
}

// Delete add a new route to the Mux with the Delete method
func (m *Mux) Delete(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["DELETE"] = append(m.Routes["DELETE"], r.Delete())
	byLength(m.Routes["DELETE"]).Sort()
}

// Head add a new route to the Mux with the Head method
func (m *Mux) Head(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["HEAD"] = append(m.Routes["HEAD"], r.Head())
	byLength(m.Routes["HEAD"]).Sort()
}

// Patch add a new route to the Mux with the Patch method
func (m *Mux) Patch(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["PATCH"] = append(m.Routes["PATCH"], r.Patch())
	byLength(m.Routes["PATCH"]).Sort()
}

// Options add a new route to the Mux with the Options method
func (m *Mux) Options(path string, handler http.Handler) {
	r := NewRoute(path, handler)
	m.Routes["OPTIONS"] = append(m.Routes["OPTIONS"], r.Options())
	byLength(m.Routes["OPTIONS"]).Sort()
}

// NotFound the mux custom 404 handler
func (m *Mux) NotFound(handler http.HandlerFunc) {
	m.notFound = handler
}
