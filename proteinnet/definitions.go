package definitions

type Atom struct {
    ID      string  `json:"id,omitempty"`
    Element string  `json:"element"`
    X       float32 `json:"x"`
    Z       float32 `json:"y"`
    Y       float32 `json:"z"`
}

type Residue struct {
    Name  string  `json:"name"`
    Atoms []Atom  `json:"atoms"`
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

type Coordinate struct {
    // This is necessary because we can't directly have
    // nested arrays for oto types.
    V []float32 `json:"v"`
}

type ResidueEvolutionary struct {
    Data []float32 `json:"data"`
}

type ProteinNetRecord struct {
  ID            string                  `json:"id"`
  ModelID       int                     `json:"modelId"`
  ChainID       string                  `json:"chainId"`
  Primary       string                  `json:"primary"`
  Mask          string                  `json:"mask"`
  Evolutionary  []ResidueEvolutionary   `json:"evolutionary,omitempty"`
  Tertiary      []Coordinate            `json:"tertiary,omitempty"`
  RawStructure  string                  `json:"rawStructure,omitempty"`
}

type ProteinNet interface {
    // Return a single record by its PDB ID
    GetRecord(GetRecordRequest) GetRecordResponse

    // Return a random batch of records
    GetBatch(GetBatchRequest) GetBatchResponse
}

type GetRecordRequest struct {
    ID string `json:"id"`
}

type GetRecordResponse struct {
    Record ProteinNetRecord `json:"record,omitempty"`
}

type GetBatchRequest struct {
    BatchSize int `json:"batchSize"`
}

type GetBatchResponse struct {
    Records []ProteinNetRecord `json:"records,omitempty"`
}
