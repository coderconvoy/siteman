package usr

import (
	"testing"
)

func t_permission(t *testing.T, u Usr, p string, exp int) {
	r := u.Permission(p)
	if r != exp {
		t.Errorf("Permission test: %s,%s. exp %d, got %d", u.Username, p, exp, r)
	}
}

func TestPermission(t *testing.T) {
	u := Usr{
		Username: "NoLog",
		Paths: map[string]int{
			"/logs": CAN_READ,
			"":      NO_READ,
		},
	}

	t_permission(t, u, "/logstik", CAN_EDIT)
	t_permission(t, u, "/logstik/poo", CAN_EDIT)
	t_permission(t, u, "/logs", CAN_READ)
	t_permission(t, u, "/logs/hello", CAN_READ)
	t_permission(t, u, "logstik", CAN_EDIT)
	t_permission(t, u, "logstik/poo", CAN_EDIT)
	t_permission(t, u, "logs", CAN_READ)
	t_permission(t, u, "logs/hello", CAN_READ)

}
