package api

import (
	"net/http"
	"strconv"

	"github.com/pressly/chi"
)

type Role string
type RoleGrant int
type Authorizer func(string) bool

const (
	OwnerRole  = "owner"
	EditorRole = "editor"
	ViewerRole = "viewer"
)

const (
	CanAssignRoles = iota
	CanRevokeRoles
	CanUpdateRoles
	CanViewProjectRoles
	CanUpdateProject
	CanDeleteProject
	CanViewProject
	CanCreateLocales
	CanUpdateLocales
	CanDeleteLocales
	CanViewLocales
)

var permissions = map[Role][]RoleGrant{
	OwnerRole: []RoleGrant{
		CanAssignRoles,
		CanRevokeRoles,
		CanUpdateRoles,
		CanViewProjectRoles,
		CanUpdateProject,
		CanDeleteProject,
		CanViewProject,
		CanCreateLocales,
		CanUpdateLocales,
		CanDeleteLocales,
		CanViewLocales,
	},
	EditorRole: []RoleGrant{
		CanViewProjectRoles,
		CanUpdateProject,
		CanViewProject,
		CanCreateLocales,
		CanUpdateLocales,
		CanDeleteLocales,
		CanViewLocales,
	},
	ViewerRole: []RoleGrant{
		CanViewProjectRoles,
		CanViewProject,
		CanViewLocales,
	},
}

func isAllowed(r Role, a RoleGrant) bool {
	actions, ok := permissions[r]
	if !ok {
		return false
	}
	for _, currentAction := range actions {
		if currentAction == a {
			return true
		}
	}
	return false
}

func getProjectUserRole(userID, projID int) (string, error) {
	users, err := store.GetProjectUserRoles(projID)
	if err != nil {
		return "", err
	}
	for _, u := range users {
		if u.UserID == userID {
			return u.Role, nil
		}
	}
	return "", ErrNotFound
}

func mustAuthorize(action RoleGrant, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
		if err != nil {
			handleError(w, ErrBadRequest)
			return
		}
		requesterID, err := getUserIDFromContext(r.Context())
		if err != nil {
			handleError(w, ErrBadRequest)
			return
		}
		requesterRole, err := getProjectUserRole(requesterID, projectID)
		if err != nil {
			handleError(w, err)
			return
		}
		if !isAllowed(Role(requesterRole), action) {
			handleError(w, ErrForbiden)
			return
		}
		next.ServeHTTP(w, r)
	}
}
