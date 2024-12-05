package access

type Role string

const (
	RoleViewer Role = "Viewer"
	RoleEditor Role = "Editor"
	RoleOwner  Role = "Owner"
)

func (r Role) IsValid() bool {
	return r == RoleViewer || r == RoleEditor || r == RoleOwner
}

func (r Role) GreaterThan(other Role) bool {
	if !r.IsValid() {
		return false
	}

	return roleRanks[r] > roleRanks[other]
}

var roleRanks = map[Role]int{
	RoleViewer: 0,
	RoleEditor: 1,
	RoleOwner:  2,
}

type CurriculumPrivileges struct {
	CurriculumId int64 `json:"-" db:"curriculum_id"`
	Role         Role  `json:"role" db:"role"`
}
