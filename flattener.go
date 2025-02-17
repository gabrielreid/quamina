package quamina

// Flattener is provided as an interface in the hope that flatterners for other
// non-JSON message formats might be implemented. How it needs to work, by JSON
// example:
//  { "a": 1, "b": "two", "c": true", "d": nil, "e": { "e1": 2, "e2":, 3.02e-5} "f": [33, "x"]} }
// should produce
//  "a", "1"
//  "b", "\"two\"",
//  "c", "true"
//  "d", "nil",
//  "e\ne1", "2"
//  "e\ne2", "3.02e-5"
//  "f", "33"
//  "f", "\"x\""
//
// Let's call the first column, eg "d" and "e\ne1", the pathSegments. For each
// step i the pathSegments, e.g. "d" and "e1", the Flattener shold call
// NameTracker.IsNameUsed(step) and if that comes back negative, not include any
// paths which don't contain that step. So in the example above, if
// nameTracker.IsNameUsed() only came back true for "a" and "f", then the output
// would be
//  "a", "1"
//  "f", "33"
//  "f", "\"x\""
type Flattener interface {
	Flatten(event []byte, tracker NameTracker) ([]Field, error)
	Copy() Flattener
}

// Arrays are invisible in the automaton.  That is to say, if an event has
//  { "a": [ 1, 2, 3 ] }
// Then the Flattener must produce a/1, a/2, and a/3 Same for  {"a": [[1, 2], 3]} or any
// other permutation.
// If we're not careful, this would create a problem. If you have
//  {"a": [ { "b": 1, "c": 2}, {"b": 3, "c": 4}] }
// then a pattern like
//  { "a": { "b": 1, "c": 4 } }
// would match. To prevent that from happening, each ArrayPos contains two
// numbers; the first identifies the array in
// the event that this name/val occurred in, the second the position in the array. We don't allow
// transitioning between field values that occur in different positions in the same array.
// See the arrays_test unit for more examples, and the jsonFlattener source code to
// see how it's implemented

// ArrayPos represents a Field's position in an Event's structure. Each array in the Event
// should get an integer which identifies it - in flattenJSON this is accomplished by keeping a counter and
// giving arrays numbers starting from 0.
// Array uniquely identifies an array in an Event
// Pos is the Field's index in the Array
type ArrayPos struct {
	Array int32
	Pos   int32
}

// Field represents a pathname/value combination, the data item which is matched against Patterns by the
// MatchesForEvent API.
// Path is \n-separated path from the event root to this field value
// Val is the value - note no information as to type
// ArrayTrail, for each array in the Path, identifies the array and the index in it
type Field struct {
	Path       []byte
	Val        []byte
	ArrayTrail []ArrayPos
}
