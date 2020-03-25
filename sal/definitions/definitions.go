package definitions

type Backend interface {
	Run(RunRequest) RunResponse
}

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
	StructureIDs []string               `json:"structureIds"`
	Backend      string                 `json:"backend"` // e.g. "gromacs"
	Config       map[string]interface{} `json:"config"`  // Backend config
	Gromacs      []GromacsConfig        `json:"gromacs"`
	//Configs      []map[string]GromacsConfig `json:"configs"`
}

type RunResponse struct {
	// Output information
}
