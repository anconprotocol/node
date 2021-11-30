package adapters

// DagTransaction is the DAG transaction
type DagTransaction struct {
	MetadataCid string
	ResultCid   string
	FromOwner   string
	ToOwner     string
	Signature   string
}
