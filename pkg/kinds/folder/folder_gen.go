// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     GoTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package folder

// Resource is the wire representation of Folder. (TODO be better)
type Resource struct {
	Metadata Metadata `json:"metadata"`
	Spec Spec `json:"spec"`
	Status Status `json:"status"`
}
