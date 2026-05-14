// Package pipeline provides a lightweight, composable processing pipeline
// for confmap configuration maps.
//
// A Pipeline is an ordered sequence of Stages. Each Stage receives the
// current config map and returns a (possibly modified) copy. Stages are
// applied sequentially; if any stage returns an error the pipeline halts
// and propagates the error with positional context.
//
// Built-in stage constructors:
//
//	- FilterKeys  – drop keys matching given prefixes
//	- SetDefaults – fill in missing keys with default values
//	- RequireKeys – assert that mandatory keys are present
//	- StageFunc   – wrap any function as a labelled stage
//
// Example:
//
//	pl := pipeline.New("main").
//	    Use(pipeline.SetDefaults(map[string]any{"timeout": 30})).
//	    Use(pipeline.RequireKeys("host", "port")).
//	    Use(pipeline.FilterKeys("internal_"))
//
//	cfg, err := pl.Run(raw)
package pipeline
