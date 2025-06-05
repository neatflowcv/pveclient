package domain

type Cluster struct {
	id   string
	url  string
	name string
}

func NewCluster(id string, url string, name string) *Cluster {
	return &Cluster{
		id:   id,
		url:  url,
		name: name,
	}
}

func (c *Cluster) ID() string {
	return c.id
}

func (c *Cluster) URL() string {
	return c.url
}

func (c *Cluster) Name() string {
	return c.name
}

type ClusterSpec struct {
	url   string
	token string
}

func NewClusterSpec(url string, token string) *ClusterSpec {
	return &ClusterSpec{
		url:   url,
		token: token,
	}
}

func (s *ClusterSpec) URL() string {
	return s.url
}

func (s *ClusterSpec) Token() string {
	return s.token
}
