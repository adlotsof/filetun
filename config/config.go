package config

// import "github.com/alecthomas/kong"

var CLI struct {
	    BackendType string `name:"backend" help:"backend to use" enum:"file" default:"file"`
		OwnName string `name:"own_name" help:"name of the tun device" required:""`
		OwnCidr string `name:"own_cidr" help:"cidr of the tun device" required:""`
		Output string `help:"path to the output file" type:path required:""`
		PeersCidr string `name:peer_cidr help:"cidrs of subnets to route to" required:""`
		Input string `help:"path to the peers output file" type:path required:""`
	}
