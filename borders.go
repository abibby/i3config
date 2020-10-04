package i3config

type BorderType string

const (
	None       BorderType = "none"
	Vertical   BorderType = "vertical"
	Horizontal BorderType = "horizontal"
	Both       BorderType = "both"
	Smart      BorderType = "smart"
)

func (c *Config) HideEdgeBorders(b BorderType) {
	c.raw("hide_edge_borders " + string(b))
}
