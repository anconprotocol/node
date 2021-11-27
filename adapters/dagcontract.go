package adapters

// DagTransaction is the DAG transaction
type DagTransaction struct {
	SchemaCid        string
	DataSourceCid    string
	Variables        string
	ContractMutation string
	Result        string
	Signature string	
}
