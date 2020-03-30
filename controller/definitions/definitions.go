package definitions

type Backend interface {
	// Run the experiment
	Run(RunRequest) RunResponse

	// Run some tests to ensure internal consistency
	Test(TestRequest) TestResponse

	// Produce a visualization
	Visualize(VisualizeRequest) VisualizeResponse
}

type Atom struct {
	ID      string  `json:"id,omitempty"`
	Element string  `json:"element"`
	X       float32 `json:"x"`
	Z       float32 `json:"y"`
	Y       float32 `json:"z"`
}

type Residue struct {
	Atoms []Atom `json:"atoms"`
}

type Chain struct {
	Residues []Residue `json:"residues"`
}

type Model struct {
	Chains map[string]Chain `json:"chains"`
}

type Structure struct {
	Models map[int]Model `json:"models"`
}

// Dataset -> Tranformation -> Transformation -> Backend
// Received AFTER transformations
type GromacsConfig struct {
	NumSteps       int     `json:"numSteps"`                 // nsteps: Number of steps to run the experiment
	Integrator     string  `json:"integrator,omitempty"`     // Algorithm (steep = steepest descent minimization). Default is "md"
	Seed           int     `json:"seed,omitempty"`           // ld-seed: Langevin dynamics seed
	DeltaTime      float64 `json:"deltaTime,omitempty"`      // dt: Deltatime in seconds. Default is 2e-15 seconds == 2 fs == 0.002 ps
	OutSteps       int     `json:"outSteps,omitempty"`       // nstxout: Output interval. 0 == output only at the end.
	CoupleIntramol bool    `json:"coupleIntramol,omitempty"` // couple-intramol: The intra-molecular Van der Waals and Coulomb interactions are also turned on/off
	StopTolerance  float64 `json:"stopTolerance,omitempty"`  // emtol: Stopping tolerance. 0 ==
}

type RunRequest struct {
	ID      string                 `json:"id"`      // Run ID
	Backend string                 `json:"backend"` // e.g. "gromacs"
	Input   string                 `json:"input"`   // Input? TODO
	Config  map[string]interface{} `json:"config"`  // Backend config
	Foo  map[int]string `json:"config"`  // Backend config
}

type RunResponse struct {
	// Experiment output information
}

type TestRequest struct {
}

type TestResponse struct {
}

type VisualizeRequest struct {
	FPS       int `json:"fps"`
	NumFrames int `json:"numFrames"`
}

type VisualizeResponse struct {
	Resource string `json:"resource"` // URL to the webm file, https://[...].webm
}
