package filesystem

type DataPaths struct {
	DataDir          string
	DatabasesDir     string
	MainDatabaseFile string
}

// @Injectable
func GetDataPaths(d *DataDirInitializer) (*DataPaths, error) {
	if err := d.Init(); err != nil {
		return nil, err
	}

	return d.Paths, nil
}
